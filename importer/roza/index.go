package roza

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/queries"
)

type IdxFile struct {
	Name         string
	SHA1         string
	Size         int64
	LastModified time.Time
}

type IdxDirectory struct {
	Name     string
	Children *treemap.Map
	Files    []*IdxFile
}

func NewIdxDirectory(name string) *IdxDirectory {
	return &IdxDirectory{
		Name:     name,
		Children: treemap.NewWithStringComparator(),
		Files:    make([]*IdxFile, 0),
	}
}

func (d *IdxDirectory) printRec(prefix string) {
	fmt.Printf("%s%s\n", prefix, d.Name)
	idntPreifx := prefix + "\t"

	d.Children.Each(func(k interface{}, v interface{}) {
		v.(*IdxDirectory).printRec(idntPreifx)
	})
}

type RozaIndex struct {
	Roots     map[string]*IdxDirectory
	FileCount int
	DirCount  int
}

func (idx *RozaIndex) Load(db *sql.DB) error {
	idx.Roots = make(map[string]*IdxDirectory)

	rows, err := queries.Raw(db, "select path, sha1, size, last_modified from roza_index").Query()
	if err != nil {
		return errors.Wrap(err, "Load roza_index table")
	}
	defer rows.Close()

	idx.FileCount = 0
	for rows.Next() {
		var path string
		var sha1b []byte
		var size int64
		var lm time.Time
		if err := rows.Scan(&path, &sha1b, &size, &lm); err != nil {
			return errors.Wrap(err, "rows.Scan")
		}
		idx.FileCount++

		idx.insertFile(path, sha1b, size, lm)
	}

	if err := rows.Err(); err != nil {
		return errors.Wrap(err, "rows.Err")
	}

	return nil
}

func (idx *RozaIndex) GetDir(path string) *IdxDirectory {
	s := strings.Split(path, "/")[1:]

	d, ok := idx.Roots[s[0]]
	if !ok || len(s) == 1 {
		return d
	}

	for i := 1; i < len(s); i++ {
		if x, ok := d.Children.Get(s[i]); ok {
			d = x.(*IdxDirectory)
		} else {
			return nil
		}
	}

	return d
}

func (idx *RozaIndex) insertFile(path string, sha1b []byte, size int64, lm time.Time) {
	li := strings.LastIndex(path, "/")

	d := idx.GetDir(path[:li])
	if d == nil {
		d = idx.mkDirAll(path[:li])
	}

	f := &IdxFile{
		Name:         path[li+1:],
		SHA1:         hex.EncodeToString(sha1b),
		Size:         size,
		LastModified: lm,
	}
	d.Files = append(d.Files, f)
}

func (idx *RozaIndex) mkDirAll(path string) *IdxDirectory {
	s := strings.Split(path, "/")[1:]

	d, ok := idx.Roots[s[0]]
	if !ok {
		d = NewIdxDirectory(s[0])
		idx.Roots[d.Name] = d
	}

	if len(s) == 1 {
		return d
	}

	var x interface{}
	for i := 1; i < len(s); i++ {
		if x, ok = d.Children.Get(s[i]); !ok {
			x = NewIdxDirectory(s[i])
			d.Children.Put(x.(*IdxDirectory).Name, x)
		}
		d = x.(*IdxDirectory)
	}

	return d
}

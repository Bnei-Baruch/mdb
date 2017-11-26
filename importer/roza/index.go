package roza

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"encoding/json"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/queries"
)

type IdxFile struct {
	Name         string
	SHA1         string
	Size         int64
	LastModified time.Time
	Directory    *IdxDirectory
}

func (f *IdxFile) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name string `json:"name"`
		Path string `json:"path"`
	}{
		Name: f.Name,
		Path: f.Directory.path()[len("/vfs/archive/Archive"):],
	})
}

type IdxDirectory struct {
	Name     string
	Parent   *IdxDirectory
	Children *treemap.Map
	Files    []*IdxFile
}

func NewIdxDirectory(name string, parent *IdxDirectory) *IdxDirectory {
	return &IdxDirectory{
		Name:     name,
		Parent:   parent,
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

func (d *IdxDirectory) path() string {
	path := d.Name
	x := d
	for x.Parent != nil {
		path = fmt.Sprintf("%s/%s", x.Parent.Name, path)
		x = x.Parent
	}
	return fmt.Sprintf("/%s", path)
}

type RozaIndex struct {
	Roots     map[string]*IdxDirectory
	FileCount int
}

func (idx *RozaIndex) Load(db *sql.DB) error {
	idx.Roots = make(map[string]*IdxDirectory)
	idx.FileCount = 0

	rows, err := queries.Raw(db, "select path, sha1, size, last_modified from roza_index").Query()
	if err != nil {
		return errors.Wrap(err, "Load roza_index table")
	}
	defer rows.Close()

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
		Directory:    d,
	}
	d.Files = append(d.Files, f)
}

func (idx *RozaIndex) mkDirAll(path string) *IdxDirectory {
	s := strings.Split(path, "/")[1:]

	d, ok := idx.Roots[s[0]]
	if !ok {
		d = NewIdxDirectory(s[0], nil)
		idx.Roots[d.Name] = d
	}

	if len(s) == 1 {
		return d
	}

	var x interface{}
	for i := 1; i < len(s); i++ {
		if x, ok = d.Children.Get(s[i]); !ok {
			x = NewIdxDirectory(s[i], d)
			d.Children.Put(x.(*IdxDirectory).Name, x)
		}
		d = x.(*IdxDirectory)
	}

	return d
}

func (idx *RozaIndex) Sha1Map() map[string][]*IdxFile {

	// all roots
	//s := make([]*IdxDirectory, len(idx.Roots))
	//i := 0
	//for k := range idx.Roots {
	//	s[i] = idx.Roots[k]
	//	i++
	//}

	// root is ____beavoda
	s := []*IdxDirectory{idx.GetDir("/vfs/archive/Archive/____beavoda")}

	sMap := make(map[string][]*IdxFile, 600000)
	var x *IdxDirectory
	for len(s) > 0 {
		x, s = s[0], s[1:]

		for i := range x.Files {
			f := x.Files[i]
			k := f.SHA1
			v, ok := sMap[k]
			if !ok {
				v = make([]*IdxFile, 0)
			}
			sMap[k] = append(v, f)
		}

		if x.Children.Empty() {
			continue
		}

		values := make([]*IdxDirectory, x.Children.Size())
		it := x.Children.Iterator()
		for i := 0; it.Next(); i++ {
			values[i] = it.Value().(*IdxDirectory)
		}
		s = append(values, s...)
	}

	return sMap
}

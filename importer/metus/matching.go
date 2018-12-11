package metus

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/hashicorp/go-immutable-radix"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/models"
)

type FilenameMatcher struct {
	tree *iradix.Tree
}

func (fm *FilenameMatcher) Load() error {
	log.Info("FilenameMatcher: Loading files from mdb")
	files, err := models.Files(mdb, qm.Select("id, name, content_unit_id")).All()
	if err != nil {
		return errors.Wrap(err, "load mdb files")
	}

	log.Info("FilenameMatcher: creating radix trie")
	t := iradix.New()
	tx := t.Txn()
	for i := range files {
		tx.Insert([]byte(fm.cleanFilename(files[i].Name)), files[i].ContentUnitID.Int64)
	}
	fm.tree = tx.Commit()

	return nil
}

func (fm *FilenameMatcher) Match(filename string) int64 {
	// clean hbru
	fName := filename
	if strings.HasPrefix(fName, "hbru") {
		fName = fmt.Sprintf("heb%s", fName[4:])
	}

	fName = fm.cleanFilename(fName)
	bFName := []byte(fName)

	if v, ok := fm.tree.Get(bFName); ok {
		return v.(int64)
	}

	fm.tree.Root().Walk(func (k []byte, v interface{}) bool {
		return true
	})

	return -1
}

func (fm *FilenameMatcher) cleanFilename(filename string) string {
	s := strings.Split(filename, ".")
	return strings.ToLower(strings.Join(s[:len(s)-1], "."))
}

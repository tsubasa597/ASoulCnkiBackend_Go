package cache

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/tidwall/buntdb"
	"github.com/tsubasa597/ASoulCnkiBackend/conf"
)

type BuntDB struct {
	db       *buntdb.DB
	mutex    *sync.Mutex
	fileName string
}

var _ Cacher = (*BuntDB)(nil)

func (b BuntDB) Get(v string) (val string, err error) {
	if err = b.db.View(func(tx *buntdb.Tx) error {
		val, err = tx.Get(v)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return
	}
	return
}

func (b BuntDB) Set(key, value string) error {
	return b.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, value, nil)
		return err
	})
}

func NewBuntDB(path string) (*BuntDB, error) {
	b := &BuntDB{
		mutex:    &sync.Mutex{},
		fileName: path,
	}

	os.Mkdir(conf.CacheFilePath, os.ModePerm)
	file, err := os.OpenFile(conf.CacheFilePath+b.fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	b.db, err = buntdb.Open(":memory:")
	if err != nil {
		return nil, err
	}

	if info, _ := file.Stat(); info.Size() != 0 {
		b.db.Load(file)
	}

	return b, nil
}

func (b BuntDB) Save() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if err := os.Remove(conf.CacheFilePath + b.fileName); err != nil {
		return err
	}

	file, err := os.OpenFile(conf.CacheFilePath+b.fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	return b.db.Save(file)
}

func (b BuntDB) Increment(id string, hashSet map[int64]struct{}) error {
	b.db.Update(func(tx *buntdb.Tx) error {
		for k := range hashSet {
			if val, err := tx.Get(fmt.Sprint(k)); err == nil {
				if strings.Contains(val, id) {
					continue
				}

				tx.Set(fmt.Sprint(k), val+","+id, nil)
				continue
			}
			tx.Set(fmt.Sprint(k), id, nil)
		}
		tx.Set("LastCommentID", id, nil)

		return nil
	})

	return nil
}

func (b BuntDB) Stop() {
	b.db.Close()
}

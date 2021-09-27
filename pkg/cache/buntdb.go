package cache

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/tidwall/buntdb"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/setting"
)

type BuntDB struct {
	db       *buntdb.DB
	mutex    *sync.Mutex
	fileName string
}

var _ Cacher = (*BuntDB)(nil)

func (b BuntDB) Get(_, v string) (val string, err error) {
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

func NewBuntDB(path string) (*BuntDB, error) {
	b := &BuntDB{
		mutex:    &sync.Mutex{},
		fileName: path,
	}

	os.Mkdir(setting.CacheFilePath, os.ModePerm)
	file, err := os.OpenFile(setting.CacheFilePath+b.fileName, os.O_RDWR|os.O_CREATE, 0755)
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

	if err := os.Remove(setting.CacheFilePath + b.fileName); err != nil {
		return err
	}

	file, err := os.OpenFile(setting.CacheFilePath+b.fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	return b.db.Save(file)
}

func (b BuntDB) Increment(_ string, field string, val interface{}) error {
	b.db.Update(func(tx *buntdb.Tx) error {
		switch v := val.(type) {
		case map[int64]struct{}:
			for k := range v {
				if val, err := tx.Get(fmt.Sprint(k)); err == nil {
					if strings.Contains(val, field) {
						continue
					}

					tx.Set(fmt.Sprint(k), val+","+field, nil)
					continue
				}
				tx.Set(fmt.Sprint(k), field, nil)
			}
			tx.Set("LastCommentID", field, nil)
		case string:
			tx.Set(field, v, nil)
			tx.Set("LastCommentID", field, nil)
		}

		return nil
	})

	return nil
}

func (b BuntDB) Stop() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.db.Close()
}

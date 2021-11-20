package cache

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/tidwall/buntdb"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/config"
)

// BuntDB 缓存
type BuntDB struct {
	db       *buntdb.DB
	mutex    *sync.Mutex
	fileName string
}

var _ Cacher = (*BuntDB)(nil)

// Get 获取缓存值
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

// NewBuntDB 实例化 BuntDB
func NewBuntDB(path string) (*BuntDB, error) {
	b := &BuntDB{
		mutex:    &sync.Mutex{},
		fileName: path,
	}

	if err := os.Mkdir(config.CacheFilePath, os.ModePerm); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(config.CacheFilePath+b.fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	b.db, err = buntdb.Open(":memory:")
	if err != nil {
		return nil, err
	}

	if info, _ := file.Stat(); info.Size() != 0 {
		if err := b.db.Load(file); err != nil {
			return nil, err
		}
	}

	return b, nil
}

// Save 持久化
func (b BuntDB) Save() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if err := os.Remove(config.CacheFilePath + b.fileName); err != nil {
		return err
	}

	file, err := os.OpenFile(config.CacheFilePath+b.fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	return b.db.Save(file)
}

// Increment 添加数据
func (b BuntDB) Increment(_ string, field string, val interface{}) error {
	return b.db.Update(func(tx *buntdb.Tx) error {
		switch v := val.(type) {
		case map[int64]struct{}:
			for k := range v {
				if val, err := tx.Get(fmt.Sprint(k)); err == nil {
					if strings.Contains(val, field) {
						continue
					}

					if _, _, err := tx.Set(fmt.Sprint(k), val+","+field, nil); err != nil {
						return err
					}
					continue
				}

				if _, _, err := tx.Set(fmt.Sprint(k), field, nil); err != nil {
					return err
				}
			}
		case string:
			if _, _, err := tx.Set(field, v, nil); err != nil {
				return err
			}
		}

		return nil
	})
}

// Stop 停止
func (b BuntDB) Stop() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.db.Close()
}

package cache

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/config"
)

// LevelDB 缓存
type LevelDB struct {
	db    *leveldb.DB
	mutex *sync.RWMutex
}

var _ Cacher = (*LevelDB)(nil)

// NewLevelDB 实例化 LevelDB
func NewLevelDB(path string) (*LevelDB, error) {
	db, err := leveldb.OpenFile(config.CacheFilePath+path, nil)
	if err != nil {
		return nil, err
	}

	return &LevelDB{
		mutex: &sync.RWMutex{},
		db:    db,
	}, nil
}

// Get 获取缓存值
func (l LevelDB) Get(_, key string) (string, error) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	val, err := l.db.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	return string(val), nil
}

// Increment 添加数据
func (l LevelDB) Increment(_ string, field string, val interface{}) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	switch v := val.(type) {
	case map[int64]struct{}:
		bID := []byte(field)
		buffer := bytes.Buffer{}
		key := bytes.Buffer{}

		for k := range v {
			key.Reset()
			key.WriteString(fmt.Sprint(k))

			if val, err := l.db.Get(key.Bytes(), nil); err == nil {
				buffer.Reset()
				buffer.Write(val)

				if bytes.Contains(val, bID) {
					continue
				}

				buffer.WriteString(",")
				buffer.Write(bID)

				if err := l.db.Put(key.Bytes(), buffer.Bytes(), nil); err != nil {
					return err
				}
				continue
			}
			if err := l.db.Put(key.Bytes(), bID, nil); err != nil {
				return err
			}
		}
	case string:
		if err := l.db.Put([]byte(field), []byte(v), nil); err != nil {
			return err
		}
	}

	return nil
}

// Save 持久化
func (l LevelDB) Save() error {
	return nil
}

// Stop 停止
func (l LevelDB) Stop() error {
	return nil
}

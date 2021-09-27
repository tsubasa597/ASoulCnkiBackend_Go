package cache

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/setting"
)

type LevelDB struct {
	db    *leveldb.DB
	mutex *sync.RWMutex
}

var _ Cacher = (*LevelDB)(nil)

func NewLevelDB(path string) (*LevelDB, error) {
	db, err := leveldb.OpenFile(setting.CacheFilePath+path, nil)
	if err != nil {
		return nil, err
	}

	return &LevelDB{
		mutex: &sync.RWMutex{},
		db:    db,
	}, nil
}

func (l LevelDB) Get(_, key string) (string, error) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	val, err := l.db.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	return string(val), nil
}

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
				l.db.Put(key.Bytes(), buffer.Bytes(), nil)
				continue
			}
			l.db.Put(key.Bytes(), bID, nil)
		}
		l.db.Put([]byte("LastCommentID"), bID, nil)
	case string:
		l.db.Put([]byte(field), []byte(v), nil)
		l.db.Put([]byte("LastCommentID"), []byte(field), nil)
	}

	return nil
}

func (l LevelDB) Save() error {
	return nil
}

func (l LevelDB) Stop() {

}

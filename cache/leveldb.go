package cache

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tsubasa597/ASoulCnkiBackend/conf"
)

type LevelDB struct {
	db    *leveldb.DB
	mutex *sync.RWMutex
}

var _ Cacher = (*LevelDB)(nil)

func NewLevelDB(path string) (*LevelDB, error) {
	db, err := leveldb.OpenFile(conf.CacheFilePath+path, nil)
	if err != nil {
		return nil, err
	}

	return &LevelDB{
		mutex: &sync.RWMutex{},
		db:    db,
	}, nil
}

func (l LevelDB) Get(key string) (string, error) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	val, err := l.db.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func (l LevelDB) Set(key, value string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	return l.db.Put([]byte(key), []byte(value), nil)
}

func (l LevelDB) Increment(id string, hashSet map[int64]struct{}) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	bID := []byte(id)
	buffer := bytes.Buffer{}
	key := bytes.Buffer{}

	for k := range hashSet {
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
	return nil
}

func (l LevelDB) Save() error {
	return nil
}

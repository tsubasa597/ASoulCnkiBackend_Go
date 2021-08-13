package cache

import (
	"fmt"
	"os"
	"sync"

	"github.com/tidwall/buntdb"
	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/db/entry"
)

type Comment struct {
	db    *buntdb.DB
	mutex *sync.Mutex
}

var _ Cacher = (*Comment)(nil)

func (c Comment) Get(v interface{}) (interface{}, error) {
	var data string

	if s, ok := v.(string); ok {
		if err := c.db.View(func(tx *buntdb.Tx) error {
			var err error

			data, err = tx.Get(s)
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("type error")
	}

	return data, nil
}

func (c Comment) Set(key, value interface{}) error {
	k, ok1 := key.(string)
	v, ok2 := value.(string)

	if ok1 && ok2 {
		c.db.Update(func(tx *buntdb.Tx) error {
			tx.Set(k, v, nil)
			return nil
		})
	} else {
		return fmt.Errorf("type error")
	}
	return nil
}

func (c Comment) Update(key interface{}, value interface{}) error {
	k, ok1 := key.(string)
	v, ok2 := value.(string)

	if ok1 && ok2 {
		c.db.Update(func(tx *buntdb.Tx) error {
			if val, err := tx.Get(k); err == nil {
				tx.Set(k, val+","+v, nil)
				return nil
			}
			tx.Set(k, v, nil)
			return nil
		})
	} else {
		return fmt.Errorf("type error")
	}
	return nil
}

func NewComment() (*Comment, error) {
	c := &Comment{
		mutex: &sync.Mutex{},
	}

	file, err := os.OpenFile(conf.CacheFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	c.db, err = buntdb.Open(":memory:")
	if err != nil {
		return nil, err
	}

	if info, _ := file.Stat(); info.Size() != 0 {
		c.db.Load(file)
	}

	return c, nil
}

func (c Comment) Save() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := os.Remove(conf.CacheFile); err != nil {
		return err
	}

	file, err := os.OpenFile(conf.CacheFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	return c.db.Save(file)
}

func (c Comment) Increment(comm entry.Comment, hashSet map[int64]struct{}) error {
	err := c.db.Update(func(tx *buntdb.Tx) error {
		for k := range hashSet {
			if val, err := tx.Get(fmt.Sprint(k)); err == nil {
				tx.Set(fmt.Sprint(k), val+","+fmt.Sprint(comm.ID), nil)
				continue
			}
			tx.Set(fmt.Sprint(k), fmt.Sprint(comm.ID), nil)
		}
		tx.Set("LastCommentID", fmt.Sprint(comm.ID), nil)

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

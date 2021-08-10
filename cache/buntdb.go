package cache

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/tidwall/buntdb"
	"github.com/tsubasa597/ASoulCnkiBackend/conf"
	"github.com/tsubasa597/ASoulCnkiBackend/db"
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

func (c Comment) Init(key, value interface{}) error {
	if comments, ok := value.(map[int64]struct{}); ok {
		for k := range comments {
			data, id := fmt.Sprint(k), fmt.Sprint(key.(uint64))
			if ids, err := c.Get(data); err == nil {
				if strings.Contains(ids.(string), id) {
					continue
				}
				c.Set(data, ids.(string)+","+id)
				continue
			}
			c.Set(data, id)
		}
	} else {
		return fmt.Errorf("type error")
	}
	return nil
}

func (c Comment) Increment(db_ db.DB, f func(string) map[int64]struct{}) error {
	val, err := c.Get("LastCommentID")
	if err != nil {
		val = "0"
	}

	comms, err := db_.Find(&entry.Comment{}, db.Param{
		Query: "id > ?",
		Args:  []interface{}{val},
		Order: "id",
	})
	if err != nil {
		return err
	}

	for _, v := range *comms.(*[]entry.Comment) {
		for k := range f(v.Comment) {
			if val, err := c.Get(fmt.Sprint(v.ID)); err == nil {
				c.Set(fmt.Sprint(k), val.(string)+","+fmt.Sprint(v.ID))
				continue
			}
			c.Set(fmt.Sprint(k), fmt.Sprint(v.ID))
		}
		c.Set("LastCommentID", fmt.Sprint(v.ID))
	}

	return c.Save()
}
package cache

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/buntdb"
)

type Comment struct {
	IsInit bool
	db     *buntdb.DB
	file   *os.File
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

func NewComment() (*Comment, error) {
	c := &Comment{}

	file, err := os.OpenFile("./cache.dat", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}

	c.file = file

	c.db, err = buntdb.Open(":memory:")
	if err != nil {
		return nil, err
	}

	if info, _ := file.Stat(); info.Size() != 0 {
		c.IsInit = true
		c.db.Load(file)
	}

	return c, nil
}

func (c Comment) Save() error {
	if c.file != nil {
		c.db.Save(c.file)
		return nil
	}
	return fmt.Errorf("file not fuond")
}

func (c Comment) Init(key, value interface{}) error {
	if c.IsInit {
		return fmt.Errorf("already init")
	}

	if comments, ok := value.(map[int64]struct{}); ok {
		for k := range comments {
			data, id := strconv.Itoa(int(k)), strconv.Itoa(int(key.(uint64)))
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

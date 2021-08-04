package cache

import (
	"fmt"
	"os"

	"github.com/tidwall/buntdb"
)

type Comment struct {
	db   *buntdb.DB
	file *os.File
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

func NewComment(f func(Cacher)) (*Comment, error) {
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
		c.db.Load(file)
	} else {
		c.Load(f)
		c.db.Save(file)
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

func (c Comment) Load(f func(Cacher)) {
	f(c)
}

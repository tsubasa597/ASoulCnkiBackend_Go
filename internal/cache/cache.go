package cache

import (
	"fmt"
	"strings"

	"github.com/tsubasa597/ASoulCnkiBackend/internal/dao"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/check"
	"github.com/tsubasa597/BILIBILI-HELPER/info"
	"gorm.io/gorm"
)

const (
	CheckKey, ContentKey string = "check", "content"
	// 初始化缓存时一次读取数据量
	initBatch int = 10000
)

// Cacher 缓存接口
type Cacher interface {
	Get(string, string) (string, error)
	Save() error
	Increment(string, string, interface{}) error
	Stop() error
}

// Cache 缓存实例
type Cache struct {
	Check   Cacher
	Content Cacher
}

var (
	cache Cache
)

// Stop 停止
func (c Cache) Stop() string {
	var sb strings.Builder

	if err := c.Check.Save(); err != nil {
		sb.WriteString(err.Error())
	}

	if err := c.Content.Save(); err != nil {
		sb.WriteString(err.Error())
	}

	if err := c.Check.Stop(); err != nil {
		sb.WriteString(err.Error())
	}

	if err := c.Content.Stop(); err != nil {
		sb.WriteString(err.Error())
	}

	return sb.String()
}

// Setup 初始化
func Setup() error {
	check, err := NewRedis()
	if err != nil {
		return err
	}

	content, err := NewRedis()
	if err != nil {
		return err
	}

	cache = Cache{
		Check:   *check,
		Content: *content,
	}

	contents := make([]dao.CommentCache, 0, initBatch)
	return dao.GetContent(initBatch, contents, func(tx *gorm.DB, batch int) error {
		for _, content := range contents {
			if err := cache.Store(content); err != nil {
				return err
			}
		}
		return nil
	})
}

func (c Cache) Store(comments interface{}) error {
	switch comments := comments.(type) {
	case []dao.CommentCache:
		for _, comment := range comments {
			if err := c.Check.Increment(CheckKey, fmt.Sprint(comment.Rpid),
				check.HashSet(comment.Content)); err != nil {
				return err
			}

			if err := c.Content.Increment(ContentKey, fmt.Sprint(comment.Rpid),
				comment.Content); err != nil {
				return err
			}
		}
	case []info.Comment:
		for _, comment := range comments {
			if err := c.Check.Increment(CheckKey, fmt.Sprint(comment.Rpid),
				check.HashSet(comment.Content)); err != nil {
				return err
			}

			if err := c.Content.Increment(ContentKey, fmt.Sprint(comment.Rpid),
				comment.Content); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("Error Type")
	}

	return nil
}

// GetCache 获取缓存实例
func GetCache() Cache {
	return cache
}

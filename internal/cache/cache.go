package cache

import (
	"fmt"
	"strings"

	"github.com/tsubasa597/ASoulCnkiBackend/internal/dao"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/check"
	"github.com/tsubasa597/BILIBILI-HELPER/info"
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

const (
	CheckKey, ContentKey string = "check", "content"
	// 初始化缓存时一次读取数据量
	_initBatch int = 100000
)

var (
	_cache Cache
)

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

	_cache = Cache{
		Check:   *check,
		Content: *content,
	}

	return dao.GetContent(_initBatch, func(contents []dao.CommentCache) error {
		if err := _cache.Store(contents); err != nil {
			return err
		}
		return nil
	})
}

// Store 保存缓存
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

// Stop 停止
func (c Cache) Stop() string {
	var sb strings.Builder

	if err := c.Check.Save(); err != nil {
		sb.WriteString(err.Error() + "\n")
	}

	if err := c.Content.Save(); err != nil {
		sb.WriteString(err.Error() + "\n")
	}

	if err := c.Check.Stop(); err != nil {
		sb.WriteString(err.Error() + "\n")
	}

	if err := c.Content.Stop(); err != nil {
		sb.WriteString(err.Error())
	}

	return sb.String()
}

// GetInstance 获取缓存实例
func GetInstance() Cache {
	return _cache
}

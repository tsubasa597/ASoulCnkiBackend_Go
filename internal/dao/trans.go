package dao

import (
	"sort"
	"sync"

	"github.com/tsubasa597/ASoulCnkiBackend/internal/model/entity"
	"github.com/tsubasa597/BILIBILI-HELPER/info"
)

var (
	dynamicPool *sync.Pool = &sync.Pool{
		New: func() interface{} {
			return &entity.Dynamic{}
		},
	}
	commentPool *sync.Pool = &sync.Pool{
		New: func() interface{} {
			return &entity.Comment{}
		},
	}
)

// TransDynamic 将爬取数据转义为数据库结构
func TransDynamic(infos []info.Dynamic) entity.Dynamics {
	dynamics := make(entity.Dynamics, 0, len(infos))
	for _, inf := range infos {
		dynamic := dynamicPool.Get().(*entity.Dynamic)
		dynamic.RID = inf.RID
		dynamic.Type = uint8(inf.Type)
		dynamic.Time = inf.Time
		dynamic.Content = inf.Content
		dynamic.Card = inf.Card
		dynamic.Name = inf.Name
		dynamic.UID = inf.UID

		dynamics = append(dynamics, *dynamic)
		dynamicPool.Put(dynamic)
	}

	// 为更新 user 表中字段排序
	sort.Sort(dynamics)
	return dynamics
}

// TransComment 将爬取数据转义为数据库结构
func TransComment(infos []info.Comment) entity.Comments {
	comments := make(entity.Comments, 0, len(infos))
	for _, inf := range infos {
		comm := commentPool.Get().(*entity.Comment)
		comm.Name = inf.Name
		comm.Time = inf.Time
		comm.DynamicUID = inf.DynamicUID
		comm.UID = inf.UID
		comm.Rpid = inf.Rpid
		comm.LikeNum = uint32(inf.LikeNum)
		comm.Content = inf.Content
		comm.DynamicID = uint64(inf.RID)

		comments = append(comments, *comm)
		commentPool.Put(comm)
	}

	return comments
}

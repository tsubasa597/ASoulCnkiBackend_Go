package listen

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	rpcv1 "github.com/tsubasa597/ASoulCnkiBackend/api/grpc/v1"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/cache"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/dao"
	"github.com/tsubasa597/ASoulCnkiBackend/internal/model/entity"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/config"
	"github.com/tsubasa597/BILIBILI-HELPER/info"
	"github.com/tsubasa597/BILIBILI-HELPER/task"
)

type Listen struct {
	Ctx     context.Context
	cancel  context.CancelFunc
	dynamic task.Corn
	comment task.Corn
	rpc     bool
	log     *logrus.Logger
}

var (
	_listen *Listen
)

func Setup() error {
	ctx, cancel := context.WithCancel(context.Background())
	lis := &Listen{
		Ctx:    ctx,
		cancel: cancel,
	}

	if !config.Enable {
		_listen = lis
		return nil
	}
	lis.dynamic = task.New(ctx)
	lis.comment = task.New(ctx)

	log, err := config.NewLogFile("/listen")
	if err != nil {
		return err
	}
	lis.log = log

	if config.RPCEnable {
		if err := rpcv1.Setup(); err != nil {
			return err
		}

		lis.rpc = true
	}

	if users, err := dao.Get(&entity.User{}); err != nil {
		return err
	} else {
		for _, user := range *users.(*[]entity.User) {
			lis.addDynamic(user.UID, user.LastDynamicTime)
		}
	}

	if dynamics, err := dao.Get(&entity.Dynamic{}); err == nil {
		for _, dynamic := range *dynamics.(*[]entity.Dynamic) {
			lis.addComment(dynamic.Time, dynamic.RID, info.Type(dynamic.Type))
		}
	}

	go lis.SaveDynamic()
	go lis.SaveComment()

	lis.comment.Start()
	lis.dynamic.Start()

	_listen = lis
	return nil
}

func (l Listen) SaveDynamic() {
	for data := range l.dynamic.Ch {
		if len(data.([]info.Dynamic)) == 0 {
			continue
		}

		dynamics := dao.TransDynamic(data.([]info.Dynamic))
		if err := dao.Add(dynamics, info.MaxPs); err != nil {
			l.log.Error(err)
			continue
		}

		for _, dynamic := range dynamics {
			l.addComment(dynamic.Time, dynamic.RID, info.Type(dynamic.Type))
		}
	}
}

func (l Listen) SaveComment() {
	for data := range l.comment.Ch {
		if len(data.([]info.Comment)) == 0 {
			continue
		}

		comments := dao.TransComment(data.([]info.Comment))
		// batchSize 为 1
		// 保证事务处理完成，并提交
		if err := dao.Add(comments, 1); err != nil {
			l.log.Error(err)
			continue
		}

		if err := cache.GetInstance().Store(data.([]info.Comment)); err != nil {
			l.log.Error(err)
		}
	}
}

func (l Listen) Stop() error {
	l.cancel()

	if l.rpc {
		return rpcv1.Stop()
	}
	return nil
}

func (l Listen) addDynamic(uid, t int64) {
	if l.rpc {
		l.dynamic.Add(uid, rpcv1.NewDynamic(l.Ctx, uid, t, l.log))
	} else {
		l.dynamic.Add(uid, task.NewDynamic(uid, t, time.Duration(config.DynamicDuration)))
	}
}

func (l Listen) addComment(t, rid int64, typ info.Type) {
	if l.rpc {
		l.comment.Add(rid, rpcv1.NewComment(l.Ctx, t, rid, info.Type(typ), l.log))
	} else {
		l.comment.Add(rid, task.NewComment(rid, t, time.Duration(1 /* timeCell 间隔时间为 1s */),
			info.Type(typ), l.log))
	}
}

func GetInstance() *Listen {
	return _listen
}

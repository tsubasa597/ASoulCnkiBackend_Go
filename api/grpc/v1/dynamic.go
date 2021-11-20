package rpcv1

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/config"
	"github.com/tsubasa597/BILIBILI-HELPER/info"
	"github.com/tsubasa597/BILIBILI-HELPER/rpc/service"
	"github.com/tsubasa597/BILIBILI-HELPER/state"
	"github.com/tsubasa597/BILIBILI-HELPER/task"
	"google.golang.org/grpc"
)

type RPCDynamic struct {
	UID      int64
	Time     int64
	Ctx      context.Context
	timeCell time.Duration
	log      *logrus.Logger
	state    state.State
}

var (
	_              task.Tasker = (*RPCDynamic)(nil)
	dynamicClient  service.DynamicClient
	dynamicReqPool = &sync.Pool{
		New: func() interface{} {
			return &service.AllDynamicRequest{
				BaseCommentRequest: &service.BaseDynamicRequest{},
			}
		},
	}
	dynamicRespPool = &sync.Pool{
		New: func() interface{} {
			return &info.Dynamic{}
		},
	}
)

func (r *RPCDynamic) Run(ch chan<- interface{}) {
	if r.state == state.Stop {
		return
	}

	req := dynamicReqPool.Get().(*service.AllDynamicRequest)
	defer dynamicReqPool.Put(req)
	req.Time = r.Time
	req.BaseCommentRequest.UID = r.UID

	stream, err := dynamicClient.GetAll(r.Ctx, req, grpc.EmptyCallOption{})
	if err != nil {
		r.log.Error(err)
	}

	dynamics := make([]info.Dynamic, 0)
	for {
		resp, err := stream.Recv()
		if err != nil {
			if err != io.EOF {
				r.log.Error(err)
			}

			if len(dynamics) != 0 {
				ch <- dynamics
			}
			return
		}

		dynamic := dynamicRespPool.Get().(*info.Dynamic)
		dynamic.Name = resp.Name
		dynamic.Time = resp.Time
		dynamic.Card = resp.Card
		dynamic.Content = resp.Content
		dynamic.RID = resp.RID
		dynamic.Type = info.Type(resp.Type)
		dynamic.UID = resp.UID

		dynamics = append(dynamics, *dynamic)
		dynamicRespPool.Put(dynamic)
	}
}

func (r RPCDynamic) State() state.State {
	return r.state
}

// Next 下次运行时间
func (r RPCDynamic) Next(t time.Time) time.Time {
	return t.Add(time.Minute * r.timeCell)
}

func NewDynamic(ctx context.Context, uid, t int64, log *logrus.Logger) *RPCDynamic {
	return &RPCDynamic{
		UID:      uid,
		Time:     t,
		Ctx:      ctx,
		timeCell: time.Duration(config.DynamicDuration),
		log:      log,
		state:    state.Runing,
	}
}

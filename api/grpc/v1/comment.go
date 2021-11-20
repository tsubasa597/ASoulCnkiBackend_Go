package rpcv1

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/BILIBILI-HELPER/info"
	"github.com/tsubasa597/BILIBILI-HELPER/rpc/service"
	"github.com/tsubasa597/BILIBILI-HELPER/state"
	"github.com/tsubasa597/BILIBILI-HELPER/task"
	"google.golang.org/grpc"
)

type RPCComment struct {
	RID      int64
	Type     info.Type
	Ctx      context.Context
	Time     int64
	timeCell time.Duration
	log      *logrus.Logger
	state    state.State
}

var (
	_              task.Tasker = (*RPCComment)(nil)
	commentClient  service.CommentClient
	commentReqPool = &sync.Pool{
		New: func() interface{} {
			return &service.AllCommentRequest{
				BaseCommentRequest: &service.BaseCommentRequest{},
			}
		},
	}
	commentRespPool = &sync.Pool{
		New: func() interface{} {
			return &info.Comment{}
		},
	}
	CommentTimeCell time.Duration = 1 * time.Hour
)

func NewComment(ctx context.Context, t, rid int64, typ info.Type, log *logrus.Logger) *RPCComment {
	return &RPCComment{
		RID:      rid,
		Type:     typ,
		Ctx:      ctx,
		Time:     t,
		timeCell: CommentTimeCell,
		log:      log,
		state:    state.Runing,
	}
}

func (r *RPCComment) Run(ch chan<- interface{}) {
	if r.State() == state.Stop {
		return
	}

	req := commentReqPool.Get().(*service.AllCommentRequest)
	defer commentReqPool.Put(req)
	req.Time = 0
	req.BaseCommentRequest.RID = r.RID
	req.BaseCommentRequest.Type = int32(r.Type)

	stream, err := commentClient.GetAll(r.Ctx, req, grpc.EmptyCallOption{})
	if err != nil {
		r.log.Error(err)
	}

	comments := make([]info.Comment, 0)
	for {
		resp, err := stream.Recv()
		if err != nil {
			if err != io.EOF {
				r.log.Error(err)
			}

			if len(comments) != 0 {
				ch <- comments
			}
			return
		}

		comment := commentRespPool.Get().(*info.Comment)
		comment.Time = resp.Time
		comment.Name = resp.Name
		comment.DynamicUID = resp.DynamicUID
		comment.RID = resp.RID
		comment.UID = resp.UID
		comment.Rpid = resp.Rpid
		comment.LikeNum = uint32(resp.LikeNum)
		comment.Content = resp.Content

		comments = append(comments, *comment)
		commentRespPool.Put(comment)
	}
}

func (r RPCComment) State() state.State {
	return r.state
}

// Next 下次运行时间
func (r RPCComment) Next(t time.Time) time.Time {
	if time.Now().AddDate(0, 0, -7).Unix() > r.Time {
		return t.Add(time.Hour * 24 * 2)
	}

	return t.Add(time.Second * r.timeCell)
}

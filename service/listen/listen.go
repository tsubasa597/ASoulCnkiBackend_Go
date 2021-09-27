package listen

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/ASoulCnkiBackend/models/entity"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/cache"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/check"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/model"
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/setting"
	"github.com/tsubasa597/BILIBILI-HELPER/listen"
)

type Listen struct {
	comment *listen.Listen
	dynamic *listen.Listen
	Enable  bool
	log     *logrus.Entry
}

func (lis Listen) Stop() {
	lis.dynamic.Stop()
	lis.comment.Stop()
	lis.Enable = false
}

func (lis Listen) Add(user entity.User) {
	if !lis.Enable {
		return
	}

	ctx, ch, err := lis.dynamic.Add(user.UID, user.LastDynamicTime,
		time.Duration(setting.DynamicDuration)*time.Minute)
	if err != nil {
		lis.log.WithField("Func", "Listen.Add").Error(err)
		return
	}

	lis.log.WithField("Func", "Listen.Add").Info(fmt.Sprintf("Listen %d", user.UID))
	go lis.SaveDyanmic(ctx, user.ID, ch)
}

func New(log *logrus.Entry) *Listen {
	ctx := context.Background()
	comm := listen.New(ctx, listen.NewComment(ctx, setting.GoroutineNum, log), nil, log)
	dyna := listen.New(ctx, listen.NewDynamic(ctx, log), nil, log)

	listen := &Listen{
		Enable:  setting.Enable,
		comment: comm,
		dynamic: dyna,
		log:     log,
	}

	if !setting.Enable {
		return listen
	}

	users, err := model.Find(&entity.User{}, model.Param{
		Order: "id asc",
		Page:  -1,
	})
	if err != nil {
		return listen
	}

	for _, user := range *users.(*[]entity.User) {
		listen.Add(user)
	}

	val, err := cache.GetCache().Content.Get("content", "LastCommentID")
	if err != nil {
		val = "0"
	}

	comms, err := model.GetContent(val)
	if err != nil {
		log.WithField("Func", "model.GetContent").Error(err)
	}

	for _, comm := range comms {
		if err := cache.GetCache().Content.Increment("content", fmt.Sprint(comm.ID), comm.Content); err != nil {
			log.WithField("Func", "cache.Set").Error(err)
		}
		cache.GetCache().Content.Increment("content", "LastCommentID", fmt.Sprint(comm.ID))
	}

	if err := cache.GetCache().Content.Save(); err != nil {
		log.WithField("Func", "cache.Save").Error(err)
	}

	val, err = cache.GetCache().Check.Get("content", "LastCommentID")
	if err != nil {
		val = "0"
	}

	comms, err = model.GetContent(val)
	if err != nil {
		log.WithField("Func", "model.GetContent").Error(err)
	}

	for _, comm := range comms {
		if err := cache.GetCache().Check.Increment("check", fmt.Sprint(comm.ID), check.HashSet(comm.Content)); err != nil {
			log.WithField("Func", "cache.Increment").Error(err)
		}
	}

	if err := cache.GetCache().Check.Save(); err != nil {
		log.WithField("Func", "cache.Save").Error(err)
	}

	return listen
}

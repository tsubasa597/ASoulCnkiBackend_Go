package comment

import (
	"github.com/sirupsen/logrus"
	"github.com/tsubasa597/BILIBILI-HELPER/api"
	"github.com/tsubasa597/BILIBILI-HELPER/info"
	"github.com/tsubasa597/BILIBILI-HELPER/listen"
)

var _ listen.Listener = (*dynamic)(nil)

type dynamic struct {
	listen.Dynamic
}

func (dynamic dynamic) Listen(uid int64, _ api.API, log *logrus.Entry) []info.Infoer {

	return nil
}

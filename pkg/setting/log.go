package setting

import (
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// Write 写入日志
func Write() {
	if RunMode == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	writer, err := rotatelogs.New(
		path.Join(LogPath, "%Y-%m-%d.log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		logrus.WithError(err).Error("unable to write logs")
		return
	}

	logrus.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.DebugLevel: writer,
			logrus.InfoLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.FatalLevel: writer,
		}, &logrus.JSONFormatter{},
	))
}

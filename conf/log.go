package conf

import (
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func WriteLog() {
	writer, err := rotatelogs.New(
		path.Join(Path, "%Y-%m-%d.log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		panic(err)
	}

	w := lfshook.WriterMap{
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
	}

	if RunMode == "debug" {
		w[logrus.ErrorLevel] = writer
	}

	logrus.AddHook(lfshook.NewHook(w, &logrus.JSONFormatter{}))
}

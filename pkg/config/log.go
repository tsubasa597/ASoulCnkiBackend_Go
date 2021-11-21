package config

import (
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// NewLogFile 将日志写入文件
func NewLogFile(fileName string) (*logrus.Logger, error) {
	writer, err := rotatelogs.New(
		path.Join(LogPath, fileName, "%Y-%m-%d.log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		logrus.WithError(err).Error("unable to write logs")
		return nil, err
	}

	logger := logrus.New()
	logger.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.DebugLevel: writer,
			logrus.InfoLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.FatalLevel: writer,
		}, &logrus.JSONFormatter{},
	))

	return logger, nil
}

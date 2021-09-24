package setting

import (
	"fmt"

	"github.com/go-ini/ini"
)

var (
	RunMode, SQL, Host, User, Password, DBName, LogPath, CacheFilePath string
	Port, DefaultK, MaxOpen, DefaultPS, MaxConn, HeapLength, Size      int
	DynamicDuration, GoroutineNum                                      int64
	DefaultB                                                           float64
	Enable                                                             bool
)

func init() {
	cfg, err := ini.Load("./conf/conf.ini")
	if err != nil {
		panic(fmt.Sprintf("Fail to parse 'conf/app.ini': %v", err))
	}
	RunMode = cfg.Section("web").Key("RUN_MODE").MustString("debug")
	Port = cfg.Section("web").Key("PORT").MustInt(8000)
	Size = cfg.Section("web").Key("SIZE").MustInt(10)

	DefaultB = cfg.Section("check").Key("DEFAULT_B").MustFloat64(2)
	DefaultK = cfg.Section("check").Key("DEFAULT_K").MustInt(8)
	HeapLength = cfg.Section("check").Key("HEAP_LENGTH").MustInt(10)

	DefaultPS = cfg.Section("dynamic").Key("DEFAULT_PS").MustInt(49)

	MaxConn = cfg.Section("sql").Key("MAX_CONN").MustInt(100)
	MaxOpen = cfg.Section("sql").Key("MAX_OPEN").MustInt(10)
	SQL = cfg.Section("sql").Key("SQL").MustString("sqllite")
	Host = cfg.Section("sql").Key("HOST").MustString("")
	User = cfg.Section("sql").Key("USER").MustString("")
	Password = cfg.Section("sql").Key("PASSWORD").MustString("")
	DBName = cfg.Section("sql").Key("DBNAME").MustString("")

	LogPath = cfg.Section("log").Key("PATH").MustString("./log")

	Enable = cfg.Section("listen").Key("ENABLE").MustBool(true)
	DynamicDuration = cfg.Section("listen").Key("DYNAMIC_DURATION").MustInt64(5)
	GoroutineNum = cfg.Section("listen").Key("GOROUTINE_NUM").MustInt64(10)

	CacheFilePath = cfg.Section("cache").Key("PATH").MustString("./cache.dat")
}

package conf

import (
	"log"

	"github.com/go-ini/ini"
)

var (
	RunMode        string
	Port, DefaultK int
	DefaultB       float64
	HeapLength     int
	DefaultPS      int
	MaxConn        int
	MaxOpen        int
	SQL            string
	Host           string
	User           string
	Password       string
	DBName         string
)

func init() {
	cfg, err := ini.Load("./conf/conf.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}
	RunMode = cfg.Section("web").Key("RunMode").MustString("debug")
	Port = cfg.Section("web").Key("Port").MustInt(8000)
	DefaultB = cfg.Section("check").Key("DefaultB").MustFloat64(2)
	DefaultK = cfg.Section("check").Key("DefaultK").MustInt(8)
	HeapLength = cfg.Section("check").Key("HEAP_LENGTH").MustInt(10)
	DefaultPS = cfg.Section("dynamic").Key("DEFAULT_PS").MustInt(49)
	MaxConn = cfg.Section("sql").Key("MAX_CONN").MustInt(100)
	MaxOpen = cfg.Section("sql").Key("MAX_OPEN").MustInt(10)
	SQL = cfg.Section("sql").Key("SQL").MustString("sqllite")
	Host = cfg.Section("sql").Key("HOST").MustString("")
	User = cfg.Section("sql").Key("USER").MustString("")
	Password = cfg.Section("sql").Key("PASSWORD").MustString("")
	DBName = cfg.Section("sql").Key("DBNAME").MustString("")
}

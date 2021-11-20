package config

import (
	"fmt"

	"github.com/go-ini/ini"
)

var (
	// RunMode 运行模式
	RunMode string
	// SQL 数据库类型
	SQL string
	// Host 服务地址
	Host string
	// User 数据库用户名称
	User string
	// Password 数据库用户密码
	Password string
	// DBName 数据库名称
	DBName string
	// LogPath 日志保存地址
	LogPath string
	// CacheFilePath 缓存存放地址
	CacheFilePath string
	// RedisADDR redis 服务地址
	RedisADDR string
	// RedisPwd redis 密码
	RedisPwd string
	// Port 运行端口
	Port int
	// DefaultK 字符窗口长度
	DefaultK int
	// MaxOpen 数据库最大空闲连接数
	MaxOpen int
	// MaxConn 数据库最大连接数
	MaxConn int
	// HeapLength 查重返回数据的数量
	HeapLength int
	// Size 分页 一页的数据量
	Size int
	// DB redis 数据库号
	DB int
	// DynamicDuration 更新动态间隔时间(Min)
	DynamicDuration int64
	// DefaultB DefaultB
	DefaultB float64
	// Enable 是否启用自动更新
	Enable bool
	// RPCPath 服务地址
	RPCPath string
	// RPC 是否启动 RPC 服务
	RPCEnable bool
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

	RedisADDR = cfg.Section("redis").Key("ADDR").MustString("localhost:6379")
	RedisPwd = cfg.Section("redis").Key("PWD").MustString("")
	DB = cfg.Section("redis").Key("DB").MustInt(0)

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

	CacheFilePath = cfg.Section("cache").Key("PATH").MustString("./cache.dat")

	RPCPath = cfg.Section("rpc").Key("PATH").MustString("127.0.0.1:8080")
	RPCEnable = cfg.Section("rpc").Key("ENABLE").MustBool(false)
}

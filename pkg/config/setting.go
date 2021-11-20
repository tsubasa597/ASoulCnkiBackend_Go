package config

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

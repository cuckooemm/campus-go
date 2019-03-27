package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

type App struct {
	ServerHost    string
	ServerApiHost string
	PageSize      int

	JwtSecret string
	PrefixUrl string

	JwtExpired time.Duration
	JwtRefresh time.Duration

	RuntimeRootPath string

	ExportSavePath string
	QrCodeSavePath string
	FontSavePath   string

	LogSavePath string
	LogFileExt  string
	TimeFormat  string
}

type Server struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	HttpsCert    string
	HttpsKey     string
}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Port        string
	DBName      string
	LogModel    bool
	TablePrefix string
}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	CacheDB     int
	IdleTimeout time.Duration
}

type Email struct {
	Account  string
	Password string
	Host     string
	SendName string
}
type SNS struct {
	AccessKeyID     string
	AccessKeySecret string
	SignName        string
	TemplateCode    string
}
type OSS struct {
	AccesskeyID     string
	AccesskeySecret string
	Endpoint        string
	Path            string
	BucketName      string
	CallbackUrl     string
	Expiration      time.Duration
}
type STS struct {
	AccesskeyID     string
	AccesskeySecret string
	RoleArn         string
	RoleSessionName string
	DurationSeconds int
	Debug           bool
}
type UmengPush struct {
	AppKey                 string
	AndroidAppMasterSecret string
	Host                   string
	Mode                   bool
}
type WeChat struct {
	AppID     string
	AppSecret string
}

var AppSetting = &App{}
var ServerSetting = &Server{}
var DatabaseSetting = &Database{}
var RedisSetting = &Redis{}
var EmailSetting = &Email{}
var SNSSetting = &SNS{}
var OSSSetting = &OSS{}
var STSSetting = &STS{}
var PushSetting = &UmengPush{}
var WeChatSetting = &WeChat{}

var cfg *ini.File

func Setup() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Print("打开'conf/app.ini' 文件失败: %v", err)
		cfg,err = ini.Load("/root/web/conf/app.ini")
		if err != nil {
			log.Fatal("打开'/root/web/conf/app.ini' 文件失败: %v", err)
		}
	}
	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)
	mapTo("email", EmailSetting)
	mapTo("sns", SNSSetting)
	mapTo("oss", OSSSetting)
	mapTo("sts", STSSetting)
	mapTo("push", PushSetting)
	mapTo("weChat", WeChatSetting)

	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
	AppSetting.JwtExpired = AppSetting.JwtExpired * time.Hour
	AppSetting.JwtRefresh = AppSetting.JwtRefresh * time.Hour
	OSSSetting.Expiration = OSSSetting.Expiration * time.Second
	OSSSetting.Path = "https://" + OSSSetting.BucketName + "." + OSSSetting.Endpoint + "/"
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo RedisSetting err: %v", err)
	}
}

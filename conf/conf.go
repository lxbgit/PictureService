package conf

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"fmt"
	"qiniupkg.com/x/log.v7"
)

type qiNiu struct {
	Domain   string
	Bucket   string
	Key      string
	Secret   string
	Pipeline string
}

type upYun struct {
	Domain     string
	Secret     string // 表单密钥
	Bucket     string // 空间名（即服务名称）
	Expiration int    //策略过期时间
}

type DatabaseConfig struct {
	Driver string
	Host   string
	Port   int
	User   string
	DBName string
	PassWd string
}

var (
	HttpPort int

	SqliteDir              string
	SqliteFileName         string
	UserPassCodeEncryptKey string

	DebugMode        = false
	WebDebugMode     bool
	LogLevel         = "INFO"
	QiniuConfig      qiNiu
	UPYunConfig      upYun
	ServiceAddr      string
	RequestLogDir    string
	RequestLogEnable bool
	ErrorLogDir      string
	OpLogDir         string

	serviceConfigFile = flag.String("service-config", "__unset__", "service config file")
	debugMode         = flag.Bool("debug", false, "debug mode")
	logLevel          = flag.String("log-level", "INFO", "DEBUG | INFO | WARN | ERROR | FATAL | PANIC")
	webDebugMode      = flag.Bool("web-debug", false, "web debug mode")
)

func init() {
	flag.Parse()
	DebugMode = *debugMode
	LogLevel = *logLevel
	WebDebugMode = *webDebugMode

	if len(os.Args) == 2 {
		if os.Args[1] == "reload" {
			wd, _ := os.Getwd()
			pidFile, err := os.Open(filepath.Join(wd, "picture.pid"))
			if err != nil {
				log.Printf("Failed to open pid file: %s", err.Error())
				os.Exit(1)
			}
			pids := make([]byte, 10)
			n, err := pidFile.Read(pids)
			if err != nil {
				log.Printf("Failed to read pid file: %s", err.Error())
				os.Exit(1)
			}
			if n == 0 {
				log.Printf("No pid in pid file: %s", err.Error())
				os.Exit(1)
			}
			_, err = exec.Command("kill", "-USR2", string(pids[:n])).Output()
			if err != nil {
				log.Printf("Failed to restart picture service: %s", err.Error())
				os.Exit(1)
			}
			pidFile.Close()
			os.Exit(0)
		}
	}

	if DebugMode {
		LogLevel = "DEBUG"
	}

	if !DebugMode {
		nullFile, _ := os.Open(os.DevNull)
		log.SetOutput(nullFile)
		os.Stdout = nullFile
	}

	if *serviceConfigFile == "__unset__" {
		*serviceConfigFile = "./conf/service.json"
	}

	appConfig := make(map[string]interface{})
	appConfFileName, err := filepath.Abs(*serviceConfigFile)
	if err != nil {
		log.Fatal("Failed to format app_config_file path")
	}

	appConfFileContent, err := ioutil.ReadFile(appConfFileName)
	if err != nil {
		log.Fatal("Failed to read ", appConfFileName, ". error: %s", err)
	}

	err = json.Unmarshal(appConfFileContent, &appConfig)
	fmt.Println(appConfig)
	if err != nil {
		log.Println(err)
		return
	}

	//    setup database config
	setupQiniuConfig(appConfig["qiniu"])
	setupUPYunConfig(appConfig["upyun"])
	HttpPort = parseJsonInt(appConfig, "http_port", "")

	SqliteDir = parseJsonString(appConfig, "sqlite_dir", "")
	ServiceAddr = parseJsonString(appConfig, "service_addr", "")
	SqliteFileName = parseJsonString(appConfig, "sqlite_filname", "")
	UserPassCodeEncryptKey = parseJsonString(appConfig, "user_passcode_encrypt_key", "")
}

func parseJsonInt(config map[string]interface{}, key string, logPrefix string) (res int) {
	if t, ok := config[key].(float64); !ok {
		log.Fatal(logPrefix + key + " not set correctly, should be int value")
	} else {
		res = int(t)
	}

	return
}

func parseJsonString(config map[string]interface{}, key string, logPrefix string) (res string) {
	if t, ok := config[key].(string); !ok {
		log.Fatal(logPrefix + key + " not set correctly, should be string value")
	} else {
		res = t
	}

	return
}

func setupQiniuConfig(conf interface{}) {
	if conf == nil {
		log.Fatalln("No config for Qiniu")
	}

	data, ok := conf.(map[string]interface{})
	if !ok {
		log.Fatalln("Qiniu config not set correct")
	}

	qn := qiNiu{}
	qn.Domain = parseJsonString(data, "domain", "")
	qn.Key = parseJsonString(data, "key", "")
	qn.Secret = parseJsonString(data, "secret", "")
	qn.Bucket = parseJsonString(data, "bucket", "")
	qn.Pipeline = parseJsonString(data, "pipeline", "")

	s, _ := json.Marshal(qn)
	log.Println("Qiniu Config:", string(s))

	QiniuConfig = qn
}

func setupUPYunConfig(conf interface{}) {
	if conf == nil {
		log.Fatal("No config for UPYun")
	}
	data, ok := conf.(map[string]interface{})
	if !ok {
		log.Fatal("UPYun config not set correct")
	}
	up := upYun{}
	up.Bucket = parseJsonString(data, "bucket", "")
	up.Secret = parseJsonString(data, "secret", "")
	up.Domain = parseJsonString(data, "domain", "")
	up.Expiration = parseJsonInt(data, "expiration", "")
	s, _ := json.Marshal(up)
	log.Println("UPYun Config:", string(s))

	UPYunConfig = up
}

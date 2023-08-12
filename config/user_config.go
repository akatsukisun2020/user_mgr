package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

// GrpcGoConfig grpc_go的静态配置信息
type GrpcGoConfig struct {
	UserConfig UserConfig `yaml:"UserConfig"` // 用户配置
}

// UserConfig 用户自定义配置
type UserConfig struct {
	WXLoginAppID                    string `yaml:"WXLoginAppID"`                    // wx小程序登录appid
	WXLoginSecret                   string `yaml:"WXLoginSecret"`                   // wx小程序登录秘钥
	LoginExpireMillSecond           int64  `yaml:"LoginExpireMillSecond"`           // 登录过期微秒数
	LoginExpireRedundanceMillSecond int64  `yaml:"LoginExpireRedundanceMillSecond"` // 登录过期冗余微秒数
	RedisAddr                       string `yaml:"RedisAddr"`                       // redis地址
	RedisPasswd                     string `yaml:"RedisPasswd"`                     // redis秘钥
}

// gUserConfig 服务用户配置
var gUserConfig *GrpcGoConfig

func init() {
	file, err := ioutil.ReadFile("grpcgo_formal.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var grpcGoConfig GrpcGoConfig
	err = yaml.Unmarshal(file, &grpcGoConfig)
	if err != nil {
		log.Fatal(err)
	}

	gUserConfig = &grpcGoConfig
	log.Printf("InitUserConfig success, grpcGoConfig:%v", grpcGoConfig)
}

// 获取用户配置
func GetUserConfig() *UserConfig {
	if gUserConfig != nil {
		return &gUserConfig.UserConfig
	}
	return &UserConfig{}
}

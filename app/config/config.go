package config

import (
	"blog/app/models"
	"blog/app/support"

	"github.com/revel/config"
	"github.com/revel/revel"
)

const (
	CONFIG_PATH = "/conf/speedy.conf"
)

var AppConfig *config.Config
var IsInstalled bool

func InitConfig() {
	file := (revel.BasePath + CONFIG_PATH)
	var err error
	//检查配置文件是否存在
	AppConfig, err = config.ReadDefault(file)
	IsInstalled = true
	// 配置文件不存在
	if err != nil {
		revel.WARN.Println("获取配置文件失败，准备安装")
		IsInstalled = false
	} else {
		// 配置文件存在
		// 检查数据库是否可以正常连接
		err = support.InitXorm(AppConfig)
		if err != nil {
			IsInstalled = false
			revel.WARN.Println("连接数据库失败，准备安装")
		} else {
			// 数据库可以正常连接，同步表结构
			revel.WARN.Println("连接数据库成功，开始同步数据库")
			go models.SyncDB()
		}
	}
	if !IsInstalled {
		AppConfig = config.New(config.DEFAULT_COMMENT, config.ALTERNATIVE_SEPARATOR, false, true)
	} else {
		revel.INFO.Println("配置加载成功...")
	}
	support.InitCache(IsInstalled, AppConfig)
}

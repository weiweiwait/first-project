package main

import (
	conf "MyFirstProject/config"
	util "MyFirstProject/pkg/utils/log"
	"MyFirstProject/pkg/utils/track"
	"MyFirstProject/repository/cache"
	"MyFirstProject/repository/db/dao"
	"MyFirstProject/routes"
	"fmt"
)

func main() {
	loading() // 加载配置
	r := routes.NewRouter()
	_ = r.Run(conf.Config.System.HttpPort)
	fmt.Println("启动配成功...")
}

// loading一些配置
func loading() {
	conf.InitConfig()
	dao.InitMySQL()
	cache.InitCache()
	track.InitJaeger()
	util.InitLog() // 如果接入ELK请进入这个func打开注释
	fmt.Println("加载配置完成...")
	go scriptStarting()
}

func scriptStarting() {
	// 启动一些脚本
}

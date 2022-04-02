package main

import (
	config "go-micro-demo/conf/config_center"

	"github.com/asim/go-micro/v3/logger"
)

//配置中心测试
func main() {
	cfg, err := config.GetConfig("127.0.0.1", 8500, "/micro/config")
	if err != nil {
		logger.Fatal(err)
	}

	// Mysql配置信息
	mysqlInfo, err := config.GetMysqlFromConsul(cfg, "mysql")
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Mysql配置信息:", mysqlInfo)
}

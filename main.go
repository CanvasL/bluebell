package main

// @title bluebell项目接口文档
// @version 1.0
// @description gin框架实现论坛项目

// @contact.name LiYangfan
// @contact.url http://www.baidu.com

// @host 127.0.0.1:8084
// @BasePath /api/v1

import (
	"fmt"
	"go.uber.org/zap"
	"go_web/gin/bluebell/controller"
	"go_web/gin/bluebell/dao/mysql"
	"go_web/gin/bluebell/dao/redis"
	"go_web/gin/bluebell/logger"
	"go_web/gin/bluebell/pkg/snowflake"
	"go_web/gin/bluebell/router"
	"go_web/gin/bluebell/setting"
)

func main() {
	//1、加载配置
	if err := setting.Init(); err != nil {
		fmt.Printf("init setting failed, err:%v\n", err)
		return
	}

	//2、初始化日志
	if err := logger.Init(setting.Conf.LogConfig, setting.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer zap.L().Sync() //把缓冲区的日志追加到我们的日志当中
	zap.L().Debug("logger init success...")

	//3、初始化MySQL连接
	if err := mysql.Init(setting.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close()

	//4、初始化Redis连接
	if err := redis.Init(setting.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()

	//5、
	if err := snowflake.Init(setting.Conf.StartTime, setting.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
	}

	//6、初始化gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init validator failed, err:%v\n", err)
		return
	}

	//6、注册路由
	r := router.SetupRouter(setting.Conf.Mode)
	err := r.Run(fmt.Sprintf(":%d", setting.Conf.Port))
	if err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}
}

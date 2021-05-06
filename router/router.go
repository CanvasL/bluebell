package router

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go_web/gin/bluebell/controller"
	_ "go_web/gin/bluebell/docs"
	"go_web/gin/bluebell/logger"
	"go_web/gin/bluebell/middleware"
	"net/http"
	"time"
	"github.com/gin-contrib/pprof"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) //gin设置成发布模式
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.LoadHTMLFiles("./templates/index.html")
	//所有以"/static"开头的路径全部映射到"./static"路径中去
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	//注册
	v1.POST("/signup", controller.SignUpHandler)
	//登录
	v1.POST("/login", controller.LoginHandler)

	v1.Use(middleware.JWTAuthMiddleware(), middleware.RateLimitMiddleware(2 * time.Second, 1)) //应用JWT认证中间件

	v1.GET("/posts", controller.GetPostListHandler)
	// 根据时间或分数获取帖子列表
	v1.GET("/posts2", controller.GetPostListHandler2)

	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.GetPostDetailHandler)



		//投票
		v1.POST("/vote", controller.PostVoteController)
	}

	r.GET("/ping", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
		//如果是登录的用户，判断请求头中是否有有效的JWT
		c.String(http.StatusOK, "ok")

	})

	pprof.Register(r)	// 注册pprof相关路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}

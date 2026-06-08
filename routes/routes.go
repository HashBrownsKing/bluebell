package routes

import (
	"bluebell/controllers"
	"bluebell/logger"
	"bluebell/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode { // gin.ReleaseMode就是"release"
		gin.SetMode(gin.ReleaseMode) // 设置为发布模式，不会打印日志了
	}
	controllers.InitTrans("zh")
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// ==========================================
	// 配置静态文件路由，访问根路径时直接渲染前端页面
	// ==========================================
	r.StaticFile("/", "./templates/index.html")

	v1 := r.Group("api/v1")
	// 注册用户
	v1.POST("/signup", controllers.SignUpHandler)

	// 登录
	v1.POST("/login", controllers.LoginHandler)

	v1.Use(middlewares.JWTAuthMiddleware())

	{
		v1.GET("/community", controllers.CommunityHandler)           // 查询所有社区
		v1.GET("/community/:id", controllers.CommunityDetailHandler) //  通过id查找社区详细信息

		v1.POST("/post", controllers.CreatePostHandler)             // 创建帖子
		v1.GET("/postdetail/:id", controllers.GetPostDetailHandler) // 通过路径参数指定帖子id查找指定帖子
		v1.GET("/posts", controllers.GetPostListHandler)
		v1.GET("/posts2", controllers.GetPostListHandler2) // 按时间或分数排序获取帖子
		// v1.GET("/communityposts", controllers.GetCommunityPostListHandler) 合并到controllers.CreatePostHandler

		v1.POST("/vote", controllers.PostVoteHandler)
	}

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "hello world",
		})
	})

	return r
}

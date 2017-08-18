// @APIVersion 1.0
// @Title API
// @Description
// @Schemes http
// @Host texasholdem.com
package main

import (
	"chess/api/components/auth"
	"chess/api/components/middleware"
	"chess/api/controllers/auth"
	"chess/api/controllers/debug"
	"chess/api/controllers/room"
	"chess/api/controllers/user"
	"chess/common/config"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"time"
)

func InitRouter() {
	if config.Api.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE, OPTIONS",
		RequestHeaders:  "Origin, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, If-Modified-Since, x-requested-with",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))

	//router.Use(system_maintenance.CheckSystemMaintenance(*configServerPath))
	//router.Use(colletion.Filter(config.C.Filter))
	router.Use(middleware.SetContext())
	router.Use(config.SetContextConfig())
	if config.Api.Debug {
		debugRouter := router.Group("/debug")
		{
			// pprof
			debugRouter.GET("/pprof/", c_debug.PprofIndex())
			debugRouter.GET("/pprof/heap", c_debug.PprofHeap())
			debugRouter.GET("/pprof/goroutine", c_debug.PprofGoroutine())
			debugRouter.GET("/pprof/block", c_debug.PprofBlock())
			debugRouter.GET("/pprof/threadcreate", c_debug.PprofThreadCreate())
			debugRouter.GET("/pprof/cmdline", c_debug.PprofCmdline())
			debugRouter.GET("/pprof/profile", c_debug.PprofProfile())
			debugRouter.GET("/pprof/symbol", c_debug.PprofSymbol())

			debugRouter.GET("/ip", c_debug.IP)
			//debugRouter.GET("/config", c_debug.Config)
		}
	}

	// @SubApi /auth - 授权相关 [/auth/]
	authRouter := router.Group("/auth")
	{

		authRouter.POST("/login", c_auth.Login)                // 账号密码登录
		authRouter.POST("/login/quick", c_auth.LoginMobile)    // 手机号快速登录
		authRouter.POST("/login/tp", c_auth.TpLogin)           // 第三方登录
		authRouter.POST("/login/tourist", c_auth.TouristLogin) // 游客登录

		authRouter.POST("/token/info", c_auth.TokenInfo)       //获取token信息
		authRouter.POST("/token/refresh", c_auth.TokenRefrash) // 刷新token
		authRouter.POST("/register/mobile", c_auth.RegisterMobile)

		authRouter.GET("/test", c_auth.Ttest)
	}
	// @SubApi /room -房间相关 [/room/]
	roomRouter := router.Group("/room")
	{
		roomRouter.GET("/list", c_room.RoomsList)
	}

	// @SubApi /user/:user_id - 用户相关 [/user/{user_id}/]
	userRouter := router.Group("/user/:user_id")
	{
	    	userRouter.GET("/logout",auth.Login(config.C.TokenSecret) ,c_user.Logout)                        // 登出，销毁token
		userRouter.GET("/info", auth.Login(config.C.TokenSecret), c_user.GetUserInfo)
		userRouter.GET("/checkin", auth.Login(config.C.TokenSecret), c_user.Checkin)
		userRouter.POST("/password/reset", auth.Login(config.C.TokenSecret), c_user.PasswordReset)
	        userRouter.GET("/exchange",auth.Login(config.C.TokenSecret),c_user.Exchange)

	}
	// @SubApi /verify - 验证码相关 [/verify/]
	verifyRouter := router.Group("/verify")
	{
		verifyRouter.POST("/send", nil)
		verifyRouter.POST("/check", nil)
	}

	// @SubApi /client - 客户端配置相关 [/client/]
	clientRouter := router.Group("/client")
	{
		clientRouter.GET("/upgrade", nil)
	}
	//router.GET("/testquery",controllers.Get)
	router.Run(config.Api.Port)
}

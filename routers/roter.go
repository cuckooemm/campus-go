package routers

import (
	"campus/controllers"
	"campus/controllers/auth"
	"campus/controllers/common"
	"campus/controllers/dynamic"
	"campus/controllers/user"
	"campus/helper/setting"
	"campus/middleware/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	gin.SetMode(setting.ServerSetting.RunMode)
	apiV1 := r.Group("/api/v1")
	oauth := apiV1.Group("/auth")
	{
		oauth.POST("/login", auth.Login)
		oauth.POST("/login/qq",auth.QQLogin)
		oauth.POST("/login/weapp",auth.WeAppLogin)
		oauth.POST("/register", auth.Register)
		oauth.PUT("/code", auth.SendCode)
		oauth.POST("/retrieve", auth.Retrieve)
		// 刷新token
		oauth.GET("/refresh/token",auth.UpdateToken)
	}

	jwtAuth := apiV1.Group("/",jwt.JWT())
	userGroup := jwtAuth.Group("/user")
	{
		userGroup.GET("/info", user.UserInfo)
		userGroup.PUT("/info", user.Profile)
		userGroup.GET("/dynamic", user.UserDynamicList)
		userGroup.DELETE("/dynamic",user.UserDynamicDelete)
	}
	{
		jwtAuth.GET("oss/sts/token",auth.STSToken)
		jwtAuth.GET("oss/sign",auth.OSSSign)
		jwtAuth.POST("/dynamic/publish",dynamic.Publish)
		jwtAuth.POST("/dynamic/comment",dynamic.CommentPublish)
		jwtAuth.PUT("/dynamic/praise",dynamic.PraisePublish)
		jwtAuth.POST("/dynamic/reply",dynamic.ReplyPublish)
		jwtAuth.GET("/message",user.CommentToMeList)
	}

	apiV1.GET("/dynamic",dynamic.List)
	apiV1.GET("/dynamic/comment",dynamic.CommentList)
	apiV1.GET("/dynamic/reply",dynamic.ReplyList)
	apiV1.POST("/feedback",common.Publish)
	r.GET("test",controllers.Test)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound,gin.H{
			"code":404,
			"msg":"这里什么都没有",
		})
	})
	// api v2
	apiV2 := r.Group("/api/v2")
	apiV2.GET("/dynamic",dynamic.ListV2)
	apiV2.GET("/dynamic/comment",dynamic.CommentListV2)
	apiV2.GET("/dynamic/reply",dynamic.ReplyListV2)

	jwtAuthV2 := apiV2.Group("/",jwt.JWT())
	userGroupV2 := jwtAuthV2.Group("/user")
	{
		userGroupV2.GET("/dynamic", user.UserDynamicListV2)
		userGroupV2.GET("/message",user.UserMessageList)
	}
	{

	}
	return r
}

package router

import (
	"ginchat/docs"
	"ginchat/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/index", service.GetIndex)
	r.GET("/user/list", service.GetUserList)
	r.GET("/user/createuser", service.CreateUser)
	r.GET("/user/deleteuser", service.DeleteUser)
	r.POST("/user/updateuser", service.UpdateUser)
	r.POST("/user/finduserbynameandpwd", service.FindUserByNameAndPwd)
	return r
}

package service

import (
	"ginchat/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserList(c *gin.Context) {
	list := make([]*models.UserBasic, 10)
	list = models.GetUserList()
	c.JSON(http.StatusOK, gin.H{
		"message": list,
	})
}

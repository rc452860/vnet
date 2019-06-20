package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rc452860/vnet/model"
	"github.com/rc452860/vnet/service"
	"github.com/rc452860/vnet/utils/langx"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func StartServer(port int) {
	r := InitRouter()
	if err := r.Run(fmt.Sprintf(":%v", port)); err != nil {
		panic(err)
	}
}
func InitRouter() *gin.Engine {
	r := gin.Default()
	r1 := r.Group("/api")
	{
		r1.POST("/user/add", UserAdd)
		r1.POST("/user/del/:uid", UserDel)
		r1.POST("/user/edit", UserEdit)
		r1.GET("/user/list", UserList)
	}

	return r
}

func UserAdd(c *gin.Context) {
	var user model.UserInfo
	if err := c.ShouldBind(&user); err != nil {
		fail(c, err)
		return
	}
	logrus.Infof("add user,uid: %v, port: %v", user.Uid, user.Port)
	if err := service.ShadowsocksrServiceInstance.AddUser(&user); err != nil {
		fail(c, err)
		return
	}
	success(c)
}

func UserDel(c *gin.Context) {
	logrus.Infof("del uid: %v \n", c.Param("uid"))
	if err := service.ShadowsocksrServiceInstance.DelUser(langx.FirstResult(strconv.Atoi, c.Param("uid")).(int)); err != nil {
		fail(c, err)
		return
	}
	success(c)
}

func UserEdit(c *gin.Context) {
	var user model.UserInfo
	if err := c.ShouldBind(&user); err != nil {
		fail(c, err)
		return
	}
	logrus.Infof("edit user,uid: %v, port: %v", user.Uid, user.Port)
	if err := service.ShadowsocksrServiceInstance.EditUser(&user); err != nil {
		fail(c, err)
		return
	}
	success(c)

}

func UserList(c *gin.Context) {
	c.JSON(http.StatusOK, service.ShadowsocksrServiceInstance.GetUserList())
}

func fail(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{"success": "false", "content": err.Error()})
}
func success(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": "true", "content": "sucess"})
}

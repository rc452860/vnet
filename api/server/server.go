package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rc452860/vnet/model"
	"github.com/rc452860/vnet/service"
	"github.com/rc452860/vnet/utils/langx"
	"net/http"
	"strconv"
)

func StartServer(){
	r := InitRouter()
	if err := r.Run(":8081");err !=nil{
		panic(err)
	}
}
func InitRouter() *gin.Engine{
	r := gin.Default()
	r1 := r.Group("/api")
	{
		r1.POST("/user/add",UserAdd)
		r1.POST("/user/del/:uid",UserDel)
		r1.POST("/user/edit",UserEdit)
	}

	return r
}

func UserAdd(c *gin.Context){
	fmt.Println("add uid:"+c.Param("uid"))
	var user model.UserInfo
	if err := c.ShouldBind(&user);err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	service.ShadowsocksrServiceInstance.AddUser(&user)
	success(c)
}


func UserDel(c *gin.Context){
	fmt.Printf("del uid: %v \n",c.Param("uid"))
	service.ShadowsocksrServiceInstance.DelUser(langx.FirstResult(strconv.Atoi,c.Param("uid")).(int))
	success(c)
}


func UserEdit(c *gin.Context){
	fmt.Println("add uid:"+c.Param("uid"))
	var user model.UserInfo
	if err := c.ShouldBind(&user);err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	service.ShadowsocksrServiceInstance.EditUser(&user)
	success(c)

}

func success(c *gin.Context){
	c.JSON(http.StatusOK,gin.H{"success":"true","content":"sucess"})
}
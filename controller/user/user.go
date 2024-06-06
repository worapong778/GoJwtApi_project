package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/worapong778/GoJwtApi_project/orm"
)

// Read AllUsers
func ReadUsersAll(c *gin.Context) {
	var usersAll []orm.Tb_users
	orm.Db.Find(&usersAll)
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "read all users",
		"Data":    usersAll,
	})
}

// Read Profile
func Profile(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var user []orm.Tb_users
	orm.Db.First(&user, userId)
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "read profile",
		"Profile": user,
	})
}

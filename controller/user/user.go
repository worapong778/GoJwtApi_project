package user

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/worapong778/GoJwtApi_project/orm"
)

func ReadUsersAll(c *gin.Context) {
	hmacSampleSecret := []byte(os.Getenv("JWT_SECRET_KEY"))

	header := c.Request.Header.Get("Authorization")
	tokenString := strings.Replace(header, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})

	// check token exp
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// data request
		var UsersAll []orm.Tb_users
		orm.Db.Find(&UsersAll)
		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"message": "User Read All Success",
			"users":   UsersAll,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Forbidden",
			"message": err.Error(),
		})
		return
	}
}

package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/worapong778/GoJwtApi_project/orm"
	"golang.org/x/crypto/bcrypt"
)

var hmacSampleSecret []byte

// Binding from Json
type RegisterBody struct {
	Fname    string `json:"fname" binding:"required"`
	Lname    string `json:"lanme" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Tel      string `json:"tel" binding:"required"`
}

//******************************************** Register ************************************************//

func Register(c *gin.Context) {
	//Check ต้องใส่ข้อมูลให้ครบตามที่กำหนดใน RegisterBody
	var json RegisterBody
	err := c.ShouldBindJSON(&json) // check binding ต้องใส่ข้อมูลเข้ามา
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Check Register ซ้ำ โดยเช็กจาก email
	var UserEmailExist orm.Tb_users
	orm.Db.Where("User_email = ?", json.Email).First(&UserEmailExist)
	if UserEmailExist.Id > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "user exists",
		})
		return
	}

	// Encry Password
	encrytedPassword, _ := bcrypt.GenerateFromPassword([]byte(json.Password), 10)

	user := orm.Tb_users{
		User_fname:    json.Fname,
		User_lname:    json.Lname,
		User_email:    json.Email,
		User_password: string(encrytedPassword),
		User_tel:      json.Tel,
	}

	// Create data to Database
	orm.Db.Create(&user)

	// Check create error
	if user.Id > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "register success",
			//"userID":  user.Id,
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "register failed",
		})
	}
}

//************************************************ Login ********************************************//

// Binding from Json
type LoginBody struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	//Check ต้องใส่ข้อมูลให้ครบตามที่กำหนดใน LoginBody
	var json LoginBody
	err := c.ShouldBindJSON(&json) // check binding ต้องใส่ข้อมูลเข้ามา
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	//Check Email login
	var UserEmailExist orm.Tb_users
	orm.Db.Where("User_email = ?", json.Email).First(&UserEmailExist)
	if UserEmailExist.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "email not found",
		})
		return
	}

	//Check Email Password
	err = bcrypt.CompareHashAndPassword([]byte(UserEmailExist.User_password), []byte(json.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "incorrect password",
		})
		return
	}

	// JWT
	hmacSampleSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": UserEmailExist.Id,                      //ฝัง userId เข้าไปใน token ด้วย
		"exp":    time.Now().Add(time.Minute * 1).Unix(), // กำหนดเวลาหมดอายุ token
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "login success",
		"token":   tokenString,
	})
}

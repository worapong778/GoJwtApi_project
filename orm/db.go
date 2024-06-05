package orm

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Tb_users struct {
	Id            int
	User_fname    string
	User_lname    string
	User_email    string
	User_password string
	User_tel      string
}

var Db *gorm.DB
var err error

func InitDB() {
	dsn := os.Getenv("MYSQL_DNS")
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	Db.AutoMigrate(&Tb_users{})
}

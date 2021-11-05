package controller

import (
	"log"
	"net/http"

	"com.mego.first/megofirst/common"
	"com.mego.first/megofirst/model"
	"com.mego.first/megofirst/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(ctx *gin.Context) {
	db := common.InitDB()
	// get params
	username := ctx.PostForm("username")
	mobile := ctx.PostForm("mobile")
	password := ctx.PostForm("password")

	if len(mobile) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "invalid mobile"})
		return
	}

	if len(password) < 8 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "password too weak"})
		return
	}

	if len(username) == 0 {
		username = util.RandomString(8)
	}

	log.Println(username, mobile, password)
	// duplicate check
	if isMobileExist(db, mobile) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "mobile already exist"})
		return
	}
	if isUserExist(db, username) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "user already exist"})
		return
	}
	newUser := model.User{
		Username: username,
		Mobile:   mobile,
		Password: password,
	}

	db.Create(&newUser)

	ctx.JSON(200, gin.H{
		"message": "success",
	})
}

func isMobileExist(db *gorm.DB, mobile string) bool {
	var user model.User
	db.Where("mobile = ?", mobile).First(&user)
	return user.ID != 0
}

func isUserExist(db *gorm.DB, username string) bool {
	var user model.User
	db.Where("username = ?", username).First(&user)
	return user.ID != 0
}

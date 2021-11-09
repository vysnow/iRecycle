package controller

import (
	"log"
	"net/http"

	"com.mego.first/megofirst/common"
	"com.mego.first/megofirst/model"
	"com.mego.first/megofirst/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

	hashedPawword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "error happening on password encryption"})
		return
	}

	newUser := model.User{
		Username: username,
		Mobile:   mobile,
		Password: string(hashedPawword),
	}

	db.Create(&newUser)

	ctx.JSON(200, gin.H{
		"message": "success",
	})
}

func Login(ctx *gin.Context) {
	db := common.InitDB()
	// get params
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	// params ok?
	if len(password) < 8 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "password too weak"})
		return
	}

	if len(username) == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "username cannot be empty"})
		return
	}

	// if exist?
	var user model.User
	db.Where("username = ?", username).First(&user)

	if user.ID == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "user not exist"})
		return
	}

	// password correct?
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "username / password not match"})
		return
	}

	// generate token
	token := "token"

	// response
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"code": 200,
			"msg":  "login success",
			"data": gin.H{
				"token": token,
			},
		},
	)
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

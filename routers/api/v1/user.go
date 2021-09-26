package v1

import (
	"net/http"
	"onlineShopping/middleware"
	"onlineShopping/models"
	"onlineShopping/models/requset"
	"onlineShopping/pkg/app"
	"onlineShopping/pkg/e"
	"onlineShopping/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type IUser interface {
	Login(c *gin.Context)
	LoginOut(c *gin.Context)
	Register(c *gin.Context)

	CheckToken(c *gin.Context)

	MyInformation(c *gin.Context)
}

type UserManager struct {
	DB *gorm.DB
}

// CheckToken 用户详情
func (u UserManager) CheckToken(c *gin.Context) {
	app.Response(c,
		http.StatusOK,
		e.SUCCESS,
		nil,
	)
}

func (u UserManager) Login(c *gin.Context) {
	var rUser requset.UserSignIn

	err := c.ShouldBindJSON(&rUser)
	if err != nil {
		app.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	exist, err := models.ExistUser(util.EncodeMD5(rUser.Telephone), util.EncodeMD5(rUser.Password))
	if err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, nil)
		return
	}
	if !exist {
		app.Response(c, http.StatusUnprocessableEntity, e.ERROR_NOT_EXIST_USER, nil)
		return
	}

	j := middleware.NewJWT()
	token, err := j.GenerateToken(models.User{Telephone: rUser.Telephone})
	if err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_GENERATE_TOKEN_FAIL, nil)
		return
	}
	app.Response(c, http.StatusOK, e.SUCCESS, gin.H{
		"token": token,
	})

}

func (u UserManager) LoginOut(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (u UserManager) Register(c *gin.Context) {
	var rUser requset.UserSignUp
	err := c.ShouldBindJSON(&rUser)
	if err != nil {
		app.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	exist, err := models.ExistUserByTelephone(rUser.Telephone)
	if err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_GET_USER_TELEPHONE_FAIL, nil)
		return
	}
	if exist {
		app.Response(c, http.StatusUnprocessableEntity, e.ERROR_ADD_EXIST_USER, nil)
		return
	}

	user := models.User{
		Username:  rUser.Username,
		Password:  util.EncodeMD5(rUser.Password),
		Telephone: util.EncodeMD5(rUser.Telephone),
	}

	if err := u.DB.Create(&user).Error; err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_ADD_USER_FAIL, nil)
		return
	}

	j := middleware.NewJWT()
	token, err := j.GenerateToken(models.User{Telephone: user.Telephone})
	if err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_GENERATE_TOKEN_FAIL, nil)
		return
	}

	app.Response(c, http.StatusOK, e.SUCCESS, gin.H{
		"token": token,
	})
}

func (u UserManager) MyInformation(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	err := u.DB.Preload("Orders").Where("id = ?", id).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		// app.Response(c, http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, nil)
		app.Response(c, http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, nil)
		return
	}

	app.Response(c, http.StatusOK, e.SUCCESS, gin.H{"user": user})
}

func NewUserManager() IUser {
	db := models.GetDB()
	db.AutoMigrate(&models.User{})
	return UserManager{DB: db}
}

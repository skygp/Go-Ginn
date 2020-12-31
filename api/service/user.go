package service

import (
	"ginn/api/dao/mysql"
	"ginn/api/dao/redis"
	"ginn/api/models"
	"ginn/config"
	"ginn/package/jwt"
	"ginn/utils"
	"time"

	"errors"

	errMsg "ginn/utils/code"
)

func RegisterUser(user *models.RequestRegisterUser) (err error) {
	err = mysql.CreateUser(user)
	return
}

func LoginUser(user *models.RequestLoginUser) (token string, err error) {
	var resultUser *models.DbUser
	username := user.UserName
	if !mysql.FindUserIsExist(username) {
		return "", errors.New(errMsg.ReturnCodeMsg(errMsg.CodeUserNotExist))
	}
	resultUser, err = mysql.FindUserName(username)
	if nil != err {
		return
	}
	err = utils.CompareHashAndPassword(resultUser.PassWord, user.PassWord)
	if nil != err {
		return
	}
	token, err = jwt.GenToken(resultUser.Uid, resultUser.UserName)
	if nil != err {
		return
	}
	err = redis.SetKey(token, resultUser.Uid, time.Duration(config.Gconfig.Auth.JwtExpire) * time.Second)
	return
}

func GetUserInfo(token string)(role, userName string, err error)  {
	var (
		uid string
		resultUser *models.DbUser
	)
	uid, err = redis.GetKey(token)
	if nil != err{
		return "","",errors.New(errMsg.ReturnCodeMsg(errMsg.CodeInvalidToken))
	}
	resultUser, err = mysql.FindUid(uid)
	if nil != err{
		return "","",errors.New(errMsg.ReturnCodeMsg(errMsg.CodeUserNotExist))
	}
	role = models.GetRole(resultUser.Role)
	userName = resultUser.UserName
	return
}

func Logout(token string) (err error) {
	err = redis.DelKey(token)
	if nil != err{
		return errors.New(errMsg.ReturnCodeMsg(errMsg.CodeInvalidToken))
	}
	return
}
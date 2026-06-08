package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

func SignUp(param *models.ParamSignUp) (err error) {
	// 检查用户是否存在
	var exist bool
	exist, err = mysql.CheckExistsUser(param)
	if err != nil {
		return err
	}
	if exist {
		return mysql.ErrorUserExist
	}
	// 生成UID
	uid := snowflake.GetID()
	// 创建一个用户实例
	user := &models.User{
		UserID:   uid,
		UserName: param.Username,
		Password: param.Password,
	}
	// 向数据库插入用户
	return mysql.InsertUser(user)
}

func Login(param *models.ParamLogin) (token string, err error) {
	user := new(models.User)
	user.UserName = param.Username
	user.Password = param.Password
	if err = mysql.Login(user); err != nil {
		return "", err
	}
	// 生成JWT
	return jwt.GenerateToken(user.UserID, user.UserName)
}

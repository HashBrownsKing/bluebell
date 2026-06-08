package mysql

import (
	"bluebell/models"
	"database/sql"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func CheckExistsUser(param *models.ParamSignUp) (exist bool, err error) {
	sqlStr := `SELECT COUNT(*) FROM bluebell_user WHERE username = ?`
	var count int
	if err = db.Get(&count, sqlStr, param.Username); err != nil {
		return false, err
	}
	return count > 0, nil
}

func InsertUser(user *models.User) (err error) {
	sqlStr := `INSERT INTO bluebell_user (user_id, username, password) VALUES (?, ?, ?)`
	var encryptPassword string
	if encryptPassword, err = EncryptPassword(user.Password); err != nil {
		return err
	}
	_, err = db.Exec(sqlStr, user.UserID, user.UserName, encryptPassword)
	return err
}

func EncryptPassword(password string) (string, error) {
	// GenerateFromPassword 自动生成盐并加密，bcrypt.DefaultCost表示默认安全系数
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func Login(user *models.User) (err error) {
	sqlStr := `SELECT user_id, username, password FROM bluebell_user WHERE username = ?`
	var u models.User
	if err = db.Get(&u, sqlStr, user.UserName); err != nil {
		if err == sql.ErrNoRows {
			return ErrorUserNotExist
		}
		return err
	}
	// 修复：把数据库查询到的 用户ID 赋值给入参！！！
	user.UserID = u.UserID
	// 可选：用户名也赋值
	user.UserName = u.UserName
	if err = CheckPassword(user.Password, u.Password); err != nil {
		return ErrorInvalidPassword
	}
	return
}

func GetUsernameByID(userID int64) (username string, err error) {
	sqlStr := `SELECT username From bluebell_user WHERE user_id = ?`
	err = db.Get(&username, sqlStr, userID)
	if err != nil {
		zap.L().Error("mysql.GetUsernameByID failed", zap.Error(err))
		return "", err
	}
	return
}

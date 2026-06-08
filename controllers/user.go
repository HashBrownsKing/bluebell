package controllers

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUpHandler - 处理注册请求的函数
func SignUpHandler(c *gin.Context) {
	// 1. 参数校验
	p := new(models.ParamSignUp)
	// ShouldBindJSON只能检验格式是否为JSON格式
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUp Params Error: ", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors) // 类型断言判断是否是参数校验的错误
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			// zap.L().Error("SignUp Params Error: ", zap.Error(err))
			// c.JSON(http.StatusOK, gin.H{
			// 	"msg": err.Error(),
			// })
			ResponseError(c, CodeInvalidParam)
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		// c.JSON(http.StatusOK, gin.H{
		// 	"msg": removeTopStruct(errs.Translate(trans)),
		// })
		return
	}
	/*
		// 检验参数的值是否正确
		// 可以使用gin框架中自带的validator库来帮组完成这种操作，只需要在需要绑定的参数的结构体中给字段绑定 binding : "xxx" 的tag
		// 默认的报错信息为英文，但是validator支持翻译
		if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.Password != p.RePassword {
			zap.L().Error("SignUp Params Error: ")
			c.JSON(http.StatusOK, gin.H{
				"msg": "Params Error",
			})
			return
		}
	*/
	// 2. 处理业务
	if err := logic.SignUp(p); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{
		// 	"msg": err.Error(),
		// })
		zap.L().Error("logic.SignUp failed: ", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
	}
	// 3. 返回响应
	// c.JSON(http.StatusOK, gin.H{
	// 	"msg": "success",
	// })
	ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	// 参数验证
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Login Params Error: ", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			// zap.L().Error("Login Params Error: ", zap.Error(err))
			// c.JSON(http.StatusOK, gin.H{
			// 	"msg": err.Error(),
			// })
			zap.L().Error("Login Params Error: ", zap.Error(err))
			ResponseError(c, CodeInvalidParam)
			return
		}
		// c.JSON(http.StatusOK, gin.H{
		// 	"msg": removeTopStruct(errs.Translate(trans)),
		// })
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 业务处理
	token, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed: ", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
		} else if errors.Is(err, mysql.ErrorInvalidPassword) {
			ResponseError(c, CodeInvalidPassword)
		} else {
			ResponseError(c, CodeServerBusy)
		}
		// c.JSON(http.StatusBadRequest, gin.H{
		// 	"msg": errors.New("用户名或密码错误").Error(),
		// })
		return
	}
	// 返回响应
	// c.JSON(http.StatusOK, gin.H{
	// 	"msg": "登录成功",
	// })
	ResponseSuccess(c, token)
}

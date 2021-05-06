package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"go_web/gin/bluebell/dao/mysql"
	"go_web/gin/bluebell/logic"
	"go_web/gin/bluebell/model"
)

// SignUpHandler 处理注册请求的函数
func SignUpHandler(c *gin.Context) {
	// 1.获取参数和参数检验
	var p model.ParamSignUp
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误，返回参数
		zap.L().Error("Signup with invalid params", zap.Error(err))
		//判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		//err是validator的类型，则需要用翻译器翻译
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": removeTopStruct(errs.Translate(trans)), //翻译字段
		//})
		return
	}

	// 2.业务处理
	if err := logic.SignUp(&p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": "注册失败",
		//})
		return
	}

	// 3.返回响应
	//c.JSON(http.StatusOK, gin.H{
	//	"msg": "success",
	//})
	ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	//1.获取请求参数及参数校验
	p := new(model.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误，返回参数
		zap.L().Error("Login with invalid params", zap.Error(err))
		//判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		//err是validator的类型，则需要用翻译器翻译
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": removeTopStruct(errs.Translate(trans)), //翻译字段
		//})
		return
	}
	//2.业务逻辑处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}

	//3.返回响应
	ResponseSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID), //防止id值大小超出前端js能表示的数值范围
		"user_name": user.Username,
		"token":     user.Token,
	})
}

package logic

import (
	"go_web/gin/bluebell/dao/mysql"
	"go_web/gin/bluebell/model"
	"go_web/gin/bluebell/pkg/jwt"
	"go_web/gin/bluebell/pkg/snowflake"
)

// 存放业务逻辑的代码

func SignUp(p *model.ParamSignUp) (err error) {
	// 1.判断用户存不存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	// 2.生成UID
	userID := snowflake.GenID()
	//构造一个User实例
	user := &model.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	// 3.密码加密

	// 4.保存进数据库
	err = mysql.InsertUser(user)
	return err
}

func Login(p *model.ParamLogin) (user *model.User, err error) {
	user = &model.User{
		Username: p.Username,
		Password: p.Password,
	}
	//传递的是指针，就能拿到user.UserID
	if err = mysql.Login(user); err != nil {
		return nil, err
	}
	//生成JWT
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}

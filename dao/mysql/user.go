package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"go_web/gin/bluebell/model"
)

// 把每一步数据库操作封装成函数
// 待logic层根据业务需求调用

const secret = "canvas.com"

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	//写sql，查询下用户数量
	sqlStr := `select count(user_id) from user where username=?`
	var count int
	if err = db.Get(&count, sqlStr, username); err != nil {
		//查询出问题
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return nil
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *model.User) (err error) {
	//对密码进行加密
	user.Password = encryptPassword(user.Password)
	//执行sql语句入库
	sqlStr := `insert into user(user_id,username,password) values (?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

// encryptPassword 密码加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

// Login 登录
func Login(user *model.User) (err error) {
	oPassword := user.Password
	sqlStr := `select user_id,username,password from user where username=?`
	if err = db.Get(user, sqlStr, user.Username); err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		//查询数据库失败
		return err
	}
	//查询成功，未报错，则去查询密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}

// GetUserByID 根据id获取用户信息
func GetUserByID(uid int64) (user *model.User, err error) {
	user = new(model.User)
	sqlStr := `select user_id,username from user where user_id=?`
	err = db.Get(user, sqlStr, uid)
	return
}

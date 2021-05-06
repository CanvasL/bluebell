package logic

import (
	"go_web/gin/bluebell/dao/mysql"
	"go_web/gin/bluebell/model"
)

func GetCommunityList() (communityList []*model.Community, err error) {
	//查找到所有的community并返回
	communityList, err = mysql.GetCommunityList()
	return
}

func GetCommunityDetail(id int64) (*model.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}

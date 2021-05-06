package mysql

import (
	"database/sql"
	"go.uber.org/zap"
	"go_web/gin/bluebell/model"
)

func GetCommunityList() (communityList []*model.Community, err error) {
	//写sql
	sqlStr := `select community_id,community_name from community`
	//执行
	if err = db.Select(&communityList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db.")
			err = nil
		}
	}
	return
}

// GetCommunityDetailByID 根据ID查询社区详情
func GetCommunityDetailByID(id int64) (community *model.CommunityDetail, err error) {
	community = new(model.CommunityDetail)
	sqlStr := `select community_id,community_name,introduction,create_time from community where community_id=?`
	if err = db.Get(community, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}
	return community, err
}

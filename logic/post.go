package logic

import (
	"go.uber.org/zap"
	"go_web/gin/bluebell/dao/mysql"
	"go_web/gin/bluebell/dao/redis"
	"go_web/gin/bluebell/model"
	"go_web/gin/bluebell/pkg/snowflake"
)

func CreatePost(p *model.Post) (err error) {
	// 1.生成post id
	p.ID = snowflake.GenID()
	// 2.保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	// 3.返回
	return
}

// GetPostByID 根据帖子id查询帖子详情数据
func GetPostByID(pid int64) (data *model.ApiPostDetail, err error) {
	//查询并组合我们接口想用的数据
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostByID failed", zap.Int64("pid", pid), zap.Error(err))
		return
	}
	//根据作者id查询作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
		return
	}
	//根据社区id查询社区详细信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
		return
	}
	data = &model.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (data []*model.ApiPostDetail, err error) {
	//查询并组合我们接口想用的数据
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList failed", zap.Error(err))
		return
	}
	data = make([]*model.ApiPostDetail, 0, len(posts)) //帖子长度作为容量
	for _, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}

		//根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}

		postDetail := &model.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

func GetPostList2(p *model.ParamPostList) (data []*model.ApiPostDetail, err error) {
	// 去redis查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	// 根据id去MySQL数据库查询帖子详细信息
	// 返回的数据还要按照我给定的id顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	//提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	// 将帖子作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}

		//根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}

		postDetail := &model.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

func GetCommunityPostList(p *model.ParamPostList) (data []*model.ApiPostDetail, err error) {
	// 去redis查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	// 根据id去MySQL数据库查询帖子详细信息
	// 返回的数据还要按照我给定的id顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return nil, err
	}
	//提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}
	// 将帖子作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}

		//根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}

		postDetail := &model.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}


// GetPostListNew 将两个查询逻辑合二为一的函数
func GetPostListNew(p *model.ParamPostList) (data []*model.ApiPostDetail, err error) {
	if p.CommunityID == 0 {
		// 获取所有数据
		data, err = GetPostList2(p)
	} else {
		// 根据社区id获取数据
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}
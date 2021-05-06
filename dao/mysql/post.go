package mysql

import (
	"github.com/jmoiron/sqlx"
	"go_web/gin/bluebell/model"
	"strings"
)

// CreatePost 创建帖子
func CreatePost(p *model.Post) (err error) {
	sqlStr := `insert into post (post_id,title,content,author_id,community_id) values (?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// GetPostByID 根据id查询单个帖子数据详情
func GetPostByID(pid int64) (post *model.Post, err error) {
	post = new(model.Post)
	sqlStr := `select post_id,title,content,author_id,community_id,create_time from post where post_id=?`
	err = db.Get(post, sqlStr, pid)
	return
}

// GetPostList 查询帖子列表函数
func GetPostList(page, size int64) (posts []*model.Post, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,create_time from post order by create_time desc limit?,?`
	posts = make([]*model.Post, 0, 2)
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}

// GetPostListByIDs 根据给定的id列表查询数据
func GetPostListByIDs(ids []string) (postList []*model.Post, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,create_time from post where post_id in (?) order by FIND_IN_SET(post_id,?)`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query) //注意Rebind的用法

	err = db.Select(&postList, query, args...) //不要忘记...
	return
}

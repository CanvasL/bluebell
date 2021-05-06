package mysql

import (
	"go_web/gin/bluebell/model"
	"go_web/gin/bluebell/setting"
	"testing"
)

func init() {
	dbCfg := &setting.MySQLConfig{
		Host: "127.0.0.1",
		User: "root",
		Password: "199732skfly",
		DB: "bluebell",
		MaxOpenConns: 200,
		MaxIdleConns: 50,
	}
	err := Init(dbCfg)
	if err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	post := &model.Post{
		ID:	10,
		AuthorID: 123,
		CommunityID: 1,
		Title:"test",
		Content: "just a test",
	}
	err := CreatePost(post)
	if err != nil {
		t.Fatalf("CreatePost failed, err:%v\n", err)
	}
	t.Logf("success")
}
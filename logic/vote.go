package logic

import (
	"go.uber.org/zap"
	"go_web/gin/bluebell/dao/redis"
	"go_web/gin/bluebell/model"
	"strconv"
)

//投票功能

//本项目使用简化版的投票分数
// 投一票就加432分 86400/200	->需要200张赞成票才能给帖子续一天

/* 投票的几种情况
direction=1时，有两种情况：
	1.之前没有投过票，现在投赞成票；		--->更新分数和投票记录
	2.之前投反对票，现在要改投赞成票；	--->更新分数和投票记录
direction=-0时，有两种情况：
	1.之前投过赞成票，现在要取消投票；	--->更新分数和投票记录
	2.之前投过反对票，现在要取消投票；	--->更新分数和投票记录
direction=-1时，有两种情况：
	1.之前没有投过票，现在投反对票；		--->更新分数和投票记录
	2.之前投过赞成票，现在要改投反对票；	--->更新分数和投票记录

投票的限制：
每个帖子自发表之日起，一个星期之内允许用户投票，超过一个星期就不允许再投票了。
	1.到期之后将redis中保存的赞成票数及反对票数存储到mysql表中
	2.到期之后删除保存的KeyPostVotedZSetPrefix
*/

// PostForVote 为帖子投票的函数
func PostForVote(userID int64, p *model.ParamVoteData) error {
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", userID),
		zap.String("postID", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}

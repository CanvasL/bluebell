package redis

import (
	"github.com/go-redis/redis"
	"go_web/gin/bluebell/model"
	"strconv"
	"time"
)

func getIDsFormKey(key string, page, size int64) ([]string, error) {
	// 2.确定查询的索引起始点
	start := (page - 1) * size
	end := start + size - 1
	// 3.ZREVRANGE 按分数从大到小的顺序查询指定数量的元素
	return client.ZRevRange(key, start, end).Result()
}

func GetPostIDsInOrder(p *model.ParamPostList) ([]string, error) {
	// 1.从redis获取id
	// 根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == model.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	return getIDsFormKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇帖子投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	data = make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZSetPF+id)
	//	//查找key中分数是1的元素的数量->统计每篇帖子赞成票的数量
	//	v := client.ZCount(key, "1", "1").val()
	//	data = append(data, v)
	//}

	//使用pipeline一次发送多条命令，减少RTT
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return data, nil
}

// GetCommunityPostIDsInOrder 按社区查询ids
func GetCommunityPostIDsInOrder(p *model.ParamPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == model.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	// 使用zinterstore把分区的帖子set与帖子分数的zset生成一个新的zset
	// 针对新的zset，按之前的逻辑取数据

	// 社区的key
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))

	// 利用缓存key减少zinterstore的执行次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(key).Val() < 1 {
		// 不存在，需要计算
		pipeline := client.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, orderKey) // zinterstore 计算
		pipeline.Expire(key, 60*time.Second) // 设置超时时间
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	// 存在的话直接根据key查询id即可
	return getIDsFormKey(key, p.Page, p.Size)
}

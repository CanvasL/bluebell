package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_web/gin/bluebell/logic"
	"strconv"
)

// 跟社区相关的

func CommunityHandler(c *gin.Context) {
	// 查询到所有的社区(community_id,community_name)以列表的形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务器报错暴露给外面
		return
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 社区分类详情
func CommunityDetailHandler(c *gin.Context) {
	// 获取社区id
	idStr := c.Param("id") //获取路径参数
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		//如果路径不对，返回参数错误
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 查询到所有的社区(community_id,community_name)以列表的形式返回
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不轻易把服务器报错暴露给外面
		return
	}
	ResponseSuccess(c, data)
}

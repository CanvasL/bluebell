package controller

import "go_web/gin/bluebell/model"

//专门用来放接口文档用到的model
//因为我们的接口文档返回的数据格式是一致的，但是具体的data类型不一致

type _ResponsePostList struct {
	Code     ResCode                `json:"code"`		// 业务的响应状态码
	Messsage string                 `json:"message"`	// 提示信息
	Data     []*model.ApiPostDetail `json:"data"`		// 数据
}

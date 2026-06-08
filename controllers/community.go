package controllers

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// --------------社区相关-------------

func CommunityHandler(c *gin.Context) {
	// 查询和所有的社区(community_id, community_name)，用Slice存储
	communityList, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed : ", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不轻易把服务错误暴露给外面
		return
	}
	ResponseSuccess(c, communityList)
}

func CommunityDetailHandler(c *gin.Context) {
	// 根据路径参数中的id查询社区的信息
	idStr := c.Param("id")
	communityID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	communityDetail, err := logic.GetCommunityDetailByID(communityID)
	if err != nil {
		if errors.Is(err, mysql.ErrorInvalidID) { // 当ID不存在的时候
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, communityDetail)
}

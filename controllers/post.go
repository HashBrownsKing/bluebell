package controllers

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func CreatePostHandler(c *gin.Context) {
	// 1. 参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("参数有误: ", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	userID, err := GetCurrentUserID(c) // 验证登录时已经将用户ID存入gin.Context上下文
	if err != nil {
		zap.L().Error("需要登录: ", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	// 2. 处理业务
	if err = logic.CreatePost(p); err != nil {
		zap.L().Error("创建帖子失败：", zap.Error(err))
		if errors.Is(err, mysql.ErrorCommunityNotExist) {
			ResponseError(c, CodeCommunityNotExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// 通过路径参数中的post_id获取帖子详情
func GetPostDetailHandler(c *gin.Context) {
	//  获取post_id
	postIDStr := c.Param("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		zap.L().Error("controllers.GetPostDetialHandler failed: ", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 获取帖子详细信息
	data, err := logic.GetPostDetailByID(postID)
	if err != nil {
		if errors.Is(err, mysql.ErrorPostNotExists) {
			zap.L().Error("controllers.GetPostDetailHandler failed: ", zap.Error(err))
			ResponseError(c, CodePostNotExists)
			return
		}
		zap.L().Error("controllers.GetPostDetailHandler failed: ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

func GetPostListHandler(c *gin.Context) {
	page, size := getPageInfo(c)

	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostListHandler failed: ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, data)
}

/*
1. 获取前端传来的参数：page, size, order
2. 根据orderType去redis里找ID（redis中有根据时间和分数排序的表）
3. 根据查到的ID去mysql中查找数据
*/
func GetPostListHandler2(c *gin.Context) {
	// 参数获取和校验
	p := &models.ParamPostList{
		Page:  1,
		Size:  5,
		Order: models.OrderTime,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("controller.GetPostListHandler2 failed: ", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errorData := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParam, errorData)
		return
	}
	zap.L().Info("接收的社区ID", zap.Int64("community_id", p.CommunityID))

	// 业务处理
	data, err := logic.GetPostListMerge(p)
	if err != nil {
		zap.L().Error("controller.GetPostListHandler2 failed: ", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 返回响应
	ResponseSuccess(c, data)
}

// func GetCommunityPostListHandler(c *gin.Context) {
// 	// 获取参数和校验参数
// 	p := &models.ParamCommunityPostList{
// 		CommunityID: 1,
// 		ParamPostList: models.ParamPostList{
// 			Page:  1,
// 			Size:  5,
// 			Order: "time",
// 		},
// 	}
// 	if err := c.ShouldBindQuery(p); err != nil {
// 		zap.L().Error("controllers.GetCommunityPostListHandler failed: ", zap.Error(err))
// 		errs, ok := err.(validator.ValidationErrors)
// 		if !ok {
// 			ResponseError(c, CodeInvalidParam)
// 			return
// 		}
// 		errorData := removeTopStruct(errs.Translate(trans))
// 		ResponseErrorWithMsg(c, CodeInvalidParam, errorData)
// 	}

// 	// 业务处理
// 	data, err := logic.GetCommunityPostList(p)
// 	if err != nil {
// 		zap.L().Error("controllers.GetCommunityPostListHandler failed: ", zap.Error(err))
// 		ResponseError(c, CodeServerBusy)
// 		return
// 	}
// 	// 返回响应
// 	ResponseSuccess(c, data)
// }

func getPageInfo(c *gin.Context) (int64, int64) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil || size < 1 { // 不传、非数字、负数、0 → 统一默认5
		size = 5
	}
	if size > 50 {
		size = 50
	}
	return page, size
}

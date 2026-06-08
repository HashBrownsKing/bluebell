package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"

	"go.uber.org/zap"
)

func CreatePost(post *models.Post) (err error) {
	var exists bool
	exists, err = mysql.CheckCommunityExists(post.CommunityID)
	if err != nil {
		return
	}
	if !exists {
		return mysql.ErrorCommunityNotExist
	}
	return mysql.CreatePost(post)
}

func GetPostDetailByID(postID int64) (data *models.ApiPostDetail, err error) {
	// 查询帖子信息
	post, err := mysql.GetPostByID(postID)
	if err != nil {
		zap.L().Error("logic.GetPostDetailByID failed: ", zap.Error(err))
		return nil, err
	}
	// 查询作者姓名
	username, err := mysql.GetUsernameByID(post.AuthorID)
	if err != nil {
		zap.L().Error("logic.GetPostDetailByID failed: ", zap.Error(err))
		return nil, err
	}
	// 查询社区信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("logic.GetPostDetailByID failed: ", zap.Error(err))
		return nil, err
	}
	data = &models.ApiPostDetail{
		AuthorName:      username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

func GetPostList(page, size int64) ([]*models.ApiPostDetail, error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList failed: ", zap.Error(err))
		return nil, err
	}
	postList := make([]*models.ApiPostDetail, 0, size)
	for _, post := range posts {
		// 查询作者姓名
		username, err := mysql.GetUsernameByID(post.AuthorID)
		if err != nil {
			zap.L().Error("logic.GetPostList failed: ", zap.Error(err))
			continue
		}
		// 查询社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("logic.GetPostList failed: ", zap.Error(err))
			continue
		}
		data := &models.ApiPostDetail{
			AuthorName:      username,
			Post:            post,
			CommunityDetail: community,
		}
		postList = append(postList, data)
	}
	return postList, err
}

func GetPostList2(p *models.ParamPostList) (postList []*models.ApiPostDetail, err error) {
	// 根据用户传入的排序规则去redis对于的ZSet中找帖子的ID
	ids, err := redis.GetPostIDInOrder(p)
	if err != nil {
		zap.L().Error("logic.GetPostList2 failed: ", zap.Error(err))
		return
	}

	// 根据ids去mysql数据库中查找post
	if len(ids) == 0 {
		zap.L().Warn("redis return 0 data")
		return nil, nil
	}
	posts, err := mysql.GetPostList2(ids)
	if err != nil {
		zap.L().Error("mysql.GetPostList2 failed", zap.Error(err))
		return nil, err
	}

	// 获取每篇帖子的赞成票
	voteNum, err := redis.GetPostVoteNum(ids)
	if err != nil {
		zap.L().Error("logic.GetPostList2 failed: ", zap.Error(err))
		return nil, err
	}

	postList = make([]*models.ApiPostDetail, 0, len(ids))

	for idx, post := range posts {
		// 查询作者姓名
		username, err := mysql.GetUsernameByID(post.AuthorID)
		if err != nil {
			zap.L().Error("logic.GetPostList failed: ", zap.Error(err))
			continue
		}
		// 查询社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("logic.GetPostList failed: ", zap.Error(err))
			continue
		}
		data := &models.ApiPostDetail{
			AuthorName:      username,
			VoteNum:         voteNum[idx],
			Post:            post,
			CommunityDetail: community,
		}
		postList = append(postList, data)
	}
	return postList, err
}

func GetCommunityPostList(p *models.ParamPostList) (postList []*models.ApiPostDetail, err error) {
	// 根据用户传入的排序规则去redis对于的ZSet中找帖子的ID
	ids, err := redis.GetCommunityPostIDInOrder(p)
	if err != nil {
		zap.L().Error("logic.GetPostList2 failed: ", zap.Error(err))
		return
	}

	// 根据ids去mysql数据库中查找post
	if len(ids) == 0 {
		zap.L().Warn("redis return 0 data")
		return nil, nil
	}
	posts, err := mysql.GetPostList2(ids)
	if err != nil {
		zap.L().Error("mysql.GetPostList2 failed", zap.Error(err))
		return nil, err
	}

	// 获取每篇帖子的赞成票
	voteNum, err := redis.GetPostVoteNum(ids)
	if err != nil {
		zap.L().Error("logic.GetPostList2 failed: ", zap.Error(err))
		return nil, err
	}

	postList = make([]*models.ApiPostDetail, 0, len(ids))

	for idx, post := range posts {
		// 查询作者姓名
		username, err := mysql.GetUsernameByID(post.AuthorID)
		if err != nil {
			zap.L().Error("logic.GetPostList failed: ", zap.Error(err))
			continue
		}
		// 查询社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("logic.GetPostList failed: ", zap.Error(err))
			continue
		}
		data := &models.ApiPostDetail{
			AuthorName:      username,
			VoteNum:         voteNum[idx],
			Post:            post,
			CommunityDetail: community,
		}
		postList = append(postList, data)
	}
	return postList, err
}

// 将GetPostList2和GetCommunityPostList逻辑合并，没有传入社区ID时查询全部，传入社区ID时在社区里面查询
func GetPostListMerge(p *models.ParamPostList) (postList []*models.ApiPostDetail, err error) {
	if p.CommunityID == 0 {
		postList, err = GetPostList2(p)
		if err != nil {
			zap.L().Error("logic.GetPostListMerge failed: ", zap.Error(err))
			return nil, err
		}
	} else {
		postList, err = GetCommunityPostList(p)
		if err != nil {
			zap.L().Error("logic.GetPostListMerge failed: ", zap.Error(err))
			return nil, err
		}
	}
	return
}

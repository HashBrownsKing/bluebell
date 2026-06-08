package mysql

import (
	"bluebell/models"
	"database/sql"
	"errors"

	"go.uber.org/zap"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "SELECT community_id, community_name FROM community"
	err = db.Select(&communityList, sqlStr)
	if err != nil {
		zap.L().Error("query community list failed", zap.Error(err))
		return nil, err
	}
	return communityList, nil
}

func CheckCommunityExists(communityID int64) (exists bool, err error) {
	sqlStr := `SELECT COUNT(*) FROM community WHERE community_id = ?`
	var count int
	err = db.Get(&count, sqlStr, communityID)
	if err != nil {
		zap.L().Error("mysql.CheckCommunityExists failed: ", zap.Error(err))
		return false, err
	}
	return count > 0, err
}

func GetCommunityDetailByID(communityID int64) (communityDetail *models.CommunityDetail, err error) {
	communityDetail = new(models.CommunityDetail)
	sqlStr := `SELECT community_id, community_name, introduction, create_time FROM community WHERE community_id = ?`
	err = db.Get(communityDetail, sqlStr, communityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID failed : ", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorInvalidID
		}
		return
	}
	return
}

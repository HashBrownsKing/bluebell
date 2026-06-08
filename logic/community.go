package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	return mysql.GetCommunityList()
}

func GetCommunityDetailByID(communityID int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(communityID)
}

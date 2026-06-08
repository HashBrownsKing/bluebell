package mysql

import "errors"

var (
	ErrorUserExist         = errors.New("用户已存在")
	ErrorUserNotExist      = errors.New("用户不存在")
	ErrorInvalidPassword   = errors.New("密码错误")
	ErrorInvalidID         = errors.New("ID不合法")
	ErrorCommunityNotExist = errors.New("社区不存在")
	ErrorPostNotExists     = errors.New("帖子不存在")
)

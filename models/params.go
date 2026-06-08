package models

type ParamSignUp struct {
	Username   string `json:"username" binding:"required"` // binding为validator的校验tag
	Password   string `json:"password" binding:"required"` // required表示该字段为必填字段
	RePassword string `json:"re_password"`                 // eqfield表示该字段要跟哪个字段相等
}

type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ParamPostVote struct {
	PostID   int64 `json:"post_id,string" binding:"required"`
	VoteInfo int8  `json:"vote_info,string" binding:"oneof=1 0 -1"` // oneof=x y z表示只能从oneof后面的数值中取值
}

const (
	OrderTime  = "time"
	OrderScore = "score"
)

type ParamPostList struct {
	CommunityID int64  `json:"community_id" form:"community_id"`
	Page        int64  `json:"page"  form:"page"  binding:"min=1"`        // 页码最小为1
	Size        int64  `json:"size"  form:"size"  binding:"min=1,max=50"` // 限制每页数量
	Order       string `json:"order" form:"order" binding:"oneof=time score"`
}

// type ParamCommunityPostList struct {
// 	ParamPostList
// 	CommunityID int64 `json:"community_id" form:"community_id"`
// }

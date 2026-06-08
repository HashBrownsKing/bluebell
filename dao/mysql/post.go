package mysql

import (
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func CreatePost(post *models.Post) (err error) {
	// 雪花算法生成帖子id
	post.ID = snowflake.GetID()

	sqlStr := `INSERT INTO post (post_id, community_id, author_id, title, content) VALUES (?, ?, ?, ?, ?)`
	_, err = db.Exec(sqlStr, post.ID, post.CommunityID, post.AuthorID, post.Title, post.Content)
	if err != nil {
		zap.L().Error("create post failed", zap.Error(err), zap.Any("post", post))
		return
	}
	postIDStr := strconv.FormatInt(post.ID, 10)
	communityIDStr := strconv.FormatInt(post.CommunityID, 10)
	err = redis.CreatePost(postIDStr, communityIDStr)
	if err != nil {
		delSQL := `DELETE FROM post WHERE post_id = ?`
		_, delErr := db.Exec(delSQL, post.ID)
		if delErr != nil {
			zap.L().Error("rollback mysql post failed", zap.Error(delErr), zap.Int64("post_id", post.ID))
		}
		zap.L().Error("create post redis failed", zap.Error(err), zap.Int64("post_id", post.ID))
		return err
	}
	return
}

func GetPostByID(post_id int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `SELECT post_id, author_id, community_id, title, content, create_time FROM post WHERE post_id = ?`
	err = db.Get(post, sqlStr, post_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorPostNotExists
		}
		zap.L().Error("mysql.GetPostByID failed : ", zap.Error(err))
		return nil, err
	}
	return post, nil
}

func GetPostList(page, size int64) (posts []*models.Post, err error) {
	posts = make([]*models.Post, 0, size)
	sqlStr := `SELECT post_id, author_id, community_id, title, content, create_time FROM post ORDER BY create_time DESC LIMIT ?, ?`
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}

func GetPostList2(ids []string) (posts []*models.Post, err error) {
	posts = make([]*models.Post, 0, len(ids))
	sqlStr := `SELECT 
		post_id, author_id, community_id, title, content, create_time 
		FROM post WHERE post_id IN (?)
		ORDER BY FIND_IN_SET(post_id, ?)`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		zap.L().Error("sqlx.In failed", zap.Error(err))
		return nil, err
	}

	query = db.Rebind(query)

	err = db.Select(&posts, query, args...)
	if err != nil {
		zap.L().Error("db.Select failed", zap.Error(err))
		return nil, err
	}

	return posts, nil
}

package controller

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Model struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type User struct {
	Id            int64    `json:"id,omitempty"`
	Name          string   `json:"name,omitempty"`
	Password      string   `json:"password,omitempty"`
	Follows       []Follow `json:"follows,omitempty" gorm:"foreignKey:FollowerId"`
	Followers     []Follow `json:"followers,omitempty" gorm:"foreignKey:FolloweeId"`
	FollowCount   int64    `json:"follow_count,omitempty" gorm:"-"`
	FollowerCount int64    `json:"follower_count,omitempty" gorm:"-"`
	IsFollow      bool     `json:"is_follow,omitempty" gorm:"-"`
	Videos        []Video  `json:"videos,omitempty" gorm:"foreignKey:AuthorId"`
	Model
}

type Video struct {
	Id            int64  `json:"id,omitempty"`
	AuthorId      int64  `json:"author_id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title,omitempty"`
	Model
}

type Follow struct {
	Id         int64 `json:"id,omitempty" gorm:"primaryKey"`
	FolloweeId int64 `json:"followee_id,omitempty"`
	FollowerId int64 `json:"follower_id,omitempty"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type Message struct {
	Id         int64  `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}

type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

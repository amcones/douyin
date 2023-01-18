package models

type UserInfo struct {
	ID            int    `json:"id"`
	Name          string `gorm:"type:varchar(32) not null;uniqueIndex" json:"name"`
	Password      string `gorm:"type:varchar(255) not null;"`
	FollowCount   int    `json:"follow-count"`
	FollowerCount int    `json:"follower-count"`
	IsFollow      bool   `json:"is-follow"`
}

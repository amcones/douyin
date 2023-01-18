package models

type UserInfo struct {
	Id            int    `gorm:"PrimaryKey;AutoIncrement;Unique" json:"id"`
	Name          string `gorm:"type:varchar(255) not null;" json:"name"`
	FollowCount   int    `json:"follow-count"`
	FollowerCount int    `json:"follower-count"`
	IsFollow      bool   `json:"is-follow"`
}

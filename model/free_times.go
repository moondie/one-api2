package model

import "errors"

type FreeTimes struct {
	Id         int   `json:"id"`
	UserId     int   `json:"user_id" gorm:"index"`
	ChangeTime int64 `json:"change_time" gorm:"bigint"`
	Times      int   `json:"times"`
}

func (times *FreeTimes) Insert() error {
	var err error
	err = DB.Create(times).Error
	return err
}

func (times *FreeTimes) Update() error {
	var err error
	err = DB.Model(times).Select("name", "change_time", "times").Updates(times).Error
	return err
}

func GetTimesByUserId(userId int) (*FreeTimes, error) {
	if userId == 0 {
		return nil, errors.New("userId 为空！")
	}
	times := FreeTimes{UserId: userId}
	var err error = nil
	err = DB.First(&times, "user_id = ?", userId).Error
	return &times, err
}

package model

import "gorm.io/gorm"

type Keywords struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

//SearchKey 存在返回true 否则返回false
func SearchKey(keyword string) (re bool, err error) {
	var keywords Keywords
	err = DB.Where("? LIKE CONCAT('%', name, '%')", keyword).First(&keywords).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username  string
	Password  string
	Telephone string
	Email     string
	Orders    []Order
}

func ExistUser(telephone, password string) (bool, error) {
	var user User
	err := db.Where("telephone = ? AND password = ?", telephone, password).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if user.ID > 0 {
		return true, nil
	}
	return false, nil
}

func ExistUserByTelephone(telephone string) (bool, error) {
	var user User
	err := db.Where("telephone = ?", telephone).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if user.ID > 0 {
		return true, nil
	}
	return false, nil
}

package models

import (
	"onlineShopping/pkg/setting"
	"strconv"

	"github.com/jinzhu/gorm"
)

type Order struct {
	gorm.Model
	UserID      int64
	ProductID   int64
	OrderStatus int
	Products    []Product `gorm:"many2many:order_products;"`
}

const (
	OrderWait = iota
	OrderSuccess
	OrderFailed
)

func ExistOrderByID(id string) (bool, error) {
	var order Order
	err := db.Where("id = ? ", id).First(&order).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if order.ID > 0 {
		return true, nil
	}
	return false, nil
}

func GetOrdersByPageNum(page string) ([]*Order, error) {
	var orders []*Order
	var offsetNum = 0

	pageI, err := strconv.Atoi(page)
	if err != nil {
		return nil, err
	}
	if pageI > 0 {
		offsetNum = (pageI - 1) * setting.TotalConfig.App.OrderPageSize
	} else {

	}

	// err = db.Offset(offsetNum).Limit(setting.TotalConfig.App.OrderPageSize).Preload("products").Find(&orders).Error
	err = db.Offset(offsetNum).Limit(setting.TotalConfig.App.OrderPageSize).Find(&orders).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return orders, nil
}

func GetOrdersTotalNum() (page int, err error) {
	var count int
	err = db.Model(&Order{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

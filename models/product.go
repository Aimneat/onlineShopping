package models

import (
	"onlineShopping/pkg/setting"
	"strconv"

	"github.com/jinzhu/gorm"
)

type Product struct {
	// Model
	gorm.Model
	ProductName  string `json:"ProductName" `
	ProductNum   int64  `json:"ProductNum" `
	ProductImage string `json:"ProductImage" `
	ProductUrl   string `json:"ProductUrl" `
	Price        int64
	Order        []Order
}

func ExistProductByID(id string) (bool, error) {
	var product Product
	err := db.Where("id = ? ", id).First(&product).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if product.ID > 0 {
		return true, nil
	}
	return false, nil
}

func GetPrdouctsByPageNum(page string) ([]*Product, error) {
	var products []*Product
	var offsetNum = 0

	pageI, err := strconv.Atoi(page)
	if err != nil {
		return nil, err
	}
	if pageI > 0 {
		offsetNum = (pageI - 1) * setting.TotalConfig.App.ProductPageSize
	} else {

	}

	err = db.Offset(offsetNum).Limit(setting.TotalConfig.App.ProductPageSize).Find(&products).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return products, nil
}

func GetPrdouctsTotalNum() (page int, err error) {
	var count int
	err = db.Model(&Product{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func UpdateProductNum(id, Num int) (ok bool, err error) {
	err = db.Model(&Product{}).Where("id = ?", id).Update("ProductNum", Num).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

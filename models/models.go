package models

import (
	"fmt"
	"log"
	"onlineShopping/pkg/setting"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const timeFormat = "2006-01-02 15:04:05"
const timezone = "Asia/Shanghai"

var db *gorm.DB

type Model struct {
	ID         int       `gorm:"primary_key" json:"id"`
	CreatedOn  time.Time `json:"created_on" gorm:"type:timestamp"`
	ModifiedOn time.Time `json:"modified_on" gorm:"type:timestamp"`
	DeletedOn  time.Time `json:"deleted_on" gorm:"type:timestamp"`
}

func Setup() {
	var err error

	driverName := setting.TotalConfig.Datasource.DriverName
	db, err = gorm.Open(driverName, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		setting.TotalConfig.Datasource.Username,
		setting.TotalConfig.Datasource.Password,
		setting.TotalConfig.Datasource.Host,
		setting.TotalConfig.Datasource.Port,
		setting.TotalConfig.Datasource.Database,
		setting.TotalConfig.Datasource.Charset,
		// setting.TotalConfig.Datasource.Loc,
	))
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}
}

func GetDB() *gorm.DB {
	return db
}

package v1

import (
	"fmt"
	"net/http"
	"onlineShopping/models"
	"onlineShopping/pkg/app"
	"onlineShopping/pkg/e"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type IOrder interface {
	Create(c *gin.Context)
	Updata(c *gin.Context)
	Delete(c *gin.Context)
	ShowAll(c *gin.Context)
	ShowByKey(c *gin.Context)
}

type OrderManager struct {
	DB *gorm.DB
}

type OrderForm struct {
	UserID      int64
	ProductId   int64
	OrderStatus int
	Products    []models.Product
}

func (o OrderManager) Create(c *gin.Context) {
	var orderForm OrderForm
	err := c.ShouldBindJSON(&orderForm)
	if err != nil {
		app.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	Order := models.Order{
		UserID:      orderForm.UserID,
		ProductID:   orderForm.ProductId,
		OrderStatus: orderForm.OrderStatus,
		Products:    orderForm.Products,
	}

	db := models.GetDB()
	tx := db.Begin()

	for _, orderProduct := range orderForm.Products {
		exists, err := models.ExistProductByID(strconv.Itoa(int(orderProduct.ID)))
		if err != nil {
			app.Response(c, http.StatusInternalServerError, e.ERROR_CHECK_EXIST_PRODUCT_FAIL, nil)
			return
		}
		if !exists {
			app.Response(c, http.StatusOK, e.ERROR_NOT_EXIST_PRODUCT, nil)
			return
		}

		var pd models.Product

		err = tx.Model(&models.Product{}).Where("id = ?", orderProduct.ID).First(&pd).Error
		if err != nil {
			app.Response(c, http.StatusInternalServerError, e.ERROR_CHECK_EXIST_PRODUCT_FAIL, nil)
			return
		}

		stock := pd.ProductNum - orderProduct.ProductNum
		fmt.Println(pd.ProductNum, "-", orderProduct.ProductNum, "=", stock)
		if stock >= 0 {
			err = tx.Model(&models.Product{}).Where("id = ?", orderProduct.ID).Update("ProductNum", stock).Error
			if err != nil {
				app.Response(c, http.StatusInternalServerError, e.ERROR_UPDATA_PRODUCT_FAIL, nil)
				return
			}

		} else {
			app.Response(c, http.StatusBadRequest, e.ERROR_OUT_OF_STOCK, nil)
			tx.Rollback()
			return
		}

	}
	tx.Commit()

	// if err := o.DB.Omit("Products").Create(&Order).Error; err != nil { //BUG：不会创建从表，关联表对应信息。
	if err := o.DB.Debug().Set("gorm:association_autoupdate", false).Set("gorm:association_autoupdate", false).Create(&Order).Error; err != nil { //BUG：不会创建从表，关联表对应信息。
		// if err := o.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&Order).Error; err != nil { //创建order后会修改关联表的值,导致超卖！
		app.Response(c, http.StatusInternalServerError, e.ERROR_ADD_ORDER_FAIL, nil)
		return
	}

	app.Response(c, http.StatusOK, e.SUCCESS, nil)
}

func (o OrderManager) Updata(c *gin.Context) {
	id := c.Param("id")

	exists, err := models.ExistOrderByID(id)
	if err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ORDER_FAIL, nil)
		return
	}
	if !exists {
		app.Response(c, http.StatusOK, e.ERROR_NOT_EXIST_ORDER, nil)
		return
	}

	var orderForm OrderForm
	err = c.ShouldBindJSON(&orderForm)
	if err != nil {
		app.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	Order := models.Order{
		UserID:      orderForm.UserID,
		ProductID:   orderForm.ProductId,
		OrderStatus: orderForm.OrderStatus,
		Products:    orderForm.Products,
	}

	if err = o.DB.Model(&models.Order{}).Where("id = ?", id).Update(&Order).Error; err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_UPDATA_ORDER_FAIL, nil)
		return
	}

	app.Response(c, http.StatusOK, e.SUCCESS, gin.H{"order": Order})
}

func (o OrderManager) Delete(c *gin.Context) {
	id := c.Param("id")

	var Order models.Order

	exists, err := models.ExistOrderByID(id)
	if err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ORDER_FAIL, nil)
		return
	}
	if !exists {
		app.Response(c, http.StatusOK, e.ERROR_NOT_EXIST_ORDER, nil)
		return
	}

	if err := o.DB.Where("id = ?", id).Delete(&Order).Error; err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_ADD_ORDER_FAIL, nil)
		return
	}

	app.Response(c, http.StatusOK, e.SUCCESS, gin.H{"order": Order})
}

func (o OrderManager) ShowAll(c *gin.Context) {
	page := c.Query("page")

	Orders, err := models.GetOrdersByPageNum(page)
	if err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_GET_ORDERS_FAIL, nil)
		return
	}

	total, err := models.GetOrdersTotalNum()
	if err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_GET_ORDERS_NUMBER_FAIL, nil)
		return
	}

	app.Response(c, http.StatusOK, e.SUCCESS, gin.H{"Order": Orders, "total": total})
}

func (o OrderManager) ShowByKey(c *gin.Context) {
	id := c.Param("id")

	var Order models.Order
	var Product models.Product

	exists, err := models.ExistOrderByID(id)
	if err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ORDER_FAIL, nil)
		return
	}
	if !exists {
		app.Response(c, http.StatusOK, e.ERROR_NOT_EXIST_ORDER, nil)
		return
	}

	if err := o.DB.Where("id = ?", id).First(&Order).Error; err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_GET_ORDER_FAIL, nil)
		return
	}

	exists, err = models.ExistProductByID(strconv.Itoa(int(Order.ProductID)))
	if err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_CHECK_EXIST_PRODUCT_FAIL, nil)
		return
	}
	if !exists {
		app.Response(c, http.StatusOK, e.ERROR_NOT_EXIST_PRODUCT_IN_ORDER, nil)
		return
	}
	if err := o.DB.Where("id = ?", Order.ProductID).First(&Product).Error; err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_GET_PRODUCT_FAIL, nil)
		return
	}

	app.Response(c, http.StatusOK, e.SUCCESS, gin.H{"order": Order, "product": Product})
}

func NewOrderManager() IOrder {
	db := models.GetDB()
	db.AutoMigrate(&models.Order{})
	return OrderManager{DB: db}
}

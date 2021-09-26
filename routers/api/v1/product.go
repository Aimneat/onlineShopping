package v1

import (
	"net/http"
	"onlineShopping/models"
	"onlineShopping/pkg/app"
	"onlineShopping/pkg/e"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type IProduct interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	ShowAll(c *gin.Context)
	ShowByKey(c *gin.Context)
}

type ProductManager struct {
	DB *gorm.DB
}

type ProductForm struct {
	ProductName  string
	ProductNum   int64
	ProductImage string
	ProductUrl   string
	Price        int64
}

func (p ProductManager) Create(c *gin.Context) {
	var productForm ProductForm
	err := c.ShouldBindJSON(&productForm)
	if err != nil {
		app.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	product := models.Product{
		ProductName:  productForm.ProductName,
		ProductNum:   productForm.ProductNum,
		ProductImage: productForm.ProductImage,
		ProductUrl:   productForm.ProductUrl,
		Price:        productForm.Price,
	}

	if err := p.DB.Create(&product).Error; err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_ADD_PRODUCT_FAIL, nil)
		return
	}

	app.Response(c, http.StatusOK, e.SUCCESS, nil)
}

func (p ProductManager) Update(c *gin.Context) {
	id := c.Param("id")

	exists, err := models.ExistProductByID(id)
	if err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_CHECK_EXIST_PRODUCT_FAIL, nil)
		return
	}
	if !exists {
		app.Response(c, http.StatusOK, e.ERROR_NOT_EXIST_PRODUCT, nil)
		return
	}

	var productForm ProductForm
	err = c.ShouldBindJSON(&productForm)
	if err != nil {
		app.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	product := models.Product{
		ProductName:  productForm.ProductName,
		ProductNum:   productForm.ProductNum,
		ProductImage: productForm.ProductImage,
		ProductUrl:   productForm.ProductUrl,
		Price:        productForm.Price,
	}

	if err = p.DB.Model(&models.Product{}).Where("id = ?", id).Update(&product).Error; err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_UPDATA_PRODUCT_FAIL, nil)
		return
	}

	app.Response(c, http.StatusOK, e.SUCCESS, gin.H{"product": product})
}

func (p ProductManager) Delete(c *gin.Context) {
	id := c.Param("id")

	var product models.Product

	exists, err := models.ExistProductByID(id)
	if err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_CHECK_EXIST_PRODUCT_FAIL, nil)
		return
	}
	if !exists {
		app.Response(c, http.StatusOK, e.ERROR_NOT_EXIST_PRODUCT, nil)
		return
	}

	if err := p.DB.Where("id = ?", id).Delete(&product).Error; err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_ADD_PRODUCT_FAIL, nil)
		return
	}

	app.Response(c, http.StatusOK, e.SUCCESS, gin.H{"product": product})
}

func (p ProductManager) ShowAll(c *gin.Context) {
	// var products []*models.Product
	page := c.Query("page")

	products, err := models.GetPrdouctsByPageNum(page)
	if err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_GET_PRODUCTS_FAIL, nil)
		return
	}
	total, err := models.GetPrdouctsTotalNum()
	if err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_GET_PRODUCTS_NUMBER_FAIL, nil)
		return
	}

	app.Response(c, http.StatusOK, e.SUCCESS, gin.H{"product": products, "total": total})
}

func (p ProductManager) ShowByKey(c *gin.Context) {
	id := c.Param("id")

	var product models.Product

	exists, err := models.ExistProductByID(id)
	if err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_CHECK_EXIST_PRODUCT_FAIL, nil)
		return
	}
	if !exists {
		app.Response(c, http.StatusOK, e.ERROR_NOT_EXIST_PRODUCT, nil)
		return
	}

	if err := p.DB.Where("id = ?", id).First(&product).Error; err != nil {
		app.Response(c, http.StatusInternalServerError, e.ERROR_GET_PRODUCT_FAIL, nil)
		return
	}

	app.Response(c, http.StatusOK, e.SUCCESS, gin.H{"product": product})
}

func NewProductManager() IProduct {
	db := models.GetDB()
	db.AutoMigrate(&models.Product{})
	return ProductManager{DB: db}
}

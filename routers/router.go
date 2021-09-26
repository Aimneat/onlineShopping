package routers

import (
	_ "onlineShopping/docs"
	"onlineShopping/middleware"
	v1 "onlineShopping/routers/api/v1"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(middleware.Cors())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiv1 := r.Group("/api/v1")
	user := v1.NewUserManager()
	apiv1.POST("/register", user.Register)
	apiv1.POST("/login", user.Login)

	apiv1.Use(middleware.JWTAuth())
	{
		//验证token
		apiv1.GET("/ping", user.CheckToken)

		apiv1.GET("/user/:id", user.MyInformation)

		product := v1.NewProductManager()
		apiv1.POST("/product", product.Create)
		apiv1.GET("/product/:id", product.ShowByKey)
		apiv1.GET("/products", product.ShowAll)
		apiv1.DELETE("/product/:id", product.Delete)
		apiv1.PUT("/product/:id", product.Update)

		order := v1.NewOrderManager()
		apiv1.POST("/order", order.Create)
		apiv1.GET("/order/:id", order.ShowByKey)
		apiv1.GET("/orders", order.ShowAll)
		apiv1.DELETE("/order/:id", order.Delete)
		apiv1.PUT("/order/:id", order.Updata)
	}

	return r
}

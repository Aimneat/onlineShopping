package e

const (
	ERROR_NOT_EXIST_PRODUCT = 10001 + iota
	ERROR_CHECK_EXIST_PRODUCT_FAIL
	ERROR_ADD_PRODUCT_FAIL
	ERROR_UPDATA_PRODUCT_FAIL
	ERROR_GET_PRODUCTS_FAIL
	ERROR_GET_PRODUCTS_NUMBER_FAIL
	ERROR_GET_PRODUCT_FAIL
	ERROR_OUT_OF_STOCK

	ERROR_NOT_EXIST_ORDER
	ERROR_CHECK_EXIST_ORDER_FAIL
	ERROR_ADD_ORDER_FAIL
	ERROR_UPDATA_ORDER_FAIL
	ERROR_GET_ORDERS_FAIL
	ERROR_GET_ORDERS_NUMBER_FAIL
	ERROR_GET_ORDER_FAIL
	ERROR_NOT_EXIST_PRODUCT_IN_ORDER

	ERROR_NOT_EXIST_USER
	ERROR_ADD_EXIST_USER
	ERROR_ADD_USER_FAIL
	ERROR_GET_USER_FAIL
	ERROR_GET_USER_TELEPHONE_FAIL

	ERROR_AUTH_CHECK_TOKEN_TIMEOUT
	ERROR_GENERATE_TOKEN_FAIL

	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400
)

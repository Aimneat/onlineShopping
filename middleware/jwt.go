package middleware

import (
	"errors"
	"net/http"
	"onlineShopping/models"
	"onlineShopping/models/requset"
	"onlineShopping/pkg/app"
	"onlineShopping/pkg/e"
	"onlineShopping/pkg/setting"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	TokenInvalid = errors.New("Couldn't handle this token")
	TokenExpired = errors.New("Token is expired")
)

type JWT struct {
	JwtSecret []byte
}

func NewJWT() *JWT {
	return &JWT{JwtSecret: []byte(setting.TotalConfig.Jwt.JwtSecret)}
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code = e.SUCCESS
		var data interface{}

		j := NewJWT()
		// 这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.Request.Header.Get("x-token")
		if token == "" {
			code = e.INVALID_PARAMS
			data = "未登录或非法访问"
		} else {
			_, err := j.ParseToken(token)
			if err != nil {
				switch err {
				case TokenExpired:
					code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
					data = TokenExpired
				default:
					code = e.INVALID_PARAMS
					data = err.Error()
				}
			}
		}

		if code != e.SUCCESS {
			app.Response(c, http.StatusUnauthorized, code, data)
			c.Abort()
			return
		}
		c.Next()
	}
}

// ParseToken parsing token
func (j *JWT) ParseToken(tokenString string) (*requset.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &requset.CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return j.JwtSecret, nil
	})
	if err != nil {
		// if ve, ok := err.(*jwt.ValidationError); ok {
		// 	if ve.Errors&jwt.ValidationErrorExpired != 0 {
		// 		return nil, TokenExpired
		// 	}
		// }

		switch err.(*jwt.ValidationError).Errors {
		case jwt.ValidationErrorExpired:
			return nil, TokenExpired
		default:
			return nil, TokenInvalid
		}

	}
	if token != nil {
		if claims, ok := token.Claims.(*requset.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}

}

// GenerateToken generate tokens used for auth
func (j *JWT) GenerateToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &requset.CustomClaims{
		Telephone: user.Telephone,
		UserName:  user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "y",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.JwtSecret)
}

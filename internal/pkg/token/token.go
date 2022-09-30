package token

import (
	"strconv"
	"sync"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type config struct {
	// 签名密钥
	secret string
	// 身份标识的key
	identityKey string
	// 过期时间，单位：小时
	expire int
}

var (
	cfg  = config{secret: "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5", identityKey: "userId", expire: 2}
	once sync.Once
)

// Parse 解析token字符串
func Parse(tokenStr string) (int, error) {
	// parse token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(cfg.secret), nil
	})

	if err != nil {
		return -1, err
	}

	var identity string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		identity = claims[cfg.identityKey].(string)
	}

	ret, _ := strconv.Atoi(identity)
	return ret, nil
}

// Sign 生成jwt令牌
func Sign(identity, expire int) (string, error) {
	if expire <= 0 {
		expire = cfg.expire
	}

	claims := jwt.MapClaims{
		cfg.identityKey: strconv.Itoa(identity),
		"nbf":           time.Now().Unix(),                                        // 生效时间
		"iat":           time.Now().Unix(),                                        // 签发时间
		"exp":           time.Now().Add(time.Duration(expire) * time.Hour).Unix(), // 过期时间
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(cfg.secret))
}

func Init(secret string) {
	once.Do(func() {
		if secret != "" {
			cfg.secret = secret
		}
	})
}

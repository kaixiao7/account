package token

import (
	"strconv"
	"sync"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

// jwt 签名配置
type config struct {
	// 签名密钥
	secret string
	// 身份标识的key
	identityKey string
	typeKey     string

	// 过期时间，单位：秒
	accessExpire  int64
	refreshExpire int64
}

const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"
)

var (
	cfg = config{
		secret:        "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5",
		identityKey:   "userId",
		typeKey:       "type",
		accessExpire:  7200,    // 默认7200秒，两个小时
		refreshExpire: 2592000, // 默认30天
	}
	once sync.Once
)

// Parse 解析token字符串
// 返回 tokenType 与 identity 的值
func Parse(tokenStr string) (string, int64, error) {
	// parse token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(cfg.secret), nil
	})

	if err != nil {
		return "", -1, err
	}

	var identity string
	var tokenType string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		identity = claims[cfg.identityKey].(string)
		tokenType = claims[cfg.typeKey].(string)
	}

	ret, _ := strconv.ParseInt(identity, 10, 64)
	return tokenType, ret, nil
}

// Sign 生成jwt令牌
func Sign(tokenType string, identity int64, expire int64) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		cfg.identityKey: strconv.FormatInt(identity, 10),
		cfg.typeKey:     tokenType,
		"nbf":           now.Unix(),                                          // 生效时间
		"iat":           now.Unix(),                                          // 签发时间
		"exp":           now.Add(time.Duration(expire) * time.Second).Unix(), // 过期时间
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(cfg.secret))
}

// GenerateAccessToken 生成accessToken
func GenerateAccessToken(identity int64) (string, error) {
	return Sign(AccessTokenType, identity, cfg.accessExpire)
}

// GenerateRefreshToken 生成refreshToken
func GenerateRefreshToken(identity int64) (string, error) {
	return Sign(RefreshTokenType, identity, cfg.refreshExpire)
}

// DecodeAccessToken 解码accessToken
// 返回identity
func DecodeAccessToken(token string) (int64, error) {
	tokenType, identity, err := Parse(token)
	if err != nil {
		return 0, err
	}
	if tokenType != AccessTokenType {
		return 0, errors.New("token type is invalid, should is [access] but is " + tokenType)
	}

	return identity, nil
}

// DecodeRefreshToken 解码refreshToken
// 返回identity
func DecodeRefreshToken(token string) (int64, error) {
	tokenType, identity, err := Parse(token)
	if err != nil {
		return 0, err
	}
	if tokenType != RefreshTokenType {
		return 0, errors.New("token type is invalid, should is [refresh] but is " + tokenType)
	}

	return identity, nil
}

// Init 初始化
func Init(secret string, accessExpire, refreshExpire int64) {
	once.Do(func() {
		if secret != "" {
			cfg.secret = secret
		}

		cfg.accessExpire = accessExpire
		cfg.refreshExpire = refreshExpire
	})
}

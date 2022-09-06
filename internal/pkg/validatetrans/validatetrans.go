package validatetrans

import (
	"errors"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

func Init() error {
	zh := zh.New()
	uni = ut.New(zh, zh)

	var ok bool
	if validate, ok = binding.Validator.Engine().(*validator.Validate); !ok {
		return errors.New("获取gin.binding.Validator失败")
	}

	if trans, ok = uni.GetTranslator("zh"); !ok {
		return errors.New("获取中文翻译器失败")
	}
	// 注册翻译器
	return zh_translations.RegisterDefaultTranslations(validate, trans)
}

func Translate(err error) map[string][]string {
	var result = make(map[string][]string)

	if errors, ok := err.(validator.ValidationErrors); ok {
		// 遍历所有错误
		for _, err := range errors {
			result[err.Field()] = append(result[err.Field()], err.Translate(trans))
		}
	}

	return result
}

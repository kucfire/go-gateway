package middleware

import (
	"gatewayDemo/public"
	"reflect"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
)

//设置Translation
func TranslationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//参照：https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go

		//设置支持语言
		en := en.New()
		zh := zh.New()

		//设置国际化翻译器
		uni := ut.New(zh, zh, en)
		val := validator.New()

		//根据参数取翻译器实例
		locale := c.DefaultQuery("locale", "zh")
		trans, _ := uni.GetTranslator(locale)

		//翻译器注册到validator
		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("en_comment")
			})
			break
		default:
			zh_translations.RegisterDefaultTranslations(val, trans)
			val.RegisterTagNameFunc(func(fld reflect.StructField) string {
				return fld.Tag.Get("comment")
			})

			//自定义验证方法
			//https://github.com/go-playground/validator/blob/v9/_examples/custom-validation/main.go
			// 验证方法：dto.AdminLoginInput.UserName
			val.RegisterValidation("vaild_username", func(fl validator.FieldLevel) bool {
				return fl.Field().String() == "admin"
			})
			// 验证方法：dto.ServiceAddHTTPInput.ServiceName
			val.RegisterValidation("valid_service_name", func(fl validator.FieldLevel) bool {
				// return fl.Field().String() == "admin"
				regexp.
			})

			//自定义验证器
			//https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go
			// 验证器：dto.AdminLoginInput.UserName
			val.RegisterTranslation("vaild_username", trans, func(ut ut.Translator) error {
				return ut.Add("vaild_username", "{0} 填写不正确哦", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("vaild_username", fe.Field())
				return t
			})
			// 验证器：dto.ServiceAddHTTPInput.ServiceName
			val.RegisterTranslation("vaild_username", trans, func(ut ut.Translator) error {
				return ut.Add("vaild_username", "{0} 填写不正确哦", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("vaild_username", fe.Field())
				return t
			})
			break
		}
		c.Set(public.TranslatorKey, trans)
		c.Set(public.ValidatorKey, val)
		c.Next()
	}
}

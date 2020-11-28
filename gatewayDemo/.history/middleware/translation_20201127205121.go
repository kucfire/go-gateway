package middleware

import (
	"gatewayDemo/public"
	"reflect"
	"regexp"
	"strings"

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
				match, _ := regexp.Match(`^[a-zA-Z0-9_]{6,128}$`, []byte(fl.Field().String()))
				return match
			})
			// 验证方法：dto.ServiceAddHTTPInput.Rule
			val.RegisterValidation("vaild_rule", func(fl validator.FieldLevel) bool {

				// return fl.Field().String() == "admin"
				match, _ := regexp.Match(`^\S+$`, []byte(fl.Field().String()))
				return match
			})
			// 验证方法：dto.ServiceAddHTTPInput.URLRewrite
			val.RegisterValidation("valid_url_rewrite", func(fl validator.FieldLevel) bool {
				// return fl.Field().String() == "admin"
				if fl.Field().String() == "" {
					return true
				}
				for _, ms := range strings.Split(fl.Field().String(), "\n") {
					if len(strings.Split(ms, " ")) != 2 {
						return false
					}
				}
				return true
			})
			// 验证方法：dto.ServiceAddHTTPInput.HeaderTransfor
			val.RegisterValidation("valid_header_transfor", func(fl validator.FieldLevel) bool {
				// return fl.Field().String() == "admin"
				if fl.Field().String() == "" {
					return true
				}
				for _, ms := range strings.Split(fl.Field().String(), "\n") {
					if len(strings.Split(ms, " ")) != 3 {
						return false
					}
				}
				return true
			})
			// 验证方法：dto.ServiceAddHTTPInput.IPList
			val.RegisterValidation("valid_ip_list", func(fl validator.FieldLevel) bool {
				// return fl.Field().String() == "admin"
				for _, ms := range strings.Split(fl.Field().String(), "\n") {
					if matched, _ := regexp.Match(`^\S+\:\d+$`, []byte(ms)); !matched {
						return false
					}
				}
				return true
			})
			// 验证方法：dto.ServiceAddHTTPInput.WeightList
			val.RegisterValidation("valid_weight_list", func(fl validator.FieldLevel) bool {
				// return fl.Field().String() == "admin"
				for _, ms := range strings.Split(fl.Field().String(), "\n") {
					if matched, _ := regexp.Match(`^\d+$`, []byte(ms)); !matched {
						return false
					}
				}
				return true
			})
			// 验证方法：dto.AppInfo.Secret
			val.RegisterValidation("valid_secret", func(fl validator.FieldLevel) bool {
				return fl.Field().String() == ""
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
			val.RegisterTranslation("valid_service_name", trans, func(ut ut.Translator) error {
				return ut.Add("valid_service_name", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_service_name", fe.Field())
				return t
			})
			// 验证器：dto.ServiceAddHTTPInput.Rule
			val.RegisterTranslation("vaild_rule", trans, func(ut ut.Translator) error {
				return ut.Add("vaild_rule", "{0} 必须是非空字符", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("vaild_rule", fe.Field())
				return t
			})
			// 验证器：dto.ServiceAddHTTPInput.URLRewrite
			val.RegisterTranslation("valid_url_rewrite", trans, func(ut ut.Translator) error {
				return ut.Add("valid_url_rewrite", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_url_rewrite", fe.Field())
				return t
			})
			// 验证器：dto.ServiceAddHTTPInput.HeaderTransfor
			val.RegisterTranslation("valid_header_transfor", trans, func(ut ut.Translator) error {
				return ut.Add("valid_header_transfor", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_header_transfor", fe.Field())
				return t
			})
			// 验证器：dto.ServiceAddHTTPInput.IPList
			val.RegisterTranslation("valid_ip_list", trans, func(ut ut.Translator) error {
				return ut.Add("valid_ip_list", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_ip_list", fe.Field())
				return t
			})
			// 验证器：dto.ServiceAddHTTPInput.WeightList
			val.RegisterTranslation("valid_weight_list", trans, func(ut ut.Translator) error {
				return ut.Add("valid_weight_list", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_weight_list", fe.Field())
				return t
			})
			// 验证器：dto.AppInfo.Secret
			val.RegisterTranslation("valid_secret", trans, func(ut ut.Translator) error {
				return ut.Add("valid_secret", "{0} 不符合输入格式", true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("valid_secret", fe.Field())
				return t
			})

			break
		}
		c.Set(public.TranslatorKey, trans)
		c.Set(public.ValidatorKey, val)
		c.Next()
	}
}

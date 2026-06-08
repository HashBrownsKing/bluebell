package controllers

import (
	"fmt"
	"reflect"
	"strings"

	"bluebell/models"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// 定义一个全局翻译器T
var trans ut.Translator

// InitTrans 初始化翻译器
func InitTrans(locale string) (err error) {
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// 注册一个获取json tag的自定义方法，使得返回错误的字段名为前端json的字段名
		v.RegisterTagNameFunc(func(fld reflect.StructField) string { // fld reflect.StructField：通过反射拿到结构体字段的所有信息（字段名、标签等）
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0] // 处理带额外参数的 json 标签（如 json:"username,omitempty"），按逗号切割成 2 段,[0]只取第一个字段，忽略后续参数
			if name == "-" {                                       // 忽略隐藏字段
				return ""
			}
			return name
		})

		// 为SignUpParam注册自定义校验方法
		v.RegisterStructValidation(SignUpParamStructLevelValidation, models.ParamSignUp{}) // 仅仅注册前一个函数无法完全将结构体字段转化为前端json字段

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}

// 将字段转化为前端字段的json字段后，还需要定义一个去除结构体名称前缀的方法
func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

// SignUpParamStructLevelValidation 自定义SignUpParam结构体校验函数
// 函数签名：固定格式！validator库规定，参数必须是 validator.StructLevel
func SignUpParamStructLevelValidation(sl validator.StructLevel) {
	// 1. 获取当前正在校验的结构体对象
	// sl.Current()：获取校验的目标结构体
	// .Interface()：转为空接口
	// .(models.ParamSignUp)：类型断言，转换成我们的注册参数结构体
	// 最终 su 就是前端传过来的 用户名/密码/重复密码 对象
	su := sl.Current().Interface().(models.ParamSignUp)

	// 2. 核心业务逻辑：判断 密码 和 重复密码 不相等
	if su.Password != su.RePassword {
		// 3. 手动抛出校验错误（最关键！控制错误提示内容）
		// ReportError 是 validator 提供的固定方法，用来报告校验失败
		sl.ReportError(
			su.RePassword, // 出错的字段：哪个字段错了（重复密码）
			"re_password", // 前端显示的字段名（JSON标签名，小写）
			"RePassword",  // 结构体原生字段名（Go内部用）
			"eqfield",     // 校验规则名称（复用翻译文案）
			"password",    // 对比的目标字段名（直接写前端JSON名！）
		)
	}
}

package gin

import "net/http"

/**
自定义gin server可选配置
*/

type HttpResponseFiled struct {
	Name      string `json:"name"`
	FieldType string `json:"type"`
}

type megaGinServerOptions struct {
	responseFields      []*HttpResponseFiled
	responseSuccessCode int
	responseFailCode    int
}

func defaultMegaGinServerOptions() *megaGinServerOptions {
	return &megaGinServerOptions{
		responseFields: []*HttpResponseFiled{
			{Name: "message", FieldType: "string"},
			{Name: "code", FieldType: "int"},
			{Name: "success", FieldType: "bool"},
		},
		responseSuccessCode: http.StatusOK,
		responseFailCode:    http.StatusBadRequest,
	}
}

type MegaGinServerOption interface {
	apply(opt *megaGinServerOptions)
}
type funcOption struct {
	f func(*megaGinServerOptions)
}
func (p *funcOption) apply(opt *megaGinServerOptions) {
	p.f(opt)
}
func NewFuncOption(f func(*megaGinServerOptions)) *funcOption {
	return &funcOption{
		f: f,
	}
}

//**********with options*******************
func WithResponseFields(fields []*HttpResponseFiled) MegaGinServerOption {
	return NewFuncOption(func(options *megaGinServerOptions) {
		options.responseFields = fields
	})
}

func WithResponseSuccessCode(code int) MegaGinServerOption{
	return NewFuncOption(func(options *megaGinServerOptions) {
		options.responseSuccessCode = code
	})
}

func WithResponseFailCode(code int) MegaGinServerOption  {
	return NewFuncOption(func(options *megaGinServerOptions) {
		options.responseFailCode = code
	})
}

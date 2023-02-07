package gin
/**
自定义gin server可选配置
 */

type HttpResponseFiled struct {
	Name      string `json:"name"`
	FieldType string `json:"type"`
}

type megaGinServerOptions struct {
	responseFields []*HttpResponseFiled

}
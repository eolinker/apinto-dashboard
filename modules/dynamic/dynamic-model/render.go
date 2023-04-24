package dynamic_model

type DynamicModuleRender struct {
	*DynamicBasicInfo
	Render map[string]interface{} `json:"render"`
}

package restful

type ResponseData[D any] struct {
	Code      int    `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	RequestId string `json:"requestId,omitempty"`
	Success   bool   `json:"success,omitempty"`
	Data      *D     `json:"data"`
}
type Response struct {
	Code      int    `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	RequestId string `json:"requestId,omitempty"`
	Success   bool   `json:"success,omitempty"`
}

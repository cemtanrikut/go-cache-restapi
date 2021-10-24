package helper

type Resp struct {
	Success bool              `json:"success"`
	Method  string            `json:"method"`
	Action  string            `json:"action,omitempty"`
	Data    interface{}       `json:"data"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// NewSuccess returns `*Success` instance.
func NewResponse() *Resp {
	return &Resp{}
}

func ResponseSuccess(success bool, action string, data interface{}, method string) (httpResp *Resp) {

	httpResp = NewResponse()
	httpResp.Success = success
	httpResp.Action = action
	httpResp.Data = data
	httpResp.Method = method
	return httpResp

}

package standard

type Response struct {
	Message string      `json:"message,omitempty"`
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
}

var (
	ParamsError = &Response{Message: "Params Error", Success: false, Code: 40000}
	AuthError   = &Response{Message: "Auth Error", Success: false, Code: 40001}
)

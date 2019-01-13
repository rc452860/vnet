package api

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	SUCCESS = "success"
	FAILURE = "failure"
	ERROR   = "error"
)

func Success(data interface{}, param ...string) *Response {
	r := &Response{
		Code: SUCCESS,
		Data: data,
	}
	if len(param) > 0 {
		r.Message = param[0]
	}
	return r
}

func Failure(message string) *Response {
	r := &Response{
		Code:    FAILURE,
		Message: message,
	}
	return r
}

func Error(err error) *Response {
	r := &Response{
		Code:    ERROR,
		Message: err.Error(),
	}
	return r
}

type SystemInfo struct {
	CPUUsage  int         `json:"cpu_usage,omitempty"`
	MemUsage  int         `json:"mem_usage,omitempty"`
	DiskUsage int         `json:"disk_usage,omitempty"`
	Network   NetworkInfo `json:"network,omitempty"`
}

type NetworkInfo struct {
	Up   uint64 `json:"up,omitempty"`
	Down uint64 `json:"down,omitempty"`
}

func NewSystemInfo(cpu, mem, disk int, up, down uint64) SystemInfo {
	sysinfo := SystemInfo{
		CPUUsage:  cpu,
		MemUsage:  mem,
		DiskUsage: disk,
		Network: NetworkInfo{
			Up:   up,
			Down: down,
		},
	}
	return sysinfo
}

package domain

// Author:Boyn
// Date:2020/9/14

type ServiceInputAdapter interface {
	ExtractServiceAndRule() (*ServiceInfo, ServiceRule)
}

type HttpServiceJsonInput struct {
	ServiceType    int    `gorm:"column:service_type" json:"service_type"`
	ServiceName    string `gorm:"column:service_name" json:"service_name"`
	ServicePort    int    `gorm:"column:service_port" json:"service_port"`
	ServiceDesc    string `gorm:"column:service_desc" json:"service_desc"`
	RoundType      int    `json:"round_type" gorm:"column:round_type"`
	Prefix         string `json:"prefix" gorm:"column:prefix"` // URL前缀
	NeedHttps      int    `json:"need_https" gorm:"column:need_https"`
	NeedWebsocket  int    `json:"need_websocket" gorm:"column:need_websocket"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor"`
}

func (h *HttpServiceJsonInput) ExtractServiceAndRule() (*ServiceInfo, ServiceRule) {
	input := &ServiceInfo{
		ServiceType: h.ServiceType,
		ServiceName: h.ServiceName,
		ServicePort: h.ServicePort,
		ServiceDesc: h.ServiceDesc,
		RoundType:   h.RoundType,
	}
	rule := &HttpRule{
		Prefix:         h.Prefix,
		NeedHttps:      h.NeedHttps,
		NeedWebsocket:  h.NeedWebsocket,
		HeaderTransfor: h.HeaderTransfor,
	}
	return input, rule
}

type GrpcServiceJsonInput struct {
	ServiceType    int    `gorm:"column:service_type" json:"service_type"`
	ServiceName    string `gorm:"column:service_name" json:"service_name"`
	ServicePort    int    `gorm:"column:service_port" json:"service_port"`
	ServiceDesc    string `gorm:"column:service_desc" json:"service_desc"`
	RoundType      int    `json:"round_type" gorm:"column:round_type"`
	Port           int    `json:"port" gorm:"column:port"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor"`
}

func (g *GrpcServiceJsonInput) ExtractServiceAndRule() (*ServiceInfo, ServiceRule) {
	input := &ServiceInfo{
		ServiceType: g.ServiceType,
		ServiceName: g.ServiceName,
		ServicePort: g.ServicePort,
		ServiceDesc: g.ServiceDesc,
		RoundType:   g.RoundType,
	}
	rule := &GrpcRule{
		Port:           g.Port,
		HeaderTransfor: g.HeaderTransfor,
	}
	return input, rule
}

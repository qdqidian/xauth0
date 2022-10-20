package xauth0

import "net/url"

type Req struct {
	Schema string `json:"schema"`

	Domain string `json:"domain"`

	Port string `json:"port"`

	Path string `json:"path"`

	Form url.Values

	Body string
}

type Res struct {
	Pass bool `json:"pass"`

	OperationCheckStrategy int32 `json:"operationCheckStrategy"`

	AllowCrossDomain bool `json:"allowCrossDomain"`

	CrossDomains []string `json:"crossDomains"`

	// 字段级数据权限定义
	FieldPermissionDefine string `json:"FieldPermissionDefine"`

	User *User `json:"user"`
}

type User struct {
	Id string `json:"id"`

	Name string `json:"name"`

	Code string `json:"code"`

	Meta map[string]string `json:"meta"`
}

type ErrorRes struct {
	Code int `json:"code"`

	Reason string `json:"reason"`

	Message string `json:"message"`

	Metadata interface{} `json:"metadata"`
}

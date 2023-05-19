package model

type User struct {
	// 结构体字段与json字符串的key对应
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}

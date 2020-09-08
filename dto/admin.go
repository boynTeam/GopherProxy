package dto

import "time"

// Author:Boyn
// Date:2020/9/8

type AdminInput struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type AdminUserSession struct {
	Id        uint      `json:"id"`
	UserName  string    `json:"user_name"`
	LoginTime time.Time `json:"login_time"`
}

package pkg

import "github.com/gorilla/sessions"

// Author:Boyn
// Date:2020/9/8
const (
	UserCookieName = "GATEWAY_USER_INFO"
)

var CookieSession = sessions.NewCookieStore([]byte("123456"))

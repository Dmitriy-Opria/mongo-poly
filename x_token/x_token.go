package x_token

import (
	"golang.org/x/net/xsrftoken"
)

var (
	key      = `key_{key}_key`
	userID   = "mongo_check"
	actionID = "check_field"
)

func GetToken() (token string) {
	token = xsrftoken.Generate(key, userID, actionID)

	return
}

func CheckValid(token string) bool {

	return xsrftoken.Valid(token, key, userID, actionID)
}

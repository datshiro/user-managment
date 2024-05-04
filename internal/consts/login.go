package consts

type LoginType int

const (
	UserNameLoginType LoginType = iota + 1
	EmailLoginType
	PhoneNumberLoginType
)

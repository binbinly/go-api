package service

type RspCheckLogin struct {
	UserId          int    `json:"user_id"`
	Phone           string `json:"phone"`
	Gold            int    `json:"gold"`
	RegisterChannel int    `json:"register_channel"`
}

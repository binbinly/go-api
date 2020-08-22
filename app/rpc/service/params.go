package service

type ParamCheckLogin struct {
	SessionId string `json:"session_id"`
}

type ParamEditGold struct {
	Event    string `json:"event"`
	Gold     int64  `json:"gold"`
	UserId   int    `json:"user_id"`
	GoldType int    `json:"gold_type"`
}

type ParamSms struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

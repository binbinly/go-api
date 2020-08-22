package agency

type Balance struct {
	UserId       int    `json:"user_id"`
	Time         string `json:"time"`
	Win          int64  `json:"win"`
	Award        int64  `json:"award"`
	ActiveUser   uint16 `json:"active_user"`
	Discounts    int64  `json:"discounts"`
	Recharge     int64  `json:"recharge"`
	Withdraw     int64  `json:"withdraw"`
	Fee          int64  `json:"fee"`
	PlatformFee  int64  `json:"platform_fee"`
	RealWin      int64  `json:"real_win"`
	Rate         uint8  `json:"rate"`
	AgencyAmount int64  `json:"agency_amount"`
	Amount       int64  `json:"amount"`
}

func (Balance) TableName() string {
	return "agency_balance"
}

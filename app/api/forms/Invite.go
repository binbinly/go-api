package forms

type InviteDetail struct {
	UserId int    `json:"user_id" binding:"omitempty,numeric"`
	Date   string `json:"date" binding:"omitempty,len=10"`
}

type BalanceDetail struct {
	Date string `json:"date" binding:"required,len=10"`
}

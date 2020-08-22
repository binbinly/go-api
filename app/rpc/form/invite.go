package form

//添加用户
type InviteUserReq struct {
	UserId       int    `json:"user_id" validate:"required"`
	ParentId     int    `json:"parent_id"`
	ChannelId    int    `json:"channel_id" validate:"required"`
	RegisterTime int    `json:"register_time" validate:"required"`
	Username     string `json:"username" validate:"required"`
	Mobile       string `json:"mobile" validate:"required"`
}

//修改用户
type InviteUserEdit struct {
	UserId   int    `json:"user_id" validate:"required"`
	Username string `json:"username"`
	Mobile   string `json:"mobile"`
}

type Currency uint8

const (
	CurrencyGuess    Currency = 0 //竞猜币
	CurrencyPlatform Currency = 1 //平台币
)

// 用户订单结算/注销产生的流水记录
type OrderFlow struct {
	UserId       int64    `json:"user_id" validate:"required"` // 用户ID
	SuperiorId   int64    `json:"superior_id"`                 // 用户上级ID
	Currency     Currency `json:"currency"`                    // 币种
	BetAmount    int64    `json:"bet_amount"`                  // 下单额
	ValidAmount  int64    `json:"valid_amount"`                // 有效下单额
	ProfitAmount int64    `json:"profit_amount"`               // 盈利额
	Flow         int64    `json:"flow"`                        // 流水
	ValidFlow    int64    `json:"valid_flow"`                  // 用户的有效流水
}

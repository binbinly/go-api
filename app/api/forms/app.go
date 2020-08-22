package forms

type ChannelConfig struct {
	ChannelId int `json:"channel" binding:"required,numeric"`
}

type VersionUpgrade struct {
	Num      string `json:"num" binding:"required,v_num"`
	Platform int    `json:"platform" binding:"required,numeric"`
}

type SendSms struct {
	Type   string `json:"type" binding:"required,oneof=login reg apply withdraw phone password"`
	Mobile string `json:"mobile" binding:"required,mobile"`
}

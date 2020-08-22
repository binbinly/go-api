package forms

type AddAgencyForm struct {
	RealName string `json:"real_name" binding:"required,max=20"`
	Mobile   string `json:"mobile" binding:"required,mobile"`
	Code     string `json:"code" binding:"required,numeric"`
	QQ       int    `json:"qq" binding:"numeric"`
	Email    string `json:"email" binding:"omitempty,email"`
	Desc     string `json:"desc" binding:"required,max=500"`
}

type AgencyDetail struct {
	Time string `json:"time" binding:"len=7"`
}

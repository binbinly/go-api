package models

const (
	TableChannel       = "seo_channel"
	TableConfig        = "system_config"
	TableVersion       = "app_version"
	TableInviteUser    = "invite_user"
	TableInviteUserDay = "invite_user_day"
	TableInviteTree    = "invite_tree"
	TableInviteLog     = "invite_log"
	TableInviteConfig  = "invite_config"
	TableInviteBalance = "invite_balance"
)

type Model struct {
	ID        int `gorm:"primary_key" json:"id"`
	CreatedAt int `json:"created_at"`
	UpdatedAt int `json:"updated_at"`
	DeletedAt int `json:"deleted_at"`
}

type ModelId struct {
	ID int `gorm:"primary_key" json:"id"`
}

type ModelCreate struct {
	ID        int `gorm:"primary_key" json:"id"`
	CreatedAt int `json:"created_at"`
}

type ModelOnlyUpdate struct {
	ID        int `gorm:"primary_key" json:"id"`
	UpdatedAt int `json:"updated_at"`
}

type ModelUpdate struct {
	ID        int `gorm:"primary_key" json:"id"`
	CreatedAt int `json:"created_at"`
	UpdatedAt int `json:"updated_at"`
}

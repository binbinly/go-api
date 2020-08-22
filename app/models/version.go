package models

type AppVersion struct {
	VersionNumber string `json:"version_number"`
	VersionName   string `json:"version_name"`
	Desc          string `json:"desc"`
	DownloadUrl   string `json:"download_url"`
	IsCompel      int    `json:"is_compel"`
	IosUrl        string `json:"ios_url"`
}

package entity

type ProjectItem struct {
	Id           string `json:"id"`
	ProjectName  string `json:"projectName"`
	State        int    `json:"state"`
	CreateTime   string `json:"createTime"`
	CreateUserId string `json:"createUserId"`
	IsDelete     int    `json:"isDelete"`
	IndexImage   string `json:"indexImage"`
	Remarks      string `json:"remarks"`
}
type ProjectContent struct {
	Id           string `json:"id"`
	ProjectName  string `json:"projectName"`
	State        int    `json:"state"`
	CreateTime   string `json:"createTime"`
	CreateUserId string `json:"createUserId"`
	IsDelete     int    `json:"isDelete"`
	IndexImage   string `json:"indexImage"`
	Remarks      string `json:"remarks"`
	Content      string `json:"content"`
}

type UploadResp struct {
	Id             string `json:"id"`
	FileName       string `json:"fileName"`
	BucketName     string `json:"bucketName"`
	FileSize       int    `json:"fileSize"`
	FileSuffix     string `json:"fileSuffix"`
	CreateUserId   string `json:"createUserId"`
	CreateUserName string `json:"createUserName"`
	CreateTime     string `json:"createTime"`
	UpdateUserId   string `json:"updateUserId"`
	UpdateUserName string `json:"updateUserName"`
	UpdateTime     string `json:"updateTime"`
}

package entity

type UploadFileReq struct {
	FileName    string            `json:"file_name"`    //原始文件名
	FileSize    int64             `json:"file_size"`    //文件大小
	Bucket      string            `json:"bucket"`       //文件桶，不填则为default, 文件的一种分类，例如 media,doc,xls...
	Meta        map[string]string `json:"meta"`         //文件其他相关元数据。例如owner等
	EncryptName string            `json:"encrypt_name"` //minio内文件名，唯一标识。带有后缀。需要保证唯一性。例如uuid或nanoid拼接
	Data        []byte            `json:"data"`         //文件数据
	AutoConvert bool              `json:"auto_convert"` //是否自动转换
	SaveOrg     bool              `json:"save_org"`     //自动转换后，是否保存原始文件
}

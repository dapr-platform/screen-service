package entity

import "github.com/dapr-platform/common"

type Topology struct {
	ID          string           `json:"id"`         //唯一标识
	CreatedBy   string           `json:"createBy"`   //创建人
	CreatedTime common.LocalTime `json:"createTime"` //创建时间
	UpdatedBy   string           `json:"updateBy"`   //更新人
	UpdatedTime common.LocalTime `json:"updateTime"` //更新时间
	GroupID     string           `json:"groupId"`    //分组id
	FileData    string           `json:"fileData"`   //文件内容
	KeyName     string           `json:"keyName"`    //标识
	Name        string           `json:"name"`       //名称
	ParentID    string           `json:"parentId"`   //父id
	Remark      string           `json:"remark"`     //备注
	Type        string           `json:"type"`       //类型

}

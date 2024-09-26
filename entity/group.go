package entity

type GroupVo struct {
	Id       string     `json:"id"`       //id
	Name     string     `json:"name"`     //名称
	ParentId string     `json:"parentId"` //parent_id
	Level    int        `json:"level"`    //level
	SortNum  int        `json:"sortNum"`
	Type     string     `json:"type"`
	Children []*GroupVo `json:"children"`
}
type GroupVoItem struct {
	Id       string `json:"id"`       //id
	Name     string `json:"name"`     //名称
	ParentId string `json:"parentId"` //parent_id
	Level    int    `json:"level"`    //level
	SortNum  int    `json:"sortNum"`
	Type     string `json:"type"`
}

func GetGroupTreeRecursive(list []*GroupVo, parentNode *GroupVo) {
	for _, v := range list {
		if v.ParentId == parentNode.Id {
			parentNode.Children = append(parentNode.Children, v)
			GetGroupTreeRecursive(list, v)

		}
	}
}

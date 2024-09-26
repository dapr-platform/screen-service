package service

import (
	"context"
	"github.com/dapr-platform/common"
	"github.com/pkg/errors"
	"screen-service/entity"
	"screen-service/model"
)

func GetGroupTrees(ctx context.Context, t string) (result []*entity.GroupVo, err error) {
	qStr := ""
	if t != "" {
		qStr = "type=" + t
	}
	list, err := common.DbQuery[model.Group](ctx, common.GetDaprClient(), model.GroupTableInfo.Name, qStr)
	if err != nil {
		err = errors.Wrap(err, "gropu query list error")
		return
	}

	voList := make([]*entity.GroupVo, 0)
	for _, v := range list {
		var vo = &entity.GroupVo{}
		vo.Id = v.ID
		vo.ParentId = v.ParentID
		vo.Name = v.Name
		vo.Type = v.Type
		vo.SortNum = int(v.SortNum)
		vo.Level = int(v.Level)
		vo.Children = make([]*entity.GroupVo, 0)
		voList = append(voList, vo)
	}

	roots := make([]*entity.GroupVo, 0)
	for _, v := range voList {
		if v.ParentId == "0" || v.ParentId == "" {
			roots = append(roots, v)
		}
	}
	common.Logger.Debug("roots len=", len(roots))
	for _, r := range roots {
		entity.GetGroupTreeRecursive(voList, r)
	}

	return roots, nil
}

package api

import (
	"github.com/dapr-platform/common"
	"github.com/go-chi/chi/v5"
	"github.com/guregu/null"
	"net/http"
	"screen-service/entity"
	"screen-service/model"
	"screen-service/service"
	"strings"
	"time"
)

var (
	_ = null.String{}
)

func InitGroupRoute(r chi.Router) {
	r.Get(common.BASE_CONTEXT+"/group/all-trees", groupTreeHandler)
	r.Post(common.BASE_CONTEXT+"/group/save", UpsertGroupHandler)
	r.Delete(common.BASE_CONTEXT+"/group/{id}", DeleteGroupHandler)
	r.Post(common.BASE_CONTEXT+"/group/batch-delete", batchDeleteGroupHandler)
	r.Post(common.BASE_CONTEXT+"/group/batch-upsert", batchUpsertGroupHandler)
}

// @Summary 分组树
// @Description 分组树
// @Tags Group
// @Param type query string false "分类"
// @Produce  json
// @Success 200 {object} common.Response{data=[]entity.GroupVo} "实例的数组"
// @Failure 500 {object} common.Response "错误code和错误信息"
// @Router /group/all-trees [get]
func groupTreeHandler(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("type")
	data, err := service.GetGroupTrees(r.Context(), t)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpResult(w, common.OK.WithData(data))

}

// @Summary 保存分组
// @Description 保存分组
// @Tags Group
// @Accept       json
// @Param item body entity.GroupVoItem true "Group全部信息"
// @Produce  json
// @Success 200 {object} common.Response{data=entity.GroupVoItem} "实例"
// @Failure 500 {object} common.Response "错误code和错误信息"
// @Router /group/save [post]
func UpsertGroupHandler(w http.ResponseWriter, r *http.Request) {
	var val entity.GroupVoItem
	err := common.ReadRequestBody(r, &val)
	if err != nil {
		common.HttpResult(w, common.ErrParam)
		return
	}
	sub, _ := common.ExtractUserSub(r)
	newFlag := false

	if val.Id == "" {
		val.Id = common.NanoId()
	}

	dbItem := model.Group{
		ID:          val.Id,
		Name:        val.Name,
		UpdatedBy:   sub,
		UpdatedTime: common.LocalTime(time.Now()),
		ParentID:    val.ParentId,
		Type:        val.Type,
		Level:       int32(val.Level),
		SortNum:     int32(val.SortNum),
	}
	if newFlag {
		dbItem.CreatedBy = sub
		dbItem.CreatedTime = common.LocalTime(time.Now())
	}
	err = common.DbUpsert[model.Group](r.Context(), common.GetDaprClient(), dbItem, model.GroupTableInfo.Name, "id")
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK.WithData(val))
}

// @Summary 删除实例
// @Description 删除实例
// @Tags Group
// @Param id  path string true "实例id"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Group} "实例"
// @Failure 500 {object} common.Response "错误code和错误信息"
// @Router /group/{id} [delete]
func DeleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	beforeHook, exists := common.GetDeleteBeforeHook("Group")
	if exists {
		_, err1 := beforeHook(r, id)
		if err1 != nil {
			common.HttpResult(w, common.ErrService.AppendMsg(err1.Error()))
			return
		}
	}
	common.CommonDelete(w, r, common.GetDaprClient(), "p_group", "id", "id")
}

// @Summary 批量删除
// @Description 批量删除
// @Tags Group
// @Accept  json
// @Param ids body []string true "id数组"
// @Produce  json
// @Success 200 {object} common.Response "成功code和成功信息"
// @Failure 500 {object} common.Response "错误code和错误信息"
// @Router /group/batch-delete [post]
func batchDeleteGroupHandler(w http.ResponseWriter, r *http.Request) {

	var ids []string
	err := common.ReadRequestBody(r, &ids)
	if err != nil {
		common.HttpResult(w, common.ErrParam)
	}
	if len(ids) == 0 {
		common.HttpResult(w, common.ErrParam)
	}
	idstr := strings.Join(ids, ",")
	err = common.DbDeleteByOps(r.Context(), common.GetDaprClient(), "p_group", []string{"id"}, []string{"in"}, []any{idstr})
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}

// @Summary 批量更新
// @Description 批量更新
// @Tags Group
// @Accept  json
// @Param entities body []map[string]any true "对象数组"
// @Produce  json
// @Success 200 {object} common.Response "成功code和成功信息"
// @Failure 500 {object} common.Response "错误code和错误信息"
// @Router /group/batch-upsert [post]
func batchUpsertGroupHandler(w http.ResponseWriter, r *http.Request) {

	var entities []map[string]any
	err := common.ReadRequestBody(r, &entities)
	if err != nil {
		common.HttpResult(w, common.ErrParam)
	}
	if len(entities) == 0 {
		common.HttpResult(w, common.ErrParam)
	}

	err = common.DbBatchUpsert[map[string]any](r.Context(), common.GetDaprClient(), entities, model.GroupTableInfo.Name, model.Group_FIELD_NAME_id)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}

	common.HttpResult(w, common.OK)
}

package api

import (
	"github.com/dapr-platform/common"
	"github.com/go-chi/chi/v5"
	"github.com/guregu/null"
	"net/http"
	"net/url"
	"screen-service/entity"
	"screen-service/model"
	"strconv"
	"strings"
	"time"
)

var (
	_ = null.String{}
)

func InitTopologyRoute(r chi.Router) {
	r.Get(common.BASE_CONTEXT+"/topology", TopologyListHandler)
	r.Get(common.BASE_CONTEXT+"/topology-query", TopologyQueryHandler)
	r.Post(common.BASE_CONTEXT+"/topology", UpsertTopologyHandler)
	r.Delete(common.BASE_CONTEXT+"/topology/{id}", DeleteTopologyHandler)
}

// @Summary 根据条件查询
// @Description 根据条件查询, 可设置 order, 以及查询条件等，例如 status=1, name=$like.%25%25CAM%25%25 等
// @Tags Topology
// @Param _order query string false "order"
// @Param _select query string false "select"
// @Param id query string false "id"
// @Param created_by query string false "created_by"
// @Param created_time query string false "created_time"
// @Param updated_by query string false "updated_by"
// @Param updated_time query string false "updated_time"
// @Param group_id query string false "group_id"
// @Param file_data query string false "file_data"
// @Param key_name query string false "key_name"
// @Param name query string false "name"
// @Param parent_id query string false "parent_id"
// @Param remark query string false "remark"
// @Param type query string false "type"
// @Produce  json
// @Success 200 {object} common.Response{data=common.Page{items=[]model.Topology}} "实例的数组"
// @Failure 500 {object} common.Response "错误code和错误信息"
// @Router /topology-query [get]
func TopologyQueryHandler(w http.ResponseWriter, r *http.Request) {
	var vars = r.URL.Query()
	var qstr = ""
	if vars != nil {
		for k, v := range vars {
			//val := strings.Replace(v[0], "%", "%25", -1)
			val := url.QueryEscape(v[0])
			qstr += k + "=" + val + "&"
		}
		qstr = strings.TrimSuffix(qstr, "&")
		common.Logger.Debug("qstr:", qstr)

	}

	data, err := common.DbQuery[model.Topology](r.Context(), common.GetDaprClient(), model.TopologyTableInfo.Name, qstr)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	entities := make([]entity.Topology, 0)
	for _, d := range data {
		entities = append(entities, entity.Topology{
			ID:          d.ID,
			CreatedTime: d.CreatedTime,
			CreatedBy:   d.CreatedBy,
			UpdatedBy:   d.UpdatedBy,
			UpdatedTime: d.UpdatedTime,
			GroupID:     d.GroupID,
			FileData:    d.FileData,
			KeyName:     d.KeyName,
			Name:        d.Name,
			ParentID:    d.ParentID,
			Remark:      d.Remark,
			Type:        d.Type,
		})
	}
	common.HttpResult(w, common.OK.WithData(entities))
}

// @Summary 分页查询topo实例
// @Description 分页查询topo实例
// @Tags Topology
// @Param _page query int false "current page"
// @Param _page_size query int false "page size"
// @Produce  json
// @Success 200 {object} common.Response{data=common.Page{items=[]entity.Topology}} "实例的数组"
// @Failure 500 {object} common.Response "错误code和错误信息"
// @Router /topology [get]
func TopologyListHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("_page")
	pageSize := r.URL.Query().Get("_page_size")
	pageI := 1
	pageSizeI := 10
	if page != "" {
		pageI, _ = strconv.Atoi(page)
	}
	if pageSize != "" {
		pageSizeI, _ = strconv.Atoi(pageSize)
	}
	data, err := common.DbPageQuery[model.Topology](r.Context(), common.GetDaprClient(), pageI, pageSizeI, "", model.TopologyTableInfo.Name, model.Topology_FIELD_NAME_id, "")
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	entities := make([]entity.Topology, 0)
	for _, d := range data.Items {
		entities = append(entities, entity.Topology{
			ID:          d.ID,
			CreatedTime: d.CreatedTime,
			CreatedBy:   d.CreatedBy,
			UpdatedBy:   d.UpdatedBy,
			UpdatedTime: d.UpdatedTime,
			GroupID:     d.GroupID,
			FileData:    d.FileData,
			KeyName:     d.KeyName,
			Name:        d.Name,
			ParentID:    d.ParentID,
			Remark:      d.Remark,
			Type:        d.Type,
		})
	}
	pageData := common.PageGeneric[entity.Topology]{
		Page:     pageI,
		PageSize: pageSizeI,
		Total:    data.Total,
		Items:    entities,
	}
	common.HttpResult(w, common.OK.WithData(pageData))

}

// @Summary 保存topology
// @Description 保存topology
// @Tags Topology
// @Accept       json
// @Param item body entity.Topology true "topology全部信息"
// @Produce  json
// @Success 200 {object} common.Response{data=entity.Topology} "实例"
// @Failure 500 {object} common.Response "错误code和错误信息"
// @Router /topology [post]
func UpsertTopologyHandler(w http.ResponseWriter, r *http.Request) {
	var val entity.Topology
	err := common.ReadRequestBody(r, &val)
	if err != nil {
		common.HttpResult(w, common.ErrParam)
		return
	}
	sub, _ := common.ExtractUserSub(r)
	if val.ID == "" {
		val.ID = common.NanoId()
		val.CreatedTime = common.LocalTime(time.Now())
		val.CreatedBy = sub
		val.CreatedTime = common.LocalTime(time.Now())
	}
	val.UpdatedBy = sub
	val.UpdatedTime = common.LocalTime(time.Now())
	dbentity := model.Topology{
		ID:          val.ID,
		CreatedTime: val.CreatedTime,
		CreatedBy:   val.CreatedBy,
		UpdatedBy:   val.UpdatedBy,
		UpdatedTime: val.UpdatedTime,
		GroupID:     val.GroupID,
		FileData:    val.FileData,
		KeyName:     val.KeyName,
		Name:        val.Name,
		ParentID:    val.ParentID,
		Remark:      val.Remark,
		Type:        val.Type,
	}
	err = common.DbUpsert[model.Topology](r.Context(), common.GetDaprClient(), dbentity, model.TopologyTableInfo.Name, "id")
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK.WithData(val))
}

// @Summary 删除topology
// @Description 删除topology
// @Tags Topology
// @Param id  path string true "实例id"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Topology} "实例"
// @Failure 500 {object} common.Response "错误code和错误信息"
// @Router /topology/{id} [delete]
func DeleteTopologyHandler(w http.ResponseWriter, r *http.Request) {
	common.CommonDelete(w, r, common.GetDaprClient(), "p_topology", "id", "id")
}

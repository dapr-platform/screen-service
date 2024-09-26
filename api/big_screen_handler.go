package api

import (
	"bufio"
	"encoding/json"
	"github.com/dapr-platform/common"
	daprc "github.com/dapr/go-sdk/client"
	"github.com/go-chi/chi/v5"
	"github.com/guregu/null"
	"net/http"
	"screen-service/entity"
	"screen-service/model"
	"strconv"
	"strings"
	"time"
)

var (
	_ = null.String{}
)

func InitBig_screenRoute(r chi.Router) {
	r.Get(common.BASE_CONTEXT+"/project/list", Big_screenListHandler)
	r.Post(common.BASE_CONTEXT+"/project/edit", UpsertBig_screenHandler)
	r.Post(common.BASE_CONTEXT+"/project/create", CreateBig_screenHandler)
	r.Post(common.BASE_CONTEXT+"/project/upload", UploadFileHandler)
	r.Put(common.BASE_CONTEXT+"/project/publish", PublishHandler)
	r.Post(common.BASE_CONTEXT+"/project/save/data", SaveBig_screenContentHandler)
	r.Get(common.BASE_CONTEXT+"/project/getData", GetDataHandler)
	r.Delete(common.BASE_CONTEXT+"/project/delete", DeleteBig_screenHandler)
}

// @Summary 查询所有实例
// @Description 查询所有实例
// @Tags Big_screen
// @Param page query int false "current page"
// @Param limit query int false "page size"
// @Produce  json
// @Success 200 {object} common.Response{data=common.Page{items=[]entity.ProjectItem}} "实例的数组"
// @Failure 500 {object} common.Response "错误code和错误信息"
// @Router /project/list [get]
func Big_screenListHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("limit")
	pageI := 1
	pageSizeI := 10
	if page != "" {
		pageI, _ = strconv.Atoi(page)
	}
	if pageSize != "" {
		pageSizeI, _ = strconv.Atoi(pageSize)
	}
	data, err := common.DbPageQuery[model.Big_screen](r.Context(), common.GetDaprClient(), pageI, pageSizeI, "", model.Big_screenTableInfo.Name, model.Big_screen_FIELD_NAME_id, "")
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	entities := make([]entity.ProjectItem, 0)
	for _, d := range data.Items {
		entities = append(entities, entity.ProjectItem{
			Id:           d.ID,
			ProjectName:  d.Name,
			State:        int(d.State),
			CreateTime:   d.CreatedTime.String(),
			CreateUserId: d.CreatedBy,
			IsDelete:     0,
			IndexImage:   d.IndexImage,
			Remarks:      d.Remarks,
		})
	}
	result := make(map[string]any, 0)
	result["status"] = common.OK.Status
	result["count"] = data.Total
	result["data"] = entities
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	bytes, _ := json.Marshal(result)
	w.Write(bytes)

}

// @Summary 查询数据
// @Description 查询数据
// @Tags Big_screen
// @Param projectId query string false "id"
// @Produce  json
// @Success 200 {object} common.Response{data=common.Page{items=[]entity.ProjectItem}} "实例的数组"
// @Failure 500 {object} common.Response "错误code和错误信息"
// @Router /project/getData [get]
func GetDataHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("projectId")

	data, err := common.DbGetOne[model.Big_screen](r.Context(), common.GetDaprClient(), model.Big_screenTableInfo.Name, model.Big_screen_FIELD_NAME_id+"="+id)
	if err != nil {
		common.HttpResult(w, common.ErrParam.AppendMsg(err.Error()))
		return
	}
	retData := entity.ProjectContent{
		Id:           data.ID,
		ProjectName:  data.Name,
		State:        int(data.State),
		CreateTime:   data.CreatedTime.String(),
		CreateUserId: data.CreatedBy,
		IsDelete:     0,
		IndexImage:   data.IndexImage,
		Remarks:      data.Remarks,
		Content:      data.Content,
	}

	common.HttpSuccess(w, common.OK.WithData(retData))
}

// @Summary 保存实例
// @Description 保存实例
// @Tags Big_screen
// @Accept       json
// @Param item body map[string]interface{} true "实例全部信息"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Big_screen} "实例"
// @Failure 500 {object} common.Response "错误code和错误信息"
// @Router /project/edit [post]
func UpsertBig_screenHandler(w http.ResponseWriter, r *http.Request) {
	var val map[string]any
	err := common.ReadRequestBody(r, &val)
	if err != nil {
		common.HttpResult(w, common.ErrParam)
		return
	}
	var dbVal = make(map[string]any, 0)
	for k, v := range val {
		if v == nil || v == "null" {
			continue
		}
		if k == "projectName" {
			dbVal["name"] = v
		} else if k == "indexImage" {
			dbVal["index_image"] = v
		} else if k == "createTime" {
			dbVal["create_time"] = v
		} else if k == "createUserId" {
			dbVal["create_by"] = v
		} else if k == "id" {
			dbVal["id"] = v
		} else if k == "state" {
			dbVal["state"] = v
		} else if k == "isDelete" {
			//dbVal["state"] = v
		} else if k == "remark" {
			dbVal["remark"] = v
		}
	}
	err = common.DbUpsert[map[string]any](r.Context(), common.GetDaprClient(), dbVal, model.Big_screenTableInfo.Name, "id")
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK)
}

// @Summary 保存实例
// @Description 保存实例
// @Tags Big_screen
// @Accept       json
// @Param item body map[string]interface{} true "实例全部信息"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Big_screen} "实例"
// @Failure 500 {object} common.Response "错误code和错误信息"
// @Router /project/create [post]
func CreateBig_screenHandler(w http.ResponseWriter, r *http.Request) {
	var val entity.ProjectItem
	err := common.ReadRequestBody(r, &val)
	if err != nil {
		common.HttpResult(w, common.ErrParam)
		return
	}
	var dbVal = model.Big_screen{ID: common.NanoId(), State: -1, CreatedTime: common.LocalTime(time.Now())}
	dbVal.Name = val.ProjectName
	val.Id = dbVal.ID

	err = common.DbUpsert[model.Big_screen](r.Context(), common.GetDaprClient(), dbVal, model.Big_screenTableInfo.Name, "id")
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK.WithData(dbVal))
}

// @Summary 保存大屏内容
// @Description 保存大屏内容
// @Tags Big_screen
// @Accept       json
// @Param projectId formData string true "id"
// @Param content formData string true "content"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Big_screen} "实例"
// @Failure 500 {object} common.Response "错误code和错误信息"
// @Router /project/save/data [post]
func SaveBig_screenContentHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("projectId")
	content := r.FormValue("content")

	var dbVal = make(map[string]any, 0)
	dbVal["id"] = id
	dbVal["content"] = content

	err := common.DbUpsert[map[string]any](r.Context(), common.GetDaprClient(), dbVal, model.Big_screenTableInfo.Name, "id")
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK.WithData(dbVal))
}

// @Summary 保存大屏内容
// @Description 保存大屏内容
// @Tags Big_screen
// @Accept       json
// @Param data body map[string]any true "data"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Big_screen} "实例"
// @Failure 500 {object} common.Response "错误code和错误信息"
// @Router /project/publish [put]
func PublishHandler(w http.ResponseWriter, r *http.Request) {
	var val map[string]any
	err := common.ReadRequestBody(r, &val)
	if err != nil {
		common.HttpResult(w, common.ErrParam)
		return
	}
	var dbVal = make(map[string]any, 0)
	for k, v := range val {
		if v == nil || v == "null" {
			continue
		}
		if k == "projectName" {
			dbVal["name"] = v
		} else if k == "indexImage" {
			dbVal["index_image"] = v
		} else if k == "createTime" {
			dbVal["create_time"] = v
		} else if k == "createUserId" {
			dbVal["create_by"] = v
		} else if k == "id" {
			dbVal["id"] = v
		} else if k == "state" {
			dbVal["state"] = v
		} else if k == "isDelete" {
			//dbVal["state"] = v
		} else if k == "remark" {
			dbVal["remark"] = v
		}
	}
	err = common.DbUpsert[map[string]any](r.Context(), common.GetDaprClient(), dbVal, model.Big_screenTableInfo.Name, "id")
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg(err.Error()))
		return
	}
	common.HttpSuccess(w, common.OK)
}

// @Summary 保存大屏封面文件
// @Description 保存大屏封面文件
// @Tags Big_screen
// @Accept       json
// @Param object formData file true "file"
// @Produce  json
// @Success 200 {object} common.Response{data=model.Big_screen} "实例"
// @Failure 500 {object} common.Response "错误code和错误信息"
// @Router /project/upload [post]
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("object")
	if err != nil {
		common.Logger.Error("fileUpload ", err)
		common.HttpResult(w, common.ErrParam.AppendMsg("read file error,"+err.Error()))
		return
	}
	common.Logger.Debug(" upload file ", handler.Filename)
	fileReq := entity.UploadFileReq{}
	bytes := make([]byte, handler.Size)
	_, err = bufio.NewReader(file).Read(bytes)
	if err != nil {
		common.Logger.Error("file read ", err)
		common.HttpResult(w, common.ErrParam.AppendMsg("read file error,"+err.Error()))
		return
	}
	fileReq.Data = bytes
	fileReq.FileSize = int64(len(bytes))
	fileReq.FileName = handler.Filename
	fileReq.Bucket = "public"
	ext := handler.Filename[strings.LastIndex(handler.Filename, ".")+1:]
	fileReq.EncryptName = common.NanoId() + "." + ext

	metaData := map[string]string{}
	fileReq.Meta = metaData
	data, err := json.Marshal(fileReq)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg("json marshal error,"+err.Error()))
		return
	}

	content := &daprc.DataContent{
		ContentType: "application/json",
		Data:        data,
	}

	_, err = common.GetDaprClient().InvokeMethodWithContent(r.Context(), common.FILE_SERVICE_NAME, "/file-upload-json", "POST", content)
	if err != nil {
		common.HttpResult(w, common.ErrService.AppendMsg("invoke file service error,"+err.Error()))
		return
	}
	result := make(map[string]string, 0)
	result["fileName"] = fileReq.EncryptName

	common.HttpSuccess(w, common.OK.WithData(result))
}

// @Summary 删除实例
// @Description 删除实例
// @Tags Big_screen
// @Param ids  query string true "实例id"
// @Produce  json
// @Success 200 {object} common.Response{data=entity.ProjectItem} "实例"
// @Failure 500 {object} common.Response "错误code和错误信息"
// @Router /big-screen/delete [delete]
func DeleteBig_screenHandler(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query().Get("ids")
	for _, id := range strings.Split(ids, ",") {
		common.DbDelete(r.Context(), common.GetDaprClient(), model.Big_screenTableInfo.Name, "id", id)
	}
	common.HttpSuccess(w, common.OK)
}

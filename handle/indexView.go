package handle

import (
	"github.com/flosch/pongo2/v4"
	"imgGo-Server/dao"
	"imgGo-Server/global"
	"imgGo-Server/model"
	"math"
	"net/http"
	"strconv"
)

func IndexView(w http.ResponseWriter, r *http.Request) {
	conf := global.Conf
	MyDebug := global.MyDebug

	MyDebug("--call IndexView--")
	var err error
	var pagination = model.PaginationModel{
		Index: 1,
		Size:  conf.Page.PageSize,
	}

	if index, err := strconv.Atoi(r.FormValue("page")); err == nil {
		pagination.Index = index
	}

	MyDebug("index ->", pagination.Index)

	//if size, err := strconv.Atoi(r.FormValue("size")); err == nil {
	//	pagination.Size = size
	//}

	var items []model.FileDbModel
	dbEngine := dao.GetDao()

	query := "status=1"
	result := dbEngine.Model(&model.FileDbModel{}).Where(query).Order("created_at desc")
	result.Count(&pagination.ItemCount)
	pagination.PageCount = int(math.Ceil(float64(pagination.ItemCount) / float64(int64(pagination.Size))))

	result.Limit(pagination.Size).Offset((pagination.Index - 1) * pagination.Size).Find(&items)

	//if len(items) == 0 {
	//	util.Http404(ctx, "没有可显示的内容了")
	//}

	//解析模板文件
	var t = pongo2.Must(pongo2.FromFile("./template/index.html"))

	//输出文件数据
	err = t.ExecuteWriter(pongo2.Context{
		"conf":            conf.Page,
		"query":           r.FormValue("query"),
		"items":           items,
		"PaginationModel": pagination,
	}, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

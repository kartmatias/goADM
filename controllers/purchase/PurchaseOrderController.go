package purchase

import (
	"bytes"
	"encoding/json"
	"goADM/controllers/base"
	md "goADM/models"
	"goADM/utils"
	"strconv"
	"strings"
)

// PurchaseOrderController purchase order
type PurchaseOrderController struct {
	base.BaseController
}

// Post request
func (ctl *PurchaseOrderController) Post() {
	action := ctl.Input().Get("action")
	switch action {
	case "validator":
		ctl.Validator()
	case "table": //bootstrap table的post请求
		ctl.PostList()
	case "create":
		ctl.PostCreate()
	default:
		ctl.PostList()
	}
}

// Put request
func (ctl *PurchaseOrderController) Put() {
	id := ctl.Ctx.Input.Param(":id")
	ctl.URL = "/purchase/order/"
	if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {
		if order, err := md.GetPurchaseOrderByID(idInt64); err == nil {
			if err := ctl.ParseForm(&order); err == nil {

				if err := md.UpdatePurchaseOrderByID(order); err == nil {
					ctl.Redirect(ctl.URL+id+"?action=detail", 302)
				}
			}
		}
	}
	ctl.Redirect(ctl.URL+id+"?action=edit", 302)

}

// Get request
func (ctl *PurchaseOrderController) Get() {
	ctl.PageName = "采购订单管理"
	action := ctl.Input().Get("action")
	switch action {
	case "create":
		ctl.Create()
	case "edit":
		ctl.Edit()
	case "detail":
		ctl.Detail()
	default:
		ctl.GetList()

	}
	// 标题合成
	b := bytes.Buffer{}
	b.WriteString(ctl.PageName)
	b.WriteString("\\")
	b.WriteString(ctl.PageAction)
	ctl.Data["PageName"] = b.String()
	ctl.URL = "/purchase/order/"
	ctl.Data["URL"] = ctl.URL

	ctl.Data["MenuPurchaseOrderActive"] = "active"
}

// Edit edit purchase order
func (ctl *PurchaseOrderController) Edit() {
	id := ctl.Ctx.Input.Param(":id")
	orderInfo := make(map[string]interface{})
	if id != "" {
		if idInt64, e := strconv.ParseInt(id, 10, 64); e == nil {
			if order, err := md.GetPurchaseOrderByID(idInt64); err == nil {
				ctl.PageAction = order.Name
				orderInfo["name"] = order.Name

			}
		}
	}
	ctl.Data["Action"] = "edit"
	ctl.Data["RecordID"] = id
	ctl.Data["order"] = orderInfo
	ctl.Layout = "base/base.html"
	ctl.TplName = "purchase/purchase_order_form.html"
}

// Create display purchase order create page
func (ctl *PurchaseOrderController) Create() {
	ctl.Data["Action"] = "create"
	ctl.Data["Readonly"] = false
	ctl.PageAction = utils.MsgCreate
	ctl.Layout = "base/base.html"
	ctl.TplName = "purchase/purchase_order_form.html"
}

// Detail display purchase order info
func (ctl *PurchaseOrderController) Detail() {
	//获取信息一样，直接调用Edit
	ctl.Edit()
	ctl.Data["Readonly"] = true
	ctl.Data["Action"] = "detail"
}

// PostCreate post request create purchase order
func (ctl *PurchaseOrderController) PostCreate() {
	order := new(md.PurchaseOrder)
	if err := ctl.ParseForm(order); err == nil {

		if id, err := md.AddPurchaseOrder(order); err == nil {
			ctl.Redirect("/purchase/order/"+strconv.FormatInt(id, 10)+"?action=detail", 302)
		} else {
			ctl.Get()
		}
	} else {
		ctl.Get()
	}
}

// Validator js valid
func (ctl *PurchaseOrderController) Validator() {
	name := ctl.GetString("name")
	name = strings.TrimSpace(name)
	recordID, _ := ctl.GetInt64("recordID")
	result := make(map[string]bool)
	obj, err := md.GetPurchaseOrderByName(name)
	if err != nil {
		result["valid"] = true
	} else {
		if obj.Name == name {
			if recordID == obj.ID {
				result["valid"] = true
			} else {
				result["valid"] = false
			}

		} else {
			result["valid"] = true
		}

	}
	ctl.Data["json"] = result
	ctl.ServeJSON()
}

//PurchaseOrderList 获得符合要求的数据
func (ctl *PurchaseOrderController) PurchaseOrderList(query map[string]interface{}, exclude map[string]interface{}, condMap map[string]map[string]interface{}, fields []string, sortby []string, order []string, offset int64, limit int64) (map[string]interface{}, error) {

	var arrs []md.PurchaseOrder
	paginator, arrs, err := md.GetAllPurchaseOrder(query, exclude, condMap, fields, sortby, order, offset, limit)
	result := make(map[string]interface{})
	if err == nil {

		//使用多线程来处理数据，待修改
		tableLines := make([]interface{}, 0, 4)
		for _, line := range arrs {
			oneLine := make(map[string]interface{})
			oneLine["name"] = line.Name
			oneLine["ID"] = line.ID
			oneLine["id"] = line.ID

			tableLines = append(tableLines, oneLine)
		}
		result["data"] = tableLines
		if jsonResult, er := json.Marshal(&paginator); er == nil {
			result["paginator"] = string(jsonResult)
			result["total"] = paginator.TotalCount
		}
	}
	return result, err
}

// PostList post request json response
func (ctl *PurchaseOrderController) PostList() {
	query := make(map[string]interface{})
	exclude := make(map[string]interface{})
	cond := make(map[string]map[string]interface{})

	fields := make([]string, 0, 0)
	sortby := make([]string, 0, 1)
	order := make([]string, 0, 1)
	offset, _ := ctl.GetInt64("offset")
	limit, _ := ctl.GetInt64("limit")
	orderStr := ctl.GetString("order")
	sortStr := ctl.GetString("sort")
	if orderStr != "" && sortStr != "" {
		sortby = append(sortby, sortStr)
		order = append(order, orderStr)
	} else {
		sortby = append(sortby, "Id")
		order = append(order, "desc")
	}
	if result, err := ctl.PurchaseOrderList(query, exclude, cond, fields, sortby, order, offset, limit); err == nil {
		ctl.Data["json"] = result
	}
	ctl.ServeJSON()

}

// GetList display purchase order with list
func (ctl *PurchaseOrderController) GetList() {
	viewType := ctl.Input().Get("view")
	if viewType == "" || viewType == "table" {
		ctl.Data["ViewType"] = "table"
	}
	ctl.PageAction = utils.MsgList
	ctl.Data["tableId"] = "table-purchase-order"
	ctl.Layout = "base/base_list_view.html"
	ctl.TplName = "purchase/purchase_order_list_search.html"
}

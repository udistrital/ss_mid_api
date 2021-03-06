package controllers

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// UtilsController operations for Beneficiarios
type UtilsController struct {
	beego.Controller
}

// URLMapping ...
func (c *UtilsController) URLMapping() {
	c.Mapping("GetActualDate", c.GetActualDate)
}

// GetActualDate ...
// @Title GetActualDate
// @Description obtiene la fecha actual del servidor
// @Param	body		body 	models.Beneficiarios	true		"body for Beneficiarios content"
// @Success 201 {string} fecha_actual
// @Failure 403 body is empty
// @router /GetActualDate [get]
func (c *UtilsController) GetActualDate() {
	t := time.Now()
	respuesta := map[string]string{"fecha_actual": t.Format("2006-01-02")}
	c.Data["json"] = respuesta
	c.ServeJSON()
}

// ImprimirError estándar para imprimir errores
func ImprimirError(mensaje string, err error) {
	logs.Error(mensaje, " => Error:", err.Error())
	//beego.Error(mensaje, " => Error:", err.Error())
}

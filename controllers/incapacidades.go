package controllers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
)

// IncapacidadesController operations for Incapacidades
type IncapacidadesController struct {
	beego.Controller
}

// URLMapping ...
func (c *IncapacidadesController) URLMapping() {
	c.Mapping("GetPersonas", c.GetPersonas)
}

// GetPersonas ...
// @Title GetPersonas
// @Description obtiene todas las personas que pueden aplicar a cualquier nómina
// @Param	documento		query	string false		"documento de la persona"
// @Success 200 {object} interface{}
// @Failure 403
// @router / [get]
func (c *IncapacidadesController) GetPersonas() {
	var proveedores, contratos, personaNatural, respuesta []map[string]interface{}
	documento := c.GetString("documento")
	try.This(func() {
		if err := getJson("http://"+beego.AppConfig.String("administrativaService")+"/informacion_proveedor?"+
			"limit=6&query=NumDocumento__icontains:"+documento+",TipoPersona:NATURAL", &proveedores); err != nil {
			panic(err)
		}

		for _, proveedor := range proveedores {
			id := strconv.Itoa(int(proveedor["Id"].(float64)))
			if err := getJson("http://"+beego.AppConfig.String("administrativaService")+"/contrato_general?"+
				"query=Estado:true,Contratista:"+id+"&fields=Id,VigenciaContrato", &contratos); err != nil {
				beego.Error("error en contrato general")
				panic(err)
			}

			if err := getJson("http://"+beego.AppConfig.String("administrativaService")+"/informacion_persona_natural?"+
				"limit=1&query=Id:"+proveedor["NumDocumento"].(string), &personaNatural); err != nil {
				beego.Error("error en contrato informacion_persona_natural")
				panic(err)
			}

			var contratosPropios []map[string]interface{} // contratos de cada proveedor

			for _, contrato := range contratos {
				numeroContrato := contrato["Id"].(string)
				vigenciaContrato := contrato["VigenciaContrato"].(float64)

				temp := map[string]interface{}{
					"NumeroContrato":   numeroContrato,
					"VigenciaContrato": vigenciaContrato,
				}
				contratosPropios = append(contratosPropios, temp)
			}

			resp := map[string]interface{}{
				"display":       proveedor["NomProveedor"],
				"value":         proveedor["NumDocumento"],
				"documento":     proveedor["NumDocumento"],
				"contratos":     contratosPropios,
				"tipoDocumento": personaNatural[0]["TipoDocumento"].(map[string]interface{})["Abreviatura"],
			}

			respuesta = append(respuesta, resp)
		}
		c.Data["json"] = respuesta

	}).Catch(func(e try.E) {
		beego.Error("Error en GetPersonas() ", e)
		c.Data["json"] = e
	})
	c.ServeJSON()
}

// IncapacidadesPorPersona ...
// @Title IncapacidadesPorPersona
// @Description obtiene todas las incapacidades activdas de una persona
// @Param	documento		query	string false		"documento de la persona"
// @Success 200 {object} interface{}
// @Failure 403
// @router /incapacidadesPersona/:contrato/:vigencia [get]
func (c *IncapacidadesController) IncapacidadesPorPersona() {
	contrato := c.Ctx.Input.Param(":contrato")
	vigencia := c.Ctx.Input.Param(":vigencia")

	var incapacidades []map[string]interface{}
	try.This(func() {
		incapacidadesLaborales, err := traerIncapacidades("incapacidad_laboral", contrato, vigencia)
		if err != nil {
			panic(err.Error())
		}

		incapacidaGenerales, err := traerIncapacidades("incapacidad_general", contrato, vigencia)
		if err != nil {
			panic(err.Error())
		}

		incapacidades = append(incapacidades, incapacidadesLaborales...)
		incapacidades = append(incapacidades, incapacidaGenerales...)
		c.Data["json"] = incapacidades
	}).Catch(func(e try.E) {
		beego.Error("Error en IncapacidadesPorPersona() ", e)
		c.Data["json"] = e
	})

	c.ServeJSON()
}

func traerIncapacidades(tipoIncapacidad, contrato, vigencia string) (incapacidades []map[string]interface{}, err error) {
	var detalleNovedad []map[string]interface{}
	err = getJson("http://"+beego.AppConfig.String("titanServicio")+"/concepto_nomina_por_persona?query=Concepto.Nombreconcepto:"+tipoIncapacidad+
		",NumeroContrato:"+contrato+",VigenciaContrato:"+vigencia+",Activo:true", &incapacidades)

	for i, v := range incapacidades {
		beego.Info(v["Id"])
		conceptoNominaPorPesona := strconv.Itoa(int(v["Id"].(float64)))
		err = getJson("http://"+beego.AppConfig.String("segSocialService")+"/detalle_novedad_seguridad_social?"+
			"query=ConceptoNominaPorPersona:"+conceptoNominaPorPesona+"&limit=1", &detalleNovedad)
		beego.Info(detalleNovedad)
		incapacidades[i]["Codigo"] = detalleNovedad[0]["Descripcion"]
	}
	return
}

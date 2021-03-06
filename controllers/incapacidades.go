package controllers

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/udistrital/seguridad_social_mid/models"
)

// IncapacidadesController operations for Incapacidades
type IncapacidadesController struct {
	beego.Controller
}

// URLMapping ...
func (c *IncapacidadesController) URLMapping() {
	c.Mapping("BuscarPersonas", c.BuscarPersonas)
	c.Mapping("IncapacidadesPorPersona", c.IncapacidadesPorPersona)
}

// BuscarPersonas ...
// @Title BuscarPersonas
// @Description obtiene todas las personas que pueden aplicar a cualquier nómina
// @Param	documento		query	string false		"documento de la persona"
// @Success 200 {object} interface{}
// @Failure 403
// @router /BuscarPersonas/:documento [get]
func (c *IncapacidadesController) BuscarPersonas() {
	var proveedores, contratos, personaNatural, respuesta []map[string]interface{}
	documento := c.Ctx.Input.Param(":documento")

	try.This(func() {

		if err := getJson("http://"+beego.AppConfig.String("administrativaService")+"informacion_proveedor?"+
			"limit=6&query=NumDocumento__icontains:"+documento+",TipoPersona:NATURAL", &proveedores); err != nil {
			panic(err)
		}

		for _, proveedor := range proveedores {
			idProveedor := strconv.Itoa(int(proveedor["Id"].(float64)))
			
			if err := getJson("http://"+beego.AppConfig.String("administrativaService")+"contrato_general?"+
				"query=Estado:true,Contratista:"+idProveedor, &contratos); err != nil {
				log.Panicln("error en contrato general")
			}

			if err := getJson("http://"+beego.AppConfig.String("administrativaService")+"informacion_persona_natural?"+
				"limit=1&query=Id:"+proveedor["NumDocumento"].(string), &personaNatural); err != nil {
				log.Panicln("error en contrato informacion_persona_natural")
			}

			var contratosPropios []map[string]interface{} // contratos de cada proveedor
			contratosActivos := revisarcontratosActivos(contratos)

			for _, contrato := range contratosActivos {
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
				"idProveedor":   idProveedor,
			}

			respuesta = append(respuesta, resp)
		}
		if respuesta == nil {
			respuesta = append(respuesta, map[string]interface{}{})
		}

		c.Data["json"] = respuesta

	}).Catch(func(e try.E) {
		log.Panicln("Error en GetPersonas() ", e)
		c.Data["json"] = models.Alert{Code: "error"}
	})
	c.ServeJSON()
}

// revisarcontratosActivos revisa si el contrato está activo, si el contrato es de DVE revisa las actas de inicio, si es de CPS revisa contrato_estado
func revisarcontratosActivos(contratos []map[string]interface{}) (contratosActivos []map[string]interface{}) {
	var actaInicio []map[string]time.Time

	for _, contrato := range contratos {
		if strings.Contains(contrato["Id"].(string), "DVE") { // Contrato de DVE

			if err := getJson("http://"+beego.AppConfig.String("administrativaService")+"acta_inicio?"+
				"query=NumeroContrato:"+contrato["Id"].(string)+"&fields=FechaFin", &actaInicio); err != nil {
				log.Panicln("error en contrato informacion_persona_natural")
			}

			for _, acta := range actaInicio {
				if time.Now().Before(acta["FechaFin"]) { // Si la fecha actual es anterior a la fecha fin, entonces el contrato todavía está activo
					contratosActivos = append(contratosActivos, contrato)
				}
			}

		} else {
			var contratoEstado interface{}

			if err := getJson("http://"+beego.AppConfig.String("administrativaService")+"contrato_estado?"+
				"query=NumeroContrato:"+contrato["Id"].(string)+",Estado.NombreEstado:Finalizado", &contratoEstado); err != nil {
				log.Panicln("error en contrato contrato_estado")
			}

			if contratoEstado == nil {
				contratosActivos = append(contratosActivos, contrato)
			}
		}

	}
	
	return
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
			log.Panicln(err.Error())
		}

		incapacidadesGenerales, err := traerIncapacidades("incapacidad_general", contrato, vigencia)
		if err != nil {
			log.Panicln(err.Error())
		}

		prorrogas, err := traerIncapacidades("prorroga_incapacidad", contrato, vigencia)
		if err != nil {
			log.Panicln(err.Error())
		}

		incapacidades = append(incapacidades, incapacidadesLaborales...)
		incapacidades = append(incapacidades, incapacidadesGenerales...)
		incapacidades = append(incapacidades, prorrogas...)
		c.Data["json"] = incapacidades
	}).Catch(func(e try.E) {
		log.Panicf("Error en IncapacidadesPorPersona() ", e)
		c.Data["json"] = map[string]interface{}{"error": e}
	})

	c.ServeJSON()
}

func traerIncapacidades(tipoIncapacidad, contrato, vigencia string) (incapacidades []map[string]interface{}, err error) {
	var detalleNovedad []map[string]interface{}

	_ = getJson("http://"+beego.AppConfig.String("titanServicio")+"/concepto_nomina_por_persona?"+
		"query=Concepto.Nombreconcepto:"+tipoIncapacidad+
		",NumeroContrato:"+contrato+",VigenciaContrato:"+vigencia+",Activo:true&limit=0", &incapacidades)

	for _, v := range incapacidades {
		if len(v) > 0 {
			conceptoNominaPorPesona := strconv.Itoa(int(v["Id"].(float64)))
			_ = getJson("http://"+beego.AppConfig.String("segSocialService")+"/detalle_novedad_seguridad_social?"+
				"query=ConceptoNominaPorPersona:"+conceptoNominaPorPesona+"&limit=1", &detalleNovedad)
			v["Codigo"] = detalleNovedad[0]["Descripcion"]
		}  else {
			incapacidades = nil
		}
	}

	return
}

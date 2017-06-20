package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/udistrital/ss_mid_api/golog"
	"github.com/udistrital/ss_mid_api/models"

	"github.com/astaxie/beego"
)

// DescSeguridadSocialController oprations for DescSeguridadSocial
type DescSeguridadSocialController struct {
	beego.Controller
}

// URLMapping ...
func (c *DescSeguridadSocialController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("CalcularSegSocial", c.CalcularSegSocial)
	c.Mapping("GetConceptosIbc", c.GetConceptosIbc)
	c.Mapping("GenerarPlanillaActivos", c.GenerarPlanillaActivos)
	c.Mapping("GenerarPlanillaPensionados", c.GenerarPlanillaPensionados)
}

// Post ...
// @Title Post
// @Description create DescSeguridadSocial
// @Param	body		body 	models.DescSeguridadSocial	true		"body for DescSeguridadSocial content"
// @Success 201 {int} models.DescSeguridadSocial
// @Failure 403 body is empty
// @router / [post]
func (c *DescSeguridadSocialController) Post() {
	var v models.DescSeguridadSocial
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddDescSeguridadSocial(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *DescSeguridadSocialController) GetConceptosIbc() {
	fmt.Println("GetConceptosIbc")
	var predicados []models.Predicado
	var conceptos []models.Concepto
	var conceptosIbc []models.ConceptosIbc
	err := getJson("http://"+beego.AppConfig.String("rulerServicio")+
		"/predicado?limit=-1&query=Nombre__startswith:conceptos_ibc,Dominio.Id:4", &predicados)
	errConceptoTitan := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/concepto?limit=-1", &conceptos)
	if err != nil && errConceptoTitan != nil {
		c.Data["json"] = err.Error() + errConceptoTitan.Error()
	} else {
		nombres := golog.GetString(FormatoReglas(predicados), "conceptos_ibc(X).", "X")
		for i := 0; i < len(predicados); i++ {
			for j := 0; j < len(conceptos); j++ {
				if nombres[i] == conceptos[j].NombreConcepto {
					aux := models.ConceptosIbc{
						Id:          predicados[i].Id,
						Nombre:      nombres[i],
						Descripcion: conceptos[j].AliasConcepto}
					conceptosIbc = append(conceptosIbc, aux)
					break
				}
			}

		}
		c.Data["json"] = conceptosIbc
	}
	c.ServeJSON()
}

func (c *DescSeguridadSocialController) CalcularSegSocial() {
	idStr := c.Ctx.Input.Param(":id")
	_, err := strconv.Atoi(idStr)
	var alertas []string
	var predicado []models.Predicado
	var buffer bytes.Buffer //objeto para concatenar strings a la variable errores

	var detalleLiquidacion []models.DetalleLiquidacion
	var pagosSeguridadSocial []models.PagosSeguridadSocial

	if err != nil {
		c.Data["json"] = err.Error()
		fmt.Println("Error en detalle liquidacion: ", err)
	} else {

		err := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_liquidacion"+
			"?limit=0&query=Liquidacion.Id:"+idStr+",Concepto.NombreConcepto:ibc_liquidado", &detalleLiquidacion)

		fmt.Println(buffer.String())

		if err != nil {
			fmt.Println("ERROR EN LA PETICION:\n", buffer.String())
			alertas = append(alertas, "error al traer detalle liquidacion")
			c.Data["json"] = alertas
		} else {

			for index := 0; index < len(detalleLiquidacion); index++ {
				predicado = append(predicado, models.Predicado{Nombre: "ibc(" + strconv.Itoa(detalleLiquidacion[index].Persona) + "," + strconv.Itoa(int(detalleLiquidacion[index].ValorCalculado)) + ", salud)."})
				predicado = append(predicado, models.Predicado{Nombre: "ibc(" + strconv.Itoa(detalleLiquidacion[index].Persona) + "," + strconv.Itoa(int(detalleLiquidacion[index].ValorCalculado)) + ", riesgos)."})
				predicado = append(predicado, models.Predicado{Nombre: "ibc(" + strconv.Itoa(detalleLiquidacion[index].Persona) + "," + strconv.Itoa(int(detalleLiquidacion[index].ValorCalculado)) + ", apf)."})
			}

			reglas := CargarReglasBase() + FormatoReglas(predicado) + CargarNovedades(idStr) +
				ValorSaludEmpleado(idStr) + ValorPensionEmpleado(idStr)

			ids := golog.GetInt64(reglas, "v_salud_ud(I,Y).", "I")
			saludUd := golog.GetFloat(reglas, "v_salud_ud(I,Y).", "Y")
			saludTotal := golog.GetInt64(reglas, "v_total_salud(X,T).", "T")
			pensionUd := golog.GetFloat(reglas, "v_pen_ud(I,Y).", "Y")
			pensionTotal := golog.GetInt64(reglas, "v_total_pen(X,T).", "T")
			arl := golog.GetInt64(reglas, "v_arl(I,Y).", "Y")
			caja := golog.GetInt64(reglas, "v_caja(I,Y).", "Y")
			icbf := golog.GetInt64(reglas, "v_icbf(I,Y).", "Y")

			for index := 0; index < len(ids); index++ {
				aux := models.PagosSeguridadSocial{
					Persona:              ids[index],
					SaludUd:              saludUd[index],
					SaludTotal:           saludTotal[index],
					PensionUd:            pensionUd[index],
					PensionTotal:         pensionTotal[index],
					Caja:                 caja[index],
					Icbf:                 icbf[index],
					IdDetalleLiquidacion: detalleLiquidacion[index].Id,
					Arl:                  arl[index]}

				pagosSeguridadSocial = append(pagosSeguridadSocial, aux)
			}

			c.Data["json"] = pagosSeguridadSocial
		}
		c.ServeJSON()
	}
}

// AproximarValor ...
// @Title Aproximar Valor
// @Description Aproxima un valor al número de aproximacion más cercano
// @Param valorInicial valor a cual quiere aproximar
// @Param	aproximacion	multiplo sobre al cual debe aproximarse el valor inicial
// @return valorAproximado el valor aproximado de acuerdo a la aproximacion
func AproximarValor(valorInicial int64, aproximacion int64) (valorAproximado int64) {
	x := float64(valorInicial) / float64(aproximacion)
	y := math.Trunc(float64(valorInicial / aproximacion))
	if (x - y) > 0 {
		valorAproximado = AproximarValor(valorInicial+1, aproximacion)
	} else {
		valorAproximado = valorInicial
	}
	return
}

// ValorSaludEmpleado ...
// @Title Valor Salud Empleado
// @Description Crea todos los hechos con la información del valor de salud
// @Param	idLiquidacion		id de la liquidacion correspondiente
func ValorSaludEmpleado(idLiquidacion string) (valorSaludEmpleado string) {
	var detalleLiquSalud []models.DetalleLiquidacion
	var predicado []models.Predicado

	errSalud := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_liquidacion"+
		"?limit=0&query=Liquidacion.Id:"+idLiquidacion+",Concepto.NombreConcepto:salud", &detalleLiquSalud)

	if errSalud != nil {
		fmt.Println("Error en ValorSaludEmpleado:\n", errSalud)
	} else {
		for index := 0; index < len(detalleLiquSalud); index++ {
			predicado = append(predicado, models.Predicado{Nombre: "v_salud_func(" + strconv.Itoa(detalleLiquSalud[index].Persona) + ", " + strconv.Itoa(int(detalleLiquSalud[index].ValorCalculado)) + ")."})
			valorSaludEmpleado += predicado[index].Nombre + "\n"
		}
	}
	return
}

// ValorPensionEmpleado ...
// @Title Valor Pensión Empleado
// @Description Crea todos los hechos con la información del valor de la pensión
// @Param	idLiquidacion		id de la liquidacion correspondiente
func ValorPensionEmpleado(idLiquidacion string) (valorPensionEmpleado string) {
	var detalleLiquPension []models.DetalleLiquidacion
	var predicado []models.Predicado

	errPension := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_liquidacion"+
		"?limit=0&query=Liquidacion.Id:"+idLiquidacion+",Concepto.NombreConcepto:pension", &detalleLiquPension)

	if errPension != nil {
		fmt.Println("Error en ValorPensionEmpleado:\n", errPension)
	} else {
		for index := 0; index < len(detalleLiquPension); index++ {
			predicado = append(predicado, models.Predicado{Nombre: "v_pen_func(" + strconv.Itoa(detalleLiquPension[index].Persona) + ", " + strconv.Itoa(int(detalleLiquPension[index].ValorCalculado)) + ")."})
			valorPensionEmpleado += predicado[index].Nombre + "\n"
		}
	}
	return
}

// CargarNovedades ...
// @Title CargarNovedades
// @Description obtiene todas los conceptos con naturaleza seguridad_social desde detalle_liquidacion por el id
// @Param	id		path 	string	true		"The key for staticblock"
func CargarNovedades(id string) (novedades string) {
	var conceptos []models.DetalleLiquidacion
	var predicado []models.Predicado

	errLincNo := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_liquidacion?limit=0&query=Liquidacion.Id:"+id+",Concepto.naturaleza:seguridad_social&fields=Concepto,Persona", &conceptos)

	if errLincNo != nil {
		fmt.Println("error en CargarNovedades()", errLincNo)
	} else {
		for index := 0; index < len(conceptos); index++ {
			predicado = append(predicado, models.Predicado{Nombre: "novedad_persona(" + conceptos[index].Concepto.NombreConcepto + ", " + strconv.Itoa(int(conceptos[index].Persona)) + ")."})
			novedades += predicado[index].Nombre + "\n"
		}
	}
	return
}

/*func (c *DescSeguridadSocialController) GenerarArchivoPlano() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	GetA

}*/

// GetOne ...
// @Title Get One
// @Description get DescSeguridadSocial by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.DescSeguridadSocial
// @Failure 403 :id is empty
// @router /:id [get]
func (c *DescSeguridadSocialController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetDescSeguridadSocialById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get DescSeguridadSocial
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.DescSeguridadSocial
// @Failure 403
// @router / [get]
func (c *DescSeguridadSocialController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllDescSeguridadSocial(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the DescSeguridadSocial
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.DescSeguridadSocial	true		"body for DescSeguridadSocial content"
// @Success 200 {object} models.DescSeguridadSocial
// @Failure 403 :id is not int
// @router /:id [put]
func (c *DescSeguridadSocialController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.DescSeguridadSocial{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateDescSeguridadSocialById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the DescSeguridadSocial
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *DescSeguridadSocialController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteDescSeguridadSocial(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *DescSeguridadSocialController) GenerarPlanillaActivos() {
	idStr := c.Ctx.Input.Param(":id")
	idDescSegSocial, _ := strconv.Atoi(idStr)
	var proveedores []models.InformacionProveedor
	var pagosSalud []models.DescSeguridadSocialDetalle
	//var pagoSalud []models.DescSeguridadSocial
	var detalleLiquidacion []models.DetalleLiquidacion
	var detalleIncapcidadLaboral []models.DetalleLiquidacion
	var diasLiquidados []models.DetalleLiquidacion
	var soloLiquidadoDetalle []models.DetalleLiquidacion
	var conceptoPersona []models.ConceptoPorPersona
	var conceptos []models.Concepto
	var personaNatural []models.InformacionPersonaNatural
	var errStrings []string

	tipoRegistro := "02"
	fila := ""

	errLiquidacion := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/detalle_liquidacion?limit=-1", &detalleLiquidacion)
	if errLiquidacion != nil {
		errStrings = append(errStrings, errLiquidacion.Error())
	}

	errProveedores := getJson("http://"+beego.AppConfig.String("agoraServicio")+
		"/informacion_proveedor?limit=0", &proveedores)
	if errProveedores != nil {
		errStrings = append(errStrings, errProveedores.Error())
	}

	errConceptos := getJson("http://"+beego.AppConfig.String("titanServicio")+
		"/concepto?limit=0", &conceptos)
	if errConceptos != nil {
		errStrings = append(errStrings, errConceptos.Error())
	}

	fmt.Println("errStrings: ", errStrings)
	if errStrings == nil {
		secuencia := 1
		for i := 0; i < len(detalleLiquidacion); i++ {
			for j := 0; j < len(proveedores); j++ {
				if detalleLiquidacion[i].Persona == proveedores[j].Id {
					if strings.Contains(fila, strconv.Itoa(int(proveedores[j].NumDocumento))) {
						break
					} else {
						//Novedades
						var ingreso = false
						var retiro = false
						var trasladoPensiones = false
						var trasladoEps = false
						var exterior = false
						var suspencionContrato = false
						var licenciaNoRem = false
						var comisionServicios = false
						var incapacidadGeneral = false
						var licenciaMaternidad = false
						var vacaciones = false
						var licenciaRem = false
						var aporteVoluntario = false
						var variacionCentroTrabajo = false
						var diasIncapcidadLaboral = 0
						//var valorIncapacidadLaboral int
						//var novedad = false

						fila += formatoDato(tipoRegistro, 2)                                    //Tipo Registro
						fila += formatoDato(completarSecuencia(secuencia, 5), 5)                //Secuencia
						fila += formatoDato("CC", 2)                                            //Tip de documento del cotizante
						fila += formatoDato(strconv.Itoa(int(proveedores[j].NumDocumento)), 16) //Número de identificación del cotizante
						fila += formatoDato(completarSecuencia(1, 2), 2)                        //Tipo Cotizante
						fila += formatoDato(completarSecuencia(1, 2), 2)                        //Subtipo de Cotizante
						fila += formatoDato("", 1)                                              //Extranjero no obligado a cotizar pensión

						errConceptoPersona := getJson("http://"+beego.AppConfig.String("titanServicio")+
							"/concepto_por_persona"+
							"?limit=0"+
							"&query=EstadoNovedad:Activo,Persona.Id:"+strconv.Itoa(detalleLiquidacion[i].Persona)+
							",Concepto.Naturaleza:seguridad_social", &conceptoPersona)

						if errConceptoPersona != nil {
							fmt.Println("errConceptoPersona: ", errConceptoPersona)
						}

						fmt.Println("Conceptos para el id: ", detalleLiquidacion[i].Persona)
						fmt.Println(conceptoPersona[0].Concepto.NombreConcepto)
						for h := 0; h < len(conceptoPersona); h++ {
							switch conceptoPersona[h].Concepto.NombreConcepto {
							case "exterior_familia":
								exterior = true
								//novedad = true
							case "suspencionContrato":
								suspencionContrato = true
								//novedad = true
							case "licenciaNoRem":
								licenciaNoRem = true
								//novedad = true
							case "comision_norem":
								comisionServicios = true
								//novedad = true
							case "incapacidad_general":
								incapacidadGeneral = true
								//novedad = true
							case "licenciaMaternidad":
								licenciaMaternidad = true
								//novedad = true
							case "vacaciones":
								vacaciones = true
								//novedad = true
							case "licencia_rem":
								licenciaRem = true
								//novedad = true
							case "aporteVoluntario":
								aporteVoluntario = true
								//novedad = true
							case "variacionCentroTrabajo":
								variacionCentroTrabajo = true
								//novedad = true
							case "incapacidad_laboral":
								errEnfermedadLaboral := getJson("http://"+beego.AppConfig.String("titanServicio")+
									"/detalle_liquidacion?limit=-1", &detalleIncapcidadLaboral)
								if errEnfermedadLaboral != nil {
									fmt.Println("errEnfermedadLaboral: ", errEnfermedadLaboral)
								}
								diasIncapcidadLaboral, _ = strconv.Atoi(detalleIncapcidadLaboral[0].DiasLiquidados)
								//valorIncapacidadLaboral = int(detalleIncapcidadLaboral[0].ValorCalculado)
								//novedad = true
							}
						}

						if exterior {
							fila += formatoDato("X", 1) //Colombiano en el exterior
							fila += formatoDato(" ", 2) //Código del departamento de la ubicación laboral
							fila += formatoDato(" ", 3) //Código del municipio de ubicación laboral
						} else {
							fila += formatoDato("", 1)    //Colombiano en el exterior
							fila += formatoDato("11", 2)  //Código del departamento de la ubicación laboral
							fila += formatoDato("001", 3) //Código del municipio de ubicación laboral
						}

						errPersonaNatural := getJson("http://"+beego.AppConfig.String("agoraServicio")+
							"/informacion_persona_natural"+
							"?limit=1"+
							"&query=Id:"+strconv.FormatFloat(proveedores[j].NumDocumento, 'E', -1, 64), &personaNatural)

						if errPersonaNatural != nil {
							fmt.Println("errPersonaNatural: ", errPersonaNatural)
						}

						fila += formatoDato(personaNatural[0].PrimerApellido, 20)  //Primer apellido
						fila += formatoDato(personaNatural[0].SegundoApellido, 30) //Segundo apellido
						fila += formatoDato(personaNatural[0].PrimerNombre, 20)    //Primer nombre
						fila += formatoDato(personaNatural[0].SegundoNombre, 30)   //Segundo nombre

						fila += formatoDato("", 1) //ING:Ingreso
						fila += formatoDato("", 1) //RET: retiro
						fila += formatoDato("", 1) //TDE: Traslado desde otra EPS o EOC
						fila += formatoDato("", 1) //TAE: Traslado a otra EPS o EOC
						//TDP: Traslado desde otra administradora de pensiones
						if trasladoPensiones {
							fila += formatoDato("X", 1)
						} else {
							fila += formatoDato("", 1)
						}

						fila += formatoDato("", 1) //TAP: Traslado a otra administradora de pensiones
						fila += formatoDato("", 1) //Variación permanente de salario
						fila += formatoDato("", 1) //Correcciones
						fila += formatoDato("", 1) //VST: Variación transitoria de salario

						//SLN: Suspención temporal del contrato de tabajo o licencia no remunerada o comisión de servicios
						if suspencionContrato {
							fila += formatoDato("X", 1)
						} else if licenciaNoRem {
							fila += formatoDato("X", 1)
						} else if comisionServicios {
							fila += formatoDato("C", 1)
						} else {
							fila += formatoDato("", 1)
						}

						//IGE: Incapacidad temporal por enfermedad general
						if incapacidadGeneral {
							fila += formatoDato("X", 1)
						} else {
							fila += formatoDato("", 1)
						}

						//LMA: Licencia de Maternidad o paternidad
						if licenciaMaternidad { //
							fila += formatoDato("X", 1)
						} else {
							fila += formatoDato("", 1)
						}

						//VAC: Vacaciones
						if vacaciones {
							fila += formatoDato("X", 1)
						} else if licenciaRem {
							fila += formatoDato("L", 1)
						} else {
							fila += formatoDato("", 1)
						}

						//AVP: Aporte voluntario
						if aporteVoluntario {
							fila += formatoDato("X", 1)
						} else {
							fila += formatoDato("", 1)
						}

						//VCT: Variación centros de trabajo
						if variacionCentroTrabajo {
							fila += formatoDato("X", 1)
						} else {
							fila += formatoDato("", 1)
						}

						//IRL: Días de incapacidad por accidente de trabajo o enfermedad laboral
						fila += formatoDato(completarSecuencia(diasIncapcidadLaboral, 2), 2)

						//Código de la administradora de fondo de pensiones a la cual pertenece el afiliado
						fila += formatoDato("231001", 6)

						//Código de la admnistradora de pensiones a la cual se traslada el afiliado
						if trasladoPensiones {
							fila += formatoDato("230301", 6)
						} else {
							fila += formatoDato(" ", 6)
						}

						//Código EPS o EOC a la cual pertenece el afiliado
						fila += formatoDato("EPS010", 6)

						//Código EPS o EOC a la cual se traslada el afiliado
						if trasladoEps {
							fila += formatoDato("EPS012", 6)
						} else {
							fila += formatoDato(" ", 6)
						}

						//Código CCF a la cual pertenece el afiliado
						fila += formatoDato("CCF04", 6)

						errDiasLiquidados := getJson("http://"+beego.AppConfig.String("titanServicio")+
							"/detalle_liquidacion?limit=0"+
							"&query=Concepto.NombreConcepto:ibc_novedad,Persona:"+
							strconv.Itoa(detalleLiquidacion[i].Persona), &diasLiquidados)
						if errDiasLiquidados != nil {
							fmt.Println("errDiasLiquidados: ", errDiasLiquidados)
						}
						diasCotizados, _ := strconv.Atoi(diasLiquidados[0].DiasLiquidados)

						if ingreso || retiro {
							fila += formatoDato(completarSecuencia(diasCotizados, 2), 2) //Número de días cotizados a pensión
							fila += formatoDato(completarSecuencia(diasCotizados, 2), 2) //Número de días cotizados a salud
						} else {
							fila += formatoDato("30", 2) //Número de días cotizados a pensión
							fila += formatoDato("30", 2) //Número de días cotizados a salud
						}

						if incapacidadGeneral || licenciaMaternidad || vacaciones || diasIncapcidadLaboral > 0 {
							fila += formatoDato(completarSecuencia(diasCotizados, 2), 2) //Número de días cotizados a ARL
							fila += formatoDato(completarSecuencia(diasCotizados, 2), 2) //Número de días cotizados a CCF
						} else {
							fila += formatoDato("30", 2) //Número de días cotizados a ARL
							fila += formatoDato("30", 2) //Número de días cotizados a CCF
						}

						errSalarioBasico := getJson("http://"+beego.AppConfig.String("titanServicio")+
							"/detalle_liquidacion?limit=0"+
							"&query=Concepto.NombreConcepto:salarioBase,Persona:"+
							strconv.Itoa(detalleLiquidacion[i].Persona), &soloLiquidadoDetalle)
						if errSalarioBasico != nil {
							fmt.Println("errSalarioBasico: ", errSalarioBasico)
						} else {
							salarioBase := strconv.FormatInt(soloLiquidadoDetalle[0].ValorCalculado, 10)
							fila += formatoDato(salarioBase, 9) //Salario básico
						}

						fila += formatoDato("", 1) //Salario integral

						errSoloLiquidado := getJson("http://"+beego.AppConfig.String("titanServicio")+
							"/detalle_liquidacion?limit=0"+
							"&query=Concepto.NombreConcepto:ibc_liquidado,Persona:"+
							strconv.Itoa(detalleLiquidacion[i].Persona), &soloLiquidadoDetalle)
						if errSoloLiquidado != nil {
							fmt.Println("errSoloLiquidado: ", errSoloLiquidado)
						} else {
							ibcLiquidado := int(soloLiquidadoDetalle[0].ValorCalculado)
							fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC pensión
							fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC salud
							fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC ARL
							fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC CCF
						}

						fila += formatoDato(completarSecuencia(16, 7), 7) //Tarifa de aportes pensiones

						errPagosSalud := getJson("http://"+beego.AppConfig.String("seguridadSocialService")+
							"/desc_seguridad_social_detalle"+
							"?query=IdDescSeguridadSocial:"+strconv.Itoa(idDescSegSocial)+
							",IdDetalleLiquidacion:"+strconv.Itoa(soloLiquidadoDetalle[0].Id), &pagosSalud)
						if errPagosSalud != nil {
							fmt.Println("errPagosSalud: ", errPagosSalud)
						}

						fmt.Println("http://" + beego.AppConfig.String("seguridadSocialService") +
							"/desc_seguridad_social_detalle" +
							"?query=IdDescSeguridadSocial:" + strconv.Itoa(idDescSegSocial) +
							",IdDetalleLiquidacion:" + strconv.Itoa(detalleLiquidacion[i].Id))

						//Cotización obligatoria a pensiones
						for _, pagoPension := range pagosSalud {
							if pagoPension.IdTipoPagoSeguridadSocial.Nombre == "Pension" {
								fila += formatoDato(completarSecuencia(int(pagoPension.Valor), 9), 9)
								break
							}
						}

						fila += formatoDato(completarSecuencia(0, 9), 9) //Aporte voluntario del afiliado al fondo de pensiones obligatorias

						//Aporte voluntario del aportante al fondo de pensiones obligatoria
						errAporteVoluntario := getJson("http://"+beego.AppConfig.String("titanServicio")+
							"/detalle_liquidacion"+
							",Persona:"+strconv.Itoa(detalleLiquidacion[i].Persona), &soloLiquidadoDetalle)

						if errAporteVoluntario != nil {
							fmt.Println("errAporteVoluntario", errAporteVoluntario)
							fila += formatoDato(completarSecuencia(0, 9), 9)
						} else {
							for _, liquidado := range soloLiquidadoDetalle {
								if liquidado.Concepto.NombreConcepto == "nombreRegla2176" {
									fila += formatoDato(strconv.FormatInt(liquidado.ValorCalculado, 10), 9)
								} else if liquidado.Concepto.NombreConcepto == "nombreRegla2178" {
									fila += formatoDato(strconv.FormatInt(liquidado.ValorCalculado, 10), 9)
								} else if liquidado.Concepto.NombreConcepto == "nombreRegla2173" {
									fila += formatoDato(strconv.FormatInt(liquidado.ValorCalculado, 10), 9)
								} else {
									break
								}
							}
						}

						fila += formatoDato(completarSecuencia(0, 9), 9) //Total cotización Sistema General de Pensiones
						fila += formatoDato(completarSecuencia(0, 9), 9) //Aportes a fondo de solidaridad pensional subcuenta de solidaridad
						fila += formatoDato(completarSecuencia(0, 9), 9) //Aportes a fondo de solidaridad pensional subcuenta de subsistencia
						fila += formatoDato(completarSecuencia(0, 9), 9) //Valor no retenido por aportes voluntarios
						fila += formatoDato("12.5", 7)                   //Tarifa de aportes salud

						//Cotización obligatoria a salud
						for _, pagoSalud := range pagosSalud {
							if pagoSalud.IdTipoPagoSeguridadSocial.Nombre == "Salud" {
								fila += formatoDato(completarSecuencia(int(pagoSalud.Valor), 9), 9)
								break
							}
						}

						fila += formatoDato(completarSecuencia(0, 9), 9) //Valor UPC Adicional
						fila += formatoDato("", 15)                      //Nº de autorización de la incapacidad por enfermedad general
						fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de la incapacidad por enfermedad general
						fila += formatoDato("", 15)                      //Nº de autorización de la licencia de maternidad o paternidad
						fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de la licencia de maternidad

						fila += formatoDato(completarSecuenciaString("0.000522", 9), 9) //Tarifa de aportes a Riegos Laborales

						fila += formatoDato(completarSecuenciaString("0", 9), 9) //Centro de trabajo CT

						//Cotización obligatoria a Sistema General de Riesgos Laborales
						for _, pagoArl := range pagosSalud {
							if pagoArl.IdTipoPagoSeguridadSocial.Nombre == "ARL" {
								fila += formatoDato(completarSecuencia(int(pagoArl.Valor), 9), 9)
								break
							}
						}

						fila += formatoDato(completarSecuenciaString("4", 7), 7) //Tarifa de aportes CCF

						//Valor aporte CCF
						for _, pagoCaja := range pagosSalud {
							if pagoCaja.IdTipoPagoSeguridadSocial.Nombre == "Caja" {
								fila += formatoDato(completarSecuencia(int(pagoCaja.Valor), 9), 9)
								break
							}
						}

						fila += formatoDato(completarSecuencia(0, 7), 7) //Tarifa de aportes SENA
						fila += formatoDato(completarSecuencia(0, 9), 9) //Valor Aportes SENA

						fila += formatoDato(completarSecuencia(3, 7), 7) //Tarifa de aportes ICBF

						//Valor aporte ICBF
						for _, pagoIcbf := range pagosSalud {
							if pagoIcbf.IdTipoPagoSeguridadSocial.Nombre == "ICBF" {
								fila += formatoDato(completarSecuencia(int(pagoIcbf.Valor), 9), 9)
							}
						}

						fila += formatoDato(completarSecuencia(0, 7), 7) //Tarifa de aportes ESAP
						fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de aporte ESAP
						fila += formatoDato(completarSecuencia(0, 7), 7) //Tarifa de aportes MEN
						fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de aporte MEN

						//fila += formatoDato(, longitud)
						fila += "\n"
						secuencia++
					}
				}
			}
		}
		fmt.Println("Filas:\n", fila)
		c.Data["json"] = fila
	}
	c.ServeJSON()
}

func (c *DescSeguridadSocialController) GenerarPlanillaPensionados() {
	var proveedores []models.InformacionProveedor
	var personasNatural []models.InformacionPersonaNatural
	//var pagosSalud []models.DescSeguridadSocialDetalle
	//var pagoSalud []models.DescSeguridadSocial
	var detalleLiquidacion []models.DetalleLiquidacion
	var errStrings []string
	tipoRegistro := "02"
	fila := ""

	errLiquidacion := getJson("http://"+beego.AppConfig.String("titanServicio")+"/detalle_liquidacion"+
		"?limit=-1", &detalleLiquidacion)
	if errLiquidacion != nil {
		errStrings = append(errStrings, errLiquidacion.Error())
	}

	errProveedores := getJson("http://"+beego.AppConfig.String("titanServicio")+"/informacion_proveedor"+
		"?limit=-1", &proveedores)
	if errProveedores != nil {
		errStrings = append(errStrings, errProveedores.Error())
	}

	errPersonaNatural := getJson("http://"+beego.AppConfig.String("titanServicio")+"/informacion_persona_natural"+
		"?limit=0", &personasNatural)
	if errPersonaNatural != nil {
		errStrings = append(errStrings, errPersonaNatural.Error())
	}

	fmt.Println("**errStrings: ", errStrings)
	if errStrings == nil {
		secuencia := 1
		for i := 0; i < len(proveedores); i++ {
			for j := 0; j < len(detalleLiquidacion); j++ {
				for k := 0; k < len(personasNatural); k++ {
					if proveedores[i].Id == detalleLiquidacion[j].Persona {
						if int(proveedores[i].NumDocumento) == personasNatural[k].Id {
							fila += formatoDato(tipoRegistro, 2)
							fila += formatoDato(completarSecuencia(secuencia, 5), 5)
							fila += formatoDato(personasNatural[k].PrimerApellido, 20)
							fila += formatoDato(personasNatural[k].SegundoApellido, 30)
							fila += formatoDato(personasNatural[k].PrimerNombre, 20)
							fila += formatoDato(personasNatural[k].SegundoApellido, 30)
							fila += formatoDato("CC", 2)
							fila += formatoDato(strconv.Itoa(int(personasNatural[k].Id)), 16)
							fila += "\n"
							secuencia++
						}
					}
				}
			}
		}
		fmt.Println("Filas:\n", fila)
		c.Data["json"] = fila
	}
	c.ServeJSON()
}

func completarSecuencia(num, cantSecuencia int) (secuencia string) {
	tamanioNum := len(strconv.Itoa(num))
	for i := 0; i < cantSecuencia-tamanioNum; i++ {
		secuencia += "0"
	}
	secuencia += strconv.Itoa(num)
	return
}

func completarSecuenciaString(num string, cantSecuencia int) (secuencia string) {
	tamanioNum := len(num)
	for i := 0; i < cantSecuencia-tamanioNum; i++ {
		secuencia += "0"
	}
	secuencia += num
	return
}

func formatoDato(texto string, longitud int) (textoEscribir string) {
	for _, r := range texto {
		textoEscribir += string(r)
	}
	for i := 0; i < longitud-len(texto); i++ {
		textoEscribir += " "
	}
	return
}

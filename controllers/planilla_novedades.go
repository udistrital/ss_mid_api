package controllers

import (
	"github.com/astaxie/beego"
)

// PlanillaNovedadesController operations for Planilla_novedades
type PlanillaNovedadesController struct {
	beego.Controller
}

// URLMapping ...
func (c *PlanillaNovedadesController) URLMapping() {
	// c.Mapping("GenerarPlanillaN", c.GenerarPlanillaN)
}

/*


func (c *PlanillaNovedadesController) GenerarPlanillaN() {
idStr := c.Ctx.Input.Param(":id")
idDescSegSocial, _ := strconv.Atoi(idStr)
var proveedores []models.InformacionProveedor
var upc []models.UpcAdicional
var detallePreliquidacion []models.DetallePreliquidacion
var detLiq []models.DetallePreliquidacion
var conceptos []models.Concepto
var personaNatural []models.InformacionPersonaNatural
var conceptosSeguridadSocial []models.Concepto
var errStrings []string
//formatoFecha := "2006-01-02"
tipoRegistro := "02"

// Se obtienen todos los conceptos de seguridad social en tabla conceptos de titan
errConceptosSs := getJson("http://"+beego.AppConfig.String("titanServicio")+
"/concepto?limit=0&query=Naturaleza:seguridad_social", &conceptosSeguridadSocial)
if errConceptosSs != nil {
fmt.Println("errConceptosSs: ", errConceptosSs)
}

errLiquidacion := getJson("http://"+beego.AppConfig.String("titanServicio")+
"/detalle_liquidacion?limit=-1", &detallePreliquidacion)
if errLiquidacion != nil {
errStrings = append(errStrings, errLiquidacion.Error())
}

errProveedores := getJson("http://"+beego.AppConfig.String("titanServicio")+
"/informacion_proveedor?limit=0", &proveedores)
/*errProveedores := getJson("http://"+beego.AppConfig.String("agoraServicio")+
"/informacion_proveedor?limit=0", &proveedores)
if errProveedores != nil {
errStrings = append(errStrings, errProveedores.Error())
}

errConceptos := getJson("http://"+beego.AppConfig.String("titanServicio")+
"/concepto?limit=0", &conceptos)
if errConceptos != nil {
errStrings = append(errStrings, errConceptos.Error())
}

errUpc := getJson("http://"+beego.AppConfig.String("seguridadSocialService")+
"/upc_adicional?limit=0", &upc)
if errUpc != nil {
errStrings = append(errStrings, errUpc.Error())
}

fmt.Println("errStrings: ", errStrings)
if errStrings == nil {
secuencia := 1
for i := 0; i < len(detallePreliquidacion); i++ {
for j := 0; j < len(proveedores); j++ {
if detallePreliquidacion[i].Persona == proveedores[j].Id {
if strings.Contains(fila, strconv.Itoa(int(proveedores[j].NumDocumento))) {
break
} else {
var ibcLiquidado = 0
var pagoSalud = 0
var pagoPension = 0
var pagoArl = 0
var pagoCaja = 0
var pagoIcbf = 0

fila += formatoDato(tipoRegistro, 2)                     //Tipo Registro
fila += formatoDato(completarSecuencia(secuencia, 5), 5) //Secuencia

fila += formatoDato("CC", 2)                                            //Tip de documento del cotizante
fila += formatoDato(strconv.Itoa(int(proveedores[j].NumDocumento)), 16) //Número de identificación del cotizante
fila += formatoDato(completarSecuencia(1, 2), 2)                        //Tipo Cotizante
fila += formatoDato(completarSecuencia(1, 2), 2)                        //Subtipo de Cotizante
fila += formatoDato("", 1)                                              //Extranjero no obligado a cotizar pensión

errPersonaNatural := getJson("http://"+beego.AppConfig.String("titanServicio")+
"/informacion_persona_natural"+
"?limit=1"+
"&query=Id:"+strconv.FormatFloat(proveedores[j].NumDocumento, 'E', -1, 64), &personaNatural)

/*errPersonaNatural := getJson("http://"+beego.AppConfig.String("agoraServicio")+
"/informacion_persona_natural"+
"?limit=1"+
"&query=Id:"+strconv.FormatFloat(proveedores[j].NumDocumento, 'E', -1, 64), &personaNatural)

if errPersonaNatural != nil {
fmt.Println("errPersonaNatural: ", errPersonaNatural)
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

fila += formatoDato(personaNatural[0].PrimerApellido, 20)  //Primer apellido
fila += formatoDato(personaNatural[0].SegundoApellido, 30) //Segundo apellido
fila += formatoDato(personaNatural[0].PrimerNombre, 20)    //Primer nombre
fila += formatoDato(personaNatural[0].SegundoNombre, 30)   //Segundo nombre

// --AQUÍ VA LA FUNCIÓN DE LAS NOVEDADES!--  //
establecerNovedades(strconv.Itoa(detallePreliquidacion[i].Persona))

//Código de la administradora de fondo de pensiones a la cual pertenece el afiliado
fila += formatoDato("231001", 6)

//Código de la admnistradora de pensiones a la cual se traslada el afiliado
// Si hay un translado, debe aparecer el nuevo código, de lo contrario será un campo vació
if trasladoPensiones {
fila += formatoDato("230301", 6)
} else {
fila += formatoDato(" ", 6)
}

//Código EPS o EOC a la cual pertenece el afiliado
fila += formatoDato("EPS010", 6)

//Código EPS o EOC a la cual se traslada el afiliado
// Si hay un translado, debe aparecer el nuevo código, de lo contrario será un campo vació
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
strconv.Itoa(detallePreliquidacion[i].Persona), &detLiq)
if errDiasLiquidados != nil {
fmt.Println("errDiasLiquidados: ", errDiasLiquidados)
}
diasCotizados, _ := strconv.Atoi(detLiq[0].DiasLiquidados)

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

fmt.Println(detallePreliquidacion[i].Persona)
errSalarioBasico := getJson("http://"+beego.AppConfig.String("titanServicio")+
"/detalle_liquidacion?limit=0"+
"&query=Concepto.NombreConcepto:salarioBase,Persona:"+
strconv.Itoa(detallePreliquidacion[i].Persona), &detLiq)
if errSalarioBasico != nil {
fmt.Println("errSalarioBasico: ", errSalarioBasico)
} else {
salarioBase := strconv.FormatInt(detLiq[0].ValorCalculado, 10)
fila += formatoDato(salarioBase, 9) //Salario básico
}

fila += formatoDato("", 1) //Salario integral

errSoloLiquidado := getJson("http://"+beego.AppConfig.String("titanServicio")+
"/detalle_liquidacion?limit=0"+
"&query=Concepto.NombreConcepto:ibc_liquidado,Persona:"+
strconv.Itoa(detallePreliquidacion[i].Persona), &detLiq)
if errSoloLiquidado != nil {
fmt.Println("errSoloLiquidado: ", errSoloLiquidado)
} else {
ibcLiquidado = int(detLiq[0].ValorCalculado)
fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC pensión
fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC salud
fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC ARL
fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC CCF
}

fila += formatoDato(completarSecuencia(16, 7), 7) //Tarifa de aportes pensiones

//Cotización obligatoria a pensiones
for _, pago := range conceptosSeguridadSocial {
switch pago.NombreConcepto {
case "pensionTotal":
pagoPension, _ = strconv.Atoi(obtenerPago(strconv.Itoa(idDescSegSocial), strconv.Itoa(detLiq[0].Id), strconv.Itoa(pago.Id)))
case "saludTotal":
pagoSalud, _ = strconv.Atoi(obtenerPago(strconv.Itoa(idDescSegSocial), strconv.Itoa(detLiq[0].Id), strconv.Itoa(pago.Id)))
case "icbf":
pagoIcbf, _ = strconv.Atoi(obtenerPago(strconv.Itoa(idDescSegSocial), strconv.Itoa(detLiq[0].Id), strconv.Itoa(pago.Id)))
case "caja":
pagoCaja, _ = strconv.Atoi(obtenerPago(strconv.Itoa(idDescSegSocial), strconv.Itoa(detLiq[0].Id), strconv.Itoa(pago.Id)))
case "arl":
pagoArl, _ = strconv.Atoi(obtenerPago(strconv.Itoa(idDescSegSocial), strconv.Itoa(detLiq[0].Id), strconv.Itoa(pago.Id)))
}
}

fila += formatoDato(completarSecuencia(pagoPension, 9), 9) // Cotización obligatoria a pensiones

fila += formatoDato(completarSecuencia(0, 9), 9) //Aporte voluntario del afiliado al fondo de pensiones obligatorias

//Aporte voluntario del aportante al fondo de pensiones obligatoria
errAporteVoluntario := getJson("http://"+beego.AppConfig.String("titanServicio")+
"/detalle_liquidacion"+
"?query=Persona:"+strconv.Itoa(detallePreliquidacion[i].Persona), &detLiq)

if errAporteVoluntario != nil {
fmt.Println("errAporteVoluntario", errAporteVoluntario)
fila += formatoDato(completarSecuencia(0, 9), 9)
} else {
for _, liquidado := range detLiq {
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

fila += formatoDato(completarSecuencia(0, 9), 9)         // Total cotización Sistema General de Pensiones
fila += formatoDato(completarSecuencia(0, 9), 9)         // Aportes a fondo de solidaridad pensional subcuenta de solidaridad
fila += formatoDato(completarSecuencia(0, 9), 9)         // Aportes a fondo de solidaridad pensional subcuenta de subsistencia
fila += formatoDato(completarSecuencia(0, 9), 9)         // Valor no retenido por aportes voluntarios
fila += formatoDato("12.5", 7)                           // Tarifa de aportes salud
fila += formatoDato(completarSecuencia(pagoSalud, 9), 9) // Cotización obligatoria a salud

fila += formatoDato(completarSecuencia(0, 9), 9) //Valor UPC Adicional
fila += formatoDato("", 15)                      //Nº de autorización de la incapacidad por enfermedad general
fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de la incapacidad por enfermedad general
fila += formatoDato("", 15)                      //Nº de autorización de la licencia de maternidad o paternidad
fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de la licencia de maternidad

fila += formatoDato(completarSecuenciaString("0.000522", 9), 9) //Tarifa de aportes a Riegos Laborales

fila += formatoDato(completarSecuenciaString("0", 9), 9) //Centro de trabajo CT
fila += formatoDato(completarSecuencia(pagoArl, 9), 9)   // Cotización obligatoria a salud

fila += formatoDato(completarSecuenciaString("4", 7), 7) //Tarifa de aportes CCF
fila += formatoDato(completarSecuencia(pagoCaja, 9), 9)  // Cotización obligatoria a salud

fila += formatoDato(completarSecuencia(0, 7), 7) //Tarifa de aportes SENA
fila += formatoDato(completarSecuencia(0, 9), 9) //Valor Aportes SENA

fila += formatoDato(completarSecuencia(3, 7), 7)        //Tarifa de aportes ICBF
fila += formatoDato(completarSecuencia(pagoIcbf, 9), 9) // Cotización obligatoria a salud

fila += formatoDato(completarSecuencia(0, 7), 7) //Tarifa de aportes ESAP
fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de aporte ESAP
fila += formatoDato(completarSecuencia(0, 7), 7) //Tarifa de aportes MEN
fila += formatoDato(completarSecuencia(0, 9), 9) //Valor de aporte MEN

//Para los registros de las UPC
/*for _, upcAdicional := range upc {
if upcAdicional.PersonaAsociada == detallePreliquidacion[i].Persona {
fila += formatoDato(texto, longitud)
}
}

// Estos campos están vacios porque solo aplican a los registros que son upc
fila += formatoDato(" ", 2)  // Tipo de documento del cotizante principal
fila += formatoDato(" ", 16) // Número de identificación del cotizante principal

fila += formatoDato("N", 1)     // Cotizante exonerado de pago de aporte salud, SENA e ICBF - Ley 1607 de 2012
fila += formatoDato("14-23", 6) // Código de la administradora de Riesgos Laborales a la cual pertenece el afiliado
fila += formatoDato("1", 1)     // Clase de Riesgo en la que se encuentra el afiliado
fila += formatoDato("", 1)      // Indicador tarifa especial pensiones (Actividades de alto riesgo, Senadores, CTI y Aviadores aplican)

//Fechas de novedades (AAAA-MM-DD)
fila += formatoDato(fechaIngreso, 10)          //Fecha ingreso
fila += formatoDato(fechaRetiro, 10)           //Fecha retiro
fila += formatoDato(fechaInicioVsp, 10)        //Fecha inicio VSP
fila += formatoDato(fechaInicioSuspencion, 10) //Fecha inicio SLN
fila += formatoDato(fechaFinSuspencion, 10)    //Fecha fin SLN
fila += formatoDato(fechaInicioIge, 10)        //Fecha inicio IGE
fila += formatoDato(fechaFinIge, 10)           //Fecha fin IGE
fila += formatoDato(fechaInicioLma, 10)        //Fecha inicio LMA
fila += formatoDato(fechaFinLma, 10)           //Fecha fin LMA
fila += formatoDato(fechaInicioVac, 10)        //Fecha inicio VAC-LR
fila += formatoDato(fechaFinVac, 10)           //Fecha fin VAC-LR
fila += formatoDato(fechaInicioVct, 10)        //Fecha inicio VCT
fila += formatoDato(fechaFinVct, 10)           //Fecha fin VCT
fila += formatoDato(fechaInicioIrl, 10)        //Fecha inicio IRL
fila += formatoDato(fechaFinIrl, 10)           //Fecha fin IRL

fila += formatoDato(completarSecuencia(ibcLiquidado, 9), 9) //IBC otros parafiscales difenrentes a CCF
fila += formatoDato("240", 3)
fila += "\n"
secuencia++
fmt.Println("aqui va uno")
}
}
}
}
c.Data["json"] = fila
}
c.ServeJSON()
}
*/

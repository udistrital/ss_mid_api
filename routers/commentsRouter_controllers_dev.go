package routers

import "github.com/astaxie/beego"

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PagoController"],
		beego.ControllerComments{
			Method:           "CalcularSegSocial",
			Router:           `CalcularSegSocial/:id`,
			AllowHTTPMethods: []string{"get"},
			Params:           nil})

	// beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PagoController"],
	// 	beego.ControllerComments{
	// 		Method:           "NovedadesPorPersona",
	// 		Router:           `NovedadesPorPersona/:persona`,
	// 		AllowHTTPMethods: []string{"get"},
	// 		Params:           nil})

	/*beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PagoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PagoController"],
	beego.ControllerComments{
		Method:           "SumarPagosSalud",
		Router:           `SumarPagosSalud/:idPeriodoPago`,
		AllowHTTPMethods: []string{"get"},
		Params:           nil})*/

	/*beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PlanillasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PlanillasController"],
	beego.ControllerComments{
		Method:           "GenerarPlanillaPensionados",
		Router:           `GenerarPlanillaPensionados/:id`,
		AllowHTTPMethods: []string{"get"},
		Params:           nil})*/

	/*beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PlanillasController"] = append(beego.GlobalControllerRouter["github.com/udistrital/seguridad_social_mid/controllers:PlanillasController"],
	beego.ControllerComments{
		Method:           "GenerarPlanillaN",
		Router:           `GenerarPlanillaN/:id`,
		AllowHTTPMethods: []string{"get"},
		Params:           nil})*/
}

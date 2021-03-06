package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/seguridad_social_mid/models"
	//"github.com/astaxie/beego"
)

func convertirMapa(arr []interface{}) map[string][]interface{} {
	var (
		arry     []interface{}
		contrato string
	)
	returnedMap := make(map[string][]interface{})

	for i := range arr {
		tempMap := arr[i].(map[string]interface{})
		for key, value := range tempMap {
			if key == "NumeroContrato" {
				arry = append(arry, value)
				contrato = value.(string)
			}
			returnedMap[contrato] = arr
		}
	}
	return returnedMap
}

// sendJson envía un json
//@Param url ruta del servicio a la que se le envía el json
//@Param trequest el tipo de petición (GET, POST, PUT, DELETE, etc)
//@Param target la respuesta de la peticón REST
//@Param json que se envia al serivico
func sendJson(url string, trequest string, target interface{}, datajson interface{}) error {
	b := new(bytes.Buffer)
	if datajson != nil {
		json.NewEncoder(b).Encode(datajson)
	}
	client := &http.Client{}
	req, err := http.NewRequest(trequest, url, b)
	r, err := client.Do(req)
	//r, err := http.Post(url, "application/json; charset=utf-8", b)
	if err != nil {
		fmt.Println("error", err)
		//beego.Error("error", err)
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

// getJsonWSO2 ...
func getJsonWSO2(urlp string, target interface{}) error {
	b := new(bytes.Buffer)
	//proxyUrl, err := url.Parse("http://10.20.4.15:3128")
	//http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlp, b)
	req.Header.Set("Accept", "application/json")
	r, err := client.Do(req)
	//r, err := http.Post(url, "application/json; charset=utf-8", b)
	if err != nil {
		fmt.Println("error", err)
		//beego.Error("error", err)
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

// getJson obtiene un json de una URL
//@Param url ruta del servicio del que se recibe el json
//@Param target json que se recibe del servicio
func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func diff(a, b time.Time) (year, month, day int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	oneDay := time.Hour * 5
	a = a.Add(oneDay)
	b = b.Add(oneDay)
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)

	// Normalize negative values
	/*if day < 0{
				day = 0
			}
			if month < 0 {
	        month = 0
	    }*/
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

// GetParametroEstandar devuelve un mapa con la información de los tipos de documento
func GetParametroEstandar() (map[int]string, error) {
	var parametrosEstandar []models.ParametroEstandar
	parametros := make(map[int]string)
	err := getJson("http://"+beego.AppConfig.String("administrativaService")+
		"/parametro_estandar?query=ClaseParametro:Tipo%20Documento&limit=-1", &parametrosEstandar)
	if err != nil {
		return nil, err
	}
	for _, value := range parametrosEstandar {
		parametros[value.Id] = value.Abreviatura
	}
	return parametros, nil
}

func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}

/*
// CargarReglasBase ...
func CargarReglasBase(dominio string) (reglas string) {
	var reglasbase string = ``
	var v []models.Predicado

	if err:= getJson("http://localhost:8080/v1/predicado?limit=0&query=Dominio.Nombre:"+dominio, &v); err == nil {
		reglasbase = reglasbase + FormatoReglas(v)
	} else {
		fmt.Println("err: ", err)
	}
	fmt.Println(reglasbase)
	return reglasbase
}*/

// CargarReglasBase carga las reglas base del ruler
func CargarReglasBase() (reglas string) {
	reglas = `
			%% 		HECHOS PARA ACTIVOS
			%%(ud, conceptoDeDescuento, porcentaje, concepto, nominaCorrespondiente, valorPorcentaje, vigencia).
			concepto(descuento, porcentaje, salud, planta, 0.085,	2017). 	%%descuento salud ud
			concepto(descuento, porcentaje, pension, planta, 0.12, 2017).	%%descuento pension ud
			concepto(descuento, porcentaje, arl, planta, 0.00522, 2017). %%descuento de ARL

			%% HECHOS PARA HCS
			concepto(descuento, porcentaje, salud, salarios, 0.085,	2017). 	%%descuento salud ud
			concepto(descuento, porcentaje, pension, salarios, 0.12, 2017).	%%descuento pension ud
			concepto(descuento, porcentaje, arl, salarios, 0.00522, 2017). %%descuento de ARL

			%% HECHOS PARA CONTRATISTAS
			%%(ud, conceptoDeDescuento, porcentaje, concepto, nominaCorrespondiente, valorPorcentaje, vigencia).
			concepto(descuento, porcentaje, salud, contratistas, 0.125,	2017). 	%%descuento salud contratista
			concepto(descuento, porcentaje, pension, contratistas, 0.16, 2017).	%%descuento pension contratista
			concepto(descuento, porcentaje, arl, contratistas, 0.00522, 2017). %%descuento de ARL contratista

			%%		HECHOS PARA PENSIONADOS
			concepto(X, devengo, porcentaje, salud, pensionado, 0.12, 2017).	%%descuento de salud pensionado

			%% Hechos para aportes parafiscales
			concepto(descuento, porcentaje, caja,	5, 0.04, 2017).	%%caja de compensa familiar
			concepto(descuento, porcentaje,	icbf, 5, 0.03, 2017).	%%ICBF

			%%		NOVEDADES
			%%(descripcion, persona)
			novedad(licencia_maternidad, salud, 0).

			novedad(licencia_maternidad, pension, 0).

			novedad(exterior_familia, arl, 0).
			novedad(vacaciones, arl, 0).
			novedad(licencia_norem, arl, 0).
			novedad(comision_norem, arl, 0).
			novedad(incapacidad_laboral, arl, 0).
			novedad(incapacidad_general, arl, 0).
			novedad(prorroga_incapacidad, arl, 0).
			novedad(licencia_maternidad, arl, 0).
			novedad(licencia_paternidad, arl, 0).

			novedad_persona(-1,-1).

			%%salario minimo legal mensual vigente
			smlmv(737717, 2017).

			%%		SALUD
			v_salud_ud(I,Y) :- concepto(Z,T,salud,planta,V,2017), ibc(I,W, salud), (novedad_persona(N,I), novedad(N,salud,U) -> Y is ((V * W) * U) approach 100; Y is (V * W) approach 100).
			v_total_salud(X,T) :- v_salud_func(X,Y), v_salud_ud(X,U), T is (Y + U) approach 100.


			%%		PENSION
			v_pen_ud(I,Y) :- concepto(Z,T,pension,planta,V,2017), ibc(I,W, salud), (novedad_persona(N,I), novedad(N,pension,U) -> Y is ((V * W) * U) approach 100; Y is (V * W) approach 100).
			v_total_pen(X,T) :- v_pen_func(X,Y), v_pen_ud(X,U), T is (Y + U) approach 100.
			v_pen_contratista(I,Y,C) :- concepto(Z,T,pension,contratista,V,2017), ibc(I,W,C,salud), Y is (V * W) approach 100.

			%%		ARL
			v_arl(I,Y) :- concepto(Z,T,arl,planta,V,2017), ibc(I,W,riesgos), (novedad_persona(N,I), novedad(N,arl,U) -> Y is ((V * W) * U) approach 100; Y is (V * W) approach 100).

			%%		FONDO DE SOLIDARIDAD
			v_fondo1(X,S,D,Y) :- ibc(X,W,apf), smlmv(M,2017),
			(S is 4*M, W@>= S, D is 16*M, W@< D -> Y is W * 0.01;
				S is 16*M, W@>= S, D is 17*M, W@< D -> Y is W * 0.012;
				S is 17*M, W@>= S, D is 18*M, W@< D -> Y is W * 0.014;
				S is 18*M, W@>= S, D is 19*M, W@< D -> Y is W * 0.016;
				S is 19*M, W@>= S, D is 20*M, W@=< D -> Y is W * 0.018;
				S is 20*M, W@> S -> Y is W * 0.02), Y is Y approach 100.	%calculo de fondo de solidaridad 1

			%% 		PAGO UPC
			v_upc(I,Y,Z) :- ibc(I,W,salud,D), upc(Z,V,I), Y is W - V.

			%%		CAJA DE COMPENSACION FAMILIAR
			v_caja(I,Y) :- concepto(Z,T,caja,X,V,2017), ibc(I,W,apf), Y is (V * W) approach 100.

			%%		ICBF
			v_icbf(I,Y) :- concepto(Z,T,icbf,X,V,2017), ibc(I,W,apf), Y is (V * W) approach 100.
	`
	return
}

// CrearReglas crea reglas para cálculos de novedades en las filas auxiliares del archivo plano
func CrearReglas(pagoCompleto bool) string {
	var reglas string
	if pagoCompleto {
		reglas = `
			concepto(salud, 0.125,	2017). 	%%descuento salud ud
			concepto(pension, 0.16, 2017).	%%descuento pension ud
			concepto(arl, 0.00522, 2017). %%descuento de ARL
			concepto(caja, 0.04, 2017).	%%caja de compensa familiar
			concepto(icbf, 0.03, 2017).	%%ICBF`
	} else {
		reglas = `
			concepto(salud,  0.085,	2017). 	%%descuento salud ud
			concepto(pension, 0.12, 2017).	%%descuento pension ud
			concepto(arl, 0, 2017). %%descuento de ARL
			concepto(caja, 0, 2017).	%%caja de compensa familiar
			concepto(icbf, 0, 2017).	%%ICBF`
	}
	reglas += `
		v_total_salud(X,T) :- concepto(salud, V, 2017), ibc(X, W, salud), T is (V * W) approach 100.
		v_total_pen(X,T) :- concepto(pension, V, 2017), ibc(X, W, salud), T is (V * W) approach 100.
		v_arl(I,Y) :- concepto(arl, V, 2017), ibc(I, W, riesgos), Y is (V * W) approach 100. 
		v_caja(I,Y) :- concepto(caja, V, 2017), ibc(I,W,apf), Y is (V * W) approach 100.
		v_icbf(I,Y) :- concepto(icbf, V, 2017), ibc(I,W,apf), Y is (V * W) approach 100.
		`
	return reglas
}

// FormatoReglas crea el formato necesario de las reglas y hechos para golog a partir de un arreglo de predicados
func FormatoReglas(v []models.Predicado) (reglas string) {
	var arregloReglas = make([]string, len(v))
	reglas = ""

	for i := 0; i < len(v); i++ {
		arregloReglas[i] = v[i].Nombre
	}

	for i := 0; i < len(arregloReglas); i++ {
		reglas = reglas + arregloReglas[i] + "\n"
	}

	return
}

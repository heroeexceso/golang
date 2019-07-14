package mongodb

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/heroeexceso/golang/httputility"
	"github.com/heroeexceso/golang/logutility"
)

// Locomotive ... es la estructura que se define para el documento asociado
type Locomotive struct {
	Model         string `json:"model"`
	PowerType     string `json:"powertype"`
	Builder       string `json:"builder"`
	BuildDate     string `json:"buildate"`
	WheelSystem   string `json:"wheelsystem"`
	MaximunSpeed  int    `json:"maximunspeed"`
	PowerOutputHP int    `json:"poweroutputhp"`
}

//  Constantes.-
const port = ":3000"
const logEnabled = false

// Init ... funcion de inicio de la API
func Init() {

	//  Crear un nuevo router.-
	router := mux.NewRouter()

	//  Configurar los recursos.-

	//	GET.-
	router.HandleFunc("/mongoDB/HeaderInfo", getHeaders).Methods("GET")
	router.HandleFunc("/mongoDB/Database/{databaseName}/Collection/{collectionName}/Locomotives", getLocomotives).Methods("GET")
	router.HandleFunc("/mongoDB/Database/{databaseName}/Collection/{collectionName}/Locomotive/{model}", getLocomotive).Methods("GET")

	//	POST.-
	router.HandleFunc("/mongoDB/Database/{databaseName}/Collection/{collectionName}/Locomotive", postLocomotive).Methods("POST")

	//	PATCH.-
	router.HandleFunc("/mongoDB/Database/{databaseName}/Collection/{collectionName}/Locomotive/{model}", patchLocomotive).Methods("PATCH")

	//	DELETE.-
	router.HandleFunc("/mongoDB/Database/{databaseName}/Collection/{collectionName}/Locomotive/{model}", deleteLocomotive).Methods("DELETE")

	//	Dejar el servicio escuchando.-
	log.Fatal(http.ListenAndServe(port, router))
}

// getHeaders ... obtener el header del request
func getHeaders(w http.ResponseWriter, r *http.Request) {
	//	Obtener la información del Request.-
	request := httputility.GetInfoRequest(r)

	//	Devolver el Json con la información.-
	httputility.GetJsonResponse(w, request)
}

// getLocomotives ... obtener todos los documentos de la Base de Datos
func getLocomotives(w http.ResponseWriter, r *http.Request) {

	//	Obtener la base de datos y colección a utilizar.-
	params := mux.Vars(r)
	databaseName := params["databaseName"]
	collectionName := params["collectionName"]

	//	Obtener el modelo a filtrar.-
	model := ""
	powerType := ""
	builder := ""
	wheelSystem := ""

	newDocument, err := getLocomotiveWithFilter(databaseName, collectionName, model, powerType, builder, wheelSystem)
	if err != nil {
		httputility.GetJsonResponseMessage(w, "getLocomotives: "+err.Error())
	} else {
		httputility.GetJsonResponse(w, newDocument)
	}
}

// getLocomotive ... obtener un documento de la Base de Datos
func getLocomotive(w http.ResponseWriter, r *http.Request) {

	//	Obtener la base de datos y colección a utilizar.-
	params := mux.Vars(r)
	databaseName := params["databaseName"]
	collectionName := params["collectionName"]

	//	Obtener el modelo a filtrar.-
	model := params["model"]

	powerType := ""
	builder := ""
	wheelSystem := ""

	newDocument, err := getLocomotiveWithFilter(databaseName, collectionName, model, powerType, builder, wheelSystem)
	if err != nil {
		httputility.GetJsonResponseMessage(w, "getLocomotive: "+err.Error())
	} else {
		httputility.GetJsonResponse(w, newDocument)
	}
}

func getLocomotiveWithFilter(databaseName string, collectionName string, model string, powerType string, builder string, wheelsystem string) (interface{}, error) {
	var document Locomotive

	//	Configurar los filtros.-
	filter := make(map[string]interface{})
	if model != "" {
		filter["model"] = model
	}
	if powerType != "" {
		filter["powertype"] = powerType
	}
	if builder != "" {
		filter["builder"] = builder
	}
	if wheelsystem != "" {
		filter["wheelsystem"] = wheelsystem
	}

	//	Obtener los documentos que cumplan el criterio
	newDocument, err := GetDocuments(databaseName, collectionName, filter, &document)
	if err != nil {
		return newDocument, logutility.GetExceptionError(err, "getDocumentWithFilter")
	}

	return newDocument, nil
}

// postLocomotive ... enviar un documento a la Base de Datos
func postLocomotive(w http.ResponseWriter, r *http.Request) {
	var document Locomotive

	//	Obtener la base de datos y colección a utilizar.-
	params := mux.Vars(r)
	databaseName := params["databaseName"]
	collectionName := params["collectionName"]

	//	Decodificar el documento json recibido y dejarlo en la variable de tipo struct.-
	_ = json.NewDecoder(r.Body).Decode(&document)

	//	Insertar la locomotora.-
	err := InsertDocument(databaseName, collectionName, document)
	if err != nil {
		httputility.GetJsonResponseMessage(w, "postLocomotive: "+err.Error())
	} else {
		httputility.GetJsonResponse(w, document)
	}
}

// patchLocomotive ... enviar un documento a la Base de Datos
func patchLocomotive(w http.ResponseWriter, r *http.Request) {
	var document Locomotive

	//	Obtener la base de datos y colección a utilizar.-
	params := mux.Vars(r)
	databaseName := params["databaseName"]
	collectionName := params["collectionName"]

	//	Obtener el modelo a filtrar.-
	model := params["model"]

	//	Configurar los filtros.-
	filter := make(map[string]interface{})
	if model != "" {
		filter["model"] = model
	}

	//	Decodificar el documento json recibido y dejarlo en la variable de tipo struct.-
	_ = json.NewDecoder(r.Body).Decode(&document)

	//	Configurar los updates.-
	update := make(map[string]interface{})
	if document.PowerType != "" {
		update["powertype"] = document.PowerType
	}
	if document.Builder != "" {
		update["builder"] = document.Builder
	}
	if document.BuildDate != "" {
		update["builddate"] = document.BuildDate
	}
	if document.WheelSystem != "" {
		update["wheelsystem"] = document.WheelSystem
	}
	if document.MaximunSpeed > 0 {
		update["maximunspeed"] = document.MaximunSpeed
	}
	if document.PowerOutputHP > 0 {
		update["poweroutputhp"] = document.PowerOutputHP
	}

	//	Actualizar la locomotora.-
	updatedCount, err := UpdateDocument(databaseName, collectionName, filter, update)
	if err != nil {
		httputility.GetJsonResponseMessage(w, "patchLocomotive: "+err.Error())
	} else {
		if updatedCount == 0 {
			httputility.GetJsonResponseMessage(w, "patchLocomotive: No se encontró ningun documento a actualizar en la Base de Datos (MongoDB).")
		} else {
			if updatedCount == 1 {
				httputility.GetJsonResponseMessage(w, "patchLocomotive: Se actualizó correctamente el documento de la Base de Datos (MongoDB).")
			} else {
				httputility.GetJsonResponseMessage(w, "patchLocomotive: Se actualizaron correctamente "+string(updatedCount)+" documentos de la Base de Datos (MongoDB).")
			}
		}
	}
}

// deleteLocomotive ... enviar un documento a la Base de Datos
func deleteLocomotive(w http.ResponseWriter, r *http.Request) {

	//	Obtener la base de datos y colección a utilizar.-
	params := mux.Vars(r)
	databaseName := params["databaseName"]
	collectionName := params["collectionName"]

	//	Obtener el modelo a filtrar.-
	model := params["model"]

	//	Configurar los filtros.-
	filter := make(map[string]interface{})
	if model != "" {
		filter["model"] = model
	}

	//	Actualizar el documento.-
	deletedCount, err := DeleteDocument(databaseName, collectionName, filter)
	if err != nil {
		//logutil.PrintMessage("getMessage: No se pudo obtener el string de conexión.")
		httputility.GetJsonResponseMessage(w, "deleteLocomotive: "+err.Error())
	} else {
		if deletedCount == 0 {
			httputility.GetJsonResponseMessage(w, "deleteLocomotive: No se encontró ningun documento a eliminar en la Base de Datos (MongoDB).")
		} else {
			if deletedCount == 1 {
				httputility.GetJsonResponseMessage(w, "deleteLocomotive: Se eliminó correctamente el documento de la Base de Datos (MongoDB).")
			} else {
				httputility.GetJsonResponseMessage(w, "deleteLocomotive: Se eliminaron correctamente "+string(deletedCount)+" documentos de la Base de Datos (MongoDB).")
			}
		}
	}
}

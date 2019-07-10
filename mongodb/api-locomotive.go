package mongodb

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/heroeexceso/golang/httputility"
)

// Locomotive ... es la estructura que se define para el documento asociado
type Locomotive struct {
	Model string `json:"model"`
	Brand string `json:"brand"`
	Year  int    `json:"year"`
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
	router.HandleFunc("/mongoDB/Database/{databaseName}/Collection/{collectionName}/Locomotive/{id}", getLocomotive).Methods("GET")

	//	POST.-
	router.HandleFunc("/mongoDB/Database/{databaseName}/Collection/{collectionName}/Locomotive", postLocomotive).Methods("POST")

	//	PATCH.-
	router.HandleFunc("/mongoDB/Database/{databaseName}/Collection/{collectionName}/Locomotive/{id}", patchLocomotive).Methods("PATCH")

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

// getLocomotive ... obtener un documento de la Base de Datos
func getLocomotive(w http.ResponseWriter, r *http.Request) {
	var document Locomotive

	//	Obtener la base de datos y colección a utilizar.-
	params := mux.Vars(r)
	databaseName := params["databaseName"]
	collectionName := params["collectionName"]
	ID := params["id"]

	newDocument, err := GetDocumentByID(databaseName, collectionName, ID, &document)
	if err != nil {
		httputility.GetJsonResponseMessage(w, "getLocomotive: "+err.Error())
	} else {
		httputility.GetJsonResponse(w, newDocument)
	}
}

// getLocomotives ... obtener todos los documentos de la Base de Datos
func getLocomotives(w http.ResponseWriter, r *http.Request) {
	var document Locomotive

	//	Obtener la base de datos y colección a utilizar.-
	params := mux.Vars(r)
	databaseName := params["databaseName"]
	collectionName := params["collectionName"]

	newDocument, err := GetDocuments(databaseName, collectionName, &document)
	if err != nil {
		httputility.GetJsonResponseMessage(w, "getLocomotives: "+err.Error())
	} else {
		httputility.GetJsonResponse(w, newDocument)
	}
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

	//	Insertar el documento (isJson = true).-
	err := InsertDocument(databaseName, collectionName, document)
	if err != nil {
		//logutil.PrintMessage("getMessage: No se pudo obtener el string de conexión.")
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
	ID := params["id"]

	//	Decodificar el documento json recibido y dejarlo en la variable de tipo struct.-
	_ = json.NewDecoder(r.Body).Decode(&document)

	//	Insertar el documento (isJson = true).-
	matchedCount, err := UpdateDocumentByID(databaseName, collectionName, ID, document)
	if err != nil {
		//logutil.PrintMessage("getMessage: No se pudo obtener el string de conexión.")
		httputility.GetJsonResponseMessage(w, "patchLocomotive: "+err.Error())
	} else {
		if matchedCount == 0 {
			httputility.GetJsonResponseMessage(w, "patchLocomotive: No se encontró ningun documento a actualizar en la Base de Datos (MongoDB).")
		} else {
			if matchedCount == 1 {
				httputility.GetJsonResponseMessage(w, "patchLocomotive: Se actualizó correctamente el documento de la Base de Datos (MongoDB).")
			} else {
				httputility.GetJsonResponseMessage(w, "patchLocomotive: Se actualizaron correctamente "+string(matchedCount)+" documentos de la Base de Datos (MongoDB).")
			}
		}
	}
}

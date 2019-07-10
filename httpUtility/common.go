package httputility

import (
	"encoding/json"
	"net/http"
)

// HeaderRequest ... estructura para devolver la información (campo-valor) asociada al header request
type HeaderRequest struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

// RequestInfo ... estructura para devolver la información asociada al request
type RequestInfo struct {
	Host       string          `json:"host"`
	Method     string          `json:"method"`
	Proto      string          `json:"proto"`
	RemoteAddr string          `json:"remoteaddr"`
	URL        string          `json:"url"`
	Header     []HeaderRequest `json:"header"`
}

// Exception ... devolver el formato del mensaje (JSON)
type Exception struct {
	Message string `json:"message"`
}

// GetHeaderKey ... obtener el valor de una Key del Header
func GetHeaderKey(req *http.Request, key string) string {
	return req.Header.Get(key)
}

// SetHeaderKey ... setear el valor de una Key del Header
func SetHeaderKey(req *http.Request, key string, value string) {
	req.Header.Set(key, value)
}

// DeleteHeaderKey ... eliminar el valor de una Key del Header
func DeleteHeaderKey(req *http.Request, key string) {
	req.Header.Del(key)
}

// GetInfoRequest ... devolver la información del Header
func GetInfoRequest(req *http.Request) interface{} {
	var request RequestInfo
	var hRequest HeaderRequest

	//	Asignar los valores genéricos.-
	request.Host = req.Host
	request.Method = req.Method
	request.Proto = req.Proto
	request.RemoteAddr = req.RemoteAddr
	request.URL = req.URL.String()

	for field, value := range req.Header {
		//	Asignar los valores del header.-
		hRequest.Field = field
		hRequest.Value = value[0]

		request.Header = append(request.Header, hRequest)
	}

	return request
}

// GetJsonResponse ... devolver en el Response lo que se recibe en la variable de interface
func GetJsonResponse(w http.ResponseWriter, value interface{}) {
	json.NewEncoder(w).Encode(&value)
}

// GetJsonResponseMessage ... devolver en el Response lo que se recibe en la variable string
func GetJsonResponseMessage(w http.ResponseWriter, message string) {
	json.NewEncoder(w).Encode(Exception{Message: message})
}

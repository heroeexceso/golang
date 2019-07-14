package logutility

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/heroeexceso/golang/httputility"
)

// GetExceptionError ... obtener el mensaje de error formateado
func GetExceptionError(err error, text string) error {
	return errors.New(text + ": " + err.Error())
}

// GetExceptionText ... obtener el mensaje de error formateado
func GetExceptionText(text string) error {
	return errors.New(text)
}

// FailOnError ... devolver un mensaje de error en pantalla
func FailOnError(err error, value string) {
	if err != nil {
		log.Fatalf("%s %s", value, err)
	}
}

// PrintMessage ... imprimir un valor de tipo string en pantalla
func PrintMessage(value string) {
	fmt.Println(value)
}

// PrintByte ... imprimir un valor array de tipo byte en pantalla
func PrintByte(value []byte) {
	fmt.Println(string(value))
}

// PrintStruct ... imprimir un valor de tipo interface en pantalla
func PrintStruct(value interface{}) {
	fmt.Printf("Valor: %+v\n", value)
}

// PrintArrayStruct ... imprimir un valor de tipo interface en pantalla
func PrintArrayStruct(value []interface{}) {
	for _, valueItem := range value {
		fmt.Printf("Valor: %+v\n", valueItem)
	}
}

// PrintHeaderRequest ... imprimir el header del request recibido
func PrintHeaderRequest(req *http.Request) {
	//	Obtener los valores del request.-
	request := httputility.GetInfoRequest(req)

	//	Imprimir la estructura.-
	PrintStruct(request)
}

// PrintInResponseWriter ... imprimir la dupla campo y valor en el response
func PrintInResponseWriter(w http.ResponseWriter, field string, value string) {
	fmt.Fprintf(w, "%v: %+v\n", field, value)
}

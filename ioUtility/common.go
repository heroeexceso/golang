package ioutility

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/heroeexceso/golang/logutility"
	"gopkg.in/yaml.v2" //"gopkg.in/yaml.v2"
)

// GetExceptionError ... obtener el mensaje de error formateado
func GetExceptionError(err error, text string) error {
	return errors.New(text + ": " + err.Error())
}

// GetContentFile ... obtener el contenido de un archivo
func GetContentFile(fileName string) (string, error) {

	//  Leer el archivo.-
	bytesReaded, err := ioutil.ReadFile(fileName)
	if err != nil {
		//logutility.FailOnError(err, "GetContentFile:")
		return "", GetExceptionError(err, "GetContentFile")
	}

	//  Obtener el contenido del mismo para informarlo.-
	content := string(bytesReaded)

	//  Devolver el contenido del archivo.-
	return content, nil
}

// GetConfigurationJSON ... obtener la configuración del archivo Json indicado.-
func GetConfigurationJSON(fileName string, value interface{}, logEnabled bool) error {
	//  Abrir el archivo.-
	jsonFile, err := os.Open(fileName)
	logutility.FailOnError(err, "GetConfigurationJson (Open):")
	defer jsonFile.Close()

	//	Convertir a bytes.-
	byteValue, err := ioutil.ReadAll(jsonFile)
	logutility.FailOnError(err, "GetConfigurationJson (ReadAll):")

	if logEnabled == true {
		logutility.PrintByte(byteValue)
		logutility.PrintStruct(value)
	}

	//	Unmarshal el json.-
	err = json.Unmarshal(byteValue, &value)
	logutility.FailOnError(err, "GetConfigurationJson (Unmarshal):")

	if logEnabled == true {
		logutility.PrintStruct(value)
	}

	//return "amqp://guest:guest@bhhdwwbueg82:5672/", nil
	//return "amqp://guest:pato2019@bhhdwwbueg82:5672/",
	return nil
}

// GetConfigurationYaml ... obtener la configuración del archivo Yaml indicado.-
func GetConfigurationYaml(fileName string, value interface{}, logEnabled bool) error {
	//  Abrir el archivo.-
	jsonFile, err := os.Open(fileName)
	logutility.FailOnError(err, "GetConfigurationYaml (Open):")
	defer jsonFile.Close()

	//	Convertir a bytes.-
	byteValue, err := ioutil.ReadAll(jsonFile)
	logutility.FailOnError(err, "GetConfigurationYaml (ReadAll):")

	if logEnabled == true {
		logutility.PrintByte(byteValue)
		logutility.PrintStruct(value)
	}

	//	Unmarshal el yaml.-
	err = yaml.Unmarshal(byteValue, value)
	logutility.FailOnError(err, "GetConfigurationYaml (Unmarshal):")

	if logEnabled == true {
		logutility.PrintStruct(value)
	}

	return nil
}

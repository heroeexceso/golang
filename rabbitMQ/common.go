package rabbitmq

import "github.com/heroeexceso/golang/ioutility"

// Constantes
const fileNameJSON = "../rabbitMQ/configuration.json"
const fileNameYaml = "../rabbitMQ/configuration.yaml"
const logInfo = false

// ConfigurationRabbitMQ ... estructura cabecera para devolver la configuración de la Cola MQ.-
type ConfigurationRabbitMQ struct {
	Rabbit ConfigurationRabbit `json:"rabbitmq" yaml:"rabbitmq"`
}

// ConfigurationRabbit ... estructura detalle para devolver la configuración de la Cola MQ.-
type ConfigurationRabbit struct {
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Server   string `json:"server" yaml:"server"`
	Port     string `json:"port" yaml:"port"`
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

// GetQueueEndPoint ... sirve para obtener la cadena de conexión para la Cola MQ
func GetQueueEndPoint(isJSON bool) (string, error) {
	var confRabbit ConfigurationRabbitMQ
	var err error
	var strConexion string

	if isJSON == true {
		//	Obtener la configuración del archivo json.-
		err = ioutility.GetConfigurationJSON(fileNameJSON, &confRabbit, logInfo)
	} else {
		err = ioutility.GetConfigurationYaml(fileNameYaml, &confRabbit, logInfo)
	}

	if err == nil {
		if confRabbit.Rabbit.User != "" {
			//	Devolver el string obtenido.-
			strConexion = "amqp://" + confRabbit.Rabbit.User + ":" + confRabbit.Rabbit.Password + "@" + confRabbit.Rabbit.Server + ":" + confRabbit.Rabbit.Port + "/"

			return strConexion, nil
		}
	}

	return "", ioutility.GetExceptionError(err, "GetQueueEndPoint")
}

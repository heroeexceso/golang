package rabbitmq

import (
	"github.com./streadway/amqp"
	"github.com/heroeexceso/golang/logutility"
)

// SendMessage ... enviar un mensaje a la Cola MQ
func SendMessage(colaMQEndPoint string, nameMQ string, exchange string, messageToSend string) bool {
	var rdo bool

	if colaMQEndPoint == "" {
		logutility.PrintMessage("RabbitMQ: No se puedo obtener el string de conexión.")
	} else {
		rdo = true
		//logutility.PrintMessage("RabbitMQ: [String de Conexión: " + colaMQEndPoint + "]")

		//  Establecer la conexión a la cola MQ.-
		conn, err := amqp.Dial(colaMQEndPoint)
		logutility.FailOnError(err, "RabbitMQ: Error al intentar conectarse.")
		defer conn.Close()

		//  Abrir un canal.-
		ch, err := conn.Channel()
		logutility.FailOnError(err, "RabbitMQ: Error al intentar abrir un canal.")
		defer ch.Close()

		//  Configurar variables para crear la cola MQ.-
		name := nameMQ
		durable := false
		autoDelete := false
		exclusive := false
		noWait := false

		//  Declarar la cola MQ.-
		q, err := ch.QueueDeclare(
			name,
			durable,
			autoDelete,
			exclusive,
			noWait,
			nil,
		)
		logutility.FailOnError(err, "RabbitMQ: Error al intentar declarar una Queue.")

		//  Configurar variables para el mensaje a enviar.-
		mandatory := false
		immediate := false

		err = ch.Publish(
			exchange,
			q.Name,
			mandatory,
			immediate,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(messageToSend),
			})
		logutility.FailOnError(err, "RabbitMQ: Error al intentar publicar un mensaje.")
	}

	return rdo
}

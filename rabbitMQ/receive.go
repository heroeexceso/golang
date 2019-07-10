package rabbitmq

import (
	"fmt"

	"github.com./streadway/amqp"
	"github.com/heroeexceso/golang/logutility"
)

// ReceiveOneMessage ... obtener un mensaje de la Cola MQ
func ReceiveOneMessage(colaMQEndPoint string, nameMQ string) (string, bool) {
	var msg string

	if colaMQEndPoint == "" {
		logutility.PrintMessage("RabbitMQ: No se puedo obtener el string de conexión.")
		return "", false
	}

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
	durable := false
	autoDelete := false
	exclusive := false
	noWait := false

	//  Declarar la cola MQ (el último valor son los argumentos).-
	q, err := ch.QueueDeclare(
		nameMQ,
		durable,
		autoDelete,
		exclusive,
		noWait,
		nil,
	)
	logutility.FailOnError(err, "RabbitMQ: Error al intentar declarar una Queue.")

	//  Configurar variables para el mensaje a recibir.-
	autoAck := true

	//	Obtener un único mensaje.-
	msgMQ, rdo, _ := ch.Get(q.Name, autoAck)
	if rdo == true {
		msg = string(msgMQ.Body)
		if len(msg) > 0 {
			return msg, true
		}
	}
	return "", false
}

// ReceiveAllMessages ... otener todos los mensajes de la Cola MQ
func ReceiveAllMessages(colaMQEndPoint string, nameMQ string) (string, bool) {

	msg := "["

	if colaMQEndPoint == "" {
		logutility.PrintMessage("RabbitMQ: No se puedo obtener el string de conexión.")
	} else {
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
		durable := false
		autoDelete := false
		exclusive := false
		noWait := false

		//  Declarar la cola MQ (el último valor son los argumentos).-
		q, err := ch.QueueDeclare(
			nameMQ,
			durable,
			autoDelete,
			exclusive,
			noWait,
			nil,
		)
		logutility.FailOnError(err, "RabbitMQ: Error al intentar declarar una Queue.")

		//  Configurar variables para el mensaje a recibir.-
		autoAck := true

		//	Recorrer la cola para obtener todos los mensajes.-
		for {
			msgMQ, rdo, _ := ch.Get(q.Name, autoAck)
			if rdo == true && string(msgMQ.Body) != "" {
				msg = msg + string(msgMQ.Body) + ","
			} else {
				break
			}
		}
	}

	//	Quitar la última "," y agregar el corchete final.-
	if len(msg) > 1 {
		msg = msg[:len(msg)-1] + "]"
		return msg, true
	}

	return "", false
}

// ReceiveAllMessageChannel ... obtener todos los mensajes de la Cola MQ y quedarse escuchando
func ReceiveAllMessageChannel(colaMQEndPoint string, nameMQ string) bool {
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
		durable := false
		autoDelete := false
		exclusive := false
		noWait := false

		//  Declarar la cola MQ (el último valor son los argumentos).-
		q, err := ch.QueueDeclare(
			nameMQ,
			durable,
			autoDelete,
			exclusive,
			noWait,
			nil,
		)
		logutility.FailOnError(err, "RabbitMQ: Error al intentar declarar una Queue.")

		//  Configurar variables para el mensaje a recibir.-
		consumer := ""
		autoAck := true
		//exclusive := false
		noLocal := false
		//noWait := false

		messageToReceive, err := ch.Consume(
			q.Name,
			consumer,
			autoAck,
			exclusive,
			noLocal,
			noWait,
			nil,
		)
		logutility.FailOnError(err, "RabbitMQ: Error al intentar consumir un mensaje.")

		//	Abrir un canal.-
		forever := make(chan bool)

		go func() {
			//	Recorrer la cola de mensajes.-
			for elemento := range messageToReceive {
				fmt.Println(elemento.Body)
			}
		}()

		//	Enviar contenido del canal.-
		<-forever

		logutility.PrintMessage("Receiving: closing App...")
	}

	return rdo
}

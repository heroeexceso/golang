package rabbitmq

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/heroeexceso/golang/httputility"
	"github.com/heroeexceso/golang/securityutility"
)

// Message ... estructura de mensaje de Cola MQ.-
type Message struct {
	ID       string `json:"id,omitempty"`
	Text     string `json:"text"`
	Consumer string `json:"consumer"`
}

// User ... algo
type User struct {
	UserName     string `json:"username"`
	UserPassword string `json:"userpassword"`
}

// JwtToken ... algo
type JwtToken struct {
	Token string `json:"token"`
}

//  Constantes.-
const port = ":3000"
const logEnabled = false

// Messages ... array para devolver los mensajes
type Messages []Message

// Init ... funcion de inicio de la API
func Init() {

	//  Crear un nuevo router.-
	router := mux.NewRouter()

	//  Configurar los recursos.-

	//	GET.-
	router.HandleFunc("/messageMQ/HeaderInfo", getHeaders).Methods("GET")
	router.HandleFunc("/messageMQ/Messages/queue/{queue}/consumer/{consumer}", getMessages).Methods("GET")
	router.HandleFunc("/messageMQ/Message/queue/{queue}/consumer/{consumer}", getMessage).Methods("GET")
	router.HandleFunc("/messageMQ/protected", getProtected).Methods("GET")

	//	POST.-
	router.HandleFunc("/messageMQ/{id:[0-9]+}/queue/{queue}/consumer/{consumer}", postMessage).Methods("POST")
	router.HandleFunc("/messageMQ/authenticate", postAuthenticate).Methods("POST")

	//	Dejar el servicio escuchando.-
	log.Fatal(http.ListenAndServe(port, router))
}

// postMessage ... enviar un mensaje a la Cola MQ.-
func postMessage(w http.ResponseWriter, r *http.Request) {
	var message Message

	//	Obtener el EndPoint (isJson = true).-
	colaMQEndPoint, err := GetQueueEndPoint(true)

	if err != nil {
		//logutil.PrintMessage("postMessage: No se puedo obtener el string de conexión.")
		httputility.GetJsonResponseMessage(w, "postMessage: No se puedo obtener el string de conexión.")
	} else {

		//  Obtener las variables recibidas.-
		params := mux.Vars(r)
		message.ID = params["id"]
		queueName := params["queue"]
		exchangeName := ""
		message.Consumer = params["consumer"]

		//	Decodificar el mensaje json recibido y dejarlo en la variable de tipo struct.-
		_ = json.NewDecoder(r.Body).Decode(&message)

		//	Realizar el marshal de la variable de tipo struct y guardarlo en []byte.-
		messageToSend, _ := json.Marshal(message)

		//  Enviar mensaje convirtiendo el []byte a string.-
		SendMessage(colaMQEndPoint, queueName, exchangeName, string(messageToSend))

		//	Devolver json .-
		httputility.GetJsonResponseMessage(w, "postMessage: La operación se ha realizado con éxito")
	}
}

// getMessage ... obtener mensaje de la Cola MQ
func getMessage(w http.ResponseWriter, r *http.Request) {
	var message Message

	//	Obtener el EndPoint (isJson = true).-
	colaMQEndPoint, err := GetQueueEndPoint(true)

	if err != nil {
		//logutil.PrintMessage("getMessage: No se puedo obtener el string de conexión.")
		httputility.GetJsonResponseMessage(w, "getMessage: No se puedo obtener el string de conexión.")
	} else {
		//  Obtener las variables recibidas.-
		params := mux.Vars(r)
		queueName := params["queue"]

		//  Recibir mensaje.-
		messageToReceive, rdo := ReceiveOneMessage(colaMQEndPoint, queueName)
		if rdo == true && messageToReceive != "" {
			_ = json.Unmarshal([]byte(messageToReceive), &message)
			//json.NewEncoder(w).Encode(message)
			httputility.GetJsonResponse(w, message)
		} else {
			httputility.GetJsonResponseMessage(w, "getMessages: No se recibieron mensajes.")
		}
	}
}

// getMessages ... obtener mensajes de la Cola MQ
func getMessages(w http.ResponseWriter, r *http.Request) {
	var message Messages

	//	Obtener el EndPoint (isJson = true).-
	colaMQEndPoint, err := GetQueueEndPoint(true)

	if err != nil {
		//logutil.PrintMessage("getMessages: No se puedo obtener el string de conexión.")
		httputility.GetJsonResponseMessage(w, "getMessages: No se puedo obtener el string de conexión.")
	} else {
		//  Obtener las variables recibidas.-
		params := mux.Vars(r)
		queueName := params["queue"]

		//  Recibir mensaje.-
		messageToReceive, rdo := ReceiveAllMessages(colaMQEndPoint, queueName)
		if rdo == true && messageToReceive != "" {
			_ = json.Unmarshal([]byte(messageToReceive), &message)
			json.NewEncoder(w).Encode(message)
		} else {
			httputility.GetJsonResponseMessage(w, "getMessages: No se recibieron mensajes.")
		}
	}
}

// getHeaders ... obtener el header del request
func getHeaders(w http.ResponseWriter, r *http.Request) {
	//	Obtener la información del Request.-
	request := httputility.GetInfoRequest(r)

	//	Devolver el Json con la información.-
	httputility.GetJsonResponse(w, request)
}

//	Sin pruebas ok!!

// postAuthenticate ... algo
func postAuthenticate(w http.ResponseWriter, req *http.Request) {
	var user User
	_ = json.NewDecoder(req.Body).Decode(&user)

	tokenString, rdo := securityutility.GetTokenString(user.UserName, user.UserPassword)
	if rdo == true {
		json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
	} else {
		httputility.GetJsonResponseMessage(w, "Invalid authorization token")
	}
}

// getProtected ... algo
func getProtected(w http.ResponseWriter, req *http.Request) {
	var user User

	params := req.URL.Query()
	tokenString := params["token"][0]

	userName, userPassword, rdo := securityutility.PostProtected(tokenString)

	if rdo == true {
		user.UserName = userName
		user.UserPassword = userPassword
		json.NewEncoder(w).Encode(user)
	} else {
		httputility.GetJsonResponseMessage(w, "Invalid authorization token")
	}
}

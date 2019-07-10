package mongodb

import (
	"context"

	"github.com/heroeexceso/golang/ioutility"
	"github.com/heroeexceso/golang/logutility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConfigurationMongoDB ... estructura cabecera para devolver la configuración de la Cola MQ.-
type ConfigurationMongoDB struct {
	Mongo ConfigurationMongo `json:"mongodb" yaml:"mongodb"`
}

// ConfigurationMongo ... estructura detalle para devolver la configuración de la Cola MQ.-
type ConfigurationMongo struct {
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Server   string `json:"server" yaml:"server"`
	Port     string `json:"port" yaml:"port"`
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

// Constantes
const fileNameJSON = "../mongodb/configuration.json"
const fileNameYaml = "../mongodb/configuration.yaml"
const logInfo = false
const isJSON = true
const cantMaxima = 5

// getDatabaseEndPoint ... sirve para obtener la cadena de conexión a la Base de Datos
func getDatabaseEndPoint() (string, error) {
	var confMongo ConfigurationMongoDB
	var err error
	var strConexion string

	if isJSON == true {
		//	Obtener la configuración del archivo json.-
		err = ioutility.GetConfigurationJSON(fileNameJSON, &confMongo, logInfo)
	} else {
		err = ioutility.GetConfigurationYaml(fileNameYaml, &confMongo, logInfo)
	}

	if err == nil {

		if confMongo.Mongo.User != "" {
			//	Devolver el string obtenido.-
			//strConexion = "mongodb+srv://user:<password>@cluster0-cd4g5.mongodb.net/test?retryWrites=true&w=majority"

			strConexion = "mongodb+srv://" + confMongo.Mongo.User + ":" + confMongo.Mongo.Password + "@" + confMongo.Mongo.Server + "retryWrites=true&w=majority"

			return strConexion, nil
		}
	}

	return "", ioutility.GetExceptionError(err, "GetQueueEndPoint")
}

// ConnectDB ... sirve para realizar la conexión a la Base de Datos
func ConnectDB() (*mongo.Client, error) {
	var client *mongo.Client

	//	Obtener el EndPoint (isJson = true).-
	dbEndPoint, err := getDatabaseEndPoint()

	if err == nil {
		//logutility.PrintMessage(dbEndPoint)

		//	Realizar la conexión a la BD.-
		clientOptions := options.Client().ApplyURI(dbEndPoint)
		client, err = mongo.Connect(context.TODO(), clientOptions)

		if err == nil {
			err = client.Ping(context.TODO(), nil)
			if err == nil {
				//logutility.PrintMessage("Conectado a la Base de Datos MongoDB")

				return client, nil
			}
		}
	}
	logutility.PrintMessage(err.Error())

	return client, ioutility.GetExceptionError(err, "ConnectDB")
}

// GetDocumentByID ... obtener un documento (por ID) de la Base de Datos
func GetDocumentByID(databaseName string, collectionName string, ID string, document interface{}) (interface{}, error) {

	//	Configurar el filtro a utilizar.-
	//filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	filter := bson.D{primitive.E{Key: "brand", Value: "ALCO"}}

	newDocument, err := getDocument(databaseName, collectionName, filter, document)
	if err != nil {
		return nil, err
	}

	return newDocument, nil
}

// getDocument ... obtener un documento de la Base de Datos
func getDocument(databaseName string, collectionName string, filter interface{}, document interface{}) (interface{}, error) {

	//	Conectar a la Base de Datos.-
	client, err := ConnectDB()
	if err != nil {
		return nil, ioutility.GetExceptionError(err, "getDocument.ConnectDB")
	}

	//	Configurar la base de datos y colección a utilizar.-
	collection := client.Database(databaseName).Collection(collectionName)

	//	Configurar el filtro a utilizar.-

	//	Obtener el documento.-
	err = collection.FindOne(context.TODO(), filter).Decode(&document)
	if err != nil {
		return nil, ioutility.GetExceptionError(err, "getDocument.FindOne")
	}

	//	Desconectar de la Base de Datos.-
	err = client.Disconnect(context.TODO())
	if err != nil {
		return nil, ioutility.GetExceptionError(err, "getDocument.Disconnect")
	}
	return document, nil
}

// GetDocuments ... obtener un documento de la Base de Datos
func GetDocuments(databaseName string, collectionName string, document interface{}) ([]interface{}, error) {
	var curr *mongo.Cursor

	//	Conectar a la Base de Datos.-
	client, err := ConnectDB()
	if err != nil {
		return nil, ioutility.GetExceptionError(err, "GetDocuments.ConnectDB")
	}

	//	Configurar la base de datos y colección a utilizar.-
	collection := client.Database(databaseName).Collection(collectionName)

	//	Configurar el filtro a utilizar.-
	findOptions := options.Find()
	findOptions.SetLimit(cantMaxima)

	//	Obtener el documento.-
	curr, err = collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		return nil, ioutility.GetExceptionError(err, "GetDocuments.Find")
	}

	//	Completar el array.-
	var documents []interface{}
	for curr.Next(context.TODO()) {
		err = curr.Decode(&document)
		if err == nil {
			documents = append(documents, document)
		}
	}

	//	Cerrar el cursor.-
	curr.Close(context.TODO())

	//	Desconectar de la Base de Datos.-
	err = client.Disconnect(context.TODO())
	if err != nil {
		return nil, ioutility.GetExceptionError(err, "GetDocuments.Disconnect")
	}

	return documents, nil
}

// InsertDocument ... insertar un documento en la Base de Datos
func InsertDocument(databaseName string, collectionName string, document interface{}) error {

	//	Conectar a la Base de Datos.-
	client, err := ConnectDB()
	if err != nil {
		return ioutility.GetExceptionError(err, "InsertDocument.ConnectDB")
	}

	//	Configurar la base de datos y colección a utilizar.-
	collection := client.Database(databaseName).Collection(collectionName)

	//	Insertar el nuevo valor en la base de datos.-
	_, err = collection.InsertOne(context.TODO(), document)
	//insertResult, err := collection.InsertOne(context.TODO(), document)
	if err != nil {
		return ioutility.GetExceptionError(err, "InsertDocument.InsertOne")
	}

	//	Obtener el ID generado.-
	//ID := insertResult.InsertedID

	//	Desconectar de la Base de Datos.-
	err = client.Disconnect(context.TODO())
	if err != nil {
		return ioutility.GetExceptionError(err, "InsertDocument.Disconnect")
	}
	return nil
}

// UpdateDocumentByID ... actualizar un documento (por ID) de la Base de Datos
func UpdateDocumentByID(databaseName string, collectionName string, ID string, document interface{}) (int64, error) {

	//	Configurar el filtro a utilizar.-
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	update := bson.D{primitive.E{Key: "model", Value: "Coupe Megane"}, {Key: "brand", Value: "Renault"}, {Key: "year", Value: 2000}}

	//"$set"

	matchedCount, err := updateDocument(databaseName, collectionName, filter, update)
	if err != nil {
		return 0, ioutility.GetExceptionError(err, "UpdateDocumentById.UpdateDocument")
	}

	return matchedCount, err
}

// updateDocument ... actualizar un documento en la Base de Datos
func updateDocument(databaseName string, collectionName string, filter interface{}, update interface{}) (int64, error) {
	var updateResult *mongo.UpdateResult

	//	Conectar a la Base de Datos.-
	client, err := ConnectDB()
	if err != nil {
		return 0, ioutility.GetExceptionError(err, "UpdateDocument.ConnectDB")
	}

	//	Configurar la base de datos y colección a utilizar.-
	collection := client.Database(databaseName).Collection(collectionName)

	//	Insertar el nuevo valor en la base de datos.-
	updateResult, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return 0, ioutility.GetExceptionError(err, "UpdateDocument.UpdateOne")
	}

	matchedCount := updateResult.MatchedCount

	//	Desconectar de la Base de Datos.-
	err = client.Disconnect(context.TODO())
	if err != nil {
		return 0, ioutility.GetExceptionError(err, "UpdateDocument.Disconnect")
	}
	return matchedCount, nil
}

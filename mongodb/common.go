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
const cantMaxima = 10

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

	return "", logutility.GetExceptionError(err, "GetQueueEndPoint")
}

// ConnectDB ... sirve para realizar la conexión a la Base de Datos
func connectDB() (*mongo.Client, error) {
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

	return client, logutility.GetExceptionError(err, "ConnectDB")
}

// GetDocuments ... todos los documentos de la Base de Datos
func GetDocuments(databaseName string, collectionName string, filter interface{}, document interface{}) ([]interface{}, error) {
	var curr *mongo.Cursor

	//	Conectar a la Base de Datos.-
	client, err := connectDB()
	if err != nil {
		return nil, logutility.GetExceptionError(err, "GetDocuments.ConnectDB")
	}

	//	Configurar la base de datos y colección a utilizar.-
	collection := client.Database(databaseName).Collection(collectionName)

	//	Configurar el filtro a utilizar.-
	findOptions := options.Find()
	findOptions.SetLimit(cantMaxima)

	//	Obtener el documento.-
	//curr, err = collection.Find(context.TODO(), bson.D{{}}, findOptions)
	curr, err = collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, logutility.GetExceptionError(err, "GetDocuments.Find")
	}

	//	Completar el array.-
	var documents []interface{}

	incluteItems := false

	for curr.Next(context.TODO()) {
		if incluteItems == false {
			incluteItems = true
		}
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
		return nil, logutility.GetExceptionError(err, "GetDocuments.Disconnect")
	}

	if incluteItems == false {
		return nil, logutility.GetExceptionText("No se han encontrado documentos con los criterios recibidos.")
	}

	return documents, nil
}

// InsertDocument ... insertar un documento en la Base de Datos
func InsertDocument(databaseName string, collectionName string, document interface{}) error {

	//	Conectar a la Base de Datos.-
	client, err := connectDB()
	if err != nil {
		return logutility.GetExceptionError(err, "InsertDocument.ConnectDB")
	}

	//	Configurar la base de datos y colección a utilizar.-
	collection := client.Database(databaseName).Collection(collectionName)

	//	Insertar el nuevo valor en la base de datos.-
	_, err = collection.InsertOne(context.TODO(), document)
	//insertResult, err := collection.InsertOne(context.TODO(), document)
	if err != nil {
		return logutility.GetExceptionError(err, "InsertDocument.InsertOne")
	}

	//	Obtener el ID generado.-
	//ID := insertResult.InsertedID

	//	Desconectar de la Base de Datos.-
	err = client.Disconnect(context.TODO())
	if err != nil {
		return logutility.GetExceptionError(err, "InsertDocument.Disconnect")
	}
	return nil
}

// UpdateDocument ... actualizar un documento en la Base de Datos
func UpdateDocument(databaseName string, collectionName string, filter interface{}, update interface{}) (int64, error) {
	var updateResult *mongo.UpdateResult

	//	Conectar a la Base de Datos.-
	client, err := connectDB()
	if err != nil {
		return 0, logutility.GetExceptionError(err, "UpdateDocument.ConnectDB")
	}

	//	Configurar la base de datos y colección a utilizar.-
	collection := client.Database(databaseName).Collection(collectionName)

	newUpdate := bson.D{primitive.E{"$set", update}}

	//	Actualizar el valor en la base de datos.-
	updateResult, err = collection.UpdateOne(context.TODO(), filter, newUpdate)
	if err != nil {
		return 0, logutility.GetExceptionError(err, "UpdateDocument.UpdateOne")
	}

	updatedCount := updateResult.MatchedCount

	//	Desconectar de la Base de Datos.-
	err = client.Disconnect(context.TODO())
	if err != nil {
		return 0, logutility.GetExceptionError(err, "UpdateDocument.Disconnect")
	}
	return updatedCount, nil
}

// DeleteDocument ... actualizar un documento en la Base de Datos
func DeleteDocument(databaseName string, collectionName string, filter interface{}) (int64, error) {
	var deleteResult *mongo.DeleteResult

	//	Conectar a la Base de Datos.-
	client, err := connectDB()
	if err != nil {
		return 0, logutility.GetExceptionError(err, "DeleteDocument.ConnectDB")
	}

	//	Configurar la base de datos y colección a utilizar.-
	collection := client.Database(databaseName).Collection(collectionName)

	//	Eliminar el valor de la base de datos.-
	deleteResult, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return 0, logutility.GetExceptionError(err, "DeleteDocument.UpdateOne")
	}

	deletedCount := deleteResult.DeletedCount

	//	Desconectar de la Base de Datos.-
	err = client.Disconnect(context.TODO())
	if err != nil {
		return 0, logutility.GetExceptionError(err, "DeleteDocument.Disconnect")
	}
	return deletedCount, nil
}

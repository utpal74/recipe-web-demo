package bootstrap

import (
	"log"

	inframongo "github.com/gin-demo/recipes-web/internal/shared/infra/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoResources struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func InitMongo(uri, dbName string) (*MongoResources, error) {
	client, err := inframongo.NewClient(uri)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	log.Println("Mongo connected:", dbName)

	return &MongoResources{
		Client:   client,
		Database: db,
	}, nil
}

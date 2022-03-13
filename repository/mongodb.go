package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/SmartTriageFiap/notification/config"
	"github.com/SmartTriageFiap/notification/entities"
	"github.com/arturmartini/envconfig"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

const (
	database           = "hvm"
	collectionPartners = "partners"
	limitTimeout       = 30 * time.Second
)

var (
	once                  sync.Once
	repo                  *mongoRepository
	ErrorPartneraNotFound = errors.New("partners not found")
)

type mongoRepository struct {
	client *mongo.Client
}

func newMongoRepository() Repository {
	once.Do(func() {
		if repo == nil {
			repo = &mongoRepository{
				client: connect(),
			}
		}
	})
	return repo
}

func connect() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), limitTimeout)
	defer cancel()

	urlCredential := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		envconfig.GetStr(config.MongoUser),
		envconfig.GetStr(config.MongoPass),
		envconfig.GetStr(config.MongoAddr),
		envconfig.GetStr(config.MongoPort))

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(urlCredential))
	if err != nil {
		client.Disconnect(ctx)
		log.WithError(err).Panic("error when try connect to database")
	}
	return client
}

func (r mongoRepository) GetPartnersAddresses() ([]entities.Partner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), limitTimeout)
	defer cancel()
	partners := []entities.Partner{}
	cursor, err := r.client.Database(database).Collection(collectionPartners).Find(ctx, bson.D{})
	defer cursor.Close(ctx)
	if err != nil {
		log.WithError(err).Warn("there was error in find partners")
		return []entities.Partner{}, ErrorPartneraNotFound
	}

	err = cursor.All(ctx, &partners)
	if err != nil {
		log.WithError(err).Warn("there was error in cursor binding all")
		return []entities.Partner{}, ErrorPartneraNotFound
	}

	return partners, nil
}

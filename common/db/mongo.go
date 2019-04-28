package db

import (
	"context"
	"errors"

	"time"

	"github.com/kuritka/break-down.io/common/data"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type db struct {
	client  *mongo.Client
	context context.Context
	options *ClientOptions
}

func (db db) Upsert(calendar data.Calendar) (interface{}, error) {

	if calendar.Name == "" {
		err := errors.New("missing calendar name")
		log.Fatal().Err(err).Msgf("missing calendar name")
	}

	collection := db.client.Database(db.options.Database).Collection(db.options.Collection)
	existing, err := db.GetByName(calendar.Name)
	if err != nil {
		return "", err
	}

	if existing != nil {
		filter := bson.M{"_id": calendar.Name}
		_, err := collection.ReplaceOne(db.context, filter, calendar)
		if err != nil {
			return nil, err
		}
		return calendar.Name, nil
	}
	res, err := collection.InsertOne(db.context, calendar)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

func (db db) GetByName(name string) (*data.Calendar, error) {
	var calendar data.Calendar

	col := db.client.Database(db.options.Database).Collection(db.options.Collection)
	filter := bson.M{"_id": name}
	err := col.FindOne(db.context, filter).Decode(&calendar)
	if err != nil {
		return nil, nil
	}
	return &calendar, nil
}

func (m *mongoImpl) Connection() db {
	ctx, _ := context.WithTimeout(context.Background(), m.Options.Timeout*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.Options.ConnectionString))
	if err != nil {
		log.Fatal().Err(err).Msgf("can't connect to %s:%s", m.Options.Database, m.Options.Collection)
	}
	return db{client, ctx, &m.Options}
}

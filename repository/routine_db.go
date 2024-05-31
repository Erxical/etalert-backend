package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type routineRepositoryDB struct {
	collection *mongo.Collection
}

func NewRoutineRepositoryDB(client *mongo.Client, dbName string, collName string) RoutineRepository {
	collection := client.Database(dbName).Collection(collName)
	return &routineRepositoryDB{collection: collection}
}

func (r *routineRepositoryDB) GetHighestOrder(gId string) (int, error) {
	ctx := context.Background()
	filter := bson.M{"googleId": gId}
	opts := options.FindOne().SetSort(bson.M{"order": -1})
	var routine Routine
	err := r.collection.FindOne(ctx, filter, opts).Decode(&routine)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, err
	}
	return routine.Order, nil
}

func (r *routineRepositoryDB) InsertRoutine(routine *Routine) error {
	ctx := context.Background()
	_, err := r.collection.InsertOne(ctx, routine)
	return err
}

func (r *routineRepositoryDB) GetRoutine(gId string) (*Routine, error) {
	ctx := context.Background()
	var routine Routine
	filter := bson.M{"googleId": gId}
	err := r.collection.FindOne(ctx, filter).Decode(&routine)
	if err != nil {
		return nil, err
	}
	return &routine, nil

}

func (r *routineRepositoryDB) UpdateRoutine(gId string, routine *Routine) error {
	ctx := context.Background()
	filter := bson.M{"googleId": gId}
	update := bson.M{
		"$set": bson.M{
			"name":     routine.Name,
			"duration": routine.Duration,
			"order":    routine.Order,
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

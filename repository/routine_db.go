package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type routineRepositoryDB struct {
	collection *mongo.Collection
}

func NewRoutineRepositoryDB(client *mongo.Client, dbName string, collName string) RoutineRepository {
	collection := client.Database(dbName).Collection(collName)
	return &routineRepositoryDB{collection: collection}
}

func (r *routineRepositoryDB) InsertRoutine(routine *Routine) error {
	ctx := context.Background()
	_, err := r.collection.InsertOne(ctx, routine)
	return err
}

func (r *routineRepositoryDB) GetAllRoutines(gId string) ([]*Routine, error) {
    ctx := context.Background()
    routines := []*Routine{}
    filter := bson.M{"googleId": gId}

    // Use Find to get all matching documents
    cursor, err := r.collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    // Ensure the cursor is closed once we're done
    defer cursor.Close(ctx)

    // Iterate through the cursor
    for cursor.Next(ctx) {
        var routine Routine
        // Decode each document into a Routine struct
        if err := cursor.Decode(&routine); err != nil {
            return nil, err
        }
        // Append the decoded routine to the slice
        routines = append(routines, &routine)
    }

    // Check if any errors occurred during iteration
    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return routines, nil
}

func (r *routineRepositoryDB) GetRoutineById(id string) (*Routine, error) {
	ctx := context.Background()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("failed to convert ID: %v", err)
	}
	filter := bson.M{"_id": objectId}
	
	routine := &Routine{}
	err = r.collection.FindOne(ctx, filter).Decode(routine)
	if err != nil {
		return nil, err
	}
	return routine, nil
}

func (r *routineRepositoryDB) UpdateRoutine(id string, routine *Routine) error {
	ctx := context.Background()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert ID: %v", err)
	}
	filter := bson.M{"_id": (objectId)}
	update := bson.M{
		"$set": bson.M{
			"name":     routine.Name,
			"duration": routine.Duration,
			"order":    routine.Order,
		},
	}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *routineRepositoryDB) DeleteRoutine(id string) error {
	ctx := context.Background()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert ID: %v", err)
	}
	filter := bson.M{"_id": objectId}
	_, err = r.collection.DeleteOne(ctx, filter)
	return err
}
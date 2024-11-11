package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type tagRepositoryDB struct {
	collection *mongo.Collection
}

func NewTagRepositoryDB(client *mongo.Client, dbName string, collName string) TagRepository {
	collection := client.Database(dbName).Collection(collName)
	return &tagRepositoryDB{collection: collection}
}

func (t *tagRepositoryDB) InsertTag(tag *Tag) error {
	ctx := context.Background()
	_, err := t.collection.InsertOne(ctx, tag)
	return err
}

func (t *tagRepositoryDB) GetAllTags(gId string) ([]*Tag, error) {
	ctx := context.Background()
	tags := []*Tag{}
	filter := bson.M{"googleId": gId}

	// Use Find to get all matching documents
	cursor, err := t.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	// Ensure the cursor is closed once we're done
	defer cursor.Close(ctx)

	// Iterate through the cursor
	for cursor.Next(ctx) {
		var tag Tag
		// Decode each document into a Routine struct
		if err := cursor.Decode(&tag); err != nil {
			return nil, err
		}
		// Append the decoded routine to the slice
		tags = append(tags, &tag)
	}

	// Check if any errors occurred during iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (t *tagRepositoryDB) GetRoutinesByTagId(id string) ([]string, error) {
	ctx := context.Background()

	if id == "" {
		return []string{}, fmt.Errorf("id is required")
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("failed to convert ID: %v", err)
	}

	filter := bson.M{"_id": objectId}
	
	var tag Tag
	err = t.collection.FindOne(ctx, filter).Decode(&tag)
	if err != nil {
		return nil, err
	}
	return tag.Routines, nil
}

func (t *tagRepositoryDB) GetTagByRoutineId(id string) (*Tag, error) {
	ctx := context.Background()

	filter := bson.M{"routines": id}

	var tag Tag
	err := t.collection.FindOne(ctx, filter).Decode(&tag)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (t *tagRepositoryDB) UpdateTag(id string, tag *Tag) error {
	ctx := context.Background()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert ID: %v", err)
	}
	filter := bson.M{"_id": (objectId)}
	update := bson.M{
		"$set": bson.M{
			"name":     tag.Name,
			"routines": tag.Routines,
		},
	}
	_, err = t.collection.UpdateOne(ctx, filter, update)
	return err
}

func (t *tagRepositoryDB) DeleteTag(id string) error {
	ctx := context.Background()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert ID: %v", err)
	}
	filter := bson.M{"_id": objectId}
	_, err = t.collection.DeleteOne(ctx, filter)
	return err
}

package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID string `bson:"_id,omitemtpy" json:"id,omitempty"`
	Name string `bson:"name" json:"name"`
	Data string `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func New(mongo *mongo.Client) Models {
	client = mongo


	return Models {
		LogEntry : LogEntry{},
	}
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs");

	_ , err := collection.InsertOne(context.TODO(),LogEntry {
		Name : entry.Name, 
		Data : entry.Data,
		CreatedAt : time.Now(),
		UpdatedAt : time.Now(),
	})

	if err != nil {
		log.Println("Error inserting",err);
		return err
	}

	return nil
}

func (l *LogEntry) All() ([]*LogEntry,error) {
	ctx , cancel := context.WithTimeout(context.Background(),15*time.Second);
	defer cancel();

	collection := client.Database("logs").Collection("logs");

	opts := options.Find();
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	
	cursor, err := collection.Find(context.TODO(),bson.D{},opts);
	if err != nil {
		log.Println("Finding all docs faild",err);
		return nil,err
	}
	defer cursor.Close(ctx);

	var logs []*LogEntry;

	for cursor.Next(ctx) {
		var item LogEntry
		err := cursor.Decode(&item)
		if err != nil {
			log.Println("Error in decoding doc",err);
			return nil,err
		} else {
			logs = append(logs,&item);
		}
	}

	if err := cursor.Err(); err != nil {
        log.Println("Cursor error", err)
        return nil, err
    }

	return logs,nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry,error) {
	ctx , cancel := context.WithTimeout(context.Background(),15*time.Second);
	defer cancel();

	collection := client.Database("logs").Collection("logs");

	docID,err := primitive.ObjectIDFromHex(id);
	if err != nil {
		log.Println("Failed to convert to ObjID",err);
		return nil,err
	}

	var entry LogEntry
	err = collection.FindOne(ctx,bson.M{"_id": docID}).Decode(&entry);
	if err != nil {
		log.Println("Failed to fetch from mongo",err);
		return nil,err
	}

	return &entry,nil
}

func (l *LogEntry) DropCollection() error {
	ctx , cancel := context.WithTimeout(context.Background(),15*time.Second);
	defer cancel();

	collection := client.Database("logs").Collection("logs");

	if err:=collection.Drop(ctx); err!=nil {
		log.Println("Error in dropping collection",err)
		return err
	}

	return nil
}

func (l *LogEntry) Update(id string) (*mongo.UpdateResult,error) {
	ctx , cancel := context.WithTimeout(context.Background(),15*time.Second);
	defer cancel();

	collection := client.Database("logs").Collection("logs");

	docID,err := primitive.ObjectIDFromHex(l.ID);
	if err != nil {
		log.Println("Failed to convert to ObjID",err);
		return nil,err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: l.Name},
				{Key: "data", Value: l.Data},
				{Key: "updated_at", Value: time.Now()},
			}},
		},
	)

	if err!=nil {
		return nil,err
	}

	return result,nil
}


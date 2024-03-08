package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MondoDB struct {
	DB         *mongo.Client
	Collection *mongo.Collection
}

func NewMongoDB() *MondoDB {
	return &MondoDB{}
}

func (s *MondoDB) ConnToDB() {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	s.DB, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Failed to connect MondoDB: " + err.Error())
	}
	fmt.Println("Successfully connected to MondoDB")
	s.Collection = s.DB.Database("testezhik").Collection("users")
}

func (s *MondoDB) UserExist(username string, user interface{}) error {
	filter := bson.M{"username": username}
	err := s.Collection.FindOne(context.Background(), filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil // User with the given username does not exist
		}
		return err // Some other error occurred
	}
	return fmt.Errorf("user with username: %s already exist", username) // User exist
}

func (s *MondoDB) EmailExist(email string, user interface{}) error {
	filter := bson.M{"email": email}
	err := s.Collection.FindOne(context.Background(), filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil // User with the given username does not exist
		}
		return err // Some other error occurred
	}
	return fmt.Errorf("user with email: %s already exist", email) // User exist
}

func (s *MondoDB) SaveUser(user interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := s.Collection.InsertOne(ctx, user)
	return err
}

func (s *MondoDB) GetUserByUsername(username string, user interface{}) error {
	filter := bson.M{"username": username}
	err := s.Collection.FindOne(context.Background(), filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("user with username '%s' not found", username)
		}
		return err // Some other error occurred
	}
	return nil
}

func (s *MondoDB) GetUserByEmail(email string, user interface{}) error {
	filter := bson.M{"email": email}
	err := s.Collection.FindOne(context.Background(), filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("user with email '%s' not found", email)
		}
		return err // Some other error occurred
	}
	return nil
}

func (s *MondoDB) GetUserByID(id any, user interface{}) error {
	var oid primitive.ObjectID
	var err error
	if oid, err = primitive.ObjectIDFromHex(id.(string)); err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	err = s.Collection.FindOne(context.Background(), filter).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("user with objectID '%v' not found", oid)
		}
		return err // Some other error occurred
	}
	return nil
}

func (s *MondoDB) VerifyAccount(email string, verTime sql.NullTime, user interface{}) error {
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"verifiedAt": verTime}}
	_, err := s.Collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (s *MondoDB) UpdatePswdHash(newPswdHash string, id any, user interface{}) error {
	var oid primitive.ObjectID
	var err error
	if oid, err = primitive.ObjectIDFromHex(id.(string)); err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	update := bson.M{"$set": bson.M{"pswdHash": newPswdHash}}
	_, err = s.Collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (s *MondoDB) UpdateVerHash(newVerHash string, t time.Time, id any, user interface{}) error {
	var oid primitive.ObjectID
	var err error
	if oid, err = primitive.ObjectIDFromHex(id.(string)); err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	update := bson.M{"$set": bson.M{"verHash": newVerHash, "timeout": t}}
	_, err = s.Collection.UpdateOne(context.Background(), filter, update)
	return err
}

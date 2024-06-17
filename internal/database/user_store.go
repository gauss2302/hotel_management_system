package database

import (
	"context"
	"github.com/gauss2302/hotel_management_system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColl = "users"

type Map map[string]any

type UserStore interface {
	GetUserByID(context.Context, string) (*models.User, error)
	GetUserByEmail(context.Context, string) (*models.User, error)
	GetUsers(context.Context) ([]*models.User, error)
	InsertUser(context.Context, *models.User) (*models.User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(ctx context.Context, filter Map, params models.UpdateUserParams) error
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database("MONGO_DB_NAME").Collection(userColl),
	}
}

func (s *MongoUserStore) InsertUser(ctx context.Context, user *models.User) (*models.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*models.User, error) {
	cur, err := s.coll.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	var users []*models.User

	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *MongoUserStore) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	var user models.User

	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil

}

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := s.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, filter Map, params models.UpdateUserParams) error {
	oid, err := primitive.ObjectIDFromHex(filter["_id"].(string))
	if err != nil {
		return err
	}

	filter["_id"] = oid
	update := bson.M{"$set": params.ToBSON()}
	_, err = s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	return nil
}

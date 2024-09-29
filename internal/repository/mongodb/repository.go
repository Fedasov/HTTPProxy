package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"Proxy/internal/domain/models"
)

type Repo struct {
	Collection *mongo.Collection
}

func NewRequestRepository(collection *mongo.Collection) *Repo {
	return &Repo{
		Collection: collection,
	}
}

func (r *Repo) GetRequestByID(ctx context.Context, id primitive.ObjectID) (*models.Request, error) {
	var request models.Request
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *Repo) GetAllRequests(ctx context.Context) ([]models.Request, error) {

	cursor, err := r.Collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var requests []models.Request
	if err := cursor.All(ctx, &requests); err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *Repo) AddRequestResponse(ctx context.Context, req models.ParsedRequest, resp models.ParsedResponse) (primitive.ObjectID, error) {
	rec := models.Request{
		Request:   req,
		Response:  resp,
		CreatedAt: time.Now(),
	}
	res, err := r.Collection.InsertOne(ctx, rec)
	if err != nil {
		return primitive.NilObjectID, err
	}

	result, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.ObjectID{}, fmt.Errorf("Error")
	}

	return result, nil
}

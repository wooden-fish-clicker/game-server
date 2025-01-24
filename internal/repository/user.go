package repository

import (
	"context"
	"game-server/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) *UserRepository {
	collection := database.Collection("users")
	return &UserRepository{collection: collection}
}

func (u *UserRepository) GetDeteil(ctx context.Context, user *models.User) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	err = u.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) UpdateUserState(ctx context.Context, userId string, hp int, points int64) error {
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	now := time.Now()
	filter := bson.M{"_id": objectId}
	update := bson.M{
		"$set": bson.M{
			"user_info.points": points,
			"user_info.hp":     hp,
			"updated_at":       now,
		},
	}

	result, err := u.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return err
	}

	return nil
}

package repositories

import (
	"7hunt-be-rest-api/internal/core/domain"
	"7hunt-be-rest-api/internal/core/port/output"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	mongoCollection *mongo.Collection
}

func NewUserRepository(mongoCollection *mongo.Database, collectionName string) output.UserRepository {
	return &userRepository{
		mongoCollection: mongoCollection.Collection("User"),
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	user.ID = primitive.NewObjectID()
	_, err := r.mongoCollection.InsertOne(ctx, user)
	return err
}

func (r *userRepository) GetUsers(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User
	cursor, err := r.mongoCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	var user domain.User
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	err = r.mongoCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, userID string, user *domain.User) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return 0, err
	}

	update := bson.M{
		"$set": bson.M{
			"name":  user.Name,
			"email": user.Email,
		},
	}

	result, err := r.mongoCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return 0, err
	}

	return result.MatchedCount, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, userID string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return 0, err
	}

	result, err := r.mongoCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.mongoCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.mongoCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

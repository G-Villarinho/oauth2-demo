package repositories

import (
	"context"
	"time"

	"github.com/aetheris-lab/aetheris-id/api/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OTPRepository interface {
	Create(ctx context.Context, otp *entities.OTP) error
	FindByID(ctx context.Context, id string) (*entities.OTP, error)
	Delete(ctx context.Context, id string) error
}

type otpRepository struct {
	collection *mongo.Collection
}

func NewOTPRepository(db *mongo.Database) OTPRepository {
	return &otpRepository{
		collection: db.Collection("otps"),
	}
}

func (r *otpRepository) Create(ctx context.Context, otp *entities.OTP) error {
	if otp.ID.IsZero() {
		otp.ID = primitive.NewObjectID()
	}

	otp.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, otp)
	if err != nil {
		return err
	}

	return nil
}

func (r *otpRepository) FindByID(ctx context.Context, id string) (*entities.OTP, error) {
	var otp entities.OTP

	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&otp)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, entities.ErrOTPNotFound
		}

		return nil, err
	}

	return &otp, nil
}

func (r *otpRepository) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/harih/pulse-poll/backend/models"
)

type PollService struct {
	Polls *mongo.Collection
	Redis *RedisService
}

func NewPollService(polls *mongo.Collection, redisService *RedisService) *PollService {
	return &PollService{
		Polls: polls,
		Redis: redisService,
	}
}

func (s *PollService) CreatePoll(
	ctx context.Context,
	question string,
	options []string,
	createdBy interface{},
) (*models.Poll, error) {

	question = strings.TrimSpace(question)

	if question == "" {
		return nil, errors.New("question is required")
	}

	if len(question) > 300 {
		return nil, errors.New("question must be 300 characters or less")
	}

	if len(options) < 2 || len(options) > 10 {
		return nil, errors.New("poll must have between 2 and 10 options")
	}

	pollOptions := make([]models.PollOption, 0, len(options))
	seen := make(map[string]bool)

	for _, option := range options {
		option = strings.TrimSpace(option)

		if option == "" {
			return nil, errors.New("options cannot be empty")
		}

		if len(option) > 100 {
			return nil, errors.New("each option must be 100 characters or less")
		}

		key := strings.ToLower(option)

		if seen[key] {
			return nil, errors.New("duplicate options are not allowed")
		}

		seen[key] = true

		pollOptions = append(pollOptions, models.PollOption{
			ID:   uuid.NewString(),
			Text: option,
		})
	}

	poll := &models.Poll{
		Question:  question,
		Options:   pollOptions,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
		IsActive:  true,
	}

	result, err := s.Polls.InsertOne(ctx, poll)
	if err != nil {
		return nil, err
	}

	poll.ID = result.InsertedID

	pollID := result.InsertedID.(bson.ObjectID).Hex()

	// Redis needs only the option IDs to initialize vote counts.
	optionIDs := make([]string, 0, len(pollOptions))

	for _, option := range pollOptions {
		optionIDs = append(optionIDs, option.ID)
	}

	if err := s.Redis.InitializePoll(ctx, pollID, optionIDs); err != nil {
		// Remove the MongoDB poll if Redis initialization fails.
		_, _ = s.Polls.DeleteOne(ctx, bson.M{"_id": result.InsertedID})
		return nil, err
	}

	return poll, nil
}

func (s *PollService) GetPoll(
	ctx context.Context,
	id string,
) (*models.Poll, error) {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid poll id")
	}

	var poll models.Poll

	err = s.Polls.FindOne(
		ctx,
		bson.M{"_id": objectID},
	).Decode(&poll)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("poll not found")
		}

		return nil, err
	}

	return &poll, nil
}
func (s *PollService) GetPollResults(ctx context.Context, pollID string) (map[string]int, error) {
	return s.Redis.GetResults(ctx, pollID)
}

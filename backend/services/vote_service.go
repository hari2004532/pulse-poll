package services

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/harih/pulse-poll/backend/models"
)

type VoteService struct {
	Votes *mongo.Collection
	Polls *mongo.Collection
	Redis *RedisService
}

func NewVoteService(
	votes *mongo.Collection,
	polls *mongo.Collection,
	redisService *RedisService,
) *VoteService {
	return &VoteService{
		Votes: votes,
		Polls: polls,
		Redis: redisService,
	}
}

func (s *VoteService) Vote(
	ctx context.Context,
	pollID string,
	optionID string,
	voterID interface{},
) (map[string]int, error) {

	objectID, err := bson.ObjectIDFromHex(pollID)
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

	if !poll.IsActive {
		return nil, errors.New("poll is closed")
	}

	validOption := false

	for _, option := range poll.Options {
		if option.ID == optionID {
			validOption = true
			break
		}
	}

	if !validOption {
		return nil, errors.New("invalid option")
	}

	// Prevent the same user from voting more than once.
	var existingVote models.Vote

	err = s.Votes.FindOne(
		ctx,
		bson.M{
			"pollId":  objectID,
			"voterId": voterID,
		},
	).Decode(&existingVote)

	if err == nil {
		return nil, errors.New("you have already voted")
	}

	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	vote := &models.Vote{
		PollID:    objectID,
		OptionID:  optionID,
		VoterID:   voterID,
		CreatedAt: time.Now(),
	}

	_, err = s.Votes.InsertOne(ctx, vote)
	if err != nil {
		return nil, err
	}

	results, err := s.Redis.IncrementVote(
		ctx,
		pollID,
		optionID,
	)

	if err != nil {
		return nil, err
	}

	if err := s.Redis.PublishVoteUpdate(
		ctx,
		pollID,
		results,
	); err != nil {
		return nil, err
	}

	return results, nil
}

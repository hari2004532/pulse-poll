package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	Client *redis.Client
}

type VoteUpdate struct {
	PollID  string         `json:"pollId"`
	Results map[string]int `json:"results"`
}

func NewRedisService(addr, password string) *RedisService {
	return &RedisService{
		Client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       0,
		}),
	}
}

func (r *RedisService) IncrementVote(
	ctx context.Context,
	pollID string,
	optionID string,
) (map[string]int, error) {

	key := fmt.Sprintf("poll:%s:results", pollID)

	_, err := r.Client.HIncrBy(
		ctx,
		key,
		optionID,
		1,
	).Result()

	if err != nil {
		return nil, err
	}

	return r.GetResults(ctx, pollID)
}

func (r *RedisService) GetResults(
	ctx context.Context,
	pollID string,
) (map[string]int, error) {

	key := fmt.Sprintf("poll:%s:results", pollID)

	values, err := r.Client.HGetAll(ctx, key).Result()

	if err != nil {
		return nil, err
	}

	results := make(map[string]int)

	for optionID, value := range values {
		count, err := strconv.Atoi(value)

		if err != nil {
			continue
		}

		results[optionID] = count
	}

	return results, nil
}

func (r *RedisService) PublishVoteUpdate(
	ctx context.Context,
	pollID string,
	results map[string]int,
) error {

	channel := fmt.Sprintf("poll:%s:updates", pollID)

	update := VoteUpdate{
		PollID:  pollID,
		Results: results,
	}

	data, err := json.Marshal(update)

	if err != nil {
		return err
	}

	return r.Client.Publish(
		ctx,
		channel,
		data,
	).Err()
}

func (r *RedisService) InitializePoll(
	ctx context.Context,
	pollID string,
	optionIDs []string,
) error {

	key := fmt.Sprintf("poll:%s:results", pollID)

	values := make(map[string]interface{})

	for _, optionID := range optionIDs {
		values[optionID] = 0
	}

	if len(values) == 0 {
		return nil
	}

	return r.Client.HSet(ctx, key, values).Err()
}

package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/graduate-work-mirea/api-gateway/config"
	"github.com/graduate-work-mirea/api-gateway/model"
	lru "github.com/hashicorp/golang-lru"
)

// CacheRepository represents a cache repository for prediction data
type CacheRepository interface {
	SavePrediction(userID uuid.UUID, prediction model.PredictionHistory) error
	GetUserPredictions(userID uuid.UUID) ([]model.PredictionHistory, bool)
	PopulateFromMap(predictions map[uuid.UUID][]model.PredictionHistory)
}

type lruCacheRepository struct {
	cache     *lru.Cache
	userCache map[uuid.UUID][]model.PredictionHistory
	mutex     sync.RWMutex
}

// NewCacheRepository creates a new cache repository
func NewCacheRepository(cfg *config.Config) (CacheRepository, error) {
	// Create LRU cache with the configured size
	cache, err := lru.New(cfg.CacheSize)
	if err != nil {
		return nil, err
	}

	return &lruCacheRepository{
		cache:     cache,
		userCache: make(map[uuid.UUID][]model.PredictionHistory),
		mutex:     sync.RWMutex{},
	}, nil
}

// SavePrediction saves a prediction to the cache
func (r *lruCacheRepository) SavePrediction(userID uuid.UUID, prediction model.PredictionHistory) error {
	// Skip saving if both predicted values are 0
	if prediction.Result.PredictedPrice == 0 && prediction.Result.PredictedSales == 0 {
		return nil
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Update user-specific cache
	predictions, exists := r.userCache[userID]
	if !exists {
		predictions = []model.PredictionHistory{}
	}

	// Add the new prediction at the beginning (for most recent first)
	predictions = append([]model.PredictionHistory{prediction}, predictions...)
	r.userCache[userID] = predictions

	// Also store in general LRU cache by prediction ID
	r.cache.Add(prediction.ID, prediction)

	return nil
}

// GetUserPredictions retrieves all predictions for a user from the cache
func (r *lruCacheRepository) GetUserPredictions(userID uuid.UUID) ([]model.PredictionHistory, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	predictions, exists := r.userCache[userID]
	if !exists {
		return predictions, exists
	}

	// Filter out predictions where both predicted values are 0
	filteredPredictions := []model.PredictionHistory{}
	for _, pred := range predictions {
		if !(pred.Result.PredictedPrice == 0 && pred.Result.PredictedSales == 0) {
			filteredPredictions = append(filteredPredictions, pred)
		}
	}

	return filteredPredictions, exists
}

// PopulateFromMap populates the cache with predictions from a map
func (r *lruCacheRepository) PopulateFromMap(predictions map[uuid.UUID][]model.PredictionHistory) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Clear existing cache
	r.userCache = make(map[uuid.UUID][]model.PredictionHistory)

	// Populate cache from map
	for userID, userPredictions := range predictions {
		r.userCache[userID] = userPredictions

		// Also add each prediction to the LRU cache
		for _, prediction := range userPredictions {
			r.cache.Add(prediction.ID, prediction)
		}
	}
}

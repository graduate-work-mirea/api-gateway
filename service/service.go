package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/graduate-work-mirea/api-gateway/config"
	"github.com/graduate-work-mirea/api-gateway/model"
	"github.com/graduate-work-mirea/api-gateway/repository"
)

// Service represents the business logic of the API Gateway
type Service interface {
	// Auth Service
	RegisterUser(request *model.UserRegisterRequest) (*model.UserRegisterResponse, error)
	LoginUser(request *model.UserLoginRequest) (*model.UserLoginResponse, error)

	// ML Service
	Predict(userID uuid.UUID, request *model.PredictionRequest) (*model.PredictionResult, error)
	PredictMinimal(userID uuid.UUID, request *model.PredictionRequestMinimal) (*model.PredictionResult, error)
	TrainModels() (*model.TrainingResult, error)
	GetModelStatus() (*model.ModelStatus, error)

	// Statistics
	GetUserStatistics(userID uuid.UUID) (*model.UserStatistics, error)
}

type service struct {
	config     *config.Config
	dbRepo     repository.DBRepository
	cacheRepo  repository.CacheRepository
	httpClient *http.Client
}

// NewService creates a new service
func NewService(cfg *config.Config, dbRepo repository.DBRepository, cacheRepo repository.CacheRepository) Service {
	log.Println("Service: Creating new service")
	return &service{
		config:     cfg,
		dbRepo:     dbRepo,
		cacheRepo:  cacheRepo,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// RegisterUser registers a new user
func (s *service) RegisterUser(request *model.UserRegisterRequest) (*model.UserRegisterResponse, error) {
	url := fmt.Sprintf("http://%s:%s/auth/register", s.config.Auth.Host, s.config.Auth.Port)
	log.Printf("Service: Registering user with email: %s at %s", request.Email, url)

	// Marshal request to JSON
	reqBody, err := json.Marshal(request)
	if err != nil {
		log.Printf("Service: Error marshaling register request: %v", err)
		return nil, err
	}

	// Send request to Auth service
	startTime := time.Now()
	log.Printf("Service: Sending request to Auth service: %s", url)
	resp, err := s.httpClient.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("Service: Error sending request to Auth service: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	log.Printf("Service: Auth service responded in %v with status code: %d", time.Since(startTime), resp.StatusCode)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Service: Error reading response body: %v", err)
		return nil, err
	}

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		var errResp model.ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			log.Printf("Service: Error unmarshaling error response: %v, status code: %d", err, resp.StatusCode)
			return nil, fmt.Errorf("auth service error: %d", resp.StatusCode)
		}
		log.Printf("Service: Auth service returned error: %s", errResp.Error)
		return nil, errors.New(errResp.Error)
	}

	// Unmarshal response
	var response model.UserRegisterResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Service: Error unmarshaling response: %v", err)
		return nil, err
	}

	log.Printf("Service: User registered successfully with ID: %s", response.UserID)
	return &response, nil
}

// LoginUser logs in a user
func (s *service) LoginUser(request *model.UserLoginRequest) (*model.UserLoginResponse, error) {
	url := fmt.Sprintf("http://%s:%s/auth/login", s.config.Auth.Host, s.config.Auth.Port)
	log.Printf("Service: Logging in user with email: %s at %s", request.Email, url)

	// Marshal request to JSON
	reqBody, err := json.Marshal(request)
	if err != nil {
		log.Printf("Service: Error marshaling login request: %v", err)
		return nil, err
	}

	// Send request to Auth service
	startTime := time.Now()
	log.Printf("Service: Sending request to Auth service: %s", url)
	resp, err := s.httpClient.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("Service: Error sending request to Auth service: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	log.Printf("Service: Auth service responded in %v with status code: %d", time.Since(startTime), resp.StatusCode)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Service: Error reading response body: %v", err)
		return nil, err
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		var errResp model.ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			log.Printf("Service: Error unmarshaling error response: %v, status code: %d", err, resp.StatusCode)
			return nil, fmt.Errorf("auth service error: %d", resp.StatusCode)
		}
		log.Printf("Service: Auth service returned error: %s", errResp.Error)
		return nil, errors.New(errResp.Error)
	}

	// Unmarshal response
	var response model.UserLoginResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Service: Error unmarshaling response: %v", err)
		return nil, err
	}

	log.Printf("Service: User logged in successfully with ID: %s", response.UserID)
	return &response, nil
}

// Predict makes a prediction using the ML service
func (s *service) Predict(userID uuid.UUID, request *model.PredictionRequest) (*model.PredictionResult, error) {
	url := fmt.Sprintf("http://%s:%s/api/v1/predict", s.config.ML.Host, s.config.ML.Port)
	log.Printf("Service: Making prediction for product: %s by user: %s at %s", request.ProductName, userID, url)

	// Marshal request to JSON
	reqBody, err := json.Marshal(request)
	if err != nil {
		log.Printf("Service: Error marshaling predict request: %v", err)
		return nil, err
	}

	// Send request to ML service
	startTime := time.Now()
	log.Printf("Service: Sending request to ML service: %s", url)
	resp, err := s.httpClient.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("Service: Error sending request to ML service: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	log.Printf("Service: ML service responded in %v with status code: %d", time.Since(startTime), resp.StatusCode)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Service: Error reading response body: %v", err)
		return nil, err
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		var errResp model.ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			log.Printf("Service: Error unmarshaling error response: %v, status code: %d", err, resp.StatusCode)
			return nil, fmt.Errorf("ml service error: %d", resp.StatusCode)
		}
		log.Printf("Service: ML service returned error: %s", errResp.Error)
		return nil, errors.New(errResp.Error)
	}

	// Unmarshal response
	var result model.PredictionResult
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Service: Error unmarshaling response: %v", err)
		return nil, err
	}

	// Create prediction history
	prediction := model.PredictionHistory{
		ID:           uuid.New(),
		UserID:       userID,
		Request:      *request,
		Result:       result,
		CreatedAt:    time.Now(),
		EndpointType: "predict",
		Minimal:      false,
	}

	// Only save predictions where both predicted values are not zero
	if !(result.PredictedPrice == 0 && result.PredictedSales == 0) {
		// Save prediction to database
		go func() {
			log.Printf("Service: Saving prediction to database for user: %s", userID)
			if err := s.dbRepo.SavePrediction(userID, request, &result, false); err != nil {
				log.Printf("Service: Error saving prediction to database: %v", err)
			} else {
				log.Printf("Service: Prediction saved to database successfully")
			}
		}()

		// Save prediction to cache
		go func() {
			log.Printf("Service: Saving prediction to cache for user: %s", userID)
			if err := s.cacheRepo.SavePrediction(userID, prediction); err != nil {
				log.Printf("Service: Error saving prediction to cache: %v", err)
			} else {
				log.Printf("Service: Prediction saved to cache successfully")
			}
		}()
	} else {
		log.Printf("Service: Skipping saving prediction with zero values for user: %s", userID)
	}

	log.Printf("Service: Prediction successful, price: %f, sales: %f", result.PredictedPrice, result.PredictedSales)
	return &result, nil
}

// PredictMinimal makes a prediction using the ML service with minimal input
func (s *service) PredictMinimal(userID uuid.UUID, request *model.PredictionRequestMinimal) (*model.PredictionResult, error) {
	url := fmt.Sprintf("http://%s:%s/api/v1/predict/minimal", s.config.ML.Host, s.config.ML.Port)

	// Marshal request to JSON
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	// Send request to ML service
	resp, err := s.httpClient.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		var errResp model.ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, fmt.Errorf("ml service error: %d", resp.StatusCode)
		}
		return nil, errors.New(errResp.Error)
	}

	// Unmarshal response
	var result model.PredictionResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// Create prediction history
	prediction := model.PredictionHistory{
		ID:           uuid.New(),
		UserID:       userID,
		Result:       result,
		CreatedAt:    time.Now(),
		EndpointType: "predict/minimal",
		Minimal:      true,
	}

	// Only save predictions where both predicted values are not zero
	if !(result.PredictedPrice == 0 && result.PredictedSales == 0) {
		// Save prediction to database
		go func() {
			if err := s.dbRepo.SavePrediction(userID, request, &result, true); err != nil {
				log.Printf("Error saving prediction to database: %v", err)
			}
		}()

		// Save prediction to cache
		go func() {
			if err := s.cacheRepo.SavePrediction(userID, prediction); err != nil {
				log.Printf("Error saving prediction to cache: %v", err)
			}
		}()
	} else {
		log.Printf("Service: Skipping saving prediction with zero values for user: %s", userID)
	}

	return &result, nil
}

// TrainModels trains the ML models
func (s *service) TrainModels() (*model.TrainingResult, error) {
	url := fmt.Sprintf("http://%s:%s/api/v1/train", s.config.ML.Host, s.config.ML.Port)

	// Send request to ML service
	resp, err := s.httpClient.Post(url, "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		var errResp model.ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, fmt.Errorf("ml service error: %d", resp.StatusCode)
		}
		return nil, errors.New(errResp.Error)
	}

	// Unmarshal response
	var result model.TrainingResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetModelStatus gets the status of the ML models
func (s *service) GetModelStatus() (*model.ModelStatus, error) {
	url := fmt.Sprintf("http://%s:%s/api/v1/status", s.config.ML.Host, s.config.ML.Port)

	// Send request to ML service
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		var errResp model.ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, fmt.Errorf("ml service error: %d", resp.StatusCode)
		}
		return nil, errors.New(errResp.Error)
	}

	// Unmarshal response
	var status model.ModelStatus
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

// GetUserStatistics gets statistics for a user
func (s *service) GetUserStatistics(userID uuid.UUID) (*model.UserStatistics, error) {
	log.Printf("Service: Getting statistics for user: %s", userID)

	// Try to get predictions from cache first
	log.Printf("Service: Attempting to get predictions from cache for user: %s", userID)
	predictions, found := s.cacheRepo.GetUserPredictions(userID)
	if found {
		log.Printf("Service: Found %d predictions in cache for user: %s", len(predictions), userID)
		return &model.UserStatistics{
			UserID:      userID,
			Predictions: predictions,
		}, nil
	}

	// If not in cache, get from database
	log.Printf("Service: No cache entry found, getting predictions from database for user: %s", userID)
	predictions, err := s.dbRepo.GetUserPredictions(userID)
	if err != nil {
		log.Printf("Service: Error getting predictions from database: %v", err)
		return nil, err
	}

	log.Printf("Service: Found %d predictions in database for user: %s", len(predictions), userID)
	return &model.UserStatistics{
		UserID:      userID,
		Predictions: predictions,
	}, nil
}

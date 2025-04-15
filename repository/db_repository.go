package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/graduate-work-mirea/api-gateway/config"
	"github.com/graduate-work-mirea/api-gateway/model"
	_ "github.com/lib/pq"
)

// DBRepository represents a PostgreSQL repository
type DBRepository interface {
	SavePrediction(userID uuid.UUID, request interface{}, result *model.PredictionResult, minimal bool) error
	GetUserPredictions(userID uuid.UUID) ([]model.PredictionHistory, error)
	Close() error
}

type postgreRepository struct {
	db *sql.DB
}

// NewPostgreRepository creates a new PostgreSQL repository
func NewPostgreRepository(cfg *config.Config) (DBRepository, error) {
	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Check connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		return nil, err
	}

	return &postgreRepository{db: db}, nil
}

// createTables creates the necessary tables if they don't exist
func createTables(db *sql.DB) error {
	// Create predictions table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS prediction_history (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL,
			request JSONB NOT NULL,
			result JSONB NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			endpoint_type VARCHAR(50) NOT NULL,
			minimal BOOLEAN NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`)
	if err != nil {
		return err
	}

	return nil
}

// SavePrediction saves a prediction request and result to the database
func (r *postgreRepository) SavePrediction(userID uuid.UUID, request interface{}, result *model.PredictionResult, minimal bool) error {
	// Skip saving if both predicted values are 0
	if result.PredictedPrice == 0 && result.PredictedSales == 0 {
		return nil
	}

	// Convert request and result to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return err
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return err
	}

	endpointType := "predict"
	if minimal {
		endpointType = "predict/minimal"
	}

	// Insert prediction history
	_, err = r.db.Exec(`
		INSERT INTO prediction_history (user_id, request, result, endpoint_type, minimal, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, userID, requestJSON, resultJSON, endpointType, minimal, time.Now())
	if err != nil {
		log.Printf("Error saving prediction: %v", err)
		return err
	}

	return nil
}

// GetUserPredictions retrieves all predictions for a user
func (r *postgreRepository) GetUserPredictions(userID uuid.UUID) ([]model.PredictionHistory, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, request, result, created_at, endpoint_type, minimal
		FROM prediction_history
		WHERE user_id = $1 
		AND (result->>'predicted_price' != '0' OR result->>'predicted_sales' != '0')
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var predictions []model.PredictionHistory
	for rows.Next() {
		var prediction model.PredictionHistory
		var requestJSON, resultJSON []byte

		err := rows.Scan(
			&prediction.ID,
			&prediction.UserID,
			&requestJSON,
			&resultJSON,
			&prediction.CreatedAt,
			&prediction.EndpointType,
			&prediction.Minimal,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal request based on minimal flag
		if prediction.Minimal {
			var minimalRequest model.PredictionRequestMinimal
			if err := json.Unmarshal(requestJSON, &minimalRequest); err != nil {
				return nil, err
			}
		} else {
			if err := json.Unmarshal(requestJSON, &prediction.Request); err != nil {
				return nil, err
			}
		}

		// Unmarshal result
		if err := json.Unmarshal(resultJSON, &prediction.Result); err != nil {
			return nil, err
		}

		predictions = append(predictions, prediction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return predictions, nil
}

// Close closes the database connection
func (r *postgreRepository) Close() error {
	return r.db.Close()
}

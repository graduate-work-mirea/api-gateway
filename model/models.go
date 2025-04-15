package model

import (
	"time"

	"github.com/google/uuid"
)

// PredictionHistory represents a saved prediction request and result
type PredictionHistory struct {
	ID           uuid.UUID         `json:"id" db:"id"`
	UserID       uuid.UUID         `json:"user_id" db:"user_id"`
	Request      PredictionRequest `json:"request" db:"request"`
	Result       PredictionResult  `json:"result" db:"result"`
	CreatedAt    time.Time         `json:"created_at" db:"created_at"`
	EndpointType string            `json:"endpoint_type" db:"endpoint_type"`
	Minimal      bool              `json:"minimal" db:"minimal"`
}

// UserStatistics represents statistics for a user's prediction requests
type UserStatistics struct {
	UserID      uuid.UUID           `json:"user_id"`
	Predictions []PredictionHistory `json:"predictions"`
}

// Auth Service Models

// UserRegisterRequest represents a request to register a new user
type UserRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserRegisterResponse represents a response to a user registration request
type UserRegisterResponse struct {
	UserID       string    `json:"user_id"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    int64     `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// UserLoginRequest represents a request to login a user
type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserLoginResponse represents a response to a user login request
type UserLoginResponse struct {
	UserID       string    `json:"user_id"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    int64     `json:"expires_at"`
	LastLoginAt  time.Time `json:"last_login_at"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// ML Service Models

// PredictionRequest represents a request to make a prediction
type PredictionRequest struct {
	ProductName               string  `json:"product_name"`
	Brand                     string  `json:"brand"`
	Category                  string  `json:"category"`
	Region                    string  `json:"region"`
	Seller                    string  `json:"seller"`
	Price                     float64 `json:"price"`
	OriginalPrice             float64 `json:"original_price"`
	DiscountPercentage        float64 `json:"discount_percentage"`
	StockLevel                float64 `json:"stock_level"`
	CustomerRating            float64 `json:"customer_rating"`
	ReviewCount               float64 `json:"review_count"`
	DeliveryDays              float64 `json:"delivery_days"`
	IsWeekend                 bool    `json:"is_weekend"`
	IsHoliday                 bool    `json:"is_holiday"`
	DayOfWeek                 int     `json:"day_of_week"`
	Month                     int     `json:"month"`
	Quarter                   int     `json:"quarter"`
	SalesQuantityLag1         float64 `json:"sales_quantity_lag_1"`
	PriceLag1                 float64 `json:"price_lag_1"`
	SalesQuantityLag3         float64 `json:"sales_quantity_lag_3"`
	PriceLag3                 float64 `json:"price_lag_3"`
	SalesQuantityLag7         float64 `json:"sales_quantity_lag_7"`
	PriceLag7                 float64 `json:"price_lag_7"`
	SalesQuantityRollingMean3 float64 `json:"sales_quantity_rolling_mean_3"`
	PriceRollingMean3         float64 `json:"price_rolling_mean_3"`
	SalesQuantityRollingMean7 float64 `json:"sales_quantity_rolling_mean_7"`
	PriceRollingMean7         float64 `json:"price_rolling_mean_7"`
}

// PredictionRequestMinimal represents a minimal request to make a prediction
type PredictionRequestMinimal struct {
	ProductName    string     `json:"product_name"`
	Region         string     `json:"region"`
	Seller         string     `json:"seller"`
	PredictionDate *time.Time `json:"prediction_date,omitempty"`
	Price          *float64   `json:"price,omitempty"`
	OriginalPrice  *float64   `json:"original_price,omitempty"`
	StockLevel     *float64   `json:"stock_level,omitempty"`
	CustomerRating *float64   `json:"customer_rating,omitempty"`
	ReviewCount    *float64   `json:"review_count,omitempty"`
	DeliveryDays   *float64   `json:"delivery_days,omitempty"`
}

// PredictionResult represents a prediction result
type PredictionResult struct {
	PredictedPrice float64 `json:"predicted_price"`
	PredictedSales float64 `json:"predicted_sales"`
}

// TrainingResult represents a training result
type TrainingResult struct {
	PriceModel struct {
		BestIteration int     `json:"best_iteration"`
		BestScore     float64 `json:"best_score"`
	} `json:"price_model"`
	SalesModel struct {
		BestIteration int     `json:"best_iteration"`
		BestScore     float64 `json:"best_score"`
	} `json:"sales_model"`
}

// ModelStatus represents the status of the prediction models
type ModelStatus struct {
	ModelsTrained bool `json:"models_trained"`
}

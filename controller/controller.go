package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graduate-work-mirea/api-gateway/middleware"
	"github.com/graduate-work-mirea/api-gateway/model"
	"github.com/graduate-work-mirea/api-gateway/service"
)

// Controller represents the HTTP request controller
type Controller struct {
	service service.Service
	router  *gin.Engine
}

// NewController creates a new controller
func NewController(service service.Service, router *gin.Engine) *Controller {
	log.Println("Controller: Creating new controller")
	return &Controller{
		service: service,
		router:  router,
	}
}

// RegisterRoutes registers all routes to the router
func (c *Controller) RegisterRoutes(authMiddleware gin.HandlerFunc) {
	log.Println("Controller: Registering routes...")

	// Auth routes
	authGroup := c.router.Group("/auth")
	{
		authGroup.POST("/register", c.registerUser)
		authGroup.POST("/login", c.loginUser)
	}
	log.Println("Controller: Auth routes registered: POST /auth/register, POST /auth/login")

	// ML routes
	mlGroup := c.router.Group("/api/v1")
	mlGroup.Use(authMiddleware)
	{
		mlGroup.POST("/predict", c.predict)
		mlGroup.POST("/predict/minimal", c.predictMinimal)
		mlGroup.POST("/train", c.trainModels)
		mlGroup.GET("/status", c.getModelStatus)
	}
	log.Println("Controller: ML routes registered with auth middleware: POST /api/v1/predict, POST /api/v1/predict/minimal, POST /api/v1/train, GET /api/v1/status")

	// Statistics routes
	statsGroup := c.router.Group("/api/v1/statistics")
	statsGroup.Use(authMiddleware)
	{
		statsGroup.GET("/user", c.getUserStatistics)
	}
	log.Println("Controller: Statistics routes registered with auth middleware: GET /api/v1/statistics/user")
	log.Println("Controller: All routes registered")
}

// registerUser handles user registration
func (c *Controller) registerUser(ctx *gin.Context) {
	log.Println("Controller: Handling registerUser request")
	var request model.UserRegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("Controller: Invalid request format: %v", err)
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request format"})
		return
	}

	log.Printf("Controller: Registering user with email: %s", request.Email)
	response, err := c.service.RegisterUser(&request)
	if err != nil {
		log.Printf("Controller: Error registering user: %v", err)
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("Controller: User registered successfully with ID: %s", response.UserID)
	ctx.JSON(http.StatusCreated, response)
}

// loginUser handles user login
func (c *Controller) loginUser(ctx *gin.Context) {
	log.Println("Controller: Handling loginUser request")
	var request model.UserLoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("Controller: Invalid request format: %v", err)
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request format"})
		return
	}

	log.Printf("Controller: Logging in user with email: %s", request.Email)
	response, err := c.service.LoginUser(&request)
	if err != nil {
		log.Printf("Controller: Error logging in user: %v", err)
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("Controller: User logged in successfully with ID: %s", response.UserID)
	ctx.JSON(http.StatusOK, response)
}

// predict handles predictions
func (c *Controller) predict(ctx *gin.Context) {
	log.Println("Controller: Handling predict request")
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		log.Printf("Controller: Unauthorized access: %v", err)
		ctx.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: err.Error()})
		return
	}

	var request model.PredictionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("Controller: Invalid request format: %v", err)
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request format"})
		return
	}

	log.Printf("Controller: Making prediction for product: %s by user: %s", request.ProductName, userID)
	result, err := c.service.Predict(userID, &request)
	if err != nil {
		log.Printf("Controller: Error making prediction: %v", err)
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("Controller: Prediction successful, price: %f, sales: %f", result.PredictedPrice, result.PredictedSales)
	ctx.JSON(http.StatusOK, result)
}

// predictMinimal handles minimal predictions
func (c *Controller) predictMinimal(ctx *gin.Context) {
	log.Println("Controller: Handling predictMinimal request")
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		log.Printf("Controller: Unauthorized access: %v", err)
		ctx.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: err.Error()})
		return
	}

	var request model.PredictionRequestMinimal
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("Controller: Invalid request format: %v", err)
		ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request format"})
		return
	}

	log.Printf("Controller: Making minimal prediction for product: %s by user: %s", request.ProductName, userID)
	result, err := c.service.PredictMinimal(userID, &request)
	if err != nil {
		log.Printf("Controller: Error making minimal prediction: %v", err)
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("Controller: Minimal prediction successful, price: %f, sales: %f", result.PredictedPrice, result.PredictedSales)
	ctx.JSON(http.StatusOK, result)
}

// trainModels handles model training
func (c *Controller) trainModels(ctx *gin.Context) {
	log.Println("Controller: Handling trainModels request")
	result, err := c.service.TrainModels()
	if err != nil {
		log.Printf("Controller: Error training models: %v", err)
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	log.Println("Controller: Models trained successfully")
	ctx.JSON(http.StatusOK, result)
}

// getModelStatus handles getting model status
func (c *Controller) getModelStatus(ctx *gin.Context) {
	log.Println("Controller: Handling getModelStatus request")
	status, err := c.service.GetModelStatus()
	if err != nil {
		log.Printf("Controller: Error getting model status: %v", err)
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("Controller: Model status retrieved, trained: %t", status.ModelsTrained)
	ctx.JSON(http.StatusOK, status)
}

// getUserStatistics handles getting user statistics
func (c *Controller) getUserStatistics(ctx *gin.Context) {
	log.Println("Controller: Handling getUserStatistics request")
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		log.Printf("Controller: Unauthorized access: %v", err)
		ctx.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("Controller: Getting statistics for user: %s", userID)
	statistics, err := c.service.GetUserStatistics(userID)
	if err != nil {
		log.Printf("Controller: Error getting user statistics: %v", err)
		ctx.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("Controller: Statistics retrieved, prediction count: %d", len(statistics.Predictions))
	ctx.JSON(http.StatusOK, statistics)
}

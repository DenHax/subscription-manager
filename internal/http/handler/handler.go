package handler

import (
	"net/http"

	"github.com/DenHax/subscription-manager/internal/service"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	Services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		Services: services,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.New()

	router.GET("/swagger", h.redirectToSwagger)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", h.CheckHealth)

	apiV1 := router.Group("/api/v1")
	{
		subscriptions := apiV1.Group("/subscriptions")
		{
			subscriptions.GET("/", h.ListSubscriptions)
			subscriptions.POST("/", h.CreateSubscription)
			subscriptions.GET("/:id", h.GetSubscriptionByID)
			subscriptions.PUT("/:id", h.UpdateSubscription)
			subscriptions.DELETE("/:id", h.DeleteSubscription)
		}
		apiV1.GET("/subscriptions/summary", h.GetSubscriptionSummary)
	}
	return router
}

func (h *Handler) redirectToSwagger(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
}

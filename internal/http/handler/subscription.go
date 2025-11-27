package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) ListSubscriptions(c *gin.Context) {
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")

	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	var limit, offset int
	var err error

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
			return
		}
	} else {
		limit = 10
	}

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset parameter"})
			return
		}
	} else {
		offset = 0
	}

	subscriptions, total, err := h.Services.GetAllSubscriptions(&userID, &serviceName, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"subscriptions": subscriptions,
		"total":         total,
	})
}

func (h *Handler) CreateSubscription(c *gin.Context) {
	var req struct {
		ServiceName string `json:"service_name" binding:"required"`
		Price       int    `json:"price" binding:"required"`
		UserID      string `json:"user_id" binding:"required"`
		StartDate   string `json:"start_date" binding:"required,datetime=01-2006"`
		EndDate     string `json:"end_date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := uuid.Parse(req.UserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
		return
	}

	subscription, err := h.Services.CreateSubscrition(req.ServiceName, req.Price, req.UserID, req.StartDate, &req.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

func (h *Handler) GetSubscriptionByID(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id format"})
		return
	}

	subscription, err := h.Services.Subscription(id)
	if err != nil {
		if err.Error() == "subscription not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

func (h *Handler) UpdateSubscription(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id format"})
		return
	}

	var req struct {
		ServiceName string `json:"service_name"`
		Price       int    `json:"price"`
		UserID      string `json:"user_id"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.UserID != "" {
		if _, err := uuid.Parse(req.UserID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
			return
		}
	}

	subscription, err := h.Services.UpdateSubscription(&id, &req.ServiceName, &req.Price, &req.UserID, &req.StartDate, &req.EndDate)
	if err != nil {
		if err.Error() == "subscription not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

func (h *Handler) DeleteSubscription(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id format"})
		return
	}

	err := h.Services.DeleteSubscription(id)
	if err != nil {
		if err.Error() == "subscription not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) GetSubscriptionSummary(c *gin.Context) {
	startDate := c.Query("start_date")
	if startDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date is required"})
		return
	}

	endDate := c.Query("end_date")
	if endDate == "" {
		endDate = startDate
	}

	userID := c.Query("user_id")
	serviceName := c.Query("service_name")

	if userID != "" {
		if _, err := uuid.Parse(userID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
			return
		}
	}

	totalCost, err := h.Services.SummarySubscription(startDate, endDate, &userID, &serviceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"total_cost": totalCost,
		"period": gin.H{
			"start_date": startDate,
			"end_date":   endDate,
		},
	}

	if userID != "" || serviceName != "" {
		filters := gin.H{}
		if userID != "" {
			filters["user_id"] = userID
		}
		if serviceName != "" {
			filters["service_name"] = serviceName
		}
		response["filters"] = filters
	}

	c.JSON(http.StatusOK, response)
}

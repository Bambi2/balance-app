package delivery

import (
	_ "github.com/bambi2/balance-app/docs"
	"github.com/bambi2/balance-app/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		balance := api.Group("/balance")
		{
			balance.POST("/", h.AddToBalance)
			balance.GET("/:userId", h.GetBalance)
			balance.POST("/send", h.SendMoney)
		}

		api.GET("/transactions/:userId/:limit/:offset", h.GetTransactions)
		api.POST("/reserve", h.Reserve)
		api.POST("/unreserve", h.Unreserve)
		api.POST("/check", h.Check)
		api.GET("/checks:year/:month", h.GetChecks)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

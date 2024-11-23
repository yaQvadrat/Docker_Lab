package httpapi

import (
	"app/internal/service"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

func ConfigureRouter(handler *echo.Echo, services *service.Services) {
	handler.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Output: setLogsFile()}))
	handler.Use(middleware.Recover())

	api := handler.Group("/api")
	{
		api.GET("/ping", func(c echo.Context) error { return c.String(http.StatusOK, "ok") })

		tenders := api.Group("/tenders")
		{
			r := newTenderRoutes(services.Tender)
			tenders.POST("/new", r.newTender)
			tenders.GET("/my", r.myTenders)
			tenders.GET("", r.tenders)
			tenders.PUT("/:tenderId/status", r.putStatus)
			tenders.GET("/:tenderId/status", r.getStatus)
			tenders.PATCH("/:tenderId/edit", r.editTender)
			tenders.PUT("/:tenderId/rollback/:version", r.rollbackTender)
		}

		bids := api.Group("/bids")
		{
			r := newBidRoutes(services.Bid)
			bids.POST("/new", r.newBid)
			bids.PUT("/:bidId/submit_decision", r.submitDecision)
			bids.PUT("/:bidId/status", r.putStatus)
			bids.GET("/:bidId/status", r.getStatus)
			bids.PATCH("/:bidId/edit", r.editBid)
			bids.PUT("/:bidId/rollback/:version", r.rollbackBid)
		}
	}
}

func setLogsFile() *os.File {
	file, err := os.OpenFile("./logs/logfile.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(fmt.Errorf("httpapi - setLogsFile - os.OpenFile: %w", err))
	}
	return file
}

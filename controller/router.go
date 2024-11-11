package controller

import (
	"github.com/rahul0tripathi/pipetg/controller/handlers"
	"github.com/rahul0tripathi/pipetg/internal/integrations/tg"
	"github.com/rahul0tripathi/pipetg/internal/services"
	"github.com/rahul0tripathi/pipetg/pkg/httpserver"
)

func Router(
	router httpserver.Router,
	wrapper *tg.Client,
	authFlowService *services.AuthFlowService,
	messagesService *services.Scraper,
) {

	handler := handlers.New(wrapper)
	base := router.Group("/pipetg")
	auth := base.Group("/auth")
	{
		auth.POST("/send", handler.MakeHandleSendCode(authFlowService))
		auth.POST("/submit", handler.MakeHandleSubmitCode(authFlowService))
	}

	messages := base.Group("/scheduler")
	{
		messages.GET("/scrape", handler.MakeFetchAllMessages(messagesService))
	}
}

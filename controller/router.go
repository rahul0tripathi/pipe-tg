package controller

import (
	"github.com/rahul0tripathi/pipetg/controller/handlers"
	"github.com/rahul0tripathi/pipetg/internal/services"
	"github.com/rahul0tripathi/pipetg/pkg/httpserver"
)

func Router(
	router httpserver.Router,
	authFlowService *services.AuthFlowService,
	messagesService *services.MessageLogger,
) {

	handler := handlers.New()
	base := router.Group("/pipetg")
	auth := base.Group("/auth")
	{
		auth.POST("/send", handler.MakeHandleSendCode(authFlowService))
		auth.POST("/submit", handler.MakeHandleSubmitCode(authFlowService))
	}

	messages := base.Group("/messages")
	{
		messages.GET("/", handler.MakeFetchAllMessages(messagesService))
	}
}

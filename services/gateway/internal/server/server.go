package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/gateway/internal/config"
)

type Server struct {
	config     *config.Config
	httpServer *http.Server
}

func NewHTTPServer(config *config.Config) (*Server, error) {
	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.HTTP.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(config.HTTP.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.HTTP.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.HTTP.IdleTimeout) * time.Second,
	}

	server := &Server{
		config:     config,
		httpServer: httpServer,
	}

	server.SetupRoutes(router)

	return server, nil
}

func (s *Server) SetupRoutes(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong from gateway",
		})
	})
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

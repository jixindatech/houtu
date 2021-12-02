package server

import (
	"admin/config"
	"admin/core/rbac"
	"admin/server/models"
	"admin/server/router"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	cfg    *config.Config
	server *http.Server
}

func (s *Server) Setup(cfg *config.Config) error {
	s.cfg = cfg

	err := models.Setup(cfg.Database)
	if err != nil {
		return err
	}

	err = rbac.Setup(cfg.Rbac)
	if err != nil {
		return err
	}

	r, err := router.Setup(cfg.RunMode)
	if err != nil {
		return err
	}

	s.server = &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:        r,
		ReadTimeout:    cfg.Server.ReadTimeout * time.Second,
		WriteTimeout:   cfg.Server.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return nil
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Close() error {
	return s.server.Close()
}

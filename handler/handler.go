package handler

import "forum/service"

type Handler struct {
	usecases *service.Service
}

func NewHandler(usecases *service.Service) *Handler {
	return &Handler{usecases: usecases}
}

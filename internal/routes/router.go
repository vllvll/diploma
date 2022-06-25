package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vllvll/diploma/internal/handlers"
	"github.com/vllvll/diploma/internal/middlewares"
	"github.com/vllvll/diploma/internal/repositories"
)

type Router struct {
	Router   chi.Router
	handlers handlers.Handler
}

func NewRouter(handlers handlers.Handler, userRepository repositories.UserInterface, tokenRepository repositories.TokenInterface) Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(middlewares.Auth(userRepository, tokenRepository))

	return Router{
		Router:   r,
		handlers: handlers,
	}
}

func (ro *Router) RegisterHandlers() {
	ro.Router.Route("/api/user/", func(r chi.Router) {
		r.Post("/register", ro.handlers.Register())
		r.Post("/login", ro.handlers.Login())
		r.Get("/balance", ro.handlers.GetBalance())

		r.Post("/orders", ro.handlers.AddOrder())
		r.Get("/orders", ro.handlers.GetOrders())

		ro.Router.Route("/balance/", func(r chi.Router) {
			r.Post("/withdraw", ro.handlers.AddWithdraw())
			r.Get("/withdrawals", ro.handlers.GetWithdrawals())
		})
	})
}

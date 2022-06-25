package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/vllvll/diploma/internal/middlewares"
	"github.com/vllvll/diploma/internal/repositories"
	"github.com/vllvll/diploma/internal/services"
	"github.com/vllvll/diploma/internal/types"
	"net/http"
	"time"
)

type Handler struct {
	userRepository  repositories.UserInterface
	tokenRepository repositories.TokenInterface
	cryptService    services.CryptInterface
}

type UserHandlers interface {
	Register() http.HandlerFunc
	Login() http.HandlerFunc
	AddOrder() http.HandlerFunc
	GetOrders() http.HandlerFunc
	GetBalance() http.HandlerFunc
	AddWithdraw() http.HandlerFunc
	GetWithdrawals() http.HandlerFunc
}

func NewHandler(userRepository repositories.UserInterface, tokenRepository repositories.TokenInterface, cryptService services.CryptInterface) *Handler {
	return &Handler{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		cryptService:    cryptService,
	}
}

//func (h Handler) isAuth(r *http.Request) bool {
//	c, err := r.Cookie("gophermart-auth-cookie")
//	if err != nil || c.Value == "" || !h.tokenRepository.IsExists(c.Value) {
//		return false
//	}
//
//	return true
//}

func (h Handler) Register() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var userItem types.UserRequest

		if err := json.NewDecoder(r.Body).Decode(&userItem); err != nil {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if h.userRepository.IsExists(userItem.Login) {
			http.Error(rw, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}

		password := h.cryptService.Hash(userItem.Password)

		userId, err := h.userRepository.CreateUser(userItem.Login, password)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		token, err := h.cryptService.GenerateRand()
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = h.tokenRepository.CreateToken(token, userId)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.SetCookie(rw, &http.Cookie{
			Name:    "gophermart-auth-cookie",
			Value:   token,
			Expires: time.Now().Add(365 * 24 * time.Hour),
		})

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(http.StatusText(http.StatusOK)))
	}
}

func (h Handler) Login() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var userItem types.UserRequest

		if err := json.NewDecoder(r.Body).Decode(&userItem); err != nil {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		user, err := h.userRepository.GetUserHashByLogin(userItem.Login)
		if err != nil || !h.cryptService.IsEqual(userItem.Password, user.Hash) {
			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		token, err := h.cryptService.GenerateRand()
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = h.tokenRepository.CreateToken(token, user.Id)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.SetCookie(rw, &http.Cookie{
			Name:    "gophermart-auth-cookie",
			Value:   token,
			Expires: time.Now().Add(365 * 24 * time.Hour),
		})

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(http.StatusText(http.StatusOK)))
	}
}

func (h Handler) GetOrders() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		user := middlewares.ForContext(r.Context())
		fmt.Println(user)

		//if !h.isAuth(r) {
		//	http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		//	return
		//}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(http.StatusText(http.StatusOK)))
	}
}

func (h Handler) AddOrder() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		//_ := middlewares.ForContext(r.Context())

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(http.StatusText(http.StatusOK)))
	}
}

func (h Handler) GetBalance() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(http.StatusText(http.StatusOK)))
	}
}

func (h Handler) AddWithdraw() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(http.StatusText(http.StatusOK)))
	}
}

func (h Handler) GetWithdrawals() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(http.StatusText(http.StatusOK)))
	}
}

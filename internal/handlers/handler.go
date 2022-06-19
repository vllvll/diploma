package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/vllvll/diploma/internal/middlewares"
	"github.com/vllvll/diploma/internal/repositories"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	userRepository repositories.UserInterface
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

func NewHandler(userRepository repositories.UserInterface) *Handler {
	return &Handler{
		userRepository: userRepository,
	}
}

//func (h Handler) Ping() http.HandlerFunc {
//	return func(rw http.ResponseWriter, r *http.Request) {
//		err := h.db.Ping()
//		if err != nil {
//			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//
//			return
//		}
//
//		rw.WriteHeader(http.StatusOK)
//		rw.Write([]byte(http.StatusText(http.StatusOK)))
//	}
//}

func (h Handler) Register() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		type User struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}

		var userItem User

		if err := json.NewDecoder(r.Body).Decode(&userItem); err != nil {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		isExists, err := h.userRepository.IsExists(userItem.Login)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if isExists {
			http.Error(rw, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}

		id, err := h.userRepository.CreateUser(userItem.Login, userItem.Password)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		//r.AddCookie(&http.Cookie{
		//	Name:    "gophermart-auth-cookie",
		//	Value:   "token",
		//	Expires: time.Now().Add(365 * 24 * time.Hour),
		//})

		http.SetCookie(rw, &http.Cookie{
			Name:    "gophermart-auth-cookie",
			Value:   "token",
			Expires: time.Now().Add(365 * 24 * time.Hour),
		})

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(strconv.Itoa(id)))
	}
}

func (h Handler) Login() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		cookie := r.Cookies()

		fmt.Println(cookie)

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(http.StatusText(http.StatusOK)))
	}
}

func (h Handler) GetOrders() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if user := middlewares.ForContext(r.Context()); user == nil {
			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(http.StatusText(http.StatusOK)))
	}
}

func (h Handler) AddOrder() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
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

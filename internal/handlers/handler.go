package handlers

import (
	"encoding/json"
	"github.com/vllvll/diploma/internal/middlewares"
	"github.com/vllvll/diploma/internal/repositories"
	"github.com/vllvll/diploma/internal/services"
	"github.com/vllvll/diploma/internal/types"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	userRepository    repositories.UserInterface
	tokenRepository   repositories.TokenInterface
	orderRepository   repositories.OrderInterface
	balanceRepository repositories.BalanceInterface
	cryptService      services.CryptInterface
	luhnService       services.LuhnInterface
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

func NewHandler(
	userRepository repositories.UserInterface,
	tokenRepository repositories.TokenInterface,
	orderRepository repositories.OrderInterface,
	balanceRepository repositories.BalanceInterface,
	cryptService services.CryptInterface,
	luhnService services.LuhnInterface,
) *Handler {
	return &Handler{
		userRepository:    userRepository,
		tokenRepository:   tokenRepository,
		orderRepository:   orderRepository,
		balanceRepository: balanceRepository,
		cryptService:      cryptService,
		luhnService:       luhnService,
	}
}

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

		err = h.balanceRepository.CreateBalance(userId)
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
		if user == nil {
			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		orders, err := h.orderRepository.GetOrdersByUser(user.Id)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if len(orders) == 0 {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusNoContent)
			rw.Write([]byte(http.StatusText(http.StatusNoContent)))
		}

		response, err := json.Marshal(orders)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(response)
	}
}

func (h Handler) AddOrder() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		user := middlewares.ForContext(r.Context())
		if user == nil {
			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		orderNumber, err := strconv.ParseInt(string(body), 10, 64)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if !h.luhnService.IsValid(orderNumber) {
			http.Error(rw, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}

		order, err := h.orderRepository.GetByNumber(orderNumber)
		if err != nil {
			err := h.orderRepository.CreateOrder(orderNumber, user.Id)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			rw.WriteHeader(http.StatusAccepted)
			rw.Write([]byte(http.StatusText(http.StatusAccepted)))
			return
		}

		if order.UserId != user.Id {
			http.Error(rw, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(http.StatusText(http.StatusOK)))
	}
}

func (h Handler) GetBalance() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		user := middlewares.ForContext(r.Context())
		if user == nil {
			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		balance, err := h.balanceRepository.GetByUserId(user.Id)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(balance)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(response)
	}
}

func (h Handler) AddWithdraw() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		user := middlewares.ForContext(r.Context())
		if user == nil {
			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		var withdraw types.WithdrawRequest
		if err := json.NewDecoder(r.Body).Decode(&withdraw); err != nil {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		orderNumber, err := strconv.ParseInt(withdraw.Order, 10, 64)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if !h.luhnService.IsValid(orderNumber) {
			http.Error(rw, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}

		isUpdate, err := h.balanceRepository.UpdateBalance(user.Id, orderNumber, withdraw.Sum)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if !isUpdate {
			http.Error(rw, http.StatusText(http.StatusPaymentRequired), http.StatusPaymentRequired)
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(http.StatusText(http.StatusOK)))
	}
}

func (h Handler) GetWithdrawals() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		user := middlewares.ForContext(r.Context())
		if user == nil {
			http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		withdrawals, err := h.balanceRepository.GetWithdrawals(user.Id)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if len(withdrawals) == 0 {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusNoContent)
			rw.Write([]byte(http.StatusText(http.StatusNoContent)))
		}

		response, err := json.Marshal(withdrawals)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(response)
	}
}

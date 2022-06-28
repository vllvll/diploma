package integrations

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/vllvll/diploma/internal/config"
	"github.com/vllvll/diploma/internal/dictionaries"
	"github.com/vllvll/diploma/internal/repositories"
	"github.com/vllvll/diploma/internal/types"
	"log"
	"net/http"
	"time"
)

type LoyaltyClient struct {
	Client            *resty.Client
	orderCh           <-chan string
	orderRepository   repositories.OrderInterface
	balanceRepository repositories.BalanceInterface
}

type LoyaltyClientInterface interface {
	Processing()
}

func NewLoyaltyClient(
	config *config.GophermartConfig,
	ch <-chan string,
	orderRepository repositories.OrderInterface,
	balanceRepository repositories.BalanceInterface,
) LoyaltyClientInterface {
	client := resty.New().
		SetBaseURL(config.AccrualSystemAddress).
		SetRetryCount(5).
		SetRetryWaitTime(time.Second).
		SetRetryMaxWaitTime(5 * time.Second).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				if r.StatusCode() == http.StatusInternalServerError ||
					r.StatusCode() == http.StatusTooManyRequests ||
					r.StatusCode() == http.StatusNoContent {

					return true
				}

				var orderLoyalty types.OrderLoyalty
				if err := json.Unmarshal(r.Body(), &orderLoyalty); err != nil {
					return true
				}

				if orderLoyalty.Status == dictionaries.OrderNew || orderLoyalty.Status == dictionaries.OrderProcessing {
					return true
				}

				return false
			},
		)

	return LoyaltyClient{
		Client:            client,
		orderCh:           ch,
		orderRepository:   orderRepository,
		balanceRepository: balanceRepository,
	}
}

func (l LoyaltyClient) getOrder(number string) (orderLoyalty types.OrderLoyalty, err error) {
	response, err := l.Client.R().
		Get("/api/orders/" + number)
	if err != nil {
		return orderLoyalty, err
	}

	if err := json.Unmarshal(response.Body(), &orderLoyalty); err != nil {
		log.Printf("Error get order processing: %v", err)

		return orderLoyalty, err
	}

	return orderLoyalty, nil
}

func (l LoyaltyClient) Processing() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic: %v", err)

			l.Processing()
		}
	}()

	for orderNumber := range l.orderCh {
		orderLoyalty, err := l.getOrder(orderNumber)
		if err != nil {
			log.Printf("Error get order: %v", err)

			continue
		}

		userID, err := l.orderRepository.UpdateOrder(orderNumber, orderLoyalty.Status, orderLoyalty.Accrual)
		if err != nil {
			log.Printf("Error update order: %v", err)

			continue
		}

		err = l.balanceRepository.UpdateBalance(userID, orderLoyalty.Accrual)
		if err != nil {
			log.Printf("Error update balance: %v", err)

			continue
		}
	}
}

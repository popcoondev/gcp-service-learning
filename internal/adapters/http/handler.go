package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/popcoondev/gcp-service-learning/internal/application"
	"github.com/popcoondev/gcp-service-learning/internal/domain"
)

type Handler struct {
	service application.OrderService
}

type createOrderRequest struct {
	CustomerID string                  `json:"customerId"`
	Items      []createOrderRequestItem `json:"items"`
}

type createOrderRequestItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type orderResponse struct {
	OrderID    string              `json:"orderId"`
	CustomerID string              `json:"customerId"`
	Items      []orderResponseItem `json:"items"`
	Status     string              `json:"status"`
	CreatedAt  string              `json:"createdAt"`
}

type orderResponseItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func NewHandler(service application.OrderService) Handler {
	return Handler{service: service}
}

func (h Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/orders", h.handleOrders)
	mux.HandleFunc("/orders/", h.handleOrderByID)
	return mux
}

func (h Handler) handleOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Message: "method not allowed"})
		return
	}

	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Message: "invalid json body"})
		return
	}

	items := make([]application.OrderItemInput, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, application.OrderItemInput{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	order, err := h.service.CreateOrder(r.Context(), application.CreateOrderCommand{
		CustomerID: req.CustomerID,
		Items:      items,
	})
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toOrderResponse(order))
}

func (h Handler) handleOrderByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Message: "method not allowed"})
		return
	}

	orderID := strings.TrimPrefix(r.URL.Path, "/orders/")
	if orderID == "" || strings.Contains(orderID, "/") {
		writeJSON(w, http.StatusNotFound, errorResponse{Message: "order not found"})
		return
	}

	order, err := h.service.GetOrder(r.Context(), orderID)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toOrderResponse(order))
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrOrderNotFound):
		writeJSON(w, http.StatusNotFound, errorResponse{Message: "order not found"})
	case errors.Is(err, domain.ErrInvalidOrder),
		errors.Is(err, domain.ErrInvalidOrderID),
		errors.Is(err, domain.ErrInvalidCustomer),
		errors.Is(err, domain.ErrInvalidOrderItem):
		writeJSON(w, http.StatusBadRequest, errorResponse{Message: err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, errorResponse{Message: "internal server error"})
	}
}

func toOrderResponse(order application.OrderDTO) orderResponse {
	items := make([]orderResponseItem, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, orderResponseItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	return orderResponse{
		OrderID:    order.OrderID,
		CustomerID: order.CustomerID,
		Items:      items,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(value)
}

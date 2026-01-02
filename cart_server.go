package main

import (
    "context"
    "encoding/json"
    "fmt"

    "go-micro.dev/v5"
    "go-micro.dev/v5/server"
    "go-micro.dev/v5/store"
    "go-micro.dev/v5/store/postgres"
    "go-micro.dev/v5/transport/grpc"
)

// ----- Request / Response -----
type AddItemRequest struct {
    UserID    string
    ProductID string
    Quantity  int32
}

type RemoveItemRequest struct {
    UserID    string
    ProductID string
}

type GetCartRequest struct {
    UserID string
}

type CartResponse struct {
    Items map[string]int32
}

// ----- Handler -----
type CartHandler struct {
    store store.Store
}

func NewCartHandler(s store.Store) *CartHandler {
    return &CartHandler{store: s}
}

// Load cart from Postgres store with JSON decoding
func (h *CartHandler) loadCart(userID string) map[string]int32 {
    records, err := h.store.Read(userID)
    if err != nil || len(records) == 0 {
        return make(map[string]int32)
    }
    items := make(map[string]int32)
    json.Unmarshal(records[0].Value, &items)
    return items
}

// Save cart to Postgres store with JSON encoding
func (h *CartHandler) saveCart(userID string, items map[string]int32) error {
    data, _ := json.Marshal(items)
    h.store.Delete(userID) // remove old cart
    return h.store.Write(&store.Record{
        Key:   userID,
        Value: data,
    })
}


// ----- Handler Methods -----
func (h *CartHandler) AddItem(ctx context.Context, req *AddItemRequest, rsp *CartResponse) error {
    items := h.loadCart(req.UserID)
    items[req.ProductID] += req.Quantity
    h.saveCart(req.UserID, items)
    rsp.Items = items
    return nil
}

func (h *CartHandler) RemoveItem(ctx context.Context, req *RemoveItemRequest, rsp *CartResponse) error {
    items := h.loadCart(req.UserID)
    delete(items, req.ProductID)
    h.saveCart(req.UserID, items)
    rsp.Items = items
    return nil
}

func (h *CartHandler) GetCart(ctx context.Context, req *GetCartRequest, rsp *CartResponse) error {
    items := h.loadCart(req.UserID)
    rsp.Items = items
    return nil
}

// ----- Main -----
func main() {
    // gRPC transport
    t := grpc.NewTransport()
    // Postgres store
    s := postgres.NewStore()

    service := micro.NewService(
        micro.Name("cart.service"),
        micro.Transport(t),
    )
    service.Init()

    cartHandler := NewCartHandler(s)
    srv := service.Server()
    srv.Handle(server.NewHandler(cartHandler))

    fmt.Println("Cart Service running with Postgres store and JSON persistence...")
    if err := service.Run(); err != nil {
        fmt.Println("Error running service:", err)
    }
}
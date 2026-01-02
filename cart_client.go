package main

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "go-micro.dev/v5"
    "go-micro.dev/v5/transport/grpc"
)

// Requests & Responses
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

func main() {
    // 1. gRPC transport
    t := grpc.NewTransport()

    // 2. Initialize Go-Micro service client
    service := micro.NewService(
        micro.Transport(t),
    )
    service.Init()

    client := service.Client()

    // ---- Add Item ----
    addReq := &AddItemRequest{UserID: "user1", ProductID: "p1", Quantity: 2}
    addResp := &CartResponse{}
    err := client.Call(context.Background(),
        client.NewRequest("cart.service", "CartHandler.AddItem", addReq),
        addResp,
    )
    if err != nil {
        fmt.Println("Error adding item:", err)
    } else {
        fmt.Println("Cart after AddItem:", addResp.Items)
    }

    // ---- Add another item ----
    addReq2 := &AddItemRequest{UserID: "user1", ProductID: "p2", Quantity: 3}
    addResp2 := &CartResponse{}
    err = client.Call(context.Background(),
        client.NewRequest("cart.service", "CartHandler.AddItem", addReq2),
        addResp2,
    )
    if err != nil {
        fmt.Println("Error adding second item:", err)
    } else {
        fmt.Println("Cart after adding second item:", addResp2.Items)
    }

    // ---- Get Cart ----
    getReq := &GetCartRequest{UserID: "user1"}
    getResp := &CartResponse{}
    err = client.Call(context.Background(),
        client.NewRequest("cart.service", "CartHandler.GetCart", getReq),
        getResp,
    )
    if err != nil {
        fmt.Println("Error getting cart:", err)
    } else {
        // Pretty print JSON to check contents
        b, _ := json.MarshalIndent(getResp.Items, "", "  ")
        fmt.Println("GetCart:", string(b))
    }

    // ---- Remove Item ----
    removeReq := &RemoveItemRequest{UserID: "user1", ProductID: "p1"}
    removeResp := &CartResponse{}
    err = client.Call(context.Background(),
        client.NewRequest("cart.service", "CartHandler.RemoveItem", removeReq),
        removeResp,
    )
    if err != nil {
        fmt.Println("Error removing item:", err)
    } else {
        fmt.Println("Cart after RemoveItem:", removeResp.Items)
    }

    fmt.Println("\nNow restart the server and run GetCart again to see persistence...")
    time.Sleep(1 * time.Second)
}
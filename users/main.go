package main

import "context"

func main() {
	store := NewStore()
	svc := NewService(store)

	_ = svc.CreateUser(context.Background())
}

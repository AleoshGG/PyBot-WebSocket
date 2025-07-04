package models

type QueueBind struct {
	Name       string
	RoutingKey string
	Exchange   string
}
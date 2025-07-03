package models

type Client struct {
	ID   string
	Send chan []byte
}
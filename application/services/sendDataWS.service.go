package services

import "PyBot-WebSocket/domain/repositories"

type SendDataWS struct {
	sd repositories.ISendDataWS
}

func NewSendDataWS(sd repositories.ISendDataWS) *SendDataWS {
	return &SendDataWS{sd: sd}
}

func (sd *SendDataWS) Run(data []byte) {
	sd.sd.SendData(data)
}
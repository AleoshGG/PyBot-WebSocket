package repositories

type ISendDataWS interface {
	SendData(data []byte)
}
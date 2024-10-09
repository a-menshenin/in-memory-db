package network

//go:generate mockgen -source=./interfaces.go -destination=./mock/handler.go -package=mock Handler
//go:generate mockgen -destination=./mock/conn.go -package=mock net Conn

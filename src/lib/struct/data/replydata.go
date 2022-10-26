package data

import (
	"fmt"
	"net"
)


type ReplyData struct {
	ReplyTo net.UDPAddr
	Conn net.UDPConn
	Msg CommData
}

func (data *CommData) ReplyData() {
	fmt.Println("=== Communication data ===")
	fmt.Println("Identifier:", data.Identifier)
	fmt.Println("SenderIP:", data.SenderIP)
	fmt.Println("ReceiverIP:", data.ReceiverIP)
	fmt.Println("Message ID:", data.MsgID)
	fmt.Println("DataType:", data.DataType)
	fmt.Println("DataValue:", data.DataValue)
}
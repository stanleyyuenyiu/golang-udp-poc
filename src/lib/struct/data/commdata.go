package data

import (
	"fmt"
	"net"
	"encoding/json"
	"log"
)

type CommData struct {
	ReplyTo net.UDPAddr 
	Identifier string
	SenderIP	string
	ReceiverIP	string
	MsgID string
	DataType string
	DataValue any
}

func (data *CommData) PrintData() {
	fmt.Println("=== Communication data ===")
	fmt.Println("Identifier:", data.Identifier)
	fmt.Println("SenderIP:", data.SenderIP)
	fmt.Println("ReceiverIP:", data.ReceiverIP)
	fmt.Println("Message ID:", data.MsgID)
	fmt.Println("DataType:", data.DataType)
	fmt.Println("DataValue:", data.DataValue)
}


func (data *CommData) Marshal() ([]byte) {
	convMsg, err := json.Marshal(*data)
	if err != nil {
		log.Printf("Convert json error: %v", err.Error())
		return nil
	}
	return convMsg
}


func (data *CommData) Unmarshal(buffer []byte)  (*CommData)  {
	err := json.Unmarshal(buffer, data);
	if err != nil {
		log.Printf("Convert json error: %v", err.Error())
	}
	return data
}

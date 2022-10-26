package data

import (
	"fmt"
	"time"
)

type ConnData struct {
	SenderIP string
	MsgID string
	SendTime time.Time
	Status string
}

func (data *ConnData) PrintData() {
	fmt.Println("=== Connection data ===")
	fmt.Println("SenderIP:", data.SenderIP)
	fmt.Println("Message ID:", data.MsgID)
	fmt.Println("Time:", data.SendTime)
	fmt.Println("Status:", data.Status)
}
package udp

import (
	"log"
	"net"
	"time"
	c "udppoc/lib/struct/data"
	u "udppoc/lib/utils"
)

const com_id = "listener_id"

func BuildMsg(receiverIP string, msgID string, dataType string, dataValue any) *c.CommData {
	return &c.CommData{
		Identifier: com_id,
		SenderIP:   u.GetLocalIP(),
		ReceiverIP: receiverIP,
		MsgID:      msgID,
		DataType:   dataType,
		DataValue:  dataValue,
	}
}

func Listen(sendCh chan c.CommData, listenAddr string, listenPort int) (<-chan c.ConnData, <-chan c.CommData) {
	commSend := make(chan c.CommData)
	commReply := make(chan c.CommData)

	commReceive := make(chan c.CommData, 1)
	commSentStatus := make(chan c.ConnData)
	receivedMsg := make(chan c.CommData)

	go listen(commReceive, commReply, listenPort, false)
	go messageFowarder(commReceive, commReply, receivedMsg, commSentStatus, commSend, sendCh)

	return commSentStatus, receivedMsg
}

func Send(sendCh chan c.CommData, targetAddr string, targetPort int) (<-chan c.ConnData, <-chan c.CommData) {
	commSend := make(chan c.CommData)
	commReply := make(chan c.CommData)
	commReceive := make(chan c.CommData, 1)
	commSentStatus := make(chan c.ConnData)
	receivedMsg := make(chan c.CommData)

	go broadcast(commSend, commReceive, targetAddr, targetPort)
	go messageFowarder(commReceive, commReply, receivedMsg, commSentStatus, commSend, sendCh)

	return commSentStatus, receivedMsg
}

func closeConn(conn *net.UDPConn) {
	log.Printf("Close connection")
	conn.Close()
}

func broadcast(commSend chan c.CommData, commReceive chan c.CommData, targetAddr string, targetPort int) {
	local := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	boardcastAddr := &net.UDPAddr{IP: net.ParseIP(targetAddr), Port: targetPort}
	conn, err := net.ListenUDP("udp", local)
	defer closeConn(conn)

	if err != nil {
		panic(err)
	}

	log.Printf("Server boardcasting from %v to %v", local, boardcastAddr)

	ch := make(chan struct{})

	go writeFromCn(conn, boardcastAddr, commSend)

	go read(conn, commReceive)

	<-ch
}

func listen(commReceive chan c.CommData, commReply chan c.CommData, port int, readOnly bool) {
	listenAddr := &net.UDPAddr{IP: net.IPv4zero, Port: port}
	conn, err := net.ListenUDP("udp", listenAddr)
	defer closeConn(conn)

	if err != nil {
		panic(err)
	}

	log.Printf("Server listening to %v", listenAddr)

	ch := make(chan struct{})

	go readAndReply(conn, commReceive)

	go reply(conn, commReply)

	<-ch
}

func messageFowarder(
	commReceive <-chan c.CommData,
	commReply chan<- c.CommData,
	receivedMsg chan<- c.CommData,
	commSentStatus chan<- c.ConnData,
	commSend chan<- c.CommData,
	sendCh <-chan c.CommData) {
	for {
		select {
		case msg := <-commReceive:

			if msg.DataType == "Received" {
				log.Printf("Receive confriming message")
				response := c.ConnData{
					SenderIP: msg.SenderIP,
					MsgID:    msg.MsgID,
					SendTime: time.Now(),
					Status:   "Received",
				}
				commSentStatus <- response
			} else {
				log.Printf("Receive message and reply ack to %v", msg.ReplyTo)
				response := c.CommData{
					Identifier: com_id,
					SenderIP:   u.GetLocalIP(),
					ReceiverIP: msg.SenderIP,
					MsgID:      msg.MsgID,
					DataType:   "Received",
					DataValue:  time.Now(),
					ReplyTo:    msg.ReplyTo,
				}
				commReply <- response
				receivedMsg <- msg
			}

		case msg := <-sendCh:
			log.Printf("Sending message")
			timeSent := c.ConnData{
				SenderIP: u.GetLocalIP(),
				MsgID:    msg.MsgID,
				SendTime: time.Now(),
				Status:   "Sent",
			}
			commSentStatus <- timeSent
			commSend <- msg
		}
	}
}

func writeFromCn(conn *net.UDPConn, to *net.UDPAddr, ch chan c.CommData) {
	for {
		msg := <-ch
		conn.WriteToUDP(msg.Marshal(), to)
	}
}

func reply(conn *net.UDPConn, ch chan c.CommData) {
	for {
		msg := <-ch
		conn.WriteToUDP(msg.Marshal(), &msg.ReplyTo)
	}
}

func read(conn *net.UDPConn, ch chan c.CommData) {
	var msg c.CommData
	for {
		buffer := make([]byte, 4096)
		length, _, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Println(err)
		}
		msg.Unmarshal(buffer[:length])
		ch <- msg
	}
}

func readAndReply(conn *net.UDPConn, ch chan c.CommData) {
	var msg c.CommData
	for {
		buffer := make([]byte, 4096)
		length, replyTo, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println(err)
		}
		msg.Unmarshal(buffer[:length])
		msg.ReplyTo = *replyTo
		ch <- msg
	}
}

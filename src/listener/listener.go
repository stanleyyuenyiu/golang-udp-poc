package listener

import (
	"flag"
	udp "udppoc/lib/udp"
	c "udppoc/lib/struct/data"
	
)

var (
	port = flag.Int("port", 6001, "The server port")
)
func Run() {
	flag.Parse()

	sendCh := make(chan c.CommData)
	connStatus, msg := udp.Listen(sendCh, "", *port)
	
	ch := make(chan struct{})
	
	go RunPrintConn(connStatus)
	
	go RunPrintMsg(msg)
	
	<-ch
}

func RunPrintConn(connStatus <-chan c.ConnData) {
	for{
		data := <- connStatus
    	data.PrintData()
	}
}

func RunPrintMsg(msg <-chan c.CommData) {
	for{
		data := <- msg
    	data.PrintData()
	}
}

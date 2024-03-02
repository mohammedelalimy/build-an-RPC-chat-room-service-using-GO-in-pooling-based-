package main
import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
	"time"
)
type Message struct {
	ID      int
	Content string
	Sender  string
	Time    time.Time
}

type MessageServer struct {
	messages []Message
	mu       sync.Mutex
}

func (ms *MessageServer) Send_Message(msg Message, reply *int) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	msg.ID = len(ms.messages) + 1
	msg.Time = time.Now()
	ms.messages = append(ms.messages, msg)
	fmt.Printf("Received message from %s: %s\n", msg.Sender, msg.Content)
	return nil
}

func (ms *MessageServer) RefreshMyChat(timestamp time.Time, messages *[]Message) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	for _, msg := range ms.messages {
		if msg.Time.After(timestamp) {
			*messages = append(*messages, msg)
		}
	}

	return nil
}

func main() {
	messageServer := new(MessageServer)
	rpc.Register(messageServer)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 1234...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

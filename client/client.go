package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"time"
)

type Message struct {
	ID      int
	Content string
	Sender  string
	Time    time.Time
}

func main() {
	client, err := rpc.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer client.Close()

	// Get user => username

	fmt.Print("username: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	username := scanner.Text()

	for {
		// Get user_message
		fmt.Print("Enter your message (to quit => type 'exit'): ")
		scanner.Scan()
		messageContent := scanner.Text()

		// if  user type exit
		if messageContent == "exit" {
			break
		}

		// Create message
		newMessage := Message{
			Content: messageContent,
			Sender:  username,
		}
		var reply int

		// Call rpc function to send message
		err := client.Call("MessageServer.Send_Message", newMessage, &reply)
		if err != nil {
			fmt.Println("Error in sending :", err)
			return
		}

// Refresh (pool mechanism) chat every 2 seconds
time.Sleep(2 * time.Second)
sinceTimestamp := time.Now().Add(-time.Minute) // Fetch messages from the last minute
var messages []Message
err = client.Call("MessageServer.RefreshMyChat", sinceTimestamp, &messages)
if err != nil {
    fmt.Println("Error fetching messages:", err)
    return
}

		fmt.Println("Received messages:")
		for _, msg := range messages {
			fmt.Printf("[%d] %s (%s): %s\n", msg.ID, msg.Time.Format("2006-01-02 15:04:05"), msg.Sender, msg.Content)
		}
	}
}
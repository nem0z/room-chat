package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nem0z/room-chat/src/client"

	"github.com/nem0z/room-chat/src/server"
)

func Handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	go server.Start(nil)

	c, err := client.New("localhost:50051", nil)
	Handle(err)

	err = c.PostMessage("general", "hello")
	Handle(err)
	log.Println("Message sent successfully!")

	time.Sleep(time.Second)

	messages, err := c.GetMessagesByTag("general")
	Handle(err)
	fmt.Println("Messages: ")
	for _, msg := range messages {
		fmt.Println(msg.Prettify())
	}

	select {}
}

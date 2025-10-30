package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	connectionString := "amqp://guest:guest@localhost:5672/"

	connection, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatalf("Can't connect to rabbitmq: %v", err)
	}
	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("Can't create channel %v", err)
	}
	defer connection.Close()
	log.Println("Connection success")

	gamelogic.PrintServerHelp()
	for {
		userInput := gamelogic.GetInput()
		if len(userInput) == 0 {
			gamelogic.PrintServerHelp()
			continue
		}

		var isPaused bool

		switch userInput[0] {
		case "pause":
			isPaused = true
		case "resume":
			isPaused = false
		case "quit":
			log.Println("Goodbuy!")
			return
		case "help":
			gamelogic.PrintServerHelp()
			continue
		default:
			log.Printf("Unknown command: %s", userInput[0])
			gamelogic.PrintServerHelp()
			continue
		}

		log.Printf("Sending %s message. debug -> %t", userInput[0], isPaused)
		err := pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
			IsPaused: isPaused,
		})
		if err != nil {
			log.Printf("can't publish message: %v", err)
		}

		_, _, declareErr := pubsub.DeclareAndBind(connection, routing.ExchangePerilTopic, routing.GameLogSlug, "game_logs", pubsub.Durable)

		if declareErr != nil {
			log.Printf("can't declare topic: %v", err)
		}
	}
}

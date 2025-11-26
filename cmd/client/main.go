package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func main() {
	fmt.Println("Starting Peril client...")

	connectionString := "amqp://guest:guest@localhost:5672/"

	connection, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatalf("Can't connect to rabbitmq: %v", err)
	}
	defer connection.Close()
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Can't get username: %v", err)
	}
	_, _, pubErr := pubsub.DeclareAndBind(connection, routing.ExchangePerilDirect, fmt.Sprintf("%s.%s", routing.PauseKey, username), routing.PauseKey, pubsub.Transient)
	if pubErr != nil {
		log.Fatalf("Error during declaration for topic: %v", pubErr)
	}

	gameState := gamelogic.NewGameState(username)
	pauseErr := pubsub.SubscribeJSON(connection, routing.ExchangePerilDirect, fmt.Sprintf("pause.%s", username), routing.PauseKey, pubsub.Transient, handlerPause(gameState))
	if pauseErr != nil {
		log.Fatalf("Can't subscribe to pause topic: %v", pauseErr)
	}
	for {
		userInput := gamelogic.GetInput()
		if len(userInput) == 0 {
			continue
		}
		switch userInput[0] {
		case "move":
			_, err := gameState.CommandMove(userInput)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case "spawn":
			err := gameState.CommandSpawn(userInput)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case "status":
			gameState.CommandStatus()

		case "help":
			gamelogic.PrintClientHelp()

		case "spam":
			fmt.Println("Spamming not allowed yet!")

		case "quit":
			gamelogic.PrintQuit()
			return

		default:
			fmt.Println("unknown command")
		}
	}
}

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(ps routing.PlayingState) {
		defer fmt.Printf("> ")
		gs.HandlePause(ps)
	}
}

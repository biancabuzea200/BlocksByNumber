package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")

	}
	WSS := os.Getenv("WSS")
	DSN := os.Getenv("DSN")
	fmt.Println("Dialing...")
	client, err := ethclient.Dial(WSS)
	if err != nil {
		log.Fatalln("Could not connect to rpc client")
	}
	fmt.Println("Successfully dialed")
	ctx := context.Background()
	headerChannel := make(chan *types.Header)
	fmt.Println("starting subscription")
	_, err = client.SubscribeNewHead(ctx, headerChannel)
	for err != nil {
		_, err := client.SubscribeNewHead(ctx, headerChannel)
		if err == nil {
			fmt.Println("subscribed!")
		}
	}
	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	for err != nil {
		fmt.Println("could not connect to db, trying again")
		db, err = gorm.Open(postgres.Open(DSN), &gorm.Config{})
	}
	fmt.Println("successfully connected to db")
	fmt.Println(db)

	for {
		head := <-headerChannel
		fmt.Println("neswest block no is  ", head.Number)
	}

}

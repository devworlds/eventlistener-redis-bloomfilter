package listener

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func Connect() *ethclient.Client {
	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/YOUR_API_KEY")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("EthereumClient Successfully connected!")
	return client
}

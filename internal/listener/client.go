package listener

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func Connect() *EthClient {
	ethC, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/YOU_API_KEY")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("EthereumClient Successfully connected!")
	return &EthClient{
		client: ethC,
	}
}

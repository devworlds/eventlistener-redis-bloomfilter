package listener

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func Connect() *EthClient {
	ethC, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/546c10bde2474f839967d30f45a53bdf")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("EthereumClient Successfully connected!")
	return &EthClient{
		client: ethC,
	}
}

package listener

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	token "github.com/devworlds/eventlistener-redis-performance/internal/listener/abi"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Listen() {
	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/YOUR_API_KEY")
	if err != nil {
		log.Fatal(err)
	}

	lastBlock, _ := client.BlockByNumber(context.Background(), nil)

	contractAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	for {
		query := ethereum.FilterQuery{
			FromBlock: lastBlock.Number(),
			Addresses: []common.Address{
				contractAddress,
			},
		}

		logs := make(chan types.Log)
		sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
		if err != nil {
			log.Printf("Erro ao se inscrever para logs: %v. Tentando reconectar...", err)
			time.Sleep(10 * time.Second) // Espera antes de tentar reconectar
			continue
		}

		logTransferSig := []byte("Transfer(address,address,uint256)")
		logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

		contractAbi, err := abi.JSON(strings.NewReader(string(token.TokenABI)))
		if err != nil {
			log.Fatal(err)
		}

		for {
			select {
			case err := <-sub.Err():
				log.Printf("Erro na subscrição: %v. Reconectando...", err)
				sub.Unsubscribe()
				time.Sleep(5 * time.Second)

			case vLog := <-logs:
				switch vLog.Topics[0].Hex() {
				case logTransferSigHash.Hex():
					fmt.Printf("Log Name: Transfer\n")
					event := struct {
						From  common.Address
						To    common.Address
						Value *big.Int
					}{}

					err := contractAbi.UnpackIntoInterface(&event, "Transfer", vLog.Data)
					if err != nil {
						log.Printf("Erro ao descompactar evento: %v", err)
						continue
					}

					event.From = common.HexToAddress(vLog.Topics[1].Hex())
					event.To = common.HexToAddress(vLog.Topics[2].Hex())

					eventLog := struct {
						From  string
						To    string
						Value string
					}{}
					eventLog.From = event.From.Hex()
					eventLog.To = event.To.Hex()
					eventLog.Value = event.Value.String()

					fmt.Printf("%+v\n", eventLog)
				}
			}
		}
	}
}
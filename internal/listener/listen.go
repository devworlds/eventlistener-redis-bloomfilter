package listener

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/devworlds/eventlistener-redis-performance/internal/db/repository"
	token "github.com/devworlds/eventlistener-redis-performance/internal/listener/abi"
	"github.com/devworlds/eventlistener-redis-performance/internal/redis"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func Listen(ethc *EthClient, filter *bloom.BloomFilter, rdb *redis.RedisService, ctx context.Context) {
	lastBlock, _ := ethc.client.BlockByNumber(context.Background(), nil)
	wallets, _ := repository.GetWallets()
	for _, addr := range wallets {
		filter.Add([]byte(addr.Address))
	}
	jsonFilter, err := filter.MarshalJSON()
	if err != nil {
		log.Fatal(err.Error())
	}
	rdb.SetWallet(string(jsonFilter))
	println("Wallets to listen: ", filter.ApproximatedSize())
	contractAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	for {
		query := ethereum.FilterQuery{
			FromBlock: lastBlock.Number(),
			Addresses: []common.Address{
				contractAddress,
			},
		}
		logs := make(chan types.Log)
		sub, err := ethc.client.SubscribeFilterLogs(context.Background(), query, logs)
		if err != nil {
			log.Printf("Error to subscribe logs: %v. Trying to reconnect...", err)
			time.Sleep(10 * time.Second)
			continue
		}

		logTransferSig := []byte("Transfer(address,address,uint256)")
		logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

		contractAbi, err := abi.JSON(strings.NewReader(string(token.TokenABI)))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Listening events...\n")
		for {
			select {
			case err := <-sub.Err():
				log.Printf("Error in subscriber: %v. Reconnecting...", err)
				sub.Unsubscribe()
				time.Sleep(5 * time.Second)

			case vLog := <-logs:
				switch vLog.Topics[0].Hex() {
				case logTransferSigHash.Hex():
					if os.Getenv("CACHE") == "true" {
						event := struct {
							From  common.Address
							To    common.Address
							Value *big.Int
						}{}

						err := contractAbi.UnpackIntoInterface(&event, "Transfer", vLog.Data)
						if err != nil {
							log.Printf("Error to event unzip: %v", err)
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

						if filter.Test([]byte(event.From.Hex())) {
							wallets, _ := repository.GetWallets()
							for _, addr := range wallets {
								if addr.Address == event.From.Hex() {
									fmt.Printf("%+v\n", eventLog)
								}
							}
						}
						if filter.Test([]byte(event.To.Hex())) {
							wallets, _ := repository.GetWallets()
							for _, addr := range wallets {
								if addr.Address == event.To.Hex() {
									fmt.Printf("%+v\n", eventLog)
								}
							}
						}
					} else if os.Getenv("CACHE") == "false" {
						wallets, _ := repository.GetWallets()

						event := struct {
							From  common.Address
							To    common.Address
							Value *big.Int
						}{}

						err := contractAbi.UnpackIntoInterface(&event, "Transfer", vLog.Data)
						if err != nil {
							log.Printf("Error to event unzip: %v", err)
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

						for _, addr := range wallets {
							if addr.Address == event.From.Hex() {
								fmt.Printf("%+v\n", eventLog)
							}
							if addr.Address == event.To.Hex() {
								fmt.Printf("%+v\n", eventLog)
							}
						}
					}
				}
			}
		}
	}
}

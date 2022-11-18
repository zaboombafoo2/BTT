package priceoracle

import (
	"context"
	"errors"
	"math/big"
	"sync"

	conabi "github.com/bittorrent/go-btfs/chain/abi"
	"github.com/bittorrent/go-btfs/transaction"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	errDecodeABI = errors.New("could not decode abi data")
)

type service struct {
	priceOracleAddress common.Address
	transactionService transaction.Service
}

type Service interface {
	// CurrentPrice CurrentRate CurrentTotalPrice get cached info from memory.
	CurrentPrice(token common.Address) (*big.Int, error)
	CurrentRate(token common.Address) (*big.Int, error)
	CurrentTotalPrice(token common.Address) (*big.Int, error)

	// CheckNewPrice retrieves latest available information from oracle
	CheckNewPrice(token common.Address) (*big.Int, error)
}

var (
	priceOracleABI = transaction.ParseABIUnchecked(conabi.OracleAbi2)

	curMutex        sync.Mutex
	mpCurPrice      = make(map[common.Address]*big.Int)
	curRate         = big.NewInt(0)
	mpCurTotalPrice = make(map[common.Address]*big.Int)
)

func New(priceOracleAddress common.Address, transactionService transaction.Service) Service {
	return &service{
		priceOracleAddress: priceOracleAddress,
		transactionService: transactionService,
	}
}

func (s *service) CurrentPrice(token common.Address) (price *big.Int, err error) {
	curMutex.Lock()
	price = big.NewInt(0).Set(mpCurPrice[token])
	curMutex.Unlock()

	return price, nil
}
func (s *service) CurrentRate(token common.Address) (rate *big.Int, err error) {
	curMutex.Lock()
	rate = big.NewInt(0).Set(curRate)
	curMutex.Unlock()

	return rate, nil
}
func (s *service) CurrentTotalPrice(token common.Address) (totalPrice *big.Int, err error) {
	curMutex.Lock()
	totalPrice = big.NewInt(0).Set(mpCurTotalPrice[token])
	curMutex.Unlock()

	return totalPrice, nil
}

func (s *service) CheckNewPrice(token common.Address) (*big.Int, error) {
	price, err := s.currentPrice(token)
	if err != nil {
		return nil, err
	}
	//fmt.Println("currentPrice ", price)

	rate, err := s.currentRate()
	if err != nil {
		return nil, err
	}
	//fmt.Println("currentRate ", rate)

	curMutex.Lock()
	defer curMutex.Unlock()
	mpCurPrice[token] = price
	curRate = rate
	mpCurTotalPrice[token] = big.NewInt(0).Mul(price, rate)

	return big.NewInt(0).Set(mpCurTotalPrice[token]), nil
}

func (s *service) currentRate() (*big.Int, error) {
	callData, err := priceOracleABI.Pack("getExchangeRate")
	if err != nil {
		return nil, err
	}
	result, err := s.transactionService.Call(context.Background(), &transaction.TxRequest{
		To:   &s.priceOracleAddress,
		Data: callData,
	})
	if err != nil {
		return nil, err
	}

	results, err := priceOracleABI.Unpack("getExchangeRate", result)
	if err != nil {
		return nil, err
	}

	if len(results) != 1 {
		return nil, errDecodeABI
	}

	rate, ok := abi.ConvertType(results[0], new(big.Int)).(*big.Int)
	if !ok || rate == nil {
		return nil, errDecodeABI
	}

	//fmt.Println("currentRate, rate = ", rate)

	return rate, nil
}

func (s *service) currentPrice(token common.Address) (*big.Int, error) {
	callData, err := priceOracleABI.Pack("getPrice", token)
	if err != nil {
		return nil, err
	}
	result, err := s.transactionService.Call(context.Background(), &transaction.TxRequest{
		To:   &s.priceOracleAddress,
		Data: callData,
	})
	if err != nil {
		return nil, err
	}
	results, err := priceOracleABI.Unpack("getPrice", result)
	if err != nil {
		return nil, err
	}
	if len(results) != 1 {
		return nil, errDecodeABI
	}
	price, ok := abi.ConvertType(results[0], new(big.Int)).(*big.Int)
	if !ok || price == nil {
		return nil, errDecodeABI
	}

	//fmt.Println("currentPrice, price = ", price)

	return price, nil
}

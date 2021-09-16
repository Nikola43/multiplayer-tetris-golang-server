package ethbasedclient

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nikola43/tetrisMultiplayer/utils/errorsutil"
	"github.com/nikola43/tetrisMultiplayer/utils/ethutil"
	"golang.org/x/crypto/sha3"
	"log"
	"math/big"
)

type EthBasedClient struct {
	Client     *ethclient.Client
	PrivateKey *ecdsa.PrivateKey
	Address    common.Address
	ChainID    *big.Int
	Transactor *bind.TransactOpts
	Nonce      uint64
}

func New(rawUrl, plainPrivateKey string) EthBasedClient {
	client, err := ethclient.Dial(rawUrl)
	errorsutil.HandleError(err)

	privateKey := ethutil.GenerateEcdsaPrivateKey(plainPrivateKey)
	ethBasedClientTemp := EthBasedClient{
		Client:     client,
		PrivateKey: privateKey,
		Address:    ethutil.GenerateAddress(privateKey),
		ChainID:    ethutil.GetChainID(client),
		Transactor: ethutil.GenerateTransactor(client, privateKey),
	}

	return ethBasedClientTemp
}

func (ethBasedClient *EthBasedClient) ConfigureTransactor(value *big.Int, gasPrice *big.Int, gasLimit uint64) {

	if value.String() != "-1" {
		ethBasedClient.Transactor.Value = value
	}

	ethBasedClient.Transactor.GasPrice = gasPrice
	ethBasedClient.Transactor.GasLimit = gasLimit
	ethBasedClient.Transactor.Nonce = ethBasedClient.PendingNonce()
	ethBasedClient.Transactor.Context = context.Background()
}

func (ethBasedClient *EthBasedClient) Balance() *big.Int {
	balance, balanceErr := ethBasedClient.Client.BalanceAt(context.Background(), ethBasedClient.Address, nil)
	errorsutil.HandleError(balanceErr)
	return balance
}

func (ethBasedClient *EthBasedClient) TransferTokens(tokenAddress common.Address, toAddress common.Address, amount *big.Int) string {

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit, err := ethBasedClient.Client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	})
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, SuggestGasPriceErr := ethBasedClient.Client.SuggestGasPrice(context.Background())
	if SuggestGasPriceErr != nil {
		log.Fatal(SuggestGasPriceErr)
	}

	tx := types.NewTransaction(ethBasedClient.Nonce, tokenAddress, amount, gasLimit, gasPrice, data)

	chainID, err := ethBasedClient.Client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, signTxErr := types.SignTx(tx, types.NewEIP155Signer(chainID), ethBasedClient.PrivateKey)
	if signTxErr != nil {
		log.Fatal(signTxErr)
	}

	err = ethBasedClient.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	return signedTx.Hash().Hex()
}

func (ethBasedClient *EthBasedClient) SuggestGasPrice() *big.Int {
	gasPrice, err := ethBasedClient.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return gasPrice
}

func (ethBasedClient *EthBasedClient) PendingNonce() *big.Int {
	nonce, nonceErr := ethBasedClient.Client.PendingNonceAt(context.Background(), ethBasedClient.Address)
	errorsutil.HandleError(nonceErr)
	return big.NewInt(int64(nonce))
}

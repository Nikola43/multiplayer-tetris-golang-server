package main

import (
	"fmt"
	"github.com/nikola43/tetrisMultiplayer/game"
	"github.com/nikola43/tetrisMultiplayer/ethbasedclient"
)

func main()  {
	//rawUrl := "https://data-eed-prebsc-1-s1.binance.org:8545/"
	rawUrl := "https://rinkeby.infura.io/v3/dd4857da75ac450e8422383943558b43"
	plainPrivateKey := "fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19"
	ethBasedClient := ethbasedclient.New(rawUrl, plainPrivateKey)



	fmt.Println(ethBasedClient)
	a := game.Game{}
	a.Initialize(":3001")
}

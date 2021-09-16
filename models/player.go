package models

type Player struct {
	UUID string
	WalletAddress string
	IsPlaying bool
}

func PlayerExists(players []*Player, walletAddress string) int {
	for i, v := range players {
		if v.WalletAddress == walletAddress {
			return i
		}
	}
	return -1
}



package loaders

import "pulse/internal/config"

type SolanaLoader struct {
	apiKeys ApiKeys
}

type ApiKeys struct {
	blockdaemon string
	coinbase    string
}

func NewSolanaLoader(cfg *config.Config) *SolanaLoader {
	apiKeys := ApiKeys{blockdaemon: cfg.ApiKeys.Blockdaemon, coinbase: cfg.ApiKeys.Coinbase}
	return &SolanaLoader{apiKeys}
}

func (s *SolanaLoader) LoadSolWallet(address string) error {
	// TODO: load the solana address and return the tokens and their balances. This should only be called when we create a new Subwallet which is of type Solana
	return nil
}

func (s *SolanaLoader) LoadCurrentSolanaPrice() (float64, error) {
	//TODO: load from coinbase and parse into struct
	//
	//
	//TODO: calculateUsdPrice here
	price, err := calculateUsdPrice(0, 0)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func (s *SolanaLoader) LoadTokenPrice(address string) (float64, error) {
	//TODO: load from blockdaemon and parse into struct
	//
	//
	//TODO: calculateUsdPrice here
	price, err := calculateUsdPrice(0, 0)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func calculateUsdPrice(price float64, decimals int) (float64, error) {
	return 0, nil
}

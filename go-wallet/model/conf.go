package model

type CommonConf struct {
	RedisTtl int `json:"redisttl"`
}

type WalletConf struct {
	CommonConf
}

type BalanceConf struct {
	CommonConf
	Threshold int `json:"threshold"`
}

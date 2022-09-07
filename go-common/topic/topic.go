package topic

import "github.com/lovoo/goka"

type EmitterTopic string

func (e EmitterTopic) String() string {
	return string(e)
}

func (e EmitterTopic) GokaStream() goka.Stream {
	return goka.Stream(e)
}

const (
	BALANCE_TRANSACTION_TOPIC EmitterTopic = "balance-transaction" // go-gateway to go-wallet
	WALLET_TRANSACTION_TOPIC  EmitterTopic = "wallet-transaction"  // go-gateway to go-wallet
)

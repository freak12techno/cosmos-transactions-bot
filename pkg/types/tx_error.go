package types

import (
	"main/pkg/data_fetcher"
)

type TxError struct {
	Error error
}

func (txError TxError) GetMessages() []Message {
	return []Message{}
}

func (txError TxError) Type() string {
	return "TxError"
}

func (txError TxError) GetHash() string {
	return "TxError"
}

func (txError *TxError) GetAdditionalData(fetcher data_fetcher.DataFetcher) {
}

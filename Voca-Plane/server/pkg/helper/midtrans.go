package helper

import (
	"crypto/sha512"
	"encoding/hex"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type MidTransClient struct {
	snap snap.Client
	key string
}

func NewMidTransClient(serverKey string, isProd bool) *MidTransClient {
	env := midtrans.Sandbox
	if isProd {
		env = midtrans.Production
	}

	var client snap.Client
	client.New(serverKey, env)

	return &MidTransClient{
		snap: client,
		key: serverKey,
	}
}

func (m *MidTransClient) CreatePayment(orderID string, amount float64) (*snap.Response, error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(amount),
		},
	}

	res, err := m.snap.CreateTransaction(req)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (m *MidTransClient) VerifySignature(orderID, statusCode, grossAmount, signature string) bool {
	data := orderID + statusCode + grossAmount + m.key

	hash := sha512.Sum512([]byte(data))
	expected := hex.EncodeToString(hash[:])

	return expected == signature
}
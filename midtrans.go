package main

import (
	"fmt"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type MidtransService struct {
	Client *coreapi.Client
}

func NewMidtrans(client_key, server_key string) MidtransService {
	var client coreapi.Client
	client.New(server_key, midtrans.Sandbox)
	client.ClientKey = client_key
	return MidtransService{ Client: &client }
}

func (s *MidtransService) CreateTransaction(req TransactionReq) (*TransactionRes, error) {
	res, err := s.Client.ChargeTransaction(&coreapi.ChargeReq{
		PaymentType: coreapi.CoreapiPaymentType(req.PaymentType),
		BankTransfer: &coreapi.BankTransferDetails{
			Bank: midtrans.Bank(req.BankName),
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID: req.OrderId,
			GrossAmt: req.Amount,
		},
	})
	if err != nil {
		return nil, err
	}

	fmt.Printf("id: %s", res.TransactionID)

	p := &TransactionRes{
		Status: PAYMENT_PENDING,
		Id: res.TransactionID,
		Message: res.StatusMessage,
		BankName: BankName(res.Bank),
		Time: res.TransactionTime,
	}

	return p, nil
}

func (s *MidtransService) CancelTransaction(req TransactionReq) (*TransactionRes, error) {
	res, err := s.Client.CancelTransaction(req.OrderId)
	if err != nil {
		return nil, err
	}

	p := &TransactionRes{
		Id: res.TransactionID,
		Message: res.StatusMessage,
		BankName: BankName(res.Bank),
		Time: res.TransactionTime,
		Status: PAYMENT_CANCELED,
	}

	return p, nil
}

func (s *MidtransService) RefundTransaction(req TransactionReq) (*TransactionRes, error) {
	res, err := s.Client.RefundTransaction(req.OrderId, &coreapi.RefundReq{
		RefundKey: "",
		Amount: req.Amount,
		Reason: "",
	})
	if err != nil {
		return nil, err
	}

	p := &TransactionRes{
		Id: res.TransactionID,
		Status: PAYMENT_REFUND,
		Message: res.StatusMessage,
		Time: res.TransactionTime,
		FraudStatus: res.FraudStatus,
	}
	return p, nil
}

func (s *MidtransService) GetTransactionStatus(TxId string) (*TransactionRes, error) {
	res, err := s.Client.CheckTransaction(TxId)
	if err != nil {
		return nil, err
	}

	p := &TransactionRes{
		Id: res.ID,
		Message: res.StatusMessage,
		BankName: BankName(res.Bank),
		FraudStatus: res.FraudStatus,
		Time: res.TransactionTime,
		SignatureKey: res.SignatureKey,
	}
	return p, nil
}
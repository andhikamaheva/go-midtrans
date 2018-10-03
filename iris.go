package midtrans

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

// IrisGateway struct
type IrisGateway struct {
	Client Client
}

// Call : base method to call Core API
func (gateway *IrisGateway) Call(method, path string, body io.Reader, v interface{}) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	path = gateway.Client.APIEnvType.IrisURL() + path
	req, err := gateway.Client.NewRequest(method, path, body, gateway.Client.ApproverKey)
	if err != nil {
		return err
	}

	return gateway.Client.ExecuteRequest(req, v)
}

// CreateBeneficiaries : Perform transaction using ChargeReq
func (gateway *IrisGateway) CreateBeneficiaries(req *BeneficiariesReq) (map[string]interface{}, error) {
	var resp map[string]interface{}
	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("POST", "api/v1/beneficiaries", bytes.NewBuffer(jsonReq), &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error create beneficiaries: ", err)
		return resp, err
	}

	return resp, nil
}

// ValidateBankAccount : get order status using order ID
func (gateway *IrisGateway) ValidateBankAccount(bankName string, account string) (map[string]interface{}, error) {
	var resp map[string]interface{}

	err := gateway.Call("GET", "api/v1/account_validation?bank="+bankName+"&account="+account, nil, &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error approving: ", err)
		return resp, err
	}

	return resp, nil
}

// CreatePayouts : Create Payout with single or multiple payouts
func (gateway *IrisGateway) CreatePayouts(req PayoutReq) (Payout, error) {
	resp := Payout{}

	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("POST", "api/v1/payouts", bytes.NewBuffer(jsonReq), &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error create payouts: ", err)
		return resp, err
	}

	return resp, nil
}

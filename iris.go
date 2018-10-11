package midtrans

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io"
	"strings"
)

// IrisGateway struct
type IrisGateway struct {
	Client Client
}

// Call : base method to call Core API
func (gateway *IrisGateway) Call(method, path string, body io.Reader, v interface{}, key string) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	path = gateway.Client.APIEnvType.IrisURL() + path
	req, err := gateway.Client.NewRequest(method, path, body, key)
	if err != nil {
		return err
	}

	return gateway.Client.ExecuteRequest(req, v)
}

// CreateBeneficiaries : Perform transaction using ChargeReq
func (gateway *IrisGateway) CreateBeneficiaries(req *BeneficiariesReq) (map[string]interface{}, error) {
	var resp map[string]interface{}
	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("POST", "api/v1/beneficiaries", bytes.NewBuffer(jsonReq), &resp, gateway.Client.ApproverKey)
	if err != nil {
		gateway.Client.Logger.Println("Error create beneficiaries: ", err)
		return resp, err
	}

	return resp, nil
}

// ValidateBankAccount : get order status using order ID
func (gateway *IrisGateway) ValidateBankAccount(bankName string, account string) (map[string]interface{}, error) {
	var resp map[string]interface{}

	err := gateway.Call("GET", "api/v1/account_validation?bank="+bankName+"&account="+account, nil, &resp, gateway.Client.ApproverKey)
	if err != nil {
		gateway.Client.Logger.Println("Error approving: ", err)
		return resp, err
	}

	return resp, nil
}

// CreatePayouts : Create Payout with single or multiple payouts
func (gateway *IrisGateway) CreatePayouts(req *PayoutReq) (Payout, error) {
	resp := Payout{}

	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("POST", "api/v1/payouts", bytes.NewBuffer(jsonReq), &resp, gateway.Client.CreatorKey)
	if err != nil {
		gateway.Client.Logger.Println("Error create payouts: ", err)
		return resp, err
	}

	return resp, nil
}

// ApprovePayouts : Approve Payout(s) with single or multiple payouts
func (gateway *IrisGateway) ApprovePayouts(req *ApprovePayoutReq) (map[string]interface{}, error) {
	var resp map[string]interface{}

	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("POST", "api/v1/payouts/approve", bytes.NewBuffer(jsonReq), &resp, gateway.Client.ApproverKey)
	if err != nil {
		gateway.Client.Logger.Println("Error approve payouts: ", err)
		return resp, err
	}

	return resp, nil
}

// RejectPayouts : Reject Payout(s) with single or multiple payouts
func (gateway *IrisGateway) RejectPayouts(req *RejectPayoutReq) (map[string]interface{}, error) {
	var resp map[string]interface{}

	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("POST", "api/v1/payouts/reject", bytes.NewBuffer(jsonReq), &resp, gateway.Client.ApproverKey)
	if err != nil {
		gateway.Client.Logger.Println("Error reject payouts: ", err)
		return resp, err
	}

	return resp, nil
}

// ValidateSignatureKey : Validate Iris Signature from Payout Notification
func (gateway *IrisGateway) ValidateSignatureKey(payload string, headerKey string) bool {
	hasher := sha512.New()
	hasher.Write([]byte(string(payload) + gateway.Client.MerchantKey))

	signatureKey := hex.EncodeToString(hasher.Sum(nil))

	if signatureKey == headerKey {
		return true
	}
	return false
}

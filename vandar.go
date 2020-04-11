package vandar

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type VandarPayment struct {
	PaymentAPIEndpoint
	APIKey string
}

func (vp *VandarPayment) CheckData() error {
	if vp.APIKey == "" {
		return errors.New("API KEY NOT SET")
	}
	if vp.VerifyApi == "" {
		return errors.New("VerificationAPI NOT SET")
	}
	if vp.PaymentApi == "" {
		return errors.New("PaymentAPI NOT SET")
	}
	if vp.RequestApi == "" {
		return errors.New("RequestAPI NOT SET")
	}
	return nil
}

type SendRequest struct {
	apiKey      string `json:"api_key"`
	Amount      int    `json:"amount"`
	CallbackURL string `json:"callback_url"`
	Mobile      string `json:"mobile_number"`
	FactorID    string `json:"factorNumber"`
	Description string `json:"description"`
}

type PaymentAPIEndpoint struct {
	RequestApi string
	PaymentApi string
	VerifyApi  string
}

type VandarResponse interface {
	errors() string
}
type VandarRequestToken struct {
	Status int      `json:"status"`
	Token  string   `json:"token"`
	Errors []string `json:"errors"`
}

func (v *VandarRequestToken) errors() string {
	return fmt.Sprintln(v.Errors)
}

type VandarPaymentVerifyRequest struct {
	APIKey string `json:"api_key"`
	Token  string `json:"token"`
}
type VandarPaymentVerfiy struct {
	Status        int      `json:"status"`
	Amount        int      `json:"amount"`
	TransactionID string   `json:"transId"`
	FactorID      string   `json:"factorNumber"`
	Mobile        string   `json:"mobile"`
	Description   string   `json:"description"`
	CardNumber    string   `json:"cardNumber"`
	Date          string   `json:"paymentDate"`
	Message       string   `josn:"message"`
	Errors        []string `json:"errors"`
}

func (v *VandarPaymentVerfiy) errors() string {
	return fmt.Sprintln(v.Errors)
}

func (vp *VandarPayment) RequestPayment(sr *SendRequest) (string, error) {
	sr.apiKey = vp.APIKey
	fmt.Println("API", sr.apiKey)
	requestBody, err := json.Marshal(sr)

	if err != nil {
		fmt.Println("RP, marshal",err)
		return "", err
	}
	paymentRequest, err := http.Post(vp.PaymentAPIEndpoint.RequestApi, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("RP, PostRequest",err)
		return "", err
	}
	defer paymentRequest.Body.Close()

	response, err := ioutil.ReadAll(paymentRequest.Body)
	if err !=nil{
		fmt.Println("RP, responseRead",err)
		return "", err
	}
	var vr VandarRequestToken
	err = json.Unmarshal(response, &vr)
	if err != nil {
		fmt.Println("RP, UnMarshalResponse",string(response),err)
		return "", err
	}
	if vr.Status == 0 {
		fmt.Println("RP, PaymentResponse",vr)
		return "", errors.New(vr.errors())
	}
	return vp.PaymentApi + vr.Token, nil
}
func (vp *VandarPayment) VerifyPayment(token string) (*VandarPaymentVerfiy, error) {
	pr := VandarPaymentVerifyRequest{
		APIKey: vp.APIKey,
		Token:  token,
	}
	requestBody, err := json.Marshal(pr)
	if err != nil {
		return nil, err
	}
	verifyRequest, err := http.Post(vp.VerifyApi, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer verifyRequest.Body.Close()
	response, err := ioutil.ReadAll(verifyRequest.Body)
	if err != nil {
		return nil, err
	}
	var vpv VandarPaymentVerfiy
	err = json.Unmarshal(response, &vpv)
	if err != nil {
		return nil, err
	}
	if vpv.Status == 0 {
		return nil, errors.New(vpv.errors())
	}
	return &vpv, nil
}

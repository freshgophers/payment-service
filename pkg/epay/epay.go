package epay

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"html/template"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"sync"
	"time"
)

type Credential struct {
	TerminalID    string
	ClientID      string
	ClientSecret  string
	OauthEndpoint string
	Endpoint      string
	JSLink        string
	BackLink      string
	PostLink      string
	Amount        string
	AccessToken   string `json:"access_token,omitempty"`
	ExpiresIn     string `json:"expires_in,omitempty"`
	ExpiresAt     int64
}

type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type Payment struct {
	Amount          decimal.Decimal `json:"amount"`
	Currency        string          `json:"currency"`
	Name            string          `json:"name"`
	TerminalID      string          `json:"terminalId"`
	InvoiceID       string          `json:"invoiceId"`
	Description     string          `json:"description"`
	AccountID       string          `json:"accountId"`
	Email           string          `json:"email"`
	Phone           string          `json:"phone"`
	BackLink        string          `json:"backLink"`
	FailureBackLink string          `json:"failureBackLink"`
	PostLink        string          `json:"postLink"`
	FailurePostLink string          `json:"failurePostLink"`
	Language        string          `json:"language"`
	PaymentType     string          `json:"paymentType"`
	CardID          struct {
		ID string `json:"id"`
	} `json:"cardId"`
	PaymentJsLink string `json:"-"`
	Token         *Token `json:"-"`
	CardSave      string `json:"-"`
	HomebankToken string `json:"-"`
}

type Invoice struct {
	ID             string          `json:"id" db:"external_id"`
	DateTime       time.Time       `json:"dateTime" db:"date_time"`
	InvoiceID      string          `json:"invoiceId" db:"invoice_id"`
	Amount         decimal.Decimal `json:"amount" db:"amount"`
	AmountBonus    decimal.Decimal `json:"amountBonus" db:"amount_bonus"`
	Currency       string          `json:"currency" db:"currency"`
	Terminal       string          `json:"terminal" db:"terminal"`
	AccountID      string          `json:"accountId" db:"account_id"`
	Description    string          `json:"description" db:"description"`
	Language       string          `json:"language" db:"language"`
	CardMask       string          `json:"cardMask" db:"card_mask"`
	CardType       string          `json:"cardType" db:"card_type"`
	Issuer         string          `json:"issuer" db:"issuer"`
	Reference      string          `json:"reference" db:"reference"`
	IntReference   string          `json:"intReference" db:"int_reference"`
	Secure         string          `json:"secure" db:"secure"`
	Secure3D       string          `json:"secure3D" db:"secure_3d"`
	TokenRecipient string          `json:"tokenRecipient" db:"token_recipient"`
	Code           string          `json:"code" db:"code"`
	Reason         string          `json:"reason" db:"reason"`
	ReasonCode     string          `json:"reasonCode" db:"reason_code"`
	Name           string          `json:"name" db:"name"`
	Email          string          `json:"email" db:"email"`
	Phone          string          `json:"phone" db:"phone"`
	IP             string          `json:"ip" db:"ip"`
	IPCountry      string          `json:"ipCountry" db:"ip_country"`
	IPCity         string          `json:"ipCity" db:"ip_city"`
	IPRegion       string          `json:"ipRegion" db:"ip_region"`
	IPDistrict     string          `json:"ipDistrict" db:"ip_district"`
	IPLongitude    decimal.Decimal `json:"ipLongitude" db:"ip_longitude"`
	IPLatitude     decimal.Decimal `json:"ipLatitude" db:"ip_latitude"`
	CardID         string          `json:"cardID" db:"card_id"`
}

type Client struct {
	client     *http.Client
	mutex      *sync.Mutex
	credential Credential
}

func NewClient(credential Credential) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	client.Timeout = 30 * time.Second

	return &Client{
		client:     client,
		credential: credential,
		mutex:      &sync.Mutex{},
	}
}

func (s *Client) GetCredential() Credential {
	return s.credential
}

func (s *Client) GetToken() (*Token, error) {
	// setup data struct
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("grant_type", "client_credentials")
	_ = writer.WriteField("scope", "webapi usermanagement email_send verification statement statistics")
	_ = writer.WriteField("client_id", s.credential.ClientID)
	_ = writer.WriteField("client_secret", s.credential.ClientSecret)

	err := writer.Close()
	if err != nil {
		return nil, err
	}

	// setup request
	url := s.credential.OauthEndpoint + "/oauth2/token"
	respBytes, code, err := s.handler("POST", url, body.Bytes(), writer, nil)
	if err != nil {
		return nil, err
	}

	// check response code
	switch code {
	case 200:
		// unmarshal response data
		token := &Token{}
		err = json.Unmarshal(respBytes, &token)
		if err != nil {
			return nil, err
		}

		return token, nil
	default:
		return nil, errors.New(string(respBytes))
	}
}

func (s *Client) getPaymentToken(payment *Payment) (*Token, error) {
	// setup data struct
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("grant_type", "client_credentials")
	_ = writer.WriteField("scope", "payment")
	_ = writer.WriteField("client_id", s.credential.ClientID)
	_ = writer.WriteField("client_secret", s.credential.ClientSecret)
	_ = writer.WriteField("postLink", payment.PostLink)
	_ = writer.WriteField("terminal", payment.TerminalID)
	_ = writer.WriteField("currency", payment.Currency)
	_ = writer.WriteField("invoiceID", payment.InvoiceID)
	_ = writer.WriteField("amount", payment.Amount.String())

	err := writer.Close()
	if err != nil {
		return nil, err
	}

	// setup request
	url := s.credential.OauthEndpoint + "/oauth2/token"
	respBytes, code, err := s.handler("POST", url, body.Bytes(), writer, nil)
	if err != nil {
		return nil, err
	}

	// check response code
	switch code {
	case 200:
		// unmarshal response data
		token := new(Token)
		err = json.Unmarshal(respBytes, &token)
		if err != nil {
		}

		return token, nil
	default:
		return nil, errors.New(string(respBytes))
	}
}

func (s *Client) PayOnTemplate(c *fiber.Ctx, cardSave, homebankToken, insuranceID string, payment *Payment) error {
	paymentDest := &Payment{
		Amount:          payment.Amount,
		Currency:        payment.Currency,
		Name:            "",
		TerminalID:      s.credential.TerminalID,
		InvoiceID:       payment.InvoiceID,
		Description:     payment.Description,
		AccountID:       payment.AccountID,
		Email:           "",
		Phone:           payment.Phone,
		BackLink:        s.credential.BackLink + "/order/" + insuranceID,
		FailureBackLink: s.credential.BackLink,
		PostLink:        s.credential.PostLink,
		FailurePostLink: s.credential.PostLink,
		Language:        "RU",
		PaymentType:     "",
		PaymentJsLink:   s.credential.JSLink,
		CardSave:        cardSave,
		HomebankToken:   homebankToken,
	}

	// get token for payment
	token, err := s.getPaymentToken(paymentDest)
	if err != nil {
		return err
	}
	paymentDest.Token = token

	tmpl, err := template.ParseFiles("/app/templates/redirect.html")
	if err != nil {
		return err
	}

	err = tmpl.Execute(c, paymentDest)
	return err
}

func (s *Client) PayByCardID(cardID, insuranceID string, payment *Payment) (*Invoice, error) {
	paymentDest := &Payment{
		Amount:          payment.Amount,
		Currency:        payment.Currency,
		Name:            "",
		TerminalID:      s.credential.TerminalID,
		InvoiceID:       payment.InvoiceID,
		Description:     payment.Description,
		AccountID:       payment.AccountID,
		Email:           "",
		Phone:           payment.Phone,
		BackLink:        s.credential.BackLink + "/order/" + insuranceID,
		FailureBackLink: s.credential.BackLink,
		PostLink:        s.credential.PostLink,
		FailurePostLink: s.credential.PostLink,
		Language:        "rus",
		PaymentType:     "cardId",
	}
	paymentDest.CardID.ID = cardID

	// get token for payment
	token, err := s.getPaymentToken(paymentDest)
	if err != nil {
		return nil, err
	}

	// marshal and read invoice body
	aByte, err := json.Marshal(paymentDest)
	if err != nil {
		return nil, err
	}

	// setup request
	path := s.credential.Endpoint + "/payments/cards/auth"
	respBytes, code, err := s.handler("POST", path, aByte, nil, token)
	if err != nil {
		return nil, err
	}
	fmt.Println("invoice confirmation by card id "+paymentDest.InvoiceID+": ", string(respBytes))

	// check response code
	switch code {
	case 200:
		// unmarshal response data
		invoiceSrc := new(Invoice)
		err = json.Unmarshal(respBytes, &invoiceSrc)
		if err != nil {
			return nil, err
		}

		return invoiceSrc, nil
	default:
		return nil, errors.New(string(respBytes))
	}
}

func (s *Client) handler(method string, url string, body []byte, writer *multipart.Writer, token *Token) ([]byte, int, error) {
	// setup request
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, 0, err
	}

	// setup request header
	req.Header.Add("Content-Type", "application/json")
	if writer != nil {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}

	if token != nil {
		req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	} else {
		req.Header.Add("Authorization", "Bearer "+s.GetCredential().AccessToken)
	}

	// send request
	res, err := s.client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	// read response body
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}

	return resBody, res.StatusCode, nil
}

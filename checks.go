package lob

import (
	"strconv"
	"time"
)

// Check represents a printed check in Lob's system.
type Check struct {
	Error                *Error              `json:"error"`
	Amount               float64             `json:"amount"`
	BankAccount          *BankAccount        `json:"bank_account"`
	CheckBottom          *string             `json:"check_bottom"`
	CheckNumber          int                 `json:"check_number"`
	Data                 map[string]string   `json:"data"`
	DateCreated          string              `json:"date_created"`
	DateModified         string              `json:"date_modified"`
	Description          string              `json:"description"`
	ExpectedDeliveryDate string              `json:"expected_delivery_date"`
	SendDate             time.Time           `json:"send_date"`
	From                 *Address            `json:"from"`
	ID                   string              `json:"id"`
	Logo                 *string             `json:"logo"`
	MailType             *string             `json:"mail_type"`
	Memo                 string              `json:"memo"`
	Message              *string             `json:"message"`
	Metadata             map[string]string   `json:"metadata"`
	Name                 string              `json:"name"`
	Object               string              `json:"object"`
	Thumbnails           []map[string]string `json:"thumbnails"`
	To                   *Address            `json:"to"`
	Tracking             *Tracking           `json:"tracking"`
	URL                  string              `json:"url"`
}

// Tracking provides information on shipment tracking for a check.
type Tracking struct {
	Carrier        string        `json:"carrier"`
	Events         []interface{} `json:"events"`
	ID             string        `json:"id"`
	Link           *string       `json:"link"`
	Object         string        `json:"object"`
	TrackingNumber string        `json:"tracking_number"`
}

// Mail types that lob supports.
const (
	MailTypeUspsFirstClass = "usps_first_class"
	MailTypeUpsNextDayAir  = "ups_next_day_air"
)

// CreateCheckRequest specifies options for creating a check.
type CreateCheckRequest struct {
	Amount        float64           `json:"amount"`
	BankAccountID string            `json:"bank_account"`
	CheckBottom   *string           `json:"check_bottom"` // 400 chars, at bottom (cannot use with message)
	CheckNumber   *string           `json:"check_number"`
	Data          map[string]string `json:"data"`
	Description   *string           `json:"description"`
	FromAddressID string            `json:"from"`
	Logo          *string           `json:"logo"` // url or multiform. Square, RGB / CMYK, >= 100x100, transparent bg, PNG or JPEG, and will be grayscaled
	MailType      *string           `json:"mail_type"`
	Memo          *string           `json:"memo"`    // 40 chars in memo line
	Message       *string           `json:"message"` // 400 chars, at top (cannot use with check_bottom)
	ToAddressID   string            `json:"to"`
}

// CreateCheck requests for a new check to be printed and mailed.
func (lob *lob) CreateCheck(req *CreateCheckRequest) (*Check, error) {
	resp := new(Check)
	if err := lob.post("checks/", json2form(*req), resp); err != nil {
		return resp, err
	}
	return resp, nil
}

type CancelCheckResponse struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// GetCheck gets information about a particulr check.
func (lob *lob) GetCheck(id string) (*Check, error) {
	resp := new(Check)
	if err := lob.get("checks/"+id, nil, resp); err != nil {
		return resp, err
	}
	return resp, nil
}

func (lob *lob) CancelCheck(id string) (*CancelCheckResponse, error) {
	resp := new(CancelCheckResponse)
	if err := lob.delete("checks/"+id, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// ListChecksResponse details all of the checks we've ever mailed and printed.
type ListChecksResponse struct {
	Data        []Check `json:"data"`
	Object      string  `json:"object"`
	NextURL     string  `json:"next_url"`
	PreviousURL string  `json:"previous_url"`
	Count       int     `json:"count"`
}

// ListChecks retrieves information on all checks we've ever made, in reverse chrono order.
func (lob *lob) ListChecks(count, offset int) (*ListChecksResponse, error) {
	if count <= 0 {
		count = 10
	}
	if offset < 0 {
		offset = 0
	}
	resp := new(ListChecksResponse)
	if err := lob.get("checks", map[string]string{
		"limit":  strconv.Itoa(count),
		"offset": strconv.Itoa(offset),
	}, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

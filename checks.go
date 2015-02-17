package lob

import "strconv"

// Check represents a printed check in Lob's system.
type Check struct {
	Amount               string              `json:"amount"`
	BankAccount          *BankAccount        `json:"bank_account"`
	CheckNumber          int                 `json:"check_number"`
	DateCreated          string              `json:"date_created"`
	DateModified         string              `json:"date_modified"`
	ExpectedDeliveryDate string              `json:"expected_delivery_date"`
	ID                   string              `json:"id"`
	Memo                 string              `json:"memo"`
	Message              string              `json:"message"`
	Name                 string              `json:"name"`
	Object               string              `json:"object"`
	Price                string              `json:"price"`
	Status               string              `json:"status"`
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
	Link           string        `json:"link"`
	Object         string        `json:"object"`
	TrackingNumber string        `json:"tracking_number"`
}

// CreateCheckRequest specifies options for creating a check.
type CreateCheckRequest struct {
	Name          string `json:"name"`
	CheckNumber   string `json:"check_number"`
	BankAccountID string `json:"bank_account"`
	ToAddressID   string `json:"to"`
	Amount        string `json:"amount"`
	Message       string `json:"message"` // 400 chars, at top
	Memo          string `json:"memo"`    // 40 chars in memo line
	Logo          string `json:"logo"`    // url or multiform. Square, RGB / CMYK, >= 100x100, transparent bg, PNG or JPEG, and will be grayscaled
}

// CreateCheck requests for a new check to be printed and mailed.
func (lob *Lob) CreateCheck(req *CreateCheckRequest) (*Check, error) {
	resp := new(Check)
	return resp, Metrics.CreateCheck.Call(func() error {
		return lob.Post("checks/", json2form(*req), resp)
	})
}

// GetCheck gets information about a particulr check.
func (lob *Lob) GetCheck(id string) (*Check, error) {
	resp := new(Check)
	return resp, Metrics.GetCheck.Call(func() error {
		return lob.Get("checks/"+id, nil, resp)
	})
}

// ListChecksResponse details all of the checks we've ever mailed and printed.
type ListChecksResponse struct {
	Data        []Check `json:"data"`
	Object      string  `json:"object"`
	NextURL     string  `json:"next_url"`
	PreviousURL string  `json:"next_url"`
	Count       int     `json:"count"`
}

// ListChecks retrieves information on all checks we've ever made, in reverse chrono order.
func (lob *Lob) ListChecks(count, offset int) (*ListChecksResponse, error) {
	if count <= 0 {
		count = 10
	}
	if offset < 0 {
		offset = 0
	}
	resp := new(ListChecksResponse)
	return resp, Metrics.ListChecks.Call(func() error {
		return lob.Get("checks", map[string]string{
			"count":  strconv.Itoa(count),
			"offset": strconv.Itoa(offset),
		}, resp)
	})
}

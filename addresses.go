package lob

import "strconv"

// Address represents an address stored in the Lob's system.
type Address struct {
	AddressCity    string `json:"address_city"`
	AddressCountry string `json:"address_country"`
	AddressLine1   string `json:"address_line1"`
	AddressLine2   string `json:"address_line2"`
	AddressState   string `json:"address_state"`
	AddressZip     string `json:"address_zip"`
	DateCreated    string `json:"date_created"`
	DateModified   string `json:"date_modified"`
	Email          string `json:"email"`
	ID             string `json:"id"`
	Name           string `json:"name"`
	Object         string `json:"object"`
	Phone          string `json:"phone"`
	Deleted        int    `json:"deleted"`
}

// CreateAddress creates an address in Lob's system.
func (lob *Lob) CreateAddress(address *Address) (*Address, error) {
	resp := new(Address)
	return resp, Metrics.CreateAddress.Call(func() error {
		return lob.Post("addresses", json2form(*address), resp)
	})
}

// GetAddress retrieves an address with the given id.
func (lob *Lob) GetAddress(id string) (*Address, error) {
	resp := new(Address)
	return resp, Metrics.GetAddress.Call(func() error {
		return lob.Get("addresses/"+id, nil, resp)
	})
}

type message struct {
	Message string `json:"message"`
}

// DeleteAddress deletes the given address from Lob's system.
func (lob *Lob) DeleteAddress(id string) (string, error) {
	resp := new(message)
	err := Metrics.DeleteAddress.Call(func() error {
		return lob.Delete("addresses/"+id, resp)
	})
	return resp.Message, err
}

// ListAddressesResponse gives the results for listing all addresses for our account.
type ListAddressesResponse struct {
	Data        []Address `json:"data"`
	Object      string    `json:"object"`
	NextURL     string    `json:"next_url"`
	PreviousURL string    `json:"next_url"`
	Count       int       `json:"count"`
}

// ListAddresses lists all addresses on this account, paginated.
func (lob *Lob) ListAddresses(count int, offset int) (*ListAddressesResponse, error) {
	if count <= 0 {
		count = 10
	}
	if offset < 0 {
		offset = 0
	}

	resp := new(ListAddressesResponse)
	return resp, Metrics.ListAddresses.Call(func() error {
		return lob.Get("addresses/", map[string]string{
			"count":  strconv.Itoa(count),
			"offset": strconv.Itoa(offset),
		}, resp)
	})
}

// AddressVerificationRequest validates the given subset of info from an address.
type AddressVerificationRequest struct {
	AddressCity    string `json:"address_city"`
	AddressCountry string `json:"address_country"`
	AddressLine1   string `json:"address_line1"`
	AddressLine2   string `json:"address_line2"`
	AddressState   string `json:"address_state"`
	AddressZip     string `json:"address_zip"`
}

// AddressVerificationResponse gives the response from attempting to verify an address.
type AddressVerificationResponse struct {
	Address Address        `json:"address"`
	Errors  []ErrorMessage `json:"errors"`
}

// ErrorMessage gives information about a failure.
type ErrorMessage struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

// VerifyAddress verifies the given address and cleans it up.
func (lob *Lob) VerifyAddress(address *Address) (*AddressVerificationResponse, error) {
	req := AddressVerificationRequest{
		AddressCity:    address.AddressCity,
		AddressCountry: address.AddressCountry,
		AddressLine1:   address.AddressLine1,
		AddressLine2:   address.AddressLine2,
		AddressState:   address.AddressState,
		AddressZip:     address.AddressZip,
	}
	resp := new(AddressVerificationResponse)
	return resp, Metrics.VerifyAddress.Call(func() error {
		return lob.Post("verify", json2form(req), resp)
	})
}

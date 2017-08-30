package lob

import (
	"errors"
	"strconv"
)

// Address represents an address stored in the Lob's system.
type Address struct {
	AddressCity    *string           `json:"address_city"`
	AddressCountry *string           `json:"address_country"`
	AddressLine1   string            `json:"address_line1"`
	AddressLine2   *string           `json:"address_line2"`
	AddressState   *string           `json:"address_state"`
	AddressZip     *string           `json:"address_zip"`
	Company        *string           `json:"company"`
	DateCreated    string            `json:"date_created"`
	DateModified   string            `json:"date_modified"`
	Deleted        *bool             `json:"deleted"`
	Description    *string           `json:"description"`
	Email          *string           `json:"email"`
	ID             string            `json:"id"`
	Metadata       map[string]string `json:"metadata"`
	Name           *string           `json:"name"`
	Object         string            `json:"object"`
	Phone          *string           `json:"phone"`
}

// CreateAddress creates an address in Lob's system.
func (lob *lob) CreateAddress(address *Address) (*Address, error) {
	resp := new(Address)
	if err := lob.post("addresses", json2form(*address), resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetAddress retrieves an address with the given id.
func (lob *lob) GetAddress(id string) (*Address, error) {
	resp := new(Address)
	if err := lob.get("addresses/"+id, nil, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

type deleteAddressResp struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// DeleteAddress deletes the given address from Lob's system.
func (lob *lob) DeleteAddress(id string) error {
	resp := new(deleteAddressResp)

	if err := lob.delete("addresses/"+id, resp); err != nil {
		return err
	}
	if !resp.Deleted {
		return errors.New("failed to delete address")
	}
	return nil
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
func (lob *lob) ListAddresses(count int, offset int) (*ListAddressesResponse, error) {
	if count <= 0 {
		count = 10
	}
	if offset < 0 {
		offset = 0
	}

	resp := new(ListAddressesResponse)
	if err := lob.get("addresses/", map[string]string{
		"limit":  strconv.Itoa(count),
		"offset": strconv.Itoa(offset),
	}, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// AddressVerificationRequest validates the given subset of info from an address.
type AddressVerificationRequest struct {
	AddressCity    *string `json:"address_city"`
	AddressCountry *string `json:"address_country"`
	AddressLine1   *string `json:"address_line1"`
	AddressLine2   *string `json:"address_line2"`
	AddressState   *string `json:"address_state"`
	AddressZip     *string `json:"address_zip"`
	Name           *string `json:"name"`
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
func (lob *lob) VerifyAddress(address *Address) (*AddressVerificationResponse, error) {
	req := AddressVerificationRequest{
		AddressCity:    address.AddressCity,
		AddressCountry: address.AddressCountry,
		AddressLine1:   &address.AddressLine1,
		AddressLine2:   address.AddressLine2,
		AddressState:   address.AddressState,
		AddressZip:     address.AddressZip,
	}
	resp := new(AddressVerificationResponse)
	if err := lob.post("verify", json2form(req), resp); err != nil {
		return nil, err
	}
	return resp, nil
}

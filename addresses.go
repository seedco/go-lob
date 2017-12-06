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
type USAddressVerificationRequest struct {
	Recipient    *string `json:"recipient"`
	AddressLine1 *string `json:"primary_line"`
	AddressLine2 *string `json:"secondary_line"`
	AddressCity  *string `json:"city"`
	AddressState *string `json:"state"`
	AddressZip   *string `json:"zip_code"`
}

// USAddressVerificationResponse gives the response from attempting to verify a US address.
type USAddressVerificationResponse struct {
	Id                     string                          `json:"id"`
	Recipient              string                          `json:"recipient"`
	PrimaryLine            string                          `json:"primary_line"`
	SecondaryLine          string                          `json:"secondary_line"`
	Urbanization           string                          `json:"urbanization,omitempty"`
	LastLine               string                          `json:"last_line"`
	Deliverability         string                          `json:"deliverability"`
	Components             USAddressComponents             `json:"components"`
	DeliverabilityAnalysis USAddressDeliverabilityAnalysis `json:"deliverability_analysis"`
	Object                 string                          `json:"object"`
}

type USAddressComponents struct {
	PrimaryNumber             string  `json:"primary_number"`
	StreetPredirection        string  `json:"street_predirection"`
	StreetName                string  `json:"street_name"`
	StreetSuffix              string  `json:"street_suffix"`
	StreetPostdirection       string  `json:"street_postdirection"`
	SecondaryDesignator       string  `json:"secondary_designator"`
	SecondaryNumber           string  `json:"secondary_number"`
	PmbDesignator             string  `json:"pmb_designator"`
	PmbNumber                 string  `json:"pmb_number"`
	ExtraSecondary_designator string  `json:"extra_secondary_designator"`
	ExtraSecondary_number     string  `json:"extra_secondary_number"`
	City                      string  `json:"city"`
	State                     string  `json:"state"`
	ZipCode                   string  `json:"zip_code"`
	ZipCodePlus_4             string  `json:"zip_code_plus_4"`
	ZipCodeType               string  `json:"zip_code_type"`
	DeliveryPoint_barcode     string  `json:"delivery_point_barcode"`
	AddressType               string  `json:"address_type"`
	RecordType                string  `json:"record_type"`
	DefaultBuildingAddress    bool    `json:"default_building_address"`
	County                    string  `json:"county"`
	CountyFips                string  `json:"county_fips"`
	CarrierRoute              string  `json:"carrier_route"`
	CarrierRouteType          string  `json:"carrier_route_type"`
	Latitude                  float64 `json:"latitude"`
	Longitude                 float64 `json:"longitude"`
}

type USAddressDeliverabilityAnalysis struct {
	DpvConfirmation string   `json:"dpv_confirmation"`
	DpvCmra         string   `json:"dpv_cmra"`
	DpvVacant       string   `json:"dpv_vacant"`
	DpvFootnotes    []string `json:"dpv_footnotes"`
	EwsMatch        bool     `json:"ews_match"`
	LacsIndicator   string   `json:"lacs_indicator"`
	LacsReturnCode  string   `json:"lacs_return_code"`
	SuiteReturnCode string   `json:"suite_return_code"`
}

// VerifyUSAddress verifies the given US address and returns the validation results.
func (lob *lob) VerifyUSAddress(address *Address) (*USAddressVerificationResponse, error) {
	req := USAddressVerificationRequest{
		Recipient:    address.Name,
		AddressLine1: &address.AddressLine1,
		AddressLine2: address.AddressLine2,
		AddressCity:  address.AddressCity,
		AddressState: address.AddressState,
		AddressZip:   address.AddressZip,
	}
	resp := new(USAddressVerificationResponse)
	if err := lob.post("us_verifications", json2form(req), resp); err != nil {
		return nil, err
	}
	return resp, nil
}

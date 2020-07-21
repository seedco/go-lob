package lob

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Error struct {
	Message    string `json:"message,omitempty"`
	StatusCode int    `json:"status_code,omitempty"`
}

// Address represents an address stored in the Lob's system.
type Address struct {
	Error          *Error            `json:"error"`
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
		return resp, err
	}
	return resp, nil
}

// GetAddress retrieves an address with the given id.
func (lob *lob) GetAddress(id string) (*Address, error) {
	resp := new(Address)
	if err := lob.get("addresses/"+id, nil, resp); err != nil {
		return resp, err
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
	PreviousURL string    `json:"previous_url"`
	Count       int       `json:"count"`
}

// ListAddresses lists all addresses on this account, paginated.
func (lob *lob) ListAddresses(count int) (*ListAddressesResponse, error) {
	if count <= 0 {
		count = 10
	}

	resp := new(ListAddressesResponse)
	if err := lob.get("addresses/", map[string]string{
		"limit": strconv.Itoa(count),
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
	Id                     string                              `json:"id"`
	Recipient              string                              `json:"recipient"`
	PrimaryLine            string                              `json:"primary_line"`
	SecondaryLine          string                              `json:"secondary_line"`
	Urbanization           string                              `json:"urbanization,omitempty"`
	LastLine               string                              `json:"last_line"`
	Deliverability         USAddressVerificationDeliverability `json:"deliverability"`
	Components             USAddressComponents                 `json:"components"`
	DeliverabilityAnalysis USAddressDeliverabilityAnalysis     `json:"deliverability_analysis"`
	Object                 string                              `json:"object"`
}

//USAddressVerificationDeliverability is the type for the deliverability of an address verified
type USAddressVerificationDeliverability string

//list of USAddressVerificationDeliverability values
var (
	USAddressVerificationDeliverabilityDeliverable     USAddressVerificationDeliverability = "deliverable"
	USAddressVerificationDeliverabilityUnnecessaryUnit USAddressVerificationDeliverability = "deliverable_unnecessary_unit"
	USAddressVerificationDeliverabilityIncorrectUnit   USAddressVerificationDeliverability = "deliverable_incorrect_unit"
	USAddressVerificationDeliverabilitydMissingUnit    USAddressVerificationDeliverability = "deliverable_missing_unit"
	USAddressVerificationDeliverabilityUndeliverable   USAddressVerificationDeliverability = "undeliverable"
)

//USAddressComponents are the components which make up a US address
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

// with special codes you can get a proper response back in test mode; this magic secondary line means we didn't
// request it with a special code so fallback to the old behavior
const testFillInLine2Required = "See Https://www.lob.com/docs#us-verification-test-environment For More Info"

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

	// in test, fill in components
	if strings.HasPrefix(lob.APIKey, "test") && resp.SecondaryLine == testFillInLine2Required {
		streetSplit := strings.Split(address.AddressLine1, " ")
		if len(streetSplit) > 2 {
			resp.Components.PrimaryNumber = streetSplit[0]
			resp.Components.StreetName = streetSplit[1]
			resp.Components.StreetSuffix = streetSplit[2]
		}
		if address.AddressLine2 != nil {
			resp.Components.SecondaryNumber = *address.AddressLine2
		}
		if address.AddressZip != nil {
			resp.Components.ZipCode = *address.AddressZip
		}
		if address.AddressCity != nil {
			resp.Components.City = *address.AddressCity
		}
		if address.AddressState != nil {
			resp.Components.State = *address.AddressState
		}
	}

	return resp, nil
}

//AddressVerificationRequestCasing states how the verified address should be returned
type AddressVerificationRequestCasing string

//possible values of AddressVerificationRequestCasing
var (
	AddressVerificationRequestCasingUpper  AddressVerificationRequestCasing = "upper"
	AddressVerificationRequestCasingProper AddressVerificationRequestCasing = "proper"
)

// VerifyUSAddress verifies the given US address and returns the validation results.
func (lob *lob) VerifyUSAddressWithCasing(address *Address, casing AddressVerificationRequestCasing) (*USAddressVerificationResponse, error) {
	req := USAddressVerificationRequest{
		Recipient:    address.Name,
		AddressLine1: String(address.AddressLine1),
		AddressLine2: address.AddressLine2,
		AddressCity:  address.AddressCity,
		AddressState: address.AddressState,
		AddressZip:   address.AddressZip,
	}

	resp := new(USAddressVerificationResponse)
	if err := lob.post(fmt.Sprintf("us_verifications?case=%s", casing), json2form(req), resp); err != nil {
		return nil, err
	}

	// in test, fill in components
	if strings.HasPrefix(lob.APIKey, "test") && resp.SecondaryLine == testFillInLine2Required {
		streetSplit := strings.Split(address.AddressLine1, " ")
		if len(streetSplit) > 2 {
			resp.Components.PrimaryNumber = streetSplit[0]
			resp.Components.StreetName = streetSplit[1]
			resp.Components.StreetSuffix = streetSplit[2]
		}
		if address.AddressLine2 != nil {
			resp.Components.SecondaryNumber = *address.AddressLine2
		}
		if address.AddressZip != nil {
			resp.Components.ZipCode = *address.AddressZip
		}
		if address.AddressCity != nil {
			resp.Components.City = *address.AddressCity
		}
		if address.AddressState != nil {
			resp.Components.State = *address.AddressState
		}
	}

	return resp, nil
}

package lob

import (
	"strconv"
	"time"
)

//#region model definition

// Letter represents a letter in lob's system
type Letter struct {
	Error                *Error                 `json:"error"`
	ID                   string                 `json:"id"`
	Description          *string                `json:"description"`
	Metadata             map[string]string      `json:"metadata"`
	To                   *Address               `json:"to"`
	From                 *Address               `json:"from"`
	Color                bool                   `json:"color"`
	DoubleSided          bool                   `json:"double_sided"`
	AddressPlacement     LetterAddressPlacement `json:"address_placement"`
	ReturnEnvelope       bool                   `json:"return_envelope"`
	PerforatedPage       *uint64                `json:"perforated_page"`
	CustomEnvelope       *CustomEnvelope        `json:"custom_envelope"`
	ExtraService         *LetterExtraService    `json:"extra_service"`
	MailType             *string                `json:"mail_type"` //MailTypeUspsFirstClass or MailTypeUspsStandard
	URL                  string                 `json:"url"`
	MergeVariables       map[string]string      `json:"merge_variables"`
	TemplateID           *string                `json:"template_id"`
	TemplateVersionID    *string                `json:"template_version_id"`
	Carrier              string                 `json:"carrier"` //value is USPS
	TrackingNumber       *string                `json:"tracking_number"`
	TrackingEvents       []TrackingEvent        `json:"tracking_events"`
	Thumbnails           []LetterThumbnail      `json:"thumbnails"`
	ExpectedDeliveryDate string                 `json:"expected_delivery_date"`
	DateCreated          time.Time              `json:"date_created"`
	DateModified         time.Time              `json:"date_modified"`
	SendDate             time.Time              `json:"send_date"`
	Deleted              bool                   `json:"deleted"`
	Object               string                 `json:"object"` //value will always be letter
}

//CreateLetterRequest is the object for creating a new letter
type CreateLetterRequest struct {
	Description      *string                `json:"description"` //must be no longer than 255 characters
	To               Address                `json:"to"`          //if you need the id, simply do lob.Address{ID: "id-here"}
	From             Address                `json:"from"`        //if you need the id, simply do lob.Address{ID: "id-here"}
	BillingGroupID   *string                `json:"billing_group_id"`
	SendDate         *time.Time             `json:"send_date"`
	Color            bool                   `json:"color"`
	File             string                 `json:"file"` //please see lob's create letter documentation: https://lob.com/docs#letters_create
	DoubleSided      bool                   `json:"double_sided"`
	AddressPlacement LetterAddressPlacement `json:"address_placement"`
	MailType         *string                `json:"mail_type"` //MailTypeUspsFirstClass or MailTypeUspsStandard
	ExtraService     *LetterExtraService    `json:"extra_service"`
	ReturnEnvelope   bool                   `json:"return_envelope"`
	PerforatedPage   *uint64                `json:"perforated_page"`
	CustomEnvelope   *CustomEnvelope        `json:"custom_envelope"`
	MergeVariables   map[string]string      `json:"merge_variables"`
	Metadata         map[string]string      `json:"metadata"`
}

//LetterAddressPlacement represents where the address should be placed on the letter
type LetterAddressPlacement string

var (
	//LetterAddressPlacementTopFirstPage is for being printed at the top of the first page
	LetterAddressPlacementTopFirstPage LetterAddressPlacement = "top_first_page"
	//LetterAddressPlacementInsertBlankPage is for inserting a blank page (does cost extra)
	LetterAddressPlacementInsertBlankPage LetterAddressPlacement = "insert_blank_page"
)

// CustomEnvelope represents a custom envelope in lob's system
type CustomEnvelope struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Object string `json:"object"` //value should always be "envelope"
}

//LetterExtraService represents the type of extra service requested
type LetterExtraService string

var (
	//LetterExtraServiceCertified is for certified mail
	LetterExtraServiceCertified LetterExtraService = "certified"
	//LetterExtraServiceCertifiedReturnReceipt is for certified mail with return receipt
	LetterExtraServiceCertifiedReturnReceipt LetterExtraService = "certified_return_receipt"
	//LetterExtraServiceRegistered is for registered mail
	LetterExtraServiceRegistered LetterExtraService = "registered"
)

//LetterThumbnail represents the thumbnails of the generated letter
type LetterThumbnail struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

//ListLettersResponse details all of the letters we've ever mailed and printed
type ListLettersResponse struct {
	Data        []Letter `json:"data"`
	Object      string   `json:"object"` //value will always be letter
	NextURL     string   `json:"next_url"`
	PreviousURL string   `json:"previous_url"`
	Count       int      `json:"count"`
}

//CancelLetterResponse is the response returned when deleting a letter
type CancelLetterResponse struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

//#endregion model definition

// CreateLetter requests for a new letter to be printed and mailed.
func (lob *lob) CreateLetter(req CreateLetterRequest) (*Letter, error) {
	resp := new(Letter)

	if err := lob.post("letters", json2form(req), resp); err != nil {
		return resp, err
	}

	return resp, nil
}

// GetLetter gets information about a particulr letter.
func (lob *lob) GetLetter(id string) (*Letter, error) {
	resp := new(Letter)

	if err := lob.get("letters/"+id, nil, resp); err != nil {
		return resp, err
	}

	return resp, nil
}

// ListLetters retrieves information on all letters we've ever made, in reverse chrono order
func (lob *lob) ListLetters(count int) (*ListLettersResponse, error) {
	if count <= 0 {
		count = 10
	}

	query := map[string]string{
		"limit": strconv.Itoa(count),
	}

	resp := new(ListLettersResponse)
	if err := lob.get("letters", query, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// CancelLetter prevents a letter from being sent, if the send date has yet to pass
func (lob *lob) CancelLetter(id string) (*CancelLetterResponse, error) {
	resp := new(CancelLetterResponse)
	if err := lob.delete("letters/"+id, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

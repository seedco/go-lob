package lob

import (
	"time"
)

// Mail types that lob supports.
const (
	MailTypeUspsStandard   = "usps_standard"
	MailTypeUspsFirstClass = "usps_first_class"
	MailTypeUpsNextDayAir  = "ups_next_day_air"
)

// Tracking provides information on shipment tracking for a check.
type Tracking struct {
	Carrier        string          `json:"carrier"`
	Events         []TrackingEvent `json:"events"`
	ID             string          `json:"id"`
	Link           *string         `json:"link"`
	Object         string          `json:"object"`
	TrackingNumber string          `json:"tracking_number"`
}

//TrackingEvent represents an event from the tracking
type TrackingEvent struct {
	ID           string                `json:"id"`
	Type         TrackingEventType     `json:"type"`
	Name         TrackingEventName     `json:"name"`
	Details      *TrackingEventDetails `json:"details"`
	Location     *string               `json:"location"`
	Time         time.Time             `json:"time"`
	DateCreated  time.Time             `json:"date_created"`
	DateModified time.Time             `json:"date_modified"`
	Object       string                `json:"object"` //value will be tracking_event
}

//TrackingEventType indicates the type of event
type TrackingEventType string

var (
	//TrackingEventTypeCertified is for certified mail; has the details property on the event
	TrackingEventTypeCertified TrackingEventType = "certified"
	//TrackingEventTypeNormal is for normal non-certified postcards, letters, and checks
	TrackingEventTypeNormal TrackingEventType = "normal"
)

//TrackingEventName is a special type for event names
type TrackingEventName string

//list of event names
var (
	TrackingEventNameMailed               TrackingEventName = "Mailed"
	TrackingEventNameInTransit            TrackingEventName = "In Transit"
	TrackingEventNameInLocalArea          TrackingEventName = "In Local Area"
	TrackingEventNameProcessedForDelivery TrackingEventName = "Processed for Delivery"
	TrackingEventNameReRouted             TrackingEventName = "Re-Routed"
	TrackingEventNameReturnedToSender     TrackingEventName = "Returned to Sender"
	TrackingEventNamePickupAvailable      TrackingEventName = "Pickup Available"
	TrackingEventNameDelivered            TrackingEventName = "Delivered"
	TrackingEventNameIssue                TrackingEventName = "Issue"
)

//TrackingEventDetails is for the details of a tracking event, only happens for certified mail
type TrackingEventDetails struct {
	Event          string `json:"event"`
	Description    string `json:"description"`
	Notes          string `json:"notes"`
	ActionRequired bool   `json:"action_required"`
}

package callback

import (
	"time"
)

type XenditPayload struct {
	Event      string             `json:"event,omitempty"`
	BusinessID string             `json:"business_id,omitempty"`
	Created    *time.Time         `json:"created,omitempty"`
	Data       *xenditPayloadData `json:"data,omitempty"`
}

type xenditPayloadData struct {
	ID               string                          `json:"id,omitempty"`
	Amount           int64                           `json:"amount,omitempty"`
	Country          string                          `json:"country,omitempty"`
	Currency         string                          `json:"currency,omitempty"`
	PaymentRequestID string                          `json:"payment_request_id,omitempty"`
	ReferenceID      string                          `json:"reference_id,omitempty"`
	Status           string                          `json:"status,omitempty"`
	CustomedID       string                          `json:"customed_id,omitempty"`
	Description      string                          `json:"description,omitempty"`
	PaymehtMethod    *xenditPayloadDataPaymentMethod `json:"paymeht_method,omitempty"`
	Items            []*xenditPayloadDataItem        `json:"items,omitempty"`
	Metadata         map[string]any                  `json:"metadata,omitempty"`
	Created          *time.Time                      `json:"created,omitempty"`
	Updated          *time.Time                      `json:"updated,omitempty"`
}

type xenditPayloadDataPaymentMethod struct {
	ID          string                                `json:"id,omitempty"`
	Type        string                                `json:"type,omitempty"`
	Reusability string                                `json:"reusability,omitempty"`
	Status      string                                `json:"status,omitempty"`
	Created     *time.Time                            `json:"created,omitempty"`
	QrCode      *xenditPayloadDataPaymentMethodQrCode `json:"qr_code,omitempty"`
}

type xenditPayloadDataPaymentMethodQrCode struct {
	Amount            int64                                                  `json:"amount,omitempty"`
	Currency          string                                                 `json:"currency,omitempty"`
	ChannelCode       string                                                 `json:"channel_code,omitempty"`
	ChannelProperties *xenditPayloadDataPaymentMethodQrCodeChannelProperties `json:"channel_properties,omitempty"`
}

type xenditPayloadDataPaymentMethodQrCodeChannelProperties struct {
	QrString string `json:"qr_string,omitempty"`
}

type xenditPayloadDataItem struct {
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
	Price       int64  `json:"price,omitempty"`
	Category    string `json:"category,omitempty"`
	Currency    string `json:"currency,omitempty"`
	Quantity    int    `json:"quantity,omitempty"`
	ReferenceID string `json:"reference_id,omitempty"`
}

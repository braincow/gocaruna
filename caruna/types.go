package caruna

import (
	"encoding/json"
	"net/url"
	"time"
)

type ParsedURL url.URL

func (m *ParsedURL) UnmarshalJSON(data []byte) error {
	// Unquote the data to get the raw string
	var u string
	err := json.Unmarshal(data, &u)
	if err != nil {
		return err
	}

	// Parse the URL
	parsedURL, err := url.Parse(u)
	if err != nil {
		return err
	}

	*m = ParsedURL(*parsedURL)
	return nil
}

type UserInfo struct {
	UserName                   string    `json:"userName"`
	UserType                   string    `json:"userType"`
	Email                      string    `json:"email"`
	FirstName                  string    `json:"FirstName"`
	LastName                   string    `json:"LastName"`
	PhoneNumber                string    `json:"PhoneNumber"`
	ProfileURL                 ParsedURL `json:"iamProfileUrl"`
	OwnCustomerNumbers         []string  `json:"ownCustomerNumbers"`
	RepresentedCustomerNumbers []string  `json:"representedCustomerNumbers"`
	HashedUserID               string    `json:"HashedUserId"`
	VisitorParams              string    `json:"giosgVisitorParams"`
}

type LoginInfo struct {
	Token                 string    `json:"token"`
	ExpiresAt             int       `json:"expiresAt"`
	User                  UserInfo  `json:"user"`
	RedirectAfterLoginURL ParsedURL `json:"redirectAfterLogin"`
}

type PostalAddress struct {
	ID                 string `json:"id"`
	StreetName         string `json:"streetName"`
	HouseNumber        string `json:"houseNumber"`
	PostOffice         string `json:"postOffice"`
	PostalCode         string `json:"postalCode"`
	LocalityType       string `json:"postalAddress"`
	InvoicingBaseCount int    `json:"invoicingBaseCount"`
}

type MeteringPointCoordinate struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type MeteringPoint struct {
	Type               string                  `json:"type"`
	ID                 string                  `json:"id"`
	Address            PostalAddress           `json:"address"`
	AssetID            string                  `json:"assetId"`
	ContractID         string                  `json:"contractId"`
	ContractType       string                  `json:"energyJunctionContract"`
	CustomerAddressId  string                  `json:"customerAddressId"`
	IsSupplierInvoiced bool                    `json:"isSupplierInvoiced"`
	FuseSize           string                  `json:"fuseSize"`
	Group              string                  `json:"group"`
	NetworkID          string                  `json:"networkId"`
	Position           MeteringPointCoordinate `json:"position"`
	Tabs               []string                `json:"tabs"`
	Use                string                  `json:"use"`
	CustomerID         string                  `json:"customerId"`
}

type ConsumptionPart struct {
	NightTime float64 `json:"nighttime"`
	DayTime   float64 `json:"daytime"`
}

type ConsumerCost struct {
	Timestamp                                 time.Time       `json:"timestamp"`
	TotalConsumption                          float64         `json:"totalConsumption"`
	InvoicedConsumption                       float64         `json:"invoicedConsumption"`
	TotalFee                                  float64         `json:"totalFee"`
	DistributionFee                           float64         `json:"distributionFee"`
	DistributionFeeBase                       float64         `json:"distributionFeeBase"`
	ElectricityTax                            float64         `json:"electricityTax"`
	ValueAddedTax                             float64         `json:"valueAddedTax"`
	Temperature                               float64         `json:"temperature"`
	InvoicedConsumptionByTransferProductParts ConsumptionPart `json:"invoicedConsumptionByTransferProductParts"`
	DistributionFeeByTransferProductParts     ConsumptionPart `json:"distributionFeeByTransferProductParts"`
}

type CustomerMarketingPermissions struct {
	ViaEmail          bool `json:"email"`
	ViaSMS            bool `json:"sms"`
	MarketingIsBanned bool `json:"ban"`
}

type CustomerContactingPermissions struct {
	ViaEmail bool `json:"email"`
}

type ElectronicInvoiceAddress struct {
	ID                 string `json:"id"`
	AddressTypeKey     string `json:"addressTypeKey"`
	Address            string `json:"address"`
	Operator1          string `json:"operator1"`
	Operator2          string `json:"operator2"`
	OvtCode            string `json:"ovtCode"`
	PaymentInstruction string `json:"paymentInstruction"`
	BuyerServiceCode   int    `json:"buyerSericeCode"`
	State              int    `json:"state"`
}

type CustomerInfo struct {
	ID                         string                        `json:"id"`
	Name                       string                        `json:"name"`
	Email                      string                        `json:"email"`
	BusinessID                 string                        `json:"businessId"`
	Phone                      string                        `json:"phone"`
	PostalAddress              PostalAddress                 `json:"postalAddress"`
	BillingAddresses           []PostalAddress               `json:"billingAddresses"`
	MarketingPermissions       CustomerMarketingPermissions  `json:"marketingPermissions"`
	ContactingPermissions      CustomerContactingPermissions `json:"contactingPermissions"`
	ElectronicInvoiceAddresses []ElectronicInvoiceAddress    `json:"eInvoiceAddresses"`
	RepresentedCustomer        bool                          `json:"representedCustomer"`
	JointCustomer              bool                          `json:"jointCustomer"`
	IsMissingRequiredInfo      bool                          `json:"isMissingRequiredInfo"`
	Language                   string                        `json:"language"`
	DefaultRefundAccount       string                        `json:"defaultRefundAccount"`
	DefaultRefundAccountBIC    string                        `json:"defaultRefundAccountBIC"`
	UpdateCampaignSeen         bool                          `json:"updateCampaignSeen"`
}

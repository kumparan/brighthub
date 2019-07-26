package brighthub

type (
	// Action :nodoc:
	Action string
	// EntityType :nodoc:
	EntityType string
	// Status :nodoc:
	Status string
)

const (
	// ActionCreate :nodoc:
	ActionCreate Action = "CREATE"
	// ActionPublish :nodoc:
	ActionPublish Action = "PUBLISH"

	// AssetEntityType :nodoc:
	AssetEntityType EntityType = "ASSET"
	// DigitalMasterEntityType :nodoc:
	DigitalMasterEntityType EntityType = "DIGITAL_MASTER"
	// DynamicRenditionEntityType :nodoc:
	DynamicRenditionEntityType EntityType = "DYNAMIC_RENDITION"
	// TitleEntityType :nodoc:
	TitleEntityType EntityType = "TITLE"

	// StatusFailed :nodoc:
	StatusFailed Status = "FAILED"
	// StatusSuccess :nodoc:
	StatusSuccess Status = "SUCCESS"
)

// Notification :nodoc:
type Notification struct {
	Entity             string     `json:"entity"`
	EntityType         EntityType `json:"entityType"`
	Version            string     `json:"version"`
	Action             Action     `json:"action"`
	JobID              string     `json:"jobId"`
	VideoID            string     `json:"videoId"`
	DynamicRenditionID string     `json:"dynamicRenditionId"`
	Language           string     `json:"language"`
	Variant            string     `json:"variant"`
	AccountID          string     `json:"accountId"`
	Status             Status     `json:"status"`
	ErrorMessage       string     `json:"errorMessage"`
}

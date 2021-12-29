package aidbox

type ResourceBase struct {
	ID string `json:"id,omitempty"`
	//Meta *ResourceBaseMeta `json:"meta,omitempty"`
}

//type ResourceBaseMeta struct {
//	CreatedAt   *time.Time `json:"createdAt,omitempty"`
//	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
//	VersionId   string     `json:"versionId,omitempty"`
//}

type Resource interface{}

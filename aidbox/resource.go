package aidbox

import "time"

type ResourceBase struct {
	ID   string            `json:"id,omitempty"`
	Meta *ResourceBaseMeta `json:"meta,omitempty"`
}

func (a *ResourceBase) GetID() string {
	return a.ID
}

type ResourceBaseMeta struct {
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	VersionId   string     `json:"versionId,omitempty"`
	Profile     []string   `json:"profile,omitempty"`
}

type Resource interface {
	GetResourcePath() string
	GetID() string
}

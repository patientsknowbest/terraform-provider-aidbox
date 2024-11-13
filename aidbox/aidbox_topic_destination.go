package aidbox

import (
	"context"
)

type AidboxTopicDestination struct {
	ResourceBase
	Topic     string                  `json:"topic"`
	Kind      string                  `json:"kind"`
	Parameter []SubscriptionParameter `json:"parameter"`
}

type SubscriptionParameter struct {
	Name        string `json:"name"`
	Url         string `json:"valueUrl,omitempty"`
	UnsignedInt int    `json:"valueUnsignedInt,omitempty"`
	String      string `json:"valueString,omitempty"`
}

func (*AidboxTopicDestination) GetResourcePath() string {
	return "fhir/AidboxTopicDestination"
}

func (apiClient *ApiClient) CreateAidboxTopicDestination(ctx context.Context, searchParameter *AidboxTopicDestination) (*AidboxTopicDestination, error) {
	c, err := apiClient.createResource(ctx, searchParameter)
	if err != nil {
		return nil, err
	}
	return c.(*AidboxTopicDestination), nil
}

func (apiClient *ApiClient) GetAidboxTopicDestination(ctx context.Context, id string) (*AidboxTopicDestination, error) {
	// TODO (AS) wtf is this, why hardcoded here but otherwise using the method for the path elsewhere??? fix this
	rr, err := apiClient.getResource(ctx, "/fhir/AidboxTopicDestination/"+id)
	if err != nil {
		return nil, err
	}
	return rr.(*AidboxTopicDestination), nil
}

func (apiClient *ApiClient) UpdateAidboxTopicDestination(ctx context.Context, q *AidboxTopicDestination) (*AidboxTopicDestination, error) {
	rr, err := apiClient.updateResource(ctx, q)
	if err != nil {
		return nil, err
	}
	return rr.(*AidboxTopicDestination), nil
}

func (apiClient *ApiClient) DeleteAidboxTopicDestination(ctx context.Context, id string) error {
	return apiClient.deleteResource(ctx, "/fhir/AidboxTopicDestination/"+id)
}

package aidbox

import (
	"context"
)

type AidboxTopicDestination struct {
	ResourceBase
	Topic     string                  `json:"topic"`
	Kind      string                  `json:"kind"`
	Content   string                  `json:"content"`
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

func (apiClient *ApiClient) CreateAidboxTopicDestination(ctx context.Context, aidboxTopicDestination *AidboxTopicDestination) (*AidboxTopicDestination, error) {
	response := &AidboxTopicDestination{}
	return response, apiClient.createResource(ctx, aidboxTopicDestination, response)
}

func (apiClient *ApiClient) GetAidboxTopicDestination(ctx context.Context, id string) (*AidboxTopicDestination, error) {
	response := &AidboxTopicDestination{}
	return response, apiClient.getResource(ctx, id, response)
}

func (apiClient *ApiClient) UpdateAidboxTopicDestination(ctx context.Context, q *AidboxTopicDestination) (*AidboxTopicDestination, error) {
	response := &AidboxTopicDestination{}
	return response, apiClient.updateResource(ctx, q, response)
}

func (apiClient *ApiClient) DeleteAidboxTopicDestination(ctx context.Context, id string) error {
	return apiClient.deleteResource(ctx, id, &AidboxTopicDestination{})
}

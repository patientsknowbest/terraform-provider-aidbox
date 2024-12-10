package aidbox

import (
	"context"
)

type AidboxSubscriptionTopic struct {
	ResourceBase
	Url     string         `json:"url"`
	Status  string         `json:"status"`
	Trigger []TopicTrigger `json:"trigger"`
}

type TopicTrigger struct {
	Resource string `json:"resource"`
}

func (*AidboxSubscriptionTopic) GetResourcePath() string {
	return "fhir/AidboxSubscriptionTopic"
}

func (apiClient *ApiClient) CreateAidboxSubscriptionTopic(ctx context.Context, aidboxSubscriptionTopic *AidboxSubscriptionTopic) (*AidboxSubscriptionTopic, error) {
	response := &AidboxSubscriptionTopic{}
	return response, apiClient.createResource(ctx, aidboxSubscriptionTopic, response)
}

func (apiClient *ApiClient) GetAidboxSubscriptionTopic(ctx context.Context, id string) (*AidboxSubscriptionTopic, error) {
	response := &AidboxSubscriptionTopic{}
	return response, apiClient.getResource(ctx, id, response)
}

func (apiClient *ApiClient) UpdateAidboxSubscriptionTopic(ctx context.Context, q *AidboxSubscriptionTopic) (*AidboxSubscriptionTopic, error) {
	response := &AidboxSubscriptionTopic{}
	return response, apiClient.updateResource(ctx, q, response)
}

func (apiClient *ApiClient) DeleteAidboxSubscriptionTopic(ctx context.Context, id string) error {
	return apiClient.deleteResource(ctx, id, &AidboxSubscriptionTopic{})
}

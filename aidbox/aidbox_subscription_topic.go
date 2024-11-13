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

func (apiClient *ApiClient) CreateAidboxSubscriptionTopic(ctx context.Context, searchParameter *AidboxSubscriptionTopic) (*AidboxSubscriptionTopic, error) {
	c, err := apiClient.createResource(ctx, searchParameter)
	if err != nil {
		return nil, err
	}
	return c.(*AidboxSubscriptionTopic), nil
}

func (apiClient *ApiClient) GetAidboxSubscriptionTopic(ctx context.Context, id string) (*AidboxSubscriptionTopic, error) {
	rr, err := apiClient.getResource(ctx, "/AidboxSubscriptionTopic/"+id)
	if err != nil {
		return nil, err
	}
	return rr.(*AidboxSubscriptionTopic), nil
}

func (apiClient *ApiClient) UpdateAidboxSubscriptionTopic(ctx context.Context, q *AidboxSubscriptionTopic) (*AidboxSubscriptionTopic, error) {
	rr, err := apiClient.updateResource(ctx, q)
	if err != nil {
		return nil, err
	}
	return rr.(*AidboxSubscriptionTopic), nil
}

func (apiClient *ApiClient) DeleteAidboxSubscriptionTopic(ctx context.Context, id string) error {
	return apiClient.deleteResource(ctx, "/AidboxSubscriptionTopic/"+id)
}

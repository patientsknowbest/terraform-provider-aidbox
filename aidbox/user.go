package aidbox

import "context"

type UserTwoFactor struct {
	Enabled bool `json:"enabled"`
}
type User struct {
	ResourceBase
	TwoFactor UserTwoFactor `json:"twoFactor"`
}

func (*User) GetResourcePath() string {
	return "User"
}

func (apiClient *ApiClient) GetUser(ctx context.Context, id string) (*User, error) {
	response := &User{}
	return response, apiClient.getResource(ctx, id, response)
}

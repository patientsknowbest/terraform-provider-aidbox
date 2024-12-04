package aidbox

import "context"

type UserTwoFactor struct {
	Enabled bool `json:"enabled"`
}
type User struct {
	ResourceBase
	TwoFactor UserTwoFactor `json:"twoFactor"`
}

func (*User) GetResourceName() string {
	return "User"
}

func (apiClient *ApiClient) GetUser(ctx context.Context, id, boxId string) (*User, error) {
	rr, err := apiClient.getResource(ctx, "/User/"+id, boxId)
	if err != nil {
		return nil, err
	}
	return rr.(*User), nil
}

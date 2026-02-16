package aidbox

import (
	"context"
)

// QuestionnaireTheme Represents the custom aidbox resource spec "QuestionnaireTheme"
// https://www.health-samurai.io/docs/aidbox/reference/system-resources-reference/sdc-module-resources#questionnairetheme
type QuestionnaireTheme struct {
	ResourceBase
	ThemeName    string `json:"theme-name,omitempty"`
	DesignSystem string `json:"design-system,omitempty"`
}

func (*QuestionnaireTheme) GetResourcePath() string {
	return "QuestionnaireTheme"
}

func (apiClient *ApiClient) CreateQuestionnaireTheme(ctx context.Context, theme *QuestionnaireTheme) (*QuestionnaireTheme, error) {
	response := &QuestionnaireTheme{}
	return response, apiClient.createResource(ctx, theme, response)
}

func (apiClient *ApiClient) GetQuestionnaireTheme(ctx context.Context, id string) (*QuestionnaireTheme, error) {
	response := &QuestionnaireTheme{}
	return response, apiClient.getResource(ctx, id, response)
}

func (apiClient *ApiClient) UpdateQuestionnaireTheme(ctx context.Context, theme *QuestionnaireTheme) (*QuestionnaireTheme, error) {
	response := &QuestionnaireTheme{}
	return response, apiClient.updateResource(ctx, theme, response)
}

func (apiClient *ApiClient) DeleteQuestionnaireTheme(ctx context.Context, id string) error {
	return apiClient.deleteResource(ctx, id, &QuestionnaireTheme{})
}

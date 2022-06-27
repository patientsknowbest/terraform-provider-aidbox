package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/patientsknowbest/terraform-provider-aidbox/internal/aidbox"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(client *aidbox.Client) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"client_id": {
					Type:        schema.TypeString,
					Description: "The client ID to access aidbox API",
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("AIDBOX_CLIENT_ID", "root"),
				},
				"client_secret": {
					Type:        schema.TypeString,
					Description: "The client secret to access aidbox API",
					Required:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("AIDBOX_CLIENT_SECRET", "secret"),
				},
				"url": {
					Type:        schema.TypeString,
					Description: "The URL of aidbox API",
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("AIDBOX_URL", "http://localhost:8888/"),
				},
				"is_multibox": {
					Type:        schema.TypeBool,
					Description: "true if this is a multibox instance",
					Default:     false,
					Optional:    true,
				},
			},
			DataSourcesMap: map[string]*schema.Resource{},
			ResourcesMap: map[string]*schema.Resource{
				"aidbox_token_introspector": resourceTokenIntrospector(),
				"aidbox_access_policy":      resourceAccessPolicy(),
				"aidbox_box":                resourceBox(),
				"aidbox_auth_client":        resourceAuthClient(),
			},
		}

		p.ConfigureContextFunc = func(cx context.Context, rd *schema.ResourceData) (interface{}, diag.Diagnostics) {
			if client != nil {
				return client, nil
			}
			var clientId, clientSecret, url string
			urlI, ok := rd.GetOk("url")
			if !ok {
				return nil, diag.Errorf("url is not supplied")
			}
			url, ok = urlI.(string)
			if !ok {
				return nil, diag.Errorf("url is wrong type")
			}
			clientIdI, ok := rd.GetOk("client_id")
			if !ok {
				return nil, diag.Errorf("client_id is not supplied")
			}
			clientId, ok = clientIdI.(string)
			if !ok {
				return nil, diag.Errorf("client_id is wrong type")
			}
			clientSecretI, ok := rd.GetOk("client_secret")
			if !ok {
				return nil, diag.Errorf("client_secret is not supplied")
			}
			clientSecret, ok = clientSecretI.(string)
			if !ok {
				return nil, diag.Errorf("client_secret is wrong type")
			}
			isMultiboxI, ok := rd.GetOk("is_multibox")
			isMultibox := ok && isMultiboxI.(bool)
			return aidbox.NewClient(url, clientId, clientSecret, isMultibox), nil
		}

		return p
	}
}

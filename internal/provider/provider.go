package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-scaffolding/internal/aidbox"
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
			},
			DataSourcesMap: map[string]*schema.Resource{
				"aidbox_token_introspector": dataSourceTokenIntrospector(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"aidbox_token_introspector": resourceTokenIntrospector(),
			},
		}

		p.ConfigureContextFunc = func(cx context.Context, rd *schema.ResourceData) (interface{}, diag.Diagnostics) {
			if client != nil {
				return client, nil
			}
			url := rd.Get("url").(string)
			username := rd.Get("username").(string)
			password := rd.Get("password").(string)
			return aidbox.NewClient(url, username, password), nil
		}

		return p
	}
}

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

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
                "username": {
                    Type:        schema.TypeString,
                    Description: "The username to access aidbox API",
					Required: true,
                },
				"password": {
					Type: schema.TypeString,
					Description: "The password to access aidbox API",
					Required: true,
					Sensitive: true,
					
				},
				"url": {
					Type: schema.TypeString,
					Description: "The URL of aidbox API",
					Required: true,
				},
            },
			DataSourcesMap: map[string]*schema.Resource{
				"aidbox_token_introspector": dataSourceTokenIntrospector(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"aidbox_token_introspector": resourceTokenIntrospector(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, rd *schema.ResourceData) (interface{}, diag.Diagnostics) {
		url := rd.Get("url").(string)
		username := rd.Get("username").(string)
		password := rd.Get("password").(string)
		return aidbox.NewClient(url, username, password), nil
	}
}

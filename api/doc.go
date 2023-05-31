package api

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/goai"
)

func InitDoc(s *ghttp.Server) error {
	//glog.Debug(ctx, "EnhanceOpenAPIDoc init ...")
	openapi := s.GetOpenApi()
	openapi.Config.CommonResponse = ghttp.DefaultHandlerResponse{}
	openapi.Config.CommonResponseDataField = `Data`
	openapi.Components = goai.Components{
		SecuritySchemes: goai.SecuritySchemes{
			//"ApiKeyAuth": goai.SecuritySchemeRef{
			//	Ref: "",
			//	Value: &goai.SecurityScheme{
			//		Type: "apiKey",
			//		In:   "header",
			//		Name: "token",
			//	},
			//},
			"BasicAuth": goai.SecuritySchemeRef{
				Ref: "",
				Value: &goai.SecurityScheme{
					Type:   "http",
					Scheme: "basic",
				},
			},
			//"BearerAuth": goai.SecuritySchemeRef{
			//	Value: &goai.SecurityScheme{
			//		Type:   "http",
			//		Scheme: "bearer",
			//	},
			//},
		},
	}
	openapi.Security = &goai.SecurityRequirements{map[string][]string{
		//"BearerAuth": {}, // BearerAuth ApiKeyAuth BasicAuth
		"BasicAuth": {},
	}}

	// API description.
	openapi.Info = goai.Info{
		Title:       "Go Plus Framework",
		Description: `Go Plus Framework is based on GoFrame `,
		Version:     "3.0.1",
		Contact: &goai.Contact{
			Name: "NetX",
			URL:  "https://goframe.org",
		},
	}

	// Sort the tags in custom sequence.
	openapi.Tags = &goai.Tags{
		//{Name: consts.OpenAPITagNameUser},
	}
	return nil
}

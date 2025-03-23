package gen

//go:generate go tool oapi-codegen -config oapi_server.config.yaml api/openapi/user/v1/api.yaml
//go:generate go tool oapi-codegen -config oapi_models.config.yaml api/openapi/user/v1/api.yaml

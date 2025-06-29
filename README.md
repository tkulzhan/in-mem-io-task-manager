Generate code for endpoints and components specified in api/openapi3/api.yaml:
```sh
go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=internal/infrastructure/http/server/config.yaml ./api/openapi3/api.yaml
```
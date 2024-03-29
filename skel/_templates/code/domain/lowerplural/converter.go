package lowerplural

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel/_templates/code/domain/api"
	//#if REPOSITORY_COCKROACH
	db "cto-github.cisco.com/NFV-BU/go-msx/sqldb/prepared"
	//#endif REPOSITORY_COCKROACH
	"github.com/google/uuid"
)

type lowerCamelSingularConverter struct{}

func (c *lowerCamelSingularConverter) FromCreateRequest(request api.UpperCamelSingularCreateRequest) lowerCamelSingular {
	return lowerCamelSingular{
		UpperCamelSingularId: uuid.New(),
		//#if TENANT_DOMAIN
		TenantId: db.ToModelUuid(request.TenantId),
		//#endif TENANT_DOMAIN
		Data: request.Data,
	}
}

func (c *lowerCamelSingularConverter) FromUpdateRequest(target lowerCamelSingular, request api.UpperCamelSingularUpdateRequest) lowerCamelSingular {
	target.Data = request.Data
	return target
}

func (c *lowerCamelSingularConverter) ToUpperCamelSingularListResponse(sources []lowerCamelSingular) (results []api.UpperCamelSingularResponse) {
	results = []api.UpperCamelSingularResponse{}
	for _, source := range sources {
		results = append(results, c.ToUpperCamelSingularResponse(source))
	}
	return
}

func (c *lowerCamelSingularConverter) ToUpperCamelSingularResponse(source lowerCamelSingular) api.UpperCamelSingularResponse {
	return api.UpperCamelSingularResponse{
		UpperCamelSingularId: db.ToApiUuid(source.UpperCamelSingularId),
		//#if TENANT_DOMAIN
		TenantId: db.ToApiUuid(source.TenantId),
		//#endif TENANT_DOMAIN
		Data: source.Data,
	}
}

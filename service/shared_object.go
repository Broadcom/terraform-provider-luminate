package service

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"context"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/antihax/optional"
	"github.com/pkg/errors"
)

type SharedObjectAPI struct {
	cli *sdk.APIClient
}

func NewSharedObjectAPI(client *sdk.APIClient) *SharedObjectAPI {
	return &SharedObjectAPI{
		cli: client,
	}
}

func (api *SharedObjectAPI) ListSharedObjects(sort string, size float64, page float64, filter string, objectType string) ([]dto.SharedObjectDTO, error) {
	options := &sdk.SharedObjectsApiListSharedObjectsOpts{
		Sort:   optional.NewString(sort),
		Size:   optional.NewFloat64(size),
		Page:   optional.NewFloat64(page),
		Filter: optional.NewString(filter),
		Type_:  optional.NewString(objectType),
	}
	sharedObjectsPage, _, err := api.cli.SharedObjectsApi.ListSharedObjects(context.Background(), options)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list shared objects")
	}

	sharedObjects := make([]dto.SharedObjectDTO, 0, len(sharedObjectsPage.Content))
	for _, obj := range sharedObjectsPage.Content {
		sharedObject := dto.SharedObjectDTO{
			ID:        obj.Id,
			Name:      obj.Name,
			Type:      obj.Type_,
			CreatedAt: obj.CreatedAt,
			UpdatedAt: obj.ModifiedOn,
		}
		sharedObject.Values = make([]interface{}, 0, len(obj.Values))
		for _, value := range obj.Values {
			sharedObject.Values = append(sharedObject.Values, value.Value)
		}

		sharedObjects = append(sharedObjects, sharedObject)
	}

	return sharedObjects, nil
}

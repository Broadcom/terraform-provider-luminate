package service

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"context"
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	serviceUtils "github.com/Broadcom/terraform-provider-luminate/service/utils"
	"github.com/Broadcom/terraform-provider-luminate/utils"
	"github.com/antihax/optional"
	"github.com/pkg/errors"
	"log"
)

type ApplicationAPI struct {
	cli *sdk.APIClient
}

func NewApplicationAPI(client *sdk.APIClient) *ApplicationAPI {
	return &ApplicationAPI{
		cli: client,
	}
}

func (api *ApplicationAPI) CreateApplication(application *dto.Application) (*dto.Application, error) {

	app := dto.ConvertFromApplicationDTO(*application)

	appOpts := sdk.ApplicationsApiCreateApplicationOpts{
		Body: optional.NewInterface(app),
	}
	log.Printf("[DEBUG] - Creating App")
	newApp, resp, err := api.cli.ApplicationsApi.CreateApplication(context.Background(), &appOpts)
	if err != nil {
		if resp != nil {
			body, _ := serviceUtils.ConvertReaderToString(resp.Body)
			return nil, errors.Wrapf(err, "received status code: %d ('%s')", resp.StatusCode, body)
		}

		return nil, err
	}
	log.Printf("[DEBUG] - Done Creating App")
	if resp != nil {
		if resp.StatusCode != 201 {
			errMsg := fmt.Sprintf("received bad status code creating application. Status Code: %d", resp.StatusCode)
			return nil, errors.New(errMsg)
		}
	} else {
		return nil, errors.New("received empty response from the server")
	}
	application.ID = newApp.Id

	createdApplication := dto.ConvertToApplicationDTO(newApp)
	createdApplication.SiteID = application.SiteID

	return &createdApplication, nil
}

func (api *ApplicationAPI) DeleteApplication(applicationID string) error {
	resp, err := api.cli.ApplicationsApi.DeleteApplication(context.Background(), applicationID)
	if err != nil {
		return err
	}
	if resp != nil {
		if resp.StatusCode != 204 {
			errMsg := fmt.Sprintf("received bad status code reading application. Status Code: %d", resp.StatusCode)
			return errors.New(errMsg)
		}
	} else {
		return errors.New("received empty response from the server")
	}

	return nil
}

func (api *ApplicationAPI) GetApplicationById(applicationID string) (*dto.Application, error) {
	app, resp, err := api.cli.ApplicationsApi.GetApplication(context.Background(), applicationID)

	if resp != nil && (resp.StatusCode == 404 || resp.StatusCode == 403) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	flatApp := dto.ConvertToApplicationDTO(app)

	return &flatApp, nil
}

func (api *ApplicationAPI) UpdateApplication(application *dto.Application) (*dto.Application, error) {
	app := dto.ConvertFromApplicationDTO(*application)

	appOpts := sdk.ApplicationsApiUpdateApplicationOpts{
		Body: optional.NewInterface(app),
	}

	log.Printf("[DEBUG] - Updating App")
	updatedApp, resp, err := api.cli.ApplicationsApi.UpdateApplication(context.Background(), application.ID, &appOpts)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] - Done Updating App")
	if resp != nil {
		if resp.StatusCode != 200 {
			errMsg := fmt.Sprintf("received bad status code updating application. Status Code: %d", resp.StatusCode)
			return nil, errors.New(errMsg)
		}
	} else {
		return nil, errors.New("received empty response from the server")
	}

	application.ID = updatedApp.Id
	return application, nil
}

func (api *ApplicationAPI) BindApplicationToSite(application *dto.Application, siteID string) error {
	log.Printf("[DEBUG] - Update Binding App")

	resp, err := api.cli.ApplicationsApi.BindApplicationToSite(context.Background(), application.ID, siteID, nil)
	// if bind fail with 400, there is chance that collection FT is enabled for this tenant, for BC we will try to link default collection
	// and bind again
	if err != nil {
		if resp.StatusCode == 400 {
			err = api.linkSiteToDefaultCollectionIfNeeded(siteID)
			if err != nil {
				return err
			}
			resp, err = api.cli.ApplicationsApi.BindApplicationToSite(context.Background(), application.ID, siteID, nil)
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
	}
	if resp != nil {
		if resp.StatusCode != 200 {
			errMsg := fmt.Sprintf("received bad status code binding application. Status Code: %d", resp.StatusCode)
			return errors.New(errMsg)
		}
	} else {
		return errors.New("received empty response from the server")
	}

	application.SiteID = siteID
	return nil
}

func (api *ApplicationAPI) linkSiteToDefaultCollectionIfNeeded(siteID string) error {
	collectionAPI := NewCollectionAPI(api.cli)

	collectionSiteLink := dto.CollectionSiteLink{
		CollectionID: utils.DefaultCollection,
		SiteID:       siteID,
	}

	_, err := collectionAPI.LinkSiteToCollection([]dto.CollectionSiteLink{collectionSiteLink})
	if err != nil {
		return err
	}

	return nil
}

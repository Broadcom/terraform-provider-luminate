package service

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"bitbucket.org/accezz-io/terraform-provider-symcsc/service/dto"
	"bitbucket.org/accezz-io/terraform-provider-symcsc/service/utils"
	"context"
	"github.com/pkg/errors"
	"fmt"
	"github.com/antihax/optional"
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

	appOpts := sdk.CreateApplicationOpts{
		Body: optional.NewInterface(app),
	}
	log.Printf("[DEBUG] - Creating App")
	log.Printf("[DEBUG APP DATA %v", app)
	newApp, resp, err := api.cli.ApplicationsApi.CreateApplication(context.Background(), &appOpts)
	if err != nil {
		if resp != nil  {
			body, _ := utils.ConvertReaderToString(resp.Body)
			return nil, errors.Wrap(err, fmt.Sprintf("received status code: %d ('%s')", resp.StatusCode, body))
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

	if resp != nil && resp.StatusCode == 404 {
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

	appOpts := sdk.UpdateApplicationOpts{
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
	resp, err := api.cli.ApplicationsApi.BindApplicationToSite(context.Background(), application.ID, siteID)
	if err != nil {
		return err
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

package service

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"context"
	"errors"
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/antihax/optional"
)

type SiteAPI struct {
	cli *sdk.APIClient
}

func NewSiteAPI(client *sdk.APIClient) *SiteAPI {
	return &SiteAPI{
		cli: client,
	}
}

func (api *SiteAPI) GetSiteByID(SiteID string) (*dto.Site, error) {
	s, resp, err := api.cli.SitesApi.GetSite(context.Background(), SiteID)
	if resp != nil && (resp.StatusCode == 403 || resp.StatusCode == 404) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	site := dto.Site{
		ID:               s.Id,
		Name:             s.Name,
		Region:           s.Region,
		MuteHealth:       s.MuteHealthNotification,
		K8SVolume:        s.KubernetesPersistentVolumeName,
		CountCollections: s.CountCollections,
	}

	if s.KerberosConfiguration != nil {
		site.Kerberos = &dto.SiteKerberosConfig{
			Domain:     s.KerberosConfiguration.Domain,
			KDCAddress: s.KerberosConfiguration.KdcAddress,
			KeytabPair: s.KerberosConfiguration.KeytabPath,
		}
	}

	for _, v := range s.ConnectorObjects {
		c := dto.Connector{
			Name:    v.Name,
			ID:      v.Id,
			Type:    ConnectorTypeFromDeployment(*v.DeploymentType),
			Enabled: v.Enabled,
		}
		site.Connectors = append(site.Connectors, c)
	}

	return &site, nil
}

func (api *SiteAPI) CreateSite(site *dto.Site) (*dto.Site, error) {

	newSite := sdk.Site{
		Name:                           site.Name,
		Region:                         site.Region,
		MuteHealthNotification:         site.MuteHealth,
		KubernetesPersistentVolumeName: site.K8SVolume,
		CountCollections:               site.CountCollections,
	}

	if site.Kerberos != nil {
		newSite.KerberosConfiguration = &sdk.KerberosConfiguration{
			Domain:     site.Kerberos.Domain,
			KdcAddress: site.Kerberos.KDCAddress,
			KeytabPath: site.Kerberos.KeytabPair,
		}
	}

	siteOpt := sdk.SitesApiCreateSiteOpts{
		Body: optional.NewInterface(newSite),
	}

	newSite, resp, err := api.cli.SitesApi.CreateSite(context.Background(), &siteOpt)
	if err != nil {
		return nil, err
	}

	if resp != nil {
		if resp.StatusCode != 201 {
			errMsg := fmt.Sprintf("received bad status code creating site. Status Code: %d", resp.StatusCode)
			return nil, errors.New(errMsg)
		}
	} else {
		return nil, errors.New("received empty response from the server")
	}
	site.ID = newSite.Id

	return site, nil
}

func (api *SiteAPI) UpdateSite(site *dto.Site, siteID string) (*dto.Site, error) {

	updateSite := sdk.Site{
		Name:                           site.Name,
		Region:                         site.Region,
		MuteHealthNotification:         site.MuteHealth,
		KubernetesPersistentVolumeName: site.K8SVolume,
	}

	if site.Kerberos != nil {
		updateSite.KerberosConfiguration = &sdk.KerberosConfiguration{
			Domain:     site.Kerberos.Domain,
			KdcAddress: site.Kerberos.KDCAddress,
			KeytabPath: site.Kerberos.KeytabPair,
		}
	}

	siteOpt := sdk.SitesApiUpdateSiteOpts{
		Body: optional.NewInterface(updateSite),
	}

	_, resp, err := api.cli.SitesApi.UpdateSite(context.Background(), siteID, &siteOpt)
	if err != nil {
		return nil, err
	}

	if resp != nil {
		if resp.StatusCode != 200 {
			errMsg := fmt.Sprintf("received bad status code updating site. Status Code: %d", resp.StatusCode)
			return nil, errors.New(errMsg)
		}
	} else {
		return nil, errors.New("received empty response from the server")
	}

	return site, nil
}

func (api *SiteAPI) DeleteSite(siteID string) error {
	resp, err := api.cli.SitesApi.DeleteSite(context.Background(), siteID)
	if err != nil {
		return err
	}
	if resp != nil {
		if resp.StatusCode != 204 {
			errMsg := fmt.Sprintf("received bad status code deleting site. Status Code: %d", resp.StatusCode)
			return errors.New(errMsg)
		}
	} else {
		return errors.New("received empty response from the server")
	}
	return nil
}

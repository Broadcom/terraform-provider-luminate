package service

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"context"
	"github.com/Broadcom/terraform-provider-luminate/service/dto"
	"github.com/Broadcom/terraform-provider-luminate/utils"
)

type DNSResiliencyAPI struct {
	cli *sdk.APIClient
}

func NewDNSResiliencyAPI(client *sdk.APIClient) *DNSResiliencyAPI {
	return &DNSResiliencyAPI{cli: client}
}

func (d *DNSResiliencyAPI) CreateDNSGroup(DNSGroupInput *dto.DNSGroupInputDTO) (*dto.DNSGroupOutputDTO, error) {
	var domainSuffix []string
	for _, domain := range DNSGroupInput.DomainSuffixes {
		domainSuffix = append(domainSuffix, domain.(string))
	}
	body := sdk.WssintegrationtenantDnsgroupsBody{
		Name:              DNSGroupInput.Name,
		DomainSuffixes:    domainSuffix,
		SendNotifications: DNSGroupInput.SendNotification,
	}
	res, _, err := d.cli.DNSResiliencyApi.CreateDnsGroup(context.Background(), body)
	if err != nil {
		return nil, err
	}
	return dto.ConvertDnsGroupOutputTODTO(res), nil
}

func (d *DNSResiliencyAPI) UpdateDNSGroup(DNSGroupInput *dto.DNSGroupInputDTO, DNSGroupID string) (*dto.DNSGroupOutputDTO, error) {
	var domainSuffix []string
	for _, domain := range DNSGroupInput.DomainSuffixes {
		domainSuffix = append(domainSuffix, domain.(string))
	}
	body := sdk.DnsgroupsDnsGroupIdBody{
		Name:              DNSGroupInput.Name,
		DomainSuffixes:    domainSuffix,
		SendNotifications: DNSGroupInput.SendNotification,
	}
	res, _, err := d.cli.DNSResiliencyApi.UpdateDNSGroup(context.Background(), body, DNSGroupID)
	if err != nil {
		return nil, err
	}
	return dto.ConvertDnsGroupOutputTODTO(res), nil
}

func (d *DNSResiliencyAPI) DeleteDNSGroup(DNSGroupID string) error {
	_, err := d.cli.DNSResiliencyApi.DeleteDNSGroup(context.Background(), DNSGroupID)
	if err != nil {
		return err
	}
	return nil
}

func (d *DNSResiliencyAPI) GetDNSGroup(DNSGroupID string) (*dto.DNSGroupOutputDTO, error) {
	res, _, err := d.cli.DNSResiliencyApi.GetDNSGroup(context.Background(), DNSGroupID)
	if err != nil {
		return nil, err
	}
	return dto.ConvertDnsGroupOutputTODTO(res), nil
}

func (d *DNSResiliencyAPI) GetDNServer(DNSServerID, DNSGroupID string) (*dto.DNSServerOutputDTO, error) {
	res, _, err := d.cli.DNSResiliencyApi.GetDNSServer(context.Background(), DNSGroupID, DNSServerID)
	if err != nil {
		return nil, err
	}
	return dto.ConvertDnsServerTODTO(res), nil
}

func (d *DNSResiliencyAPI) ListDNServer(DNSGroupID string) ([]*dto.DNSServerOutputDTO, error) {
	res, _, err := d.cli.DNSResiliencyApi.GetDNSServersOfGroup(context.Background(), DNSGroupID)
	if err != nil {
		return nil, err
	}
	var DNSServers []*dto.DNSServerOutputDTO
	for i := 0; i < len(res); i++ {
		DNSServers = append(DNSServers, dto.ConvertDnsServerTODTO(res[i]))
	}
	return DNSServers, nil
}

func (d *DNSResiliencyAPI) CreateDNServer(DNSServerInput *dto.DNSServerInputDTO, DNSGroupID string) (*dto.DNSServerOutputDTO, error) {
	body := sdk.DnsGroupIdServersBody{
		Name:            DNSServerInput.Name,
		InternalAddress: DNSServerInput.InternalAddress,
		SiteId:          DNSServerInput.SiteID,
		GroupId:         DNSServerInput.GroupID,
	}
	res, _, err := d.cli.DNSResiliencyApi.CreateDnsServer(context.Background(), body, DNSGroupID)
	if err != nil {
		return nil, err
	}
	return dto.ConvertDnsServerTODTO(res), nil
}

func (d *DNSResiliencyAPI) UpdateDNServer(DNSServerInput *dto.DNSServerInputDTO, DNSGroupID, DNSServerID string) (*dto.DNSServerOutputDTO, error) {
	body := sdk.ServersServerIdBody{
		Name:            DNSServerInput.Name,
		InternalAddress: DNSServerInput.InternalAddress,
		SiteId:          DNSServerInput.SiteID,
		GroupId:         DNSServerInput.GroupID,
	}
	res, _, err := d.cli.DNSResiliencyApi.UpdateDNSServer(context.Background(), body, DNSGroupID, DNSServerID)
	if err != nil {
		return nil, err
	}
	return dto.ConvertDnsServerTODTO(res), nil
}

func (d *DNSResiliencyAPI) DeleteDNSServer(DNSServerIDs []string, DNSGroupID string) error {
	ids := sdk.ServersDeletebyidsBody{
		DnsServerIds: DNSServerIDs,
	}
	_, err := d.cli.DNSResiliencyApi.DeleteDNSServers(context.Background(), ids, DNSGroupID)
	if err != nil {
		return utils.ParseSwaggerError(err)
	}
	return nil
}

func (d *DNSResiliencyAPI) UpdateDNSServersOrder(DNSServerIDs []string, DNSGroupID string) error {
	ids := sdk.DnsGroupIdServerorderBody{
		DnsServerIds: DNSServerIDs,
	}
	_, err := d.cli.DNSResiliencyApi.UpdateDNSServersOrder(context.Background(), ids, DNSGroupID)
	if err != nil {
		return err
	}
	return nil
}

func (d *DNSResiliencyAPI) EnableDisableDNSGroups(dto *dto.EnableDisableDNSGroupDTO) error {
	body := sdk.DnsgroupsEnableByIdsBody{
		Enable:   dto.Enable,
		GroupIds: dto.GroupIDs,
	}
	_, err := d.cli.DNSResiliencyApi.EnableDisableDNSGroups(context.Background(), body)
	if err != nil {
		return err
	}
	return nil
}

package dto

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"github.com/google/uuid"
	"time"
)

type Site struct {
	ID               string
	Name             string
	MuteHealth       bool
	K8SVolume        string
	Kerberos         *SiteKerberosConfig
	Connectors       []Connector
	CountCollections int32
	Region           string
}

type SiteKerberosConfig struct {
	Domain     string
	KDCAddress string
	KeytabPair string
}

type Connector struct {
	ID        string
	Name      string
	Type      string
	Enabled   bool
	Command   string
	K8SVolume string
	SiteID    string
	OTP       string
}
type Vpc struct {
	IntegrationId string
	Region        string
	Vpc           string
	CidrBlock     string
}

type CloudIntegrationData struct {
	Tags      map[string]string
	SegmentId string
	Vpcs      []Vpc
}

type AwsIntegration struct {
	Name                 string
	Id                   string
	LuminateAwsAccountId string `json:"luminate_aws_account_id"`
	AwsExternalId        string `json:"aws_external_id"`
}

type AwsIntegrationBind struct {
	Name                 string
	Id                   string
	AwsRoleArn           string `json:"aws_role_arn"`
	LuminateAwsAccountId string `json:"luminate_aws_account_id"`
	AwsExternalId        string `json:"aws_external_id"`
}

type Application struct {
	ID                    string
	Name                  string
	SiteID                string
	CollectionID          string
	Type                  string
	SubType               string
	Icon                  string
	Visible               bool
	NotificationsEnabled  bool
	InternalAddress       string
	ExternalAddress       string
	LuminateAddress       string
	Subdomain             string
	CustomExternalAddress string
	//HTTP
	CustomRootPath                    string
	HealthURL                         string
	HealthMethod                      string
	DefaultContentRewriteRulesEnabled bool
	DefaultHeaderRewriteRulesEnabled  bool
	UseExternalAddressForHostAndSni   bool
	LinkedApplications                []*string
	HeaderCustomization               map[string]interface{}
	// SSH-GW
	CloudIntegrationData *CloudIntegrationData
	//TCP
	Targets             []TCPTarget
	WildcardCertificate string
	WildcardPrivateKey  string
	//SEGMENT
	SegmentSettings         *SegmentSettings
	MultipleSegmentSettings []*SegmentSettings
	//DNS
	DnsSettings *DnsSettings
}

type SegmentSettings struct {
	OriginalIP string `json:"original_ip"`
}

type DnsSettings struct {
	DomainSuffixes []string `json:"domainSuffixes"`
}

type TCPTarget struct {
	Address     string
	Ports       []int32
	PortMapping []int32
}

type Validators struct {
	WebVerification bool
	MFA             bool
}

type Device struct {
	OpswatMetaAccess           bool
	SymantecCloudSoc           bool
	SymantecWebSecurityService bool
}

type Conditions struct {
	SourceIp        []string
	SharedIpList    []string
	Location        []string
	ManagedDevice   Device
	UnmanagedDevice Device
}

type ActivityRule struct {
	Action     string
	Conditions *RuleConditions
}

type RuleConditions struct {
	FileDownloaded bool
	FileUploaded   bool
	UriAccessed    bool
	HttpCommand    bool
	Arguments      *RuleConditionArguments
}

type RuleConditionArguments struct {
	UriList  []string
	Commands []string
}

const (
	IpUuid           = "IP_LIST"
	SharedIpListUuid = "SHARED_IP_LIST"
	CountriesUuid    = "COUNTRIES"
	Authentication   = "AUTHENTICATION"
)

const (
	ValidatorWebVerification = "VALIDATOR_WEB_VERIFICATION"
	MFA                      = "VALIDATOR_MFA"
)

const (
	IpCondition                            = "IP_CONDITION"
	LocationRestrictionCondition           = "LOCATION_RESTRICTION"
	ManagedDeviceCondition                 = "IS_DEVICE_COMPLIANCE"
	ManagedDeviceCloudSocConditionArgument = "CloudSOC"
	ManagedDeviceOpswatConditionArgument   = "OPSWAT"
	ManagedDeviceWssConditionArgument      = "IsWSSIp"
	UnmanagedDeviceCondition               = "IS_NOT_WSS_IP"
)

const (
	FileDownloadedCondition = "FILE_DOWNLOADED"
	FileUploadedCondition   = "FILE_UPLOADED"
	URICondition            = "URI_CONDITION"
	HTTPCommandCondition    = "METHOD_CONDITION"
)

const (
	URIListRuleConditionArgument     = "URI_LIST"
	HTTPCommandRuleConditionArgument = "HTTP_METHODS"
)

const (
	AllowAction          = "ALLOW"
	BlockAction          = "BLOCK"
	BlockUserAction      = "BLOCK_USER"
	DisconnectUserAction = "DISCONNECT_USER"
)

type Policy struct {
	TargetProtocol    string
	Id                string
	Enabled           bool
	CreatedAt         time.Time
	Name              string
	DirectoryEntities []DirectoryEntity
	Applications      []string
	CollectionID      string
	Conditions        *Conditions
}

type AccessPolicy struct {
	Policy
	//ACCESS
	Validators  *Validators
	RdpSettings *PolicyRdpSettings
	SshSettings *PolicySshSettings
	TcpSettings *PolicyTcpSettings
}

type ActivityPolicy struct {
	Policy
	ActivityRules []ActivityRule
}

type DirectoryEntity struct {
	IdentifierInProvider string
	IdentityProviderId   string
	EntityType           string
	IdentityProviderType string
	DisplayName          string
}

type PolicyRdpSettings struct {
	LongTermPassword bool
}

type PolicySshSettings struct {
	Accounts             []string
	AutoMapping          bool
	FullUPNAutoMapping   bool
	AgentForward         bool
	AcceptTemporaryToken bool
	AcceptCertificate    bool
}

type PolicyTcpSettings struct {
	AcceptTemporaryToken bool
	AcceptCertificate    bool
}

type CollectionSiteLink struct {
	CollectionID string
	SiteID       string
}

type Collection struct {
	ID               uuid.UUID
	Name             string
	ParentId         uuid.UUID
	CountResources   int32
	CountLinkedSites int32
	Fqdn             string
}

type ListCollectionsRequest struct {
	Sort          string
	Size          float64
	Page          float64
	Name          string
	ApplicationId uuid.UUID
	SiteId        uuid.UUID
	PolicyId      uuid.UUID
}

type CreateRoleDTO struct {
	Role     string
	Entities []DirectoryEntity
}

type CreateCollectionRoleDTO struct {
	CreateRoleDTO
	CollectionID string
}

type CreateSiteRoleDTO struct {
	CreateRoleDTO
	SiteID string
}

type RoleBinding struct {
	ID            string
	EntityIDInIDP string
	EntityIDPID   string
	EntityType    string
	RoleType      string
	CollectionID  string
	ResourceID    string
}

type Group struct {
	Name               string
	ID                 string
	IdentityProviderId string
}

func EntityDTOToEntityModel(entities []DirectoryEntity) []sdk.DirectoryEntity {
	var directoryEntities []sdk.DirectoryEntity
	for _, directoryEntity := range entities {
		identityProviderType, err := ConvertIdentityProviderTypeToEnum(directoryEntity.IdentityProviderType)
		if err == nil {
			directoryEntities = append(directoryEntities, sdk.DirectoryEntity{
				IdentifierInProvider: directoryEntity.IdentifierInProvider,
				IdentityProviderId:   directoryEntity.IdentityProviderId,
				DisplayName:          directoryEntity.DisplayName,
				IdentityProviderType: &identityProviderType,
				Type_:                ToModelType(directoryEntity.EntityType),
			})
		}
	}
	return directoryEntities
}

func EntityModelEntityDTO(directoryEntitiesDTO []sdk.DirectoryEntity) []DirectoryEntity {
	var directoryEntities []DirectoryEntity
	for _, directoryEntityDto := range directoryEntitiesDTO {
		directoryEntities = append(directoryEntities, DirectoryEntity{
			IdentifierInProvider: directoryEntityDto.IdentifierInProvider,
			IdentityProviderId:   directoryEntityDto.IdentityProviderId,
			DisplayName:          directoryEntityDto.DisplayName,
			IdentityProviderType: ConvertIdentityProviderTypeToString(directoryEntityDto.IdentityProviderType),
			EntityType:           FromModelType(*directoryEntityDto.Type_),
		})
	}
	return directoryEntities
}

type DNSGroupInputDTO struct {
	Name             string
	IsEnabled        bool          `json:"isEnabled"`
	DomainSuffixes   []interface{} `json:"domainSuffixes"`
	SendNotification bool          `json:"sendNotification"`
}

type DNSGroupOutputDTO struct {
	ID                  string
	Name                string
	DomainSuffixes      []string `json:"domainSuffixes"`
	SendNotification    bool     `json:"sendNotification"`
	Status              string
	ServerInUse         string `json:"serverInUsed"`
	ActiveServerAddress string
	Servers             []string
	IsEnabled           bool `json:"isEnabled"`
}

type DNSServerInputDTO struct {
	Name            string
	InternalAddress string
	SiteID          string
	GroupID         string
}

type DNSServerOutputDTO struct {
	GroupID         string
	ID              string
	SiteID          string
	Name            string
	InternalAddress string
	HealthStatus    string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type EnableDisableDNSGroupDTO struct {
	Enable   bool
	GroupIDs []string
}

func ConvertDnsGroupOutputTODTO(dto sdk.DnsGroupOutput) *DNSGroupOutputDTO {
	return &DNSGroupOutputDTO{
		ID:                  dto.Id,
		Name:                dto.Name,
		DomainSuffixes:      dto.DomainSuffixes,
		SendNotification:    dto.SendNotifications,
		Status:              dto.Status,
		ServerInUse:         dto.ServerInUsed,
		ActiveServerAddress: dto.ActiveServerAddress,
		Servers:             dto.Servers,
	}
}

func ConvertDnsServerTODTO(dto sdk.DnsServerOutput) *DNSServerOutputDTO {
	return &DNSServerOutputDTO{
		GroupID:         dto.GroupId,
		ID:              dto.Id,
		SiteID:          dto.SiteId,
		Name:            dto.Name,
		InternalAddress: dto.InternalAddress,
		HealthStatus:    dto.HealthStatus,
		CreatedAt:       dto.CreatedAt,
		UpdatedAt:       dto.UpdatedAt,
	}
}

package dto

import "time"

type Site struct {
	ID         string
	Name       string
	MuteHealth bool
	K8SVolume  string
	Kerberos   *SiteKerberosConfig
	Connectors []Connector
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
	HeaderCustomization               []map[string]string
	// SSH-GW
	CloudIntegrationData *CloudIntegrationData
	//TCP
	Targets             []TCPTarget
	WildcardCertificate string
	WildcardPrivateKey  string
	//SEGMENT
	SegmentSettings *SegmentSettings
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
	Address string
	Ports   []float64
}

type Validators struct {
	ComplianceCheck bool
	WebVerification bool
}

type ManagedDevice struct {
	OpswatMetaAccess           bool
	SymantecCloudSoc           bool
	SymantecWebSecurityService bool
}
type Conditions struct {
	SourceIp        []string
	SharedIpList    []string
	Location        []string
	ManagedDevice   ManagedDevice
	UnmanagedDevice bool
}

const (
	IpUuid            = "IP_LIST"
	SharedIpListUuid  = "SHARED_IP_LIST"
	CountriesUuid     = "COUNTRIES"
	ManagedDeviceUuid = "AUTHENTICATION"
)

const (
	ValidatorComplianceCheck = "VALIDATOR_COMPLIANCE_CHECK"
	ValidatorWebVerification = "VALIDATOR_WEB_VERIFICATION"
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

type AccessPolicy struct {
	TargetProtocol    string
	Id                string
	Enabled           bool
	CreatedAt         time.Time
	Name              string
	DirectoryEntities []DirectoryEntity
	Applications      []string
	Conditions        *Conditions
	Validators        *Validators
	RdpSettings       *PolicyRdpSettings
	SshSettings       *PolicySshSettings
	TcpSettings       *PolicyTcpSettings
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

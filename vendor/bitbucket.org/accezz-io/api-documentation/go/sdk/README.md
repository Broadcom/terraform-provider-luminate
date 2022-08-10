# Go API client for swagger

 ## Introduction  Secure Access Cloud API uses common RESTful resourced based URL conventions and JSON as the exchange format. <br> Properties names are case-sensitive. <br> Some of Secure Access Cloud API calls omit None values from the API response.  The base-URL is `api.`&lt;`tenant-name`&gt;`.luminatesec.com`. For example, if your administration portal URL is _admin.acme.luminatesec.com_, then your API base-URL is _api.acme.luminatesec.com_.  All examples below are performed on a tenant called acme.  ## Common Operations Steps  Below you may find a list of common operations and the relevant API calls for each. Each of these operations can also be performed by using the administrative portal at https://admin.acme.luminatesec.com.  <ol>   <li>     Creating a site and deploying a connector:     <ol type=\"a\">       <li>Creating a new site using the <a href=\"#operation/createSite\">Create site API</a>.<br></li>       <li>         Once a site is created you can use its Id (returned in the response of the Create Site request)         and call the <a href=\"#operation/createConnector\">Create connector API</a>. <br>       </li>       <li>         Deploy the Secure Access Cloud connector:         <ol type=\"i\">           <li>Retrieve the deployment command using the <a href=\"#operation/getConnectorCommand\">Connector Deployment Command API.</a> <br> </li>           <li>Execute the command on the target machine.</li>         </ol>       </li>     </ol>   </li>   <li>     Creating an application:       <ol type=\"a\">         <li>           An application is always associated with a specific site for routing the traffic to the application via the connectors associated with the same site.           In order to create the application, call the <a href=\"#operation/createApplication\">Create Application API</a>         </li>         <li>           Once the application is created, you *must* assign the application to a specific site in order to make it accessible. Assign the application to the required site           using the <a href=\"#operation/BindApplicationToSite\">Bind Application to Site API</a>.         </li>         <li>           In order to grant access to the application for specific entities (users/groups), you should assign the application to the access policy using the <a href=\"#tag/Access-and-Activity-Policies\">Access and Activity Policy API</a>         </li>       </ol>   </li> </ol>  ## Object Model The object model of the API is built around the following: <ol>   <li><a href=\"#tag/Sites\">Sites</a> - Site is a representation of the physical or virtual data center your applications reside in.</li>   <li><a href=\"#tag/Connectors\">Connectors</a> - A connector is a lightweight piece of software connecting your site to the Secure Access Cloud platform.</li>   <li><a href=\"#tag/Applications\">Applications</a>  - Application is the internal resource you would like to publish using Secure Access Cloud. </li>   <li>     <a href=\"#tag/Access-and-Activity-Policies\">Access and Activity Policies</a> - Secure Access Cloud continuously authorize each user request for the contextual access and activity,     in order to control access to resources and restrict user’s actions within resources, based on the user/device context (such as the user’s group membership, user’s location,     MFA status and managed/unmanaged device status) and the requested resource.   <li>     <a href=\"#tag/Cloud-Integration\">Cloud Integration</a> - Integration with Cloud Providers like Amazon Web Services to provide a smoother and cloud-native integration with SIEM solutions      and to allow access to resources based on their associated tags.   <li>     Logs - Secure Access Cloud internal logs for audit and forensics purposes:     <ol>       <li><a href=\"#tag/Audit-Logs\">Audit Logs</a> audit all operations done through the administration portal</li>       <li><a href=\"#tag/Forensics-Logs\">Forensics Logs</a> audit any user's access to any application as well as user's activity for any application.</li>     </ol>   </li> </ol>   ## Authentication  Authentication is done using [OAuth2](https://tools.ietf.org/html/rfc6749) with the [Bearer authentication scheme](https://tools.ietf.org/html/rfc6750).  <!-- ReDoc-Inject: <security-definitions> -->  The Secure Access Cloud API is available to Secure Access Cloud users who have administrative privileges in their Secure Access Cloud tenant. An administrator should create an API client through the Secure Access Cloud Admin portal, check the ‘Allow access to Secure Access Cloud management API’ permission and copy the ‘Client Id’ and the ‘Client Secret’.  Retrieving the API access token is done using Basic-Authentication scheme, POST of a Base64 encoded Client-ID and Client-Secret: <B>   ``` curl -X POST \\  https://api.acme.luminatesec.com/v1/oauth/token \\  -u yourApiClientId:yourApiClientSecret   ``` </B>  This call returns the following JSON: {     \"access_token\":\"edfe22e3-eb4c-4c83-8ce3-3152e6a2XXX\",   \"expires_in\":86400,   \"scope\":\"luminate-scope\",   \"token_type\":\"Bearer\",   \"error\":\"\",   \"error_description\":\"\"}  All further API calls should include the ‘Authorization’ header with value “Bearer AccessToken”  For example: <B>   ```   curl -H \"Authorization: Bearer edfe22e3-eb4c-4c83-8ce3-3152e6a2XXX\" \"https://api.acme.luminatesec.com/v2/applications\"   ``` </B>  ## Versioning and Compatibility  The latest Major Version is `v2`.  The Major Version is included in the URL path (e.g. /v2/applications ) and it denotes breaking changes to the API. Minor and Patch versions are transparent to the client.  ## Pagination   Some of our API responses are paginated, meaning that only a certain number of items are returned at a time.  The default number of items returned in a single page is 50.  You can override this by passing a size parameter to set the maximum number of results, but cannot exceed 100.  Specifying the page number sets the starting point for the result set, allowing you to fetch subsequent items  that are not in the initial set of results. The sort order for returned data can be controlled using the sort parameter.<br>  You can constrain the results by using a filter. <br><br>  **Note:** Most methods that support pagination use the approach specified above. However, some methods use varied   versions of pagination. The individual documentation for each API method is your source of truth for which pattern the method follows.  ## Auditing  All authentication operations and modify operations (POST, PUT, DELETE) are audited.   ## Rate-limiting The API has a rate limit of 5 requests per second. If you have hit the rate limit, then a ‘429’ status code will be returned. In such cases, you should back-off from submitting new requests for 1 second before resuming.  Note that rate-limitation applies to the accumulated requests of **all** of your clients. For example, if you have 6 clients submitting requests simultaneously at a rate of 1 request per second for each one then one of them is likely to get a 429 status code.  ## Support  For additional help you may refer to our support at https://support.luminate.io.  Each request submitted to the API returns a unique request ID that is generated by the API. The request ID will be returned in header `x-lum-request-id`. If you need to contact us about any specific request then this ID will serve as a reference to the given request. 

## Overview
This API client was generated by the [swagger-codegen](https://github.com/swagger-api/swagger-codegen) project.  By using the [swagger-spec](https://github.com/swagger-api/swagger-spec) from a remote server, you can easily generate an API client.

- API version: V2
- Package version: 1.0.0
- Build package: io.swagger.codegen.languages.GoClientCodegen

## Installation
Put the package under your project folder and add the following in import:
```golang
import "./swagger"
```

## Documentation for API Endpoints

All URIs are relative to *https://api.acme.luminatesec.com/v2*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*AccessAndActivityPoliciesApi* | [**AssignApplicationToPolicies**](docs/AccessAndActivityPoliciesApi.md#assignapplicationtopolicies) | **Post** /policies/by-app-id/{application-id} | AssignApplicationToPolicies
*AccessAndActivityPoliciesApi* | [**CreatePolicy**](docs/AccessAndActivityPoliciesApi.md#createpolicy) | **Post** /policies | createPolicy
*AccessAndActivityPoliciesApi* | [**DeletePolicy**](docs/AccessAndActivityPoliciesApi.md#deletepolicy) | **Delete** /policies/{policy-id} | Delete Policy
*AccessAndActivityPoliciesApi* | [**GetAllPolicies**](docs/AccessAndActivityPoliciesApi.md#getallpolicies) | **Get** /policies | getAllPolicies
*AccessAndActivityPoliciesApi* | [**GetApplicationAssignedPolicies_**](docs/AccessAndActivityPoliciesApi.md#getapplicationassignedpolicies_) | **Get** /policies/by-app-id/{application-id} | Get Application Assigned Policies.
*AccessAndActivityPoliciesApi* | [**GetPolicy**](docs/AccessAndActivityPoliciesApi.md#getpolicy) | **Get** /policies/{policy-id} | Get Policy
*AccessAndActivityPoliciesApi* | [**GetSupportedActions**](docs/AccessAndActivityPoliciesApi.md#getsupportedactions) | **Get** /policies/config/action-types | getSupportedActions
*AccessAndActivityPoliciesApi* | [**GetSupportedConditions**](docs/AccessAndActivityPoliciesApi.md#getsupportedconditions) | **Get** /policies/config/condition-definitions | getSupportedConditions
*AccessAndActivityPoliciesApi* | [**GetSupportedValidators**](docs/AccessAndActivityPoliciesApi.md#getsupportedvalidators) | **Get** /policies/config/validator-types | getSupportedValidators
*AccessAndActivityPoliciesApi* | [**Removeapplicationfrompolicies**](docs/AccessAndActivityPoliciesApi.md#removeapplicationfrompolicies) | **Delete** /policies/by-app-id/{application-id} | Remove application from policies
*AccessAndActivityPoliciesApi* | [**UpdatePolicy**](docs/AccessAndActivityPoliciesApi.md#updatepolicy) | **Put** /policies/{policy-id} | Update Policy
*AccessAndActivityPoliciesApi* | [**Updateapplicationinpolicies**](docs/AccessAndActivityPoliciesApi.md#updateapplicationinpolicies) | **Put** /policies/by-app-id/{application-id} | Update application in policies
*ApplicationsApi* | [**BindApplicationToSite**](docs/ApplicationsApi.md#bindapplicationtosite) | **Put** /applications/{application-id}/site-binding/{site-id} | BindApplicationToSite
*ApplicationsApi* | [**CreateApplication**](docs/ApplicationsApi.md#createapplication) | **Post** /applications | createApplication
*ApplicationsApi* | [**DeleteApplication**](docs/ApplicationsApi.md#deleteapplication) | **Delete** /applications/{application-id} | Delete Application
*ApplicationsApi* | [**GetAllApps**](docs/ApplicationsApi.md#getallapps) | **Get** /applications | getAllApps
*ApplicationsApi* | [**GetApplication**](docs/ApplicationsApi.md#getapplication) | **Get** /applications/{application-id} | Get Application
*ApplicationsApi* | [**GetApplicationHealthStatus**](docs/ApplicationsApi.md#getapplicationhealthstatus) | **Get** /applications/{application-id}/status | Get Application Health Status
*ApplicationsApi* | [**UpdateApplication**](docs/ApplicationsApi.md#updateapplication) | **Put** /applications/{application-id} | Update Application
*AuditLogsApi* | [**SearchAuditlogs**](docs/AuditLogsApi.md#searchauditlogs) | **Post** /logs/audit | Search Audit logs
*CloudIntegrationApi* | [**CreateCloudIntegration**](docs/CloudIntegrationApi.md#createcloudintegration) | **Post** /cloud-integrations/integrations | createCloudIntegration
*CloudIntegrationApi* | [**DeleteCloudIntegration**](docs/CloudIntegrationApi.md#deletecloudintegration) | **Delete** /cloud-integrations/integrations/{cloud-integration-id} | deleteCloudIntegration
*CloudIntegrationApi* | [**GetCloudIntegration**](docs/CloudIntegrationApi.md#getcloudintegration) | **Get** /cloud-integrations/integrations/{cloud-integration-id} | getCloudIntegration
*CloudIntegrationApi* | [**ListCloudIntegrations**](docs/CloudIntegrationApi.md#listcloudintegrations) | **Get** /cloud-integrations/integrations | listCloudIntegrations
*CloudIntegrationApi* | [**UpdateCloudIntegration**](docs/CloudIntegrationApi.md#updatecloudintegration) | **Put** /cloud-integrations/integrations/{cloud-integration-id} | updateCloudIntegration
*ConnectorsApi* | [**CreateConnector**](docs/ConnectorsApi.md#createconnector) | **Post** /connectors | createConnector
*ConnectorsApi* | [**DeleteConnector**](docs/ConnectorsApi.md#deleteconnector) | **Delete** /connectors/{connector-id} | Delete Connector
*ConnectorsApi* | [**GetAllConnectors**](docs/ConnectorsApi.md#getallconnectors) | **Get** /connectors | getAllConnectors
*ConnectorsApi* | [**GetConnector**](docs/ConnectorsApi.md#getconnector) | **Get** /connectors/{connector-id} | Get Connector
*ConnectorsApi* | [**GetConnectorCommand**](docs/ConnectorsApi.md#getconnectorcommand) | **Get** /connectors/{connector-id}/command | getConnectorCommand
*ConnectorsApi* | [**GetConnectorEnvironmentVariables**](docs/ConnectorsApi.md#getconnectorenvironmentvariables) | **Get** /connectors/{connector-id}/environment_variables | Get Connector Environment Variables
*ConnectorsApi* | [**GetConnectorVersion**](docs/ConnectorsApi.md#getconnectorversion) | **Get** /connectors/version | Get Connector Version
*ForensicsLogsApi* | [**SearchForensicslogs**](docs/ForensicsLogsApi.md#searchforensicslogs) | **Post** /logs/forensics | Search Forensics logs
*GroupsApi* | [**AssignUserToGroup**](docs/GroupsApi.md#assignusertogroup) | **Put** /identities/local/groups/{group-id}/users/{user-id} | Assign User To Group
*GroupsApi* | [**GetGroup**](docs/GroupsApi.md#getgroup) | **Get** /identities/{identity-provider-id}/groups/{entity-id} | Get Group
*GroupsApi* | [**ListAssignedUsers**](docs/GroupsApi.md#listassignedusers) | **Get** /identities/{identity-provider-id}/groups/{entity-id}/users | List Assigned Users
*GroupsApi* | [**RemoveUserFromGroup**](docs/GroupsApi.md#removeuserfromgroup) | **Delete** /identities/local/groups/{group-id}/users/{user-id} | Remove User From Group
*GroupsApi* | [**SearchGroupsbyIdp**](docs/GroupsApi.md#searchgroupsbyidp) | **Get** /identities/{identity-provider-id}/groups | SearchGroupsbyIdp
*IdentityProvidersApi* | [**GetIdentityProviderById**](docs/IdentityProvidersApi.md#getidentityproviderbyid) | **Get** /identities/settings/identity-providers/{identity-provider-id} | GetIdentityProviderById
*IdentityProvidersApi* | [**ListIdentityProviders**](docs/IdentityProvidersApi.md#listidentityproviders) | **Get** /identities/settings/identity-providers | ListIdentityProviders
*SCIMApi* | [**CreateSCIMUser**](docs/SCIMApi.md#createscimuser) | **Post** /identities/{identity-provider-id}/scim/users | Create SCIM User
*SCIMApi* | [**CreateaSCIMGroup**](docs/SCIMApi.md#createascimgroup) | **Post** /identities/{identity-provider-id}/scim/groups | Create a SCIM Group
*SCIMApi* | [**DeleteSCIMGroup**](docs/SCIMApi.md#deletescimgroup) | **Delete** /identities/{identity-provider-id}/scim/groups/{group-id} | Delete SCIM Group
*SCIMApi* | [**DeleteSCIMUser**](docs/SCIMApi.md#deletescimuser) | **Delete** /identities/{identity-provider-id}/scim/users/{user-id} | Delete SCIM User
*SCIMApi* | [**GetSCIMGroup**](docs/SCIMApi.md#getscimgroup) | **Get** /identities/{identity-provider-id}/scim/groups/{group-id} | Get SCIM Group
*SCIMApi* | [**GetSCIMUser**](docs/SCIMApi.md#getscimuser) | **Get** /identities/{identity-provider-id}/scim/users/{user-id} | Get SCIM User
*SCIMApi* | [**ListSCIMGroupsAPI**](docs/SCIMApi.md#listscimgroupsapi) | **Get** /identities/{identity-provider-id}/scim/groups | listSCIMGroupsAPI
*SCIMApi* | [**ListSCIMUsersAPI**](docs/SCIMApi.md#listscimusersapi) | **Get** /identities/{identity-provider-id}/scim/users | listSCIMUsersAPI
*SCIMApi* | [**ModifyaSCIMGroup**](docs/SCIMApi.md#modifyascimgroup) | **Patch** /identities/{identity-provider-id}/scim/groups/{group-id} | Modify a SCIM Group
*SCIMApi* | [**UpdateSCIMGroup**](docs/SCIMApi.md#updatescimgroup) | **Put** /identities/{identity-provider-id}/scim/groups/{group-id} | Update SCIM Group
*SCIMApi* | [**UpdateSCIMUser**](docs/SCIMApi.md#updatescimuser) | **Put** /identities/{identity-provider-id}/scim/users/{user-id} | Update SCIM User
*SSHClientsApi* | [**GetAllSshClients**](docs/SSHClientsApi.md#getallsshclients) | **Get** /ssh-clients | getAllSshClients
*SharedObjectsApi* | [**CreateSharedObject**](docs/SharedObjectsApi.md#createsharedobject) | **Post** /policies/shared-objects | createSharedObject
*SharedObjectsApi* | [**DeleteSharedObject**](docs/SharedObjectsApi.md#deletesharedobject) | **Delete** /policies/shared-objects/{shared-object-id} | Delete Shared Object
*SharedObjectsApi* | [**GetSharedObject**](docs/SharedObjectsApi.md#getsharedobject) | **Get** /policies/shared-objects/{shared-object-id} | Get Shared Object
*SharedObjectsApi* | [**ListSharedObjects**](docs/SharedObjectsApi.md#listsharedobjects) | **Get** /policies/shared-objects | listSharedObjects
*SharedObjectsApi* | [**UpdateSharedObject**](docs/SharedObjectsApi.md#updatesharedobject) | **Put** /policies/shared-objects/{shared-object-id} | Update Shared Object
*SitesApi* | [**CreateSite**](docs/SitesApi.md#createsite) | **Post** /sites | createSite
*SitesApi* | [**DeleteSite**](docs/SitesApi.md#deletesite) | **Delete** /sites/{site-id} | Delete Site
*SitesApi* | [**GetAllSites**](docs/SitesApi.md#getallsites) | **Get** /sites | getAllSites
*SitesApi* | [**GetSite**](docs/SitesApi.md#getsite) | **Get** /sites/{site-id} | Get Site
*SitesApi* | [**GetSiteHealthStatus**](docs/SitesApi.md#getsitehealthstatus) | **Get** /sites/{site-id}/status | Get Site Health Status
*SitesApi* | [**UpdateSite**](docs/SitesApi.md#updatesite) | **Put** /sites/{site-id} | Update Site
*UsersApi* | [**BlockUser**](docs/UsersApi.md#blockuser) | **Post** /identities/{identity-provider-id}/users/{entity-id}/block | Block User
*UsersApi* | [**CreateLocalUser**](docs/UsersApi.md#createlocaluser) | **Post** /identities/local/users | Create Local User
*UsersApi* | [**DeleteLocalUser**](docs/UsersApi.md#deletelocaluser) | **Delete** /identities/local/users/{entity-id} | Delete Local User
*UsersApi* | [**GetUser**](docs/UsersApi.md#getuser) | **Get** /identities/{identity-provider-id}/users/{entity-id} | Get User
*UsersApi* | [**ListBlockedUsers**](docs/UsersApi.md#listblockedusers) | **Get** /identities/settings/blocked-users | List Blocked Users
*UsersApi* | [**SearchUsersbyIdp**](docs/UsersApi.md#searchusersbyidp) | **Get** /identities/{identity-provider-id}/users | SearchUsersbyIdp
*UsersApi* | [**UnblockUser**](docs/UsersApi.md#unblockuser) | **Delete** /identities/{identity-provider-id}/users/{entity-id}/block | Unblock User
*UsersApi* | [**UpdateLocalUser**](docs/UsersApi.md#updatelocaluser) | **Put** /identities/local/users/{entity-id} | Update Local User


## Documentation For Models

 - [ApplicationBase](docs/ApplicationBase.md)
 - [ApplicationCloudIntegrationData](docs/ApplicationCloudIntegrationData.md)
 - [ApplicationCloudIntegrationDataProperties](docs/ApplicationCloudIntegrationDataProperties.md)
 - [ApplicationCloudIntegrationTag](docs/ApplicationCloudIntegrationTag.md)
 - [ApplicationConnectionSettings](docs/ApplicationConnectionSettings.md)
 - [ApplicationConnectionSettingsRdp](docs/ApplicationConnectionSettingsRdp.md)
 - [ApplicationConnectionSettingsSsh](docs/ApplicationConnectionSettingsSsh.md)
 - [ApplicationConnectionSettingsTcp](docs/ApplicationConnectionSettingsTcp.md)
 - [ApplicationConnectorsStatus](docs/ApplicationConnectorsStatus.md)
 - [ApplicationCore](docs/ApplicationCore.md)
 - [ApplicationDynamicSsh](docs/ApplicationDynamicSsh.md)
 - [ApplicationHealth](docs/ApplicationHealth.md)
 - [ApplicationHttp](docs/ApplicationHttp.md)
 - [ApplicationLinkTranslationSettings](docs/ApplicationLinkTranslationSettings.md)
 - [ApplicationRdp](docs/ApplicationRdp.md)
 - [ApplicationRequestCustomizationSettings](docs/ApplicationRequestCustomizationSettings.md)
 - [ApplicationSort](docs/ApplicationSort.md)
 - [ApplicationSsh](docs/ApplicationSsh.md)
 - [ApplicationSubType](docs/ApplicationSubType.md)
 - [ApplicationTcp](docs/ApplicationTcp.md)
 - [ApplicationTcpTarget](docs/ApplicationTcpTarget.md)
 - [ApplicationTcpTunnelSettings](docs/ApplicationTcpTunnelSettings.md)
 - [ApplicationToPoliciesMapping](docs/ApplicationToPoliciesMapping.md)
 - [ApplicationType](docs/ApplicationType.md)
 - [ApplicationType1](docs/ApplicationType1.md)
 - [ApplicationTypeWithRuleIds](docs/ApplicationTypeWithRuleIds.md)
 - [ApplicationVpcData](docs/ApplicationVpcData.md)
 - [ApplicationsPage](docs/ApplicationsPage.md)
 - [AuditLogResult](docs/AuditLogResult.md)
 - [AuditLogSearchResults](docs/AuditLogSearchResults.md)
 - [BlockedUser](docs/BlockedUser.md)
 - [BulkApiErrorResponse](docs/BulkApiErrorResponse.md)
 - [BulkApiResponse](docs/BulkApiResponse.md)
 - [CloudIntegration](docs/CloudIntegration.md)
 - [CloudIntegrationBase](docs/CloudIntegrationBase.md)
 - [CloudIntegrationError](docs/CloudIntegrationError.md)
 - [CloudIntegrationHealth](docs/CloudIntegrationHealth.md)
 - [CloudIntegrationPost](docs/CloudIntegrationPost.md)
 - [CloudIntegrationProvider](docs/CloudIntegrationProvider.md)
 - [CloudIntegrationPut](docs/CloudIntegrationPut.md)
 - [Connector](docs/Connector.md)
 - [ConnectorDeploymentCommand](docs/ConnectorDeploymentCommand.md)
 - [ConnectorEnvironmentVariables](docs/ConnectorEnvironmentVariables.md)
 - [ConnectorLastSeen](docs/ConnectorLastSeen.md)
 - [ConnectorStatus](docs/ConnectorStatus.md)
 - [ConnectorVersion](docs/ConnectorVersion.md)
 - [ConnectorsPage](docs/ConnectorsPage.md)
 - [Data](docs/Data.md)
 - [Data1](docs/Data1.md)
 - [DeploymentType](docs/DeploymentType.md)
 - [Direction](docs/Direction.md)
 - [DirectoryEntity](docs/DirectoryEntity.md)
 - [DirectoryEntityBinding](docs/DirectoryEntityBinding.md)
 - [DirectoryEntityPaginatedResponseBase](docs/DirectoryEntityPaginatedResponseBase.md)
 - [DirectoryProvider](docs/DirectoryProvider.md)
 - [DirectoryProviderHealth](docs/DirectoryProviderHealth.md)
 - [DirectoryProviderInstructionsLdap](docs/DirectoryProviderInstructionsLdap.md)
 - [DirectoryProviderInstructionsOneLoginOrOkta](docs/DirectoryProviderInstructionsOneLoginOrOkta.md)
 - [DirectoryProviderSettingsAzure](docs/DirectoryProviderSettingsAzure.md)
 - [DirectoryProviderSettingsBase](docs/DirectoryProviderSettingsBase.md)
 - [DirectoryProviderSettingsGApps](docs/DirectoryProviderSettingsGApps.md)
 - [DirectoryProviderSettingsLdap](docs/DirectoryProviderSettingsLdap.md)
 - [DirectoryProviderSettingsOkta](docs/DirectoryProviderSettingsOkta.md)
 - [DirectoryProviderSettingsOneLogin](docs/DirectoryProviderSettingsOneLogin.md)
 - [ElementType](docs/ElementType.md)
 - [EntityIdentifier](docs/EntityIdentifier.md)
 - [EnvironmentVariable](docs/EnvironmentVariable.md)
 - [FieldMatch](docs/FieldMatch.md)
 - [ForensicsLogResult](docs/ForensicsLogResult.md)
 - [ForensicsLogSearchResults](docs/ForensicsLogSearchResults.md)
 - [Group](docs/Group.md)
 - [GroupsPage](docs/GroupsPage.md)
 - [HealthMethod](docs/HealthMethod.md)
 - [IdentityProviderType](docs/IdentityProviderType.md)
 - [KerberosConfiguration](docs/KerberosConfiguration.md)
 - [LogGeoIp](docs/LogGeoIp.md)
 - [LogLocation](docs/LogLocation.md)
 - [LogSearch](docs/LogSearch.md)
 - [LogUserAgentFull](docs/LogUserAgentFull.md)
 - [ModelApiResponse](docs/ModelApiResponse.md)
 - [ModelType](docs/ModelType.md)
 - [NullHandling](docs/NullHandling.md)
 - [PaginatedResponseAdvanced](docs/PaginatedResponseAdvanced.md)
 - [PaginatedResponseBase](docs/PaginatedResponseBase.md)
 - [PoliciesPage](docs/PoliciesPage.md)
 - [PolicyAccess](docs/PolicyAccess.md)
 - [PolicyActionType](docs/PolicyActionType.md)
 - [PolicyActivity](docs/PolicyActivity.md)
 - [PolicyCondition](docs/PolicyCondition.md)
 - [PolicyConditionDefinition](docs/PolicyConditionDefinition.md)
 - [PolicyConditionParameter](docs/PolicyConditionParameter.md)
 - [PolicyConditionParameterEnumSettings](docs/PolicyConditionParameterEnumSettings.md)
 - [PolicyConditionParameterNumberSettings](docs/PolicyConditionParameterNumberSettings.md)
 - [PolicyConditionParameterStringSettings](docs/PolicyConditionParameterStringSettings.md)
 - [PolicyConditionTypeMapping](docs/PolicyConditionTypeMapping.md)
 - [PolicyCore](docs/PolicyCore.md)
 - [PolicyRdpSettings](docs/PolicyRdpSettings.md)
 - [PolicyRule](docs/PolicyRule.md)
 - [PolicySshSettings](docs/PolicySshSettings.md)
 - [PolicyTargetProtocol](docs/PolicyTargetProtocol.md)
 - [PolicyTcpSettings](docs/PolicyTcpSettings.md)
 - [PolicyTimeSettings](docs/PolicyTimeSettings.md)
 - [PolicyType](docs/PolicyType.md)
 - [PolicyUsage](docs/PolicyUsage.md)
 - [PolicyValidatorType](docs/PolicyValidatorType.md)
 - [PortRanges](docs/PortRanges.md)
 - [Property](docs/Property.md)
 - [Schema](docs/Schema.md)
 - [Schema1](docs/Schema1.md)
 - [Schema2](docs/Schema2.md)
 - [Schema3](docs/Schema3.md)
 - [Schema4](docs/Schema4.md)
 - [ScimApiErrorResponse](docs/ScimApiErrorResponse.md)
 - [ScimEntityPaginatedResponseBase](docs/ScimEntityPaginatedResponseBase.md)
 - [ScimGroup](docs/ScimGroup.md)
 - [ScimGroupMembers](docs/ScimGroupMembers.md)
 - [ScimGroupOperation](docs/ScimGroupOperation.md)
 - [ScimGroupPatch](docs/ScimGroupPatch.md)
 - [ScimGroupsPage](docs/ScimGroupsPage.md)
 - [ScimUser](docs/ScimUser.md)
 - [ScimUserEmail](docs/ScimUserEmail.md)
 - [ScimUserFullName](docs/ScimUserFullName.md)
 - [ScimUsersPage](docs/ScimUsersPage.md)
 - [SearchQuery](docs/SearchQuery.md)
 - [SecurityRole](docs/SecurityRole.md)
 - [SharedObject](docs/SharedObject.md)
 - [SharedObjectValue](docs/SharedObjectValue.md)
 - [Site](docs/Site.md)
 - [SiteSettings](docs/SiteSettings.md)
 - [SiteStatus](docs/SiteStatus.md)
 - [SitesPage](docs/SitesPage.md)
 - [Sort](docs/Sort.md)
 - [SshClient](docs/SshClient.md)
 - [SshClientsPage](docs/SshClientsPage.md)
 - [SshUserAccount](docs/SshUserAccount.md)
 - [SshUserAccountStrategy](docs/SshUserAccountStrategy.md)
 - [Status](docs/Status.md)
 - [Status1](docs/Status1.md)
 - [Status2](docs/Status2.md)
 - [TimeSettingsType](docs/TimeSettingsType.md)
 - [UpdateStatus](docs/UpdateStatus.md)
 - [User](docs/User.md)
 - [UserBase](docs/UserBase.md)
 - [UsersPage](docs/UsersPage.md)
 - [Value](docs/Value.md)
 - [Application](docs/Application.md)
 - [DirectoryProviderAdfs](docs/DirectoryProviderAdfs.md)
 - [DirectoryProviderAzure](docs/DirectoryProviderAzure.md)
 - [DirectoryProviderAzureAd](docs/DirectoryProviderAzureAd.md)
 - [DirectoryProviderGApps](docs/DirectoryProviderGApps.md)
 - [DirectoryProviderLdap](docs/DirectoryProviderLdap.md)
 - [DirectoryProviderOkta](docs/DirectoryProviderOkta.md)
 - [DirectoryProviderOneLogin](docs/DirectoryProviderOneLogin.md)


## Documentation For Authorization

## OAuth
- **Type**: OAuth
- **Flow**: application
- **Authorization URL**: 
- **Scopes**: N/A

Example
```golang
auth := context.WithValue(context.Background(), sw.ContextAccessToken, "ACCESSTOKENSTRING")
r, err := client.Service.Operation(auth, args)
```

Or via OAuth2 module to automatically refresh tokens and perform user authentication.
```golang
import "golang.org/x/oauth2"

/* Perform OAuth2 round trip request and obtain a token */

tokenSource := oauth2cfg.TokenSource(createContext(httpClient), &token)
auth := context.WithValue(oauth2.NoContext, sw.ContextOAuth2, tokenSource)
r, err := client.Service.Operation(auth, args)
```

## Author




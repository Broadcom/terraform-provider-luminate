# Go API client for swagger

 ## What's New  <table>   <thead>     <tr>       <th scope=\"col\">Effective from</th>       <th scope=\"col\">Change</th>     </tr>   </thead>   <tbody>     <tr>       <td>March 24th, 2025</th>       <td>         You can now choose the authentication mode for registering Connectors</br>         by specifying the <code>authentication_mode</code> field during site creation.</br>         The available options are: <code>['Connector', 'Site']</code>.</br>         This value is immutable and must be set at creation.</br>         <h4>Default Behavior (<code>Connector</code> Mode)</h4>         By default, the Connector mode is used, which aligns with the previous behavior.         In this mode:</br>         • Connectors are created and bound via API calls.</br>         • During Connector creation, the API returns a one-time password (OTP).</br>         • The OTP is passed as an environment variable and used by the Connector Container for registration.</br>         • Once the registration is complete, the OTP becomes invalid, and the Connector’s persistent storage must be maintained to ensure resiliency.</br>         <h4>New Behavior (<code>Site</code> Mode)</h4>         The Site mode is primarily designed for managed container orchestrators (e.g., Kubernetes, Fargate),<br />         but is also compatible with environments where Connector mode is currently being used.         In this mode:</br>         • A registration key acts as a long-lived token associated with the site.<br />         • The key is shown only once upon creation, and it is the user’s responsibility to store it securely in a secret manager.<br />         • The token is reusable and allows the creation of new Connector entities upon registration.<br />         • The Connector Container uses this token, passed as an environment variable, to register itself.<br />         • Resiliency is ensured as the container dynamically handles the Connector creation and deletion<br />         • Persistent storage is not required, and direct Connector creation via API is disabled in this mode.<br />         <br />New APIs for Managing Site Registration Keys:<br />           <a href=\"#operation/Get Site Registration Keys\">Get site registration keys</a><br />           <a href=\"#operation/Rotate Site Registration Keys\">Rotate site registration key</a><br />           <a href=\"#operation/Delete Site Registration Keys\">Delete site registration keys</a><br />            </td>     </tr>     <tr>       <td>Dec 18th, 2024</th>       <td>The APIs <a href=\"#operation/Get an Application\">Get Application</a>, <a href=\"#operation/Create an Application\">Create Application</a> and <a href=\"#operation/Update an Application\">Update Application</a> have been updated by removing DNS type from both response and request, You can now perform all necessary operations in the new DNS Resiliency section <a href=\"#tag/DNS%20Resiliency\">DNS Resiliency</a>. </td>     </tr>   </tbody> </table>   ## Introduction  Symantec ZTNA API uses common RESTful resourced based URL conventions and JSON as the exchange format. <br> Properties names are case-sensitive. <br> Some of Symantec ZTNA API calls omit None values from the API response.  The base-URL is `api.`&lt;`tenant-name`&gt;`.luminatesec.com`. For example, if your administration portal URL is _admin.acme.luminatesec.com_, then your API base-URL is _api.acme.luminatesec.com_.  All examples below are performed on a tenant called acme.  ## Common Operations Steps  Below you may find a list of common operations and the relevant API calls for each. Each of these operations can also be performed by using the administrative portal at https://admin.acme.luminatesec.com.  <ol>   <li>     Creating a site and deploying a connector:     <ol type=\"a\">       <li>Creating a new site using the <a href=\"#operation/Create a Site\">Create site API</a>.<br></li>       <li>         Once a site is created you can use its Id (returned in the response of the Create Site request)         and call the <a href=\"#operation/Create a Connector\">Create connector API</a>. <br>       </li>       <li>         Deploy the Symantec ZTNA connector:         <ol type=\"i\">           <li>Retrieve the deployment command using the <a href=\"#operation/Get the Connector Deployment Command\">Connector Deployment Command API.</a> <br> </li>           <li>Execute the command on the target machine.</li>         </ol>       </li>     </ol>   </li>   <li>     Creating an application:       <ol type=\"a\">         <li>           An application is always associated with a specific site for routing the traffic to the application via the connectors associated with the same site.           In order to create the application, call the <a href=\"#operation/Create an Application\">Create Application API</a>         </li>         <li>           Once the application is created, you *must* assign the application to a specific site in order to make it accessible. Assign the application to the required site           using the <a href=\"#operation/Bind an Application to a Site\">Bind Application to Site API</a>.         </li>         <li>           In order to grant access to the application for specific entities (users/groups), you should assign the application to the access policy using the <a href=\"#tag/Access-and-Activity-Policies\">Access and Activity Policy API</a>         </li>       </ol>   </li> </ol>  ## Object Model The object model of the API is built around the following: <ol>   <li><a href=\"#tag/Sites\">Sites</a> - Site is a representation of the physical or virtual data center your applications reside in.</li>   <li><a href=\"#tag/Connectors\">Connectors</a> - A connector is a lightweight piece of software connecting your site to the Symantec ZTNA platform.</li>   <li><a href=\"#tag/Applications\">Applications</a>  - Application is the internal resource you would like to publish using Symantec ZTNA. </li>   <li>     <a href=\"#tag/Access-and-Activity-Policies\">Access and Activity Policies</a> - Symantec ZTNA continuously authorize each user request for the contextual access and activity,     in order to control access to resources and restrict user’s actions within resources, based on the user/device context (such as the user’s group membership, user’s location,     MFA status and managed/unmanaged device status) and the requested resource.   <li>     <a href=\"#tag/Cloud-Integration\">Cloud Integration</a> - Integration with Cloud Providers like Amazon Web Services to provide a smoother and cloud-native integration with SIEM solutions      and to allow access to resources based on their associated tags.   <li>     Logs - Symantec ZTNA internal logs for audit and forensics purposes:     <ol>       <li><a href=\"#tag/Audit-Logs\">Audit Logs</a> audit all operations done through the administration portal</li>       <li><a href=\"#tag/Forensics-Logs\">Forensics Logs</a> audit any user's access to any application as well as user's activity for any application.</li>     </ol>   </li> </ol>   ## Authentication  Authentication is done using [OAuth2](https://tools.ietf.org/html/rfc6749) with the [Bearer authentication scheme](https://tools.ietf.org/html/rfc6750).  <!-- ReDoc-Inject: <security-definitions> -->  The Symantec ZTNA API is available to Symantec ZTNA users who have administrative privileges in their Symantec ZTNA tenant. An administrator should create an API client through the Symantec ZTNA Admin portal and copy the ‘Client Id’ and the ‘Client Secret’. Then the administrator should assign the API client an appropriate role in 'Tenant Roles' page.  Retrieving the API access token is done using Basic-Authentication scheme, POST of a Base64 encoded Client-ID and Client-Secret: <B>   ``` curl -X POST \\  https://api.acme.luminatesec.com/v1/oauth/token \\  -u yourApiClientId:yourApiClientSecret   ``` </B>  This call returns the following JSON: {     \"access_token\":\"edfe22e3-eb4c-4c83-8ce3-3152e6a2XXX\",   \"expires_in\":3600,   \"scope\":\"luminate-scope\",   \"token_type\":\"Bearer\",   \"error\":\"\",   \"error_description\":\"\"}  All further API calls should include the ‘Authorization’ header with value “Bearer AccessToken”  For example: <B>   ```   curl -H \"Authorization: Bearer edfe22e3-eb4c-4c83-8ce3-3152e6a2XXX\" \"https://api.acme.luminatesec.com/v2/applications\"   ``` </B>   <br> <br> <br>  <h2>Token</h2>     The Symantec ZTNA API is available to Symantec ZTNA users who have administrative privileges in their Symantec ZTNA tenant. An administrator can create a token through the Symantec ZTNA Admin portal, with type 'Token' and copy the ‘Client Token'.<br> **Make sure to copy the token once it's generated, it won't be presented again!**<br>  Then, the administrator should assign the token an appropriate role in 'Tenant Roles' page. <br> To enforce the new role, the administrator must click the 'Enforce Roles' button on the token entity page. <br>  All further API calls should include the ‘Authorization’ header with value “Bearer &lt;client-token-value&gt;”  For example: <B>   ```   curl -H \"Authorization: Bearer 2d902c58adc30cfaa12c83b7a524fdf753b7f6adb9d13a5d8e4122f643c2c6ca8c597ade37d198e799b338d2448018dddecd341b770d495305d93f39e53e6f6ab691456e25768e76ac7e53282b2b7a24a9d2f3636acca2dd894b1279e4e93aa4db010a0922f612f7253af71b3a414ad5435489bc7987cd6648eeb05ec6643e00ba33c920e2eea9d55cbb4167bc5d58aecfbc0acfc82be613e1176894fd6ded83943374e30fa5724d4b4088494f59eefe2164c6c6373163029b551c195c9251c633ccda6a498a0d48a43fXXX01bd88da0edc5dba6XXXc86d828256d0b8de0af11e3e1e81cf1d1651657c88af8cab57358886d680f36de53608f48b0a9b769aXXX\" \"https://api.acme.luminatesec.com/v2/sites\"   ``` </B>  ## Versioning and Compatibility  The latest Major Version is `v2`.  The Major Version is included in the URL path (e.g. /v2/applications ) and it denotes breaking changes to the API. Minor and Patch versions are transparent to the client.  ## Pagination   Some of our API responses are paginated, meaning that only a certain number of items are returned at a time.  The default number of items returned in a single page is 50.  You can override this by passing a size parameter to set the maximum number of results, but cannot exceed 100.  Specifying the page number sets the starting point for the result set, allowing you to fetch subsequent items  that are not in the initial set of results. The sort order for returned data can be controlled using the sort parameter.<br>  You can constrain the results by using a filter. <br><br>  **Note:** Most methods that support pagination use the approach specified above. However, some methods use varied   versions of pagination. The individual documentation for each API method is your source of truth for which pattern the method follows.  ## Auditing  All authentication operations and modify operations (POST, PUT, DELETE) are audited.   ## Rate-limiting The API has a rate limit of 5 requests per second. If you have hit the rate limit, then a ‘429’ status code will be returned. In such cases, you should back-off from submitting new requests for 1 second before resuming.  Note that rate-limitation applies to the accumulated requests of **all** of your clients. For example, if you have 6 clients submitting requests simultaneously at a rate of 1 request per second for each one then one of them is likely to get a 429 status code.  ## Support  For additional help you may refer to our support at https://support.broadcom.com  Each request submitted to the API returns a unique request ID that is generated by the API. The request ID will be returned in header `x-lum-request-id`. If you need to contact us about any specific request then this ID will serve as a reference to the given request. 

## Overview
This API client was generated by the [swagger-codegen](https://github.com/swagger-api/swagger-codegen) project.  By using the [swagger-spec](https://github.com/swagger-api/swagger-spec) from a remote server, you can easily generate an API client.

- API version: V2
- Package version: 1.0.0
- Build package: io.swagger.codegen.v3.generators.go.GoClientCodegen

## Installation
Put the package under your project folder and add the following in import:
```golang
import "./swagger"
```

## Documentation for API Endpoints

All URIs are relative to *https://api.acme.luminatesec.com/v2*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*AccessAndActivityPoliciesApi* | [**AssignAnApplicationToPolicies**](docs/AccessAndActivityPoliciesApi.md#assignanapplicationtopolicies) | **Post** /policies/by-app-id/{application-id} | Assign Application to policies
*AccessAndActivityPoliciesApi* | [**CreateAnAccessOrActivityPolicy**](docs/AccessAndActivityPoliciesApi.md#createanaccessoractivitypolicy) | **Post** /policies | Create Policy
*AccessAndActivityPoliciesApi* | [**DeleteAPolicy**](docs/AccessAndActivityPoliciesApi.md#deleteapolicy) | **Delete** /policies/{policy-id} | Delete Policy
*AccessAndActivityPoliciesApi* | [**GetAPolicy**](docs/AccessAndActivityPoliciesApi.md#getapolicy) | **Get** /policies/{policy-id} | Get Policy
*AccessAndActivityPoliciesApi* | [**GetAllPolicies**](docs/AccessAndActivityPoliciesApi.md#getallpolicies) | **Get** /policies | List Policies
*AccessAndActivityPoliciesApi* | [**GetPoliciesAssignedToAnApplication**](docs/AccessAndActivityPoliciesApi.md#getpoliciesassignedtoanapplication) | **Get** /policies/by-app-id/{application-id} | Get Application Assigned Policies.
*AccessAndActivityPoliciesApi* | [**GetSupportedActions**](docs/AccessAndActivityPoliciesApi.md#getsupportedactions) | **Get** /policies/config/action-types | Get Supported Rules Actions
*AccessAndActivityPoliciesApi* | [**GetSupportedConditions**](docs/AccessAndActivityPoliciesApi.md#getsupportedconditions) | **Get** /policies/config/condition-definitions | Get Supported Conditions Definitions
*AccessAndActivityPoliciesApi* | [**GetSupportedValidators**](docs/AccessAndActivityPoliciesApi.md#getsupportedvalidators) | **Get** /policies/config/validator-types | Get Supported Validators
*AccessAndActivityPoliciesApi* | [**RemoveAnApplicationFromPolicies**](docs/AccessAndActivityPoliciesApi.md#removeanapplicationfrompolicies) | **Delete** /policies/by-app-id/{application-id} | Remove application from policies
*AccessAndActivityPoliciesApi* | [**UpdateAPolicy**](docs/AccessAndActivityPoliciesApi.md#updateapolicy) | **Put** /policies/{policy-id} | Update Policy
*AccessAndActivityPoliciesApi* | [**UpdateApplicationsAssignedToPolicies**](docs/AccessAndActivityPoliciesApi.md#updateapplicationsassignedtopolicies) | **Put** /policies/by-app-id/{application-id} | Update application in policies
*ApplicationsApi* | [**BindAnApplicationToASite**](docs/ApplicationsApi.md#bindanapplicationtoasite) | **Put** /applications/{application-id}/site-binding/{site-id} | Bind Application to Site
*ApplicationsApi* | [**CreateAnApplication**](docs/ApplicationsApi.md#createanapplication) | **Post** /applications | Create Application
*ApplicationsApi* | [**DeleteAnApplication**](docs/ApplicationsApi.md#deleteanapplication) | **Delete** /applications/{application-id} | Delete Application
*ApplicationsApi* | [**GetAllApplications**](docs/ApplicationsApi.md#getallapplications) | **Get** /applications | List Applications
*ApplicationsApi* | [**GetAnApplication**](docs/ApplicationsApi.md#getanapplication) | **Get** /applications/{application-id} | Get Application
*ApplicationsApi* | [**GetTheApplicationHealth**](docs/ApplicationsApi.md#gettheapplicationhealth) | **Get** /applications/{application-id}/status | Get Application Health Status
*ApplicationsApi* | [**UpdateAnApplication**](docs/ApplicationsApi.md#updateanapplication) | **Put** /applications/{application-id} | Update Application
*AuditLogsApi* | [**SearchAuditLogs**](docs/AuditLogsApi.md#searchauditlogs) | **Post** /logs/audit | Search Audit logs
*CloudIntegrationApi* | [**CreateACloudIntegration**](docs/CloudIntegrationApi.md#createacloudintegration) | **Post** /cloud-integrations/integrations | Create Cloud Integration Configuration
*CloudIntegrationApi* | [**DeleteACloudIntegration**](docs/CloudIntegrationApi.md#deleteacloudintegration) | **Delete** /cloud-integrations/integrations/{cloud-integration-id} | Delete Cloud Integration
*CloudIntegrationApi* | [**GetCloudIntegration**](docs/CloudIntegrationApi.md#getcloudintegration) | **Get** /cloud-integrations/integrations/{cloud-integration-id} | Get Cloud Integration Configuration
*CloudIntegrationApi* | [**ListCloudIntegrations**](docs/CloudIntegrationApi.md#listcloudintegrations) | **Get** /cloud-integrations/integrations | List Cloud Integration Configurations
*CloudIntegrationApi* | [**UpdateACloudIntegration**](docs/CloudIntegrationApi.md#updateacloudintegration) | **Put** /cloud-integrations/integrations/{cloud-integration-id} | Update Cloud Integration Configuration
*CollectionsApi* | [**CreateACollection**](docs/CollectionsApi.md#createacollection) | **Post** /collection | Create Collection
*CollectionsApi* | [**CreateACollectionRoleBinding**](docs/CollectionsApi.md#createacollectionrolebinding) | **Post** /collection/collection-role-bindings | Create Collection Role Binding
*CollectionsApi* | [**CreateASiteRoleBinding**](docs/CollectionsApi.md#createasiterolebinding) | **Post** /collection/site-role-bindings | Create Site Role Binding
*CollectionsApi* | [**CreateATenantRoleBinding**](docs/CollectionsApi.md#createatenantrolebinding) | **Post** /collection/tenant-role-bindings | Create Tenant Role Binding
*CollectionsApi* | [**DeleteACollection**](docs/CollectionsApi.md#deleteacollection) | **Delete** /collection/{collection-id} | Delete Collection
*CollectionsApi* | [**DeleteARoleBinding**](docs/CollectionsApi.md#deletearolebinding) | **Post** /collection/role-bindings/delete | Delete Role Binding
*CollectionsApi* | [**GetACollection**](docs/CollectionsApi.md#getacollection) | **Get** /collection/{collection-id} | Get Collection
*CollectionsApi* | [**GetCollectionsBySite**](docs/CollectionsApi.md#getcollectionsbysite) | **Get** /collection/site/{site-id} | Get Collections by Site
*CollectionsApi* | [**GetSitesLinkedToTheCollection**](docs/CollectionsApi.md#getsiteslinkedtothecollection) | **Get** /collection/site-links/{collection-id} | Get Collection Linked Sites
*CollectionsApi* | [**LinkASiteToACollection**](docs/CollectionsApi.md#linkasitetoacollection) | **Post** /collection/site-links | Link Site to Collection
*CollectionsApi* | [**ListCollections**](docs/CollectionsApi.md#listcollections) | **Get** /collection | List Collections
*CollectionsApi* | [**ListRoleBindings**](docs/CollectionsApi.md#listrolebindings) | **Get** /collection/role-bindings | List Role Bindings
*CollectionsApi* | [**UnlinkASiteFromACollection**](docs/CollectionsApi.md#unlinkasitefromacollection) | **Delete** /collection/site-links/{collection-id}/{site-id} | Unlinks site from collection
*CollectionsApi* | [**UpdateACollection**](docs/CollectionsApi.md#updateacollection) | **Put** /collection/{collection-id} | Update Collection
*ConnectorsApi* | [**CreateAConnector**](docs/ConnectorsApi.md#createaconnector) | **Post** /connectors | Create Connector
*ConnectorsApi* | [**DeleteAConnector**](docs/ConnectorsApi.md#deleteaconnector) | **Delete** /connectors/{connector-id} | Delete Connector
*ConnectorsApi* | [**GetAConnector**](docs/ConnectorsApi.md#getaconnector) | **Get** /connectors/{connector-id} | Get Connector
*ConnectorsApi* | [**GetAllConnectors**](docs/ConnectorsApi.md#getallconnectors) | **Get** /connectors | List Connectors
*ConnectorsApi* | [**GetEnvironmentVariablesForConnectorDeployment**](docs/ConnectorsApi.md#getenvironmentvariablesforconnectordeployment) | **Get** /connectors/{connector-id}/environment_variables | Get Connector Environment Variables
*ConnectorsApi* | [**GetTheConnectorDeploymentCommand**](docs/ConnectorsApi.md#gettheconnectordeploymentcommand) | **Get** /connectors/{connector-id}/command | Get Connector Deployment Command
*ConnectorsApi* | [**GetTheLatestConnectorVersion**](docs/ConnectorsApi.md#getthelatestconnectorversion) | **Get** /connectors/version | Get Connector Version
*DNSResiliencyApi* | [**CreateADNSGroup**](docs/DNSResiliencyApi.md#createadnsgroup) | **Post** /wss-integration-tenant/dns-groups | Create New DNS Group
*DNSResiliencyApi* | [**CreateADNSServer**](docs/DNSResiliencyApi.md#createadnsserver) | **Post** /wss-integration-tenant/dns-groups/{dnsGroupId}/servers/ | Create New DNS Server in Group
*DNSResiliencyApi* | [**DeleteADNSGroup**](docs/DNSResiliencyApi.md#deleteadnsgroup) | **Delete** /wss-integration-tenant/dns-groups/{dnsGroupId} | Delete a group by ID
*DNSResiliencyApi* | [**DeleteDNSServers**](docs/DNSResiliencyApi.md#deletednsservers) | **Post** /wss-integration-tenant/dns-groups/{dnsGroupId}/servers/delete-by-ids | Delete DNS Servers By Ids
*DNSResiliencyApi* | [**EnableOrDisableDNSGroups**](docs/DNSResiliencyApi.md#enableordisablednsgroups) | **Post** /wss-integration-tenant/dns-groups/enableByIds | Enable/Disable DNS Groups
*DNSResiliencyApi* | [**GetADNSGroup**](docs/DNSResiliencyApi.md#getadnsgroup) | **Get** /wss-integration-tenant/dns-groups/{dnsGroupId} | Get DNS Group By ID
*DNSResiliencyApi* | [**GetADNSServer**](docs/DNSResiliencyApi.md#getadnsserver) | **Get** /wss-integration-tenant/dns-groups/{dnsGroupId}/servers/{serverId} | Get DNS Server By Id
*DNSResiliencyApi* | [**GetAllServersOfADNSGroup**](docs/DNSResiliencyApi.md#getallserversofadnsgroup) | **Get** /wss-integration-tenant/dns-groups/{dnsGroupId}/servers/ | Get All Servers Of a DNS Group
*DNSResiliencyApi* | [**ListDNSGroups**](docs/DNSResiliencyApi.md#listdnsgroups) | **Get** /wss-integration-tenant/dns-groups | List DNS Groups
*DNSResiliencyApi* | [**UpdateADNSGroup**](docs/DNSResiliencyApi.md#updateadnsgroup) | **Put** /wss-integration-tenant/dns-groups/{dnsGroupId} | Update a group by ID
*DNSResiliencyApi* | [**UpdateADNSServer**](docs/DNSResiliencyApi.md#updateadnsserver) | **Put** /wss-integration-tenant/dns-groups/{dnsGroupId}/servers/{serverId} | Update DNS Server By Id
*DNSResiliencyApi* | [**UpdateDNSServerResiliencyOrder**](docs/DNSResiliencyApi.md#updatednsserverresiliencyorder) | **Put** /wss-integration-tenant/dns-groups/{dnsGroupId}/server-order | Update DNS Servers Order
*ForensicsLogsApi* | [**SearchForensicsLogs**](docs/ForensicsLogsApi.md#searchforensicslogs) | **Post** /logs/forensics | Search Forensics logs
*GroupsApi* | [**AssignAUserToAGroup**](docs/GroupsApi.md#assignausertoagroup) | **Put** /identities/local/groups/{group-id}/users/{user-id} | Assign User To Group
*GroupsApi* | [**CreateAGroup**](docs/GroupsApi.md#createagroup) | **Post** /identities/{identity-provider-id}/groups | Create Group
*GroupsApi* | [**DeleteAGroup**](docs/GroupsApi.md#deleteagroup) | **Delete** /identities/{identity-provider-id}/groups/{entity-id} | Delete Group
*GroupsApi* | [**GetAGroup**](docs/GroupsApi.md#getagroup) | **Get** /identities/{identity-provider-id}/groups/{entity-id} | Get Group
*GroupsApi* | [**ListAGroupsAssignedUsers**](docs/GroupsApi.md#listagroupsassignedusers) | **Get** /identities/{identity-provider-id}/groups/{entity-id}/users | List Assigned Users
*GroupsApi* | [**RemoveAUserFromAGroup**](docs/GroupsApi.md#removeauserfromagroup) | **Delete** /identities/local/groups/{group-id}/users/{user-id} | Remove User From Group
*GroupsApi* | [**SearchGroupsByTheIdP**](docs/GroupsApi.md#searchgroupsbytheidp) | **Get** /identities/{identity-provider-id}/groups | Search Groups By Identity Provider
*IdentityProvidersApi* | [**ListIdentityProviders**](docs/IdentityProvidersApi.md#listidentityproviders) | **Get** /identities/settings/identity-providers | List Identity Providers
*SCIMApi* | [**CreateASCIMGroup**](docs/SCIMApi.md#createascimgroup) | **Post** /identities/{identity-provider-id}/scim/groups | Create a SCIM Group
*SCIMApi* | [**CreateASCIMUser**](docs/SCIMApi.md#createascimuser) | **Post** /identities/{identity-provider-id}/scim/users | Create SCIM User
*SCIMApi* | [**DeleteASCIMGroup**](docs/SCIMApi.md#deleteascimgroup) | **Delete** /identities/{identity-provider-id}/scim/groups/{group-id} | Delete SCIM Group
*SCIMApi* | [**DeleteASCIMUser**](docs/SCIMApi.md#deleteascimuser) | **Delete** /identities/{identity-provider-id}/scim/users/{user-id} | Delete SCIM User
*SCIMApi* | [**GetASCIMGroup**](docs/SCIMApi.md#getascimgroup) | **Get** /identities/{identity-provider-id}/scim/groups/{group-id} | Get SCIM Group
*SCIMApi* | [**GetASCIMUser**](docs/SCIMApi.md#getascimuser) | **Get** /identities/{identity-provider-id}/scim/users/{user-id} | Get SCIM User
*SCIMApi* | [**ListSCIMGroups**](docs/SCIMApi.md#listscimgroups) | **Get** /identities/{identity-provider-id}/scim/groups | List SCIM Groups
*SCIMApi* | [**ListSCIMUsers**](docs/SCIMApi.md#listscimusers) | **Get** /identities/{identity-provider-id}/scim/users | List SCIM Users
*SCIMApi* | [**ModifyASCIMGroup**](docs/SCIMApi.md#modifyascimgroup) | **Patch** /identities/{identity-provider-id}/scim/groups/{group-id} | Modify a SCIM Group
*SCIMApi* | [**UpdateASCIMGroup**](docs/SCIMApi.md#updateascimgroup) | **Put** /identities/{identity-provider-id}/scim/groups/{group-id} | Update SCIM Group
*SCIMApi* | [**UpdateASCIMUser**](docs/SCIMApi.md#updateascimuser) | **Put** /identities/{identity-provider-id}/scim/users/{user-id} | Update SCIM User
*SSHClientsApi* | [**GetAllSSHClients**](docs/SSHClientsApi.md#getallsshclients) | **Get** /ssh-clients | List SSH Clients
*SharedObjectsApi* | [**CreateASharedObject**](docs/SharedObjectsApi.md#createasharedobject) | **Post** /policies/shared-objects | Create Shared Object
*SharedObjectsApi* | [**DeleteASharedObject**](docs/SharedObjectsApi.md#deleteasharedobject) | **Delete** /policies/shared-objects/{shared-object-id} | Delete Shared Object
*SharedObjectsApi* | [**GetASharedObject**](docs/SharedObjectsApi.md#getasharedobject) | **Get** /policies/shared-objects/{shared-object-id} | Get Shared Object
*SharedObjectsApi* | [**ListSharedObjects**](docs/SharedObjectsApi.md#listsharedobjects) | **Get** /policies/shared-objects | List Shared Objects
*SharedObjectsApi* | [**UpdateASharedObject**](docs/SharedObjectsApi.md#updateasharedobject) | **Put** /policies/shared-objects/{shared-object-id} | Update Shared Object
*SiteRegistrationKeysApi* | [**DeleteSiteRegistrationKeys**](docs/SiteRegistrationKeysApi.md#deletesiteregistrationkeys) | **Delete** /sites/{site-id}/registration_keys | Clean Site Registration Keys
*SiteRegistrationKeysApi* | [**GetSiteRegistrationKeys**](docs/SiteRegistrationKeysApi.md#getsiteregistrationkeys) | **Get** /sites/{site-id}/registration_keys | List Site Registration Keys
*SiteRegistrationKeysApi* | [**RotateSiteRegistrationKeys**](docs/SiteRegistrationKeysApi.md#rotatesiteregistrationkeys) | **Post** /sites/{site-id}/registration_keys | Rotate Site Registration Keys
*SitesApi* | [**CreateASite**](docs/SitesApi.md#createasite) | **Post** /sites | Create Site
*SitesApi* | [**DeleteASite**](docs/SitesApi.md#deleteasite) | **Delete** /sites/{site-id} | Delete Site
*SitesApi* | [**GetARegion**](docs/SitesApi.md#getaregion) | **Get** /regions/{region_name} | Get Region By Name
*SitesApi* | [**GetASite**](docs/SitesApi.md#getasite) | **Get** /sites/{site-id} | Get Site
*SitesApi* | [**GetASitesHealthStatus**](docs/SitesApi.md#getasiteshealthstatus) | **Get** /sites/{site-id}/status | Get Site Health Status
*SitesApi* | [**GetAllSites**](docs/SitesApi.md#getallsites) | **Get** /sites | List Sites
*SitesApi* | [**ListRegions**](docs/SitesApi.md#listregions) | **Get** /regions | List Regions
*SitesApi* | [**UpdateASite**](docs/SitesApi.md#updateasite) | **Put** /sites/{site-id} | Update Site
*UsersApi* | [**BlockUser**](docs/UsersApi.md#blockuser) | **Post** /identities/{identity-provider-id}/users/{entity-id}/block | Block User
*UsersApi* | [**CreateLocalUser**](docs/UsersApi.md#createlocaluser) | **Post** /identities/local/users | Create Local User
*UsersApi* | [**DeleteLocalUser**](docs/UsersApi.md#deletelocaluser) | **Delete** /identities/local/users/{entity-id} | Delete Local User
*UsersApi* | [**DeleteUser**](docs/UsersApi.md#deleteuser) | **Delete** /identities/{identity-provider-id}/users/{entity-id} | Delete User
*UsersApi* | [**GetUser**](docs/UsersApi.md#getuser) | **Get** /identities/{identity-provider-id}/users/{entity-id} | Get User
*UsersApi* | [**ListBlockedUsers**](docs/UsersApi.md#listblockedusers) | **Get** /identities/settings/blocked-users | List Blocked Users
*UsersApi* | [**SearchUsersByIdP**](docs/UsersApi.md#searchusersbyidp) | **Get** /identities/{identity-provider-id}/users | Search Users By Identity Provider
*UsersApi* | [**UnblockUser**](docs/UsersApi.md#unblockuser) | **Delete** /identities/{identity-provider-id}/users/{entity-id}/block | Unblock User
*UsersApi* | [**UpdateLocalUser**](docs/UsersApi.md#updatelocaluser) | **Put** /identities/local/users/{entity-id} | Update Local User

## Documentation For Models

 - [AnyOfSearchAfterItems](docs/AnyOfSearchAfterItems.md)
 - [Application](docs/Application.md)
 - [ApplicationBase](docs/ApplicationBase.md)
 - [ApplicationByType](docs/ApplicationByType.md)
 - [ApplicationCloudIntegrationData](docs/ApplicationCloudIntegrationData.md)
 - [ApplicationCloudIntegrationDataProperties](docs/ApplicationCloudIntegrationDataProperties.md)
 - [ApplicationCloudIntegrationTag](docs/ApplicationCloudIntegrationTag.md)
 - [ApplicationConnectionSettings](docs/ApplicationConnectionSettings.md)
 - [ApplicationConnectionSettingsRdp](docs/ApplicationConnectionSettingsRdp.md)
 - [ApplicationConnectionSettingsSegment](docs/ApplicationConnectionSettingsSegment.md)
 - [ApplicationConnectionSettingsSegmentData](docs/ApplicationConnectionSettingsSegmentData.md)
 - [ApplicationConnectionSettingsSegmentToDeprecate](docs/ApplicationConnectionSettingsSegmentToDeprecate.md)
 - [ApplicationConnectionSettingsSsh](docs/ApplicationConnectionSettingsSsh.md)
 - [ApplicationConnectionSettingsTcp](docs/ApplicationConnectionSettingsTcp.md)
 - [ApplicationCore](docs/ApplicationCore.md)
 - [ApplicationDynamicSsh](docs/ApplicationDynamicSsh.md)
 - [ApplicationHealth](docs/ApplicationHealth.md)
 - [ApplicationHttp](docs/ApplicationHttp.md)
 - [ApplicationLinkTranslationSettings](docs/ApplicationLinkTranslationSettings.md)
 - [ApplicationRdp](docs/ApplicationRdp.md)
 - [ApplicationRequestCustomizationSettings](docs/ApplicationRequestCustomizationSettings.md)
 - [ApplicationSegment](docs/ApplicationSegment.md)
 - [ApplicationSort](docs/ApplicationSort.md)
 - [ApplicationSsh](docs/ApplicationSsh.md)
 - [ApplicationSubType](docs/ApplicationSubType.md)
 - [ApplicationTcp](docs/ApplicationTcp.md)
 - [ApplicationTcpTarget](docs/ApplicationTcpTarget.md)
 - [ApplicationTcpTargetPortRanges](docs/ApplicationTcpTargetPortRanges.md)
 - [ApplicationTcpTunnelSettings](docs/ApplicationTcpTunnelSettings.md)
 - [ApplicationToPoliciesMapping](docs/ApplicationToPoliciesMapping.md)
 - [ApplicationType](docs/ApplicationType.md)
 - [ApplicationTypeWithRuleIds](docs/ApplicationTypeWithRuleIds.md)
 - [ApplicationVpcData](docs/ApplicationVpcData.md)
 - [ApplicationsApplicationidBody](docs/ApplicationsApplicationidBody.md)
 - [ApplicationsBody](docs/ApplicationsBody.md)
 - [ApplicationsPage](docs/ApplicationsPage.md)
 - [AuditLogResult](docs/AuditLogResult.md)
 - [AuditLogResultData](docs/AuditLogResultData.md)
 - [AuditLogSearchResults](docs/AuditLogSearchResults.md)
 - [BlockedUser](docs/BlockedUser.md)
 - [BulkApiErrorResponse](docs/BulkApiErrorResponse.md)
 - [BulkApiResponse](docs/BulkApiResponse.md)
 - [ByappidApplicationidBody](docs/ByappidApplicationidBody.md)
 - [ByappidApplicationidBody1](docs/ByappidApplicationidBody1.md)
 - [ByappidApplicationidBody2](docs/ByappidApplicationidBody2.md)
 - [CloudIntegration](docs/CloudIntegration.md)
 - [CloudIntegrationBase](docs/CloudIntegrationBase.md)
 - [CloudIntegrationError](docs/CloudIntegrationError.md)
 - [CloudIntegrationHealth](docs/CloudIntegrationHealth.md)
 - [CloudIntegrationPost](docs/CloudIntegrationPost.md)
 - [CloudIntegrationProvider](docs/CloudIntegrationProvider.md)
 - [CloudIntegrationPut](docs/CloudIntegrationPut.md)
 - [CloudintegrationsIntegrationsBody](docs/CloudintegrationsIntegrationsBody.md)
 - [Collection](docs/Collection.md)
 - [CollectionBody](docs/CollectionBody.md)
 - [CollectionCollectionidBody](docs/CollectionCollectionidBody.md)
 - [CollectionCollectionidBody1](docs/CollectionCollectionidBody1.md)
 - [CollectionCollectionrolebindingsBody](docs/CollectionCollectionrolebindingsBody.md)
 - [CollectionIdSiteIdTuple](docs/CollectionIdSiteIdTuple.md)
 - [CollectionIdSiteIdTupleList](docs/CollectionIdSiteIdTupleList.md)
 - [CollectionIds](docs/CollectionIds.md)
 - [CollectionPage](docs/CollectionPage.md)
 - [CollectionRequest](docs/CollectionRequest.md)
 - [CollectionRoleType](docs/CollectionRoleType.md)
 - [CollectionSiteLink](docs/CollectionSiteLink.md)
 - [CollectionSiteLinks](docs/CollectionSiteLinks.md)
 - [CollectionSitelinksBody](docs/CollectionSitelinksBody.md)
 - [CollectionSiterolebindingsBody](docs/CollectionSiterolebindingsBody.md)
 - [CollectionTenantrolebindingsBody](docs/CollectionTenantrolebindingsBody.md)
 - [CollectionUpdateRequest](docs/CollectionUpdateRequest.md)
 - [CollectionidSiteidBody](docs/CollectionidSiteidBody.md)
 - [Connector](docs/Connector.md)
 - [ConnectorDeploymentCommand](docs/ConnectorDeploymentCommand.md)
 - [ConnectorEnvironmentVariables](docs/ConnectorEnvironmentVariables.md)
 - [ConnectorEnvironmentVariablesEnvironmentVariables](docs/ConnectorEnvironmentVariablesEnvironmentVariables.md)
 - [ConnectorLastSeen](docs/ConnectorLastSeen.md)
 - [ConnectorVersion](docs/ConnectorVersion.md)
 - [ConnectorsBody](docs/ConnectorsBody.md)
 - [ConnectorsPage](docs/ConnectorsPage.md)
 - [DeploymentType](docs/DeploymentType.md)
 - [DirectoryEntity](docs/DirectoryEntity.md)
 - [DirectoryEntityBinding](docs/DirectoryEntityBinding.md)
 - [DirectoryEntityPaginatedResponseBase](docs/DirectoryEntityPaginatedResponseBase.md)
 - [DirectoryProvider](docs/DirectoryProvider.md)
 - [DirectoryProviderAdfs](docs/DirectoryProviderAdfs.md)
 - [DirectoryProviderAzure](docs/DirectoryProviderAzure.md)
 - [DirectoryProviderAzureAd](docs/DirectoryProviderAzureAd.md)
 - [DirectoryProviderGApps](docs/DirectoryProviderGApps.md)
 - [DirectoryProviderHealth](docs/DirectoryProviderHealth.md)
 - [DirectoryProviderInstructionsLdap](docs/DirectoryProviderInstructionsLdap.md)
 - [DirectoryProviderInstructionsOneLoginOrOkta](docs/DirectoryProviderInstructionsOneLoginOrOkta.md)
 - [DirectoryProviderLdap](docs/DirectoryProviderLdap.md)
 - [DirectoryProviderNoInstructions](docs/DirectoryProviderNoInstructions.md)
 - [DirectoryProviderNoSettings](docs/DirectoryProviderNoSettings.md)
 - [DirectoryProviderOkta](docs/DirectoryProviderOkta.md)
 - [DirectoryProviderOneLogin](docs/DirectoryProviderOneLogin.md)
 - [DirectoryProviderSettingsAzure](docs/DirectoryProviderSettingsAzure.md)
 - [DirectoryProviderSettingsBase](docs/DirectoryProviderSettingsBase.md)
 - [DirectoryProviderSettingsGApps](docs/DirectoryProviderSettingsGApps.md)
 - [DirectoryProviderSettingsLdap](docs/DirectoryProviderSettingsLdap.md)
 - [DirectoryProviderSettingsOkta](docs/DirectoryProviderSettingsOkta.md)
 - [DirectoryProviderSettingsOneLogin](docs/DirectoryProviderSettingsOneLogin.md)
 - [DnsGroupIdServerorderBody](docs/DnsGroupIdServerorderBody.md)
 - [DnsGroupIdServersBody](docs/DnsGroupIdServersBody.md)
 - [DnsGroupInput](docs/DnsGroupInput.md)
 - [DnsGroupOutput](docs/DnsGroupOutput.md)
 - [DnsGroupOutputPage](docs/DnsGroupOutputPage.md)
 - [DnsServerIds](docs/DnsServerIds.md)
 - [DnsServerInput](docs/DnsServerInput.md)
 - [DnsServerOutput](docs/DnsServerOutput.md)
 - [DnsgroupsDnsGroupIdBody](docs/DnsgroupsDnsGroupIdBody.md)
 - [DnsgroupsEnableByIdsBody](docs/DnsgroupsEnableByIdsBody.md)
 - [EntityIdentifier](docs/EntityIdentifier.md)
 - [EntityType](docs/EntityType.md)
 - [FieldMatch](docs/FieldMatch.md)
 - [ForensicsLogResult](docs/ForensicsLogResult.md)
 - [ForensicsLogResultData](docs/ForensicsLogResultData.md)
 - [ForensicsLogSearchResults](docs/ForensicsLogSearchResults.md)
 - [Group](docs/Group.md)
 - [GroupCreateRequest](docs/GroupCreateRequest.md)
 - [GroupsGroupidBody](docs/GroupsGroupidBody.md)
 - [GroupsGroupidBody1](docs/GroupsGroupidBody1.md)
 - [GroupsPage](docs/GroupsPage.md)
 - [IdentityProviderType](docs/IdentityProviderType.md)
 - [IntegrationsCloudintegrationidBody](docs/IntegrationsCloudintegrationidBody.md)
 - [LogGeoIp](docs/LogGeoIp.md)
 - [LogLocation](docs/LogLocation.md)
 - [LogSearch](docs/LogSearch.md)
 - [LogUserAgentFull](docs/LogUserAgentFull.md)
 - [LogsAuditBody](docs/LogsAuditBody.md)
 - [LogsForensicsBody](docs/LogsForensicsBody.md)
 - [ModelApiResponse](docs/ModelApiResponse.md)
 - [PageOffset](docs/PageOffset.md)
 - [PaginatedResponseAdvanced](docs/PaginatedResponseAdvanced.md)
 - [PaginatedResponseBase](docs/PaginatedResponseBase.md)
 - [PoliciesBody](docs/PoliciesBody.md)
 - [PoliciesPage](docs/PoliciesPage.md)
 - [PoliciesPolicyidBody](docs/PoliciesPolicyidBody.md)
 - [PoliciesSharedobjectsBody](docs/PoliciesSharedobjectsBody.md)
 - [Policy](docs/Policy.md)
 - [PolicyAccess](docs/PolicyAccess.md)
 - [PolicyActionType](docs/PolicyActionType.md)
 - [PolicyActivity](docs/PolicyActivity.md)
 - [PolicyByType](docs/PolicyByType.md)
 - [PolicyCondition](docs/PolicyCondition.md)
 - [PolicyConditionDefinition](docs/PolicyConditionDefinition.md)
 - [PolicyConditionParameter](docs/PolicyConditionParameter.md)
 - [PolicyConditionParameterEnumSettings](docs/PolicyConditionParameterEnumSettings.md)
 - [PolicyConditionParameterEnumSettingsValues](docs/PolicyConditionParameterEnumSettingsValues.md)
 - [PolicyConditionParameterNumberSettings](docs/PolicyConditionParameterNumberSettings.md)
 - [PolicyConditionParameterStringSettings](docs/PolicyConditionParameterStringSettings.md)
 - [PolicyConditionTypeMapping](docs/PolicyConditionTypeMapping.md)
 - [PolicyCore](docs/PolicyCore.md)
 - [PolicyRdpSettings](docs/PolicyRdpSettings.md)
 - [PolicyRule](docs/PolicyRule.md)
 - [PolicySshSettings](docs/PolicySshSettings.md)
 - [PolicyTargetProtocol](docs/PolicyTargetProtocol.md)
 - [PolicyTargetProtocolSubType](docs/PolicyTargetProtocolSubType.md)
 - [PolicyTcpSettings](docs/PolicyTcpSettings.md)
 - [PolicyTimeSettings](docs/PolicyTimeSettings.md)
 - [PolicyType](docs/PolicyType.md)
 - [PolicyUsage](docs/PolicyUsage.md)
 - [PolicyValidatorType](docs/PolicyValidatorType.md)
 - [Region](docs/Region.md)
 - [Resource](docs/Resource.md)
 - [Role](docs/Role.md)
 - [RoleBinding](docs/RoleBinding.md)
 - [RoleBindings](docs/RoleBindings.md)
 - [RoleBindingsCreateRequestBase](docs/RoleBindingsCreateRequestBase.md)
 - [RoleBindingsDeleteRequest](docs/RoleBindingsDeleteRequest.md)
 - [RoleBindingsPage](docs/RoleBindingsPage.md)
 - [RoleType](docs/RoleType.md)
 - [RolebindingsDeleteBody](docs/RolebindingsDeleteBody.md)
 - [RotateKeyRequestPostBody](docs/RotateKeyRequestPostBody.md)
 - [RotateKeyResponse](docs/RotateKeyResponse.md)
 - [ScimApiErrorResponse](docs/ScimApiErrorResponse.md)
 - [ScimEntityPaginatedResponseBase](docs/ScimEntityPaginatedResponseBase.md)
 - [ScimGroup](docs/ScimGroup.md)
 - [ScimGroupMembers](docs/ScimGroupMembers.md)
 - [ScimGroupOperation](docs/ScimGroupOperation.md)
 - [ScimGroupPatch](docs/ScimGroupPatch.md)
 - [ScimGroupsBody](docs/ScimGroupsBody.md)
 - [ScimGroupsPage](docs/ScimGroupsPage.md)
 - [ScimUser](docs/ScimUser.md)
 - [ScimUserEmail](docs/ScimUserEmail.md)
 - [ScimUserFullName](docs/ScimUserFullName.md)
 - [ScimUsersBody](docs/ScimUsersBody.md)
 - [ScimUsersPage](docs/ScimUsersPage.md)
 - [SearchQuery](docs/SearchQuery.md)
 - [SecurityRole](docs/SecurityRole.md)
 - [ServersDeletebyidsBody](docs/ServersDeletebyidsBody.md)
 - [ServersServerIdBody](docs/ServersServerIdBody.md)
 - [SharedObject](docs/SharedObject.md)
 - [SharedObjectValue](docs/SharedObjectValue.md)
 - [SharedObjectsPage](docs/SharedObjectsPage.md)
 - [SharedobjectsSharedobjectidBody](docs/SharedobjectsSharedobjectidBody.md)
 - [Site](docs/Site.md)
 - [SiteAuthenticationMode](docs/SiteAuthenticationMode.md)
 - [SiteRegistrationKey](docs/SiteRegistrationKey.md)
 - [SiteRegistrationKeysResponse](docs/SiteRegistrationKeysResponse.md)
 - [SiteRoleType](docs/SiteRoleType.md)
 - [SiteSettings](docs/SiteSettings.md)
 - [SiteStatus](docs/SiteStatus.md)
 - [SitebindingSiteidBody](docs/SitebindingSiteidBody.md)
 - [SiteidRegistrationKeysBody](docs/SiteidRegistrationKeysBody.md)
 - [SitesBody](docs/SitesBody.md)
 - [SitesPage](docs/SitesPage.md)
 - [SitesSiteidBody](docs/SitesSiteidBody.md)
 - [SshClient](docs/SshClient.md)
 - [SshClientsPage](docs/SshClientsPage.md)
 - [SshUserAccount](docs/SshUserAccount.md)
 - [SubjectType](docs/SubjectType.md)
 - [TenantRoleType](docs/TenantRoleType.md)
 - [ToggleDnsGroupIdsRequest](docs/ToggleDnsGroupIdsRequest.md)
 - [User](docs/User.md)
 - [UserBase](docs/UserBase.md)
 - [UsersPage](docs/UsersPage.md)
 - [UsersUseridBody](docs/UsersUseridBody.md)
 - [WebRdpSettings](docs/WebRdpSettings.md)
 - [WorkingHoursSettings](docs/WorkingHoursSettings.md)
 - [WssintegrationtenantDnsgroupsBody](docs/WssintegrationtenantDnsgroupsBody.md)

## Documentation For Authorization

## OAuth
- **Type**: OAuth
- **Flow**: application
- **Authorization URL**: 
- **Scopes**: 

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



//go:generate mockery --inpackage --name=Api --structname=MockManage

package manage

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

type Api interface {
	GetAdminHealth() (*integration.MsxResponse, error)

	GetSubscription(subscriptionId string) (*integration.MsxResponse, error)
	GetSubscriptionsV3(serviceType string, page, pageSize int) (*integration.MsxResponse, error)
	CreateSubscription(tenantId, serviceType string, subscriptionName *string,
		subscriptionAttribute, offerDefAttribute, offerSelectionDetail, costAttribute map[string]string) (*integration.MsxResponse, error)
	UpdateSubscription(subscriptionId, serviceType string, subscriptionName *string,
		subscriptionAttribute, offerDefAttribute, offerSelectionDetail, costAttribute map[string]string) (*integration.MsxResponse, error)
	DeleteSubscription(subscriptionId string) (*integration.MsxResponse, error)

	// CreateServiceOrder
	// UpdateServiceOrder

	GetServiceInstance(serviceInstanceId string) (*integration.MsxResponse, error)
	GetSubscriptionServiceInstances(subscriptionId string, page, pageSize int) (*integration.MsxResponse, error)
	CreateServiceInstance(subscriptionId, serviceInstanceId string, serviceAttribute, serviceDefAttribute, status map[string]string) (*integration.MsxResponse, error)
	UpdateServiceInstance(serviceInstanceId string, serviceAttribute, serviceDefAttribute, status map[string]string) (*integration.MsxResponse, error)
	DeleteServiceInstance(serviceInstanceId string) (*integration.MsxResponse, error)

	//Deprecated: Use v3 Endpoint Instead
	GetSite(siteId string) (*integration.MsxResponse, error)
	//Deprecated: Use v3 Endpoint Instead
	CreateSite(subscriptionId, serviceInstanceId string, siteId, siteName, siteType, displayName *string, siteAttributes, siteDefAttributes map[string]string, devices []string) (*integration.MsxResponse, error)
	//Deprecated: Use v3 Endpoint Instead
	UpdateSite(siteId string, siteType, displayName *string, siteAttributes, siteDefAttributes map[string]string, devices []string) (*integration.MsxResponse, error)
	//Deprecated: Use v3 Endpoint Instead
	DeleteSite(siteId string) (*integration.MsxResponse, error)

	GetSitesV3(siteFilters SiteQueryFilter, page, pageSize int) (*integration.MsxResponse, error)
	GetSiteV3(siteId string, showImage string) (*integration.MsxResponse, error)
	CreateSiteV3(siteRequest SiteCreateRequest) (*integration.MsxResponse, error)
	UpdateSiteV3(siteRequest SiteUpdateRequest, siteId string, notification string) (*integration.MsxResponse, error)
	DeleteSiteV3(siteId string) (*integration.MsxResponse, error)
	AddDeviceToSiteV3(deviceId string, siteId string, notification string) (*integration.MsxResponse, error)
	DeleteDeviceFromSiteV3(deviceId string, siteId string) (*integration.MsxResponse, error)
	UpdateSiteStatusV3(siteStatus SiteStatusUpdateRequest, siteId string) (*integration.MsxResponse, error)

	//Deprecated: User v4 Endpoint Instead
	CreateManagedDevice(tenantId string, deviceModel, deviceOnboardType string, deviceOnboardInfo map[string]string) (*integration.MsxResponse, error)
	//Deprecated: User v4 Endpoint Instead
	DeleteManagedDevice(deviceInstanceId string) (*integration.MsxResponse, error)
	GetDeviceConfig(deviceInstanceId string) (*integration.MsxResponse, error)
	//Deprecated: User v4 Endpoint Instead
	GetDevice(deviceInstanceId string) (*integration.MsxResponse, error)
	//Deprecated: User v4 Endpoint Instead
	GetDevices(deviceInstanceId, subscriptionId, serialKey, tenantId *string, page, pageSize int) (*integration.MsxResponse, error)
	//Deprecated: User v4 Endpoint Instead
	CreateDevice(subscriptionId string, deviceInstanceId *string, deviceAttribute, deviceDefAttribute, status map[string]string) (*integration.MsxResponse, error)
	//Deprecated: User v4 Endpoint Instead
	UpdateDevice(deviceInstanceId string, deviceAttribute, deviceDefAttribute, status map[string]string) (*integration.MsxResponse, error)
	//Deprecated: User v4 Endpoint Instead
	DeleteDevice(deviceInstanceId string) (*integration.MsxResponse, error)

	CreateDeviceV4(deviceRequest DeviceCreateRequest) (*integration.MsxResponse, error)
	DeleteDeviceV4(deviceId string, force string) (*integration.MsxResponse, error)
	GetDevicesV4(requestQuery map[string][]string, page int, pageSize int) (*integration.MsxResponse, error)
	GetDeviceV4(deviceId string) (*integration.MsxResponse, error)
	UpdateDeviceV4(deviceRequest DeviceUpdateRequest, deviceId string) (*integration.MsxResponse, error)
	UpdateDeviceStatusV4(deviceStatus DeviceStatusUpdateRequest, deviceId string) (*integration.MsxResponse, error)

	ListDeviceTemplates(serviceType string, tenantId *types.UUID) (*integration.MsxResponse, error)
	GetDeviceTemplate(templateId types.UUID) (*integration.MsxResponse, error)
	AddDeviceTemplate(deviceTemplateCreateRequest DeviceTemplateCreateRequest) (*integration.MsxResponse, error)
	DeleteDeviceTemplate(templateId types.UUID) (*integration.MsxResponse, error)

	GetDeviceTemplateHistory(deviceInstanceId string) (*integration.MsxResponse, error)
	AttachDeviceTemplates(deviceId string, attachTemplateRequest AttachTemplateRequest) (*integration.MsxResponse, error)
	DetachDeviceTemplates(deviceId string) (*integration.MsxResponse, error)
	DetachDeviceTemplate(deviceId string, templateId types.UUID) (*integration.MsxResponse, error)
	UpdateTemplateAccess(templateId string, deviceTemplateAccess DeviceTemplateAccess) (*integration.MsxResponse, error)

	CreateDeviceActions(deviceActionList DeviceActionCreateRequests) (*integration.MsxResponse, error)
	UpdateDeviceActions(deviceActionList DeviceActionCreateRequests) (*integration.MsxResponse, error)

	GetAllControlPlanes(tenantId *string) (*integration.MsxResponse, error)
	CreateControlPlane(tenantId, name, url, resourceProvider, authenticationType string, tlsInsecure bool, attributes map[string]string) (*integration.MsxResponse, error)
	GetControlPlane(controlPlaneId string) (*integration.MsxResponse, error)
	UpdateControlPlane(controlPlaneId, tenantId, name, url, resourceProvider, authenticationType string, tlsInsecure bool, attributes map[string]string) (*integration.MsxResponse, error)
	DeleteControlPlane(controlPlaneId string) (*integration.MsxResponse, error)
	ConnectControlPlane(controlPlaneId string) (*integration.MsxResponse, error)
	ConnectUnmanagedControlPlane(username, password, url, resourceProvider string, tlsInsecure bool) (*integration.MsxResponse, error)

	CreateDeviceConnection(deviceConnection DeviceConnectionCreateRequest) (*integration.MsxResponse, *DeviceConnectionResponse, error)
	DeleteDeviceConnection(deviceConnectionId string) (*integration.MsxResponse, error)

	GetEntityShard(entityId string) (*integration.MsxResponse, error)
}

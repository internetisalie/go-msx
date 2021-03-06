// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

//go:generate mockery --inpackage --name=Api --structname=MockAwsResourceProvider

package awsresourceprovider

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

type Api interface {
	Connect(request AwsConnectRequest) (*integration.MsxResponse, error)
	//DEPRECATED use v2 instead
	GetRegions(controlPlaneId types.UUID) (*integration.MsxResponse, error)
	GetRegionsV2(controlPlaneId types.UUID, amiName *string) (*integration.MsxResponse, error) // optional amiName
	GetAvailabilityZones(controlPlaneId types.UUID, region string) (*integration.MsxResponse, error)
	GetResources(serviceConfigurationApplicationId types.UUID) (*integration.MsxResponse, error)
	GetVpnConnectionDetails(controlPlaneId types.UUID, vpnConnectionIds []string, region string) (*integration.MsxResponse, error)
	GetEc2InstanceStatus(controlPlaneId types.UUID, region string, instanceId string) (*integration.MsxResponse, error)
	GetTransitGatewayStatus(controlPlaneId types.UUID, region string, transitGatewayIds []string) (*integration.MsxResponse, error)
	GetTransitGatewayAttachmentStatus(controlPlaneId types.UUID, region string, transitGatewayAttachmentIds []string, resourceIds []string) (*integration.MsxResponse, error)
	GetTransitVPCStatus(controlPlaneId types.UUID, region string, transitVPCIds []string) (*integration.MsxResponse, error)
	GetStackOutputs(applicationId types.UUID) (*integration.MsxResponse, error)
	CheckStatus(applicationId types.UUID, request *CheckStatusRequest) (*integration.MsxResponse, error)
	GetInstanceType(controlPlaneId types.UUID, region string, availabilityZone string, instanceType string) (*integration.MsxResponse, error)
	GetAmiInformation(controlPlaneId types.UUID, amiName string, region string) (*integration.MsxResponse, error)
	GetRouteTableInformation(controlPlaneId types.UUID, region string, vpcId string) (*integration.MsxResponse, error)
	GetSecrets(controlPlaneId types.UUID, secretName string, region string) (*integration.MsxResponse, error)
}

// Ensure mock is up-to-date
var _ Api = new(MockAwsResourceProvider)

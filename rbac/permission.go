// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package rbac

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"github.com/pkg/errors"
)

const (
	PermissionIsApiAdmin                = "IS_API_ADMIN"
	PermissionAccessAllTenants          = "ACCESS_ALL_TENANTS"
	PermissionViewServices              = "VIEW_SERVICES"
	PermissionManageServices            = "MANAGE_SERVICES"
	PermissionViewContact               = "VIEW_CONTACT"
	PermissionManageContact             = "MANAGE_CONTACT"
	PermissionViewLocaleString          = "VIEW_LOCALE_STRING"
	PermissionManageLocaleString        = "MANAGE_LOCALE_STRING"
	PermissionViewIntegration           = "VIEW_INTEGRATION"
	PermissionManageIntegration         = "MANAGE_INTEGRATION"
	PermissionViewMaintenanceInfo       = "VIEW_MAINTENANCE_INFO"
	PermissionManageMaintenanceInfo     = "MANAGE_MAINTENANCE_INFO"
	PermissionViewMetadata              = "VIEW_METADATA"
	PermissionManageMetadata            = "MANAGE_METADATA"
	PermissionViewScheduledTask         = "VIEW_SCHEDULE_TASK"
	PermissionManageScheduledTask       = "MANAGE_SCHEDULE_TASK"
	PermissionViewTenantScheduledTask   = "VIEW_TENANT_SCHEDULE_TASK"
	PermissionManageTenantScheduledTask = "MANAGE_TENANT_SCHEDULE_TASK"
	PermissionViewBulkImport            = "VIEW_BULK_IMPORT"
	PermissionManageBulkImport          = "MANAGE_BULK_IMPORT"
	PermissionViewAlertService          = "VIEW_ALERT_SERVICE"
	PermissionManageAlertService        = "MANAGE_ALERT_SERVICE"
	PermissionViewThemes                = "VIEW_THEMES"
	PermissionManageThemes              = "MANAGE_THEMES"
	PermissionViewSlm                   = "VIEW_SLM"
	PermissionManageSlm                 = "MANAGE_SLM"
	PermissionViewNotification          = "VIEW_NOTIFICATION"
	PermissionManageNotification        = "MANAGE_NOTIFICATION"
	PermissionViewRegion                = "VIEW_REGION"
	PermissionManageRegion              = "MANAGE_REGION"
	PermissionViewVpn                   = "VIEW_VPN"
	PermissionManageVpn                 = "MANAGE_VPN"
	PermissionViewPnp                   = "VIEW_PNP"
	PermissionManagePnp                 = "MANAGE_PNP"
	PermissionViewPricePlan             = "VIEW_PRICE_PLAN"
	PermissionManagePricePlan           = "MANAGE_PRICE_PLAN"
	PermissionImportService             = "IMPORT_SERVICE"
	PermissionExportService             = "EXPORT_SERVICE"
	PermissionViewTermsConditions       = "VIEW_TERMS_CONDITIONS"
	PermissionManageTermsConditions     = "MANAGE_TERMS_CONDITIONS"
	PermissionViewIncident              = "VIEW_INCIDENT"
	PermissionManageIncident            = "MANAGE_INCIDENT"
	PermissionViewIncidentConfig        = "VIEW_INCIDENT_CONFIG"
	PermissionManageIncidentConfig      = "MANAGE_INCIDENT_CONFIG"
	PermissionViewStandardConfig        = "VIEW_STANDARD_CONFIG"
	PermissionManageStandardConfig      = "MANAGE_STANDARD_CONFIG"
)

var ErrUserDosNotHavePermission = errors.New("User does not have any of the required permissions")

func HasPermission(ctx context.Context, required []string) error {
	logger.WithContext(ctx).Debugf("Verifying permissions %q", required)

	userContextDetails, err := security.NewUserContextDetails(ctx)
	if err != nil {
		return err
	}

	permissions := userContextDetails.Permissions

	// bypass permission check if user has permission: IS_API_ADMIN
	for _, permission := range permissions {
		if PermissionIsApiAdmin == permission {
			return nil
		}
	}

	for _, p := range required {
		for _, permission := range permissions {
			if p == permission {
				return nil
			}
		}
	}

	return ErrUserDosNotHavePermission
}

func HasAccessAllTenants(ctx context.Context) (bool, error) {
	err := HasPermission(ctx, []string{PermissionAccessAllTenants})
	if err == ErrUserDosNotHavePermission {
		return false, nil
	}
	if err == nil {
		return true, nil
	}

	return false, err
}

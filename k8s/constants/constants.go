//go:build !azure
// +build !azure

/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 * SPDX-License-Identifier: MIT
 */

package constants

import _ "embed"

const (
	FullGroupName                = "symphony"
	AzureOperationIdKey          = "management.azure.com/operationId"
	AzureCorrelationIdKey        = "management.azure.com/correlationId"
	RunningAzureCorrelationIdKey = "management.azure.com/runningCorrelationId"
	AzureResourceIdKey           = "management.azure.com/resourceId"
	AzureSystemDataKey           = "management.azure.com/systemData"
	AzureTenantIdKey             = "management.azure.com/tenantId"
	AzureLocationKey             = "management.azure.com/location"
	AzureEdgeLocationKey         = "management.azure.com/customLocation"
	SummaryJobIdKey              = "SummaryJobIdKey"
	AzureCreatedByKey            = "createdBy"
	DefaultScope                 = "default"
	K8S                          = "symphony-k8s"
	OperationStartTimeKeyPostfix = FullGroupName + "/started-at"
	FinalizerPostfix             = FullGroupName + "/finalizer"
	ResourceSeperator            = "-v-"
	ReferenceSeparator           = ":"
	ActivityOperation_Write      = "Write"
	ActivityOperation_Read       = "Read"
	ActivityOperation_Delete     = "Delete"

	SolutionContainerOperationNamePrefix = "solutioncontainers.solution." + FullGroupName
	SolutionOperationNamePrefix          = "solutions.solution." + FullGroupName
	TargetOperationNamePrefix            = "targets.fabric." + FullGroupName
	InstanceOperationNamePrefix          = "instances.solution." + FullGroupName
	ActivationOperationNamePrefix        = "activations.workflow." + FullGroupName
	CatalogOperationNamePrefix           = "catalogs.federation." + FullGroupName
	CatalogEvalOperationNamePrefix       = "catalogevalexpression.federation." + FullGroupName
	CampaignOperationNamePrefix          = "campaigns.workflow." + FullGroupName
	CampaignContainerOperationNamePrefix = "campaigncontainers.workflow." + FullGroupName
	DiagnosticsOperationNamePrefix       = "diagnostics.monitor." + FullGroupName
)

// Environment variables keys
const (
	SymphonyAPIUrlEnvName = "SYMPHONY_API_URL"
	ConfigName            = "CONFIG_NAME"
	ApiCertEnvName        = "API_SERVING_CA"
	DeploymentFinalizer   = "DEPLOYMENT_FINALIZER"
)

// Eula Message
var (
	//go:embed eula.txt
	EulaMessage string
)

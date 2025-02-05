package google

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	googleoauth "golang.org/x/oauth2/google"
)

// Global MutexKV
var mutexKV = mutexkv.NewMutexKV()

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"credentials": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_CREDENTIALS",
					"GOOGLE_CLOUD_KEYFILE_JSON",
					"GCLOUD_KEYFILE_JSON",
				}, nil),
				ValidateFunc: validateCredentials,
			},

			"access_token": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_OAUTH_ACCESS_TOKEN",
				}, nil),
				ConflictsWith: []string{"credentials"},
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_PROJECT",
					"GOOGLE_CLOUD_PROJECT",
					"GCLOUD_PROJECT",
					"CLOUDSDK_CORE_PROJECT",
				}, nil),
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_REGION",
					"GCLOUD_REGION",
					"CLOUDSDK_COMPUTE_REGION",
				}, nil),
			},

			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_ZONE",
					"GCLOUD_ZONE",
					"CLOUDSDK_COMPUTE_ZONE",
				}, nil),
			},

			"scopes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"batching": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"send_after": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "10s",
							ValidateFunc: validateNonNegativeDuration(),
						},
						"enable_batching": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},

			"user_project_override": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			// Generated Products
			"access_context_manager_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_ACCESS_CONTEXT_MANAGER_CUSTOM_ENDPOINT",
				}, AccessContextManagerDefaultBasePath),
			},
			"app_engine_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_APP_ENGINE_CUSTOM_ENDPOINT",
				}, AppEngineDefaultBasePath),
			},
			"big_query_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_BIG_QUERY_CUSTOM_ENDPOINT",
				}, BigQueryDefaultBasePath),
			},
			"bigquery_data_transfer_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_BIGQUERY_DATA_TRANSFER_CUSTOM_ENDPOINT",
				}, BigqueryDataTransferDefaultBasePath),
			},
			"bigtable_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_BIGTABLE_CUSTOM_ENDPOINT",
				}, BigtableDefaultBasePath),
			},
			"binary_authorization_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_BINARY_AUTHORIZATION_CUSTOM_ENDPOINT",
				}, BinaryAuthorizationDefaultBasePath),
			},
			"cloud_build_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_CLOUD_BUILD_CUSTOM_ENDPOINT",
				}, CloudBuildDefaultBasePath),
			},
			"cloud_scheduler_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_CLOUD_SCHEDULER_CUSTOM_ENDPOINT",
				}, CloudSchedulerDefaultBasePath),
			},
			"compute_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_COMPUTE_CUSTOM_ENDPOINT",
				}, ComputeDefaultBasePath),
			},
			"dns_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_DNS_CUSTOM_ENDPOINT",
				}, DnsDefaultBasePath),
			},
			"filestore_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_FILESTORE_CUSTOM_ENDPOINT",
				}, FilestoreDefaultBasePath),
			},
			"firestore_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_FIRESTORE_CUSTOM_ENDPOINT",
				}, FirestoreDefaultBasePath),
			},
			"kms_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_KMS_CUSTOM_ENDPOINT",
				}, KmsDefaultBasePath),
			},
			"logging_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_LOGGING_CUSTOM_ENDPOINT",
				}, LoggingDefaultBasePath),
			},
			"ml_engine_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_ML_ENGINE_CUSTOM_ENDPOINT",
				}, MLEngineDefaultBasePath),
			},
			"monitoring_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_MONITORING_CUSTOM_ENDPOINT",
				}, MonitoringDefaultBasePath),
			},
			"pubsub_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_PUBSUB_CUSTOM_ENDPOINT",
				}, PubsubDefaultBasePath),
			},
			"redis_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_REDIS_CUSTOM_ENDPOINT",
				}, RedisDefaultBasePath),
			},
			"resource_manager_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_RESOURCE_MANAGER_CUSTOM_ENDPOINT",
				}, ResourceManagerDefaultBasePath),
			},
			"security_center_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_SECURITY_CENTER_CUSTOM_ENDPOINT",
				}, SecurityCenterDefaultBasePath),
			},
			"source_repo_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_SOURCE_REPO_CUSTOM_ENDPOINT",
				}, SourceRepoDefaultBasePath),
			},
			"spanner_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_SPANNER_CUSTOM_ENDPOINT",
				}, SpannerDefaultBasePath),
			},
			"sql_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_SQL_CUSTOM_ENDPOINT",
				}, SqlDefaultBasePath),
			},
			"storage_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_STORAGE_CUSTOM_ENDPOINT",
				}, StorageDefaultBasePath),
			},
			"tpu_custom_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateCustomEndpoint,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOOGLE_TPU_CUSTOM_ENDPOINT",
				}, TpuDefaultBasePath),
			},

			// Handwritten Products / Versioned / Atypical Entries
			CloudBillingCustomEndpointEntryKey:           CloudBillingCustomEndpointEntry,
			ComposerCustomEndpointEntryKey:               ComposerCustomEndpointEntry,
			ComputeBetaCustomEndpointEntryKey:            ComputeBetaCustomEndpointEntry,
			ContainerCustomEndpointEntryKey:              ContainerCustomEndpointEntry,
			ContainerBetaCustomEndpointEntryKey:          ContainerBetaCustomEndpointEntry,
			DataprocCustomEndpointEntryKey:               DataprocCustomEndpointEntry,
			DataprocBetaCustomEndpointEntryKey:           DataprocBetaCustomEndpointEntry,
			DataflowCustomEndpointEntryKey:               DataflowCustomEndpointEntry,
			DnsBetaCustomEndpointEntryKey:                DnsBetaCustomEndpointEntry,
			IamCredentialsCustomEndpointEntryKey:         IamCredentialsCustomEndpointEntry,
			ResourceManagerV2Beta1CustomEndpointEntryKey: ResourceManagerV2Beta1CustomEndpointEntry,
			RuntimeconfigCustomEndpointEntryKey:          RuntimeconfigCustomEndpointEntry,
			IAMCustomEndpointEntryKey:                    IAMCustomEndpointEntry,
			ServiceManagementCustomEndpointEntryKey:      ServiceManagementCustomEndpointEntry,
			ServiceNetworkingCustomEndpointEntryKey:      ServiceNetworkingCustomEndpointEntry,
			ServiceUsageCustomEndpointEntryKey:           ServiceUsageCustomEndpointEntry,
			CloudFunctionsCustomEndpointEntryKey:         CloudFunctionsCustomEndpointEntry,
			CloudIoTCustomEndpointEntryKey:               CloudIoTCustomEndpointEntry,
			StorageTransferCustomEndpointEntryKey:        StorageTransferCustomEndpointEntry,
			BigtableAdminCustomEndpointEntryKey:          BigtableAdminCustomEndpointEntry,
		},

		DataSourcesMap: map[string]*schema.Resource{
			"google_active_folder":                            dataSourceGoogleActiveFolder(),
			"google_billing_account":                          dataSourceGoogleBillingAccount(),
			"google_dns_managed_zone":                         dataSourceDnsManagedZone(),
			"google_client_config":                            dataSourceGoogleClientConfig(),
			"google_client_openid_userinfo":                   dataSourceGoogleClientOpenIDUserinfo(),
			"google_cloudfunctions_function":                  dataSourceGoogleCloudFunctionsFunction(),
			"google_composer_image_versions":                  dataSourceGoogleComposerImageVersions(),
			"google_compute_address":                          dataSourceGoogleComputeAddress(),
			"google_compute_backend_service":                  dataSourceGoogleComputeBackendService(),
			"google_compute_default_service_account":          dataSourceGoogleComputeDefaultServiceAccount(),
			"google_compute_forwarding_rule":                  dataSourceGoogleComputeForwardingRule(),
			"google_compute_image":                            dataSourceGoogleComputeImage(),
			"google_compute_instance":                         dataSourceGoogleComputeInstance(),
			"google_compute_global_address":                   dataSourceGoogleComputeGlobalAddress(),
			"google_compute_instance_group":                   dataSourceGoogleComputeInstanceGroup(),
			"google_compute_lb_ip_ranges":                     dataSourceGoogleComputeLbIpRanges(),
			"google_compute_network":                          dataSourceGoogleComputeNetwork(),
			"google_compute_network_endpoint_group":           dataSourceGoogleComputeNetworkEndpointGroup(),
			"google_compute_node_types":                       dataSourceGoogleComputeNodeTypes(),
			"google_compute_regions":                          dataSourceGoogleComputeRegions(),
			"google_compute_region_instance_group":            dataSourceGoogleComputeRegionInstanceGroup(),
			"google_compute_subnetwork":                       dataSourceGoogleComputeSubnetwork(),
			"google_compute_zones":                            dataSourceGoogleComputeZones(),
			"google_compute_vpn_gateway":                      dataSourceGoogleComputeVpnGateway(),
			"google_compute_ssl_policy":                       dataSourceGoogleComputeSslPolicy(),
			"google_compute_ssl_certificate":                  dataSourceGoogleComputeSslCertificate(),
			"google_container_cluster":                        dataSourceGoogleContainerCluster(),
			"google_container_engine_versions":                dataSourceGoogleContainerEngineVersions(),
			"google_container_registry_repository":            dataSourceGoogleContainerRepo(),
			"google_container_registry_image":                 dataSourceGoogleContainerImage(),
			"google_iam_policy":                               dataSourceGoogleIamPolicy(),
			"google_iam_role":                                 dataSourceGoogleIamRole(),
			"google_kms_secret":                               dataSourceGoogleKmsSecret(),
			"google_kms_secret_ciphertext":                    dataSourceGoogleKmsSecretCiphertext(),
			"google_kms_key_ring":                             dataSourceGoogleKmsKeyRing(),
			"google_kms_crypto_key":                           dataSourceGoogleKmsCryptoKey(),
			"google_kms_crypto_key_version":                   dataSourceGoogleKmsCryptoKeyVersion(),
			"google_folder":                                   dataSourceGoogleFolder(),
			"google_folder_organization_policy":               dataSourceGoogleFolderOrganizationPolicy(),
			"google_netblock_ip_ranges":                       dataSourceGoogleNetblockIpRanges(),
			"google_organization":                             dataSourceGoogleOrganization(),
			"google_project":                                  dataSourceGoogleProject(),
			"google_projects":                                 dataSourceGoogleProjects(),
			"google_project_organization_policy":              dataSourceGoogleProjectOrganizationPolicy(),
			"google_project_services":                         dataSourceGoogleProjectServices(),
			"google_service_account":                          dataSourceGoogleServiceAccount(),
			"google_service_account_access_token":             dataSourceGoogleServiceAccountAccessToken(),
			"google_service_account_key":                      dataSourceGoogleServiceAccountKey(),
			"google_storage_bucket_object":                    dataSourceGoogleStorageBucketObject(),
			"google_storage_object_signed_url":                dataSourceGoogleSignedUrl(),
			"google_storage_project_service_account":          dataSourceGoogleStorageProjectServiceAccount(),
			"google_storage_transfer_project_service_account": dataSourceGoogleStorageTransferProjectServiceAccount(),
			"google_tpu_tensorflow_versions":                  dataSourceTpuTensorflowVersions(),
		},

		ResourcesMap: ResourceMap(),

		ConfigureFunc: providerConfigure,
	}
}

// Generated resources: 72
// Generated IAM resources: 6
// Total generated resources: 78
func ResourceMap() map[string]*schema.Resource {
	resourceMap, _ := ResourceMapWithErrors()
	return resourceMap
}

func ResourceMapWithErrors() (map[string]*schema.Resource, error) {
	return mergeResourceMaps(
		map[string]*schema.Resource{
			"google_access_context_manager_access_policy":     resourceAccessContextManagerAccessPolicy(),
			"google_access_context_manager_access_level":      resourceAccessContextManagerAccessLevel(),
			"google_access_context_manager_service_perimeter": resourceAccessContextManagerServicePerimeter(),
			"google_app_engine_firewall_rule":                 resourceAppEngineFirewallRule(),
			"google_bigquery_dataset":                         resourceBigQueryDataset(),
			"google_bigquery_data_transfer_config":            resourceBigqueryDataTransferConfig(),
			"google_bigtable_app_profile":                     resourceBigtableAppProfile(),
			"google_binary_authorization_attestor":            resourceBinaryAuthorizationAttestor(),
			"google_binary_authorization_policy":              resourceBinaryAuthorizationPolicy(),
			"google_cloudbuild_trigger":                       resourceCloudBuildTrigger(),
			"google_cloud_scheduler_job":                      resourceCloudSchedulerJob(),
			"google_compute_address":                          resourceComputeAddress(),
			"google_compute_autoscaler":                       resourceComputeAutoscaler(),
			"google_compute_backend_bucket":                   resourceComputeBackendBucket(),
			"google_compute_backend_bucket_signed_url_key":    resourceComputeBackendBucketSignedUrlKey(),
			"google_compute_backend_service":                  resourceComputeBackendService(),
			"google_compute_region_backend_service":           resourceComputeRegionBackendService(),
			"google_compute_backend_service_signed_url_key":   resourceComputeBackendServiceSignedUrlKey(),
			"google_compute_disk":                             resourceComputeDisk(),
			"google_compute_firewall":                         resourceComputeFirewall(),
			"google_compute_forwarding_rule":                  resourceComputeForwardingRule(),
			"google_compute_global_address":                   resourceComputeGlobalAddress(),
			"google_compute_global_forwarding_rule":           resourceComputeGlobalForwardingRule(),
			"google_compute_http_health_check":                resourceComputeHttpHealthCheck(),
			"google_compute_https_health_check":               resourceComputeHttpsHealthCheck(),
			"google_compute_health_check":                     resourceComputeHealthCheck(),
			"google_compute_image":                            resourceComputeImage(),
			"google_compute_interconnect_attachment":          resourceComputeInterconnectAttachment(),
			"google_compute_network":                          resourceComputeNetwork(),
			"google_compute_network_endpoint":                 resourceComputeNetworkEndpoint(),
			"google_compute_network_endpoint_group":           resourceComputeNetworkEndpointGroup(),
			"google_compute_node_group":                       resourceComputeNodeGroup(),
			"google_compute_node_template":                    resourceComputeNodeTemplate(),
			"google_compute_region_autoscaler":                resourceComputeRegionAutoscaler(),
			"google_compute_region_disk":                      resourceComputeRegionDisk(),
			"google_compute_route":                            resourceComputeRoute(),
			"google_compute_router":                           resourceComputeRouter(),
			"google_compute_snapshot":                         resourceComputeSnapshot(),
			"google_compute_ssl_certificate":                  resourceComputeSslCertificate(),
			"google_compute_ssl_policy":                       resourceComputeSslPolicy(),
			"google_compute_subnetwork":                       resourceComputeSubnetwork(),
			"google_compute_target_http_proxy":                resourceComputeTargetHttpProxy(),
			"google_compute_target_https_proxy":               resourceComputeTargetHttpsProxy(),
			"google_compute_target_instance":                  resourceComputeTargetInstance(),
			"google_compute_target_ssl_proxy":                 resourceComputeTargetSslProxy(),
			"google_compute_target_tcp_proxy":                 resourceComputeTargetTcpProxy(),
			"google_compute_vpn_gateway":                      resourceComputeVpnGateway(),
			"google_compute_url_map":                          resourceComputeUrlMap(),
			"google_compute_vpn_tunnel":                       resourceComputeVpnTunnel(),
			"google_dns_managed_zone":                         resourceDnsManagedZone(),
			"google_filestore_instance":                       resourceFilestoreInstance(),
			"google_firestore_index":                          resourceFirestoreIndex(),
			"google_kms_key_ring":                             resourceKmsKeyRing(),
			"google_kms_crypto_key":                           resourceKmsCryptoKey(),
			"google_logging_metric":                           resourceLoggingMetric(),
			"google_ml_engine_model":                          resourceMLEngineModel(),
			"google_monitoring_alert_policy":                  resourceMonitoringAlertPolicy(),
			"google_monitoring_group":                         resourceMonitoringGroup(),
			"google_monitoring_notification_channel":          resourceMonitoringNotificationChannel(),
			"google_monitoring_uptime_check_config":           resourceMonitoringUptimeCheckConfig(),
			"google_pubsub_topic":                             resourcePubsubTopic(),
			"google_pubsub_topic_iam_binding":                 ResourceIamBinding(PubsubTopicIamSchema, PubsubTopicIamUpdaterProducer, PubsubTopicIdParseFunc),
			"google_pubsub_topic_iam_member":                  ResourceIamMember(PubsubTopicIamSchema, PubsubTopicIamUpdaterProducer, PubsubTopicIdParseFunc),
			"google_pubsub_topic_iam_policy":                  ResourceIamPolicy(PubsubTopicIamSchema, PubsubTopicIamUpdaterProducer, PubsubTopicIdParseFunc),
			"google_pubsub_subscription":                      resourcePubsubSubscription(),
			"google_redis_instance":                           resourceRedisInstance(),
			"google_resource_manager_lien":                    resourceResourceManagerLien(),
			"google_scc_source":                               resourceSecurityCenterSource(),
			"google_sourcerepo_repository":                    resourceSourceRepoRepository(),
			"google_sourcerepo_repository_iam_binding":        ResourceIamBinding(SourceRepoRepositoryIamSchema, SourceRepoRepositoryIamUpdaterProducer, SourceRepoRepositoryIdParseFunc),
			"google_sourcerepo_repository_iam_member":         ResourceIamMember(SourceRepoRepositoryIamSchema, SourceRepoRepositoryIamUpdaterProducer, SourceRepoRepositoryIdParseFunc),
			"google_sourcerepo_repository_iam_policy":         ResourceIamPolicy(SourceRepoRepositoryIamSchema, SourceRepoRepositoryIamUpdaterProducer, SourceRepoRepositoryIdParseFunc),
			"google_spanner_instance":                         resourceSpannerInstance(),
			"google_spanner_database":                         resourceSpannerDatabase(),
			"google_sql_database":                             resourceSqlDatabase(),
			"google_storage_object_access_control":            resourceStorageObjectAccessControl(),
			"google_storage_default_object_access_control":    resourceStorageDefaultObjectAccessControl(),
			"google_tpu_node":                                 resourceTpuNode(),
		},
		map[string]*schema.Resource{
			"google_app_engine_application":                resourceAppEngineApplication(),
			"google_bigquery_table":                        resourceBigQueryTable(),
			"google_bigtable_instance":                     resourceBigtableInstance(),
			"google_bigtable_instance_iam_binding":         ResourceIamBinding(IamBigtableInstanceSchema, NewBigtableInstanceUpdater, BigtableInstanceIdParseFunc),
			"google_bigtable_instance_iam_member":          ResourceIamMember(IamBigtableInstanceSchema, NewBigtableInstanceUpdater, BigtableInstanceIdParseFunc),
			"google_bigtable_instance_iam_policy":          ResourceIamPolicy(IamBigtableInstanceSchema, NewBigtableInstanceUpdater, BigtableInstanceIdParseFunc),
			"google_bigtable_table":                        resourceBigtableTable(),
			"google_billing_account_iam_binding":           ResourceIamBinding(IamBillingAccountSchema, NewBillingAccountIamUpdater, BillingAccountIdParseFunc),
			"google_billing_account_iam_member":            ResourceIamMember(IamBillingAccountSchema, NewBillingAccountIamUpdater, BillingAccountIdParseFunc),
			"google_billing_account_iam_policy":            ResourceIamPolicy(IamBillingAccountSchema, NewBillingAccountIamUpdater, BillingAccountIdParseFunc),
			"google_cloudfunctions_function":               resourceCloudFunctionsFunction(),
			"google_cloudiot_registry":                     resourceCloudIoTRegistry(),
			"google_composer_environment":                  resourceComposerEnvironment(),
			"google_compute_attached_disk":                 resourceComputeAttachedDisk(),
			"google_compute_instance":                      resourceComputeInstance(),
			"google_compute_instance_from_template":        resourceComputeInstanceFromTemplate(),
			"google_compute_instance_group":                resourceComputeInstanceGroup(),
			"google_compute_instance_group_manager":        resourceComputeInstanceGroupManager(),
			"google_compute_instance_iam_binding":          ResourceIamBinding(IamComputeInstanceSchema, NewComputeInstanceIamUpdater, ComputeInstanceIdParseFunc),
			"google_compute_instance_iam_member":           ResourceIamMember(IamComputeInstanceSchema, NewComputeInstanceIamUpdater, ComputeInstanceIdParseFunc),
			"google_compute_instance_iam_policy":           ResourceIamPolicy(IamComputeInstanceSchema, NewComputeInstanceIamUpdater, ComputeInstanceIdParseFunc),
			"google_compute_instance_template":             resourceComputeInstanceTemplate(),
			"google_compute_network_peering":               resourceComputeNetworkPeering(),
			"google_compute_project_default_network_tier":  resourceComputeProjectDefaultNetworkTier(),
			"google_compute_project_metadata":              resourceComputeProjectMetadata(),
			"google_compute_project_metadata_item":         resourceComputeProjectMetadataItem(),
			"google_compute_region_instance_group_manager": resourceComputeRegionInstanceGroupManager(),
			"google_compute_router_interface":              resourceComputeRouterInterface(),
			"google_compute_router_nat":                    resourceComputeRouterNat(),
			"google_compute_router_peer":                   resourceComputeRouterPeer(),
			"google_compute_security_policy":               resourceComputeSecurityPolicy(),
			"google_compute_shared_vpc_host_project":       resourceComputeSharedVpcHostProject(),
			"google_compute_shared_vpc_service_project":    resourceComputeSharedVpcServiceProject(),
			"google_compute_subnetwork_iam_binding":        ResourceIamBinding(IamComputeSubnetworkSchema, NewComputeSubnetworkIamUpdater, ComputeSubnetworkIdParseFunc),
			"google_compute_subnetwork_iam_member":         ResourceIamMember(IamComputeSubnetworkSchema, NewComputeSubnetworkIamUpdater, ComputeSubnetworkIdParseFunc),
			"google_compute_subnetwork_iam_policy":         ResourceIamPolicy(IamComputeSubnetworkSchema, NewComputeSubnetworkIamUpdater, ComputeSubnetworkIdParseFunc),
			"google_compute_target_pool":                   resourceComputeTargetPool(),
			"google_container_cluster":                     resourceContainerCluster(),
			"google_container_node_pool":                   resourceContainerNodePool(),
			"google_dataflow_job":                          resourceDataflowJob(),
			"google_dataproc_cluster":                      resourceDataprocCluster(),
			"google_dataproc_cluster_iam_binding":          ResourceIamBinding(IamDataprocClusterSchema, NewDataprocClusterUpdater, DataprocClusterIdParseFunc),
			"google_dataproc_cluster_iam_member":           ResourceIamMember(IamDataprocClusterSchema, NewDataprocClusterUpdater, DataprocClusterIdParseFunc),
			"google_dataproc_cluster_iam_policy":           ResourceIamPolicy(IamDataprocClusterSchema, NewDataprocClusterUpdater, DataprocClusterIdParseFunc),
			"google_dataproc_job":                          resourceDataprocJob(),
			"google_dataproc_job_iam_binding":              ResourceIamBinding(IamDataprocJobSchema, NewDataprocJobUpdater, DataprocJobIdParseFunc),
			"google_dataproc_job_iam_member":               ResourceIamMember(IamDataprocJobSchema, NewDataprocJobUpdater, DataprocJobIdParseFunc),
			"google_dataproc_job_iam_policy":               ResourceIamPolicy(IamDataprocJobSchema, NewDataprocJobUpdater, DataprocJobIdParseFunc),
			"google_dns_record_set":                        resourceDnsRecordSet(),
			"google_endpoints_service":                     resourceEndpointsService(),
			"google_folder":                                resourceGoogleFolder(),
			"google_folder_iam_binding":                    ResourceIamBinding(IamFolderSchema, NewFolderIamUpdater, FolderIdParseFunc),
			"google_folder_iam_member":                     ResourceIamMember(IamFolderSchema, NewFolderIamUpdater, FolderIdParseFunc),
			"google_folder_iam_policy":                     ResourceIamPolicy(IamFolderSchema, NewFolderIamUpdater, FolderIdParseFunc),
			"google_folder_organization_policy":            resourceGoogleFolderOrganizationPolicy(),
			"google_logging_billing_account_sink":          resourceLoggingBillingAccountSink(),
			"google_logging_billing_account_exclusion":     ResourceLoggingExclusion(BillingAccountLoggingExclusionSchema, NewBillingAccountLoggingExclusionUpdater, billingAccountLoggingExclusionIdParseFunc),
			"google_logging_organization_sink":             resourceLoggingOrganizationSink(),
			"google_logging_organization_exclusion":        ResourceLoggingExclusion(OrganizationLoggingExclusionSchema, NewOrganizationLoggingExclusionUpdater, organizationLoggingExclusionIdParseFunc),
			"google_logging_folder_sink":                   resourceLoggingFolderSink(),
			"google_logging_folder_exclusion":              ResourceLoggingExclusion(FolderLoggingExclusionSchema, NewFolderLoggingExclusionUpdater, folderLoggingExclusionIdParseFunc),
			"google_logging_project_sink":                  resourceLoggingProjectSink(),
			"google_logging_project_exclusion":             ResourceLoggingExclusion(ProjectLoggingExclusionSchema, NewProjectLoggingExclusionUpdater, projectLoggingExclusionIdParseFunc),
			"google_kms_key_ring_iam_binding":              ResourceIamBinding(IamKmsKeyRingSchema, NewKmsKeyRingIamUpdater, KeyRingIdParseFunc),
			"google_kms_key_ring_iam_member":               ResourceIamMember(IamKmsKeyRingSchema, NewKmsKeyRingIamUpdater, KeyRingIdParseFunc),
			"google_kms_key_ring_iam_policy":               ResourceIamPolicy(IamKmsKeyRingSchema, NewKmsKeyRingIamUpdater, KeyRingIdParseFunc),
			"google_kms_crypto_key_iam_binding":            ResourceIamBinding(IamKmsCryptoKeySchema, NewKmsCryptoKeyIamUpdater, CryptoIdParseFunc),
			"google_kms_crypto_key_iam_member":             ResourceIamMember(IamKmsCryptoKeySchema, NewKmsCryptoKeyIamUpdater, CryptoIdParseFunc),
			"google_service_networking_connection":         resourceServiceNetworkingConnection(),
			"google_spanner_instance_iam_binding":          ResourceIamBinding(IamSpannerInstanceSchema, NewSpannerInstanceIamUpdater, SpannerInstanceIdParseFunc),
			"google_spanner_instance_iam_member":           ResourceIamMember(IamSpannerInstanceSchema, NewSpannerInstanceIamUpdater, SpannerInstanceIdParseFunc),
			"google_spanner_instance_iam_policy":           ResourceIamPolicy(IamSpannerInstanceSchema, NewSpannerInstanceIamUpdater, SpannerInstanceIdParseFunc),
			"google_spanner_database_iam_binding":          ResourceIamBinding(IamSpannerDatabaseSchema, NewSpannerDatabaseIamUpdater, SpannerDatabaseIdParseFunc),
			"google_spanner_database_iam_member":           ResourceIamMember(IamSpannerDatabaseSchema, NewSpannerDatabaseIamUpdater, SpannerDatabaseIdParseFunc),
			"google_spanner_database_iam_policy":           ResourceIamPolicy(IamSpannerDatabaseSchema, NewSpannerDatabaseIamUpdater, SpannerDatabaseIdParseFunc),
			"google_sql_database_instance":                 resourceSqlDatabaseInstance(),
			"google_sql_ssl_cert":                          resourceSqlSslCert(),
			"google_sql_user":                              resourceSqlUser(),
			"google_organization_iam_binding":              ResourceIamBinding(IamOrganizationSchema, NewOrganizationIamUpdater, OrgIdParseFunc),
			"google_organization_iam_custom_role":          resourceGoogleOrganizationIamCustomRole(),
			"google_organization_iam_member":               ResourceIamMember(IamOrganizationSchema, NewOrganizationIamUpdater, OrgIdParseFunc),
			"google_organization_iam_policy":               ResourceIamPolicy(IamOrganizationSchema, NewOrganizationIamUpdater, OrgIdParseFunc),
			"google_organization_policy":                   resourceGoogleOrganizationPolicy(),
			"google_project":                               resourceGoogleProject(),
			"google_project_iam_policy":                    resourceGoogleProjectIamPolicy(),
			"google_project_iam_binding":                   ResourceIamBindingWithBatching(IamProjectSchema, NewProjectIamUpdater, ProjectIdParseFunc, IamBatchingEnabled),
			"google_project_iam_member":                    ResourceIamMemberWithBatching(IamProjectSchema, NewProjectIamUpdater, ProjectIdParseFunc, IamBatchingEnabled),
			"google_project_iam_audit_config":              ResourceIamAuditConfigWithBatching(IamProjectSchema, NewProjectIamUpdater, ProjectIdParseFunc, IamBatchingEnabled),
			"google_project_service":                       resourceGoogleProjectService(),
			"google_project_iam_custom_role":               resourceGoogleProjectIamCustomRole(),
			"google_project_organization_policy":           resourceGoogleProjectOrganizationPolicy(),
			"google_project_usage_export_bucket":           resourceProjectUsageBucket(),
			"google_project_services":                      resourceGoogleProjectServices(),
			"google_pubsub_subscription_iam_binding":       ResourceIamBinding(IamPubsubSubscriptionSchema, NewPubsubSubscriptionIamUpdater, PubsubSubscriptionIdParseFunc),
			"google_pubsub_subscription_iam_member":        ResourceIamMember(IamPubsubSubscriptionSchema, NewPubsubSubscriptionIamUpdater, PubsubSubscriptionIdParseFunc),
			"google_pubsub_subscription_iam_policy":        ResourceIamPolicy(IamPubsubSubscriptionSchema, NewPubsubSubscriptionIamUpdater, PubsubSubscriptionIdParseFunc),
			"google_runtimeconfig_config":                  resourceRuntimeconfigConfig(),
			"google_runtimeconfig_variable":                resourceRuntimeconfigVariable(),
			"google_service_account":                       resourceGoogleServiceAccount(),
			"google_service_account_iam_binding":           ResourceIamBinding(IamServiceAccountSchema, NewServiceAccountIamUpdater, ServiceAccountIdParseFunc),
			"google_service_account_iam_member":            ResourceIamMember(IamServiceAccountSchema, NewServiceAccountIamUpdater, ServiceAccountIdParseFunc),
			"google_service_account_iam_policy":            ResourceIamPolicy(IamServiceAccountSchema, NewServiceAccountIamUpdater, ServiceAccountIdParseFunc),
			"google_service_account_key":                   resourceGoogleServiceAccountKey(),
			"google_storage_bucket":                        resourceStorageBucket(),
			"google_storage_bucket_acl":                    resourceStorageBucketAcl(),
			// Legacy roles such as roles/storage.legacyBucketReader are automatically added
			// when creating a bucket. For this reason, it is better not to add the authoritative
			// google_storage_bucket_iam_policy resource.
			"google_storage_bucket_iam_binding": ResourceIamBinding(IamStorageBucketSchema, NewStorageBucketIamUpdater, StorageBucketIdParseFunc),
			"google_storage_bucket_iam_member":  ResourceIamMember(IamStorageBucketSchema, NewStorageBucketIamUpdater, StorageBucketIdParseFunc),
			"google_storage_bucket_iam_policy":  ResourceIamPolicy(IamStorageBucketSchema, NewStorageBucketIamUpdater, StorageBucketIdParseFunc),
			"google_storage_bucket_object":      resourceStorageBucketObject(),
			"google_storage_object_acl":         resourceStorageObjectAcl(),
			"google_storage_default_object_acl": resourceStorageDefaultObjectAcl(),
			"google_storage_notification":       resourceStorageNotification(),
			"google_storage_transfer_job":       resourceStorageTransferJob(),
		},
	)
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Project:             d.Get("project").(string),
		Region:              d.Get("region").(string),
		Zone:                d.Get("zone").(string),
		UserProjectOverride: d.Get("user_project_override").(bool),
	}

	// Add credential source
	if v, ok := d.GetOk("access_token"); ok {
		config.AccessToken = v.(string)
	} else if v, ok := d.GetOk("credentials"); ok {
		config.Credentials = v.(string)
	}

	scopes := d.Get("scopes").([]interface{})
	if len(scopes) > 0 {
		config.Scopes = make([]string, len(scopes), len(scopes))
	}
	for i, scope := range scopes {
		config.Scopes[i] = scope.(string)
	}

	batchCfg, err := expandProviderBatchingConfig(d.Get("batching"))
	if err != nil {
		return nil, err
	}
	config.BatchingConfig = batchCfg

	// Generated products
	config.AccessContextManagerBasePath = d.Get("access_context_manager_custom_endpoint").(string)
	config.AppEngineBasePath = d.Get("app_engine_custom_endpoint").(string)
	config.BigQueryBasePath = d.Get("big_query_custom_endpoint").(string)
	config.BigqueryDataTransferBasePath = d.Get("bigquery_data_transfer_custom_endpoint").(string)
	config.BigtableBasePath = d.Get("bigtable_custom_endpoint").(string)
	config.BinaryAuthorizationBasePath = d.Get("binary_authorization_custom_endpoint").(string)
	config.CloudBuildBasePath = d.Get("cloud_build_custom_endpoint").(string)
	config.CloudSchedulerBasePath = d.Get("cloud_scheduler_custom_endpoint").(string)
	config.ComputeBasePath = d.Get("compute_custom_endpoint").(string)
	config.DnsBasePath = d.Get("dns_custom_endpoint").(string)
	config.FilestoreBasePath = d.Get("filestore_custom_endpoint").(string)
	config.FirestoreBasePath = d.Get("firestore_custom_endpoint").(string)
	config.KmsBasePath = d.Get("kms_custom_endpoint").(string)
	config.LoggingBasePath = d.Get("logging_custom_endpoint").(string)
	config.MLEngineBasePath = d.Get("ml_engine_custom_endpoint").(string)
	config.MonitoringBasePath = d.Get("monitoring_custom_endpoint").(string)
	config.PubsubBasePath = d.Get("pubsub_custom_endpoint").(string)
	config.RedisBasePath = d.Get("redis_custom_endpoint").(string)
	config.ResourceManagerBasePath = d.Get("resource_manager_custom_endpoint").(string)
	config.SecurityCenterBasePath = d.Get("security_center_custom_endpoint").(string)
	config.SourceRepoBasePath = d.Get("source_repo_custom_endpoint").(string)
	config.SpannerBasePath = d.Get("spanner_custom_endpoint").(string)
	config.SqlBasePath = d.Get("sql_custom_endpoint").(string)
	config.StorageBasePath = d.Get("storage_custom_endpoint").(string)
	config.TpuBasePath = d.Get("tpu_custom_endpoint").(string)

	// Handwritten Products / Versioned / Atypical Entries

	config.CloudBillingBasePath = d.Get(CloudBillingCustomEndpointEntryKey).(string)
	config.ComposerBasePath = d.Get(ComposerCustomEndpointEntryKey).(string)
	config.ComputeBetaBasePath = d.Get(ComputeBetaCustomEndpointEntryKey).(string)
	config.ContainerBasePath = d.Get(ContainerCustomEndpointEntryKey).(string)
	config.ContainerBetaBasePath = d.Get(ContainerBetaCustomEndpointEntryKey).(string)
	config.DataprocBasePath = d.Get(DataprocCustomEndpointEntryKey).(string)
	config.DataprocBetaBasePath = d.Get(DataprocBetaCustomEndpointEntryKey).(string)
	config.DataflowBasePath = d.Get(DataflowCustomEndpointEntryKey).(string)
	config.DnsBetaBasePath = d.Get(DnsBetaCustomEndpointEntryKey).(string)
	config.IamCredentialsBasePath = d.Get(IamCredentialsCustomEndpointEntryKey).(string)
	config.ResourceManagerV2Beta1BasePath = d.Get(ResourceManagerV2Beta1CustomEndpointEntryKey).(string)
	config.RuntimeconfigBasePath = d.Get(RuntimeconfigCustomEndpointEntryKey).(string)
	config.IAMBasePath = d.Get(IAMCustomEndpointEntryKey).(string)
	config.ServiceManagementBasePath = d.Get(ServiceManagementCustomEndpointEntryKey).(string)
	config.ServiceNetworkingBasePath = d.Get(ServiceNetworkingCustomEndpointEntryKey).(string)
	config.ServiceUsageBasePath = d.Get(ServiceUsageCustomEndpointEntryKey).(string)
	config.CloudFunctionsBasePath = d.Get(CloudFunctionsCustomEndpointEntryKey).(string)
	config.CloudIoTBasePath = d.Get(CloudIoTCustomEndpointEntryKey).(string)
	config.StorageTransferBasePath = d.Get(StorageTransferCustomEndpointEntryKey).(string)
	config.BigtableAdminBasePath = d.Get(BigtableAdminCustomEndpointEntryKey).(string)

	if err := config.LoadAndValidate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func validateCredentials(v interface{}, k string) (warnings []string, errors []error) {
	if v == nil || v.(string) == "" {
		return
	}
	creds := v.(string)
	// if this is a path and we can stat it, assume it's ok
	if _, err := os.Stat(creds); err == nil {
		return
	}
	if _, err := googleoauth.CredentialsFromJSON(context.Background(), []byte(creds)); err != nil {
		errors = append(errors,
			fmt.Errorf("JSON credentials in %q are not valid: %s", creds, err))
	}

	return
}

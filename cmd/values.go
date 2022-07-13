/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

/* This struc includes clusterDomain, clusterInternalDomain,
spec and imageHub
*/
type ClusterDomain struct {
	ClusterDomain         string
	ClusterInternalDomain string `default:"cluster.local"`
	Spec                  string `default:"allinone"`
	ImageHub              string `default:"docker.io/cnvrg"`
}

type Labels struct {
	Key   string
	Value string
}

type Annotations struct {
	Key   string
	Value string
}

// Parent struct for the Backup values
type Backup struct {
	Enabled  bool
	Rotation int
	Period   string
}

type Gpu struct {
	NvidiaEnable bool
	HabanaEnable bool
}

type Dbs struct {
	CvatEnable        bool
	EsEnable          bool
	EsStorageSize     string
	EsStorageClass    string
	EsPatchNodes      bool
	EsNodeSelector    string
	CleanUpAll        string
	CleanUpApp        string
	CleanUpJobs       string
	CleanUpEndpoints  string
	MinioEnable       bool
	MinioStorageSize  string
	MinioStorageClass string
	MinioNodeSelector string
	PgEnable          bool
	PgStorageSize     string
	PgStorageClass    string
	PgNodeSelector    string
	PgPagesEnable     bool
	PgPagesSize       string
	PgPagesMemory     string
	RedisEnable       bool
	RedisStorageSize  string
	RedisStorageClass string
	RedisNodeSelector string
}

type Logging struct {
	FluentbitEnable    bool
	ElastalertEnable   bool
	ElastaStorageSize  string
	ElastaStorageClass string
	ElastaNodeSelector string
	KibanaEnable       bool
	KibanaSvcName      string
}

//Parent struct for the Capsule values
type Capsule struct {
	Enabled bool
	Image   string `default:"cnvrg-capsule:1.0.2"`
}

// Parent level of ConfigReloader struct
type ConfigReloader struct {
	Enabled bool
}

// Parent level of Registry struct
type Registry struct {
	User     string
	Password string
	Url      string `default:"docker.io"`
	Enabled  bool
}

//Parent level of Tenancy struct
type Tenancy struct {
	Enabled bool
	Key     string
	Value   string
}

// Parent level of SSO struct
type Sso struct {
	Enabled       bool
	AdminUser     string
	Provider      string
	EmailDomain   []string
	ClientId      string
	ClientSecret  string
	AzureTenant   string
	OidcIssuerUrl string
}

// Parent level of Storage struct
type Storage struct {
	Hostpath Hostpath
	Nfs      Nfs
}

// Used in the Storage struct
type Hostpath struct {
	Enabled       bool
	DefaultSc     bool
	Path          string
	ReclaimPolicy string
	NodeSelector  []string // Need to figure out how to define "{ }"
}

// Used in the Storage struct
type Nfs struct {
	Enabled       bool
	Server        string
	Path          string
	DefaultSc     bool
	ReclaimPolicy string
}

type Monitoring struct {
	DcgmExportEnable         bool
	HabanaExportEnable       bool
	NodeExportEnable         bool
	KubeStateMetricEnable    bool
	GrafanaEnable            bool
	GrafanaSvcName           string
	PrometheusOperatorEnable bool
	PrometheusEnable         bool
	PrometheusStorageSize    string
	PrometheusStorageClass   string
	PrometheusNodeSelector   string
	DefaultSvcMonitorsEnable bool
	CnvrgIdleMetricsEnable   bool
	CnvrgIdleMetricsLabels   string
}

// Parent level of the Networking struct
type Networking struct {
	Https   HttpsValues
	Proxy   Proxy
	Ingress Ingress
	Istio   Istio
}

// Used in the Networking struct
type HttpsValues struct {
	Enabled    bool
	CertSecret string `yaml:"certSecret"`
}

// Used in the Networking struct
type Proxy struct {
	Enabled    bool
	HttpProxy  []string
	HttpsProxy []string
	NoProxy    []string
}

// Used in the Networking struct
type Ingress struct {
	Type           string `default:"istio"`
	IstioGwEnabled bool   `default:"true"`
	IstioGwName    string `yaml:"IstioGwName"`
	External       bool
}

// Used in the Networking struct
type Istio struct {
	Enabled               bool
	ExternalIp            []string
	IngressSvcAnnotations string //had {} need to figure out
	IngressSvcExtraPorts  []string
	LbSourceRanges        []string
}

// Template struct for the values.tmpl file
type Template struct {
	ClusterDomain  ClusterDomain
	Labels         Labels
	Annotations    Annotations
	Registry       Registry
	Network        Networking
	Sso            Sso
	Storage        Storage
	Tenancy        Tenancy
	ConfigReloader ConfigReloader
	Capsule        Capsule
	Backup         Backup
	Gpu            Gpu
	Logging        Logging
	Monitoring     Monitoring
	Dbs            Dbs
}

/* function used to leverage the ClusterDomain struct
and to prompt user for all clusterDomain and image settings. This
function will return a struct.
*/
func gatherClusterDomain(cluster *ClusterDomain) {
	fmt.Println("In the gatherClusterDomain func")
	var clusterDomain string
	var clusterInternalDomain string = "cluster.local"

	// Ask what the wildcard domain
	fmt.Print("What is your wildcard domain? ")
	fmt.Scan(&clusterDomain)
	cluster.ClusterDomain = clusterDomain

	fmt.Printf("Do you want to change the internal cluster domain? yes/no [ default is %v ]? ", clusterInternalDomain)
	fmt.Scan(&clusterInternalDomain)
	if clusterInternalDomain == "no" {
		cluster.ClusterInternalDomain = "cluster.local"
	} else {
		cluster.ClusterInternalDomain = clusterInternalDomain
	}

}

/* function used to leverage the Labels struct
and to prompt user for all Labels settings this
will return a struct
*/
func gatherLabels(labels *Labels) {
	fmt.Println("In the gatherLabels func")
	var key string
	var value string

	fmt.Print("Do you want to add a label [format- key: value]? ")
	fmt.Scanf("%s %s", &key, &value)
	labels.Key = key
	labels.Value = value
}

/* function used to leverage the Logging struct
and to prompt user for all Logging settings this
will return a struct
*/
func gatherMonitoring(monitoring *Monitoring) {
	fmt.Println("In the gatherLabels func")
	var disableDcgmExport string
	var disableHabana string
	var disableNodeExport string
	var disableKubeState string
	var disableGrafana string
	var disablePromOperator string
	var disablePrometheus string
	var disableDefaultSvcMonitor string
	var disableCnvrgIDMetrics string

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable dcgm Export? ")
	fmt.Scan(&disableDcgmExport)
	if disableDcgmExport == "yes" {
		monitoring.DcgmExportEnable = false
	}

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable Habana? ")
	fmt.Scan(&disableHabana)
	if disableHabana == "yes" {
		monitoring.HabanaExportEnable = false
	}

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable Kibana? ")
	fmt.Scan(&disableNodeExport)
	if disableNodeExport == "yes" {
		monitoring.NodeExportEnable = false
	}

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable Kibana? ")
	fmt.Scan(&disableKubeState)
	if disableKubeState == "yes" {
		monitoring.KubeStateMetricEnable = false
	}

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable Kibana? ")
	fmt.Scan(&disableGrafana)
	if disableGrafana == "yes" {
		monitoring.GrafanaEnable = false
	}

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable Kibana? ")
	fmt.Scan(&disablePromOperator)
	if disablePromOperator == "yes" {
		monitoring.PrometheusOperatorEnable = false
	}
	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable Kibana? ")
	fmt.Scan(&disablePrometheus)
	if disablePrometheus == "yes" {
		monitoring.PrometheusEnable = false
	}

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable Kibana? ")
	fmt.Scan(&disableDefaultSvcMonitor)
	if disableDefaultSvcMonitor == "yes" {
		monitoring.DefaultSvcMonitorsEnable = false
	}

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable Kibana? ")
	fmt.Scan(&disableCnvrgIDMetrics)
	if disableCnvrgIDMetrics == "yes" {
		monitoring.CnvrgIdleMetricsEnable = false
	}

}

func gatherDbs(dbs *Dbs) {
	fmt.Println("In the gatherLabels func")
	var disableCvat string
	var disableEs string
	var disableMinio string
	var disablePg string

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to enable CVAT? ")
	fmt.Scan(&disableCvat)
	if disableCvat == "yes" {
		dbs.CvatEnable = true
	}

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable Elastic Search? ")
	fmt.Scan(&disableEs)
	if disableEs == "yes" {
		dbs.EsEnable = false
	}

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable Minio? ")
	fmt.Scan(&disableMinio)
	if disableMinio == "yes" {
		dbs.MinioEnable = false
	}

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable Postgres? ")
	fmt.Scan(&disablePg)
	if disablePg == "yes" {
		dbs.PgEnable = false
	}

}

/* function used to leverage the Logging struct
and to prompt user for all Logging settings this
will return a struct
*/
func gatherLogging(logging *Logging) {
	fmt.Println("In the gatherLabels func")
	var disableFluentbit string
	var disableElastalert string
	var disableKibana string

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable Fluentbit? ")
	fmt.Scan(&disableFluentbit)
	if disableFluentbit == "yes" {
		logging.FluentbitEnable = false
	}

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable Elastalert? ")
	fmt.Scan(&disableElastalert)
	if disableElastalert == "yes" {
		logging.ElastalertEnable = false
	}

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable Kibana? ")
	fmt.Scan(&disableKibana)
	if disableKibana == "yes" {
		logging.KibanaEnable = false
	}

}

/* function used to leverage the Gpu struct
and to prompt user for all Gpu settings this
will return a struct
*/
func gatherGpu(gpu *Gpu) {
	fmt.Println("In the gatherLabels func")
	var disableNvidia string
	var disableHabana string

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable Nvidia GPU? ")
	fmt.Scan(&disableNvidia)
	if disableNvidia == "yes" {
		gpu.NvidiaEnable = false
	}

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable Habana GPU? ")
	fmt.Scan(&disableHabana)
	if disableHabana == "yes" {
		gpu.HabanaEnable = false
	}

}

/* function used to leverage the Annotations struct
and to prompt user for all Annotations settings this
will return a struct
*/
func gatherAnnotations(annotations *Annotations) {
	fmt.Println("In the gatherLabels func")
	var key string
	var value string

	fmt.Print("Do you want to add an Annotation [format- key: value]? ")
	fmt.Scanf("%s %s", &key, &value)
	annotations.Key = key
	annotations.Value = value
}

/* function used to leverage the Backup struct
and to prompt user for all Backup settings this
will return a struct
*/
func gatherBackup(backup *Backup) {
	fmt.Println("In the gatherCapsule func")
	var disableBackup string

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable backups? ")
	fmt.Scan(&disableBackup)
	if disableBackup == "yes" {
		backup.Enabled = false
	}
	if disableBackup == "no" {
		backup.Enabled = true
		backup.Rotation = 5
		backup.Period = "24h"
	}
}

/* function used to leverage the Capsule struct
and to prompt user for all Capsule settings this
will return a struct
*/
func gatherCapsule(capsule *Capsule) {
	fmt.Println("In the gatherCapsule func")
	var disableCapsule string

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable capsule? ")
	fmt.Scan(&disableCapsule)
	if disableCapsule == "yes" {
		capsule.Enabled = false
	}
	if disableCapsule == "no" {
		capsule.Enabled = true
	}
}

/* function used to leverate the Tenancy struct
and to prompt user for all Tenancy settings this
will return a struct
*/
func gatherConfigReloader(configReloader *ConfigReloader) {
	fmt.Println("In the gatherTenancy func")
	var enableConfigReloader string

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to enable ConfigReloader? ")
	fmt.Scan(&enableConfigReloader)
	if enableConfigReloader == "no" {
		configReloader.Enabled = false
	}
	if enableConfigReloader == "yes" {
		configReloader.Enabled = true
	}
}

/* function used to leverate the Registry struct
and to prompt user for all Registry settings this
will return a struct
*/
func gatherRegistry(registry *Registry) {
	fmt.Println("In the gatherRegistry func")
	var enableRegistry string

	// Ask if they want to enable SSO skip if "no"
	fmt.Print("Do you want to include specific registry credentials? ")
	fmt.Scan(&enableRegistry)
	if enableRegistry == "no" {
		registry.Enabled = false
	}
	if enableRegistry == "yes" {
		registry.Enabled = true
		registry.Url = "docker.io"
		registry.User = "dockeruser"
		registry.Password = "dockerpassword"
	}
}

/* function used to leverate the Tenancy struct
and to prompt user for all Tenancy settings this
will return a struct
*/
func gatherTenancy(tenancy *Tenancy) {
	fmt.Println("In the gatherTenancy func")
	var enableTenancy string

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to enable Tenancy? ")
	fmt.Scan(&enableTenancy)
	if enableTenancy == "no" {
		tenancy.Enabled = false
	}
	if enableTenancy == "yes" {
		tenancy.Enabled = true
		tenancy.Key = "purpose"
		tenancy.Value = "cnvrg-control-plane"
	}
}

/* function used to leverate the Sso struct
and to prompt user for all Storage settings this
will return a struct
*/
func gatherStorage(storage *Storage) {
	fmt.Println("In the gatherStorage func")
	var enableHostpath string
	var enableNfs string

	// Ask if they want to enable Hostpath for storage skip if "no"
	fmt.Print("Do you want to enable Hostpath for storage? ")
	fmt.Scan(&enableHostpath)
	if enableHostpath == "no" {
		storage.Hostpath.Enabled = false
	}
	if enableHostpath == "yes" {

		storage.Hostpath.Enabled = true
		storage.Hostpath.DefaultSc = false
		storage.Hostpath.Path = "/cnvrg-hostpath-storage"
		storage.Hostpath.ReclaimPolicy = "Retain"
		storage.Hostpath.NodeSelector = []string{}
	}

	// Ask if they want to enable NFS for storage skip if "no"
	fmt.Print("Do you want to enable NFS for storage? ")
	fmt.Scan(&enableNfs)
	if enableNfs == "no" {
		storage.Nfs.Enabled = false
	}
	if enableNfs == "yes" {

		storage.Nfs.Enabled = true
		storage.Nfs.Server = ""
		storage.Nfs.Path = ""
		storage.Nfs.DefaultSc = false
		storage.Nfs.ReclaimPolicy = "Retain"

	}
}

/* function used to leverate the Sso struct
and to prompt user for all SSO settings this
will return a struct
*/
func gatherSso(sso *Sso) {
	fmt.Println("In the gatherSso func")
	var enableSso string

	// Ask if they want to enable SSO skip if "no"
	fmt.Print("Do you want to enable SSO? ")
	fmt.Scan(&enableSso)
	if enableSso == "no" {
		sso.Enabled = false
	}
	if enableSso == "yes" {
		sso.Enabled = true
		sso.AdminUser = ""
		sso.Provider = ""
		sso.EmailDomain = []string{"10.2.3.8,", "192.168.1.5"}
		sso.ClientId = ""
		sso.ClientSecret = ""
		sso.AzureTenant = ""
		sso.OidcIssuerUrl = ""
	}
}

/* function used to leverate the Networking struct
and to prompt user for all networking settings this
will return a struct
*/
func gatherNetworking(network *Networking) {
	fmt.Println("In the gatherNetworking func")
	var enableHttps string
	var enableProxy string
	var externalIngress string
	var diffIngress string

	// Ask if they want to enable https and skip if "no"
	fmt.Print("Do you want to enable https? ")
	fmt.Scan(&enableHttps)
	if enableHttps == "no" {
		network.Https.Enabled = false
	}
	if enableHttps == "yes" {
		network.Https.Enabled = true
		network.Https.CertSecret = "my secret"
	}
	// Ask for Proxy details and skip if answer is "no"
	fmt.Print("Do you want to enable a Proxy? ")
	fmt.Scan(&enableProxy)
	if enableProxy == "yes" {
		network.Proxy.Enabled = true
		network.Proxy.HttpProxy = []string{"hello,", "hows it going"}
		network.Proxy.HttpsProxy = []string{"10.2.3.8,", "192.168.1.5"}
		network.Proxy.NoProxy = []string{"proxy1,", "proxy2"}
	}
	if enableProxy == "no" {
		network.Proxy.Enabled = false
	}
	// Ask for external Istio and skip if answer is "no"
	fmt.Print("Do you have an external istio ingress controller? ")
	fmt.Scan(&externalIngress)
	if externalIngress == "yes" {
		network.Ingress.Type = "istio"
		network.Ingress.IstioGwEnabled = true
		network.Ingress.IstioGwName = "istio-gw"
		network.Ingress.External = true
	}
	if externalIngress == "no" {
		network.Ingress.External = false
	}

	// Ask if they are using a different ingress controll skip if "no"
	fmt.Print("Do you want to disable the istio deployment? ")
	fmt.Scan(&diffIngress)
	if diffIngress == "yes" {
		network.Istio.Enabled = false

	}
	if diffIngress == "no" {
		network.Istio.Enabled = true
		network.Istio.ExternalIp = []string{"10.0.2.5,", "17.1.9.1"}
		network.Istio.IngressSvcAnnotations = "istio-gw"
		network.Istio.IngressSvcExtraPorts = []string{"10.0.2.5,", "17.1.9.1"}
		network.Istio.LbSourceRanges = []string{"10.0.2.5,", "17.1.9.1"}
	}

}

// valuesCmd represents the values command
var valuesCmd = &cobra.Command{
	Use:   "values",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("values called")
		clusterdomain := ClusterDomain{}
		gatherClusterDomain(&clusterdomain)
		labels := Labels{}
		gatherLabels(&labels)
		annotations := Annotations{}
		gatherAnnotations(&annotations)
		registry := Registry{}
		gatherRegistry(&registry)
		network := Networking{}
		gatherNetworking(&network)
		sso := Sso{}
		gatherSso(&sso)
		storage := Storage{}
		gatherStorage(&storage)
		tenancy := Tenancy{}
		gatherTenancy(&tenancy)
		configreloader := ConfigReloader{}
		gatherConfigReloader(&configreloader)
		capsule := Capsule{}
		gatherCapsule(&capsule)
		backup := Backup{}
		gatherBackup(&backup)
		gpu := Gpu{}
		gatherGpu(&gpu)
		logging := Logging{}
		gatherLogging(&logging)
		monitoring := Monitoring{}
		gatherMonitoring(&monitoring)
		dbs := Dbs{}
		gatherDbs(&dbs)

		finaltemp := Template{clusterdomain, labels, annotations, registry, network, sso, storage,
			tenancy, configreloader, capsule, backup, gpu, logging, monitoring, dbs}
		err := temp.Execute(os.Stdout, finaltemp)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// Figure out why this is a thing
var temp *template.Template

func init() {
	createCmd.AddCommand(valuesCmd)
	temp = template.Must(template.ParseFiles("values.tmpl"))
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// valuesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// valuesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

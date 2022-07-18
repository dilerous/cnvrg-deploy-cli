/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

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
	CvatEnable bool

	EsEnable         bool
	EsStorageSize    string
	EsStorageClass   string
	EsPatchNodes     bool
	EsNodeSelector   string
	CleanUpAll       string
	CleanUpApp       string
	CleanUpJobs      string
	CleanUpEndpoints string

	MinioEnable       bool
	MinioStorageSize  string
	MinioStorageClass string
	MinioNodeSelector string

	PgEnable       bool
	PgStorageSize  string
	PgStorageClass string
	PgNodeSelector string
	PgPagesEnable  bool
	PgPagesSize    string
	PgPagesMemory  string

	RedisEnable       bool
	RedisStorageSize  string
	RedisStorageClass string
	RedisNodeSelector string
}

type ControlPlane struct {
	Image string

	BaseConfigAgentTag        string
	BaseConfigIntercom        bool
	BaseConfigFeatureFlags    string
	BaseConfigCnvrgPrivileged bool

	HyperEnable bool

	CnvrgScheduleEnable bool

	CnvrgClusterProvisionerEnable bool

	ObjectStorageType            string
	ObjectStorageBucket          string
	ObjectStorageRegion          string
	ObjectStorageAccessKey       string
	ObjectStorageSecretKey       string
	ObjectStorageEndpoint        string
	ObjectStorageAzureAcountName string
	ObjectStorageAzureContainer  string
	ObjectStorageGcpSecretRef    string
	ObjectStorageGcpProject      string

	SearchkiqEnable         bool
	SearchkiqHpaEnable      bool
	SearchkiqHpaMaxReplicas int

	SidekiqEnable         bool
	SidekiqSplit          bool
	SidekiqHpaEnable      bool
	SidekiqHpaMaxReplicas int

	CnvrgRouterEnable bool
	CnvrgRouterImage  string

	SmtpServer      string
	SmtpPort        int
	SmtpUsername    string
	SmtpPassword    string
	SmtpDomain      string
	SmtpOpenSslMode string
	SmtpSender      string

	SystemkiqEnable         bool
	SystemkiqHpaEnable      bool
	SystemkiqHpaMaxReplicas int

	WebappEnable         bool
	WebappSvcName        string
	WebappReplicas       int
	WebappHpaEnable      bool
	WebappHpaMaxReplicas int

	MpiEnable           bool
	MpiImage            string
	MpiKubectlImage     string
	MpiExtraArgs        string
	MpiRegistryUrl      string
	MpiRegistryUser     string
	MpiRegistryPassword string
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

// Template struct for the values.tmpl file
type Template struct {
	ClusterDomain ClusterDomain
	Labels        Labels
	Annotations   Annotations
	Network       Networking
	/*Registry       Registry
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
	ControlPlane   ControlPlane
	*/
}

/* This struc includes clusterDomain, clusterInternalDomain,
spec and imageHub used with gatherClusterDomain function.
*/
type ClusterDomain struct {
	ClusterDomain         string
	ClusterInternalDomain string `default:"cluster.local"`
	Spec                  string `default:"allinone"`
	ImageHub              string `default:"docker.io/cnvrg"`
}

/* function used to leverage the ClusterDomain struct
and to prompt user for all clusterDomain wildcard dns entry
and if they want to modify the internal cluster domain.
This function will return a struct.
*/
func gatherClusterDomain(cluster *ClusterDomain) {
	log.Println("In the gatherClusterDomain function")
	var clusterDomain string

	// Ask what the wildcard domain is
	fmt.Print("What is your wildcard domain? ")
	fmt.Scan(&clusterDomain)
	cluster.ClusterDomain = clusterDomain

	for {
		consoleReader := bufio.NewReader(os.Stdin)
		fmt.Print("Do you want to change the internal cluster domain? yes/no: ")
		input, _ := consoleReader.ReadString('\n')
		input = strings.ToLower(input)

		if strings.HasPrefix(input, "yes") {
			fmt.Print("Please enter the internal cluster domain: ")
			clusterInput, _ := consoleReader.ReadString('\n')
			clusterInput = strings.ToLower(clusterInput)
			cluster.ClusterInternalDomain = clusterInput
			break
		}
		if strings.HasPrefix(input, "no") {
			cluster.ClusterInternalDomain = "cluster.local"
			break
		}
	}

}

type Labels struct {
	Key       []string
	Stringify string
}

/* function used to leverage the Labels struct
and to prompt user for all Labels settings this
will return a struct
*/
func gatherLabels(labels *Labels) {
	log.Println("In the gatherLabels function")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("To add a Label enter with the format: [ key: value ]")
		scanner.Scan()

		text := scanner.Text()

		if len(text) != 0 {
			labels.Key = append(labels.Key, text)
		} else {
			break
		}
	}
	for _, v := range labels.Key {
		labels.Stringify += fmt.Sprintf("%s, ", v)
	}
}

type Annotations struct {
	Key       []string
	Stringify string
}

/* function used to leverage the Annotations struct
and to prompt user for all Annotations settings this
will return a struct
*/
func gatherAnnotations(annotations *Annotations) {
	log.Println("In the gatherAnnotations function")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("To add an Annotation enter with the format: [format- key: value] ")
		scanner.Scan()

		text := scanner.Text()

		if len(text) != 0 {
			annotations.Key = append(annotations.Key, text)
		} else {
			break
		}
	}

	for _, v := range annotations.Key {
		annotations.Stringify += fmt.Sprintf("%s, ", v)
	}
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

/* function used to leverate the Networking struct
and to prompt user for all networking settings this
will return a struct
*/
func gatherNetworking(network *Networking) {
	log.Println("In the gatherNetworking function")
	var enableProxy string
	var externalIngress string
	var diffIngress string

	for {
		consoleReader := bufio.NewReader(os.Stdin)
		fmt.Println("Do you want to modify Network settings? ")
		fmt.Print("[Settings include; Proxy, Istio Deployment, Ingress or HTTPS.] yes/no: ")
		input, _ := consoleReader.ReadString('\n')
		input = strings.ToLower(input)
		if strings.HasPrefix(input, "no") {
			fmt.Print("Making no changes")
			break
		}
		if strings.HasPrefix(input, "yes") {
			fmt.Println("Press '1' for Proxy Settings")
			fmt.Println("Press '2' for Ingress Settings")
			fmt.Println("Press '3' for HTTPS Settings")
			fmt.Println("Press '4' for Istio Settings")
			fmt.Print("Please make your selection: ")
			var selection int
			fmt.Scan(&selection)
			fmt.Printf("You have selected \n")
			switch selection {
			case 1:
				fmt.Println("In Case 1")
			case 2:
				fmt.Println("In Case 2")
			case 3:
				log.Println("In case statement 3 - HTTPS")
				// Ask if they want to enable https and skip if "no"
				fmt.Print("Do you want to enable HTTPS? ")
				input, _ := consoleReader.ReadString('\n')
				input = strings.ToLower(input)
				if input == "yes\n" {
					network.Https.Enabled = true
					fmt.Printf("The HTTPS network setting is %v", network.Https.Enabled)
				}
				if input == "no" {
					network.Https.Enabled = false
					fmt.Printf("The HTTPS network setting is %v", network.Https.Enabled)
				}
				fmt.Print("Do you want to add a Certificate Secret? ")
				certinput, _ := consoleReader.ReadString('\n')
				certinput = strings.ToLower(certinput)
				fmt.Print(certinput)
				if certinput == "yes" {
					fmt.Print("What do you want to name the secret?")
					var certName string
					fmt.Scan(&certName)
					network.Https.CertSecret = certName
					fmt.Printf("The secret name is %s", certName)
				}
			case 4:
				fmt.Println("In Case 4")
			default:
				fmt.Print("In the default")

			}
			break
		}
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

func gatherControlPlane(controlplane *ControlPlane) {
	fmt.Println("In the gatherLabels func")
	var disableHyper string
	var disableCnvrgScheduler string
	var disableCnvrgClusterProvisioner string
	var disableSearchkiq string
	var disableSidekiq string
	var disableSystemkiq string
	var disableWebapp string
	var disableMpi string

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable Hyper? ")
	fmt.Scan(&disableHyper)
	if disableHyper == "yes" {
		controlplane.HyperEnable = false
	}

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable cnvrg Scheduler? ")
	fmt.Scan(&disableCnvrgScheduler)
	if disableCnvrgScheduler == "yes" {
		controlplane.CnvrgScheduleEnable = false
	}

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable the cnvrg cluster provisioner? ")
	fmt.Scan(&disableCnvrgClusterProvisioner)
	if disableCnvrgClusterProvisioner == "yes" {
		controlplane.CnvrgClusterProvisionerEnable = false
	}

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable Searchkiq? ")
	fmt.Scan(&disableSearchkiq)
	if disableSearchkiq == "yes" {
		controlplane.SearchkiqEnable = false
	}

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable Sidekiq? ")
	fmt.Scan(&disableSidekiq)
	if disableSidekiq == "yes" {
		controlplane.SidekiqEnable = false
	}

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable Searchkiq? ")
	fmt.Scan(&disableSystemkiq)
	if disableSystemkiq == "yes" {
		controlplane.SystemkiqEnable = false
	}

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable Webapp? ")
	fmt.Scan(&disableWebapp)
	if disableWebapp == "yes" {
		controlplane.WebappEnable = false
	}

	// Ask if they want to enable Tenancy skip if "no"
	fmt.Print("Do you want to disable MPI? ")
	fmt.Scan(&disableMpi)
	if disableMpi == "yes" {
		controlplane.MpiEnable = false
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
		log.Println("You are in the values main function")
		clusterdomain := ClusterDomain{}
		gatherClusterDomain(&clusterdomain)
		labels := Labels{}
		gatherLabels(&labels)
		annotations := Annotations{}
		gatherAnnotations(&annotations)
		network := Networking{}
		gatherNetworking(&network)
		/*

			registry := Registry{}
			gatherRegistry(&registry)

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
			controlplane := ControlPlane{}
			gatherControlPlane(&controlplane)
		*/
		finaltemp := Template{clusterdomain, labels, annotations, network} /*registry, sso, storage,
		tenancy, configreloader, capsule, backup, gpu, logging, monitoring, dbs, controlplane */
		err := temp.Execute(os.Stdout, finaltemp)
		if err != nil {
			fmt.Print(err)
		}
	},
}

// Figure out why this is a thing
var temp *template.Template

func init() {
	createCmd.AddCommand(valuesCmd)
	temp = template.Must(template.ParseFiles("values.tmpl"))

}

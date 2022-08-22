/*
Copyright © 2022 BRAD SOPER	BRADLEY.SOPER@CNVRG.IO
*/

package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

// Global Variables
var (
	temp *template.Template

	// Set colors for text
	colorBlue   = "\033[34m"
	colorWhite  = "\033[37m"
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"

	// Set variables for error handling
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	createCmd.AddCommand(valuesCmd)
	temp = template.Must(template.ParseFiles("values.tmpl"))
	// Create and configure a log.txt file to capture all errors and logs
	file, error := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if error != nil {
		log.Fatal(error)
	}
	log.SetOutput(file)

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

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
	Image   string
}

// Parent level of ConfigReloader struct
type ConfigReloader struct {
	Enabled bool
}

// Parent level of Registry struct
type Registry struct {
	User     string
	Password string
	Url      string
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
	EmailDomain   string
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
	NodeSelector  string
}

// Used in the Storage struct
type Nfs struct {
	Enabled       bool
	Server        string
	Path          string
	DefaultSc     bool
	ReclaimPolicy string
	Image         string
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
	ClusterDomain        ClusterDomain
	ClusterInteralDomain ClusterInteralDomain
	Labels               Labels
	Annotations          Annotations
	Network              Networking
	Logging              Logging
	Registry             Registry
	Tenancy              Tenancy
	Sso                  Sso
	Storage              Storage
	ConfigReloader       ConfigReloader
	Capsule              Capsule
	Backup               Backup
	Gpu                  Gpu
	Monitoring           Monitoring
	ControlPlane         ControlPlane
	Dbs                  Dbs
}

/* This struc includes clusterDomain, clusterInternalDomain,
spec and imageHub used with gatherClusterDomain function.
*/
type ClusterDomain struct {
	ClusterDomain string
	Spec          string
	ImageHub      string
}

type ClusterInteralDomain struct {
	Domain string
}

/* function used to leverage the ClusterDomain struct
and to prompt user for all clusterDomain wildcard dns entry
and if they want to modify the internal cluster domain.
*/
func gatherClusterDomain(cluster *ClusterDomain) {
	InfoLogger.Println("In the gatherClusterDomain function")

	// Ask what the wildcard domain is
	fmt.Print((colorBlue), "What is your wildcard domain? ")
	clusterDomain := formatInput()
	cluster.ClusterDomain = clusterDomain

}

func gatherInternalDomain(domain *ClusterInteralDomain) {
	InfoLogger.Println("In the gatherInternalDomain function")

	for {
		fmt.Println((colorBlue), "Press '1' to modify Internal Cluster Domain [default: cluster.local]")
		fmt.Println((colorBlue), "Press '2' to Save and Exit")
		fmt.Print((colorWhite), "Please make your selection: ")
		input := formatInput()
		intVar, _ := strconv.Atoi(input)
		if intVar == 1 {
			fmt.Print((colorWhite), "Please enter the internal cluster domain: ")
			clusterInput := formatInput()
			domain.Domain = clusterInput
			InfoLogger.Printf("Setting the internal cluster domain to %v\n", domain.Domain)
		}
		if intVar == 2 {
			fmt.Println((colorYellow), "Saving and Exiting Internal Domain menu")
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
	InfoLogger.Println("In the gatherLabels function")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print((colorWhite), "Add Label, format [key: value]; 'return' when done: ")
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
	InfoLogger.Println("In the gatherAnnotations function")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print((colorWhite), "Add Annotation, format [key: value]; 'return' when done: ")
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
	CertSecret string
}

// Used in the Networking struct
type Proxy struct {
	Enabled    bool
	HttpProxy  string
	HttpsProxy string
	NoProxy    string
}

// Used in the Networking struct
type Ingress struct {
	Type           string
	IstioGwEnabled bool
	IstioGwName    string
	External       bool
}

// Used in the Networking struct
type Istio struct {
	Enabled               bool
	ExternalIp            string
	IngressSvcAnnotations string
	IngressSvcExtraPorts  string
	LbSourceRanges        string
}

// This function will format strings to lowercase and remove
// any whitespace around the string the value is returned.
func formatInput() string {
	consoleReader := bufio.NewReader(os.Stdin)
	input, _ := consoleReader.ReadString('\n')
	input = strings.ToLower(input)
	input = strings.TrimSpace(input)
	return input
}

// Outputs to std.out the helm commands which need to be ran for installation
func outputHelm() {
	fmt.Println()
	fmt.Println((colorGreen), "---------Helm Repo Commands---------")
	fmt.Println((colorGreen), "Run the following Helm command to install add cnvrg repo")
	fmt.Println((colorWhite), "helm repo add cnvrgv3 https://charts.v3.cnvrg.io")
	fmt.Println((colorWhite), "helm repo update")
	fmt.Println((colorWhite), "helm search repo cnvrgv3/cnvrg -l")
	fmt.Println()
	fmt.Println((colorGreen), "---------Helm Install Command---------")
	fmt.Println((colorGreen), "Run the following Helm command to install cnvrg.io")
	fmt.Println((colorWhite), "helm install cnvrg cnvrgv3/cnvrg --create-namespace -n cnvrg --timeout 1500s --wait --values ./values.yaml")
}

// The function prompts for a key value value
// Takes the key value and returns a string
func createArray() string {
	InfoLogger.Println("In the createArray function")
	scanner := bufio.NewScanner(os.Stdin)
	var key []string
	var stringify string

	for {
		fmt.Print((colorWhite), "Format [key: value]; 'return' when done: ")
		scanner.Scan()
		text := scanner.Text()

		if len(text) != 0 {
			key = append(key, text)
		} else {
			break
		}
	}
	for _, v := range key {
		stringify += fmt.Sprintf("%s, ", v)

	}
	return stringify

}

// This function will return a slice as a string. You can enter
// any number of values one line at a time.
func createSlice() string {
	InfoLogger.Println("In the createSlice function")
	fmt.Println((colorWhite), "Enter 1 value per line. Press 'return' when done: ")
	consoleScanner := bufio.NewScanner(os.Stdin)
	var slice []string
	var finalSlice string

	for {
		consoleScanner.Scan()
		text := consoleScanner.Text()
		if len(text) != 0 {
			slice = append(slice, text)
		} else {
			break
		}
	}
	for _, v := range slice {
		finalSlice += fmt.Sprintf("%s, ", v)
	}
	return finalSlice
}

/* function used to leverate the Networking struct
and to prompt user for all networking settings this
will return a struct
*/
func gatherNetworking(network *Networking) {
	InfoLogger.Println("In the gatherNetworking function")

	for {
		fmt.Println()
		fmt.Println((colorGreen), "----Networking Menu----")
		fmt.Println((colorGreen), "Update Networking values")
		fmt.Println((colorBlue), "Press '1' for Proxy Settings")
		fmt.Println((colorBlue), "Press '2' for Ingress Settings")
		fmt.Println((colorBlue), "Press '3' for HTTPS Settings")
		fmt.Println((colorBlue), "Press '4' for Istio Settings")
		fmt.Println((colorBlue), "Press '5' to Save and Exit Network Menu")
		fmt.Print((colorWhite), "Please make your selection: ")
		caseInput := formatInput()
		intVar, _ := strconv.Atoi(caseInput)
		switch intVar {
		case 1:
			InfoLogger.Println("In case statement 1 - Proxy")
			for {
				fmt.Println()
				fmt.Println((colorGreen), "----Proxy Menu----")
				fmt.Println((colorGreen), "Update Proxy values")
				fmt.Println((colorBlue), "Press '1' to enable Proxy")
				fmt.Println((colorBlue), "Press '2' to input HTTP proxies to use")
				fmt.Println((colorBlue), "Press '3' to input HTTPS proxies to use")
				fmt.Println((colorBlue), "Press '4' to input extra No Proxy values to use")
				fmt.Println((colorBlue), "Press '5' to Save and Exit Proxy settings")
				fmt.Print((colorWhite), "Please make your selection: ")
				caseInput := formatInput()
				intVar, _ := strconv.Atoi(caseInput)
				switch intVar {
				case 1:
					network.Proxy.Enabled = true
					fmt.Println((colorYellow), "Proxy enabled")
					InfoLogger.Printf("Network Proxy set to %v\n", network.Proxy.Enabled)
				case 2:
					fmt.Println((colorBlue), "Please enter a list of HTTP proxies")
					slice := createSlice()
					network.Proxy.HttpProxy = slice
					network.Proxy.Enabled = true
				case 3:
					fmt.Println((colorBlue), "Please enter a list of HTTPS proxies")
					slice := createSlice()
					network.Proxy.HttpsProxy = slice
					network.Proxy.Enabled = true
				case 4:
					fmt.Println((colorBlue), "Please enter a list of No proxies")
					slice := createSlice()
					network.Proxy.NoProxy = slice
					network.Proxy.Enabled = true
				}
				if intVar == 5 {
					fmt.Println((colorYellow), "Saving and Exiting Proxy section")
					break
				}
			}
		case 2:
			InfoLogger.Println("In Case statement 2 - Ingress")
			for {
				fmt.Println()
				fmt.Println((colorGreen), "----Ingress Menu----")
				fmt.Println((colorGreen), "Update Ingress values")
				fmt.Println((colorBlue), "Press '1' to modify Ingress Type")
				fmt.Println((colorBlue), "Press '2' to Save and Exit Ingress Menu")
				fmt.Print((colorWhite), "Please make your selection: ")
				caseInput := formatInput()
				intVar, _ := strconv.Atoi(caseInput)
				switch intVar {
				case 1:
					fmt.Print((colorWhite), "What is the ingress type [istio|ingress|openshift|nodeport]?: ")
					ingressType := formatInput()
					if ingressType == "istio" {
						for {
							fmt.Println((colorGreen), "----Istio Menu----")
							fmt.Println((colorGreen), "Update Istio values")
							fmt.Println((colorBlue), "Press '1' to modify External IP")
							fmt.Println((colorBlue), "Press '2' to modify Service Annotations")
							fmt.Println((colorBlue), "Press '3' to modify Service Extra Ports")
							fmt.Println((colorBlue), "Press '4' to modify Load Balance Source Ranges")
							fmt.Println((colorBlue), "Press '5' to Save and Exit")
							fmt.Print((colorWhite), "Please make your selection: ")
							caseInput := formatInput()
							intVar, _ := strconv.Atoi(caseInput)
							switch intVar {
							case 1:
								fmt.Print((colorWhite), "Input External IPs")
								input := createSlice()
								network.Istio.ExternalIp = input
							case 2:
								fmt.Print((colorWhite), "Input Service Annotations")
								input := createSlice()
								network.Istio.IngressSvcAnnotations = input
							case 3:
								fmt.Print((colorWhite), "Input Service Extra Ports")
								input := createSlice()
								network.Istio.IngressSvcExtraPorts = input
							case 4:
								fmt.Print((colorWhite), "Input Load Balance Source Ranges")
								input := createSlice()
								network.Istio.LbSourceRanges = input
							}
							if intVar == 5 {
								fmt.Println((colorYellow), "Saving and Exiting Istio menu")
								break
							}
						}
					}
					if ingressType == "ingress" {
						network.Ingress.Type = "ingress"
						network.Istio.Enabled = false
						fmt.Printf("Set Ingress to '%v' and Disabled Istio\n", ingressType)

					}
					if ingressType == "nodeport" {
						network.Ingress.Type = "nodeport"
						network.Istio.Enabled = false
						fmt.Printf("Set Ingress to '%v' and Disabled Istio\n", ingressType)
					}
				}
				if intVar == 2 {
					fmt.Println((colorYellow), "Saving and Exiting Ingress menu")
					break
				}
			}
		case 3:
			InfoLogger.Println("In case statement 3 - HTTPS")
			for {
				// Ask if they want to enable https and skip if "no"
				fmt.Print("Do you want to enable HTTPS? (yes/no): ")
				caseThreeInput := formatInput()

				if caseThreeInput == "yes" {
					network.Https.Enabled = true
					InfoLogger.Printf("The HTTPS network setting is %v\n", network.Https.Enabled)
					break
				}
				if caseThreeInput == "no" {
					network.Https.Enabled = false
					InfoLogger.Printf("The HTTPS network setting is %v\n", network.Https.Enabled)
					break
				}
			}
			for {
				fmt.Print((colorWhite), "Do you want to add a Certificate? (yes/no)")
				certinput := formatInput()
				if certinput == "yes" {
					fmt.Print("What do you want to name the Certificate secret? ")
					var certName string
					fmt.Scan(&certName)
					network.Https.CertSecret = certName
					network.Https.Enabled = true
					InfoLogger.Printf("The secret name is %s \n", certName)
					break
				}
				if certinput == "no" {
					InfoLogger.Println("Breaking for loop, not setting Certificate name")
					break
				}
			}
		case 4:
			for {
				fmt.Println()
				fmt.Println((colorGreen), "----Istio Menu----")
				fmt.Println((colorGreen), "Update Istio values")
				fmt.Println((colorBlue), "Press '1' to disable Istio")
				fmt.Println((colorBlue), "Press '2' list IPs to use for istio ingress service")
				fmt.Println((colorBlue), "Press '3' list extra ports for istio ingress service")
				fmt.Println((colorBlue), "Press '4' list extra LB sources ranges")
				fmt.Println((colorBlue), "Press '5' map of strings for Istio SVC annotations")
				fmt.Println((colorBlue), "Press '6' to Save and Exit Istio menu")
				fmt.Print((colorWhite), "Please make your selection: ")
				caseInput := formatInput()
				intVar, _ := strconv.Atoi(caseInput)
				switch intVar {
				case 1:
					network.Istio.Enabled = false
					fmt.Println((colorYellow), "Istio is disabled")
					InfoLogger.Printf("Istio set to %v", network.Istio.Enabled)
				case 2:
					fmt.Println((colorWhite), "Please enter a list of IPs to use for Istio ingress service: ")
					slice := createSlice()
					network.Istio.ExternalIp = slice
					network.Istio.Enabled = true
				case 3:
					fmt.Println((colorWhite), "Please enter a list extra ports for Istio ingress service: ")
					slice := createSlice()
					network.Istio.IngressSvcExtraPorts = slice
					network.Istio.Enabled = true
				case 4:
					fmt.Println((colorWhite), "Please enter a list of extra LB sources ranges: ")
					slice := createSlice()
					network.Istio.LbSourceRanges = slice
					network.Istio.Enabled = true
				case 5:
					fmt.Println((colorWhite), "Please enter Istio SVC annotations: ")
					slice := createArray()
					network.Istio.IngressSvcAnnotations = slice
					network.Istio.Enabled = true
				}
				if intVar == 6 {
					fmt.Println((colorYellow), "Saving and Exiting Istio menu")
					break
				}
			}
		}
		if intVar == 5 {
			fmt.Println((colorYellow), "Saving and Exiting Networking menu")
			break
		}
	}
}

/* function used to leverage the Logging struct
and to prompt user for all Logging settings this
will return a struct
*/
func gatherMonitoring(monitoring *Monitoring) {
	InfoLogger.Println("In the gatherMonitoring function")

	for {
		fmt.Println()
		fmt.Println((colorGreen), "----Monitoring Menu----")
		fmt.Println((colorGreen), "Update Monitoring values")
		fmt.Println((colorBlue), "Press '1' To disable dcgm Export Monitoring")
		fmt.Println((colorBlue), "Press '2' To disable Habana Monitoring")
		fmt.Println((colorBlue), "Press '3' To disable Node Export Monitoring")
		fmt.Println((colorBlue), "Press '4' To disable Kube State Metric Monitoring")
		fmt.Println((colorBlue), "Press '5' To disable Grafana Monitoring")
		fmt.Println((colorBlue), "Press '6' To disable the Prometheus Operator")
		fmt.Println((colorBlue), "Press '7' To disable Prometheus")
		fmt.Println((colorBlue), "Press '8' To disable Default Svc Monitoring")
		fmt.Println((colorBlue), "Press '9' To disable cnvrg Idle Metrics")
		fmt.Println((colorBlue), "Press '10' To Save and Exit")
		fmt.Print((colorWhite), "Please make your selection: ")
		caseInput := formatInput()
		intVar, _ := strconv.Atoi(caseInput)
		switch intVar {
		case 1:
			monitoring.DcgmExportEnable = false
			fmt.Println((colorYellow), "DCGM Export Disabled")
		case 2:
			monitoring.HabanaExportEnable = false
			fmt.Println((colorYellow), "Habana Export Disabled")
		case 3:
			monitoring.NodeExportEnable = false
			fmt.Println((colorYellow), "Node Export Disabled")
		case 4:
			monitoring.KubeStateMetricEnable = false
			fmt.Println((colorYellow), "Kube State Metrics Disabled")
		case 5:
			monitoring.GrafanaEnable = false
			fmt.Println((colorYellow), "Grafana Disabled")
		case 6:
			monitoring.PrometheusOperatorEnable = false
			fmt.Println((colorYellow), "Prometheus Operator Disabled")
		case 7:
			monitoring.PrometheusEnable = false
			fmt.Println((colorYellow), "Prometheus Disabled")
		case 8:
			monitoring.DefaultSvcMonitorsEnable = false
			fmt.Println((colorYellow), "Default svc Monitor Disabled")
		case 9:
			monitoring.CnvrgIdleMetricsEnable = false
			fmt.Println((colorYellow), "cnvrg Idle Metrics Disabled")
		}
		if intVar == 10 {
			fmt.Println((colorYellow), "Saving and Exiting Menu")
			break
		}
	}
}

func gatherControlPlane(controlplane *ControlPlane) {
	InfoLogger.Println("In the gatherControlPlane function")

	for {
		fmt.Println()
		fmt.Println((colorGreen), "----ControlPlane Menu----")
		fmt.Println((colorGreen), "Update ControlPlane values")
		fmt.Println((colorBlue), "Press '1' To disable Hyper")
		fmt.Println((colorBlue), "Press '2' To disable cnvrg Scheduler")
		fmt.Println((colorBlue), "Press '3' To disable cnvrg Cluster Provisioner")
		fmt.Println((colorBlue), "Press '4' To disable Searchkiq")
		fmt.Println((colorBlue), "Press '5' To disable Sidekiq")
		fmt.Println((colorBlue), "Press '6' To disable Systemkiq")
		fmt.Println((colorBlue), "Press '7' To disable Webapp")
		fmt.Println((colorBlue), "Press '8' To disable MPI")
		fmt.Println((colorBlue), "Press '9' To Save and Exit")
		fmt.Print((colorWhite), "Please make your selection: ")
		caseInput := formatInput()
		intVar, _ := strconv.Atoi(caseInput)
		switch intVar {
		case 1:
			controlplane.HyperEnable = false
			fmt.Println((colorYellow), "Hyper Disabled")
		case 2:
			controlplane.CnvrgScheduleEnable = false
			fmt.Println((colorYellow), "cnvrg.io Scheduler Disabled")
		case 3:
			controlplane.CnvrgClusterProvisionerEnable = false
			fmt.Println((colorYellow), "cnvrg.io Cluster Provisioner Disabled")
		case 4:
			controlplane.SearchkiqEnable = false
			fmt.Println((colorYellow), "Searchkiq Disabled")
		case 5:
			controlplane.SidekiqEnable = false
			fmt.Println((colorYellow), "Sidekiq Disabled")
		case 6:
			controlplane.SystemkiqEnable = false
			fmt.Println((colorYellow), "Systemkiq Disabled")
		case 7:
			controlplane.WebappEnable = false
			fmt.Println((colorYellow), "Webapp Disabled")
		case 8:
			controlplane.MpiEnable = false
			fmt.Println((colorYellow), "MPI Disabled")

		}
		if intVar == 9 {
			fmt.Println((colorYellow), "Saving and Exiting ControlPlane Settings")
			break
		}
	}
}

/* Function used to gather Database values through menu driven options.
This function used the Dbs struct
*/
func gatherDbs(dbs *Dbs) {
	InfoLogger.Println("In the gatherLabels func")

	for {
		fmt.Println()
		fmt.Println((colorGreen), "----Database Menu----")
		fmt.Println((colorGreen), "Update Database values")
		fmt.Println((colorBlue), "Press '1' To enable CVAT")
		fmt.Println((colorBlue), "Press '2' To modify Elastic Search")
		fmt.Println((colorBlue), "Press '3' To modify Minio")
		fmt.Println((colorBlue), "Press '4' To modify Postgres")
		fmt.Println((colorBlue), "Press '5' To modify Redis")
		fmt.Println((colorBlue), "Press '6' To Save and Exit")
		fmt.Print((colorWhite), "Please make your selection: ")
		caseInput := formatInput()
		intVar, _ := strconv.Atoi(caseInput)
		switch intVar {
		case 1:
			dbs.CvatEnable = true
			fmt.Println((colorYellow), "CVAT enabled")
		case 2:
			for {
				fmt.Println()
				fmt.Println((colorGreen), "----Elastic Search Menu----")
				fmt.Println((colorGreen), "Update Elastic Search values")
				fmt.Println((colorBlue), "Press '1' To disable Elastic Search")
				fmt.Println((colorBlue), "Press '2' To modify Storage Size [default: 80Gi]")
				fmt.Println((colorBlue), "Press '3' To modify Storage Class")
				fmt.Println((colorBlue), "Press '4' To disable Patch Elastic Search Nodes")
				fmt.Println((colorBlue), "Press '5' To modify Node Selector")
				fmt.Println((colorBlue), "Press '6' To Save and Exit Elastic Search menu")
				fmt.Print((colorWhite), "Please make your selection: ")
				caseInput := formatInput()
				intVar, _ := strconv.Atoi(caseInput)
				switch intVar {
				case 1:
					dbs.EsEnable = false
					fmt.Println((colorYellow), "Elastic Search disabled")
				case 2:
					fmt.Print((colorWhite), "Input Storage Size [default: 80Gi]: ")
					caseInput := formatInput()
					dbs.EsStorageSize = caseInput + "Gi"
					dbs.EsEnable = true
				case 3:
					fmt.Print((colorWhite), "Input Storage Class: ")
					caseInput := formatInput()
					dbs.EsStorageSize = caseInput
					dbs.EsEnable = true
				case 4:
					dbs.EsPatchNodes = false
					fmt.Println((colorYellow), "Elastic Search Patch Nodes disabled")
					dbs.EsEnable = true
				case 5:
					fmt.Print((colorWhite), "Input Node Selector values")
					node := createSlice()
					dbs.EsNodeSelector = node
					dbs.EsEnable = true
				}
				if intVar == 6 {
					fmt.Println((colorYellow), "Saving and Exiting Elastic Search settings")
					break
				}
			}

		case 3:
			for {
				fmt.Println()
				fmt.Println((colorGreen), "----Minio Menu----")
				fmt.Println((colorGreen), "Update Minio values")
				fmt.Println((colorBlue), "Press '1' To disable Minio")
				fmt.Println((colorBlue), "Press '2' To modify Storage Size [default: 100Gi]")
				fmt.Println((colorBlue), "Press '3' To modify Storage Class")
				fmt.Println((colorBlue), "Press '4' To modify Node Selector")
				fmt.Println((colorBlue), "Press '5' To Save and Exit Elastic Search menu")
				fmt.Print((colorWhite), "Please make your selection: ")
				caseInput := formatInput()
				intVar, _ := strconv.Atoi(caseInput)
				switch intVar {
				case 1:
					dbs.MinioEnable = false
					fmt.Println((colorYellow), "Minio disabled")
				case 2:
					fmt.Print((colorWhite), "Input Storage Size [default: 100Gi]: ")
					caseInput := formatInput()
					dbs.MinioStorageSize = caseInput + "Gi"
				case 3:
					fmt.Print((colorWhite), "Input Storage Class: ")
					caseInput := formatInput()
					dbs.MinioStorageClass = caseInput
				case 4:
					fmt.Print((colorWhite), "Input Node Selector values")
					node := createSlice()
					dbs.MinioNodeSelector = node
				}
				if intVar == 5 {
					fmt.Println((colorYellow), "Saving and Exiting Minio settings")
					break
				}
			}
		case 4:
			for {
				fmt.Println()
				fmt.Println((colorGreen), "----Postgres Menu----")
				fmt.Println((colorGreen), "Update Postgres values")
				fmt.Println((colorBlue), "Press '1' To disable Postgres")
				fmt.Println((colorBlue), "Press '2' To modify Storage Size [default: 100Gi]")
				fmt.Println((colorBlue), "Press '3' To modify Storage Class")
				fmt.Println((colorBlue), "Press '4' To modify Node Selector")
				fmt.Println((colorBlue), "Press '5' To Save and Exit Postgres menu")
				fmt.Print((colorWhite), "Please make your selection: ")
				caseInput := formatInput()
				intVar, _ := strconv.Atoi(caseInput)
				switch intVar {
				case 1:
					dbs.PgEnable = false
					fmt.Println((colorYellow), "Postgres disabled")
				case 2:
					fmt.Print((colorWhite), "Input Storage Size [default: 80Gi]: ")
					caseInput := formatInput()
					dbs.PgStorageSize = caseInput + "Gi"
				case 3:
					fmt.Print((colorWhite), "Input Storage Class: ")
					caseInput := formatInput()
					dbs.PgStorageClass = caseInput
				case 4:
					fmt.Print((colorWhite), "Input Node Selector values")
					node := createSlice()
					dbs.PgNodeSelector = node
				}
				if intVar == 5 {
					fmt.Println((colorYellow), "Saving and Exiting Postgres settings")
					break
				}
			}
		case 5:
			for {
				fmt.Println()
				fmt.Println((colorGreen), "----Redis Menu----")
				fmt.Println((colorGreen), "Update Redis values")
				fmt.Println((colorBlue), "Press '1' To disable Redis")
				fmt.Println((colorBlue), "Press '2' To modify Storage Size [default: 10Gi]")
				fmt.Println((colorBlue), "Press '3' To modify Storage Class")
				fmt.Println((colorBlue), "Press '4' To modify Node Selector")
				fmt.Println((colorBlue), "Press '5' To Save and Exit Redis menu")
				fmt.Print((colorWhite), "Please make your selection: ")
				caseInput := formatInput()
				intVar, _ := strconv.Atoi(caseInput)
				switch intVar {
				case 1:
					dbs.RedisEnable = false
					fmt.Println((colorYellow), "Postgres Redis")
				case 2:
					fmt.Print((colorWhite), "Input Storage Size [default: 10Gi]: ")
					caseInput := formatInput()
					dbs.RedisStorageSize = caseInput + "Gi"
				case 3:
					fmt.Print((colorWhite), "Input Storage Class: ")
					caseInput := formatInput()
					dbs.RedisStorageClass = caseInput
				case 4:
					fmt.Print((colorWhite), "Input Node Selector values")
					node := createSlice()
					dbs.RedisNodeSelector = node
				}
				if intVar == 5 {
					fmt.Println((colorYellow), "Saving and Exiting Redis settings")
					break
				}
			}

		}
		if intVar == 6 {
			fmt.Println((colorYellow), "Saving and Exiting Database Settings")
			break
		}
	}
}

/* Function used to gather Logging values through menu driven options.
This function uses the Logging struct
*/
func gatherLogging(logging *Logging) {
	InfoLogger.Println("In the gatherLogging function")

	for {
		fmt.Println()
		fmt.Println((colorGreen), "----Logging Menu----")
		fmt.Println((colorGreen), "Update Logging values")
		fmt.Println((colorBlue), "Press '1' To disable Fluentbit")
		fmt.Println((colorBlue), "Press '2' To disable Kibana")
		fmt.Println((colorBlue), "Press '3' To configure Elastalert")
		fmt.Println((colorBlue), "Press '4' To Save and Exit")
		fmt.Print((colorWhite), "Please make your selection: ")
		caseInput := formatInput()
		intVar, _ := strconv.Atoi(caseInput)
		switch intVar {
		case 1:
			logging.FluentbitEnable = false
			fmt.Println((colorYellow), "Fluentbit is disabled")
			InfoLogger.Printf("Fluentbit Enabled set to %v\n", logging.FluentbitEnable)
		case 2:
			logging.KibanaEnable = false
			fmt.Println((colorYellow), "Kibana is disabled")
			InfoLogger.Printf("Kibana Enabled set to %v\n", logging.KibanaEnable)
		case 3:
			for {
				fmt.Println()
				fmt.Println((colorGreen), "----Elastic Search Menu----")
				fmt.Println((colorGreen), "Update Elastic Search values")
				fmt.Println((colorBlue), "Press '1' to disable Elastalert")
				fmt.Println((colorBlue), "Press '2' to change the Storage Size")
				fmt.Println((colorBlue), "Press '3' to change the Storage Class")
				fmt.Println((colorBlue), "Press '4' to change the node Selector")
				fmt.Println((colorBlue), "Press '5' to Save and Exit")
				fmt.Print((colorWhite), "Please make your selection: ")
				caseInput := formatInput()
				intVar, _ := strconv.Atoi(caseInput)
				switch intVar {
				case 1:
					logging.ElastalertEnable = false
					fmt.Println((colorYellow), "Elastalert is disabled")
				case 2:
					fmt.Print((colorWhite), "Input Storage Size [default: 30Gi]: ")
					var storageSize string
					fmt.Scan(&storageSize)
					logging.ElastaStorageSize = storageSize + "Gi"
					logging.ElastalertEnable = true
				case 3:
					fmt.Print((colorWhite), "Please enter the new Storage Class: ")
					var storageClass string
					fmt.Scan(&storageClass)
					logging.ElastaStorageClass = storageClass
					logging.ElastalertEnable = true
				case 4:
					fmt.Print((colorWhite), "Please enter the new Node Selector: ")
					storageClass := createSlice()
					logging.ElastaNodeSelector = storageClass
					logging.ElastalertEnable = true
				}
				if intVar == 5 {
					fmt.Println((colorYellow), "Saving and Exiting Elastic Search menu")
					break
				}
			}
		}
		if intVar == 4 {
			fmt.Println((colorYellow), "Saving and Exiting Logging menu")
			break
		}
	}
}

/* Function used to gather GPU values through menu driven options.
This function uses the Gpu struct
*/
func gatherGpu(gpu *Gpu) {
	InfoLogger.Println("In the gatherGpu function")

	for {
		fmt.Println()
		fmt.Println((colorGreen), "----GPU Menu----")
		fmt.Println((colorGreen), "Update GPU values")
		fmt.Println((colorBlue), "Press '1' to disable Nvidia GPU")
		fmt.Println((colorBlue), "Press '2' to disable Habana GPU")
		fmt.Println((colorBlue), "Press '3' to Save and Exit")
		fmt.Print((colorWhite), "Please make your selection: ")
		caseInput := formatInput()
		intVar, _ := strconv.Atoi(caseInput)
		switch intVar {
		case 1:
			gpu.NvidiaEnable = false
			fmt.Println((colorYellow), "Nvidia GPU is disabled")
		case 2:
			gpu.HabanaEnable = false
			fmt.Println((colorYellow), "Habana GPU is disabled")
		}
		if intVar == 3 {
			fmt.Println((colorYellow), "Saving and Exiting GPU menu")
			break
		}
	}
}

/* Function used to gather Backup values through menu driven options.
This function uses the Backup struct
*/
func gatherBackup(backup *Backup) {
	InfoLogger.Println("In the gatherBackup function")

	for {
		fmt.Println()
		fmt.Println((colorGreen), "----Backup Menu----")
		fmt.Println((colorGreen), "Update Backup values")
		fmt.Println((colorBlue), "Press '1' to disable Backups")
		fmt.Println((colorBlue), "Press '2' to modify Backup Rotation [default: 5]")
		fmt.Println((colorBlue), "Press '3' to modify Backup Period [default: 24h]")
		fmt.Println((colorBlue), "Press '4' to Save and Exit")
		fmt.Print((colorWhite), "Please make your selection: ")
		caseInput := formatInput()
		intVar, _ := strconv.Atoi(caseInput)
		switch intVar {
		case 1:
			backup.Enabled = false
			fmt.Println((colorYellow), "Backup is disabled")
		case 2:
			fmt.Print((colorBlue), "Input Backup Rotation [default: 5]: ")
			caseInput := formatInput()
			intVar, _ := strconv.Atoi(caseInput)
			backup.Rotation = intVar
		case 3:
			fmt.Print((colorBlue), "Input Backup Period [default: 24h]: ")
			caseInput := formatInput()
			if caseInput == "" {
				backup.Period = "24h"
			} else {
				backup.Period = caseInput
			}
		}
		if intVar == 4 {
			fmt.Println((colorYellow), "Saving and Exiting Backup menu")
			break
		}
	}

}

/* Function used to gather Capsule values through menu driven options.
This function uses the Capsule struct
*/
func gatherCapsule(capsule *Capsule) {
	InfoLogger.Println("In the gatherCapsule function")

	for {
		fmt.Println()
		fmt.Println((colorGreen), "----Capsule Menu----")
		fmt.Println((colorGreen), "Update Capsule values")
		fmt.Println((colorBlue), "Press '1' to disable Capsule")
		fmt.Println((colorBlue), "Press '2' to modify Capsule image")
		fmt.Println((colorBlue), "Press '3' to Save and Exit")
		fmt.Print((colorWhite), "Please make your selection: ")
		caseInput := formatInput()
		intVar, _ := strconv.Atoi(caseInput)
		switch intVar {
		case 1:
			capsule.Enabled = false
			fmt.Println((colorYellow), "Capsule is disabled")
		case 2:
			fmt.Print((colorBlue), "Please enter new image: ")
			caseInput := formatInput()
			capsule.Image = caseInput
		}
		if intVar == 3 {
			fmt.Println((colorYellow), "Saving and Exiting Capsule menu")
			break
		}
	}
}

/* function used to leverate the ConfigReloader struct
and to prompt user for all ConfigReloader settings
*/
func gatherConfigReloader(configReloader *ConfigReloader) {
	InfoLogger.Println("In the gatherConfigReloader func")

	fmt.Println((colorYellow), "Config Reload is disabled")
	configReloader.Enabled = false
}

/* function used to leverate the Registry struct
and to prompt user for all Registry settings this
will return a struct
*/
func gatherRegistry(registry *Registry) {
	InfoLogger.Println("In the gatherRegistry function")

	var password string

	for {
		fmt.Println()
		fmt.Println((colorGreen), "----Registry Menu----")
		fmt.Println((colorGreen), "Update Registry values")
		fmt.Println((colorBlue), "Press '1' to update Registry URL")
		fmt.Println((colorBlue), "Press '2' to update Registry User Name")
		fmt.Println((colorBlue), "Press '3' to update Registry Password")
		fmt.Println((colorBlue), "Press '4' to Save and Exit")
		fmt.Print((colorWhite), "Please make your selection: ")
		caseInput := formatInput()
		intVar, _ := strconv.Atoi(caseInput)
		switch intVar {
		case 1:
			fmt.Print((colorBlue), "Input the registry URL [default docker.io]: ")
			url := formatInput()
			if url == "" {
				registry.Url = "docker.io"
			} else {
				registry.Enabled = true
				registry.Url = url
			}
		case 2:
			fmt.Print("Input the registry User Name: ")
			user := formatInput()
			registry.User = user
			registry.Enabled = true
		case 3:
			fmt.Print("Input the registry Password: ")
			fmt.Scanln(&password)
			registry.Password = password
			registry.Enabled = true
		}
		if intVar == 4 {
			fmt.Println((colorYellow), "Saving and Exiting Registry menu")
			break
		}
	}
}

/* function used to leverate the Tenancy struct
and to prompt user for all Tenancy settings this
will return a struct
*/
func gatherTenancy(tenancy *Tenancy) {
	InfoLogger.Println("In the gatherTenancy function")

	for {
		fmt.Println()
		fmt.Println((colorGreen), "----Tenancy Menu----")
		fmt.Println((colorGreen), "Update Tenancy values")
		fmt.Println((colorBlue), "Press '1' to enable Tenancy")
		fmt.Println((colorBlue), "Press '2' to add Tenancy node selector key")
		fmt.Println((colorBlue), "Press '3' to add Tenancy node selector value")
		fmt.Println((colorBlue), "Press '4' to Save and Exit")
		fmt.Print((colorWhite), "Please make your selection: ")
		caseInput := formatInput()
		intVar, _ := strconv.Atoi(caseInput)
		switch intVar {
		case 1:
			tenancy.Enabled = true
			fmt.Println((colorYellow), "Tenancy Enabled")
			InfoLogger.Printf("Tenancy enabled set to %v\n", tenancy.Enabled)
		case 2:
			fmt.Print((colorBlue), "Please enter the Tenancy node selector key: ")
			key := formatInput()
			tenancy.Key = key
			tenancy.Enabled = true
		case 3:
			fmt.Print((colorBlue), "Please enter the Tenancy node selector value: ")
			value := formatInput()
			tenancy.Value = value
			tenancy.Enabled = true
		}
		if intVar == 4 {
			fmt.Println((colorYellow), "Saving and Exiting Tenancy menu")
			break
		}
	}
}

/* function used to leverate the Sso struct
and to prompt user for all Storage settings this
will return a struct
*/
func gatherStorage(storage *Storage) {
	InfoLogger.Println("In the gatherStorage function")

	for {
		fmt.Println()
		fmt.Println((colorGreen), "----Storage Menu----")
		fmt.Println((colorGreen), "Update Storage values")
		fmt.Println((colorBlue), "Press '1' to modify HostPath settings")
		fmt.Println((colorBlue), "Press '2' to modify NFS settings")
		fmt.Println((colorBlue), "Press '3' to Save and Exit")
		fmt.Print((colorWhite), "Please make your selection: ")
		caseInput := formatInput()
		intVar, _ := strconv.Atoi(caseInput)
		switch intVar {
		case 1:
			for {
				fmt.Println()
				fmt.Println((colorGreen), "----HostPath Menu----")
				fmt.Println((colorGreen), "Update HostPath values")
				fmt.Println((colorBlue), "Press '1' to set HostPath as the Default Storage Class")
				fmt.Println((colorBlue), "Press '2' to modify the Path")
				fmt.Println((colorBlue), "Press '3' to modify Reclaim Policy [default: Retain]")
				fmt.Println((colorBlue), "Press '4' to modify Node Selector")
				fmt.Println((colorBlue), "Press '5' to Save and Exit")
				fmt.Print((colorWhite), "Please make your selection: ")
				caseInput := formatInput()
				intVar, _ := strconv.Atoi(caseInput)
				switch intVar {
				case 1:
					storage.Hostpath.Enabled = true
					storage.Hostpath.DefaultSc = true
					fmt.Println((colorYellow), "HostPath set as default Storage Class")
				case 2:
					fmt.Print((colorBlue), "Input the path [default: /cnvrg-hostpath-storage]: ")
					caseInput := formatInput()
					storage.Hostpath.Path = caseInput
					storage.Hostpath.Enabled = true
				case 3:
					var input string
					var policy = []string{"Retain", "Delete", "Recycle"}
					done := true
					for done {
						fmt.Print((colorBlue), "Set the Reclaim Policy (Retain, Delete or Recycle): ")
						fmt.Scanln(&input)
						for _, s := range policy {
							if input == s {
								storage.Hostpath.ReclaimPolicy = input
								done = false
							}
						}
					}
					storage.Hostpath.Enabled = true
				case 4:
					fmt.Print((colorBlue), "Set the Node Selector")
					nodeselector := createSlice()
					storage.Hostpath.NodeSelector = nodeselector
					storage.Hostpath.Enabled = true
				}
				if intVar == 5 {
					fmt.Println("Saving and Exiting HostPath menu")
					break
				}
			}
		case 2:
			for {
				fmt.Println()
				fmt.Println((colorGreen), "----NFS Menu----")
				fmt.Println((colorGreen), "Update NFS values")
				fmt.Println((colorBlue), "Press '1' to modify Server IP address")
				fmt.Println((colorBlue), "Press '2' to modify NFS export path")
				fmt.Println((colorBlue), "Press '3' to set NFS as default Storage Class")
				fmt.Println((colorBlue), "Press '4' to modify Reclaim Policy [default: Retain]")
				fmt.Println((colorBlue), "Press '5' to Save and Exit")
				fmt.Print((colorWhite), "Please make your selection: ")
				caseInput := formatInput()
				intVar, _ := strconv.Atoi(caseInput)
				switch intVar {
				case 1:
					fmt.Print((colorWhite), "Input the NFS server IP address: ")
					ip := formatInput()
					storage.Nfs.Server = ip
					storage.Nfs.Enabled = true
				case 2:
					fmt.Print((colorWhite), "Input the NFS export path: ")
					path := formatInput()
					storage.Nfs.Path = path
					storage.Nfs.Enabled = true
				case 3:
					storage.Nfs.Enabled = true
					storage.Nfs.DefaultSc = true
					fmt.Println((colorYellow), "NFS set as default Storage Class")
				case 4:
					var input string
					var policy = []string{"Retain", "Delete", "Recycle"}
					done := true
					for done {
						fmt.Print((colorBlue), "Set the Reclaim Policy (Retain, Delete or Recycle): ")
						fmt.Scanln(&input)
						for _, s := range policy {
							if input == s {
								storage.Nfs.ReclaimPolicy = input
								done = false
							}
						}
					}
					storage.Nfs.Enabled = true
				}
				if intVar == 5 {
					fmt.Println((colorYellow), "Saving and Exiting NFS menu")
					break
				}
			}
		}
		if intVar == 3 {
			fmt.Println((colorYellow), "Saving and Exiting Storage menu")
			break
		}
	}
}

/* Function used to gather Single Sign On values through menu driven options.
This function uses the Sso struct
*/
func gatherSso(sso *Sso) {
	InfoLogger.Println("In the gatherSso function")

	for {
		fmt.Println()
		fmt.Println((colorGreen), "----Single Sign On Menu----")
		fmt.Println((colorGreen), "Update Single Sign On values")
		fmt.Println((colorBlue), "Press '1' to enable Single Sign On")
		fmt.Println((colorBlue), "Press '2' to modify Admin User")
		fmt.Println((colorBlue), "Press '3' to modify SSO Provider")
		fmt.Println((colorBlue), "Press '4' to modify Email Domain")
		fmt.Println((colorBlue), "Press '5' to modify Client ID")
		fmt.Println((colorBlue), "Press '6' to modify Client Secret")
		fmt.Println((colorBlue), "Press '7' to modify Azure Tenant")
		fmt.Println((colorBlue), "Press '8' to modify OIDC Issuer URL")
		fmt.Println((colorBlue), "Press '9' to Save and Exit Single Sign On menu")
		fmt.Print((colorWhite), "Please make your selection: ")
		caseInput := formatInput()
		intVar, _ := strconv.Atoi(caseInput)
		switch intVar {
		case 1:
			sso.Enabled = true
			fmt.Println((colorYellow), "Single Sign On Enabled")
			InfoLogger.Printf("Single Sign on Enable set to %v", sso.Enabled)
		case 2:
			fmt.Print((colorWhite), "Input the Admin User: ")
			admin := formatInput()
			sso.AdminUser = admin
			sso.Enabled = true
		case 3:
			fmt.Print((colorWhite), "Input the SSO Provider: ")
			provider := formatInput()
			sso.Provider = provider
			sso.Enabled = true
		case 4:
			fmt.Print((colorWhite), "Input the Email Domain: ")
			domain := createSlice()
			sso.EmailDomain = domain
			sso.Enabled = true
		case 5:
			fmt.Print((colorWhite), "Input the Client ID: ")
			clientid := formatInput()
			sso.ClientId = clientid
			sso.Enabled = true
		case 6:
			fmt.Print((colorWhite), "Input the Client Secret: ")
			var clientsecret string
			fmt.Scan(&clientsecret)
			sso.ClientSecret = clientsecret
			sso.Enabled = true
		case 7:
			fmt.Print((colorWhite), "Input the Azure Tenant: ")
			azure := formatInput()
			sso.AzureTenant = azure
			sso.Enabled = true
		case 8:
			fmt.Print((colorWhite), "Input the OIDC Issuer URL: ")
			oidc := formatInput()
			sso.OidcIssuerUrl = oidc
			sso.Enabled = true
		}
		if intVar == 9 {
			fmt.Println((colorYellow), "Saving and Exiting Single Sign On menu")
			break
		}
	}
}

// Function that will take a name and create a file
// in the root directory from Template
func createFile(name string, template *Template) {

	// Creating an empty file
	// Using Create() function
	myfile, e := os.Create(name)
	if e != nil {
		log.Fatal(e)
	}
	log.Println(myfile)

	// Execute the template and write to the file which was previously created
	f := temp.Execute(myfile, template)
	if f != nil {
		log.Print(f)
	}
	myfile.Close()
}

// valuesCmd represents the values command
var valuesCmd = &cobra.Command{
	Use:   "values",
	Short: "Command to generate a values file through user input",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// set variables for each struct defined above - Used in gather functions for each menu item
		// This also sets any defaults needed for templating
		internalDomain := ClusterInteralDomain{Domain: "cluster.local"}
		labels := Labels{}
		annotations := Annotations{}
		network := Networking{Istio: Istio{Enabled: true}, Ingress: Ingress{IstioGwEnabled: true}}
		logging := Logging{FluentbitEnable: true, ElastalertEnable: true, KibanaEnable: true}
		registry := Registry{}
		tenancy := Tenancy{}
		sso := Sso{}
		storage := Storage{Hostpath: Hostpath{Path: "/cnvrg-hostpath-storage"}}
		gpu := Gpu{NvidiaEnable: true, HabanaEnable: true}
		backup := Backup{Enabled: true}
		capsule := Capsule{Enabled: true}
		configreloader := ConfigReloader{Enabled: true}
		monitoring := Monitoring{DcgmExportEnable: true, HabanaExportEnable: true, NodeExportEnable: true, KubeStateMetricEnable: true,
			GrafanaEnable: true, PrometheusOperatorEnable: true, PrometheusEnable: true, DefaultSvcMonitorsEnable: true, CnvrgIdleMetricsEnable: true}
		controlplane := ControlPlane{HyperEnable: true, CnvrgScheduleEnable: true, SearchkiqEnable: true, SidekiqEnable: true, SystemkiqEnable: true,
			WebappEnable: true, MpiEnable: true, SearchkiqHpaEnable: true, SidekiqHpaEnable: true, SystemkiqHpaEnable: true, WebappHpaEnable: true}
		dbs := Dbs{EsEnable: true, MinioEnable: true, PgEnable: true, RedisEnable: true}

		//Start of program to ask user for Input
		InfoLogger.Println((colorWhite), "You are in the values main function")
		fmt.Println((colorGreen), "********************** Welcome **********************")
		fmt.Println((colorGreen), "We will gather your information to build a values file")
		clusterdomain := ClusterDomain{}
		gatherClusterDomain(&clusterdomain)
		for {
			fmt.Println()
			fmt.Println((colorGreen), "------------------------------- Main Menu -------------------------------")
			fmt.Println((colorGreen), "Please make a selection to modify the values file for the cnvrg.io install")
			fmt.Println((colorBlue), "Press '1' To modify Labeling---------------->[ Labels, Annotations or Internal Domain ]")
			fmt.Println((colorBlue), "Press '2' To modify Networks settings------->[ Istio, NodePort, HTTPS ]")
			fmt.Println((colorBlue), "Press '3' To modify Logging settings-------->[ Kibana, ElasticAlert, Fluentbit ]")
			fmt.Println((colorBlue), "Press '4' To modify Registry settings------->[ URL, Username, Password ]")
			fmt.Println((colorBlue), "Press '5' To modify Tenancy settings-------->[ Node Selector ]")
			fmt.Println((colorBlue), "Press '6' To modify Single Sign On settings->[ Admin, Provider, Azure Tenant ]")
			fmt.Println((colorBlue), "Press '7' To modify Storage settings-------->[ NFS, Hostpath ] ")
			fmt.Println((colorBlue), "Press '8' To modify Miscellaneous settings-->[ Backup, GPU, ConfigLoader, Capsule ]")
			fmt.Println((colorBlue), "Press '9' To modify Monitoring settings----->[ Prometheus, Grafana, Exporters ]")
			fmt.Println((colorBlue), "Press '10' To modify Control Plane settings->[ CP Image, CP Services, SMTP ]")
			fmt.Println((colorBlue), "Press '11' To modify Database settings------>[ Minio, Postgres, Redis ]")
			fmt.Println((colorBlue), "Press '12' To Exit and generate Values file")
			fmt.Print((colorWhite), "Please make your selection: ")
			caseInput := formatInput()
			intVar, _ := strconv.Atoi(caseInput)
			switch intVar {
			case 1:
				for {
					fmt.Println()
					fmt.Println((colorGreen), "----Labels, Annotations Internal Domain Menu----")
					fmt.Println((colorGreen), "Update Labels, Annotations or Internal Domain values")
					fmt.Println((colorBlue), "Press '1' To modify Labels")
					fmt.Println((colorBlue), "Press '2' To modify Annotations")
					fmt.Println((colorBlue), "Press '3' To modify Internal Domain")
					fmt.Println((colorBlue), "Press '4' To Save and Exit")
					fmt.Print((colorWhite), "Please make your selection: ")
					caseInput := formatInput()
					intVar, _ := strconv.Atoi(caseInput)
					switch intVar {
					case 1:
						gatherLabels(&labels)
					case 2:
						gatherAnnotations(&annotations)
					case 3:
						gatherInternalDomain(&internalDomain)
					}
					if intVar == 4 {
						fmt.Println((colorYellow), "Saving and Exiting menu")
						break
					}
				}
			case 2:
				gatherNetworking(&network)
			case 3:
				gatherLogging(&logging)
			case 4:
				gatherRegistry(&registry)
			case 5:
				gatherTenancy(&tenancy)
			case 6:
				gatherSso(&sso)
			case 7:
				gatherStorage(&storage)
			case 8:
				for {
					fmt.Println()
					fmt.Println((colorGreen), "----Backup, GPU, Capsule and GPU Menu----")
					fmt.Println((colorGreen), "Update Backup, GPU, Capsule and GPU values")
					fmt.Println((colorBlue), "Press '1' To modify Backup settings")
					fmt.Println((colorBlue), "Press '2' To modify Capsule settings")
					fmt.Println((colorBlue), "Press '3' To disable NvidiaDp or HabanaDp GPU")
					fmt.Println((colorBlue), "Press '4' To disable ConfigReloader")
					fmt.Println((colorBlue), "Press '5' To Exit modifying settings")
					fmt.Print((colorBlue), "Please make your selection: ")
					caseInput := formatInput()
					intVar, _ := strconv.Atoi(caseInput)
					switch intVar {
					case 1:
						gatherBackup(&backup)
					case 2:
						gatherCapsule(&capsule)
					case 3:
						gatherGpu(&gpu)
					case 4:
						gatherConfigReloader(&configreloader)
					}
					if intVar == 5 {
						fmt.Println((colorYellow), "Saving changes and exiting")
						break
					}
				}
			case 9:
				gatherMonitoring(&monitoring)
			case 10:
				gatherControlPlane(&controlplane)
			case 11:
				gatherDbs(&dbs)
			}
			if intVar == 12 {
				fmt.Println((colorWhite), "Exiting and generating the values.yaml file")
				break
			}
		}

		finaltemp := Template{clusterdomain, internalDomain, labels, annotations, network, logging, registry, tenancy,
			sso, storage, configreloader, capsule, backup, gpu, monitoring, controlplane, dbs}
		err := temp.Execute(os.Stdout, finaltemp)
		if err != nil {
			log.Print(err)
		}
		createFile("values.yaml", &finaltemp)
		outputHelm()
	},
}

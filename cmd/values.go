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

// Parent level of Registry struct
type RegCred struct {
	Username string
	Password string
	Url      string `default:"docker.io"`
	RegCred  bool
}

// Parent level of SSO struct
type Sso struct {
	Sso SsoValues
}

// Used in the Sso parent struct
type SsoValues struct {
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

type Nfs struct {
	Enabled       bool
	Server        string
	Path          string
	DefaultSc     bool
	ReclaimPolicy string
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
	Type           string `default: "istio"`
	IstioGwEnabled bool   `default: true, yaml:"IstioGwEnabled"`
	IstioGwName    string `yaml: "IstioGwName"`
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
	Registry RegCred
	Network  Networking
	Sso      Sso
	Storage  Storage
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
		sso.Sso.Enabled = false
	}
	if enableSso == "yes" {
		sso.Sso.Enabled = true
		sso.Sso.AdminUser = ""
		sso.Sso.Provider = ""
		sso.Sso.EmailDomain = []string{"10.2.3.8,", "192.168.1.5"}
		sso.Sso.ClientId = ""
		sso.Sso.ClientSecret = ""
		sso.Sso.AzureTenant = ""
		sso.Sso.OidcIssuerUrl = ""
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

	// testing example remove later
	pleasework := []string{"hello,", "hows it going"}

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
		network.Proxy.HttpProxy = pleasework
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
		registry := RegCred{}
		fmt.Println("enable the registry?")
		registry.RegCred = true
		registry.Username = "brad"
		registry.Password = "Password123!"
		network := Networking{}
		gatherNetworking(&network)
		sso := Sso{}
		gatherSso(&sso)
		storage := Storage{}
		gatherStorage(&storage)
		finaltemp := Template{registry, network, sso, storage}
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

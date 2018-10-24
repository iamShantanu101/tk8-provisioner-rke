package cluster

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/kubernauts/tk8-provisioner-rke/internal/templates"
	"github.com/spf13/viper"
)

type RKEConfig struct {
	ClusterName         string
	AWSRegion           string
	RKENodeInstanceType string
	NodeCount           int
	SSHKeyPath          string
	CloudProvider       string
}

func GetRKEConfig() RKEConfig {
	ReadViperConfigFile("config")
	return RKEConfig{
		ClusterName:         viper.GetString("rke.cluster_name"),
		AWSRegion:           viper.GetString("rke.rke_aws_region"),
		RKENodeInstanceType: viper.GetString("rke.rke_node_instance_type"),
		NodeCount:           viper.GetInt("rke.node_count"),
		SSHKeyPath:          viper.GetString("rke.ssh_key_path"),
		CloudProvider:       viper.GetString("rke.cloud_provider"),
	}
}

// Install is used to setup the Kubernetes Cluster with RKE
func Install() {
	var Name string
	config := GetRKEConfig()
	Name = config.ClusterName
	os.MkdirAll("./inventory/"+Name+"/provisioner", 0755)
	exec.Command("cp", "-rfp", "./provisioner/rke/", "./inventory/"+Name+"/provisioner").Run()
	ParseTemplate(templates.VariablesRKE, "./inventory/"+Name+"/provisioner/variables.tfvars", GetRKEConfig())
	ParseTemplate(templates.Credentials, "./inventory/"+Name+"/provisioner/credentials.tfvars", GetCredentials())
	// Check if a terraform state file aclready exists
	if _, err := os.Stat("./inventory/" + Name + "/provisioner/terraform.tfstate"); err == nil {
		log.Println("There is an existing cluster, please remove terraform.tfstate file or delete the installation before proceeding")
	} else {
		log.Println("starting terraform init")

		terrInit := exec.Command("terraform", "init")
		terrInit.Dir = "./inventory/" + Name + "/provisioner/"
		out, _ := terrInit.StdoutPipe()
		terrInit.Start()
		scanInit := bufio.NewScanner(out)
		for scanInit.Scan() {
			m := scanInit.Text()
			fmt.Println(m)
		}

		terrInit.Wait()
	}

	log.Println("starting terraform apply")
	terrSet := exec.Command("terraform", "apply", "-var-file=credentials.tfvars", "-auto-approve")
	terrSet.Dir = "./inventory/" + Name + "/provisioner/"
	stdout, err := terrSet.StdoutPipe()
	terrSet.Stderr = terrSet.Stdout
	terrSet.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}

	terrSet.Wait()
	if err != nil {
		panic(err)
	}
	// // Export KUBECONFIG file to the installation folder
	log.Println("Moving kubeconfig and rke cluster config files to the installation folder")
	errmvKubeconfig := os.Rename("./inventory"+Name+"/provisioner/kube_config_cluster.yml", "./kube_config_cluster.yml")

	if errmvKubeconfig != nil {
		fmt.Println(err)
	}

	errmvRkeConfig := os.Rename("./inventory"+Name+"/provisioner/rancher-cluster.yml", "./rancher-cluster.yml")

	if errmvRkeConfig != nil {
		fmt.Println(err)
	}

	// log.Println("Voila! Kubernetes cluster created with RKE is up and running")

	os.Exit(0)

}

// Reset is used to reset the  Kubernetes Cluster back to rollout on the infrastructure.
func RKEReset() {
	NotImplemented()
}

// Remove is used to remove the Kubernetes Cluster from the infrastructure
func RKERemove() {
	NotImplemented()
	log.Println("Removing rke cluster")
	rkeRemove := exec.Command("rke", "remove", "--config rancher-cluster.yml")
	stdout, err := rkeRemoveSet.StdoutPipe()
	rkeRemoveSet.Stderr = rkeRemoveSet.Stdout
	rkeRemoveSet.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}

	terrSet.Wait()
	if err != nil {
		panic(err)
	}
}

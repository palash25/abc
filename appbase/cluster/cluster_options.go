package cluster

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/AlecAivazis/survey.v1"
)

var clusterOptions = []*survey.Question{
	{
		Name:     "name",
		Prompt:   &survey.Input{Message: "Enter the name of the cluster"},
		Validate: survey.Required,
	},
	{
		Name: "location",
		Prompt: &survey.Select{
			Message: "Enter the location of the cluster",
			Options: []string{
				"eastus",
				"westeurope",
				"centralus",
				"canadacentral",
				"canadaeast",
				"australiaeast",
				"eastus2",
				"japaneast",
				"northeurope",
				"southeastasia",
				"uksouth",
				"westus2",
				"westus",
				"us-east1-b",
				"europe-west1-b",
				"us-central1-b",
				"australia-southeast1-b",
				"us-east4-b",
				"southamerica-east1-c",
				"northamerica-northeast1-b",
				"europe-north1-b",
				"asia-southeast1-b",
				"asia-east1-b",
				"asia-northeast1-a",
			},
		},
		Validate: survey.Required,
	},
	{
		Name:     "vm_size",
		Prompt:   &survey.Input{Message: "Enter the VM size of the cluster"},
		Validate: survey.Required,
	},
	{
		Name:     "ssh_key",
		Prompt:   &survey.Input{Message: "Enter the ssh public key of the cluster"},
		Validate: survey.Required,
	},
	{
		Name: "plan",
		Prompt: &survey.Select{
			Message: "Enter the pricing plan",
			Options: []string{"sandbox", "hobby", "production1", "production2", "production3"},
		},
		Validate: survey.Required,
	},
	{
		Name: "provider",
		Prompt: &survey.Select{
			Message: "Choose the provider:",
			Options: []string{"azure", "gke"},
		},
		Validate: survey.Required,
	},
}

func buildClusterObjectString() string {
	fmt.Println("Enter the cluster details")

	answers := make(map[string]interface{})

	err := survey.Ask(clusterOptions, &answers)
	if err != nil {
		fmt.Println(err.Error())
	}

	clusterObject := "\"cluster\": {\n    "
	for key, value := range answers {
		clusterObject = clusterObject + "\"" + key + "\": " + value.(string) + ",\n    "
	}
	clusterObject = clusterObject + "},\n"
	return clusterObject

}

// TODO add validation checks
var esOptions = []*survey.Question{
	{
		Name:     "nodes",
		Prompt:   &survey.Input{Message: "Enter the number of ES nodes"},
		Validate: survey.Required,
	},
	{
		Name:     "version",
		Prompt:   &survey.Input{Message: "Enter ES version"},
		Validate: survey.Required,
	},
	{
		Name:     "volume_size",
		Prompt:   &survey.Input{Message: "Enter the volume size"},
		Validate: survey.Required,
	},
	{
		Name:   "config_url",
		Prompt: &survey.Input{Message: "Enter the config url"},
	},
	{
		Name:   "heap_size",
		Prompt: &survey.Input{Message: "Enter the heap size"},
	},
	{
		Name: "plugins",
		Prompt: &survey.MultiSelect{
			Message: "Enter plugins",
			Options: []string{
				"ICU Analysis",
				"Japanese Analysis",
				"Phonetic Analysis",
				"Smart Chinese Analysis",
				"Ukrainian Analysis",
				"Stempel Polish Analysis",
				"Ingest Attachment Processor",
				"Ingest User Agent Processor",
				"Mapper Size",
				"Mapper Murmur3",
				"X-Pack",
			},
		},
	},
	{
		Name:   "backup",
		Prompt: &survey.Confirm{Message: "Do you want backup?"},
	},
	{
		Name:   "env",
		Prompt: &survey.Input{Message: "Enter env"},
	},
}

// GetESDetails ....
func buildESObjectString() string {
	fmt.Println("Enter the cluster details")
	answers := make(map[string]interface{})

	err := survey.Ask(esOptions, &answers)
	if err != nil {
		fmt.Println(err.Error())
	}

	esObject := "\"elasticsearch\": {\n    "

	// TODO add case for env and checks for empty keys
Next:
	for key, value := range answers {
		if key == "nodes" {
		} else if key == "backup" {
			value = strconv.FormatBool(value.(bool))
		} else if key == "plugins" {
			tempStr := ""
			for _, str := range value.([]string) {
				tempStr = tempStr + "\"" + str + "\", "
			}
			value = "[" + strings.Trim(tempStr, ", ") + "]"
		} else {
			esObject = esObject + "\"" + key + "\": " + "\"" + value.(string) + "\"," + "\n    "
			continue Next
		}

		esObject = esObject + "\"" + key + "\": " + value.(string) + ",\n    "
	}

	esObject = esObject + "},\n"
	return esObject
}

var logstashKibanaOptions = []*survey.Question{
	{
		Name:   "create_node",
		Prompt: &survey.Confirm{Message: "Do you want to create node?"},
	},
	{
		Name:     "version",
		Prompt:   &survey.Input{Message: "Enter a valid ES version"},
		Validate: survey.Required,
	},
	{
		Name:   "heap_size",
		Prompt: &survey.Input{Message: "Enter the heap size"},
	},
	{
		Name:   "env",
		Prompt: &survey.Input{Message: "Enter env"},
	},
}

func buildLogstashObjectString() string {
	fmt.Println("Enter the cluster details")
	answers := struct {
		CreateNode bool              `survey:"create_node"`
		Version    string            `survey:"version"`
		HeapSize   string            `survey:"heap_size"`
		Env        map[string]string `survey:"env"`
	}{}

	err := survey.Ask(logstashKibanaOptions, &answers)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(answers.CreateNode)

	return ""
}

func buildKibanaObjectString() string {
	fmt.Println("Enter the cluster details")
	answers := struct {
		CreateNode bool              `survey:"create_node"`
		Version    string            `survey:"version"`
		HeapSize   string            `survey:"heap_size"`
		Env        map[string]string `survey:"env"`
	}{}

	err := survey.Ask(logstashKibanaOptions, &answers)
	if err != nil {
		fmt.Println(err.Error())
	}

	return ""
}

var addonsOptions = []*survey.Question{
	{
		Name: "name",
		Prompt: &survey.Select{
			Message: "Choose an addon from the following list:",
			Options: []string{"dejavu", "elasticsearch-hq", "mirage"},
		},
		Validate: survey.Required,
	},
	{
		Name:     "image",
		Prompt:   &survey.Input{Message: "Enter image"},
		Validate: survey.Required,
	},
	{
		Name:     "exposed_port",
		Prompt:   &survey.Input{Message: "Enter the exposed port"},
		Validate: survey.Required,
	},
	{
		Name:   "env",
		Prompt: &survey.Input{Message: "Enter env"},
	},
	{
		Name:   "path",
		Prompt: &survey.Input{Message: "Enter path"},
	},
}

func buildAddonsObjectString() string {

	fmt.Print("Enter the cluster details")
	answers := make(map[string]interface{})

	err := survey.Ask(addonsOptions, &answers)
	if err != nil {
		fmt.Println(err.Error())
	}

	return ""

}

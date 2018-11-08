package cluster

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/AlecAivazis/survey"
	"github.com/appbaseio/abc/appbase/common"
	"github.com/appbaseio/abc/appbase/session"
	"github.com/appbaseio/abc/appbase/spinner"
)

type status struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type cluster struct {
	Name      string    `json:"name"`
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Message   string    `json:"message"`
	Provider  string    `json:"provider"`
}

type createClusterRespBody struct {
	Status  status  `json:"status"`
	Cluster cluster `json:"cluster"`
}

// RunClusterCreate runs cluster create command
func RunClusterCreate(clusterName string) error {

	body := strings.NewReader("{\n  \"elasticsearch\": {\n    \"nodes\": 1,\n    \"version\": \"5.6.10\",\n    \"volume_size\": 15,\n    \"env\": {\n      \"cluster.name\": \"beta-73\"\n    }\n  },\n  \"cluster\": {\n    \"name\": \"beta-73\",\n    \"vm_size\": \"Standard_B2s\",\n    \"location\": \"centralus\",\n    \"pricing_plan\": \"Growth\",\n    \"ssh_public_key\": \"ssh-rsa y68qVAM1cEgl4qlotnpS7LuxtmTD/6HR9WjKioJNEUcL2RZzkPnzM3MNoybFbv6Gu5cnZJMTOPfKvHjM/K9s9V6oRuT56HRDLbSb7s9v91Qi3DT4jMoCrC0Y9nDecBs0mBu+ijY0Kjdkjddfke5O20B8u7gE7DNfD+w0PZ3uNYdV5l3EuqVBeRATjAAj2DD7TomIfXbPVGldg0ou2jdU5VmJe75Zmb1CgnGJDkp8RGgXQkVj2DVu8Ia3 alpha.gama@gmail.com\"\n  }\n}")

	spinner.Start()
	defer spinner.Stop()

	req, err := http.NewRequest("POST", common.AccAPIURL+"/v1/_deploy", body)
	if err != nil {
		return err
	}
	resp, err := session.SendRequest(req)
	if err != nil {
		return err
	}
	spinner.Stop()
	// status code not 200
	if resp.StatusCode != 202 {
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("There was an error %s", string(bodyBytes))
	}

	return nil
}

var additionalChoices = []*survey.Question{
	{
		Name:     "logstash",
		Prompt:   &survey.Confirm{Message: "Would you like to provide Logstash options to your cluster deployment?"},
		Validate: survey.Required,
	},
	{
		Name:     "kibana",
		Prompt:   &survey.Confirm{Message: "Would you like to provide Kibana options to your cluster deployment?"},
		Validate: survey.Required,
	},
	{
		Name:     "addons",
		Prompt:   &survey.Confirm{Message: "Would you like to add add-ons to your cluster deployment?"},
		Validate: survey.Required,
	},
}

func buildRequestBody() string {
	answers := make(map[string]interface{})

	err := survey.Ask(additionalChoices, &answers)
	if err != nil {
		fmt.Println(err.Error())
	}

	respBodyString := "{\n    " + buildESObjectString() + buildClusterObjectString()

	/*
		if answers["logstash"] == true {
			respBodyString = respBodyString + buildLogstashObjectString()
		}
		if answers["kibana"] == true {
			respBodyString = respBodyString + buildKibanaObjectString()
		}
		if answers["addons"] == true {
			name := ""
			prompt := &survey.Input{
				Message: "How many addons do you want to add (Choose between 1-3)?",
			}
			survey.AskOne(prompt, &name, nil)

		}*/
	return respBodyString + "}"
}

func RunClusterCreateInteractive() error {
	//str := buildRequestBody()
	//payload := strings.NewReader(str)
	//fmt.Println(str)
	payl := strings.NewReader("{\n  \"elasticsearch\": {\n    \"nodes\": 1,\n    \"version\": \"5.6.10\",\n    \"volume_size\": 15,\n    \"env\": {\n      \"cluster.name\": \"gke-beta-169\"\n    }\n  },\n  \"cluster\": {\n    \"name\": \"gke\",\n    \"vm_size\": \"n1-standard-1\",\n    \"location\": \"australia-southeast1-b\",\n    \"pricing_plan\": \"Growth\",\n    \"provider\": \"gke\",\n    \"ssh_public_key\": \"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC/eM7PAC5A2IOeAm3GXkHX6fHBoO8eVt3KqsO8g6LlVTDekqSffJyyN5PtdE+I3a2PLhBog1ghBoIQ1JV9+uuxjtX+DgLqy68qVAM1cEgl4qlotnpS7LuxtmTD/6HR9WjKioJNEUcL2RZzkPnzM3MNoybFbv6Gu5cnZJMTOPfKvHjM/K9s9V6oRuT56HRDLbSb7s9v91Qi3DT4jMoCrC0Y9nDecBs0mBu+ijY0ADrgxMCCgW5O20B8u7gE7DNfD+w0PZ3uNYdV5l3EuqVBeRATjAAj2DD7TomIfXbPVGldg0ou2jdU5VmJe75Zmb1CgnGJDkp8RGgXQkVj2DVu8Ia3 raaz.crzy@gmail.com\"\n  }\n}")

	spinner.Start()
	defer spinner.Stop()

	req, err := http.NewRequest("POST", common.AccAPIURL+"/v1/_deploy", payl)
	if err != nil {
		fmt.Println("Request error block")
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := session.SendRequest(req)
	if err != nil {
		fmt.Println("send error", err)
		return err
	}
	spinner.Stop()

	// status code not 200
	if resp.StatusCode != 202 {
		fmt.Println("ENTER S CODE", resp.StatusCode)
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("There was an error %s", string(bodyBytes))
	}

	fmt.Println(resp.StatusCode)

	var res createClusterRespBody
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&res)
	if err != nil {
		fmt.Println("Decode err")
		return err
	}

	// output
	fmt.Printf("ID:    %s\n", res.Cluster.ID)
	fmt.Printf("Name:  %s\n", res.Cluster.Name)
	fmt.Printf("Status:  %s\n", res.Cluster.Status)
	fmt.Printf("Provider:  %s\n", res.Cluster.Provider)
	fmt.Printf("Created at:  %s\n", res.Cluster.CreatedAt)
	fmt.Printf("Message:  %s\n", res.Cluster.Message)

	return nil
}

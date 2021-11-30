package pkg

import (
	"bytes"
	"fmt"
	"github.com/qovery/qovery-cli/utils"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func DeployById(clusterId string, dryRunDisabled bool){
	utils.CheckAdminUrl()

	utils.DryRunPrint(dryRunDisabled)
	if utils.Validate("deployment") {
		res := deploy(os.Getenv("ADMIN_URL") + "/cluster/deploy/" + clusterId, http.MethodPost, dryRunDisabled )

		if !strings.Contains(res.Status, "200") {
			result, _ := ioutil.ReadAll(res.Body)
			log.Errorf("Could not deploy cluster : %s. %s", res.Status, string(result) )
		} else {
			fmt.Println("Cluster " + clusterId + " deploying.")
		}
	}
}

func DeployAll(dryRunDisabled bool) {
	utils.CheckAdminUrl()

	utils.DryRunPrint(dryRunDisabled)
	if utils.Validate("deployment") {
		res := deploy(os.Getenv("ADMIN_URL") + "/cluster/deploy", http.MethodPost, dryRunDisabled )

		if !strings.Contains(res.Status, "200") {
			result, _ := ioutil.ReadAll(res.Body)
			log.Errorf("Could not deploy clusters : %s. %s", res.Status, string(result) )
		} else {
			fmt.Println("Clusters deploying.")
		}
	}
}

func deploy(url string, method string, dryRunDisabled bool) *http.Response {
	authToken, tokenErr := utils.GetAccessToken()
	if tokenErr != nil {
		utils.PrintlnError(tokenErr)
		os.Exit(0)
	}

	var body *bytes.Buffer

	if dryRunDisabled {
		body = bytes.NewBuffer([]byte( `{ "metadata": { "dry_run_deploy": true } }`))
	}

	req, err  := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer " + strings.TrimSpace(string(authToken)))
	if dryRunDisabled {
		req.Header.Set("Content-Type", "application/json")
	}


	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return res
}


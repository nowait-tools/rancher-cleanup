package main

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/client"
)

var (
	debug           = flag.Bool("debug", false, "Debug")
	apiEndpoint     = ""
	apiAccessKey    = ""
	apiSecretKey    = ""
	apiTimeout      = 10
	removalInterval = 30

	log *logrus.Logger
	api *client.RancherClient
)

func init() {
	flag.Parse()

	log = logrus.New()

	if *debug {
		log.Level = logrus.DebugLevel
	}

	setEnvVariables()
	initAPIClient()
}

func setEnvVariables() {
	if len(os.Getenv("RANCHER_API_ENDPOINT")) > 0 {
		apiEndpoint = os.Getenv("RANCHER_API_ENDPOINT")
	}

	if len(os.Getenv("RANCHER_API_ACCESSKEY")) > 0 {
		apiAccessKey = os.Getenv("RANCHER_API_ACCESSKEY")
	}

	if len(os.Getenv("RANCHER_API_SECRETKEY")) > 0 {
		apiSecretKey = os.Getenv("RANCHER_API_SECRETKEY")
	}

	if len(os.Getenv("RANCHER_API_TIMEOUT")) > 0 {
		t := os.Getenv("RANCHER_API_TIMEOUT")
		timeout, err := strconv.Atoi(t)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("Unable to parse api timeout value")
		}
		apiTimeout = timeout
	}

	if len(os.Getenv("NOWAIT_HOST_REMOVAL_INTERVAL")) > 0 {
		i := os.Getenv("NOWAIT_HOST_REMOVAL_INTERVAL")
		interval, err := strconv.Atoi(i)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("Unable to parse api timeout value")
		}
		removalInterval = interval
	}

	log.WithFields(logrus.Fields{
		"RANCHER_API_ENDPOINT":  apiEndpoint,
		"RANCHER_API_ACCESSKEY": apiAccessKey,
		"RANCHER_API_SECRETKEY": apiSecretKey,
		"RANCHER_API_TIMEOUT":   apiTimeout,
	}).Debug("Environment variables set")
}

func initAPIClient() {

	apiClient, err := client.NewRancherClient(&client.ClientOpts{
		Url:       apiEndpoint,
		AccessKey: apiAccessKey,
		SecretKey: apiSecretKey,
		Timeout:   time.Duration(apiTimeout) * time.Second,
	})

	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Unable to initialize rancher api client")
	}

	api = apiClient
}

func main() {

	cleanupReconnectingHosts()
	tickChan := time.NewTicker(time.Duration(removalInterval) * time.Second).C

	for {
		select {
		case <-tickChan:
			cleanupReconnectingHosts()
		}
	}
}

func cleanupReconnectingHosts() {
	deactivateHosts()
	removeHosts()
	purgeHosts()
}

func deactivateHosts() {
	hosts, err := api.Host.List(&client.ListOpts{Filters: map[string]interface{}{
		"state":      "active",
		"agentState": "reconnecting",
	}})

	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Unable to retrieve list of reconnecting hosts from api")
	}

	for _, host := range hosts.Data {
		_, err := api.Host.ActionDeactivate(&host)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
				"host":  host,
			}).Warnf("Unable to deactivate reconnecting host: %v", host.Resource.Id)
		}
		log.Infof("Deactivated host: %s", host.Resource.Id)
	}

	log.Infof("Deactivated %v hosts", len(hosts.Data))
}

func removeHosts() {

	hosts, err := api.Host.List(&client.ListOpts{Filters: map[string]interface{}{
		"state":      "inactive",
		"agentState": "reconnecting",
	}})

	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Unable to retrieve list of inactive hosts from api")
	}

	for _, host := range hosts.Data {
		_, err := api.Host.ActionRemove(&host)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
				"host":  host,
			}).Warnf("Unable to remove inactive host: %v", host.Resource.Id)
		}

		log.Infof("Removed host: %s", host.Resource.Id)
	}
}

func purgeHosts() {

	hosts, err := api.Host.List(&client.ListOpts{Filters: map[string]interface{}{
		"state":      "removed",
		"agentState": "reconnecting",
	}})

	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Unable to retrieve list of removed hosts from api")
	}

	for _, host := range hosts.Data {
		_, err := api.Host.ActionPurge(&host)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
				"host":  host,
			}).Warnf("Unable to purge host: %v", host.Resource.Id)
		}

		log.Infof("Purged host: %s", host.Resource.Id)
	}
}

package scan

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/bass3m/swarmeus/config"
	"github.com/fsouza/go-dockerclient"
	"io/ioutil"
	"regexp"
	"strconv"
	"time"
)

var client *docker.Client

func Initialize(c *docker.Client) {
	client = c
}

type prometheusTarget struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

func targetsToJson(targets []prometheusTarget) []byte {
	jsonTargets, err := json.MarshalIndent(targets, "", "  ")
	if err != nil {
		panic(err)
	}
	log.Debugf("JSON Targets: %+v", string(jsonTargets))
	return jsonTargets
}

func writeSDFile(targets []prometheusTarget, filePath string) {
	log.Debugf("Writing Targets: %+v to %+v", targets, filePath)
	jsonTargets := targetsToJson(targets)

	err := ioutil.WriteFile(filePath, jsonTargets, 0644)
	if err != nil {
		panic(err)
	}
}

func findTargets(mode string, network string, targets []config.Target) ([]prometheusTarget, error) {
	var pts []prometheusTarget
	cs, err := client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		log.Warnf("No containers found. error: %v", err)
		return pts, err
	}
	for _, c := range cs {
		for key, value := range c.Labels {
			if key == mode {
				for _, t := range targets {
					reg, err := regexp.Compile(t.InstanceRegex)
					if err == nil {
						jobName := reg.FindString(value)
						if jobName != "" {
							// found a match
							log.Debugf("Found matching docker container %+v", c)
							ip := c.Networks.Networks[network].IPAddress
							if ip == "" {
								log.Warnf("IP addr not set yet for container ID %+v", c.ID)
								continue
							}
							uri := ip + ":" + strconv.FormatInt(t.Port, 10)
							pts = append(pts,
								prometheusTarget{Labels: map[string]string{"job": jobName},
									Targets: []string{uri}})
						}
					}
				}
			}
		}
	}
	return pts, nil
}

func Scan(cfg config.Config, cancel <-chan struct{}) {
	log.Infoln("Starting prometheus target scan")
Loop:
	for {
		select {
		case <-cancel:
			log.Infoln("Received cancel event. Cancelling scan")
			return
		case <-time.After(time.Second * cfg.Swarmeus.ScanInterval):
			targets, err := findTargets(cfg.Swarmeus.DockerMode, cfg.Swarmeus.Network, cfg.Targets)
			if err != nil {
				log.Errorf("Failed to find targets, retrying. error: %v", err)
			}
			log.Debugf("Scanned Prometheus Targets: %+v", targets)
			// write sd file
			writeSDFile(targets, cfg.Swarmeus.SDFilePath)
			log.Debugf("Write SD file to: %+v", cfg.Swarmeus.SDFilePath)
			continue Loop
		}
	}
}

func GetTargets(cfg config.Config) ([]byte, error) {
	targets, err := findTargets(cfg.Swarmeus.DockerMode, cfg.Swarmeus.Network, cfg.Targets)
	if err != nil {
		log.Errorf("Failed to find targets, retrying. error: %v", err)
		return []byte{}, err
	}
	jsonTargets := targetsToJson(targets)
	return jsonTargets, nil
}

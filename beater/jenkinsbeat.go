package beater

import (
	"fmt"
	"time"
	"os"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/bndr/gojenkins"

	"github.com/tzankov/jenkinsbeat/config"
)

type Jenkinsbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}


// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Jenkinsbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

func (bt *Jenkinsbeat) Run(b *beat.Beat) error {
	logp.Info("jenkinsbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		logp.Info("Getting All Jobs")
		bt.getAllJobs()

		logp.Info("Getting Current Build Queue")
		bt.getQueueJobs()

		logp.Info("Getting Latest Builds")
		bt.getLatestBuildOfJobs()
		
		logp.Info("Events sent")
		}
	}
}

func (bt *Jenkinsbeat) getAllJobs() {
	jenkins := gojenkins.CreateJenkins(nil, os.Getenv("JENKINS_URL"), os.Getenv("JENKINS_USER"), os.Getenv("JENKINS_PASS"))

	jobs, err := jenkins.GetAllJobs()
		if err != nil {
			logp.Err("Failed to collect jobs list, got :", err)
			return
		}

		for _, element := range jobs {
			bt.client.Publish(bt.newJobsEvent(element))
		}

}

func (bt *Jenkinsbeat) getQueueJobs() {
	jenkins := gojenkins.CreateJenkins(nil, os.Getenv("JENKINS_URL"), os.Getenv("JENKINS_USER"), os.Getenv("JENKINS_PASS"))

	buildQueue, err := jenkins.GetQueue()
		if err != nil {
			logp.Err("Failed to collect queue list, got :", err)
			return
		}

	bt.client.Publish(bt.newQueueEvent(buildQueue))

}

func (bt *Jenkinsbeat) getLatestBuildOfJobs() {
	jenkins := gojenkins.CreateJenkins(nil, os.Getenv("JENKINS_URL"), os.Getenv("JENKINS_USER"), os.Getenv("JENKINS_PASS"))

	jobs, err := jenkins.GetAllJobs()
		if err != nil {
			logp.Err("Failed to collect latest list of build jobs, got :", err)
			return
		}

		for _, element := range jobs {
			bt.client.Publish(bt.newLatestBuildEvent(element))
		}

}

func (Jenkinsbeat) newJobsEvent(element *gojenkins.Job) beat.Event {

		event := beat.Event{
		Timestamp: time.Now(),
		Fields: common.MapStr{
			"type": "jenkinsbeat",
			"job":  element.Base,
			},
		}

	return event

}

func (Jenkinsbeat) newLatestBuildEvent(build *gojenkins.Job) beat.Event {
	event := beat.Event{
	Timestamp: time.Now(),
	Fields: common.MapStr{
		"type":   "jenkinsbeat",
		"latestBuildEvent":  build.Base,
		},
	}

	return event

}

func (Jenkinsbeat) newQueueEvent(queue *gojenkins.Queue) beat.Event {
	event := beat.Event{
	Timestamp: time.Now(),
	Fields: common.MapStr{
		"type":   "jenkinsbeat",
		"queueEvent":  queue.Base,
		"jenkins": queue.Jenkins,
		},
	}

	return event

}

func (bt *Jenkinsbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

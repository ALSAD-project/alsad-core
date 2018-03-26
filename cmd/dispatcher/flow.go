package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/communicator"
	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/config"
	"github.com/ALSAD-project/alsad-core/pkg/dispatcher/flow"
)

func makeFlow(config config.Config) (flow.Flow, error) {
	flowName := strings.ToLower(strings.TrimSpace(config.EnterMode))
	switch flowName {
	case "training":
		return makeTrainingFlow(config)
	case "detect":
		return makeDetectFlow(config)
	case "feedback":
		return makeFeedbackFlow(config)
	default:
		return nil, fmt.Errorf("unknown flow name: %s", config.EnterMode)
	}
}

func makeTrainingFlow(config config.Config) (flow.Flow, error) {
	log.Print("Starting a training flow")
	feederComm, err := communicator.NewHTTPCommunicator(config.FeederUrl)
	if err != nil {
		return nil, err
	}

	baComm, err := communicator.NewHTTPCommunicator(config.BaUrl)
	if err != nil {
		return nil, err
	}

	uslComm, err := communicator.NewHTTPCommunicator(config.UslUrl)
	if err != nil {
		return nil, err
	}

	expertSender, err := communicator.NewFSRedisQueueCommunicator(
		config.FsqDir,
		config.FsqExpertInputQueue,
		config.FsqRedisAddr,
	)
	if err != nil {
		return nil, err
	}
	expertFetcher, err := communicator.NewFSRedisQueueCommunicator(
		config.FsqDir,
		config.FsqExpertOutputQueue,
		config.FsqRedisAddr,
	)
	if err != nil {
		return nil, err
	}

	expertComm, err := communicator.NewCombineCommunicator(expertFetcher, expertSender)
	if err != nil {
		return nil, err
	}

	slComm, err := communicator.NewHTTPCommunicator(config.SlUrl)
	if err != nil {
		return nil, err
	}

	return newTrainingFlowFromCommunicators(
		feederComm,
		baComm,
		uslComm,
		expertComm,
		slComm,
		config.BasicRateLimit,
	)
}

func makeDetectFlow(config config.Config) (flow.Flow, error) {
	log.Print("Starting a detect flow")
	feederComm, err := communicator.NewHTTPCommunicator(config.FeederUrl)
	if err != nil {
		return nil, err
	}

	baComm, err := communicator.NewHTTPCommunicator(config.BaUrl)
	if err != nil {
		return nil, err
	}

	uslComm, err := communicator.NewHTTPCommunicator(config.UslUrl)
	if err != nil {
		return nil, err
	}

	slComm, err := communicator.NewHTTPCommunicator(config.SlUrl)
	if err != nil {
		return nil, err
	}

	expertComm, err := communicator.NewFSRedisQueueCommunicator(
		config.FsqDir,
		config.FsqExpertInputQueue,
		config.FsqRedisAddr,
	)
	if err != nil {
		return nil, err
	}

	return newDetectFlowFromCommunicators(
		feederComm,
		baComm,
		uslComm,
		slComm,
		expertComm,
		config.BasicRateLimit,
	)
}

func makeFeedbackFlow(config config.Config) (flow.Flow, error) {
	log.Print("Starting a feedback flow")
	expertComm, err := communicator.NewFSRedisQueueCommunicator(
		config.FsqDir,
		config.FsqExpertOutputQueue,
		config.FsqRedisAddr,
	)
	if err != nil {
		return nil, err
	}

	slComm, err := communicator.NewHTTPCommunicator(config.SlUrl)
	if err != nil {
		return nil, err
	}

	return newFeedbackFlowFromCommunicators(
		expertComm,
		slComm,
		config.BasicRateLimit,
	)
}

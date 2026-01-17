package config

import (
	logClient "github.com/EsanSamuel/sensory/LogClient"
)

var Logger, _ = logClient.New("7ed173cfa5dca9a4be50d9c6b6717b7a3b151839b5042ad592a59265fd8b68cf", ":9000")

package provider

import (
	"os"
	"yamda_go/internal/config"
	"yamda_go/internal/jsonlog"
)

// stuff common to unit tests
var (
	envConfigs, _ = config.New("./../../../debug.env")
	logger        = jsonlog.New(os.Stdout, jsonlog.LevelInfo)
)

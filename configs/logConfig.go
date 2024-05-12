package configs

import (
	"encoding/json"

	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func LoggerConfig() {
	rawJSON := []byte(`{
	  "level": "info",
	  "encoding": "json",
	  "outputPaths": ["stdout", "/tmp/logs"],
	  "errorOutputPaths": ["stderr","/tmp/error"],
	  "initialFields": {"name": "ui-backend"},
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger := zap.Must(cfg.Build())
	sugar := logger.Sugar()
	defer logger.Sync()
	Log = sugar

}

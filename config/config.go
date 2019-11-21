package config

const Version = "heplify-server 1.11"

var Setting HeplifyServer

type HeplifyServer struct {
	BrokerAddr                string   `default:"127.0.0.1:4222"`
	BrokerTopic               string   `default:"heplify.server.metric.1"`
	BrokerQueue               string   `default:"hep.metric.queue.1"`
	TargetIP                  string   `default:""`
	SIPErrorCode              string   `default:""`
	St2URL                    string   `default:""`
	St2ApiKey                 string   `default:""`
	LogDbg                    string   `default:""`
	LogLvl                    string   `default:"info"`
	LogStd                    bool     `default:"false"`
	LogSys                    bool     `default:"false"`
	Config                    string   `default:"./heplify-server.toml"`
}

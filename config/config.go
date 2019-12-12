package config

const Version = "heplify-server 1.11"

var Setting HeplifyServer

type HeplifyServer struct {
	BrokerAddr                string   `default:"127.0.0.1:4222"`
	BrokerTopic               string   `default:"heplify.server.metric.1"`
	BrokerQueue               string   `default:"hep.metric.queue.1"`
	Target                    string   `default:""`
	St2URL                    string   `default:""`
	St2ApiKey                 string   `default:""`
	LogDbg                    string   `default:""`
	LogLvl                    string   `default:"info"`
	LogStd                    bool     `default:"false"`
	LogSys                    bool     `default:"false"`
	Config                    string   `default:"./heplify-server.toml"`
	
	/*
	Target will be split up into three portion 
	1st will be the IP address to track
	2nd multiple SIP error code to match for that particular IP address
	3rd ; to mark the end
	IP address,SIPErrorCode,SIPErrorCode,SIPErrorCode,....;IP address,SIPErrorCode,SIPErrorCode,SIPErrorCode,....;IP address,SIPErrorCode,SIPErrorCode,SIPErrorCode,....
	example: 172.10.10.10,422,503,604;192.168.1.1,303,333,404;
	*/
}

package metric

import (
	"fmt"
	"strings"
	"crypto/tls"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/games130/logp"
	"github.com/games130/heplify-server-st2webhook/config"
	"github.com/games130/heplify-server-st2webhook/decoder"

)

type Prometheus struct {
	Target          []string
	TargetMap        map[string]map[string]string
}

func (p *Prometheus) setup() (err error) {
	//data coming in should be example: 172.10.10.10,422,503,604;192.168.1.1,303,333,404;
	//after split you will have map of 172.10.10.10,422,503,604  and   192.168.1.1,303,333,404
	p.Target = strings.Split(cutSpace(config.Setting.Target), ";")
	
	if p.Target != nil {
		p.TargetMap = make(map[string]map[string]string)
		for i := 0; i < len(p.Target); i++ {
			//after split you will have array of 172.10.10.10 422 503 604  and   array of 192.168.1.1 303 333 404
			tempSIPErrorCode = strings.Split(cutSpace(p.Target[i]), ",")
			tempSIPErrorCodeMap = make(map[string]string)
			for k := 1; k < len(tempSIPErrorCode); k++ {
				tempSIPErrorCodeMap[tempSIPErrorCode[k]] = tempSIPErrorCode[k]
			}
			p.TargetMap[tempSIPErrorCode[0]] = tempSIPErrorCodeMap
		}
		
		//in the end your map data should look like this:
		//map[172.10.10.10:map[422:422 503:503 604:604]    192.168.1.1:map[303:303 333:333 404:404]]
		//can query:
		//fmt.Println(mapData["172.10.10.10"])
		//fmt.Println(mapData["172.10.10.10"]["422"])
	} else {
		logp.Info("Target cannot be empty")
		return fmt.Errorf("faulty Target")
	}
	return err
}

func (p *Prometheus) expose(hCh chan *decoder.HEP) {
	for pkt := range hCh {
		if pkt != nil && pkt.ProtoType == 1 {
			//If source IP matches what we want to track proceed to check ErrorCode
			_, sOk := p.TargetMap[pkt.SrcIP]
			if sOk{
				//if Error code matches what we want to track proceed to generate WebHook to StackStorm
				_, codeOk := p.TargetMap[pkt.SrcIP][pkt.FirstMethod]
				if codeOk {
					p.generateWebhook(pkt)
				}
			}
			
			//If destination IP matches what we want to track proceed to check ErrorCode
			_, dOk := p.TargetMap[pkt.DstIP]
			if dOk {
				//if Error code matches what we want to track proceed to generate WebHook to StackStorm
				_, codeOk := p.TargetMap[pkt.DstIP][pkt.FirstMethod]
				if codeOk {
					p.generateWebhook(pkt)
				}
			}
		}
	}
}

func (p *Prometheus) generateWebhook(pkt *decoder.HEP) {
	transCfg := &http.Transport{
                 TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
    }
	var stringData string
	stringData = fmt.Sprintf("{\"SrcIP\": \"%s\",
							   \"DstIP\": \"%s\",
							   \"Tsec\": \"%s\",
							   \"CID\": \"%s\",
							   \"CseqMethod\": \"%s\",
							   \"FirstMethod\": \"%s\",
							   \"CallID\": \"%s\",
							   \"FromUser\": \"%s\",
							   \"ReasonVal\": \"%s\",
							   \"ToUser\": \"%s\",
							   \"Timestamp\": \"%s\",
		                       \"XCallID\": \"%s\",
							   \"CseqMethod\": \"%s\
							   "}",
							   pkt.SrcIP,
							   pkt.DstIP,
							   pkt.Tsec,
							   pkt.CID,
							   pkt.CseqMethod,
							   pkt.FirstMethod,
							   pkt.CallID,
							   pkt.FromUser,
							   pkt.ReasonVal,
							   pkt.ToUser,
							   pkt.Timestamp,
							   pkt.XCallID,
							   pkt.CseqMethod,
							   )
	var data = []byte(stringData)

	req, err := http.NewRequest("POST", config.Setting.St2URL, bytes.NewBuffer(data))
	req.Header.Set("St2-Api-Key", config.Setting.St2ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Transport: transCfg}
	debug(httputil.DumpRequestOut(req, true))
	resp, err := client.Do(req)
	if err != nil {
		logp.Err("%s\n\n", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	logp.Info("Response: ", string(body))
	resp.Body.Close()
}

func debug(data []byte, err error) {
    if err == nil {
        logp.Info("%s\n\n", data)
    } else {
        logp.Err("%s\n\n", err)
    }
}

package metric

import (
	"encoding/binary"
	"fmt"
	"strings"
	"sync"
	"regexp"
	"time"
	"os"
	"bufio"
	"crypto/tls"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/games130/logp"
	"github.com/games130/heplify-server-metric/config"
	"github.com/games130/heplify-server-metric/decoder"
	
)

type Prometheus struct {
	TargetIP           []string
	SIPErrorCode       []string
	TargetIPMap        map[string]string
	SIPErrorCodeMap    map[string]string
}

func (p *Prometheus) setup() (err error) {
	p.TargetIP = strings.Split(cutSpace(config.Setting.TargetIP), ",")
	p.SIPErrorCode = strings.Split(cutSpace(config.Setting.SIPErrorCode), ",")
	
	if p.TargetIP != nil && p.SIPErrorCode != nil {
		p.TargetIPMap = make(map[string]string)
		for i := 0; i < len(p.TargetIP); i++ {
			p.TargetIPMap[p.TargetIP[i]] = p.TargetIP[i]
		}
		
		p.SIPErrorCodeMap = make(map[string]string)
		for i := 0; i < len(p.SIPErrorCode); i++ {
			p.SIPErrorCodeMap[p.SIPErrorCode[i]] = p.SIPErrorCode[i]
		}
	} else {
		logp.Info("TargetIP and SIPErrorCode cannot be empty")
		return fmt.Errorf("faulty TargetIP or SIPErrorCode")
	}
	return err
}

func (p *Prometheus) expose(hCh chan *decoder.HEP) {
	for pkt := range hCh {
		if pkt != nil && pkt.ProtoType == 1 {
			dt, dOk := p.TargetIPMap[pkt.DstIP]
				if dOk {
					code, codeOk := p.SIPErrorCodeMap[pkt.FirstMethod]
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
	stringData = fmt.Sprintf("{\"FirstMethod\": \"%s\",\"ToUser\": \"%s\",\"FromUser\": \"%s\",\"CseqMethod\": \"%s\"}", pkt.FirstMethod, pkt.ToUser, pkt.FromUser, pkt.CseqMethod)
	var data = []byte(stringData)
	
	req, err := http.NewRequest("POST", config.Setting.St2URL, bytes.NewBuffer(data))
	req.Header.Set("St2-Api-Key", config.Setting.St2ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Transport: transCfg}
	debug(httputil.DumpRequestOut(req, true))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	logp.Info("Response: ", string(body))
	resp.Body.Close()	
}

func debug(data []byte, err error) {
    if err == nil {
        logp.Info("%s\n\n", data)
    } else {
        logp.Error("%s\n\n", err)
    }
}

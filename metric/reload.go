package metric

import (
	"io/ioutil"
	"strings"
	"unicode"

	"github.com/games130/logp"
	"github.com/games130/heplify-server-st2webhook/config"
)

func cutSpace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func (p *Prometheus) reload() {
	var fsTargetIP []string
	var fsTargetName []string

	fb, err := ioutil.ReadFile(config.Setting.Config)
	if err != nil {
		logp.Err("%v", err)
		return
	}

	fs := cutSpace(string(fb))

	if si := strings.Index(fs, "TargetIP=\""); si > -1 {
		s := si + len("TargetIP=\"")
		e := strings.Index(fs[s:], "\"")
		if e >= 7 {
			fsTargetIP = strings.Split(fs[s:s+e], ",")
		}
	}
	if si := strings.Index(fs, "SIPErrorCode=\""); si > -1 {
		s := si + len("SIPErrorCode=\"")
		e := strings.Index(fs[s:], "\"")
		if e > 0 {
			fsTargetName = strings.Split(fs[s:s+e], ",")
		}
	}

	if fsTargetIP != nil && fsTargetName != nil && len(fsTargetIP) == len(fsTargetName) {
		//p.TargetConf.Lock()  //not sure what this is for
		p.TargetIP = fsTargetIP
		p.SIPErrorCode = fsTargetName
		
		p.TargetIPMap = make(map[string]string)
		for i := 0; i < len(p.TargetIP); i++ {
			p.TargetIPMap[p.TargetIP[i]] = p.TargetIP[i]
		}
		
		p.SIPErrorCodeMap = make(map[string]string)
		for i := 0; i < len(p.SIPErrorCode); i++ {
			p.SIPErrorCodeMap[p.SIPErrorCode[i]] = p.SIPErrorCode[i]
		}
		
		//p.TargetConf.Unlock()   //not sure what this is for
		logp.Info("successfully reloaded TargetIP: %#v", fsTargetIP)
		logp.Info("successfully reloaded SIPErrorCode: %#v", fsTargetName)
	} else {
		logp.Info("failed to reload TargetIP: %#v", fsTargetIP)
		logp.Info("failed to reload SIPErrorCode: %#v", fsTargetName)
	}
}

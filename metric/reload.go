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
	var fsTarget []string

	fb, err := ioutil.ReadFile(config.Setting.Config)
	if err != nil {
		logp.Err("%v", err)
		return
	}

	fs := cutSpace(string(fb))

	if si := strings.Index(fs, "Target=\""); si > -1 {
		s := si + len("Target=\"")
		e := strings.Index(fs[s:], "\"")
		if e >= 7 {
			fsTarget = strings.Split(fs[s:s+e], ",")
		}
	}

	if fsTarget != nil {
		//p.TargetConf.Lock()  //not sure what this is for
		p.Target = fsTarget
		p.TargetMap = make(map[string]map[string]string)
		for i := 0; i < len(p.Target); i++ {
			//after split you will have array of 172.10.10.10 422 503 604  and   array of 192.168.1.1 303 333 404
			tempSIPErrorCode := strings.Split(cutSpace(p.Target[i]), ",")
			tempSIPErrorCodeMap := make(map[string]string)
			for k := 1; k < len(tempSIPErrorCode); k++ {
				tempSIPErrorCodeMap[tempSIPErrorCode[k]] = tempSIPErrorCode[k]
			}
			p.TargetMap[tempSIPErrorCode[0]] = tempSIPErrorCodeMap
		}


		
		//p.TargetConf.Unlock()   //not sure what this is for
		logp.Info("successfully reloaded Target: %#v", fsTarget)
	} else {
		logp.Info("failed to reload Target: %#v", fsTarget)
	}
}

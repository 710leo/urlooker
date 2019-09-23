package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/710leo/urlooker/dataobj"
	"github.com/710leo/urlooker/modules/web/g"

	"github.com/miekg/dns"
	"github.com/toolkits/sys"
)

func LookupIP(domain string, timeout int) ([]string, error) {
	var c dns.Client
	var err error
	ips := []string{}
	domain = strings.TrimRight(domain, ".") + "."
	c.DialTimeout = time.Duration(timeout) * time.Millisecond
	c.ReadTimeout = time.Duration(timeout) * time.Millisecond
	c.WriteTimeout = time.Duration(timeout) * time.Millisecond

	m := new(dns.Msg)
	m.SetQuestion(domain, dns.TypeA)
	if g.Config.DNS == "" {
		return ips, errors.New("dns is nil")
	}

	ret, _, err := c.Exchange(m, g.Config.DNS)
	if err != nil {
		domain = strings.TrimRight(domain, ".")
		e := fmt.Sprintf("lookup error: %s, %s", domain, err.Error())
		return ips, errors.New(e)
	}

	for _, i := range ret.Answer {
		result := strings.Split(i.String(), "\t")
		if result[3] == "A" && IsIP(result[4]) {
			ips = append(ips, result[4])
		}
	}

	return ips, err
}

func InternalDns(domain string) ([]dataobj.IpIdc, error) {
	ret := []dataobj.IpIdc{}
	out, err := sys.CmdOut(g.Config.InternalDns.CMD, domain) //InternalDns.CMD 文件内容需要自己定制
	if err != nil {
		return ret, err
	}
	lines := strings.Split(out, "\n")
	for _, l := range lines {
		arr := strings.Split(l, ",")
		if len(arr) == 2 {
			tmp := dataobj.IpIdc{
				Idc: arr[0],
				Ip:  arr[1],
			}
			ret = append(ret, tmp)
		}
	}
	return ret, err
}

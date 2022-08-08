package utils

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/deweppro/go-errors"
)

var (
	ErrInvalidIP = errors.New("invalid ip")
)

func ValidateDNS(ip string) (string, error) {
	if _, _, err := net.SplitHostPort(ip); err != nil {
		if v := net.ParseIP(ip); v != nil {
			return net.JoinHostPort(ip, "53"), nil
		}
		return "", ErrInvalidIP
	}
	return ip, nil
}

func ValidateDNSs(list []string) (result []string) {
	for _, ip := range list {
		if v, er := ValidateDNS(ip); er == nil {
			result = append(result, v)
		}
	}
	return
}

func DecodeIPs(data string) (ip4, ip6 []string) {
	list := strings.Split(data, ",")
	for _, host := range list {
		host = strings.TrimSpace(host)
		ip, port, err := net.SplitHostPort(host)
		if err != nil {
			ip = host
		}
		v := net.ParseIP(ip)
		if v == nil {
			continue
		}
		ip = v.String()
		if len(port) > 0 {
			host = net.JoinHostPort(ip, port)
		}
		if strings.Contains(ip, ":") {
			ip6 = append(ip6, host)
			continue
		}
		ip4 = append(ip4, host)
	}
	return
}

func EncodeIPs(ip4, ip6 []string) string {
	return strings.Join(append(ip4, ip6...), ", ")
}

func Interval(ctx context.Context, interval time.Duration, call func(context.Context)) {
	call(ctx)

	go func() {
		tick := time.NewTicker(interval)
		defer tick.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				call(ctx)
			}
		}
	}()
}

func StringError(err error) string {
	if err == nil {
		return "<nil>"
	}
	return err.Error()
}

func Tag(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

var domainRex = regexp.MustCompile(`^(?i)[a-z0-9-]+(\.[a-z0-9-]+)+\.?$`)

func ValidateDomain(domain string) (string, error) {
	domain = strings.TrimSpace(domain)
	if !domainRex.MatchString(domain) {
		return "", fmt.Errorf("invalid domain")
	}
	domain = strings.TrimRight(domain, ".")
	domain = strings.ToLower(domain)
	return domain + ".", nil
}

func Retry(count int, ttl time.Duration, call func() error) {
	for i := 0; i < count; i++ {
		err := call()
		if err != nil {
			time.Sleep(ttl)
			continue
		}
		return
	}
	return
}

func ReadClose(r io.ReadCloser) ([]byte, error) {
	b, err := ioutil.ReadAll(r)
	err = errors.Wrap(err, r.Close())
	return b, err
}

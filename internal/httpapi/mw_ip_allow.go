package httpapi

import (
	"net"
	"net/http"
	"strings"
)

var privateNetworks = mustParseCIDRs([]string{
	"127.0.0.0/8",
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"::1/128",
	"fc00::/7",
})

func AllowPrivateNetworks() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIPKey(r)
			parsed := net.ParseIP(strings.TrimSpace(ip))
			if parsed == nil || !ipInAnyCIDR(parsed, privateNetworks) {
				writeError(w, http.StatusForbidden, "forbidden")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func mustParseCIDRs(values []string) []*net.IPNet {
	out := make([]*net.IPNet, 0, len(values))
	for _, v := range values {
		_, n, err := net.ParseCIDR(v)
		if err == nil {
			out = append(out, n)
		}
	}
	return out
}

func ipInAnyCIDR(ip net.IP, nets []*net.IPNet) bool {
	for _, n := range nets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}


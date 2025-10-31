package strategy

import (
	"fmt"
	"github.com/eolinker/eosc"
	"net"
	"strings"
)

type IFilterChecker interface {
	Name() string
	Check(value []string) error
}

type FilterChecker struct {
	checks eosc.Untyped[string, IFilterChecker]
}

func (f *FilterChecker) Set(check IFilterChecker) {
	f.checks.Set(check.Name(), check)
}

func (f *FilterChecker) Get(name string) (IFilterChecker, bool) {
	return f.checks.Get(name)
}

func NewFilterChecker() *FilterChecker {
	return &FilterChecker{checks: eosc.BuildUntyped[string, IFilterChecker]()}
}

var (
	filterChecker = NewFilterChecker()
)

func RegisterChecker(check IFilterChecker) {
	filterChecker.Set(check)
}

func GetChecker(name string) (IFilterChecker, bool) {
	return filterChecker.Get(name)
}

type ipFilterChecker struct {
	name string
}

func newIpFilterChecker() *ipFilterChecker {
	return &ipFilterChecker{name: "ip"}
}

func (c *ipFilterChecker) Name() string {
	return c.name
}

func (c *ipFilterChecker) Check(value []string) error {
	for _, v := range value {
		// 判断v是否是IP或CIDR格式的IP
		if !isIPv4(v) && !isCIDR(v) {
			return fmt.Errorf("invalid ip: %s", v)
		}
	}
	return nil
}

func isIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && strings.Count(ip, ":") < 1
}

// isCIDR checks if a string is a valid CIDR notation
func isCIDR(cidr string) bool {
	_, _, err := net.ParseCIDR(cidr)
	return err == nil
}

type apiFilterChecker struct {
	name string
}

func newApiFilterChecker() *apiFilterChecker {
	return &apiFilterChecker{name: "api"}
}

func (c *apiFilterChecker) Name() string {
	return c.name
}

func (c *apiFilterChecker) Check(value []string) error {
	return nil
}

type serviceFilterChecker struct {
	name string
}

func newServiceFilterChecker() *serviceFilterChecker {
	return &serviceFilterChecker{name: "service"}
}

func (c *serviceFilterChecker) Name() string {
	return c.name
}

func (c *serviceFilterChecker) Check(value []string) error {
	return nil
}

func init() {
	RegisterChecker(newIpFilterChecker())
	RegisterChecker(newApiFilterChecker())
	RegisterChecker(newServiceFilterChecker())
}

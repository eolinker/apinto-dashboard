package enum

import "fmt"

type ServiceApiMethod int

const (
	ServiceApiMethodNone = iota
	ServiceApiMethodPOST
	ServiceApiMethodGET
	ServiceApiMethodPUT
	ServiceApiMethodDELETE
	ServiceApiMethodHEAD
	ServiceApiMethodOPTIONS
	ServiceApiMethodPATCH
	ServiceApiMethodALL
)

var (
	serviceApiMethodNames = map[ServiceApiMethod]string{
		ServiceApiMethodPOST:    "POST",
		ServiceApiMethodGET:     "GET",
		ServiceApiMethodPUT:     "PUT",
		ServiceApiMethodDELETE:  "DELETE",
		ServiceApiMethodHEAD:    "HEAD",
		ServiceApiMethodOPTIONS: "OPTIONS",
		ServiceApiMethodPATCH:   "PATCH",
	}
	ServiceApiMethodIndex = map[string]ServiceApiMethod{
		"POST":    ServiceApiMethodPOST,
		"GET":     ServiceApiMethodGET,
		"PUT":     ServiceApiMethodPUT,
		"DELETE":  ServiceApiMethodDELETE,
		"HEAD":    ServiceApiMethodHEAD,
		"OPTIONS": ServiceApiMethodOPTIONS,
		"PATCH":   ServiceApiMethodPATCH,
	}
)

func (s ServiceApiMethod) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint("\"", s.String(), "\"")), nil
}

func (s ServiceApiMethod) String() string {
	if s >= ServiceApiMethodALL {
		return "unknown"
	}
	return serviceApiMethodNames[s]
}

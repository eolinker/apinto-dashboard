package enum

import "fmt"

type EnumEnv int

const (
	EnumEnvNone = iota
	EnumEnvPRO
	EnumEnvFAT
	EnumEnvDEV
	EnumEnvUAT
	EnumEnvAll
)

var (
	enumValueNames = map[EnumEnv]string{
		EnumEnvNone: "NONE",
		EnumEnvPRO:  "PRO",
		EnumEnvFAT:  "FAT",
		EnumEnvDEV:  "DEV",
		EnumEnvUAT:  "UAT",
	}
	EnumValueList = []EnumEnv{EnumEnvPRO, EnumEnvFAT, EnumEnvDEV, EnumEnvUAT}
)

type EnumEnvOut struct {
	Name  string  `json:"name"`
	Value EnumEnv `json:"value"`
}

func (e EnumEnv) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint("\"", e.String(), "\"")), nil
}

func (e EnumEnv) String() string {
	if e >= EnumEnvAll {
		return "unknown"
	}
	return enumValueNames[e]
}

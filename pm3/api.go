package pm3

type ApiAuthority int

func (a ApiAuthority) String() string {
	switch a {
	case Internal:
		return "internal"
	case Private:
		return "private"
	case Public:
		return "public"
	case Anonymous:
		return "anonymous"

	}
	return "unset"
}

const (
	UnSet ApiAuthority = iota
	Internal
	Public
	Private
	Anonymous
	maxApiAuthority
)

var (
	authorityValues = make(map[string]ApiAuthority)
)

func init() {
	for i := UnSet; i < maxApiAuthority; i++ {
		authorityValues[i.String()] = i
	}
}
func ParseApiAuthority(v string) ApiAuthority {
	authority, h := authorityValues[v]
	if h {
		return authority
	}
	return UnSet
}

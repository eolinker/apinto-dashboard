package machine_code

const Salt = "eolink-apserver-"

var machineCode string

func SetMachineCode(mac string) {
	machineCode = mac
}

func GetMachineCode() string {
	return machineCode
}

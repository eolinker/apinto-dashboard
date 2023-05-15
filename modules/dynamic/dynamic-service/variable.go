package dynamic_service

import (
	"strings"
)

const (
	dollarSign = 36  // $
	leftSign   = 123 // {
	rightSign  = 125 // }
)

const (

	//CurrentStatus 普通状态
	CurrentStatus = iota

	//ReadyStatus 预备状态
	ReadyStatus

	//InputStatus 输入状态
	InputStatus

	//EndInputStatus 结束输入状态
	EndInputStatus
)

func parseVariables(str string) []string {
	varBuilder := strings.Builder{}
	variableMap := make(map[string]struct{})
	status := CurrentStatus
	for _, s := range str {
		oldStatus := status
		status = toggleStatus(status, s)
		switch status {
		case InputStatus:
			if oldStatus == ReadyStatus {
				// 刚切换状态，忽略此时的字符
				continue
			}
			varBuilder.WriteRune(s)
		case EndInputStatus:
			variableMap[strings.TrimSpace(varBuilder.String())] = struct{}{}
			varBuilder.Reset()
			status = CurrentStatus
		}
	}
	variables := make([]string, 0, len(variableMap))
	for key := range variableMap {
		variables = append(variables, key)
	}
	return variables
}

func toggleStatus(status int, c rune) int {
	switch status {
	case CurrentStatus, EndInputStatus:
		if c == dollarSign {
			return ReadyStatus
		}
		return CurrentStatus
	case ReadyStatus:
		if c == leftSign {
			return InputStatus
		}
		return CurrentStatus
	case InputStatus:
		if c == rightSign {
			return EndInputStatus
		}
		return InputStatus
	}
	return status
}

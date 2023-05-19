/*
 * Copyright (c) 2023. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package version

import (
	"bytes"
	"github.com/go-basic/uuid"

	"os"
)

var (
	dashboardId string
)

func init() {
	data, err := os.ReadFile("../work/apserver.id")
	if err != nil {
		dashboardId = uuid.New()
		os.WriteFile("../work/apserver.id", []byte(dashboardId), 0)
		return
	}

	dashboardId = string(bytes.Split(data, []byte("\n"))[0])

}
func DashboardId() string {
	return dashboardId
}

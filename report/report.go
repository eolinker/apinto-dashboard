package report

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/config"
	"github.com/eolinker/eosc/common/bean"
	"github.com/go-basic/uuid"
	"net/http"
	"os"
	"path/filepath"
)

var (
	myId string
)

func init() {
	bean.AddInitializingBeanFunc(func() {
		Deploy(config.GetLogDir())
	})
}
func Deploy(dir string) {
	id, isNew := createId(dir)
	myId = id
	go func() {
		if isNew {
			http.Get(fmt.Sprintf("%s/report/deploy/c?id=%s", reportAddr, id))
		}
	}()

}

func createId(dir string) (string, bool) {
	abs, _ := filepath.Abs(dir)

	file := filepath.Join(abs, ".console.id")
	data, err := os.ReadFile(file)

	if err != nil && os.IsNotExist(err) {
		newId := uuid.New()
		f, err := os.Create(file)
		if err != nil {
			return newId, true
		}
		defer f.Close()
		f.WriteString(myId)
		return newId, true
	}
	return string(data), false
}

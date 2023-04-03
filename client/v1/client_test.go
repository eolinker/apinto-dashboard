package v1

import (
	"fmt"
	"testing"
)

func TestWorker(t *testing.T) {

	client, err := NewClient([]string{"http://172.27.75.146:9400"})

	err = client.ForApp().Patch("lGmKT7RvbhCRlhI5", map[string]interface{}{"disable": true})
	fmt.Println(err)

	//key, err := client.ForVariable().GetByKey(namespace, "test1")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(key)

	//all, err := client.ForVariable().GetAll()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(all)

}

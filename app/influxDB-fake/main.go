package main

import "log"

func main() {
	//创建服务
	//err := createService(50)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//批量创建API
	//err := createAPIBatch(900)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//获取api信息并且写入文件
	//err := getApiInfoAndSave()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	err := WriteFakeDataToInfluxDB()
	if err != nil {
		log.Fatal(err)
	}
	//err := writeToLineProtocol()
	//if err != nil {
	//	log.Fatal(err)
	//}

}

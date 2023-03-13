该目录是用来生成测试数据的, 具体测试数据可以配置config.yaml文件

可用的函数：

* 生成指定数量的api
* 生成指定数量的service
* 通过influxDB-client写入数据（太慢
* writeToCsv， 生成含有测试数据的csv文件，但是通过influx write 写入数据库会缺失tag，可能是官方bug
* writeToLineProtocol, 生成line-protocol格式的测试数据文件，生成目录在./expot/line_protocol。可以通过influx write写入数据库，比用client写入快很多倍。



建议使用writeToLineProtocol，生成xxx.txt 测试数据文件，然后用gzip压缩一下，再把生成的xxx.txt.gz 文件拉到influxdb容器内，在容器内执行以下语句：

```
influx write -b "test_apinto" -f ./test_request_data.txt.gz  --org="eolink" --token="zpS_8Zuf9JQ0wOG_cAZAp8heU2PyN9sat-wtwfgSH90CUsRzIxrWYGSqZAobW4wejH5BiT8DkVFCxP5PwHklvw=="
```

参考链接：https://docs.influxdata.com/influxdb/cloud/reference/cli/influx/write/
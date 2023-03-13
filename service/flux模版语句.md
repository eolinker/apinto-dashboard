调用量统计折线趋势图模版：

输入： 开始、结束时间戳，时间间隔

返回的列表： 时间戳，请求总量，请求成功数，转发总数，转发成功数，状态码4xx数，状态码5xx



若要看某个api/app/upstream的趋势图，或者{api}-{app}，添加filter条件即可

```
		import "join"
		req_data = from(bucket: "apinto")
	          |> range(start: 1670810987, stop: 1671674987)
	          |> filter(fn: (r) => r._measurement == "request" and r._field == "timing")
			  |> group()
			  |> aggregateWindow(
				every: 1s, 
				fn: (column, tables=<-) => tables |> reduce(
            		identity: {req_total: 0, req_success_total: 0, status_4xx: 0, status_5xx: 0},
            		fn: (r, accumulator) => ({
              	 	 req_total: accumulator.req_total + 1,
               	 	 req_success_total: if r.status =~ /[1-3]/ then accumulator.req_success_total + 1 else accumulator.req_success_total,
            		 status_4xx: if r.status =~ /4/ then accumulator.status_4xx + 1 else accumulator.status_4xx,
            		 status_5xx: if r.status =~ /5/ then accumulator.status_5xx + 1 else accumulator.status_5xx,
					}),
         	  	)
			 )

		proxy_data = from(bucket: "apinto")
	          |> range(start: 1670810987, stop: 1671674987)
	          |> filter(fn: (r) => r._measurement == "proxy" and r._field == "timing")
			  |> group()
			  |> aggregateWindow(
				every: 1s, 
				fn: (column, tables=<-) => tables |> reduce(
            		identity: {proxy_total: 0, proxy_success_total: 0},
            		fn: (r, accumulator) => ({
              	 	 proxy_total: accumulator.proxy_total + 1,
               	 	 proxy_success_total: if r.status =~ /[1-3]/ then accumulator.proxy_success_total + 1 else accumulator.proxy_success_total,
					}),
         	  	)
			 )
		join.left(
			left: req_data,
			right: proxy_data,
			on: (l, r) => l._time == r._time,
			as: (l, r) => ({l with req_total: l.req_total, req_success_total: l.req_success_total, proxy_total: r.proxy_total, proxy_success_total: r.proxy_success_total, status_4xx: l.status_4xx, status_5xx: l.status_5xx}),
		)
```



以api调用统计为例的模版：

输入：请求表名，转发表名， 开始时间戳，响应时间戳

返回： apiID、请求总数、请求成功数、转发总数、转发成功数、失败状态码数、平均响应时间(ms)、最大响应时间(ms)、最小响应时间(ms)、平均请求流量（KB）、最大请求流量（KB）、最小请求流量（KB）、总请求流量、总响应流量

若要其他模块的总调用统计， 将下列模版的api替换成相应模块即可。

若要针对某个api下的应用，就加多个过滤。r.api == "api1" ,然后以app为分组。

```
		import "join"
		request_data = from(bucket: "apinto")
		     |> range(start: 1670810987, stop: 1671674987)
	         |> filter(fn: (r) => r._measurement == "request")
			 |> group(columns:["api"])
			 |> pivot(rowKey: ["_time","requestId","status"], columnKey: ["_field"], valueColumn: "_value")

		resp_timing_data = request_data
			|> duplicate(column: "timing", as: "_value")
		req_message_data = request_data
			|> duplicate(column: "request", as: "_value")
		resp_message_data = request_data
			|> duplicate(column: "response", as: "_value")

		proxy_data = from(bucket: "apinto")
		     |> range(start: 1670810987, stop: 1671674987)
	         |> filter(fn: (r) => r._measurement == "proxy" and r._field == "timing")
			 |> group(columns:["api"])

		request_total = req_message_data
			 |> count()
		request_success_total = req_message_data
			|> filter(fn: (r) => r.status =~ /[1-3]/)
			|> count()

		proxy_total = proxy_data
			 |> count()
		proxy_success_total = proxy_data
			|> filter(fn: (r) => r.status =~ /[1-3]/)
			|> count()

		status_fail_total = req_message_data
			|> filter(fn: (r) => r.status !~ /[1-3]/)
			|> count()

		avg_resp_timing = resp_timing_data
			|> mean()
		max_resp_timing = resp_timing_data
			|> max()
		min_resp_timing = resp_timing_data
			|> min()

		avg_req_mess = req_message_data
			|> mean()
		max_req_mess = req_message_data
			|> max()
		min_req_mess = req_message_data
			|> min()
		req_mess_total = req_message_data
			|> sum()

		resp_mess_total = resp_message_data
			|> sum()

		merge1 = join.left(
			left: request_total,
			right: request_success_total,
			on: (l, r) => l.api == r.api,
			as: (l, r) => ({l with request_total: l._value, req_success_total: r._value}),
		)
		merge2 = join.left(
			left: merge1,
			right: proxy_total,
			on: (l, r) => l.api == r.api,
			as: (l, r) => ({l with proxy_total: r._value}),
		)
		merge3 = join.left(
			left: merge2,
			right: proxy_success_total,
			on: (l, r) => l.api == r.api,
			as: (l, r) => ({l with proxy_success_total: r._value}),
		)
		merge4 = join.left(
			left: merge3,
			right: status_fail_total,
			on: (l, r) => l.api == r.api,
			as: (l, r) => ({l with status_fail_total: r._value}),
		)		
		merge5 = join.left(
			left: merge4,
			right: avg_resp_timing,
			on: (l, r) => l.api == r.api,
			as: (l, r) => ({l with avg_resp_timing: r._value}),
		)		
		merge6 = join.left(
			left: merge5,
			right: max_resp_timing,
			on: (l, r) => l.api == r.api,
			as: (l, r) => ({l with max_resp_timing: r._value}),
		)		
		merge7 = join.left(
			left: merge6,
			right: min_resp_timing,
			on: (l, r) => l.api == r.api,
			as: (l, r) => ({l with min_resp_timing: r._value}),
		)		
		merge8 = join.left(
			left: merge7,
			right: avg_req_mess,
			on: (l, r) => l.api == r.api,
			as: (l, r) => ({l with avg_req_mess: r._value}),
		)		
		merge9 = join.left(
			left: merge8,
			right: max_req_mess,
			on: (l, r) => l.api == r.api,
			as: (l, r) => ({l with max_req_mess: r._value}),
		)		
		merge10 = join.left(
			left: merge9,
			right: min_req_mess,
			on: (l, r) => l.api == r.api,
			as: (l, r) => ({l with min_req_mess: r._value}),
		)		
		merge11 = join.left(
			left: merge10,
			right: req_mess_total,
			on: (l, r) => l.api == r.api,
			as: (l, r) => ({l with min_req_mess: r._value}),
		)		
		join.left(
			left: merge11,
			right: resp_mess_total,
			on: (l, r) => l.api == r.api,
			as: (l, r) => ({l with resp_mess_total: r._value}),
		)
```




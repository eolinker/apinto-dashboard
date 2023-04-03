package common

import (
	"context"
	"fmt"
	"github.com/eolinker/eosc/log"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

func CommFluxGroup(ctx context.Context, queryApi api.QueryAPI, group, filter string, startTime, endTime int64) (map[string]map[string]interface{}, error) {

	requestSql := fmt.Sprintf(`import "join"
		request_data = from(bucket: "apinto")
		     |> range(start: %d, stop: %d)
	         |> filter(fn: (r) => r._measurement == "request")
			 %s	
			 |> group(columns:["%s"])
			 |> pivot(rowKey: ["_time","requestId","status"], columnKey: ["_field"], valueColumn: "_value")`, startTime, endTime, filter, group)

	proxySql := fmt.Sprintf(`proxy_data = from(bucket: "apinto")
		     |> range(start: %d, stop: %d)
	         |> filter(fn: (r) => r._measurement == "proxy" and r._field == "timing")
             %s
			 |> group(columns:["%s"])`, startTime, endTime, filter, group)

	querySql := fmt.Sprintf(`
		%s
		resp_timing_data = request_data
			|> duplicate(column: "timing", as: "_value")
		req_message_data = request_data
			|> duplicate(column: "request", as: "_value")
		resp_message_data = request_data
			|> duplicate(column: "response", as: "_value")

		%s

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
			on: (l, r) => l.%s == r.%s,
			as: (l, r) => ({l with request_total: l._value, req_success_total: r._value}),
		)
		merge2 = join.left(
			left: merge1,
			right: proxy_total,
			on: (l, r) => l.%s == r.%s,
			as: (l, r) => ({l with proxy_total: r._value}),
		)
		merge3 = join.left(
			left: merge2,
			right: proxy_success_total,
			on: (l, r) => l.%s == r.%s,
			as: (l, r) => ({l with proxy_success_total: r._value}),
		)
		merge4 = join.left(
			left: merge3,
			right: status_fail_total,
			on: (l, r) => l.%s == r.%s,
			as: (l, r) => ({l with status_fail_total: r._value}),
		)		
		merge5 = join.left(
			left: merge4,
			right: avg_resp_timing,
			on: (l, r) => l.%s == r.%s,
			as: (l, r) => ({l with avg_resp_timing: r._value}),
		)		
		merge6 = join.left(
			left: merge5,
			right: max_resp_timing,
			on: (l, r) => l.%s == r.%s,
			as: (l, r) => ({l with max_resp_timing: r._value}),
		)		
		merge7 = join.left(
			left: merge6,
			right: min_resp_timing,
			on: (l, r) => l.%s == r.%s,
			as: (l, r) => ({l with min_resp_timing: r._value}),
		)		
		merge8 = join.left(
			left: merge7,
			right: avg_req_mess,
			on: (l, r) => l.%s == r.%s,
			as: (l, r) => ({l with avg_req_mess: r._value}),
		)		
		merge9 = join.left(
			left: merge8,
			right: max_req_mess,
			on: (l, r) => l.%s == r.%s,
			as: (l, r) => ({l with max_req_mess: r._value}),
		)		
		merge10 = join.left(
			left: merge9,
			right: min_req_mess,
			on: (l, r) => l.%s == r.%s,
			as: (l, r) => ({l with min_req_mess: r._value}),
		)		
		merge11 = join.left(
			left: merge10,
			right: req_mess_total,
			on: (l, r) => l.%s == r.%s,
			as: (l, r) => ({l with min_req_mess: r._value}),
		)		
		join.left(
			left: merge11,
			right: resp_mess_total,
			on: (l, r) => l.%s == r.%s,
			as: (l, r) => ({l with resp_mess_total: r._value}),
		)`, requestSql, proxySql, group, group, group, group, group, group, group, group, group, group, group, group, group, group, group, group, group, group, group, group, group, group, group, group)

	log.Info("flux sql=", querySql)
	result, err := queryApi.Query(ctx, querySql)
	if err != nil {
		log.Errorf("flux err=", err)
		return nil, err
	}
	resultMap := make(map[string]map[string]interface{})
	for result.Next() {
		key := ""
		if v, ok := result.Record().Values()[group]; ok {
			key = v.(string)
		}
		resultMap[key] = result.Record().Values()
	}
	return resultMap, nil

}

func ServiceFluxGroup(ctx context.Context, queryApi api.QueryAPI, group, filter string, startTime, endTime int64) (map[string]map[string]interface{}, error) {

	querySql := fmt.Sprintf(`import "join"
             proxy_data = from(bucket: "apinto")
			 |> range(start: %d, stop: %d)
			  %s	
			 |> filter(fn: (r) => r._measurement == "proxy")
			 |> group(columns:["%s"])
			 |> pivot(rowKey: ["_time","requestId","status"], columnKey: ["_field"], valueColumn: "_value")

			proxy_timing_data = proxy_data
			   |> duplicate(column: "timing", as: "_value")
			req_message_data = proxy_data
			   |> duplicate(column: "request", as: "_value")
			resp_message_data = proxy_data
			   |> duplicate(column: "response", as: "_value")

			proxy_total = proxy_timing_data
				|> count()
			proxy_success_total = proxy_timing_data
			   |> filter(fn: (r) => r.status =~ /[1-3]/)
			   |> count()
			status_fail_total = proxy_timing_data
			   |> filter(fn: (r) => r.status !~ /[1-3]/)
			   |> count()

			avg_resp_timing = proxy_timing_data
			   |> mean()
			max_resp_timing = proxy_timing_data
			   |> max()
			min_resp_timing = proxy_timing_data
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
			   left: proxy_total,
			   right: proxy_success_total,
			   on: (l, r) => l.%s == r.%s,
			   as: (l, r) => ({l with proxy_total: l._value, proxy_success_total: r._value}),
			)
			merge2 = join.left(
			   left: merge1,
			   right: status_fail_total,
			   on: (l, r) => l.%s == r.%s,
			   as: (l, r) => ({l with status_fail_total: r._value}),
			)     
			merge3 = join.left(
			   left: merge2,
			   right: avg_resp_timing,
			   on: (l, r) => l.%s == r.%s,
			   as: (l, r) => ({l with avg_resp_timing: r._value}),
			)     
			merge4 = join.left(
			   left: merge3,
			   right: max_resp_timing,
			   on: (l, r) => l.%s == r.%s,
			   as: (l, r) => ({l with max_resp_timing: r._value}),
			)     
			merge5 = join.left(
			   left: merge4,
			   right: min_resp_timing,
			   on: (l, r) => l.%s == r.%s,
			   as: (l, r) => ({l with min_resp_timing: r._value}),
			)     
			merge6 = join.left(
			   left: merge5,
			   right: avg_req_mess,
			   on: (l, r) => l.%s == r.%s,
			   as: (l, r) => ({l with avg_req_mess: r._value}),
			)     
			merge7 = join.left(
			   left: merge6,
			   right: max_req_mess,
			   on: (l, r) => l.%s == r.%s,
			   as: (l, r) => ({l with max_req_mess: r._value}),
			)     
			merge8 = join.left(
			   left: merge7,
			   right: min_req_mess,
			   on: (l, r) => l.%s == r.%s,
			   as: (l, r) => ({l with min_req_mess: r._value}),
			)     
			merge9 = join.left(
			   left: merge8,
			   right: req_mess_total,
			   on: (l, r) => l.%s == r.%s,
			   as: (l, r) => ({l with min_req_mess: r._value}),
			)     
			join.left(
			   left: merge9,
			   right: resp_mess_total,
			   on: (l, r) => l.%s == r.%s,
			   as: (l, r) => ({l with resp_mess_total: r._value}),
			)`, startTime, endTime, filter, group, group, group, group, group, group, group, group, group, group, group, group, group, group, group, group, group, group, group, group, group)

	log.Info("flux sql=", querySql)
	result, err := queryApi.Query(ctx, querySql)
	if err != nil {
		log.Errorf("flux err=", err)
		return nil, err
	}
	resultMap := make(map[string]map[string]interface{})
	for result.Next() {
		key := ""
		if v, ok := result.Record().Values()[group]; ok {
			key = v.(string)
		}
		resultMap[key] = result.Record().Values()
	}
	return resultMap, nil

}

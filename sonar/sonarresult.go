package sonar

import (
	//"encoding/json"
	//"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	//"reflect"
	"strings"
)

func GetSonarResult(jobname string) map[string]string {
	result := make(map[string]string)

	client := &http.Client{}
	base_url := "http://10.132.47.15:9001/api/measures/component?componentKey=jobname_temp&metricKeys=metric_temp"
	new_url := strings.Replace(base_url, "jobname_temp", jobname, 1)

	metric := [...]string{"ncloc", "ncloc_language_distribution", "files", "coverage",
		"public_documented_api_density", "public_undocumented_api", "comment_lines_density", "comment_lines",
		"duplicated_lines_density", "duplicated_blocks", "duplicated_files", "duplicated_lines",
		"violations", "blocker_violations", "critical_violations", "major_violations", "minor_violations",
		"sqale_index", "reliability_remediation_effort", "security_remediation_effort"}

	for m := range metric {
		req_url := strings.Replace(new_url, "metric_temp", metric[m], 1)
		request, err := http.NewRequest("GET", req_url, nil)
		if err != nil {
			panic(err)
		}
		response, _ := client.Do(request)
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			// handle error
		}

		//string(body)为字符串{"component":{"id":"AV2itJweVUXXYvpYUm3K","key":"python-pattern","name":"python-pattern","qualifier":"TRK","measures":[{"metric":"ncloc","value":"2740","periods":[{"index":1,"value":"0"}]}]}}
		js, err := simplejson.NewJson([]byte(string(body)))
		arr, _ := js.Get("component").Get("measures").Array()
		for _, u := range arr {
			newdi, _ := u.(map[string]interface{})
			value := newdi["value"].(string)
			//fmt.Println(reflect.TypeOf(value), value)
			//result[metric[m]] = string(value)
			result[metric[m]] = value
		}

		/*
			var dat map[string]interface{}
			if err := json.Unmarshal([]byte(string(body)), &dat); err == nil {
				//fmt.Println(dat)
			} else {
				fmt.Println(err)
			}

			//    inter是一个interface,不能直接取inter的值，比如inter["key"]
			inter := dat["component"]
			// 将interface转换为map
			mapper1 := inter.(map[string]interface{})
			   mapper1["measures"] 为interface  [map[metric:ncloc value:2740 periods:[map[index:1 value:0]]]]
			fmt.Println(mapper1["measures"])

			for i, u := range inter {
				fmt.Println(i, u)
			}

			slice1 := mapper1.(map[string]interface{})
			fmt.Println(slice1)

				var dat2 []map[string]interface{}
				if err := json.Unmarshal([]byte(mapper1), &dat2); err == nil {
					fmt.Println(dat)
				} else {
					fmt.Println(err)
				}
				fmt.Println(dat2[0])
			result[metric[m]] = dat["component"]
			result[metric[m]] = string(body)["component"]["measures"][0]["value"]
		*/
	}
	return result

	//for k, v := range result {
	//	fmt.Printf("       %s： %s\n", k, v)
}

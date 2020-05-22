package main

import (
	"testing"
)
var apiData ApiData
var config AppConfig
func TestParseYaml(t *testing.T) {
	config, err := ParseYaml("test.yml")
	if err != nil {
		t.Errorf("%v", err)
	}
	t.Log(config.TempPath, config.User)
	if config.TempPath != "sms.tmpl" || config.User != "daiwj" {
		t.Errorf("解析失败!")
	}

}

func TestGenApiData(t *testing.T) {
	var config = &AppConfig{
		User:       "daiwj",
		Pwd:        "12345",
		PhoneNumer: "1595717xxxx",
		TempPath:   "sms.tmpl",
		UserId:     "1000",
		Url: "http://127.0.0.1:8080/sms2",
	}
	var a1 = Alert{
		Status:       "ok",
		Labels:       map[string]string{"warnPhone": "123", "test1": "test1"},
		Annotations:  nil,
		StartsAt:     "2121313",
		EndsAt:       "12121",
		GeneratorURL: "test.com",
	}
	var a2 = Alert{
		Status:       "not ok",
		Labels:       map[string]string{"warnPhone": "ssss", "test2": "test2"},
		Annotations:  nil,
		StartsAt:     "11111",
		EndsAt:       "22222",
		GeneratorURL: "test33333.com",
	}
	var promData = PromData{
		Status:       "Fired",
		Alerts:       []Alert{a1, a2},
		CommonLabels: map[string]string{"warnPhone": "159"},
		GroupKey:     "sssss",
		ExternalURL:  "ssss",
	}

	apiData, err := GenApiData(config, promData)
	if err != nil {
		t.Errorf("%v", err)
	}
	if apiData.PhoneNumer != "159" {
		t.Errorf("解析失败")
	}
	//json_data, _ := json.Marshal(res)
	t.Log("Con 内容是:" + apiData.Con)
	//t.Log(string(json_data))
	res,err:=PostSms(config,apiData)
	if err!=nil{
		t.Errorf("%v",err)
	}
	t.Log(res.Body)
}

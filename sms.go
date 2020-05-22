package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"net/http"
	"net/url"
)

type PromData struct {
	//普罗米修斯 webhook 发送的格式
	//https://prometheus.io/docs/alerting/configuration/
	Version string `json:"version"`
	Alerts []Alert `json:"alerts"`
	Status string `json:"status"`
	Receiver string `json:"receiver"`
	GroupLabels map[string]string `json:"groupLabels"`
	CommonLabels map[string]string `json:commonLabels`
	ExternalURL string `json:"externalURL"`
	GroupKey string `json:"groupKey"`
}

type Alert struct {
	//报警报文的格式
	//https://prometheus.io/docs/alerting/configuration/
	Status string `json:"status"`
	Labels map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartsAt string `json:"startAt"`
	EndsAt string `json:"endsAt"`
	GeneratorURL string `json:"generatorURL"`
}
type Sms struct {
	//测试
	Version string `json:"version"`
	Alerts string `json:"alerts"`
}


type AppConfig struct {
	//短信接口配置文件,
	User string `json:"user"`
	Pwd string `json:"pwd"`
	PhoneNumer string `json:"phoneNumber"`
	Ext string `json:"ext"`
	Priority string `json:"priority"`
	UserId string `json:"userid"`
	Url string `json:"url"`
	PhoneNumberFromAlert string `json:"phoneNumberFromAlert"`
	TempPath string `yaml:"tempPath" json:"tempPath"` //模板文件路径
}

type ApiData struct {
	//发送给短信接口的数据
	User string `json:"user"`
	Pwd string `json:"pwd"`
	PhoneNumer string `json:"phoneNumber"`
	Ext string `json:"ext"`
	Priority string `json:"priority"`
	UserId string `json:"userid"`
	Con string `json:"con"`
}

func ParseYaml(path string)(*AppConfig,error){
	//解析yml配置文件,返回 appconfig
	config := &AppConfig{}
	if f,err:=os.Open(path);err!=nil{
		log.Fatal(err)
	}else {
		yaml.NewDecoder(f).Decode(config)
	}
	return config,nil
}

func GenApiData(config *AppConfig,promData PromData)(*ApiData, error){
	//使用模板生成发送给短信接口的内容
	apiData :=&ApiData{}
	if promData.CommonLabels["warnPhone"]!=""{
		//如果alert 中指定了 warnPhone 这个值,则把短信发送给这个号码
		apiData.PhoneNumer=promData.CommonLabels["warnPhone"]
	}else {
		apiData.PhoneNumer=config.PhoneNumer
	}
	apiData.Ext=config.Ext
	apiData.Priority=config.Priority
	apiData.Pwd=config.Pwd
	apiData.User=config.User
	apiData.UserId=config.UserId

	//解析短信模板
	t1,err:=template.ParseFiles(config.TempPath)
	if err!=nil{
		log.Fatal(err)
	}
	//Execute 返回 io.writer,通过 buffer 来接收
	buf :=new(bytes.Buffer)
	t1.Execute(buf,promData.Alerts)
	//var res  string
	////将文字内容返回
	//res =buf.String()
	//fmt.Println(res)
	//把 string 内容放到 apidata 中,作为数据发送给短信接口
	apiData.Con=buf.String()

	return apiData,nil
}

//想短信接口发送数据
func PostSms(config *AppConfig, apiData *ApiData) (resp *http.Response,err error) {
	res,err:= http.PostForm(config.Url,url.Values{
		"user":{apiData.User},
		"pwd":{apiData.Pwd},
		"phone":{apiData.PhoneNumer},
		"con":{apiData.Con},
		"ext":{apiData.Ext},
		"priority":{apiData.Priority},
		"userid":{apiData.UserId},
	})
	if err !=nil{
		log.Fatal(err)
	}
	//fmt.Println(resp.Body)
	defer res.Body.Close()
	//return resp.Body
	return res,nil
}
func main() {
	var configPath = flag.String("ConfigPath","config.yml","配置文件")
	var listenPort = flag.String("HttpPort",":8080","http 监听端口")
	flag.Usage()
	//启动前,先检查配置文件
	config,err :=ParseYaml(*configPath)
	if err!=nil{
		log.Fatal(err)
	}
	configData,err:=json.Marshal(config)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("配置文件如下:",string(configData))
	//启动前读取配置文件
	//=====================

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		//测试函数

		c.String(200, "Hi! I am ok!")
	})

	r.POST("/sms",func(c *gin.Context){
		//接收 alertmanger 发送的消息,并将 alert 内容发送到短信接口

		var promdata PromData
		//把收到的 alerts 转换成 promdata 类型
		if err := c.ShouldBindJSON(&promdata);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"Error":err.Error(),
			})
		}
		//通过转换函数生成发往短信 api 的数据
		apiData,err:=GenApiData(config,promdata)
		if err!=nil{
			log.Fatal(err)
		}

		// 往短信 API 发送数据
		PostSms(config,apiData)

		//把从 alertmanger 发过来的数据解析成 json 格式
		jsonData,_:=json.Marshal(promdata)
		//返回 json 格式的内容,接口调试的时候使用.生产环境可以改成直接发送ok
		c.JSON(http.StatusOK, string(jsonData))

	})
	r.POST("/sms2",func(c *gin.Context){
		var sms ApiData
		con := c.PostForm("con")
		sms.Con=con
		//fmt.Println(con)
		c.JSON(http.StatusOK,sms)
	})
	r.Run(*listenPort) // listen and serve on 0.0.0.0:8080
}

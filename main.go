package main

import (
	"cass_open_api_sdk_golang/signer"
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	type response struct {
		Response struct {
			NotifyType int    `json:"notifyType"`
			NotifyTime string `json:"notifyTime"`
			Charset    string `json:"charset"`
			SignType   string `json:"signType"`
			Content    string `json:"content"`
		} `json:"response"`
		Sign string `json:"sign"`
	}
	r := response{}
	err := json.Unmarshal(body, &r)
	if err != nil {
		fmt.Printf("err0: %s \n", err.Error())
	}
	fmt.Printf("request: %s \n", r)
	params := map[string]interface{}{}
	bytes, err := json.Marshal(r.Response)
	if err != nil {
		fmt.Printf("err1: %s \n", err.Error())
		return
	}
	err = json.Unmarshal(bytes, &params)
	if err != nil {
		fmt.Printf("err2: %s \n", err.Error())
		return
	}

	// 对参数进行排序
	sortedKeys := make([]string, 0)
	for k := range params {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	sortFilter := map[string]interface{}{}

	for _, key := range sortedKeys {
		v := params[key]
		// 去除空值
		if v == "" || v == "{}" || v == "[]" {
			continue
		}
		// 去除 content 中的空格
		if key == "content" {
			params[key] = strings.ReplaceAll(params[key].(string), " ", "")
		}
		sortFilter[key] = params[key]
	}

	fmt.Printf("sorted params: %+v \n", sortFilter)

	t, _ := json.Marshal(sortFilter)
	waitSignStr := string(t)
	waitSignStr = strings.ReplaceAll(waitSignStr, " ", "")
	// 兼容 PHP 的 \ 字符在 json_encode 中的转义
	waitSignStr = strings.ReplaceAll(waitSignStr, `\/`, `\\/`)
	fmt.Printf("replace space str: %s \n", waitSignStr)
	waitSignStr = url.QueryEscape(waitSignStr)
	fmt.Printf("url encode str: %s \n", waitSignStr)

	// 进行验签
	s, err := signer.New(os.Getenv("VZHUO_PRIVATE_KEY_STR"), os.Getenv("VZHUO_PUBLIC_KEY_STR"))
	if err != nil {
		fmt.Printf("err3: %s \n", err.Error())
		return
	}

	// 取得请求中的 sign 二进制字数据
	sign, err := base64.StdEncoding.DecodeString(r.Sign)

	if err != nil {
		fmt.Printf("err4: %s \n", err)
	}

	// (Debug 可删除) 测试一下自己签名的结果
	selfSign, err := s.Sign([]byte(waitSignStr), crypto.SHA256)

	fmt.Printf("\n%s\n%s\n", selfSign, sign)

	// 验签
	err = s.Verify([]byte(waitSignStr), sign, crypto.SHA256)
	if err != nil {
		fmt.Printf("err4: %s \n", err.Error())
	}

	fmt.Printf("Verify success \n")

	_, _ = fmt.Fprintf(w, "Hello,"+req.URL.Path[1:])
}

func main() {
	_ = godotenv.Load(".env")
	http.HandleFunc("/", HelloServer)
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

package Cass

import (
	"cass_open_api_sdk_golang/Singer"
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type Request struct {
	signer   Singer.Signer
	BizParam map[string]interface{}
	Sign     string `json:"sign"`
	Params   *params
}

// 创建 Request
func NewRequest(privateKey string, APPID string, format string, charset string, version string, signType string) (Request, error) {
	var err error
	var request Request
	request.Params = &params{
		Method:   "",
		APPID:    APPID,
		Format:   format,
		Charset:  charset,
		Datetime: "",
		Version:  version,
		SignType: signType,
		BizParam: "",
		Sign:     "",
	}
	request.signer, err = Singer.NewSign(privateKey)
	if err != nil {
		return request, err
	}
	return request, nil
}

type params struct {
	Method   string `json:"method"`
	APPID    string `json:"APPID"`
	Format   string `json:"format"`
	Charset  string `json:"charset"`
	Datetime string `json:"datetime"`
	Version  string `json:"version"`
	SignType string `json:"signType"`
	BizParam string `json:"bizParam"`
	Sign     string `json:"sign"`
}

func (params *params) toMap() (map[string]interface{}, error) {
	var waitSignParams = map[string]interface{}{}

	j, _ := json.Marshal(params)
	err := json.Unmarshal(j, &waitSignParams)
	if err != nil {
		return waitSignParams, err
	}
	return waitSignParams, nil
}

func (params params) BuildQuery() (string, error) {
	str := new(strings.Builder)
	waitBuildQueryParams, err := params.toMap()
	if err != nil {
		return str.String(), err
	}
	if len(waitBuildQueryParams) != 0 {
		for key, val := range waitBuildQueryParams {
			str.WriteString(fmt.Sprintf("%s=%s&", key, val))
		}
	}
	return strings.TrimRight(str.String(), "&"), nil
}

// 构建请求参数对象
func (request *Request) BuildParams() error {
	bizParamBytes, err := json.Marshal(request.BizParam)
	if err != nil {
		return err
	}
	request.Params.BizParam = string(bizParamBytes)
	request.Params.Datetime = time.Now().Format("2006-01-2 15:04:05")

	err = request.makeSign()
	if err != nil {
		return err
	}
	return nil
}

func (request *Request) makeSign() error {
	var err error
	// 将 bizParam 转为json, 其中的中文不要转为 unicode 编码, 保持中文字符
	waitSignParams, err := request.Params.toMap()
	if err != nil {
		return err
	}

	// 过滤 Request 中的空字符: '', null, '[]', '{}'
	for s, v := range waitSignParams {
		if v == "" || v == "{}" || v == "[]" || v == nil || v == "null" {
			delete(waitSignParams, s)
		}
	}

	fmt.Printf("wait sign params: %s \n", waitSignParams)

	// 将 key 按照升序排序

	// 将 request 转换为 json

	jsonBytes, err := json.Marshal(waitSignParams)
	if err != nil {
		return err
	}

	// 将 request json str 中的空格 (ASCII 码空格) 去掉

	jsonStr := strings.ReplaceAll(string(jsonBytes), " ", "")

	// 将 request json str 进行 urlencode 编码, 产生待签名字符串

	urlEncodeStr := url.QueryEscape(jsonStr)
	fmt.Printf("wait sign str: %s \n", urlEncodeStr)
	// 通过字符串生成签名
	signBytes, err := request.signer.Sign([]byte(urlEncodeStr), crypto.SHA256)
	sign := base64.StdEncoding.EncodeToString(signBytes)
	fmt.Printf("sign string: %s \n", sign)
	if err != nil {
		return err
	}
	waitSignParams["sign"] = sign
	request.Sign = sign
	request.Params.Sign = sign
	return nil
}

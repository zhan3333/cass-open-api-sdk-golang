package Cass_test

import (
	"cass_open_api_sdk_golang/Cass"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	godotenv.Load("../.env")
	m.Run()
}

func TestNewRequest(t *testing.T) {
	request, err := Cass.NewRequest(
		os.Getenv("PRIVATE_KEY_STR"),
		os.Getenv("APPID"),
		"JSON",
		"UTF-8",
		"1.0",
		"RSA2",
	)
	assert.Nil(t, err)
	request.Params.Method = "Vzhuo.BcBalance.Get"
	client := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       10 * time.Second,
	}
	err = request.BuildParams()
	assert.Nil(t, err)
	buildQuery, err := request.Params.BuildQuery()
	fmt.Println(request.Params.Datetime)
	fmt.Printf("build query: %s \n", buildQuery)
	assert.Nil(t, err)
	post, err := client.Post(os.Getenv("API_URL"), "application/html; charset=utf-8", strings.NewReader(buildQuery))
	//fmt.Printf("response: %s", post.Body)
	fmt.Printf("err: %v", err)
	assert.Nil(t, err)
	assert.Equal(t, 200, post.StatusCode)
	bites, err := ioutil.ReadAll(post.Body)
	assert.Nil(t, err)
	fmt.Printf("Response: %s", bites)
}

func TestOneBankPay(t *testing.T) {
	request, err := Cass.NewRequest(
		os.Getenv("PRIVATE_KEY_STR"),
		os.Getenv("APPID"),
		"JSON",
		"UTF-8",
		"1.0",
		"RSA2",
	)
	assert.Nil(t, err)
	request.Params.Method = "Vzhuo.OneBankRemit.Pay"
	request.BizParam = map[string]interface{}{
		"payChannelK":      "1",
		"payeeChannelType": "2",
		"orderData": [1]interface{}{
			map[string]interface{}{
				"orderSN":          uuid.New().String(),
				"receiptFANO":      "13517210601",
				"payeeAccount":     "詹光",
				"requestPayAmount": "0.01",
				"notifyUrl":        "http://www.baidu.com",
			},
		},
	}
	client := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       10 * time.Second,
	}
	err = request.BuildParams()
	assert.Nil(t, err)
	buildQuery, err := request.Params.BuildQuery()
	fmt.Println(request.Params.Datetime)
	fmt.Printf("build query: %s \n", buildQuery)
	assert.Nil(t, err)
	post, err := client.Post(os.Getenv("API_URL"), "application/html; charset=utf-8", strings.NewReader(buildQuery))
	//fmt.Printf("response: %s", post.Body)
	fmt.Printf("err: %v", err)
	assert.Nil(t, err)
	assert.Equal(t, 200, post.StatusCode)
	bites, err := ioutil.ReadAll(post.Body)
	assert.Nil(t, err)
	fmt.Printf("Response: %s \n", bites)
	response := map[string]interface{}{}
	err = json.Unmarshal(bites, &response)
	if err != nil {
		fmt.Printf("%v \n", err)
	}
	assert.Nil(t, err)
	fmt.Printf("response: %s \n", response)
}

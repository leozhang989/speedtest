package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"speedtest/models"
	"strconv"
	"time"
)

//apple order handler
func DealAppleOrder(LatestReceipt string) (map[string]string, error){
	isSandBoxRes, _ := models.GetSettingsBySettingKey("isSandBox")
	appPasswordRes, _ := models.GetSettingsBySettingKey("appPassword")
	bundleIdRes, _ := models.GetSettingsBySettingKey("bundleId")
	appleParams := `{"receipt-data":"` + LatestReceipt + `", "password":"` + appPasswordRes.SettingValue + `"}`
	appleVerifyHost := "https://buy.itunes.apple.com/verifyReceipt"
	if intSandbox,_ := strconv.Atoi(isSandBoxRes.SettingValue); intSandbox == 1 {
		appleVerifyHost = "https://sandbox.itunes.apple.com/verifyReceipt";
	}

	appleResponseData,err := HttpPostJson(appleParams, appleVerifyHost)
	if err != nil {
		return nil,errors.New("get result from apple error")
	}

	var appleResponse AppleResult
	json.Unmarshal(appleResponseData, &appleResponse)
	//fmt.Println(appleResponse)

	switch {
	case appleResponse.Status == 21007:
		appleResponseData,err = HttpPostJson(appleParams, "https://sandbox.itunes.apple.com/verifyReceipt")

	case appleResponse.Status == 21008:
		appleResponseData,err = HttpPostJson(appleParams, "https://buy.itunes.apple.com/verifyReceipt")

	case appleResponse.Status >= 21100 && appleResponse.Status <= 21199:
		appleResponseData,err = HttpPostJson(appleParams, appleVerifyHost)
	}

	products := []string{"com.speed.1month", "com.speed.1year"}
	if appleResponse.Status == 0 {
		in_product_flag := false
		for _,v := range products{
			if v == appleResponse.Receipt.In_app[0].Product_id{
				in_product_flag = true
			}
		}
		if len(appleResponse.Receipt.Bundle_id) != 0 && appleResponse.Receipt.Bundle_id == bundleIdRes.SettingValue && len(appleResponse.Receipt.In_app) != 0 && in_product_flag {
			receiptInfoCount := len(appleResponse.Latest_receipt_info)
			if receiptInfoCount != 0 {
				lastestOrder := appleResponse.Latest_receipt_info[receiptInfoCount - 1]
				originalTransactionId := lastestOrder.Original_transaction_id
				//userInfos, _ := models.GetUsersByOtid(originalTransactionId)
				now := time.Now().Unix()
				if string(now * 1000) > lastestOrder.Expires_date_ms {
					return nil, errors.New("续订已过期，请重新购买")
				}

				returnRes := map[string]string{"originalTransactionId":originalTransactionId, "latestReceipt":appleResponse.Latest_receipt, "expiresDateS":lastestOrder.Expires_date_ms[:len(lastestOrder.Expires_date_ms)-3]}
				return returnRes,nil
			}
		}
	}

	return nil, errors.New("无需处理")
}


//发起POST请求
func HttpPostJson(requestJsonString string, requestUrl string) ([]byte, error) {
	jsonByteStr :=[]byte(requestJsonString)
	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(jsonByteStr))
	if err != nil {
		return nil,errors.New("request error")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil,errors.New("request error")
	}
	defer resp.Body.Close()

	//statuscode := resp.StatusCode
	//header := resp.Header
	body,_ := ioutil.ReadAll(resp.Body)
	//response := map[string]string{"code":strconv.Itoa(statuscode),"body":string(body[:])}
	return body,nil
}
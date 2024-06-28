package slzx

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dot-xiaoyuan/srun-demo/common/slzx"
	"github.com/dot-xiaoyuan/srun-demo/floger"
	"github.com/dot-xiaoyuan/srun-demo/helper"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type RespData struct {
	PayFee       string `json:"payfee"`
	TransferList []struct {
		TransferItem string `json:"transferItem"`
	} `json:"transferList"`
	MerChantOrderId string `json:"merchantorderid"`
	ReturnMsg       string `json:"returnmsg"`
	Reserved        string `json:"reserved"`
	Postscript      string `json:"postscript"`
	ReturnCode      string `json:"returncode"`
	OrderFee        string `json:"orderfee"`
	ChannelSerialNo string `json:"channelserialno"`
	Status          string `json:"status"`
}

type Pay struct {
	Model PayModel
	Conn  *gorm.DB
}

type PayResponse struct {
	ReturnCode   string      `json:"returncode"`
	ReturnMsg    string      `json:"returnmsg"`
	ReturnType   string      `json:"returntype"`
	SLIBusIId    string      `json:"slibusiid"`
	ResponseTime string      `json:"responstretime"`
	respData     interface{} `json:"respdata"`
	data         string      `json:"data"`
	sign         string      `json:"sign"`
	ResData      Data        `json:"resdata"`
}

type Data struct {
	CodeUrl     string `json:"code_url"`
	Status      string `json:"status"`
	CodeType    string `json:"code_type"`
	PathOrderId string `json:"pathorderid"`
}

func (p *Pay) UnifiedOrder() (res interface{}, err error) {
	res, err = p.NativePayRequest()

	// TODO 4k system create order
	return
}

// NativePayRequest 扫码支付
func (p *Pay) NativePayRequest() (res *PayResponse, err error) {
	// 接口地址
	nativePayRequestUrl := fmt.Sprintf("%s%s", slzx.ServiceUrl, "sltf-outside/inter/nativePayRequest")
	// code type
	codeType := "PAYLINK"

	// 组装业务参数
	var transferList []map[string]interface{}
	transferItemMap := make(map[string]interface{})
	transferItemMap["transferMerId"] = slzx.TransferMerId
	transferItemMap["transferItem"] = slzx.TransferItem
	transferList = append(transferList, transferItemMap)

	money, _ := strconv.ParseFloat(p.Model.Money, 64)
	reqData := map[string]interface{}{
		"txtype":          "03",
		"transferList":    transferList,
		"useridType":      "0",
		"code_type":       codeType,
		"userid":          p.Model.Username,
		"username":        p.Model.Username,
		"device_info":     p.Model.Username,
		"merchantorderid": p.Model.OutTradeNo,
		"orderfee":        int(money * 100), // 订单金额(分)
		"fronturl":        slzx.CallbackUrl,
		"receiveurl":      slzx.NotifyUrl,
	}

	reqDataJson, _ := json.Marshal(reqData)

	floger.Debug5("des 加密前的数据reqDataJson:", string(reqDataJson))

	reqDataDes, _ := EncryptDES(slzx.Key, string(reqDataJson))

	floger.Debug5("des 加密后的数据reqDataDes:", reqDataDes)

	signContent := reqDataDes + "&key=" + slzx.Md5Key

	floger.Debug5("md5 签名前 signContent:", signContent)

	signMsg := strings.ToLower(helper.Md5(signContent))

	floger.Debug5("md5 签名后 signMsg:", signMsg)

	headMap := map[string]interface{}{
		"version":         "V1.0",
		"charset":         "1",
		"platid":          slzx.PlatId,
		"produid":         "PM2000",
		"channelserialno": p.Model.OutTradeNo,
		"channeltime":     time.Now().Format("20060102150405"),
		"reqdata":         reqDataDes,
		"sign":            signMsg,
	}

	sendJson, _ := json.Marshal(headMap)

	floger.Debug5("发送数据 sendJson:", string(sendJson))
	resp, err := doRequest(nativePayRequestUrl, http.MethodPost, sendJson)
	if err != nil {
		floger.Error("Failed to request slzx, err:", err)
		return res, err
	}
	defer resp.Body.Close()

	floger.Debug5(resp.StatusCode)
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		floger.Error("Failed to json_decode response:", err)
		return res, err
	}
	floger.Debug5("返回结果 resultStr:", res)
	// 解密
	if res.ReturnCode != "0000" {
		return res, errors.New(res.ReturnMsg)
	}
	res.data, _ = DecryptDES(slzx.Key, res.respData.(string))

	_ = json.Unmarshal([]byte(res.data), &res.ResData)
	if res.ResData.Status != "05" {
		return res, errors.New("发起失败")
	}
	return res, nil
}

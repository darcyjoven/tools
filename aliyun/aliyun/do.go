package aliyun

import (
	"fmt"
	"net"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/spf13/viper"
)

var (
	client *alidns.Client
)

func ipv6() (string, error) {
	conn, err := net.Dial("udp", "[2001:db8::1]:80")
	if err != nil {
		fmt.Println("Error", err)
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}

// 更新阿里云记录
func updateAliyun(ip, id string) (err error) {
	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"

	request.RecordId = id
	request.RR = viper.GetString("rr")
	request.Type = "AAAA"
	request.Value = ip

	_, err = client.UpdateDomainRecord(request)
	return err
}

// 获取当前IP值
func getAliyunIP() (ip, id string, err error) {
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"
	request.DomainName = viper.GetString("domain")
	request.RRKeyWord = viper.GetString("rr")
	request.Type = "AAAA"
	response, err := client.DescribeDomainRecords(request)

	if err != nil {
		return "", "", err
	}
	if len(response.DomainRecords.Record) == 0 {
		return "", "", fmt.Errorf("未查到记录")
	}
	return response.DomainRecords.Record[0].Value, response.DomainRecords.Record[0].RecordId, err
}

// 初始化客户端
func initClient() (err error) {
	// new client
	client, err = alidns.NewClientWithAccessKey(
		"cn-hangzhou",
		viper.GetString("accesskeyid"),
		viper.GetString("accesskeysecret"),
	)
	if err != nil {
		return err
	}
	return
}

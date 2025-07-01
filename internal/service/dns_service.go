package service

import (
	"github.com/alibabacloud-go/alidns-20150109/v2/client"
	env "github.com/alibabacloud-go/darabonba-env/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
)

// DNSService 提供 DNS 相关的服务
type DNSService struct {
	client *client.Client
}

// NewDNSService 创建新的 DNS 服务实例
func NewDNSService(accessKeyId, accessKeySecret, regionId *string) (*DNSService, error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		RegionId:        regionId,
	}

	dnsClient, err := client.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &DNSService{
		client: dnsClient,
	}, nil
}

// Run 运行 DNS 更新服务
func Run(args []*string) error {
	if len(args) < 8 {
		return tea.NewSDKError(map[string]interface{}{
			"message": "参数不足，需要提供：domainName, RR, recordType, value, recordId, groupName, groupId, regionId",
		})
	}

	regionId := args[7]
	domainName := args[0]
	RR := args[1]
	recordType := args[2]
	value := args[3]
	recordId := args[4]
	groupName := args[5]
	groupId := args[6]

	// 初始化服务
	service, err := NewDNSService(
		env.GetEnv(tea.String("ACCESS_KEY_ID")),
		env.GetEnv(tea.String("ACCESS_KEY_SECRET")),
		regionId,
	)
	if err != nil {
		return err
	}

	// 执行所有操作
	operations := []func() error{
		func() error { return service.DescribeDomains() },
		func() error { return service.AddDomain(domainName) },
		func() error { return service.DescribeDomainRecords(domainName) },
		func() error { return service.DescribeRecordLogs(domainName) },
		func() error { return service.DescribeDomainRecordByRecordId(recordId) },
		func() error { return service.DescribeDomainInfo(domainName) },
		func() error { return service.AddDomainRecord(domainName, RR, recordType, value) },
		func() error { return service.UpdateDomainRecord(recordId, RR, recordType, value) },
		func() error { return service.SetDomainRecordStatus(recordId, tea.String("ENABLE")) },
		func() error { return service.DeleteDomainRecord(recordId) },
		func() error { return service.DescribeDomainGroups() },
		func() error { return service.AddDomainGroup(groupName) },
		func() error { return service.UpdateDomainGroup(groupId, groupName) },
		func() error { return service.DeleteDomainGroup(groupId) },
	}

	for _, op := range operations {
		if err := op(); err != nil {
			return err
		}
	}

	return nil
}

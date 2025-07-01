package service

import (
	dns "github.com/alibabacloud-go/alidns-20150109/v2/client"
	console "github.com/alibabacloud-go/tea-console/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
)

// DescribeDomains 查询账户下域名
func (s *DNSService) DescribeDomains() error {
	req := &dns.DescribeDomainsRequest{}
	console.Log(tea.String("查询域名列表(json)↓"))

	resp, err := s.client.DescribeDomains(req)
	if err != nil {
		console.Log(tea.String(err.Error()))
		return err
	}

	console.Log(util.ToJSONString(tea.ToMap(resp)))
	return nil
}

// AddDomain 阿里云云解析添加域名
func (s *DNSService) AddDomain(domainName *string) error {
	req := &dns.AddDomainRequest{
		DomainName: domainName,
	}
	console.Log(tea.String("云解析添加域名(" + tea.StringValue(domainName) + ")的结果(json)↓"))

	resp, err := s.client.AddDomain(req)
	if err != nil {
		console.Log(tea.String(err.Error()))
		return err
	}

	console.Log(util.ToJSONString(tea.ToMap(resp)))
	return nil
}

// DescribeDomainRecords 查询域名解析记录
func (s *DNSService) DescribeDomainRecords(domainName *string) error {
	req := &dns.DescribeDomainRecordsRequest{
		DomainName: domainName,
	}
	console.Log(tea.String("查询域名(" + tea.StringValue(domainName) + ")的解析记录(json)↓"))

	resp, err := s.client.DescribeDomainRecords(req)
	if err != nil {
		console.Log(tea.String(err.Error()))
		return err
	}

	console.Log(util.ToJSONString(tea.ToMap(resp)))
	return nil
}

// DescribeRecordLogs 查询域名解析记录日志
func (s *DNSService) DescribeRecordLogs(domainName *string) error {
	req := &dns.DescribeRecordLogsRequest{
		DomainName: domainName,
	}
	console.Log(tea.String("查询域名(" + tea.StringValue(domainName) + ")的解析记录日志(json)↓"))

	resp, err := s.client.DescribeRecordLogs(req)
	if err != nil {
		console.Log(tea.String(err.Error()))
		return err
	}

	console.Log(util.ToJSONString(tea.ToMap(resp)))
	return nil
}

// DescribeDomainRecordByRecordId 查询域名解析记录信息
func (s *DNSService) DescribeDomainRecordByRecordId(recordId *string) error {
	req := &dns.DescribeDomainRecordInfoRequest{
		RecordId: recordId,
	}
	console.Log(tea.String("查询RecordId:" + tea.StringValue(recordId) + "的域名解析记录信息(json)↓"))

	resp, err := s.client.DescribeDomainRecordInfo(req)
	if err != nil {
		console.Log(tea.String(err.Error()))
		return err
	}

	console.Log(util.ToJSONString(tea.ToMap(resp)))
	return nil
}

// DescribeDomainInfo 查询域名信息
func (s *DNSService) DescribeDomainInfo(domainName *string) error {
	req := &dns.DescribeDomainInfoRequest{
		DomainName: domainName,
	}
	console.Log(tea.String("查询域名:" + tea.StringValue(domainName) + "的信息(json)↓"))

	resp, err := s.client.DescribeDomainInfo(req)
	if err != nil {
		console.Log(tea.String(err.Error()))
		return err
	}

	console.Log(util.ToJSONString(tea.ToMap(resp)))
	return nil
}

// AddDomainRecord 添加域名解析记录
func (s *DNSService) AddDomainRecord(domainName, RR, recordType, value *string) error {
	req := &dns.AddDomainRecordRequest{
		DomainName: domainName,
		RR:         RR,
		Type:       recordType,
		Value:      value,
	}

	resp, err := s.client.AddDomainRecord(req)
	if err != nil {
		console.Log(tea.String(err.Error()))
		return err
	}

	console.Log(tea.String("添加域名解析记录的结果(json)↓"))
	console.Log(util.ToJSONString(tea.ToMap(resp)))
	return nil
}

// UpdateDomainRecord 更新域名解析记录
func (s *DNSService) UpdateDomainRecord(recordId, RR, recordType, value *string) error {
	req := &dns.UpdateDomainRecordRequest{
		RecordId: recordId,
		RR:       RR,
		Type:     recordType,
		Value:    value,
	}

	resp, err := s.client.UpdateDomainRecord(req)
	if err != nil {
		console.Log(tea.String(err.Error()))
		return err
	}

	console.Log(tea.String("更新域名解析记录的结果(json)↓"))
	console.Log(util.ToJSONString(tea.ToMap(resp)))
	return nil
}

// SetDomainRecordStatus 设置域名解析状态
func (s *DNSService) SetDomainRecordStatus(recordId, status *string) error {
	req := &dns.SetDomainRecordStatusRequest{
		RecordId: recordId,
		Status:   status,
	}

	resp, err := s.client.SetDomainRecordStatus(req)
	if err != nil {
		console.Log(tea.String(err.Error()))
		return err
	}

	console.Log(tea.String("设置域名解析状态的结果(json)↓"))
	console.Log(util.ToJSONString(tea.ToMap(resp)))
	return nil
}

// DeleteDomainRecord 删除域名解析记录
func (s *DNSService) DeleteDomainRecord(recordId *string) error {
	req := &dns.DeleteDomainRecordRequest{
		RecordId: recordId,
	}

	resp, err := s.client.DeleteDomainRecord(req)
	if err != nil {
		console.Log(tea.String(err.Error()))
		return err
	}

	console.Log(tea.String("删除域名解析记录的结果(json)↓"))
	console.Log(util.ToJSONString(tea.ToMap(resp)))
	return nil
}

// DescribeDomainGroups 查询域名组
func (s *DNSService) DescribeDomainGroups() error {
	req := &dns.DescribeDomainGroupsRequest{}

	resp, err := s.client.DescribeDomainGroups(req)
	if err != nil {
		console.Log(tea.String(err.Error()))
		return err
	}

	console.Log(tea.String("查询域名组(json)↓"))
	console.Log(util.ToJSONString(tea.ToMap(resp)))
	return nil
}

// AddDomainGroup 添加域名组
func (s *DNSService) AddDomainGroup(groupName *string) error {
	req := &dns.AddDomainGroupRequest{
		GroupName: groupName,
	}

	resp, err := s.client.AddDomainGroup(req)
	if err != nil {
		console.Log(tea.String(err.Error()))
		return err
	}

	console.Log(tea.String("添加域名组的结果(json)↓"))
	console.Log(util.ToJSONString(tea.ToMap(resp)))
	return nil
}

// UpdateDomainGroup 更新域名组名称
func (s *DNSService) UpdateDomainGroup(groupId, groupName *string) error {
	req := &dns.UpdateDomainGroupRequest{
		GroupId:   groupId,
		GroupName: groupName,
	}

	resp, err := s.client.UpdateDomainGroup(req)
	if err != nil {
		console.Log(tea.String(err.Error()))
		return err
	}

	console.Log(tea.String("更新域名组的结果(json)↓"))
	console.Log(util.ToJSONString(tea.ToMap(resp)))
	return nil
}

// DeleteDomainGroup 删除域名组
func (s *DNSService) DeleteDomainGroup(groupId *string) error {
	req := &dns.DeleteDomainGroupRequest{
		GroupId: groupId,
	}

	resp, err := s.client.DeleteDomainGroup(req)
	if err != nil {
		console.Log(tea.String(err.Error()))
		return err
	}

	console.Log(tea.String("删除域名组的结果(json)↓"))
	console.Log(util.ToJSONString(tea.ToMap(resp)))
	return nil
}

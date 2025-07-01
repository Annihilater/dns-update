package service

import (
	"github.com/alibabacloud-go/alidns-20150109/v2/client"
	dns "github.com/alibabacloud-go/alidns-20150109/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
)

// Domain 域名信息
type Domain struct {
	DomainName string `json:"domain_name"`
	DomainId   string `json:"domain_id"`
	PunyCode   string `json:"puny_code"`
	AliDomain  bool   `json:"ali_domain"`
}

// DomainRecord DNS解析记录
type DomainRecord struct {
	RecordId string `json:"record_id"`
	RR       string `json:"rr"`
	Type     string `json:"type"`
	Value    string `json:"value"`
	Status   string `json:"status"`
	Locked   bool   `json:"locked"`
	Line     string `json:"line"`
	Priority int64  `json:"priority"`
	TTL      int64  `json:"ttl"`
}

// DNSService 提供 DNS 相关的服务
type DNSService struct {
	client *client.Client
}

// NewDNSService 创建新的 DNS 服务实例
func NewDNSService(accessKeyId, accessKeySecret *string) (*DNSService, error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		RegionId:        tea.String("cn-hangzhou"),
	}

	dnsClient, err := client.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &DNSService{
		client: dnsClient,
	}, nil
}

// ListDomains 获取所有域名列表
func (s *DNSService) ListDomains() ([]Domain, error) {
	req := &dns.DescribeDomainsRequest{}
	resp, err := s.client.DescribeDomains(req)
	if err != nil {
		return nil, err
	}

	domains := make([]Domain, 0)
	for _, d := range resp.Body.Domains.Domain {
		domains = append(domains, Domain{
			DomainName: tea.StringValue(d.DomainName),
			DomainId:   tea.StringValue(d.DomainId),
			PunyCode:   tea.StringValue(d.PunyCode),
			AliDomain:  tea.BoolValue(d.AliDomain),
		})
	}

	return domains, nil
}

// ListDomainRecords 获取指定域名的解析记录
func (s *DNSService) ListDomainRecords(domainName string) ([]DomainRecord, error) {
	req := &dns.DescribeDomainRecordsRequest{
		DomainName: tea.String(domainName),
	}
	resp, err := s.client.DescribeDomainRecords(req)
	if err != nil {
		return nil, err
	}

	records := make([]DomainRecord, 0)
	for _, r := range resp.Body.DomainRecords.Record {
		records = append(records, DomainRecord{
			RecordId: tea.StringValue(r.RecordId),
			RR:       tea.StringValue(r.RR),
			Type:     tea.StringValue(r.Type),
			Value:    tea.StringValue(r.Value),
			Status:   tea.StringValue(r.Status),
			Locked:   tea.BoolValue(r.Locked),
			Line:     tea.StringValue(r.Line),
			Priority: tea.Int64Value(r.Priority),
			TTL:      tea.Int64Value(r.TTL),
		})
	}

	return records, nil
}

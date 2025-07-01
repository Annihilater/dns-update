package service

import (
	"dns-update/pkg/logger"

	"github.com/alibabacloud-go/alidns-20150109/v2/client"
	dns "github.com/alibabacloud-go/alidns-20150109/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
	"go.uber.org/zap"
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

// ListDomainRecordsOptions 获取域名解析记录的选项
type ListDomainRecordsOptions struct {
	PageSize int64 // 每页记录数，默认5000
}

// DefaultListDomainRecordsOptions 默认的获取域名解析记录选项
var DefaultListDomainRecordsOptions = ListDomainRecordsOptions{
	PageSize: 20,
}

// SearchDomainRecordsOptions 查询域名解析记录的选项
type SearchDomainRecordsOptions struct {
	DomainName string // 域名
	RecordId   string // 解析记录ID
	RR         string // 主机记录
	Type       string // 记录类型
	Status     string // 状态
	PageSize   int64  // 每页记录数
}

// DefaultSearchDomainRecordsOptions 默认的查询域名解析记录选项
var DefaultSearchDomainRecordsOptions = SearchDomainRecordsOptions{
	PageSize: 20,
}

// DNSService 提供 DNS 相关的服务
type DNSService struct {
	client *client.Client
	log    *zap.Logger
}

// NewDNSService 创建新的 DNS 服务实例
func NewDNSService(accessKeyId, accessKeySecret *string, regionId string) (*DNSService, error) {
	log := logger.GetLogger()
	log.Info("初始化 DNS 服务",
		zap.String("accessKeyId", *accessKeyId),
		zap.String("regionId", regionId),
	)
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		RegionId:        tea.String(regionId),
	}

	dnsClient, err := client.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &DNSService{
		client: dnsClient,
		log:    logger.GetLogger(),
	}, nil
}

// ListDomains 获取所有域名列表
func (s *DNSService) ListDomains() ([]Domain, error) {
	s.log.Info("正在获取域名列表")

	req := &dns.DescribeDomainsRequest{}
	resp, err := s.client.DescribeDomains(req)
	if err != nil {
		s.log.Error("获取域名列表失败", zap.Error(err))
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

	s.log.Info("获取域名列表成功", zap.Int("count", len(domains)))
	return domains, nil
}

// ListDomainRecords 获取指定域名的解析记录
func (s *DNSService) ListDomainRecords(domainName string, opts *ListDomainRecordsOptions) ([]DomainRecord, error) {
	if opts == nil {
		opts = &DefaultListDomainRecordsOptions
	}

	// 确保PageSize有默认值
	if opts.PageSize == 0 {
		opts.PageSize = DefaultListDomainRecordsOptions.PageSize
	}

	s.log.Info("正在获取域名解析记录",
		zap.String("domain", domainName),
		zap.Int64("page_size", opts.PageSize),
	)

	var allRecords []DomainRecord
	pageNumber := int64(1)
	pageSize := opts.PageSize

	for {
		req := &dns.DescribeDomainRecordsRequest{
			DomainName: tea.String(domainName),
			PageSize:   tea.Int64(pageSize),
			PageNumber: tea.Int64(pageNumber),
		}

		resp, err := s.client.DescribeDomainRecords(req)
		if err != nil {
			s.log.Error("获取域名解析记录失败",
				zap.String("domain", domainName),
				zap.Int64("page", pageNumber),
				zap.Error(err),
			)
			return nil, err
		}

		// 处理当前页的记录
		for _, r := range resp.Body.DomainRecords.Record {
			allRecords = append(allRecords, DomainRecord{
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

		// 检查是否还有下一页
		totalCount := tea.Int64Value(resp.Body.TotalCount)
		totalPages := (totalCount + pageSize - 1) / pageSize

		s.log.Debug("获取域名解析记录分页信息",
			zap.String("domain", domainName),
			zap.Int64("current_page", pageNumber),
			zap.Int64("total_pages", totalPages),
			zap.Int("current_records", len(allRecords)),
			zap.Int64("total_records", totalCount),
		)

		if pageNumber >= totalPages {
			break
		}
		pageNumber++
	}

	s.log.Info("获取域名解析记录成功",
		zap.String("domain", domainName),
		zap.Int("count", len(allRecords)),
	)
	return allRecords, nil
}

// SearchDomainRecords 根据条件查询域名解析记录
func (s *DNSService) SearchDomainRecords(opts *SearchDomainRecordsOptions) ([]DomainRecord, error) {
	if opts == nil {
		opts = &DefaultSearchDomainRecordsOptions
	}

	// 确保PageSize有默认值
	if opts.PageSize == 0 {
		opts.PageSize = DefaultSearchDomainRecordsOptions.PageSize
	}

	s.log.Info("正在查询域名解析记录",
		zap.String("domain", opts.DomainName),
		zap.String("record_id", opts.RecordId),
		zap.String("rr", opts.RR),
		zap.String("type", opts.Type),
		zap.String("status", opts.Status),
		zap.Int64("page_size", opts.PageSize),
	)

	// 如果指定了RR（子域名），使用DescribeSubDomainRecords接口
	if opts.RR != "" {
		subDomain := opts.RR + "." + opts.DomainName
		req := &dns.DescribeSubDomainRecordsRequest{
			SubDomain: tea.String(subDomain),
			PageSize:  tea.Int64(opts.PageSize),
			Type:      tea.String(opts.Type),
		}

		resp, err := s.client.DescribeSubDomainRecords(req)
		if err != nil {
			s.log.Error("查询子域名解析记录失败",
				zap.String("sub_domain", subDomain),
				zap.Error(err),
			)
			return nil, err
		}

		var records []DomainRecord
		for _, r := range resp.Body.DomainRecords.Record {
			// 如果指定了RecordId，只返回匹配的记录
			if opts.RecordId != "" && tea.StringValue(r.RecordId) != opts.RecordId {
				continue
			}
			// 如果指定了Status，只返回匹配的记录
			if opts.Status != "" && tea.StringValue(r.Status) != opts.Status {
				continue
			}

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

		s.log.Info("查询子域名解析记录成功",
			zap.String("sub_domain", subDomain),
			zap.Int("count", len(records)),
		)
		return records, nil
	}

	// 如果没有指定RR，返回空记录
	s.log.Info("未指定子域名，返回空记录列表")
	return []DomainRecord{}, nil
}

// GetDomainRecordById 根据记录ID查询解析记录
func (s *DNSService) GetDomainRecordById(recordId string) (*DomainRecord, error) {
	s.log.Info("正在查询解析记录",
		zap.String("record_id", recordId),
	)

	req := &dns.DescribeDomainRecordInfoRequest{
		RecordId: tea.String(recordId),
	}

	resp, err := s.client.DescribeDomainRecordInfo(req)
	if err != nil {
		s.log.Error("查询解析记录失败",
			zap.String("record_id", recordId),
			zap.Error(err),
		)
		return nil, err
	}

	record := &DomainRecord{
		RecordId: tea.StringValue(resp.Body.RecordId),
		RR:       tea.StringValue(resp.Body.RR),
		Type:     tea.StringValue(resp.Body.Type),
		Value:    tea.StringValue(resp.Body.Value),
		Status:   tea.StringValue(resp.Body.Status),
		Locked:   tea.BoolValue(resp.Body.Locked),
		Line:     tea.StringValue(resp.Body.Line),
		Priority: tea.Int64Value(resp.Body.Priority),
		TTL:      tea.Int64Value(resp.Body.TTL),
	}

	s.log.Info("查询解析记录成功",
		zap.String("record_id", recordId),
		zap.String("rr", record.RR),
		zap.String("type", record.Type),
	)

	return record, nil
}

// GetDomainRecordsByStatus 获取指定域名下特定状态的所有解析记录
func (s *DNSService) GetDomainRecordsByStatus(domainName, status string, pageSize int64) ([]DomainRecord, error) {
	if pageSize == 0 {
		pageSize = DefaultSearchDomainRecordsOptions.PageSize
	}

	s.log.Info("正在获取域名解析记录",
		zap.String("domain", domainName),
		zap.String("status", status),
		zap.Int64("page_size", pageSize),
	)

	var allRecords []DomainRecord
	pageNumber := int64(1)

	for {
		req := &dns.DescribeDomainRecordsRequest{
			DomainName: tea.String(domainName),
			PageSize:   tea.Int64(pageSize),
			PageNumber: tea.Int64(pageNumber),
			Status:     tea.String(status),
		}

		resp, err := s.client.DescribeDomainRecords(req)
		if err != nil {
			s.log.Error("获取域名解析记录失败",
				zap.String("domain", domainName),
				zap.String("status", status),
				zap.Int64("page", pageNumber),
				zap.Error(err),
			)
			return nil, err
		}

		// 处理当前页的记录
		for _, r := range resp.Body.DomainRecords.Record {
			allRecords = append(allRecords, DomainRecord{
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

		// 检查是否还有下一页
		totalCount := tea.Int64Value(resp.Body.TotalCount)
		totalPages := (totalCount + pageSize - 1) / pageSize

		s.log.Debug("获取域名解析记录分页信息",
			zap.String("domain", domainName),
			zap.String("status", status),
			zap.Int64("current_page", pageNumber),
			zap.Int64("total_pages", totalPages),
			zap.Int("current_records", len(allRecords)),
			zap.Int64("total_records", totalCount),
		)

		if pageNumber >= totalPages {
			break
		}
		pageNumber++
	}

	s.log.Info("获取域名解析记录成功",
		zap.String("domain", domainName),
		zap.String("status", status),
		zap.Int("count", len(allRecords)),
	)
	return allRecords, nil
}

// GetDomainRecordsByType 获取指定域名下特定类型的所有解析记录
func (s *DNSService) GetDomainRecordsByType(domainName, recordType string, pageSize int64) ([]DomainRecord, error) {
	if pageSize == 0 {
		pageSize = DefaultSearchDomainRecordsOptions.PageSize
	}

	s.log.Info("正在获取域名解析记录",
		zap.String("domain", domainName),
		zap.String("type", recordType),
		zap.Int64("page_size", pageSize),
	)

	var allRecords []DomainRecord
	pageNumber := int64(1)

	for {
		req := &dns.DescribeDomainRecordsRequest{
			DomainName: tea.String(domainName),
			PageSize:   tea.Int64(pageSize),
			PageNumber: tea.Int64(pageNumber),
			Type:       tea.String(recordType),
		}

		resp, err := s.client.DescribeDomainRecords(req)
		if err != nil {
			s.log.Error("获取域名解析记录失败",
				zap.String("domain", domainName),
				zap.String("type", recordType),
				zap.Int64("page", pageNumber),
				zap.Error(err),
			)
			return nil, err
		}

		// 处理当前页的记录
		for _, r := range resp.Body.DomainRecords.Record {
			allRecords = append(allRecords, DomainRecord{
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

		// 检查是否还有下一页
		totalCount := tea.Int64Value(resp.Body.TotalCount)
		totalPages := (totalCount + pageSize - 1) / pageSize

		s.log.Debug("获取域名解析记录分页信息",
			zap.String("domain", domainName),
			zap.String("type", recordType),
			zap.Int64("current_page", pageNumber),
			zap.Int64("total_pages", totalPages),
			zap.Int("current_records", len(allRecords)),
			zap.Int64("total_records", totalCount),
		)

		if pageNumber >= totalPages {
			break
		}
		pageNumber++
	}

	s.log.Info("获取域名解析记录成功",
		zap.String("domain", domainName),
		zap.String("type", recordType),
		zap.Int("count", len(allRecords)),
	)
	return allRecords, nil
}

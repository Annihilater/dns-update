package handler

import (
	"net/http"
	"strconv"
	"strings"

	"dns-update/internal/service"

	"github.com/gin-gonic/gin"
)

// DNSHandler 处理DNS相关的HTTP请求
type DNSHandler struct {
	dnsService *service.DNSService
}

// NewDNSHandler 创建新的DNS处理器
func NewDNSHandler(dnsService *service.DNSService) *DNSHandler {
	return &DNSHandler{
		dnsService: dnsService,
	}
}

// ListDomains godoc
// @Summary      获取域名列表
// @Description  获取账户下所有的域名列表
// @Tags         domains
// @Accept       json
// @Produce      json
// @Success      200  {array}   service.Domain
// @Failure      500  {object}  string
// @Router       /domains [get]
func (h *DNSHandler) ListDomains(c *gin.Context) {
	domains, err := h.dnsService.ListDomains()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domains)
}

// ListDomainRecords godoc
// @Summary      获取域名解析记录
// @Description  获取指定域名的所有解析记录
// @Tags         records
// @Accept       json
// @Produce      json
// @Param        domain     path      string  true   "域名"
// @Param        page_size  query     integer false  "每页记录数，默认20"  minimum(1)  maximum(500)
// @Success      200    {array}   service.DomainRecord
// @Failure      500    {object}  string
// @Router       /domains/{domain}/records [get]
func (h *DNSHandler) ListDomainRecords(c *gin.Context) {
	domain := c.Param("domain")

	// 解析page_size参数
	opts := service.DefaultListDomainRecordsOptions
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "page_size必须是有效的整数"})
			return
		}
		if pageSize < 1 || pageSize > 500 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "page_size必须在1-500之间"})
			return
		}
		opts.PageSize = pageSize
	}

	records, err := h.dnsService.ListDomainRecords(domain, &opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

// SearchDomainRecords godoc
// @Summary      搜索域名解析记录
// @Description  根据条件搜索域名解析记录
// @Tags         records
// @Accept       json
// @Produce      json
// @Param        domain      path      string  true   "域名"
// @Param        record_id   query     string  false  "解析记录ID"
// @Param        rr          query     string  false  "主机记录"
// @Param        type        query     string  false  "记录类型"
// @Param        status      query     string  false  "状态(Enable/Disable)"
// @Param        page_size   query     integer false  "每页记录数，默认20"  minimum(1)  maximum(500)
// @Success      200    {array}   service.DomainRecord
// @Failure      400    {object}  string
// @Failure      500    {object}  string
// @Router       /domains/{domain}/records/search [get]
func (h *DNSHandler) SearchDomainRecords(c *gin.Context) {
	domain := c.Param("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "域名不能为空"})
		return
	}

	// 创建查询选项
	opts := service.SearchDomainRecordsOptions{
		DomainName: domain,
		RecordId:   c.Query("record_id"),
		RR:         c.Query("rr"),
		Type:       c.Query("type"),
		Status:     c.Query("status"),
	}

	// 解析page_size参数
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "page_size必须是有效的整数"})
			return
		}
		if pageSize < 1 || pageSize > 500 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "page_size必须在1-500之间"})
			return
		}
		opts.PageSize = pageSize
	}

	// 验证status参数
	if opts.Status != "" && opts.Status != "Enable" && opts.Status != "Disable" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status必须是Enable或Disable"})
		return
	}

	records, err := h.dnsService.SearchDomainRecords(&opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

// SearchDomainRecordsByRecordId godoc
// @Summary      根据记录ID查询解析记录
// @Description  根据记录ID查询域名解析记录
// @Tags         records
// @Accept       json
// @Produce      json
// @Param        domain      path      string  true   "域名"
// @Param        record_id   path      string  true   "解析记录ID"
// @Success      200    {object}   service.DomainRecord
// @Failure      400    {object}  string
// @Failure      404    {object}  string
// @Failure      500    {object}  string
// @Router       /domains/{domain}/records/id/{record_id} [get]
func (h *DNSHandler) SearchDomainRecordsByRecordId(c *gin.Context) {
	domain := c.Param("domain")
	recordId := c.Param("record_id")

	if domain == "" || recordId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "域名和记录ID不能为空"})
		return
	}

	record, err := h.dnsService.GetDomainRecordById(recordId)
	if err != nil {
		if strings.Contains(err.Error(), "DomainRecordNotFound") {
			c.JSON(http.StatusNotFound, gin.H{"error": "解析记录不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 验证记录是否属于指定域名
	recordDomain := record.RR + "." + domain
	if !strings.HasSuffix(recordDomain, domain) {
		c.JSON(http.StatusNotFound, gin.H{"error": "解析记录不属于指定域名"})
		return
	}

	c.JSON(http.StatusOK, record)
}

// SearchDomainRecordsByRR godoc
// @Summary      根据主机记录查询解析记录
// @Description  根据主机记录查询域名解析记录
// @Tags         records
// @Accept       json
// @Produce      json
// @Param        domain   path      string  true   "域名"
// @Param        rr       path      string  true   "主机记录"
// @Success      200    {array}   service.DomainRecord
// @Failure      400    {object}  string
// @Failure      500    {object}  string
// @Router       /domains/{domain}/records/rr/{rr} [get]
func (h *DNSHandler) SearchDomainRecordsByRR(c *gin.Context) {
	domain := c.Param("domain")
	rr := c.Param("rr")

	if domain == "" || rr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "域名和主机记录不能为空"})
		return
	}

	// 验证RR的格式
	if len(rr) > 255 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "主机记录长度不能超过255个字符"})
		return
	}

	opts := service.SearchDomainRecordsOptions{
		DomainName: domain,
		RR:         rr,
	}

	records, err := h.dnsService.SearchDomainRecords(&opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

// SearchDomainRecordsByType godoc
// @Summary      根据记录类型查询解析记录
// @Description  根据记录类型查询域名解析记录
// @Tags         records
// @Accept       json
// @Produce      json
// @Param        domain     path      string  true   "域名"
// @Param        type       path      string  true   "记录类型"
// @Param        page_size  query     integer false  "每页记录数，默认20"  minimum(1)  maximum(500)
// @Success      200    {array}   service.DomainRecord
// @Failure      400    {object}  string
// @Failure      500    {object}  string
// @Router       /domains/{domain}/records/type/{type} [get]
func (h *DNSHandler) SearchDomainRecordsByType(c *gin.Context) {
	domain := c.Param("domain")
	recordType := c.Param("type")

	if domain == "" || recordType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "域名和记录类型不能为空"})
		return
	}

	// 解析page_size参数
	var pageSize int64 = 0
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		var err error
		pageSize, err = strconv.ParseInt(pageSizeStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "page_size必须是有效的整数"})
			return
		}
		if pageSize < 1 || pageSize > 500 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "page_size必须在1-500之间"})
			return
		}
	}

	records, err := h.dnsService.GetDomainRecordsByType(domain, recordType, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

// SearchDomainRecordsByStatus godoc
// @Summary      根据状态查询解析记录
// @Description  查询指定域名下所有特定状态的解析记录
// @Tags         records
// @Accept       json
// @Produce      json
// @Param        domain     path      string  true   "域名"
// @Param        status     path      string  true   "状态(Enable/Disable)"
// @Param        page_size  query     integer false  "每页记录数，默认20"  minimum(1)  maximum(500)
// @Success      200    {array}   service.DomainRecord
// @Failure      400    {object}  string
// @Failure      500    {object}  string
// @Router       /domains/{domain}/records/status/{status} [get]
func (h *DNSHandler) SearchDomainRecordsByStatus(c *gin.Context) {
	domain := c.Param("domain")
	status := c.Param("status")

	if domain == "" || status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "域名和状态不能为空"})
		return
	}

	if status != "Enable" && status != "Disable" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "状态必须是Enable或Disable"})
		return
	}

	// 解析page_size参数
	var pageSize int64 = 0
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		var err error
		pageSize, err = strconv.ParseInt(pageSizeStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "page_size必须是有效的整数"})
			return
		}
		if pageSize < 1 || pageSize > 500 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "page_size必须在1-500之间"})
			return
		}
	}

	records, err := h.dnsService.GetDomainRecordsByStatus(domain, status, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

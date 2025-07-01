package handler

import (
	"net/http"
	"strconv"

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
// @Param        page_size  query     integer false  "每页记录数，默认5000"  minimum(1)  maximum(5000)
// @Success      200    {array}   service.DomainRecord
// @Failure      500    {object}  string
// @Router       /domains/{domain}/records [get]
func (h *DNSHandler) ListDomainRecords(c *gin.Context) {
	domain := c.Param("domain")

	// 解析page_size参数
	var opts service.ListDomainRecordsOptions
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "page_size必须是有效的整数"})
			return
		}
		if pageSize < 1 || pageSize > 5000 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "page_size必须在1-5000之间"})
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

package handler

import (
	"net/http"

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
// @Param        domain  path      string  true  "域名"
// @Success      200    {array}   service.DomainRecord
// @Failure      500    {object}  string
// @Router       /domains/{domain}/records [get]
func (h *DNSHandler) ListDomainRecords(c *gin.Context) {
	domain := c.Param("domain")
	records, err := h.dnsService.ListDomainRecords(domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

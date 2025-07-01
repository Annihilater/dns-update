package handler

import (
	"net/http"

	"dns-update/internal/service"

	"github.com/gin-gonic/gin"
)

type DNSHandler struct {
	dnsService *service.DNSService
}

func NewDNSHandler(dnsService *service.DNSService) *DNSHandler {
	return &DNSHandler{
		dnsService: dnsService,
	}
}

// ListDomains 获取所有域名列表
func (h *DNSHandler) ListDomains(c *gin.Context) {
	domains, err := h.dnsService.ListDomains()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"domains": domains,
	})
}

// ListDomainRecords 获取指定域名的解析记录
func (h *DNSHandler) ListDomainRecords(c *gin.Context) {
	domainName := c.Param("domain")
	records, err := h.dnsService.ListDomainRecords(domainName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"records": records,
	})
}

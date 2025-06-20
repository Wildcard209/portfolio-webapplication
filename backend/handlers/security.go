package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CSPReport struct {
	CSPReport struct {
		DocumentURI        string `json:"document-uri"`
		Referrer           string `json:"referrer"`
		ViolatedDirective  string `json:"violated-directive"`
		EffectiveDirective string `json:"effective-directive"`
		OriginalPolicy     string `json:"original-policy"`
		Disposition        string `json:"disposition"`
		BlockedURI         string `json:"blocked-uri"`
		LineNumber         int    `json:"line-number"`
		ColumnNumber       int    `json:"column-number"`
		SourceFile         string `json:"source-file"`
		StatusCode         int    `json:"status-code"`
		ScriptSample       string `json:"script-sample"`
	} `json:"csp-report"`
}

func CSPReportHandler(c *gin.Context) {
	var report CSPReport

	if err := c.ShouldBindJSON(&report); err != nil {
		log.Printf("CSP Report - Invalid JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	log.Printf("CSP Violation Report: Document=%s, Directive=%s, BlockedURI=%s, SourceFile=%s:%d:%d",
		report.CSPReport.DocumentURI,
		report.CSPReport.ViolatedDirective,
		report.CSPReport.BlockedURI,
		report.CSPReport.SourceFile,
		report.CSPReport.LineNumber,
		report.CSPReport.ColumnNumber,
	)

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "portfolio-webapplication",
	})
}

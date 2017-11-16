package v1

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/astronomerio/clickstream-ingestion-api/pkg/logging"
	v1types "github.com/astronomerio/clickstream-ingestion-api/pkg/types/v1"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func gzipToBatch(b []byte) (batch v1types.Batch, err error) {
	gzData, err := gzip.NewReader(bytes.NewBuffer(b))
	if err != nil {
		return
	}
	defer gzData.Close()
	d, err := ioutil.ReadAll(gzData)
	if err != nil {
		return
	}
	err = json.Unmarshal(d, &batch)
	return
}

func (h *RouteHandler) batchHandler(c *gin.Context) {
	logger := logging.GetLogger().WithFields(logrus.Fields{"package": "v1", "function": "batchHandler"})
	c.Set("profile", true)
	c.Set("type", "batch")
	c.Set("action", "batch")

	rd, err := c.GetRawData()

	if err != nil {
		c.Set("error", err.Error())
		c.Set("stage", "1")
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusOK, returnJSON)
		return
	}

	var batch v1types.Batch

	if c.GetHeader("Content-Encoding") == "gzip" {
		batch, err = gzipToBatch(rd)
		if err != nil {
			logger.Error("issue with gzip")
			logger.Error(err.Error())
			c.Set("error", err.Error())
			c.Set("stage", "2")
			c.AbortWithStatusJSON(http.StatusOK, returnJSON)
			return
		}
	} else {
		err = json.Unmarshal(rd, &batch)
		if err != nil {
			logger.Error(err.Error())
			c.Set("error", err.Error())
			c.Set("stage", "2")
			c.AbortWithStatusJSON(http.StatusOK, returnJSON)
			return
		}
	}

	md := v1types.GetRequestMetadata(c)
	for _, m := range batch.Messages {
		m.SentAt = batch.SentAt
		m.ApplyMetadata(md)
		m.SkewTimestamp()
		h.ingestionHandler.ProcessMessage(m.String(), m.PartitionKey())
	}

	c.AbortWithStatusJSON(http.StatusOK, returnJSON)
}

package elasticsearch

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic"
	"github.com/teris-io/shortid"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type DocumentRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

var serverPort string

func InitServer() {
	serverPort = os.Getenv("SERVER_PORT")
}

func StartServer() {
	r := gin.Default()
	r.POST("/documents", createDocumentsEndpoint)
	r.GET("/search", searchEndpoint)

	if err := r.Run(serverPort); err != nil {
		log.Fatalln("[ERR] can't connect to the server")
	}
}

func createDocumentsEndpoint(c *gin.Context) {
	var docs []DocumentRequest
	if err := c.BindJSON(&docs); err != nil {
		errorResponse(c, http.StatusBadRequest, "Malformed request body")
		return
	}

	bulk := elasticClient.
		Bulk().
		Index(elasticIndexName).
		Type(elasticTypeName)
	for _, d := range docs {
		doc := Document{
			ID:        shortid.MustGenerate(), // the generation of short, fully unique, non-sequential and by default URL friendly Ids
			Title:     d.Title,
			CreatedAt: time.Now().UTC(),
			Content:   d.Content,
		}
		bulk.Add(elastic.NewBulkIndexRequest().Id(doc.ID).Doc(doc))
	}
	if _, err := bulk.Do(c.Request.Context()); err != nil {
		log.Println(err)
		errorResponse(c, http.StatusInternalServerError, "Failed to create documents")
		return
	}
	c.Status(http.StatusOK)
}

type SearchResponse struct {
	Time      string
	Hits      string
	Documents []DocumentResponse
}

type DocumentResponse struct {
}

func searchEndpoint(c *gin.Context) {
	// Parse request
	query := c.Query("query")
	if query == "" {
		errorResponse(c, http.StatusBadRequest, "Query not specified")
		return
	}
	skip := 0
	take := 10
	if i, err := strconv.Atoi(c.Query("skip")); err == nil {
		skip = i
	}
	if i, err := strconv.Atoi(c.Query("take")); err == nil {
		take = i
	}

	esQuery := elastic.NewMultiMatchQuery(query, "title", "content").
		Fuzziness("2").
		MinimumShouldMatch("2")
	result, err := elasticClient.Search().
		Index(elasticIndexName).
		Query(esQuery).
		From(skip).Size(take).
		Do(c.Request.Context())
	if err != nil {
		log.Println(err)
		errorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	res := SearchResponse{
		Time: fmt.Sprintf("%d", result.TookInMillis),
		Hits: fmt.Sprintf("%d", result.Hits.TotalHits),
	}

	docs := make([]DocumentResponse, 0)
	for _, hit := range result.Hits.Hits {
		var doc DocumentResponse
		err := json.Unmarshal(*hit.Source, &doc)
		if err != nil {
			return
		}
		docs = append(docs, doc)
	}
	res.Documents = docs
	c.JSON(http.StatusOK, res)
}

func errorResponse(c *gin.Context, code int, err string) {
	c.JSON(code, gin.H{
		"error": err,
	})
}

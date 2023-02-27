package elasticsearch

import (
	"github.com/olivere/elastic"
	"log"
	"os"
	"time"
)

const (
	elasticIndexName = "documents"
	elasticTypeName  = "document"
)

var (
	elasticClient *elastic.Client
	port          string
)

type Document struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
}

func InitElastic() {
	port = os.Getenv("EL_PORT")
}

func Start() error {
	var err error
	for {
		elasticClient, err = elastic.NewClient(
			elastic.SetURL("http://elasticsearch:"+port),
			elastic.SetSniff(false),
		)
		if err != nil {
			log.Println(err)
			time.Sleep(3 * time.Second)
		} else {
			return nil
		}
	}
	return err
}

package services

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ambalabanov/scanner/dao"
	"github.com/ambalabanov/scanner/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func Parse(d models.Documents) {
	var wg sync.WaitGroup
	defer wg.Wait()
	for _, doc := range d {
		wg.Add(1)
		go parse(doc, &wg)
	}
}

func parse(d models.Document, wg *sync.WaitGroup) {
	defer wg.Done()
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	r, err := client.Get(d.URL)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	d.Method = r.Request.Method
	d.Scheme = r.Request.URL.Scheme
	d.Host = r.Request.Host
	d.Status = r.StatusCode
	d.Header = r.Header
	d.ID = primitive.NewObjectID()
	d.CreatedAt = time.Now()
	d.Body = body
	log.Println("Parse body")
	parseLinks(&d, ioutil.NopCloser(bytes.NewBuffer(body)))
	parseTitle(&d, ioutil.NopCloser(bytes.NewBuffer(body)))
	d.UpdatedAt = time.Now()
	if err := dao.InsertOne(d); err != nil {
		log.Println(err.Error())
	}
}

func parseLinks(d *models.Document, b io.Reader) {
	var links, forms []string
	tokenizer := html.NewTokenizer(b)
	for tokenType := tokenizer.Next(); tokenType != html.ErrorToken; {
		token := tokenizer.Token()
		if tokenType == html.StartTagToken {
			if token.DataAtom == atom.A || token.DataAtom == atom.Form {
				for _, attr := range token.Attr {
					switch attr.Key {
					case "href":
						links = append(links, attr.Val)
					case "action":
						forms = append(forms, attr.Val)
					}
				}
			}
		}
		tokenType = tokenizer.Next()
	}
	d.Links = links
	d.Forms = forms
}

func parseTitle(d *models.Document, b io.Reader) {
	tokenizer := html.NewTokenizer(b)
	for tokenType := tokenizer.Next(); tokenType != html.ErrorToken; {
		token := tokenizer.Token()
		if tokenType == html.StartTagToken {
			if token.DataAtom == atom.Title {
				tokenType = tokenizer.Next()
				if tokenType == html.TextToken {
					d.Title = tokenizer.Token().Data
					break
				}
			}
		}
		tokenType = tokenizer.Next()
	}
}

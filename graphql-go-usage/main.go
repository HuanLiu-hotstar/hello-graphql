package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
)

var fields = graphql.Fields{
	"hello": &graphql.Field{
		Type: graphql.String,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return "world data", nil
		},
	},
}

var rootQuery = graphql.ObjectConfig{
	Name:   "RootQuery",
	Fields: fields,
}
var schemaConfig = graphql.SchemaConfig{
	Query: graphql.NewObject(rootQuery),
}

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	log.Printf("query:%s", query)
	if len(result.Errors) > 0 {
		log.Printf("wrong result:%v", result.Errors)
	}
	return result
}

func main() {

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("err:%s", err)
	}

	http.HandleFunc("/graphql", func(w http.ResponseWriter, req *http.Request) {
		result := executeQuery(req.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)

	})
	log.Printf("listen :8081")
	usage()
	doclient(schema)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("err:%s", err)
	}
}
func doclient(schema graphql.Schema) {
	query := `{hello}`
	params := graphql.Params{
		Schema:        schema,
		RequestString: query,
	}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Printf("err:%v", r.Errors)
	}
	json.NewEncoder(os.Stdout).Encode(r)
}

func usage() {
	log.Printf(`
		curl "localhost:8081/graphql?query=\{hello\}"
	`)
}

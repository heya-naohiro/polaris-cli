package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/blugelabs/bluge"
	"github.com/ikawaha/blugeplugin/analysis/lang/ja"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var searchCmd = &cobra.Command{
	Use:   "query",
	Short: "query",
	Args:  cobra.MinimumNArgs(1),
	// メイン処理
	RunE: func(cmd *cobra.Command, args []string) error {
		// フラグを取得
		query := args[0]
		path := viper.GetString("indexdir")
		s := NewSearch(path)
		s.QueryPrint(query)
		return nil
	},
}

var indexCmd = &cobra.Command{
	Use:   "input",
	Short: "input",
	Args:  cobra.MinimumNArgs(1),
	// メイン処理
	RunE: func(cmd *cobra.Command, args []string) error {
		// フラグを取得
		doc_path := args[0]
		index_path := viper.GetString("indexdir")
		s := NewSearch(index_path)
		s.AddDocument(doc_path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(indexCmd)
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	viper.SetDefault("indexdir", filepath.Join(home, ".polaris_index"))
}

type seacher struct {
	Config bluge.Config
	Writer *bluge.Writer
}

func NewSearch(path string) *seacher {
	config := bluge.DefaultConfig(path)
	writer, err := bluge.OpenWriter(config)
	if err != nil {
		log.Fatalf("error opening writer: %v", err)
	}
	return &seacher{config, writer}
}

func (s *seacher) AddDocument(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("error open document: %v", err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("error read document: %v", err)
	}
	//fmt.Println(string(data))
	doc := bluge.NewDocument(path).
		AddField(bluge.NewTextField("body", string(data)).WithAnalyzer(ja.Analyzer()))

	err = s.Writer.Update(doc.ID(), doc)
	if err != nil {
		log.Fatalf("error updating document: %v", err)
	}
}

func (s *seacher) QueryPrint(q string) {
	fmt.Println("Query Start")
	reader, err := s.Writer.Reader()
	if err != nil {
		log.Fatalf("error getting index reader: %v", err)
	}
	defer reader.Close()

	query := bluge.NewMatchQuery(q).SetAnalyzer(ja.Analyzer()).SetField("body")
	request := bluge.NewTopNSearch(10, query).
		WithStandardAggregations()
	documentMatchIterator, err := reader.Search(context.Background(), request)
	if err != nil {
		log.Fatalf("error executing search: %v", err)
	}
	match, err := documentMatchIterator.Next()
	for err == nil && match != nil {
		err = match.VisitStoredFields(func(field string, value []byte) bool {
			if field == "_id" {
				fmt.Printf("match: %s\n", string(value))
			}
			return true
		})
		if err != nil {
			log.Fatalf("error loading stored fields: %v", err)
		}
		match, err = documentMatchIterator.Next()
	}
	if err != nil {
		log.Fatalf("error iterator document matches: %v", err)
	}
}

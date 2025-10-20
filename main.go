package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/philippgille/chromem-go"
)

func main() {
	// Define flags
	sourceFileToChunk := flag.String("file", "input.txt", "source to chunk and build embedding from")
	question := flag.String("question", "no question", "specify your question about the source")

	// Parse them
	flag.Parse()

	ctx := context.Background()

	db := chromem.NewDB()

	// For doing the embedding locally, use ollama, uncomment below lines
	embed := chromem.NewEmbeddingFuncOllama("nomic-embed-text", "http://localhost:11434/api")
	c, err := db.CreateCollection("knowledge-base", nil, embed)
	// Passing nil as embedding function leads to OpenAI being used and requires
	// "OPENAI_API_KEY" env var to be set.
	// c, err := db.CreateCollection("knowledge-base", nil, nil)
	if err != nil {
		panic(err)
	}

	chunks, err := splitText(*sourceFileToChunk)
	if err != nil {
		panic(err)
	}
	for i, cu := range chunks {
		err = c.AddDocument(ctx, chromem.Document{
			ID:      fmt.Sprintf("%d", i),
			Content: cu.Text,
		})
		if err != nil {
			panic(err)
		}
	}
	if err != nil {
		panic(err)
	}

	// Warm up Ollama, in case the model isn't loaded yet
	log.Println("Warming up")
	_ = askLLM(ctx, nil, "Hello!")

	start := time.Now()
	res, err := c.Query(ctx, *question, 5, nil, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ID: %v\nSimilarity: %v\nContent: %v\n", res[0].ID, res[0].Similarity, res[0].Content)
	fmt.Printf("ID: %v\nSimilarity: %v\nContent: %v\n", res[1].ID, res[1].Similarity, res[1].Content)
	log.Println("Search (incl query embedding) took", time.Since(start))
	// Print the retrieved documents and their similarity to the question.
	contexts := []string{}
	for i, res := range res {
		// Cut off the prefix we added before adding the document (see comment above).
		// This is specific to the "nomic-embed-text" model.
		contexts = append(contexts, res.Content)
		log.Printf("Document %d (similarity: %f): \"%s\"\n", i+1, res.Similarity, res.Content)

	}

	// Now we can ask the LLM again, augmenting the question with the knowledge we retrieved.
	// In this example we just use both retrieved documents as context.
	log.Println("Asking LLM with augmented question...")
	reply := askLLM(ctx, contexts, *question)
	log.Printf("Reply after augmenting the question with knowledge:\n, %s", reply)
}

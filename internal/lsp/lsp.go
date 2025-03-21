package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/lazydiv/lazyLang-compiler/internal/lexer"
	"github.com/lazydiv/lazyLang-compiler/internal/parser"
	"github.com/sourcegraph/jsonrpc2"
)

// Server struct handles LSP requests.
type Server struct{}

// Handle processes JSON-RPC requests.
func (s *Server) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	switch req.Method {
	case "initialize":
		// Respond with capabilities.
		capabilities := map[string]interface{}{
			"textDocumentSync": 1,
		}
		response := map[string]interface{}{
			"capabilities": capabilities,
		}
		conn.Reply(ctx, req.ID, response)

	case "textDocument/didOpen":
		// Correctly extract TextDocument params.
		var params struct {
			TextDocument struct {
				URI  string `json:"uri"`
				Text string `json:"text"`
			} `json:"textDocument"`
		}
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			log.Printf("Error unmarshalling didOpen params: %v", err)
			return
		}

		// Run diagnostics
		diagnostics := runDiagnostics(params.TextDocument.Text)

		// Send diagnostics notification.
		conn.Notify(ctx, "textDocument/publishDiagnostics", map[string]interface{}{
			"uri":         params.TextDocument.URI,
			"diagnostics": diagnostics,
		})

	default:
		log.Printf("Unhandled method: %s", req.Method)
	}
}

// runDiagnostics checks for errors using the lexer and parser.
func runDiagnostics(source string) []map[string]interface{} {
	lex := lexer.NewLexer(source)
	p := parser.NewParser(lex)
	prog := p.ParseProgram()

	// If the program is empty or parsing failed, return a syntax error.
	if prog == nil || len(prog.Statements) == 0 {
		return []map[string]interface{}{
			{
				"range": map[string]interface{}{
					"start": map[string]int{"line": 0, "character": 0},
					"end":   map[string]int{"line": 0, "character": 1},
				},
				"severity": 1,
				"message":  "Syntax error",
			},
		}
	}
	return []map[string]interface{}{}
}

func main() {
	// Use stdin/stdout for LSP communication.
	stream := jsonrpc2.NewBufferedStream(os.Stdin, jsonrpc2.VSCodeObjectCodec{})
	conn := jsonrpc2.NewConn(context.Background(), stream, &Server{})

	fmt.Println("LazyLang LSP server started.")

	// Instead of `conn.Wait()` or `conn.Run()`, manually listen for requests
	<-conn.DisconnectNotify()
}


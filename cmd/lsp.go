package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"yaml-validator/internal/config"
	"yaml-validator/internal/validator"
	"yaml-validator/pkg"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
)

// runLSP запускает Language Server Protocol (режим stdio)
func runLSP() {
	ctx := context.Background()
	stream := jsonrpc2.NewStream(&stdioReadWriteCloser{stdin: os.Stdin, stdout: os.Stdout})
	conn := jsonrpc2.NewConn(stream)
	docs := make(map[string][]byte) // URI -> content
	baseCfg := loadConfig()

	handler := func(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
		switch req.Method() {
		case "initialize":
			var params protocol.InitializeParams
			if err := unmarshalParams(req.Params(), &params); err != nil {
				return reply(ctx, nil, err)
			}
			result := protocol.InitializeResult{
				Capabilities: protocol.ServerCapabilities{
					TextDocumentSync: protocol.TextDocumentSyncOptions{
						OpenClose: true,
						Change:    protocol.TextDocumentSyncKindFull,
					},
				},
				ServerInfo: &protocol.ServerInfo{
					Name:    "yaml-validator",
					Version: version,
				},
			}
			return reply(ctx, result, nil)

		case "initialized":
			return reply(ctx, nil, nil)

		case "textDocument/didOpen":
			var params protocol.DidOpenTextDocumentParams
			if err := unmarshalParams(req.Params(), &params); err != nil {
				return nil
			}
			uri := string(params.TextDocument.URI)
			docs[uri] = []byte(params.TextDocument.Text)
			sendDiagnostics(conn, uri, params.TextDocument.Text, uint64(params.TextDocument.Version), baseCfg, docs)

		case "textDocument/didChange":
			var params protocol.DidChangeTextDocumentParams
			if err := unmarshalParams(req.Params(), &params); err != nil {
				return nil
			}
			uri := string(params.TextDocument.URI)
			if len(params.ContentChanges) > 0 {
				ch := params.ContentChanges[0]
				docs[uri] = []byte(ch.Text)
			}
			content := string(docs[uri])
			sendDiagnostics(conn, uri, content, uint64(params.TextDocument.Version), baseCfg, docs)
		}
		return nil
	}

	conn.Go(ctx, handler)
	<-conn.Done()
	if err := conn.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "LSP error: %v\n", err)
		os.Exit(1)
	}
}

func unmarshalParams(raw json.RawMessage, v interface{}) error {
	if raw == nil || len(raw) == 0 {
		return nil
	}
	return json.Unmarshal(raw, v)
}

func sendDiagnostics(conn jsonrpc2.Conn, docURI, content string, version uint64, baseCfg *config.Config, _ map[string][]byte) {
	filePath := uriToPath(docURI)
	if !strings.HasSuffix(strings.ToLower(filePath), ".yaml") && !strings.HasSuffix(strings.ToLower(filePath), ".yml") {
		return
	}
	cfg := config.ConfigForFile(baseCfg, filePath)
	errors := validator.ValidateFromContent(docURI, []byte(content), cfg)
	diagnostics := errorsToDiagnostics(errors)
	conn.Notify(context.Background(), "textDocument/publishDiagnostics", protocol.PublishDiagnosticsParams{
		URI:         protocol.DocumentURI(docURI),
		Version:     uint32(version),
		Diagnostics: diagnostics,
	})
}

func errorsToDiagnostics(errors []pkg.Error) []protocol.Diagnostic {
	out := make([]protocol.Diagnostic, 0, len(errors))
	for _, e := range errors {
		sev := protocol.DiagnosticSeverityError
		if e.Severity == "warning" {
			sev = protocol.DiagnosticSeverityWarning
		}
		d := protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{Line: uint32(max(0, e.Line-1)), Character: uint32(max(0, e.Column-1))},
				End:   protocol.Position{Line: uint32(max(0, e.Line-1)), Character: uint32(max(0, e.Column))},
			},
			Severity: sev,
			Message:  e.Message,
			Source:   "yaml-validator",
			Code:     e.Type,
		}
		if e.Suggestion != "" {
			d.Message = e.Message + " (To fix: " + e.Suggestion + ")"
		}
		out = append(out, d)
	}
	return out
}

func uriToPath(u string) string {
	if strings.HasPrefix(u, "file:///") {
		return strings.TrimPrefix(u, "file:///")
	}
	if strings.HasPrefix(u, "file://") {
		return strings.TrimPrefix(u, "file://")
	}
	return u
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type stdioReadWriteCloser struct {
	stdin  *os.File
	stdout *os.File
}

func (s *stdioReadWriteCloser) Read(p []byte) (n int, err error)  { return s.stdin.Read(p) }
func (s *stdioReadWriteCloser) Write(p []byte) (n int, err error) { return s.stdout.Write(p) }
func (s *stdioReadWriteCloser) Close() error {
	_ = s.stdin.Close()
	_ = s.stdout.Close()
	return nil
}

var _ io.ReadWriteCloser = (*stdioReadWriteCloser)(nil)

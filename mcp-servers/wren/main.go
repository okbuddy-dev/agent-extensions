package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func wrenBin() string {
	if bin := os.Getenv("WREN_BIN"); bin != "" {
		return bin
	}
	return "wren"
}

func runWren(args ...string) (string, error) {
	cmd := exec.Command(wrenBin(), args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w\n%s", err, string(out))
	}
	return string(out), nil
}

func toolErr(err error) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}},
		IsError: true,
	}
}

func toolText(s string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{mcp.TextContent{Type: "text", Text: s}},
	}
}

func bodyLimitHandler(h http.Handler, limit int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, limit)
		h.ServeHTTP(w, r)
	})
}

var (
	addrFlag string
)

func init() {
	flag.StringVar(&addrFlag, "addr", "", "Listen address for HTTP transport (e.g., 127.0.0.1:8080 or :8080)")
}

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "-h" || arg == "--help" {
			fmt.Println("mcp-wren v0.2.0 — MCP server for wren task management")
			fmt.Println("Supports stdio and HTTP transports.")
			fmt.Println()
			fmt.Println("Usage:")
			fmt.Println("  ./mcp-wren           Run on stdio transport (default)")
			fmt.Println("  TRANSPORT=http ./mcp-wren  Run on HTTP transport")
			fmt.Println()
			fmt.Println("Flags:")
			fmt.Println("  --addr <address>  Listen address for HTTP transport")
			fmt.Println("                    (overrides MCP_WREN_ADDR env var, default :8080)")
			fmt.Println()
			fmt.Println("Environment:")
			fmt.Println("  WREN_BIN       Path to wren binary (default: wren)")
			fmt.Println("  TRANSPORT      Transport mode: stdio or http (default: stdio)")
			fmt.Println("  MCP_WREN_ADDR  HTTP listen address (default: :8080)")
			os.Exit(0)
		}
	}

	flag.Parse()

	s := server.NewMCPServer("wren", "0.2.0",
		server.WithToolCapabilities(true),
	)

	s.AddTool(
		mcp.NewTool("list_tasks",
			mcp.WithDescription("List active wren tasks. Optionally filter by keyword."),
			mcp.WithString("filter",
				mcp.Description("Substring to filter tasks by title (optional)"),
			),
		),
		func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := []string{"-l"}
			if f := req.GetString("filter", ""); f != "" {
				args = append(args, f)
			}
			out, err := runWren(args...)
			if err != nil {
				return toolErr(err), nil
			}
			return toolText(out), nil
		},
	)

	s.AddTool(
		mcp.NewTool("list_done_tasks",
			mcp.WithDescription("List completed wren tasks."),
		),
		func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			out, err := runWren("-d")
			if err != nil {
				return toolErr(err), nil
			}
			return toolText(out), nil
		},
	)

	s.AddTool(
		mcp.NewTool("create_task",
			mcp.WithDescription("Create a new wren task. The title becomes the filename."),
			mcp.WithString("title",
				mcp.Required(),
				mcp.Description("Task title"),
			),
		),
		func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			title := req.GetString("title", "")
			out, err := runWren(title)
			if err != nil {
				return toolErr(err), nil
			}
			return toolText(out), nil
		},
	)

	s.AddTool(
		mcp.NewTool("complete_task",
			mcp.WithDescription("Mark a wren task as done. Use the exact task title to avoid interactive disambiguation."),
			mcp.WithString("title",
				mcp.Required(),
				mcp.Description("Exact task title (filename)"),
			),
		),
		func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			title := req.GetString("title", "")
			out, err := runWren("-d", title)
			if err != nil {
				return toolErr(err), nil
			}
			return toolText(out), nil
		},
	)

	s.AddTool(
		mcp.NewTool("read_task",
			mcp.WithDescription("Read the body content of a wren task."),
			mcp.WithString("keyword",
				mcp.Required(),
				mcp.Description("Keyword to match against task title"),
			),
		),
		func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			keyword := req.GetString("keyword", "")
			out, err := runWren("-r", keyword)
			if err != nil {
				return toolErr(err), nil
			}
			return toolText(out), nil
		},
	)

	s.AddTool(
		mcp.NewTool("prepend_task",
			mcp.WithDescription("Prepend a string to a wren task's filename (title)."),
			mcp.WithString("prefix",
				mcp.Required(),
				mcp.Description("String to prepend to the task title"),
			),
			mcp.WithString("task",
				mcp.Required(),
				mcp.Description("Task title to prepend to"),
			),
		),
		func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			prefix := req.GetString("prefix", "")
			task := req.GetString("task", "")
			out, err := runWren("--prepend", prefix, task)
			if err != nil {
				return toolErr(err), nil
			}
			return toolText(out), nil
		},
	)

	s.AddTool(
		mcp.NewTool("get_random_task",
			mcp.WithDescription("Return one random active wren task."),
		),
		func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			out, err := runWren("-o")
			if err != nil {
				return toolErr(err), nil
			}
			return toolText(out), nil
		},
	)

	if os.Getenv("TRANSPORT") == "http" {
		addr := addrFlag
		if addr == "" {
			addr = os.Getenv("MCP_WREN_ADDR")
			if addr == "" {
				addr = ":8080"
			}
		}

		srv := server.NewStreamableHTTPServer(s, server.WithStateLess(true))
		httpSrv := &http.Server{
			Addr:              addr,
			Handler:           bodyLimitHandler(srv, 1<<20),
			ReadHeaderTimeout: 10e9,
			ReadTimeout:       30e9,
			IdleTimeout:       120e9,
		}
		fmt.Fprintf(os.Stderr, "mcp-wren v0.2.0 listening on %s/mcp\n", addr)
		if err := httpSrv.ListenAndServe(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

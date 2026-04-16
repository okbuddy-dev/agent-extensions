package main

import (
	"context"
	"fmt"
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

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Println("mcp-wren v0.1.0 — MCP server for wren task management")
		fmt.Println("Runs on MCP stdio transport. Point your MCP client at this binary.")
		fmt.Println()
		fmt.Println("Environment:")
		fmt.Println("  WREN_BIN  Path to wren binary (default: wren)")
		os.Exit(0)
	}

	s := server.NewMCPServer("wren", "0.1.0",
		server.WithToolCapabilities(true),
	)

	// list_tasks
	s.AddTool(
		mcp.NewTool("list_tasks",
			mcp.WithDescription("List active wren tasks. Optionally filter by keyword."),
			mcp.WithString("filter",
				mcp.Description("Substring to filter tasks by title (optional)"),
			),
		),
		func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := []string{"-l"}
			if f, ok := req.Params.Arguments["filter"].(string); ok && f != "" {
				args = append(args, f)
			}
			out, err := runWren(args...)
			if err != nil {
				return toolErr(err), nil
			}
			return toolText(out), nil
		},
	)

	// list_done_tasks
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

	// create_task
	s.AddTool(
		mcp.NewTool("create_task",
			mcp.WithDescription("Create a new wren task. The title becomes the filename."),
			mcp.WithString("title",
				mcp.Required(),
				mcp.Description("Task title"),
			),
		),
		func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			title, _ := req.Params.Arguments["title"].(string)
			out, err := runWren(title)
			if err != nil {
				return toolErr(err), nil
			}
			return toolText(out), nil
		},
	)

	// complete_task
	s.AddTool(
		mcp.NewTool("complete_task",
			mcp.WithDescription("Mark a wren task as done. Use the exact task title to avoid interactive disambiguation."),
			mcp.WithString("title",
				mcp.Required(),
				mcp.Description("Exact task title (filename)"),
			),
		),
		func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			title, _ := req.Params.Arguments["title"].(string)
			out, err := runWren("-d", title)
			if err != nil {
				return toolErr(err), nil
			}
			return toolText(out), nil
		},
	)

	// read_task
	s.AddTool(
		mcp.NewTool("read_task",
			mcp.WithDescription("Read the body content of a wren task."),
			mcp.WithString("keyword",
				mcp.Required(),
				mcp.Description("Keyword to match against task title"),
			),
		),
		func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			keyword, _ := req.Params.Arguments["keyword"].(string)
			out, err := runWren("-r", keyword)
			if err != nil {
				return toolErr(err), nil
			}
			return toolText(out), nil
		},
	)

	// prepend_task
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
			prefix, _ := req.Params.Arguments["prefix"].(string)
			task, _ := req.Params.Arguments["task"].(string)
			out, err := runWren("--prepend", prefix, task)
			if err != nil {
				return toolErr(err), nil
			}
			return toolText(out), nil
		},
	)

	// get_random_task
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

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

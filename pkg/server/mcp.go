package server

import (
	"encoding/json"

	// Verified Imports for v1.1
	"runiq/pkg/driver"
	"runiq/pkg/engine"
)

type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      *int            `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      *int        `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func HandleRequest(req MCPRequest) *MCPResponse {
	// 1. INITIALIZE (The Handshake)
	if req.Method == "initialize" {
		// New Architecture: This warms up the Global Browser Body in the background.
		// It creates the process so it's ready when the first tool is called.
		go engine.EnsureBrowser()

		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]any{
				"protocolVersion": "2024-11-05",
				"capabilities":    map[string]any{"tools": map[string]any{}},
				"serverInfo": map[string]string{
					"name":    "runiq",
					"version": "1.1.0", // BUMPED VERSION (Stable Engine)
				},
			},
		}
	}

	if req.Method == "notifications/initialized" {
		return nil
	}

	// 2. TOOL DEFINITIONS (The Menu)
	if req.Method == "tools/list" {
		return &MCPResponse{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{"tools": []map[string]any{
			{"name": "runiq_navigate", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{"url": map[string]string{"type": "string"}}}},
			{"name": "runiq_click", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{"selector": map[string]string{"type": "string"}}}},
			{"name": "runiq_type", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{"selector": map[string]string{"type": "string"}, "text": map[string]string{"type": "string"}}}},
			{"name": "runiq_screenshot", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{}}},
			{"name": "runiq_inspect_web", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{}}},
			{"name": "runiq_read_file", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{"path": map[string]string{"type": "string"}}}},
			{"name": "runiq_write_file", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{"path": map[string]string{"type": "string"}, "content": map[string]string{"type": "string"}}}},
			{"name": "runiq_launch_app", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{"name": map[string]string{"type": "string"}}}},
			{"name": "runiq_type_global", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{"text": map[string]string{"type": "string"}}}},
			{"name": "runiq_screenshot_desktop", "inputSchema": map[string]any{"type": "object", "properties": map[string]any{}}},
		}}}
	}

	// 3. TOOL EXECUTION (The Action)
	if req.Method == "tools/call" || req.Method == "call/tool" {
		var p struct {
			Name string            `json:"name"`
			Args map[string]string `json:"arguments"`
		}
		json.Unmarshal(req.Params, &p)

		var out, img string

		// --- SECURITY GUARD (The "Sudo" Layer) ---
		requiresApproval := false
		details := ""

		switch p.Name {
		case "runiq_write_file":
			requiresApproval = true
			details = "Write to: " + p.Args["path"]
		case "runiq_type_global":
			requiresApproval = true
			details = "Type keys: " + p.Args["text"]
		case "runiq_click":
			requiresApproval = true
			details = "Click selector: " + p.Args["selector"]
		case "runiq_launch_app":
			requiresApproval = true
			details = "Launch: " + p.Args["name"]
		}

		if requiresApproval {
			// This calls the Engine's Native Popup system
			allowed := engine.AskUser(p.Name, details)
			if !allowed {
				return &MCPResponse{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{"content": []map[string]string{{"type": "text", "text": "Error: User Denied Action"}}}}
			}
		}
		// ----------------------------------------

		// ROUTING
		switch p.Name {
		// Browser Tools (Now backed by the Global Body)
		case "runiq_navigate":
			out = driver.Navigate(p.Args["url"])
		case "runiq_inspect_web":
			out = driver.Inspect()
		case "runiq_click":
			out = driver.Click(p.Args["selector"])
		case "runiq_type":
			out = driver.Type(p.Args["selector"], p.Args["text"])
		case "runiq_screenshot":
			out, img = driver.Screenshot()

		// System Tools (These use your existing filesystem/desktop drivers)
		case "runiq_read_file":
			out = driver.ReadFile(p.Args["path"])
		case "runiq_write_file":
			out = driver.WriteFile(p.Args["path"], p.Args["content"])
		case "runiq_launch_app":
			out = driver.LaunchApp(p.Args["name"])
		case "runiq_type_global":
			out = driver.TypeGlobal(p.Args["text"])
		case "runiq_screenshot_desktop":
			out, img = driver.ScreenshotDesktop()
		}

		content := []map[string]string{{"type": "text", "text": out}}
		if img != "" {
			content = append(content, map[string]string{"type": "image", "data": img, "mimeType": "image/png"})
		}

		return &MCPResponse{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{"content": content}}
	}
	return nil
}

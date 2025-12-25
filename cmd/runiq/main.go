package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"runiq/pkg/engine"
	"runiq/pkg/server"
)

func main() {
	// 1. LIFECYCLE MANAGEMENT
	// Ensures the browser closes cleanly when Runiq exits (e.g. user quits the chat)
	defer engine.CloseBrowser()

	// 2. MAIN EVENT LOOP
	// Runiq sits here waiting for JSON-RPC commands from the AI (via Stdin)
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := scanner.Bytes()

		// Parse Incoming Request
		var req server.MCPRequest
		if err := json.Unmarshal(input, &req); err != nil {
			// If we get garbage text (logs/noise), ignore it to keep the pipe clean
			continue
		}

		// Execute Logic
		resp := server.HandleRequest(req)

		// Send Response
		if resp != nil {
			bytes, _ := json.Marshal(resp)
			fmt.Println(string(bytes))
		}
	}
}

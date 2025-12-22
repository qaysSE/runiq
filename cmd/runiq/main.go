package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	
	// Import our new modular packages
	"runiq/pkg/server"
)

func main() {
	// standard input loop (MCP Protocol)
	scanner := bufio.NewScanner(os.Stdin)
	
	// Increase buffer size for large JSON payloads (like screenshots)
	const maxCapacity = 10 * 1024 * 1024 // 10MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		var req server.MCPRequest
		if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
			continue
		}

		// Hand off to the Server Package
		if res := server.HandleRequest(req); res != nil {
			bytes, _ := json.Marshal(res)
			fmt.Println(string(bytes))
		}
	}
}

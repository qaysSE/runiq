package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	
	"runiq/pkg/engine" // Import Engine to close browser
	"runiq/pkg/server"
)

func main() {
	// 1. Ensure Browser Closes when Runiq exits
	defer engine.CloseBrowser()

	scanner := bufio.NewScanner(os.Stdin)
	
	// Increase buffer for images
	const maxCapacity = 10 * 1024 * 1024 
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		var req server.MCPRequest
		if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
			continue
		}

		if res := server.HandleRequest(req); res != nil {
			bytes, _ := json.Marshal(res)
			fmt.Println(string(bytes))
		}
	}
}

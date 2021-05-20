package announce

import (
	"fmt"
	"log"
	"net/http"
)

func handleResponse(resp *http.Response, err error) error {
	if err != nil {
		return fmt.Errorf("failed sending request: %w", err)
	}
	if resp.StatusCode/100 != 2 {
		bufSize := 100
		buf := make([]byte, bufSize)
		n, err := resp.Body.Read(buf)
		if err != nil {
			log.Printf("failed reading error response: %v\n", err)
		}
		if bufSize < n {
			log.Printf("huge error response (%d < len)\n", bufSize)
		}
		log.Printf(string(buf))
		return fmt.Errorf("error response: %v", resp.Status)
	}
	return nil
}

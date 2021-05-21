//	Copyright 2021 Fabian Bergström
//	
//	Licensed under the Apache License, Version 2.0 (the "License");
//	you may not use this file except in compliance with the License.
//	You may obtain a copy of the License at
//	
//			http://www.apache.org/licenses/LICENSE-2.0
//	
//	Unless required by applicable law or agreed to in writing, software
//	distributed under the License is distributed on an "AS IS" BASIS,
//	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	See the License for the specific language governing permissions and
//	limitations under the License.

package announce

import (
	"fmt"
	"log"
	"net/http"
)

// Announcer announces messages for the public
type Announcer interface {
	Announce(string) error
}

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

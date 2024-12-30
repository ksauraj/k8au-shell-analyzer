// internal/gemini/gemini.go
package gemini

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type WrappedResponse struct {
	Sections []Section `json:"sections"`
}

type Section struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Animation   []string `json:"animation"`
	Quotes      []string `json:"quotes,omitempty"`
}

var apiKey string

/*
// Make a .env file while compiling on your local machine with your GEMINI_API_KEY
func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file. Please ensure it exists with GEMINI_API_KEY")
	}

	// Get API key from environment
	apiKey = os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		panic("GEMINI_API_KEY not found in .env file")
	}
}

*/

const (
	geminiAPIURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent"
)

func GenerateWrapped(data string) (WrappedResponse, error) {
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"text": fmt.Sprintf(`Analyze the following shell data and generate a summary with insights, quotes, and animations in the following JSON format:

{
  "sections": [
    {
      "title": "Section Title",
      "description": "Section description.",
      "animation": ["RowAnimation1", "RowAnimation2", ...],
      "quotes": ["Quote1", "Quote2", ...]
    },
    ...
  ]
}

Shell data: %s`, data),
					},
				},
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return WrappedResponse{}, fmt.Errorf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", geminiAPIURL+"?key="+apiKey, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return WrappedResponse{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return WrappedResponse{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	rawResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return WrappedResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}

	// Log the raw response
	if err := logResponse(rawResponse); err != nil {
		return WrappedResponse{}, fmt.Errorf("failed to log response: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(rawResponse, &result); err != nil {
		return WrappedResponse{}, fmt.Errorf("failed to decode response: %v", err)
	}

	if candidates, ok := result["candidates"].([]interface{}); ok && len(candidates) > 0 {
		if firstCandidate, ok := candidates[0].(map[string]interface{}); ok {
			if content, ok := firstCandidate["content"].(map[string]interface{}); ok {
				if parts, ok := content["parts"].([]interface{}); ok && len(parts) > 0 {
					if firstPart, ok := parts[0].(map[string]interface{}); ok {
						if text, ok := firstPart["text"].(string); ok {
							// Log the extracted text
							if err := logResponse([]byte("Extracted text: " + text)); err != nil {
								return WrappedResponse{}, fmt.Errorf("failed to log extracted text: %v", err)
							}

							// Remove the ```json``` markers
							jsonText := strings.TrimPrefix(text, "```json\n")
							jsonText = strings.TrimSuffix(jsonText, "\n```")

							// Remove any remaining backticks
							jsonText = strings.ReplaceAll(jsonText, "`", "")

							// Remove the note at the end of the JSON text
							noteIndex := strings.Index(jsonText, "**Note:**")
							if noteIndex != -1 {
								jsonText = jsonText[:noteIndex]
							}

							// Log the final JSON text before parsing
							if err := logResponse([]byte("Final jsonText: " + jsonText)); err != nil {
								return WrappedResponse{}, fmt.Errorf("failed to log final jsonText: %v", err)
							}

							var wrappedResp WrappedResponse

							// Log the JSON text before parsing
							if err := logResponse([]byte("JSON text to be parsed: " + jsonText)); err != nil {
								return WrappedResponse{}, fmt.Errorf("failed to log JSON text: %v", err)
							}

							// Parse the JSON text
							if err := json.Unmarshal([]byte(jsonText), &wrappedResp); err != nil {
								// Log the error
								if logErr := logResponse([]byte(fmt.Sprintf("Failed to parse text as JSON: %v\nJSON text: %s", err, jsonText))); logErr != nil {
									return WrappedResponse{}, fmt.Errorf("failed to log JSON parsing error: %v", logErr)
								}
								return WrappedResponse{}, fmt.Errorf("failed to parse text as JSON: %v", err)
							}

							// Log the successfully parsed response
							if err := logResponse([]byte(fmt.Sprintf("Successfully parsed WrappedResponse: %v", wrappedResp))); err != nil {
								return WrappedResponse{}, fmt.Errorf("failed to log parsed response: %v", err)
							}

							return wrappedResp, nil
						}
					}
				}
			}
		}
	}

	// Log the invalid response format
	if err := logResponse([]byte("Invalid response format")); err != nil {
		return WrappedResponse{}, fmt.Errorf("failed to log invalid response format: %v", err)
	}

	return WrappedResponse{}, fmt.Errorf("invalid response format")
}

func logResponse(response []byte) error {
	// Define log file path
	logPath := "gemini_response.log"

	// Open the file in append mode or create it if it doesn't exist
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// Log the error to a separate error log file
		errorLogPath := "gemini_error.log"
		errorLogFile, err := os.OpenFile(errorLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open error log file: %v", err)
		}
		defer errorLogFile.Close()
		_, errWrite := errorLogFile.WriteString(fmt.Sprintf("Error writing to log file: %v\n", err))
		if errWrite != nil {
			return fmt.Errorf("failed to write error to log file: %v", errWrite)
		}
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	// Write the response to the log file
	_, err = file.Write(response)
	if err != nil {
		// Log the error to a separate error log file
		errorLogPath := "gemini_error.log"
		errorLogFile, err := os.OpenFile(errorLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open error log file: %v", err)
		}
		defer errorLogFile.Close()
		_, errWrite := errorLogFile.WriteString(fmt.Sprintf("Error writing to log file: %v\n", err))
		if errWrite != nil {
			return fmt.Errorf("failed to write error to log file: %v", errWrite)
		}
		return fmt.Errorf("failed to write to log file: %v", err)
	}

	return nil
}

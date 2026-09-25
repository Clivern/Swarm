// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

package swarm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func readOut(outDir string) (summary string, totalTokens int, err error) {
	f, err := os.Open(filepath.Join(outDir, "pi.jsonl"))
	if err != nil {
		return "", 0, fmt.Errorf("read pi.jsonl: %w", err)
	}
	defer f.Close()

	var assistant strings.Builder

	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)
	for sc.Scan() {
		line := bytes.TrimSpace(sc.Bytes())
		if len(line) == 0 {
			continue
		}

		switch {
		case bytes.Contains(line, []byte(`"type":"message_update"`)):
			var ev struct {
				Usage *struct {
					TotalTokens int `json:"totalTokens"`
				} `json:"usage"`
			}
			if json.Unmarshal(line, &ev) == nil && ev.Usage != nil {
				totalTokens = ev.Usage.TotalTokens
			}
		case bytes.Contains(line, []byte(`"type":"message_end"`)):
			var ev struct {
				Message json.RawMessage `json:"message"`
			}
			if json.Unmarshal(line, &ev) != nil {
				continue
			}
			if text := assistantText(ev.Message); text != "" {
				if assistant.Len() > 0 {
					assistant.WriteString("\n\n")
				}
				assistant.WriteString(text)
			}
		}
	}
	if err := sc.Err(); err != nil {
		return "", 0, fmt.Errorf("read pi jsonl: %w", err)
	}
	return strings.TrimSpace(assistant.String()), totalTokens, nil
}

func assistantText(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var msg struct {
		Role    string          `json:"role"`
		Content json.RawMessage `json:"content"`
	}
	if err := json.Unmarshal(raw, &msg); err != nil {
		return ""
	}
	if msg.Role != "" && msg.Role != "assistant" {
		return ""
	}

	var s string
	if err := json.Unmarshal(msg.Content, &s); err == nil {
		return strings.TrimSpace(s)
	}
	var blocks []struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(msg.Content, &blocks); err != nil {
		return ""
	}
	var b strings.Builder
	for _, block := range blocks {
		if block.Text == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		b.WriteString(block.Text)
	}
	return strings.TrimSpace(b.String())
}

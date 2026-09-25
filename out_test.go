// Copyright 2026 Swarm. All rights reserved.
// License can be found in the LICENSE file.

package swarm

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitReadOut(t *testing.T) {
	dir := t.TempDir()
	const jsonl = `{"type":"message_update","usage":{"totalTokens":105}}
{"type":"message_end","message":{"role":"assistant","content":[{"type":"text","text":"Done editing const.py."}]}}
`
	assert.NoError(t, os.WriteFile(filepath.Join(dir, "pi.jsonl"), []byte(jsonl), 0o644))

	summary, total, err := readOut(dir)
	assert.NoError(t, err)
	assert.Equal(t, "Done editing const.py.", summary)
	assert.Equal(t, 105, total)
}

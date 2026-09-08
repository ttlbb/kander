// Package rules embeds the published Kander markdown rule files.
//
// The English originals live in en/ and are the source of truth. Chinese translations live in
// cn/ as matching *.md files. Names() enumerates en/ only. Requesting cn falls back to the English
// original when the translation is missing; every other language reads en/ directly. The language
// an agent uses to talk to the user is the agent_language config value, which the rules
// themselves instruct the agent to honor; the rule language only selects which text is installed.
package rules

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
)

//go:embed en cn
var content embed.FS

const (
	LangCN = "cn"
	LangEN = "en"
)

// LangFor maps a kander interface language (cn/en/ja) to the rule language shipped for it.
// Only Chinese has a translation; everything else reads the English originals.
func LangFor(language string) string {
	if language == LangCN {
		return LangCN
	}
	return LangEN
}

// Names returns the markdown rule files under en/, sorted.
func Names() []string {
	entries, err := fs.ReadDir(content, LangEN)
	if err != nil {
		return nil
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)
	return names
}

// File returns the bytes for name in lang and reports the language actually served. When lang is
// cn and cn/name is missing, it returns the English original and reports en. Any other lang reads
// the English original without fallback.
func File(lang, name string) (data []byte, actual string, err error) {
	if err := validateName(name); err != nil {
		return nil, "", err
	}
	if lang == LangCN {
		data, err = content.ReadFile(path.Join(LangCN, name))
		if err == nil {
			return data, LangCN, nil
		}
	}
	data, err = content.ReadFile(path.Join(LangEN, name))
	if err != nil {
		return nil, "", fmt.Errorf("rules: missing %s: %w", name, err)
	}
	return data, LangEN, nil
}

// Hash returns the SHA-256 hex digest of File(lang, name) and the language actually served.
func Hash(lang, name string) (digest, actual string, err error) {
	data, actual, err := File(lang, name)
	if err != nil {
		return "", "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), actual, nil
}

func validateName(name string) error {
	if name == "" || name != path.Base(name) || strings.Contains(name, "\\") || strings.Contains(name, "..") {
		return fmt.Errorf("rules: invalid name %q", name)
	}
	if !strings.HasSuffix(name, ".md") {
		return fmt.Errorf("rules: not a markdown rule file: %q", name)
	}
	return nil
}

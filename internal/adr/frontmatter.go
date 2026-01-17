package adr

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

const frontmatterDelimiter = "---"

// ExtractFrontmatter extracts YAML frontmatter from markdown content.
// Returns the frontmatter string (without delimiters) and the remaining body.
func ExtractFrontmatter(content string) (frontmatter string, body string, err error) {
	reader := bufio.NewReader(strings.NewReader(content))

	// Read first line - must be "---"
	firstLine, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", "", fmt.Errorf("reading first line: %w", err)
	}
	firstLine = strings.TrimSpace(firstLine)
	if firstLine != frontmatterDelimiter {
		return "", content, nil // No frontmatter
	}

	// Read until closing "---"
	var fmBuf bytes.Buffer
	var bodyBuf bytes.Buffer
	inFrontmatter := true

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", "", fmt.Errorf("reading content: %w", err)
		}

		if inFrontmatter {
			trimmed := strings.TrimSpace(line)
			if trimmed == frontmatterDelimiter {
				inFrontmatter = false
			} else {
				fmBuf.WriteString(line)
			}
		} else {
			bodyBuf.WriteString(line)
		}

		if err == io.EOF {
			break
		}
	}

	if inFrontmatter {
		return "", "", fmt.Errorf("unclosed frontmatter: missing closing '---'")
	}

	return fmBuf.String(), bodyBuf.String(), nil
}

// ParseFrontmatter parses YAML frontmatter into a Frontmatter struct.
func ParseFrontmatter(yamlContent string) (*Frontmatter, error) {
	var fm Frontmatter
	if err := yaml.Unmarshal([]byte(yamlContent), &fm); err != nil {
		return nil, fmt.Errorf("parsing frontmatter YAML: %w", err)
	}
	return &fm, nil
}

// ParseADR parses a complete ADR from markdown content.
func ParseADR(content string, filename string, filePath string) (*ADR, error) {
	fmStr, body, err := ExtractFrontmatter(content)
	if err != nil {
		return nil, fmt.Errorf("extracting frontmatter: %w", err)
	}

	if fmStr == "" {
		return nil, fmt.Errorf("no frontmatter found in %s", filename)
	}

	fm, err := ParseFrontmatter(fmStr)
	if err != nil {
		return nil, fmt.Errorf("parsing frontmatter in %s: %w", filename, err)
	}

	return &ADR{
		Frontmatter: *fm,
		Body:        body,
		Filename:    filename,
		FilePath:    filePath,
	}, nil
}

// SerializeFrontmatter converts a Frontmatter struct to YAML string with delimiters.
func SerializeFrontmatter(fm *Frontmatter) (string, error) {
	data, err := yaml.Marshal(fm)
	if err != nil {
		return "", fmt.Errorf("marshaling frontmatter: %w", err)
	}
	return fmt.Sprintf("---\n%s---\n", string(data)), nil
}

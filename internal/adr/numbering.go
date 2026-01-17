package adr

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/sventorben/decider/internal/validate"
)

var adrFilenameRegex = regexp.MustCompile(`^(\d{4})-.*\.md$`)

// FindNextNumber scans the ADR directory and returns the next available ADR number.
func FindNextNumber(adrDir string) (int, error) {
	entries, err := os.ReadDir(adrDir)
	if err != nil {
		if os.IsNotExist(err) {
			return 1, nil
		}
		return 0, fmt.Errorf("reading directory %s: %w", adrDir, err)
	}

	maxNum := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		matches := adrFilenameRegex.FindStringSubmatch(entry.Name())
		if matches == nil {
			continue
		}
		num, err := strconv.Atoi(matches[1])
		if err != nil {
			continue
		}
		if num > maxNum {
			maxNum = num
		}
	}

	return maxNum + 1, nil
}

// ListADRFiles returns all ADR files in the directory sorted by number.
func ListADRFiles(adrDir string) ([]string, error) {
	entries, err := os.ReadDir(adrDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading directory %s: %w", adrDir, err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if adrFilenameRegex.MatchString(entry.Name()) {
			files = append(files, entry.Name())
		}
	}

	sort.Slice(files, func(i, j int) bool {
		numI := extractNumberFromFilename(files[i])
		numJ := extractNumberFromFilename(files[j])
		return numI < numJ
	})

	return files, nil
}

func extractNumberFromFilename(filename string) int {
	matches := adrFilenameRegex.FindStringSubmatch(filename)
	if matches == nil {
		return 0
	}
	num, _ := strconv.Atoi(matches[1])
	return num
}

// ToKebabCase converts a string to kebab-case suitable for filenames.
func ToKebabCase(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces and underscores with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")

	// Remove non-alphanumeric characters except hyphens
	var result strings.Builder
	prevHyphen := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
			prevHyphen = false
		} else if r == '-' && !prevHyphen {
			result.WriteRune(r)
			prevHyphen = true
		}
	}

	// Trim leading and trailing hyphens
	return strings.Trim(result.String(), "-")
}

// GenerateFilename creates an ADR filename from a number and title.
func GenerateFilename(number int, title string) string {
	kebab := ToKebabCase(title)
	return fmt.Sprintf("%04d-%s.md", number, kebab)
}

// ValidateFilename checks if a filename matches the expected pattern
// and if the number matches the ADR ID.
func ValidateFilename(filename string, adrID string) error {
	matches := adrFilenameRegex.FindStringSubmatch(filename)
	if matches == nil {
		return fmt.Errorf("filename %q does not match pattern NNNN-*.md", filename)
	}

	fileNum, _ := strconv.Atoi(matches[1])
	idNum, err := ExtractNumber(adrID)
	if err != nil {
		return fmt.Errorf("cannot extract number from ADR ID: %w", err)
	}

	if fileNum != idNum {
		return fmt.Errorf("filename number %04d does not match ADR ID %s", fileNum, adrID)
	}

	return nil
}

// LoadADR loads and parses an ADR from a file path.
func LoadADR(filePath string) (*ADR, error) {
	// Check file size before reading
	if err := validate.CheckFileSize(filePath); err != nil {
		return nil, fmt.Errorf("file size check failed for %s: %w", filePath, err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", filePath, err)
	}

	filename := filepath.Base(filePath)
	return ParseADR(string(content), filename, filePath)
}

// LoadAllADRs loads all ADRs from a directory.
func LoadAllADRs(adrDir string) ([]*ADR, error) {
	files, err := ListADRFiles(adrDir)
	if err != nil {
		return nil, err
	}

	var adrs []*ADR
	for _, file := range files {
		adr, err := LoadADR(filepath.Join(adrDir, file))
		if err != nil {
			return nil, fmt.Errorf("loading %s: %w", file, err)
		}
		adrs = append(adrs, adr)
	}

	return adrs, nil
}

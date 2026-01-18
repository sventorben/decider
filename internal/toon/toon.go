// Package toon provides Token-Oriented Object Notation encoding/decoding.
// TOON is a compact encoding of the JSON data model optimized for token efficiency.
package toon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// Marshal encodes a value to TOON format with deterministic output.
func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalIndent encodes a value to TOON format with indentation.
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	enc.SetIndent(prefix, indent)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshal decodes TOON data into a value.
func Unmarshal(data []byte, v interface{}) error {
	dec := NewDecoder(bytes.NewReader(data))
	return dec.Decode(v)
}

// Encoder writes TOON values to an output stream.
type Encoder struct {
	w      io.Writer
	indent string
	prefix string
	depth  int
}

// NewEncoder creates a new TOON encoder.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// SetIndent sets the indentation for encoding.
func (e *Encoder) SetIndent(prefix, indent string) {
	e.prefix = prefix
	e.indent = indent
}

// Encode writes the TOON encoding of v to the stream.
func (e *Encoder) Encode(v interface{}) error {
	if err := e.encode(v); err != nil {
		return err
	}
	_, err := e.w.Write([]byte("\n"))
	return err
}

func (e *Encoder) encode(v interface{}) error {
	switch val := v.(type) {
	case nil:
		_, err := e.w.Write([]byte("null"))
		return err
	case bool:
		if val {
			_, err := e.w.Write([]byte("true"))
			return err
		}
		_, err := e.w.Write([]byte("false"))
		return err
	case int:
		_, err := e.w.Write([]byte(strconv.Itoa(val)))
		return err
	case int64:
		_, err := e.w.Write([]byte(strconv.FormatInt(val, 10)))
		return err
	case float64:
		_, err := e.w.Write([]byte(strconv.FormatFloat(val, 'f', -1, 64)))
		return err
	case string:
		return e.encodeString(val)
	case []interface{}:
		return e.encodeArray(val)
	case map[string]interface{}:
		return e.encodeObject(val)
	default:
		// Convert via JSON for struct types
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		var generic interface{}
		if err := json.Unmarshal(data, &generic); err != nil {
			return err
		}
		return e.encode(generic)
	}
}

func (e *Encoder) encodeString(s string) error {
	// Use bare word for simple identifiers, quoted otherwise
	if e.isBareWord(s) {
		_, err := e.w.Write([]byte(s))
		return err
	}
	// Quote the string
	quoted := strconv.Quote(s)
	_, err := e.w.Write([]byte(quoted))
	return err
}

func (e *Encoder) isBareWord(s string) bool {
	if s == "" || s == "true" || s == "false" || s == "null" {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if !unicode.IsLetter(r) && r != '_' {
				return false
			}
		} else {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
				return false
			}
		}
	}
	return true
}

func (e *Encoder) encodeArray(arr []interface{}) error {
	if _, err := e.w.Write([]byte("[")); err != nil {
		return err
	}
	e.depth++
	for i, item := range arr {
		if i > 0 {
			if _, err := e.w.Write([]byte(" ")); err != nil {
				return err
			}
		}
		if e.indent != "" && len(arr) > 0 {
			if i == 0 {
				if _, err := e.w.Write([]byte("\n")); err != nil {
					return err
				}
			}
			if _, err := e.w.Write([]byte(e.prefix + strings.Repeat(e.indent, e.depth))); err != nil {
				return err
			}
		}
		if err := e.encode(item); err != nil {
			return err
		}
		if e.indent != "" && i < len(arr)-1 {
			if _, err := e.w.Write([]byte("\n")); err != nil {
				return err
			}
		}
	}
	e.depth--
	if e.indent != "" && len(arr) > 0 {
		if _, err := e.w.Write([]byte("\n" + e.prefix + strings.Repeat(e.indent, e.depth))); err != nil {
			return err
		}
	}
	_, err := e.w.Write([]byte("]"))
	return err
}

func (e *Encoder) encodeObject(obj map[string]interface{}) error {
	if _, err := e.w.Write([]byte("{")); err != nil {
		return err
	}

	// Sort keys for deterministic output
	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	e.depth++
	for i, key := range keys {
		if i > 0 {
			if _, err := e.w.Write([]byte(" ")); err != nil {
				return err
			}
		}
		if e.indent != "" && len(keys) > 0 {
			if i == 0 {
				if _, err := e.w.Write([]byte("\n")); err != nil {
					return err
				}
			}
			if _, err := e.w.Write([]byte(e.prefix + strings.Repeat(e.indent, e.depth))); err != nil {
				return err
			}
		}
		if err := e.encodeString(key); err != nil {
			return err
		}
		if _, err := e.w.Write([]byte(":")); err != nil {
			return err
		}
		if err := e.encode(obj[key]); err != nil {
			return err
		}
		if e.indent != "" && i < len(keys)-1 {
			if _, err := e.w.Write([]byte("\n")); err != nil {
				return err
			}
		}
	}
	e.depth--
	if e.indent != "" && len(keys) > 0 {
		if _, err := e.w.Write([]byte("\n" + e.prefix + strings.Repeat(e.indent, e.depth))); err != nil {
			return err
		}
	}
	_, err := e.w.Write([]byte("}"))
	return err
}

// Decoder reads TOON values from an input stream.
type Decoder struct {
	r    io.Reader
	data []byte
	pos  int
}

// NewDecoder creates a new TOON decoder.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Decode reads the next TOON value from the input.
func (d *Decoder) Decode(v interface{}) error {
	if d.data == nil {
		data, err := io.ReadAll(d.r)
		if err != nil {
			return err
		}
		d.data = data
	}

	val, err := d.decodeValue()
	if err != nil {
		return err
	}

	// Marshal to JSON and unmarshal into target
	jsonData, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, v)
}

func (d *Decoder) skipWhitespace() {
	for d.pos < len(d.data) && unicode.IsSpace(rune(d.data[d.pos])) {
		d.pos++
	}
}

func (d *Decoder) decodeValue() (interface{}, error) {
	d.skipWhitespace()
	if d.pos >= len(d.data) {
		return nil, io.EOF
	}

	switch d.data[d.pos] {
	case '{':
		return d.decodeObject()
	case '[':
		return d.decodeArray()
	case '"':
		return d.decodeQuotedString()
	default:
		return d.decodeAtom()
	}
}

func (d *Decoder) decodeObject() (map[string]interface{}, error) {
	d.pos++ // skip '{'
	obj := make(map[string]interface{})

	for {
		d.skipWhitespace()
		if d.pos >= len(d.data) {
			return nil, fmt.Errorf("unexpected end of input in object")
		}
		if d.data[d.pos] == '}' {
			d.pos++
			return obj, nil
		}

		// Read key
		key, err := d.decodeValue()
		if err != nil {
			return nil, err
		}
		keyStr, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("object key must be string, got %T", key)
		}

		d.skipWhitespace()
		if d.pos >= len(d.data) || d.data[d.pos] != ':' {
			return nil, fmt.Errorf("expected ':' after object key")
		}
		d.pos++ // skip ':'

		// Read value
		val, err := d.decodeValue()
		if err != nil {
			return nil, err
		}
		obj[keyStr] = val
	}
}

func (d *Decoder) decodeArray() ([]interface{}, error) {
	d.pos++ // skip '['
	var arr []interface{}

	for {
		d.skipWhitespace()
		if d.pos >= len(d.data) {
			return nil, fmt.Errorf("unexpected end of input in array")
		}
		if d.data[d.pos] == ']' {
			d.pos++
			return arr, nil
		}

		val, err := d.decodeValue()
		if err != nil {
			return nil, err
		}
		arr = append(arr, val)
	}
}

func (d *Decoder) decodeQuotedString() (string, error) {
	d.pos++ // skip opening '"'
	start := d.pos

	for d.pos < len(d.data) {
		if d.data[d.pos] == '\\' {
			d.pos += 2
			continue
		}
		if d.data[d.pos] == '"' {
			s := string(d.data[start:d.pos])
			d.pos++ // skip closing '"'
			// Unescape the string
			unquoted, err := strconv.Unquote(`"` + s + `"`)
			if err != nil {
				return s, nil
			}
			return unquoted, nil
		}
		d.pos++
	}
	return "", fmt.Errorf("unterminated string")
}

func (d *Decoder) decodeAtom() (interface{}, error) {
	start := d.pos
	for d.pos < len(d.data) {
		c := d.data[d.pos]
		if unicode.IsSpace(rune(c)) || c == '{' || c == '}' || c == '[' || c == ']' || c == ':' {
			break
		}
		d.pos++
	}

	atom := string(d.data[start:d.pos])
	if atom == "" {
		return nil, fmt.Errorf("empty atom")
	}

	// Check for literals
	switch atom {
	case "null":
		return nil, nil
	case "true":
		return true, nil
	case "false":
		return false, nil
	}

	// Try to parse as number
	if i, err := strconv.ParseInt(atom, 10, 64); err == nil {
		return i, nil
	}
	if f, err := strconv.ParseFloat(atom, 64); err == nil {
		return f, nil
	}

	// Treat as bare string
	return atom, nil
}

// ToJSON converts TOON data to JSON.
func ToJSON(toonData []byte) ([]byte, error) {
	var v interface{}
	if err := Unmarshal(toonData, &v); err != nil {
		return nil, err
	}
	return json.MarshalIndent(v, "", "  ")
}

// FromJSON converts JSON data to TOON.
func FromJSON(jsonData []byte) ([]byte, error) {
	var v interface{}
	if err := json.Unmarshal(jsonData, &v); err != nil {
		return nil, err
	}
	return MarshalIndent(v, "", "  ")
}

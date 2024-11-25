package TemplateParser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Our basic object types we can handle
const (
	OBJECT_TYPE_STRING = iota
	OBJECT_TYPE_INTEGER
	OBJECT_TYPE_BOOLEAN
)

// ObjectType
// represents a generic type that can hold multiple kinds
// of data including integers, strings, and booleans.
type ObjectType struct {
	ObjectTypeId     int
	ObjectValue      interface{}
	ObjectDescriptor string
}

// SetString
// sets the ObjectType instance to hold a string value.
func (obj *ObjectType) SetString(s string, desc string) {
	obj.ObjectTypeId = OBJECT_TYPE_STRING
	obj.ObjectDescriptor = desc
	obj.ObjectValue = s
}

// SetInteger
// sets the ObjectType instance to hold an integer value.
func (obj *ObjectType) SetInteger(i uint64, desc string) {
	obj.ObjectTypeId = OBJECT_TYPE_INTEGER
	obj.ObjectValue = i
	obj.ObjectDescriptor = desc
}

// SetBoolean
// sets the ObjectType instance to hold a boolean value.
func (obj *ObjectType) SetBoolean(b bool, desc string) {
	obj.ObjectTypeId = OBJECT_TYPE_BOOLEAN
	obj.ObjectValue = b
	obj.ObjectDescriptor = desc
}

// GetString
// retrieves the string value and descriptor if the ObjectType holds a string, otherwise returns an error message.
func (obj *ObjectType) GetString() (bool, string, string) {
	if obj.ObjectTypeId != OBJECT_TYPE_STRING {
		return false, "Mismatched object type", ""
	} else {
		return true, obj.ObjectValue.(string), obj.ObjectDescriptor
	}

}

// GetInteger
// returns a boolean indicating success, the integer value, and an error message if the object type is not integer.
func (obj *ObjectType) GetInteger() (bool, uint64, string) {
	if obj.ObjectTypeId != OBJECT_TYPE_INTEGER {
		return false, 0, "Mismatch object type"
	}
	return true, obj.ObjectValue.(uint64), ""
}

// GetBoolean
// retrieves the boolean value and an error message if the ObjectType is not a boolean. Returns a success flag, value, and error.
func (obj *ObjectType) GetBoolean() (bool, bool, string) {
	if obj.ObjectTypeId != OBJECT_TYPE_BOOLEAN {
		return false, false, "Mismatch object type"
	}
	return true, obj.ObjectValue.(bool), ""
}

// Constants that are tags for the objects we recognize.
// For everything you want to recognize add a constnat for it.
const (
	TokenIdentifier   = 0 // A textual identifier (not a quoted string) Must start with two alpha characters
	TokenQuotedString = 1 // Quoted string
	TokenUint64       = 2 // 64-bit unsigned integer
	TokenUint32       = 3 // 32-bit unsigned integer
	TokenUint16       = 4 // 16-bit unsigned integer
	TokenUint8        = 5 // 8 bit unsigned integer
	TokenRegister     = 6 // A register object "r"number
	TokenMacro        = 7 // A macro identifier (@identifier)

	// TokenUnknown represents an unknown or unrecognized token type in the tokenization process.
	TokenUnknown = 255
)

// Names for the above constants
// For each object, this is its name
var TokenNames = []string{
	"Identifier",
	"QuotedString",
	"Uint64",
	"Uint32",
	"Uint16",
	"Uint8",
	"Register",
	"Macro",
}

// Token
// Represents a lexical token with a type and value.
type Token struct {
	Type          int
	ValueReceived string
}

// TemplateObject
// is a structure that contains template specifications for parsing input data.
type TemplateObject struct {
	TemplateType  int
	TemplateValue ObjectType
	TemplateError string
}

// Tokenize
// Scans the input string and generates a slice of tokens based on predefined patterns.
func Tokenize(input string) []Token {
	patterns := []struct {
		regex     *regexp.Regexp
		tokenType int
	}{
		{regexp.MustCompile(`^"([^"]*)"`), TokenQuotedString},
		{regexp.MustCompile(`^@[a-zA-Z][a-zA-z0-9_]*`), TokenMacro},
		{regexp.MustCompile(`^[a-zA-Z][a-zA-z][a-zA-Z0-9_]*`), TokenIdentifier},
		{regexp.MustCompile(`^[0-9a-fA-F]{9,16}`), TokenUint64},
		{regexp.MustCompile(`^[0-9a-fA-F]{5,8}`), TokenUint32},
		{regexp.MustCompile(`^[0-9a-fA-F]{3,4}`), TokenUint16},
		{regexp.MustCompile(`^[0-9a-fA-F]{1,2}`), TokenUint8},
		{regexp.MustCompile(`^r[0-9a-fA-F]*`), TokenRegister},
	}

	tokens := []Token{}
	offset := 0
	length := len(input)

	for offset < length {
		remaining := input[offset:]
		found := false

		for _, pattern := range patterns {
			matches := pattern.regex.FindStringSubmatch(remaining)
			if len(matches) > 0 {
				tokens = append(tokens, Token{pattern.tokenType, matches[0]})
				offset += len(matches[0])
				found = true
				break
			}
		}

		if !found {
			tokens = append(tokens, Token{TokenUnknown, string(remaining[0])})
			offset++
		}
	}

	return tokens
}

// EatComments
// Removes comments from the input string by truncating text at the first occurrence of a semicolon.
func EatComments(txt string) string {
	pos := strings.Index(txt, ";")
	if pos > -1 {
		return txt[:pos]
	}
	return txt
}

// ParseLine
// parses a line of text and attempts to match tokens against a list of template objects.
func ParseLine(txt string, templateList []TemplateObject) ([]ObjectType, bool, string) {
	// Create a list of objects
	objList := make([]ObjectType, 0)
	input := EatComments(strings.ToLower(txt))
	tokens := Tokenize(input)
	// If we have no tokens, stop here
	if len(tokens) == 0 {
		return nil, false, "No tokens found"
	}
	// For each token, process it and load an object
	for _, token := range tokens {
		switch token.Type {
		case TokenIdentifier:
			objList = append(objList,
				ObjectType{TokenIdentifier, token.ValueReceived, ""})
		case TokenMacro:
			objList = append(objList, ObjectType{TokenMacro, token.ValueReceived, ""})
		case TokenQuotedString:
			objList = append(objList, ObjectType{TokenQuotedString, token.ValueReceived, ""})
		case TokenUint64:
			val, err := strconv.ParseUint(token.ValueReceived, 16, 64)
			if err != nil {
				objList = append(objList, ObjectType{TokenUint64, 0, "The value of the register is not a valid hex number"})
				return objList, false, "Invalid number"
			} else {
				objList = append(objList, ObjectType{TokenUint64, val, ""})
			}
		case TokenUint32:
			val, err := strconv.ParseUint(token.ValueReceived, 16, 64)
			if err != nil {
				objList = append(objList, ObjectType{TokenUint32, 0, "The value of the register is not a valid hex number"})
				return objList, false, "Invalid number"
			} else {
				objList = append(objList, ObjectType{TokenUint32, val, ""})
			}
		case TokenUint16:
			val, err := strconv.ParseUint(token.ValueReceived, 16, 64)
			if err != nil {
				objList = append(objList, ObjectType{TokenUint16, 0, "The value of the register is not a valid hex number"})
				return objList, false, "Invalid number"
			} else {
				objList = append(objList, ObjectType{TokenUint16, val, ""})
			}
		case TokenUint8:
			val, err := strconv.ParseUint(token.ValueReceived, 16, 64)
			if err != nil {
				objList = append(objList, ObjectType{TokenUint8, 0, "The value of the register is not a valid hex number"})
				return objList, false, "Invalid number"
			} else {
				objList = append(objList, ObjectType{TokenUint8, val, ""})
			}
		case TokenUnknown:
			continue
		case TokenRegister:
			val, err := strconv.ParseUint(token.ValueReceived[1:], 16, 64)
			if err != nil {
				objList = append(objList, ObjectType{TokenRegister, 0, "The value of the register is not a valid hex number"})
				return objList, false, "Invalid number"
			} else {
				objList = append(objList, ObjectType{TokenRegister, val, ""})
			}
		}
	}
	// If we find our objects and tokens don't match, let us know.
	// It means this parsing is completely wrong
	if len(objList) != len(templateList) {
		return nil, false, "Object list and template list length do not match"
	}
	for idx, _ := range objList {
		if objList[idx].ObjectTypeId != templateList[idx].TemplateType {
			ot := objList[idx].ObjectTypeId
			tt := templateList[idx].TemplateType
			return objList, false, fmt.Sprintf("Expected type (%d)%s but got type (%d)%s: %s",
				tt, TokenNames[tt], ot, TokenNames[ot],
				templateList[idx].TemplateError)
		}
	}
	return objList, true, ""
}

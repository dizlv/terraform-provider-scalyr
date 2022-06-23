package sdk

import (
	"context"
	"encoding/json"
	"fmt"
)

const ParserConfigurationFilePath = "/logParsers/"

// Rewrite rule searches a field for a specified regular expression, and replaces the first match with a different string.
// Each rule is applied in order
type Rewrite struct {
	// Input field
	Input string `json:"input"`

	// Output field
	Output string `json:"output"`

	// Matching pattern
	Match string `json:"match"`

	// Data to replace with
	Replace string `json:"replace"`

	// Should replace all occurrances
	ReplaceAll bool `json:"replace_all"`
}

type Format struct {
	Name     string
	Format   string     `json:"format"`
	Halt     bool       `json:"halt,omitempty"`
	Rewrites []*Rewrite `json:"rewrites,omitempty"`
	Repeat   bool       `json:"repeat,omitempty"`
	Discard  bool       `json:"discard,omitempty"`
	// Each format can have dynamically named key with some value - those are the attributes that will be attached to a each parsed log
}

type LineGrouper struct {
	Start           string `json:"start"`
	ContinueThrough string `json:"continue_through"`
	ContinuePast    string `json:"continue_past"`
	HaltBefore      string `json:"halt_before"`
	HaltWith        string `json:"halt_with"`
	MaxChars        int    `json:"max_chars"`
	MaxLines        int    `json:"max_lines"`
}

type LineGroupers []LineGrouper

type Attributes = map[string]string
type Patterns = map[string]string

type Formats []Format

type CreateParserInput struct {
	Name                   string
	TimeZone               string        `json:"timezone,omitempty"`
	IntermittentTimestamps bool          `json:"intermittent_timestamps,omitempty"`
	Formats                Formats       `json:"formats"`
	AliasTo                string        `json:"alias_to,omitempty"`
	Attributes             Attributes    `json:"attributes,omitempty"`
	LineGroupers           []LineGrouper `json:"line_groupers,omitempty"`
	Patterns               Patterns      `json:"patterns,omitempty"`
}

type CreateParserOutput struct {
	Name string `json:"name"`
}

func (scalyr *ScalyrConfig) CreateParser(ctx context.Context, input *CreateParserInput) (*CreateParserOutput, error) {
	path := formatPathFromName(input.Name)

	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	_, err = scalyr.PutFile(ctx, path, fmt.Sprintf("%s", data))
	if err != nil {
		return nil, err
	}

	return &CreateParserOutput{
		Name: input.Name,
	}, nil
}

type ReadParserInput struct {
	Name string
}

type ReadParserOutput struct{}

func (scalyr *ScalyrConfig) ReadParser(ctx context.Context, input *ReadParserInput) (*ReadParserOutput, error) {
	path := formatPathFromName(input.Name)

	// todo: this has to bee unmarshalled
	_, err := scalyr.GetFile(ctx, path)
	if err != nil {
		return nil, err
	}

	return &ReadParserOutput{}, nil
}

type UpdateParserInput struct {
	CreateParserInput
}

type UpdateParserOutput struct {
	Path string
}

func (scalyr *ScalyrConfig) UpdateParser(ctx context.Context, input *UpdateParserInput) (*UpdateParserOutput, error) {
	path := formatPathFromName(input.Name)

	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	// The output must be unmarshalled once again
	output, err := scalyr.PutFile(ctx, path, fmt.Sprintf("%s", data))
	if err != nil {
		return nil, err
	}

	return &UpdateParserOutput{
		Path: output.Path,
	}, nil
}

type DeleteParserInput struct {
	Name string
}

type DeleteParserOutput struct{}

func (scalyr *ScalyrConfig) DeleteParser(ctx context.Context, input *DeleteParserInput) (*DeleteParserOutput, error) {
	path := formatPathFromName(input.Name)

	err := scalyr.DeleteFile(ctx, path)
	if err != nil {
		return nil, err
	}

	return &DeleteParserOutput{}, nil
}

func formatPathFromName(name string) string {
	return fmt.Sprintf("/logParsers/%s", name)
}

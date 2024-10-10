// Code generated by "enumer --values --type=Worker --linecomment --output worker_string.go --json --yaml --sql"; DO NOT EDIT.

package federated

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

const _WorkerName = "mastodon"

var _WorkerIndex = [...]uint8{0, 8}

const _WorkerLowerName = "mastodon"

func (i Worker) String() string {
	i -= 1
	if i < 0 || i >= Worker(len(_WorkerIndex)-1) {
		return fmt.Sprintf("Worker(%d)", i+1)
	}
	return _WorkerName[_WorkerIndex[i]:_WorkerIndex[i+1]]
}

func (Worker) Values() []string {
	return WorkerStrings()
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _WorkerNoOp() {
	var x [1]struct{}
	_ = x[Mastodon-(1)]
}

var _WorkerValues = []Worker{Mastodon}

var _WorkerNameToValueMap = map[string]Worker{
	_WorkerName[0:8]:      Mastodon,
	_WorkerLowerName[0:8]: Mastodon,
}

var _WorkerNames = []string{
	_WorkerName[0:8],
}

// WorkerString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func WorkerString(s string) (Worker, error) {
	if val, ok := _WorkerNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _WorkerNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Worker values", s)
}

// WorkerValues returns all values of the enum
func WorkerValues() []Worker {
	return _WorkerValues
}

// WorkerStrings returns a slice of all String values of the enum
func WorkerStrings() []string {
	strs := make([]string, len(_WorkerNames))
	copy(strs, _WorkerNames)
	return strs
}

// IsAWorker returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Worker) IsAWorker() bool {
	for _, v := range _WorkerValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for Worker
func (i Worker) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Worker
func (i *Worker) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Worker should be a string, got %s", data)
	}

	var err error
	*i, err = WorkerString(s)
	return err
}

// MarshalYAML implements a YAML Marshaler for Worker
func (i Worker) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for Worker
func (i *Worker) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	var err error
	*i, err = WorkerString(s)
	return err
}

func (i Worker) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *Worker) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of Worker: %[1]T(%[1]v)", value)
	}

	val, err := WorkerString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}

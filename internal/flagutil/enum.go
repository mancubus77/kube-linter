package flagutil

import (
	"fmt"
	"github.com/spf13/pflag"

	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/set"
	"golang.stackrox.io/kube-linter/internal/utils"
)

// EnumValue allows setting a list of values.
type EnumValue struct {
	flagType      string
	allowedValues set.FrozenStringSet

	currentValue string
}

// String implements pflag.Value.
// It is the preferred method of retrieving the set value for the enum.
func (e *EnumValue) String() string {
	return e.currentValue
}

// Set implements pflag.Value.
func (e *EnumValue) Set(input string) error {
	if !e.allowedValues.Contains(input) {
		return errors.Errorf("%q is not a valid option (valid options are %v)", input, e.getAllowedValuesSorted())
	}
	e.currentValue = input
	return nil
}

// Type implements pflag.Value.
func (e *EnumValue) Type() string {
	return e.flagType
}

// Check that EnumValue implements pflag.Value interface.
var _ pflag.Value = (*EnumValue)(nil)

// Usage returns a string that can be used as help text for this flag.
// It will include the flag type and the list of allowed values.
func (e *EnumValue) Usage() string {
	return fmt.Sprintf("%s (allowed values: %v)", e.flagType, e.getAllowedValuesSorted())
}

func (e *EnumValue) getAllowedValuesSorted() []string {
	return e.allowedValues.AsSortedSlice(func(i, j string) bool {
		return i < j
	})
}

// NewEnumValueFactory returns a factory that can generate enum flag values with the given flag type
// and allowedValues. Concrete flag values are generated by calling the factory with a particular defaultValue.
func NewEnumValueFactory(flagType string, allowedValues []string) func(defaultValue string) *EnumValue {
	allowedValuesSet := set.NewFrozenStringSet(allowedValues...)
	return func(defaultValue string) *EnumValue {
		partialValue := &EnumValue{flagType: flagType, allowedValues: allowedValuesSet}
		utils.Must(partialValue.Set(defaultValue))
		return partialValue
	}
}

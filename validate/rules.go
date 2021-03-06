// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package validate

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/pkg/errors"
	"regexp"
	"time"
)

var (
	Required      = validation.Required
	NilOrNotEmpty = validation.NilOrNotEmpty
	Match         = validation.Match
	OptionalUuid  = []validation.Rule{validation.NilOrNotEmpty, IfNotNil(is.UUID)}
	ValidScope    = []validation.Rule{validation.Required, validation.In([]interface{}{
		"controlPlaneId",
		"deviceId",
		"deviceType",
		"deviceSubType",
		"providerId",
		"serialKey",
		"serviceId",
		"serviceType",
		"shardId",
		"siteId",
		"subscriptionId",
		"templateId",
		"tenantGroupId",
		"tenantId",
	}...)}

	IsDuration = RuleFunc(CheckDuration)

	RegExpEmail = regexp.MustCompile(`^[_A-Za-z0-9-\+]+(\.[_A-Za-z0-9-]+)*@[A-Za-z0-9-]+(\.[A-Za-z0-9]+)*(\.[A-Za-z]{2,})$`)
)

type Rule interface {
	Validate(value interface{}) error
}

type RuleFunc func(value interface{}) error

func (f RuleFunc) Validate(value interface{}) error {
	return f(value)
}

var Self = RuleFunc(func(value interface{}) error {
	if validatable, ok := value.(Validatable); ok {
		return Validate(validatable)
	}
	return nil
})

func Composite(rules ...validation.Rule) RuleFunc {
	return func(value interface{}) error {
		result := types.ErrorList{}
		for _, rule := range rules {
			result = append(result, rule.Validate(value))
		}
		return result.Filter()
	}
}

func IfNotNil(rules ...validation.Rule) RuleFunc {
	return func(value interface{}) error {
		if value != nil {
			return Composite(rules...).Validate(value)
		}
		return nil
	}
}

func Iff(truth bool, rules ...validation.Rule) RuleFunc {
	return func(value interface{}) error {
		if truth {
			return Composite(rules...).Validate(value)
		}
		return nil
	}
}

func countConditions(conditions ...bool) int {
	var count int
	for _, condition := range conditions {
		if condition {
			count++
		}
	}
	return count
}

func Exactly(n int, fail error, conditions ...bool) RuleFunc {
	return func(value interface{}) error {
		count := countConditions(conditions...)
		if count != n {
			return fail
		}
		return nil
	}
}

func OneOf(fail error, conditions ...bool) RuleFunc {
	return Exactly(1, fail, conditions...)
}

func AllOf(fail error, conditions ...bool) RuleFunc {
	return Exactly(len(conditions), fail, conditions...)
}

func AnyOf(fail error, conditions ...bool) RuleFunc {
	return func(value interface{}) error {
		count := countConditions(conditions...)
		if count == 0 {
			return fail
		}
		return nil
	}
}

func CheckDuration(value interface{}) error {
	valueString, ok := value.(string)
	if !ok {
		return errors.New("Duration is not a string")
	}

	_, err := time.ParseDuration(valueString)
	return err
}

func ExactlyOneOptional(err error, values ...interface{}) RuleFunc {
	return func(value interface{}) error {
		nonNilValues := 0
		for _, value := range values {
			if value != nil {
				nonNilValues++
			}
		}
		if nonNilValues != 1 {
			return err
		}
		return nil
	}
}

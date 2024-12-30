package event

import (
	"errors"
	"fmt"
	"log"
	"maps"
	"slices"
	"strings"

	"github.com/neilsf/sitemongoose/internal/pkg/alert"
	"github.com/tidwall/gjson"
)

const (
	TRIGGER_TYPE_STATUS_CODE   = "status_code"
	TRIGGER_TYPE_RESPONSE_TIME = "response_time"
	TRIGGER_TYPE_JSON_RULE     = "json_rule"
)

var validTriggerTypes = map[string]bool{
	TRIGGER_TYPE_STATUS_CODE:   true,
	TRIGGER_TYPE_RESPONSE_TIME: true,
	TRIGGER_TYPE_JSON_RULE:     true,
}

const (
	JSON_RULE_VALID        = "valid"
	JSON_RULE_EXISTS       = "exists"
	JSON_RULE_EQUALS       = "eq"
	JSON_RULE_NOT_EQUAL    = "ne"
	JSON_RULE_LESS_THAN    = "lt"
	JSON_RULE_GREATER_THAN = "gt"
	JSON_RULE_REGEXP       = "regexp"
)

var validJSONRuleConditions = map[string]bool{
	JSON_RULE_VALID:        true,
	JSON_RULE_EXISTS:       true,
	JSON_RULE_EQUALS:       true,
	JSON_RULE_NOT_EQUAL:    true,
	JSON_RULE_LESS_THAN:    true,
	JSON_RULE_GREATER_THAN: true,
	JSON_RULE_REGEXP:       true,
}

type JSONRule struct {
	JSONPath  string `yaml:"json_path"`
	Condition string `yaml:"condition"`
	Value     string `yaml:"value"`
}

// An event is a condition that triggers an alert. It evaluates a server response based on TriggerType.
// If the condition is met, the event triggers an alert.
// If an alert is already active and the condition is no longer met, the event resolves the alert.
// In every other case, the event does nothing.
type Event struct {
	TriggerType            string    `yaml:"evaluate"`
	ExpectedStatusCode     int       `yaml:"expected_status_code"`
	ExpectedResponseTimeMS int       `yaml:"expected_response_time_ms"`
	JSONRule               *JSONRule `yaml:"json_rule"`
	active                 bool      // Active means it's currently triggered
	Alerts                 []alert.Alert
	MonitorName            string
}

// CheckTrigger evaluates the response based on the TriggerType and dispatches an alert if the condition is met.
func (e *Event) CheckTrigger(statusCode int, responseTimeMS int, body []byte) {
	switch e.TriggerType {
	case TRIGGER_TYPE_STATUS_CODE:
		e.checkStatusCode(statusCode)
	case TRIGGER_TYPE_RESPONSE_TIME:
		e.checkResponseTime(responseTimeMS)
	case TRIGGER_TYPE_JSON_RULE:
		e.checkJSONRule(body)
	}
}

// Validate checks if the event is correctly configured.
func (e *Event) Validate() (bool, error) {
	if !validTriggerTypes[e.TriggerType] {
		return false, errors.New("evaluate must be one of: " + strings.Join(slices.Collect(maps.Keys(validTriggerTypes)), ", "))
	}
	switch e.TriggerType {
	case TRIGGER_TYPE_STATUS_CODE:
		if e.ExpectedStatusCode < 100 || e.ExpectedStatusCode > 599 {
			return false, errors.New("expected_status_code must be between 100 and 599")
		}
	case TRIGGER_TYPE_RESPONSE_TIME:
		if e.ExpectedResponseTimeMS <= 0 {
			return false, errors.New("expected_response_time_ms must be greater than 0")
		}
	case TRIGGER_TYPE_JSON_RULE:
		if e.JSONRule == nil {
			return false, errors.New("json_rule is required")
		}
		if e.JSONRule.JSONPath == "" {
			return false, errors.New("json_path is required")
		}
		if !validJSONRuleConditions[e.JSONRule.Condition] {
			return false, errors.New("condition must be one of " + strings.Join(slices.Collect(maps.Keys(validJSONRuleConditions)), ", "))
		}
		if e.JSONRule.Value == "" {
			return false, errors.New("value is required")
		}
	}
	return true, nil
}

func (e *Event) trigger() {
	e.active = true
	for _, a := range e.Alerts {
		go alert.GetAlerter(a).SendAlert()
	}
	log.Printf("Triggered event: %s/%s\n", e.MonitorName, e.TriggerType)
}

func (e *Event) resolve() {
	e.active = false
	for _, a := range e.Alerts {
		go alert.GetAlerter(a).SendResolution()
	}
	log.Printf("Resolved event: %s/%s\n", e.MonitorName, e.TriggerType)
}

func (e *Event) dispatch(condition bool) {
	if e.active && !condition {
		e.resolve()
	} else if !e.active && condition {
		e.trigger()
	}
}

func (e *Event) checkStatusCode(statusCode int) {
	e.dispatch(statusCode != e.ExpectedStatusCode)
}

func (e *Event) checkResponseTime(responseTimeMS int) {
	e.dispatch(responseTimeMS > e.ExpectedResponseTimeMS)
}

func (e *Event) checkJSONRule(bodyBytes []byte) {
	if !gjson.ValidBytes(bodyBytes) {
		if e.JSONRule.Condition == JSON_RULE_VALID {
			e.dispatch(true)
		} else {
			log.Print("Error: Response is not valid JSON\n")
		}
	}
	value := gjson.Get(string(bodyBytes), e.JSONRule.JSONPath)
	if value.Exists() {
		e.dispatch(evaluateCondition(value, e.JSONRule.Condition, e.JSONRule.Value))
	} else if e.JSONRule.Condition == JSON_RULE_EXISTS {
		e.dispatch(true)
	} else {
		log.Printf("Error: JSON path %s does not exist in response\n", e.JSONRule.JSONPath)
	}
}

func getValueAsString(value gjson.Result) string {
	switch value.Type {
	case gjson.String:
		return value.Str
	case gjson.Number:
		return fmt.Sprintf("%v", value.Num)
	case gjson.True:
		return "true"
	case gjson.False:
		return "false"
	case gjson.Null:
		return "null"
	default:
		return value.Raw
	}
}

func evaluateCondition(value gjson.Result, condition string, expectedValue string) bool {
	valueAsString := getValueAsString(value)
	switch condition {
	case JSON_RULE_EQUALS:
		return valueAsString != expectedValue
	case JSON_RULE_NOT_EQUAL:
		return valueAsString == expectedValue
	case JSON_RULE_GREATER_THAN:
		return value.Float() <= gjson.Parse(expectedValue).Float()
	case JSON_RULE_LESS_THAN:
		return value.Float() >= gjson.Parse(expectedValue).Float()
	default:
		return false
	}
}

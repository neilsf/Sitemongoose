package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate_Valid(t *testing.T) {
	event := Event{"status_code", 200, 0, nil, false, nil, "test"}
	valid, err := event.Validate()
	assert.True(t, valid)
	assert.Nil(t, err)
}

func TestValidate_InvalidTriggerType(t *testing.T) {
	event := Event{"invalid", 200, 1000, nil, false, nil, "test"}
	valid, err := event.Validate()
	assert.False(t, valid)
	assert.NotNil(t, err)
}

func TestValidate_InvalidStatusCode(t *testing.T) {
	event := Event{"status_code", 99, 1000, nil, false, nil, "test"}
	valid, err := event.Validate()
	assert.False(t, valid)
	assert.NotNil(t, err)
	event.ExpectedStatusCode = 600
	valid, err = event.Validate()
	assert.False(t, valid)
	assert.NotNil(t, err)
}

func TestValidate_InvalidResponseTime(t *testing.T) {
	event := Event{"response_time", 200, 0, nil, false, nil, "test"}
	valid, err := event.Validate()
	assert.False(t, valid)
	assert.NotNil(t, err)
}

func TestValidate_NoJSONRule(t *testing.T) {
	event := Event{"json_rule", 200, 1000, nil, false, nil, "test"}
	valid, err := event.Validate()
	assert.False(t, valid)
	assert.NotNil(t, err)
}

func TestValidate_NoJSONPath(t *testing.T) {
	event := Event{"json_rule", 200, 1000, &JSONRule{"", "eq", "value"}, false, nil, "test"}
	valid, err := event.Validate()
	assert.False(t, valid)
	assert.NotNil(t, err)
}

func TestValidate_InvalidJSONCondition(t *testing.T) {
	event := Event{"json_rule", 200, 1000, &JSONRule{"path", "invalid", "value"}, false, nil, "test"}
	valid, err := event.Validate()
	assert.False(t, valid)
	assert.NotNil(t, err)
}

func TestValidate_NoJSONValue(t *testing.T) {
	event := Event{"json_rule", 200, 1000, &JSONRule{"path", "eq", ""}, false, nil, "test"}
	valid, err := event.Validate()
	assert.False(t, valid)
	assert.NotNil(t, err)
	event.JSONRule.Condition = "valid"
	event.JSONRule.Value = ""
	event.JSONRule.JSONPath = ""
	valid, err = event.Validate()
	assert.True(t, valid)
	assert.Nil(t, err)
}

func TestValidate_ValidJSONRule(t *testing.T) {
	event := Event{"json_rule", 200, 1000, &JSONRule{"path.to.value", "eq", "value"}, false, nil, "test"}
	valid, err := event.Validate()
	assert.True(t, valid)
	assert.Nil(t, err)
}

func TestValidate_Numeric(t *testing.T) {
	event := Event{"json_rule", 200, 1000, &JSONRule{"path.to.value", "lt", "value"}, false, nil, "test"}
	valid, err := event.Validate()
	assert.False(t, valid)
	assert.NotNil(t, err)
	event.JSONRule.Condition = "gt"
	valid, err = event.Validate()
	assert.False(t, valid)
	assert.NotNil(t, err)
	event.JSONRule.Value = "123"
	valid, err = event.Validate()
	assert.True(t, valid)
	assert.Nil(t, err)
}

func TestValidate_InvalidRegexp(t *testing.T) {
	event := Event{"json_rule", 200, 1000, &JSONRule{"path.to.value", "regexp", "([a-zA-Z"}, false, nil, "test"}
	valid, err := event.Validate()
	assert.False(t, valid)
	assert.NotNil(t, err)
}

func TestCheckTrigger_StatusCode(t *testing.T) {
	event := Event{"status_code", 200, 0, nil, false, nil, "test"}
	event.CheckTrigger(200, 0, nil)
	assert.False(t, event.active)
	event.CheckTrigger(500, 0, nil)
	assert.True(t, event.active)
	event.CheckTrigger(200, 0, nil)
	assert.False(t, event.active)
}

func TestCheckTrigger_ResponseTime(t *testing.T) {
	event := Event{"response_time", 200, 1000, nil, false, nil, "test"}
	event.CheckTrigger(200, 1000, nil)
	assert.False(t, event.active)
	event.CheckTrigger(200, 1001, nil)
	assert.True(t, event.active)
	event.CheckTrigger(200, 1000, nil)
	assert.False(t, event.active)
}

func TestCheckTrigger_JSONRuleEq(t *testing.T) {
	var validJSONBody = []byte(`{"key": "value", "server_status": {"status": "ok", "code": 200, "up": true}}`)
	event1 := Event{"json_rule", 200, 1000, &JSONRule{"server_status.up", JSON_RULE_EQUALS, "true"}, false, nil, "test"}
	event1.CheckTrigger(200, 500, validJSONBody)
	assert.False(t, event1.active)
	validJSONBody = []byte(`{"key": "value", "server_status": {"status": "ok", "code": 200, "up": false}}`)
	event1.CheckTrigger(200, 500, validJSONBody)
	assert.True(t, event1.active)
}

func TestCheckTrigger_JSONRuleNeq(t *testing.T) {
	var validJSONBody = []byte(`{"key": "value", "server_status": {"status": "ok", "code": 200, "up": true}}`)
	event1 := Event{"json_rule", 200, 1000, &JSONRule{"server_status.up", JSON_RULE_NOT_EQUAL, "true"}, false, nil, "test"}
	event1.CheckTrigger(200, 500, validJSONBody)
	assert.True(t, event1.active)
	validJSONBody = []byte(`{"key": "value", "server_status": {"status": "ok", "code": 200, "up": false}}`)
	event1.CheckTrigger(200, 500, validJSONBody)
	assert.False(t, event1.active)
}

func TestCheckTrigger_InvalidJSON(t *testing.T) {
	var invalidJSONBody = []byte(`{thisis==NotValidJSON__/*`)
	event := Event{"json_rule", 200, 1000, &JSONRule{"path.to.value", JSON_RULE_EQUALS, "value"}, false, nil, "test"}
	event.CheckTrigger(200, 500, invalidJSONBody)
	assert.False(t, event.active)
	event = Event{"json_rule", 200, 1000, &JSONRule{"", "valid", ""}, true, nil, "test"}
	event.CheckTrigger(200, 500, invalidJSONBody)
	assert.True(t, event.active)
}

func TestCheckTrigger_JSONRuleExists(t *testing.T) {
	var validJSONBody = []byte(`{"key": "value", "server_status": {"status": "ok", "code": 200, "up": true}}`)
	event := Event{"json_rule", 200, 1000, &JSONRule{"server_status.up", JSON_RULE_EXISTS, ""}, false, nil, "test"}
	event.CheckTrigger(200, 500, validJSONBody)
	assert.False(t, event.active)
	event.JSONRule.JSONPath = "server_status.notexist"
	event.CheckTrigger(200, 500, validJSONBody)
	assert.True(t, event.active)
}

func TestCheckTrigger_JSONRuleLt(t *testing.T) {
	var validJSONBody = []byte(`{"key": "value", "server_status": {"status": "ok", "connections": 200, "up": true}}`)
	event1 := Event{"json_rule", 200, 1000, &JSONRule{"server_status.connections", JSON_RULE_LESS_THAN, "200"}, false, nil, "test"}
	event1.CheckTrigger(200, 500, validJSONBody)
	assert.True(t, event1.active)
	validJSONBody = []byte(`{"key": "value", "server_status": {"status": "ok", "connections": 199.9, "up": false}}`)
	event1.CheckTrigger(200, 500, validJSONBody)
	assert.False(t, event1.active)
}

func TestCheckTrigger_JSONRuleGt(t *testing.T) {
	var validJSONBody = []byte(`{"key": "value", "server_status": {"status": "ok", "connections": 200, "up": true}}`)
	event1 := Event{"json_rule", 200, 1000, &JSONRule{"server_status.connections", JSON_RULE_GREATER_THAN, "200"}, false, nil, "test"}
	event1.CheckTrigger(200, 500, validJSONBody)
	assert.True(t, event1.active)
	validJSONBody = []byte(`{"key": "value", "server_status": {"status": "ok", "connections": 200.1, "up": false}}`)
	event1.CheckTrigger(200, 500, validJSONBody)
	assert.False(t, event1.active)
}

func TestCheckTrigger_JSONRuleRegexp(t *testing.T) {
	var validJSONBody = []byte(`{"key": "value", "server_status": {"color": "gray"}}`)
	event1 := Event{"json_rule", 200, 1000, &JSONRule{"server_status.color", JSON_RULE_REGEXP, "gr(a|e)y"}, false, nil, "test"}
	event1.CheckTrigger(200, 500, validJSONBody)
	assert.False(t, event1.active)
	validJSONBody = []byte(`{"key": "value", "server_status": {"color": "grey"}}`)
	event1.CheckTrigger(200, 500, validJSONBody)
	assert.False(t, event1.active)
	event1.JSONRule.Value = "go*gle"
	event1.CheckTrigger(200, 500, validJSONBody)
	assert.True(t, event1.active)
}

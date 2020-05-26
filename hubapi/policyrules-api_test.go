// Copyright 2020 Synopsys, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hubapi

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePolicyRulesBOMComponent(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/policyRulesBOMComponent.json")
	assert.NoError(t, err)
	assert.NotEmpty(t, content)

	policyRules := BomComponentPolicyRulesList{}
	err = json.Unmarshal(content, &policyRules)
	assert.NoError(t, err)

	assert.Equal(t, 2, policyRules.TotalCount)
	assert.Len(t, policyRules.Items, 2)

	rule := policyRules.Items[1]
	assert.NotEmpty(t, rule.CreatedAt)
	assert.Len(t, rule.Expression.Expressions, 2)

	expression := rule.Expression.Expressions[0]
	assert.Equal(t, "HIGH_SEVERITY_VULN_COUNT", expression.Name)
	assert.Len(t, expression.Parameters.Values, 1)
	assert.Equal(t, "0", expression.Parameters.Values[0])
	assert.Len(t, expression.Parameters.Data, 1)
	assert.Equal(t, float64(0), expression.Parameters.Data[0]["data"])

	expression = rule.Expression.Expressions[1]
	assert.Equal(t, "VULN_LEVEL_EXPLOIT", expression.Name)
	assert.Len(t, expression.Parameters.Values, 1)
	assert.Equal(t, "TRUE", expression.Parameters.Values[0])
	assert.Len(t, expression.Parameters.Data, 1)
	assert.Equal(t, true, expression.Parameters.Data[0]["data"])
}

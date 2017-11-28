// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package inject

import (
	"testing"

	"istio.io/istio/pilot/proxy"
)

func TestCreateMeshTemplate(t *testing.T) {
	mesh := proxy.DefaultMeshConfig()
	cases := []struct {
		mc *Config
	}{
		{ // Custom generated
			mc: &Config{
				Init{"test", []string{"testarg1", "testarg1"}, "testpull"},
				Container{"test", []string{"- testarg1", serviceDefault}, "testpull", []string{"testenv"}},
				mesh,
			},
		},
		{ // Default generated
			mc: &Config{
				Init{"", []string{}, ""},
				Container{"", []string{}, "", []string{}},
				mesh,
			},
		},
	}

	for _, test := range cases {
		err := DefaultMeshTemplate(test.mc)
		if err != nil {
			t.Fatal(err.Error())
		}

	}
}

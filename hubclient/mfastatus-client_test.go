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

package hubclient

import (
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestFetchMfaStatus(t *testing.T) {
	client, err := NewWithSession("https://localhost", HubClientDebugTimings, 5*time.Second)
	if err != nil {
		t.Errorf("unable to instantiate client: %s", err.Error())
		return
	}
	mfaStatus, err := client.MfaStatus()
	if err != nil {
		t.Errorf("unable to get MFA status: %s", err.Error())
		return
	}

	log.Infof("successfully fetched MFA status, MFA Enabled? %+v", mfaStatus.MfaEnabled)
}

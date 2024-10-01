// Copyright © 2024 Ha Nguyen <captainnemot1k60@gmail.com>
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
//

package helper

import (
	"strings"
)

var (
	// IsGnuTar show tar type
	IsGnuTar = false
)

func init() {
	checkIsGnuTar()
}

func checkIsGnuTar() {
	out, _ := Exec("tar", "--version")
	IsGnuTar = strings.Contains(out, "GNU")
}

// CleanHost clean host url ftp://foo.bar.com -> foo.bar.com
func CleanHost(host string) string {
	// ftp://ftp.your-host.com -> ftp.your-host.com
	if strings.Contains(host, "://") {
		return strings.Split(host, "://")[1]
	}

	return host
}

// FormatEndpoint to add `https://` prefix if not present
func FormatEndpoint(endpoint string) string {
	if !strings.HasPrefix(endpoint, "http") {
		endpoint = "https://" + endpoint
	}

	return endpoint
}

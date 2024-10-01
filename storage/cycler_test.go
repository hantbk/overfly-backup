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

package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCycler_add(t *testing.T) {
	cycler := Cycler{}
	cycler.add("foo", []string{})
	cycler.add("bar", []string{})

	assert.Equal(t, len(cycler.packages), 2)
}

func TestCycler_shiftByKeep(t *testing.T) {
	cycler := Cycler{
		packages: PackageList{
			Package{
				FileKey:   "p1",
				CreatedAt: time.Now(),
			},
			Package{
				FileKey:   "p2",
				CreatedAt: time.Now(),
			},
		},
	}
	cycler.add("p3", []string{})
	cycler.add("p4", []string{})
	cycler.add("p5", []string{})
	cycler.add("p6", []string{})

	pkg := cycler.shiftByKeep(2)
	assert.Equal(t, len(cycler.packages), 5)
	assert.Equal(t, pkg.FileKey, "p1")
	pkg = cycler.shiftByKeep(2)
	assert.Equal(t, len(cycler.packages), 4)
	assert.Equal(t, pkg.FileKey, "p2")
	pkg = cycler.shiftByKeep(4)
	assert.Equal(t, len(cycler.packages), 4)
	assert.Nil(t, pkg)
}

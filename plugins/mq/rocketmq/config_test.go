// Copyright 2023 The Ryan SU Authors (https://github.com/suyuan32). All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rocketmq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConf_Validate(t *testing.T) {
	c := &ProducerConf{
		NsResolver:                 []string{"127.0.0.1:9876"},
		GroupName:                  "",
		Namespace:                  "",
		InstanceName:               "",
		MsgTimeOut:                 0,
		DefaultTopicQueueNums:      0,
		CreateTopicKey:             "",
		CompressMsgBodyOverHowMuch: 0,
		CompressLevel:              0,
		Retry:                      0,
	}

	err := c.Validate()
	assert.Nil(t, err)
}

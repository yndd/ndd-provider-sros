/*
Copyright 2021 NDD.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package collector

import (
	"fmt"
	"strings"

	"github.com/karimra/gnmic/utils"
	"github.com/openconfig/gnmi/proto/gnmi"
)

// CreateSubscriptionRequest create a gnmi subscription
func CreateSubscriptionRequest(target, subName string, paths []*gnmi.Path) (*gnmi.SubscribeRequest, error) {
	// create subscription

	gnmiPrefix, err := utils.CreatePrefix(subName, target)
	if err != nil {
		return nil, fmt.Errorf("create prefix failed")
	}
	modeVal := gnmi.SubscriptionList_Mode_value[strings.ToUpper("STREAM")]
	qos := &gnmi.QOSMarking{Marking: 21}

	subscriptions := make([]*gnmi.Subscription, len(paths))
	for i, p := range paths {
		subscriptions[i] = &gnmi.Subscription{Path: p}
		switch gnmi.SubscriptionList_Mode(modeVal) {
		case gnmi.SubscriptionList_STREAM:
			mode := gnmi.SubscriptionMode_value[strings.Replace(strings.ToUpper("ON_CHANGE"), "-", "_", -1)]
			subscriptions[i].Mode = gnmi.SubscriptionMode(mode)
		}
	}
	req := &gnmi.SubscribeRequest{
		Request: &gnmi.SubscribeRequest_Subscribe{
			Subscribe: &gnmi.SubscriptionList{
				Prefix:       gnmiPrefix,
				Mode:         gnmi.SubscriptionList_Mode(modeVal),
				Encoding:     46, // "JSON_IETF_CONFIG_ONLY"
				Subscription: subscriptions,
				Qos:          qos,
			},
		},
	}
	return req, nil
}

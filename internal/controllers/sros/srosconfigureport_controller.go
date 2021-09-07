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

package sros

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/karimra/gnmic/target"
	gnmitypes "github.com/karimra/gnmic/types"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi_ext"
	"github.com/pkg/errors"
	ndrv1 "github.com/yndd/ndd-core/apis/dvr/v1"
	nddv1 "github.com/yndd/ndd-runtime/apis/common/v1"
	"github.com/yndd/ndd-runtime/pkg/event"
	"github.com/yndd/ndd-runtime/pkg/gext"
	"github.com/yndd/ndd-runtime/pkg/gvk"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-runtime/pkg/reconciler/managed"
	"github.com/yndd/ndd-runtime/pkg/resource"
	"github.com/yndd/ndd-runtime/pkg/utils"
	"github.com/yndd/ndd-yang/pkg/parser"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	cevent "sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	srosv1alpha1 "github.com/yndd/ndd-provider-sros/apis/sros/v1alpha1"
)

const (
	// Errors
	errUnexpectedConfigurePort       = "the managed resource is not a ConfigurePort resource"
	errKubeUpdateFailedConfigurePort = "cannot update ConfigurePort"
	errReadConfigurePort             = "cannot read ConfigurePort"
	errCreateConfigurePort           = "cannot create ConfigurePort"
	erreUpdateConfigurePort          = "cannot update ConfigurePort"
	errDeleteConfigurePort           = "cannot delete ConfigurePort"

	// resource information
	levelConfigurePort = 2
	// resourcePrefixConfigurePort = "sros.ndd.yndd.io.v1alpha1.ConfigurePort"
)

var resourceRefPathsConfigurePort = []*gnmi.Path{
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "access"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "access"},
			{Name: "egress"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "pool", Key: map[string]string{"name": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "pool", Key: map[string]string{"name": ""}},
			{Name: "resv-cbs"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "pool", Key: map[string]string{"name": ""}},
			{Name: "resv-cbs"},
			{Name: "amber-alarm-action"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "access"},
			{Name: "ingress"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "access"},
			{Name: "ingress"},
			{Name: "pool", Key: map[string]string{"name": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "access"},
			{Name: "ingress"},
			{Name: "pool", Key: map[string]string{"name": ""}},
			{Name: "resv-cbs"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "access"},
			{Name: "ingress"},
			{Name: "pool", Key: map[string]string{"name": ""}},
			{Name: "resv-cbs"},
			{Name: "amber-alarm-action"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "connector"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "dist-cpu-protection"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "dwdm"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "dwdm"},
			{Name: "coherent"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "dwdm"},
			{Name: "coherent"},
			{Name: "report-alarm"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "dwdm"},
			{Name: "coherent"},
			{Name: "sweep"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "dwdm"},
			{Name: "wavetracker"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "dwdm"},
			{Name: "wavetracker"},
			{Name: "encode"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "dwdm"},
			{Name: "wavetracker"},
			{Name: "power-control"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "dwdm"},
			{Name: "wavetracker"},
			{Name: "report-alarm"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "aggregate-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "host-match"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "host-match"},
			{Name: "int-dest-id", Key: map[string]string{"destination-string": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "hsmda-queue-overrides"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "hsmda-queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "adaptation-rule"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "drop-tail"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "drop-tail"},
			{Name: "low"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "monitor-queue-depth"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "parent"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "queue-override-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "queue-override-rate"},
			{Name: "percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "queue-override-rate"},
			{Name: "percent-rate"},
			{Name: "percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "queue-override-rate"},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "queue-override-rate"},
			{Name: "rate"},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "scheduler-policy"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "scheduler-policy"},
			{Name: "overrides"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "scheduler-policy"},
			{Name: "overrides"},
			{Name: "scheduler", Key: map[string]string{"scheduler-name": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "scheduler-policy"},
			{Name: "overrides"},
			{Name: "scheduler", Key: map[string]string{"scheduler-name": ""}},
			{Name: "parent"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "scheduler-policy"},
			{Name: "overrides"},
			{Name: "scheduler", Key: map[string]string{"scheduler-name": ""}},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "virtual-port", Key: map[string]string{"vport-name": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "virtual-port", Key: map[string]string{"vport-name": ""}},
			{Name: "aggregate-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "virtual-port", Key: map[string]string{"vport-name": ""}},
			{Name: "host-match"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "egress"},
			{Name: "virtual-port", Key: map[string]string{"vport-name": ""}},
			{Name: "host-match"},
			{Name: "int-dest-id", Key: map[string]string{"destination-string": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "ingress"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "ingress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "ingress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "ingress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "ingress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "adaptation-rule"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "ingress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "drop-tail"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "ingress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "drop-tail"},
			{Name: "low"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "ingress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "monitor-queue-depth"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "ingress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "ingress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "scheduler-policy"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "ingress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "scheduler-policy"},
			{Name: "overrides"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "ingress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "scheduler-policy"},
			{Name: "overrides"},
			{Name: "scheduler", Key: map[string]string{"scheduler-name": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "ingress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "scheduler-policy"},
			{Name: "overrides"},
			{Name: "scheduler", Key: map[string]string{"scheduler-name": ""}},
			{Name: "parent"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "access"},
			{Name: "ingress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "scheduler-policy"},
			{Name: "overrides"},
			{Name: "scheduler", Key: map[string]string{"scheduler-name": ""}},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "crc-monitor"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "crc-monitor"},
			{Name: "signal-degrade"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "crc-monitor"},
			{Name: "signal-failure"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "dampening"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "dot1x"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "dot1x"},
			{Name: "macsec"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "dot1x"},
			{Name: "macsec"},
			{Name: "exclude-protocol"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "dot1x"},
			{Name: "macsec"},
			{Name: "sub-port", Key: map[string]string{"sub-port-id": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "dot1x"},
			{Name: "macsec"},
			{Name: "sub-port", Key: map[string]string{"sub-port-id": ""}},
			{Name: "encap-match"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "dot1x"},
			{Name: "macsec"},
			{Name: "sub-port", Key: map[string]string{"sub-port-id": ""}},
			{Name: "encap-match"},
			{Name: "encap"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "dot1x"},
			{Name: "macsec"},
			{Name: "sub-port", Key: map[string]string{"sub-port-id": ""}},
			{Name: "encap-match"},
			{Name: "encap"},
			{Name: "all-match"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "dot1x"},
			{Name: "macsec"},
			{Name: "sub-port", Key: map[string]string{"sub-port-id": ""}},
			{Name: "encap-match"},
			{Name: "encap"},
			{Name: "double-tag"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "dot1x"},
			{Name: "macsec"},
			{Name: "sub-port", Key: map[string]string{"sub-port-id": ""}},
			{Name: "encap-match"},
			{Name: "encap"},
			{Name: "single-tag"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "dot1x"},
			{Name: "macsec"},
			{Name: "sub-port", Key: map[string]string{"sub-port-id": ""}},
			{Name: "encap-match"},
			{Name: "encap"},
			{Name: "untagged"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "dot1x"},
			{Name: "per-host-authentication"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "dot1x"},
			{Name: "per-host-authentication"},
			{Name: "allowed-source-macs"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "dot1x"},
			{Name: "per-host-authentication"},
			{Name: "allowed-source-macs"},
			{Name: "mac-address", Key: map[string]string{"mac": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "dot1x"},
			{Name: "radius-server-policy-config"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "dot1x"},
			{Name: "radius-server-policy-config"},
			{Name: "common"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "dot1x"},
			{Name: "radius-server-policy-config"},
			{Name: "split"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "dot1x"},
			{Name: "re-authentication"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "down-on-internal-error"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "down-when-looped"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "efm-oam"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "efm-oam"},
			{Name: "discovery"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "efm-oam"},
			{Name: "discovery"},
			{Name: "advertise-capabilities"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "efm-oam"},
			{Name: "link-monitoring"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "efm-oam"},
			{Name: "link-monitoring"},
			{Name: "errored-frame"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "efm-oam"},
			{Name: "link-monitoring"},
			{Name: "errored-frame-period"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "efm-oam"},
			{Name: "link-monitoring"},
			{Name: "errored-frame-seconds"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "efm-oam"},
			{Name: "link-monitoring"},
			{Name: "errored-symbols"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "efm-oam"},
			{Name: "link-monitoring"},
			{Name: "local-sf-action"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "efm-oam"},
			{Name: "link-monitoring"},
			{Name: "local-sf-action"},
			{Name: "info-notification"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "efm-oam"},
			{Name: "peer-rdi-rx"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "expanded-secondary-shaper", Key: map[string]string{"secondary-shaper-name": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "expanded-secondary-shaper", Key: map[string]string{"secondary-shaper-name": ""}},
			{Name: "aggregate-burst"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "expanded-secondary-shaper", Key: map[string]string{"secondary-shaper-name": ""}},
			{Name: "class", Key: map[string]string{"class-number": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "hs-scheduler-policy"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "hs-scheduler-policy"},
			{Name: "overrides"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "hs-scheduler-policy"},
			{Name: "overrides"},
			{Name: "group", Key: map[string]string{"group-id": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "hs-scheduler-policy"},
			{Name: "overrides"},
			{Name: "scheduling-class", Key: map[string]string{"class-number": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "hs-secondary-shaper", Key: map[string]string{"secondary-shaper-name": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "hs-secondary-shaper", Key: map[string]string{"secondary-shaper-name": ""}},
			{Name: "aggregate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "hs-secondary-shaper", Key: map[string]string{"secondary-shaper-name": ""}},
			{Name: "class", Key: map[string]string{"class-number": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "port-qos-policy"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "percent-rate"},
			{Name: "percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "rate"},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
			{Name: "rate-or-percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
			{Name: "rate-or-percent-rate"},
			{Name: "percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
			{Name: "rate-or-percent-rate"},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "elmi"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "eth-cfm"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "eth-cfm"},
			{Name: "mep", Key: map[string]string{"md-admin-name": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "eth-cfm"},
			{Name: "mep", Key: map[string]string{"md-admin-name": ""}},
			{Name: "ais"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "eth-cfm"},
			{Name: "mep", Key: map[string]string{"md-admin-name": ""}},
			{Name: "alarm-notification"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "eth-cfm"},
			{Name: "mep", Key: map[string]string{"md-admin-name": ""}},
			{Name: "csf"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "eth-cfm"},
			{Name: "mep", Key: map[string]string{"md-admin-name": ""}},
			{Name: "eth-bn"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "eth-cfm"},
			{Name: "mep", Key: map[string]string{"md-admin-name": ""}},
			{Name: "eth-test"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "eth-cfm"},
			{Name: "mep", Key: map[string]string{"md-admin-name": ""}},
			{Name: "eth-test"},
			{Name: "test-pattern"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "eth-cfm"},
			{Name: "mep", Key: map[string]string{"md-admin-name": ""}},
			{Name: "grace"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "eth-cfm"},
			{Name: "mep", Key: map[string]string{"md-admin-name": ""}},
			{Name: "grace"},
			{Name: "eth-ed"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "eth-cfm"},
			{Name: "mep", Key: map[string]string{"md-admin-name": ""}},
			{Name: "grace"},
			{Name: "eth-vsm-grace"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "hold-time"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "hsmda-scheduler-overrides"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "hsmda-scheduler-overrides"},
			{Name: "group", Key: map[string]string{"group-id": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "hsmda-scheduler-overrides"},
			{Name: "scheduling-class", Key: map[string]string{"class-number": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "ingress"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "lldp"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "lldp"},
			{Name: "dest-mac", Key: map[string]string{"mac-type": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "lldp"},
			{Name: "dest-mac", Key: map[string]string{"mac-type": ""}},
			{Name: "tx-mgmt-address", Key: map[string]string{"mgmt-address-system-type": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "lldp"},
			{Name: "dest-mac", Key: map[string]string{"mac-type": ""}},
			{Name: "tx-tlvs"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "loopback"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "network"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "network"},
			{Name: "egress"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "network"},
			{Name: "egress"},
			{Name: "port-queues"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "network"},
			{Name: "egress"},
			{Name: "port-queues"},
			{Name: "overrides"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "network"},
			{Name: "egress"},
			{Name: "port-queues"},
			{Name: "overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "network"},
			{Name: "egress"},
			{Name: "port-queues"},
			{Name: "overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "monitor-queue-depth"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "network"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "network"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "aggregate-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "network"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "network"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "network"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "adaptation-rule"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "network"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "drop-tail"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "network"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "drop-tail"},
			{Name: "low"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "network"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "monitor-queue-depth"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "network"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "queue-override-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "network"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "queue-override-rate"},
			{Name: "percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "network"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "queue-override-rate"},
			{Name: "percent-rate"},
			{Name: "percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "network"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "queue-override-rate"},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "network"},
			{Name: "egress"},
			{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
			{Name: "queue-overrides"},
			{Name: "queue", Key: map[string]string{"queue-id": ""}},
			{Name: "queue-override-rate"},
			{Name: "rate"},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "report-alarm"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "ssm"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "symbol-monitor"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "symbol-monitor"},
			{Name: "signal-degrade"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "ethernet"},
			{Name: "symbol-monitor"},
			{Name: "signal-failure"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "gnss"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "gnss"},
			{Name: "constellation"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "hybrid-buffer-allocation"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "hybrid-buffer-allocation"},
			{Name: "egress-weight"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "hybrid-buffer-allocation"},
			{Name: "ingress-weight"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "modify-buffer-allocation"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "modify-buffer-allocation"},
			{Name: "percentage-of-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "network"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "network"},
			{Name: "egress"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "network"},
			{Name: "egress"},
			{Name: "pool", Key: map[string]string{"name": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "network"},
			{Name: "egress"},
			{Name: "pool", Key: map[string]string{"name": ""}},
			{Name: "resv-cbs"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "network"},
			{Name: "egress"},
			{Name: "pool", Key: map[string]string{"name": ""}},
			{Name: "resv-cbs"},
			{Name: "amber-alarm-action"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "fine-granularity-ber"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "fine-granularity-ber"},
			{Name: "signal-degrade"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "fine-granularity-ber"},
			{Name: "signal-degrade"},
			{Name: "clear"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "fine-granularity-ber"},
			{Name: "signal-degrade"},
			{Name: "raise"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "fine-granularity-ber"},
			{Name: "signal-failure"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "fine-granularity-ber"},
			{Name: "signal-failure"},
			{Name: "clear"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "fine-granularity-ber"},
			{Name: "signal-failure"},
			{Name: "raise"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "path-monitoring"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "path-monitoring"},
			{Name: "trail-trace-identifier"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "path-monitoring"},
			{Name: "trail-trace-identifier"},
			{Name: "expected"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "path-monitoring"},
			{Name: "trail-trace-identifier"},
			{Name: "expected"},
			{Name: "expected"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "path-monitoring"},
			{Name: "trail-trace-identifier"},
			{Name: "expected"},
			{Name: "expected"},
			{Name: "auto-generated"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "path-monitoring"},
			{Name: "trail-trace-identifier"},
			{Name: "expected"},
			{Name: "expected"},
			{Name: "bytes"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "path-monitoring"},
			{Name: "trail-trace-identifier"},
			{Name: "expected"},
			{Name: "expected"},
			{Name: "string"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "path-monitoring"},
			{Name: "trail-trace-identifier"},
			{Name: "transmit"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "path-monitoring"},
			{Name: "trail-trace-identifier"},
			{Name: "transmit"},
			{Name: "transmit"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "path-monitoring"},
			{Name: "trail-trace-identifier"},
			{Name: "transmit"},
			{Name: "transmit"},
			{Name: "auto-generated"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "path-monitoring"},
			{Name: "trail-trace-identifier"},
			{Name: "transmit"},
			{Name: "transmit"},
			{Name: "bytes"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "path-monitoring"},
			{Name: "trail-trace-identifier"},
			{Name: "transmit"},
			{Name: "transmit"},
			{Name: "string"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "payload-structure-identifier"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "payload-structure-identifier"},
			{Name: "payload"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "report-alarm"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "section-monitoring"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "section-monitoring"},
			{Name: "trail-trace-identifier"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "section-monitoring"},
			{Name: "trail-trace-identifier"},
			{Name: "expected"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "section-monitoring"},
			{Name: "trail-trace-identifier"},
			{Name: "expected"},
			{Name: "expected"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "section-monitoring"},
			{Name: "trail-trace-identifier"},
			{Name: "expected"},
			{Name: "expected"},
			{Name: "auto-generated"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "section-monitoring"},
			{Name: "trail-trace-identifier"},
			{Name: "expected"},
			{Name: "expected"},
			{Name: "bytes"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "section-monitoring"},
			{Name: "trail-trace-identifier"},
			{Name: "expected"},
			{Name: "expected"},
			{Name: "string"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "section-monitoring"},
			{Name: "trail-trace-identifier"},
			{Name: "transmit"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "section-monitoring"},
			{Name: "trail-trace-identifier"},
			{Name: "transmit"},
			{Name: "transmit"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "section-monitoring"},
			{Name: "trail-trace-identifier"},
			{Name: "transmit"},
			{Name: "transmit"},
			{Name: "auto-generated"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "section-monitoring"},
			{Name: "trail-trace-identifier"},
			{Name: "transmit"},
			{Name: "transmit"},
			{Name: "bytes"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "otu"},
			{Name: "section-monitoring"},
			{Name: "trail-trace-identifier"},
			{Name: "transmit"},
			{Name: "transmit"},
			{Name: "string"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "group", Key: map[string]string{"group-index": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "hold-time"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "path", Key: map[string]string{"path-index": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "path", Key: map[string]string{"path-index": ""}},
			{Name: "egress"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "path", Key: map[string]string{"path-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "path", Key: map[string]string{"path-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "path", Key: map[string]string{"path-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "path", Key: map[string]string{"path-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "path", Key: map[string]string{"path-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "path", Key: map[string]string{"path-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "percent-rate"},
			{Name: "percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "path", Key: map[string]string{"path-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "path", Key: map[string]string{"path-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "rate"},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "path", Key: map[string]string{"path-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "path", Key: map[string]string{"path-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
			{Name: "rate-or-percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "path", Key: map[string]string{"path-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
			{Name: "rate-or-percent-rate"},
			{Name: "percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "path", Key: map[string]string{"path-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
			{Name: "rate-or-percent-rate"},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "path", Key: map[string]string{"path-index": ""}},
			{Name: "network"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "path", Key: map[string]string{"path-index": ""}},
			{Name: "ppp"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "path", Key: map[string]string{"path-index": ""}},
			{Name: "ppp"},
			{Name: "keepalive"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "path", Key: map[string]string{"path-index": ""}},
			{Name: "report-alarm"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "report-alarm"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "section-trace"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "section-trace"},
			{Name: "section-trace"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "section-trace"},
			{Name: "section-trace"},
			{Name: "byte"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "section-trace"},
			{Name: "section-trace"},
			{Name: "increment-z0"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "sonet-sdh"},
			{Name: "section-trace"},
			{Name: "section-trace"},
			{Name: "string"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "ber-threshold"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "percent-rate"},
			{Name: "percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "rate"},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
			{Name: "rate-or-percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
			{Name: "rate-or-percent-rate"},
			{Name: "percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
			{Name: "rate-or-percent-rate"},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "idle-payload-fill"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "idle-payload-fill"},
			{Name: "idle-payload-fill-choice"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "idle-payload-fill"},
			{Name: "idle-payload-fill-choice"},
			{Name: "all-ones"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "idle-payload-fill"},
			{Name: "idle-payload-fill-choice"},
			{Name: "pattern"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "idle-signal-fill"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "idle-signal-fill"},
			{Name: "idle-signal-fill-choice"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "idle-signal-fill"},
			{Name: "idle-signal-fill-choice"},
			{Name: "all-ones"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "idle-signal-fill"},
			{Name: "idle-signal-fill-choice"},
			{Name: "pattern"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "network"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "ppp"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "ppp"},
			{Name: "compress"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "ppp"},
			{Name: "keepalive"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "hold-time"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
			{Name: "report-alarm"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
			{Name: "egress"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "percent-rate"},
			{Name: "percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "rate"},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
			{Name: "rate-or-percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
			{Name: "rate-or-percent-rate"},
			{Name: "percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
			{Name: "rate-or-percent-rate"},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
			{Name: "maintenance-data-link"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
			{Name: "maintenance-data-link"},
			{Name: "transmit-message-type"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
			{Name: "network"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
			{Name: "ppp"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
			{Name: "ppp"},
			{Name: "keepalive"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
			{Name: "report-alarm"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
			{Name: "subrate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "ber-threshold"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "percent-rate"},
			{Name: "percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "rate"},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
			{Name: "rate-or-percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
			{Name: "rate-or-percent-rate"},
			{Name: "percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
			{Name: "rate-or-percent-rate"},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "idle-payload-fill"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "idle-payload-fill"},
			{Name: "idle-payload-fill-choice"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "idle-payload-fill"},
			{Name: "idle-payload-fill-choice"},
			{Name: "all-ones"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "idle-payload-fill"},
			{Name: "idle-payload-fill-choice"},
			{Name: "pattern"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "idle-signal-fill"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "idle-signal-fill"},
			{Name: "idle-signal-fill-choice"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "idle-signal-fill"},
			{Name: "idle-signal-fill-choice"},
			{Name: "all-ones"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "idle-signal-fill"},
			{Name: "idle-signal-fill-choice"},
			{Name: "pattern"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "network"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "ppp"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "ppp"},
			{Name: "compress"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
			{Name: "ppp"},
			{Name: "keepalive"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "hold-time"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "national-bits"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e1", Key: map[string]string{"e1-index": ""}},
			{Name: "report-alarm"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e3", Key: map[string]string{"e3-index": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e3", Key: map[string]string{"e3-index": ""}},
			{Name: "egress"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e3", Key: map[string]string{"e3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e3", Key: map[string]string{"e3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e3", Key: map[string]string{"e3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e3", Key: map[string]string{"e3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e3", Key: map[string]string{"e3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e3", Key: map[string]string{"e3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "percent-rate"},
			{Name: "percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e3", Key: map[string]string{"e3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e3", Key: map[string]string{"e3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "level", Key: map[string]string{"priority-level": ""}},
			{Name: "rate-or-percent-rate"},
			{Name: "rate"},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e3", Key: map[string]string{"e3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e3", Key: map[string]string{"e3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
			{Name: "rate-or-percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e3", Key: map[string]string{"e3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
			{Name: "rate-or-percent-rate"},
			{Name: "percent-rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e3", Key: map[string]string{"e3-index": ""}},
			{Name: "egress"},
			{Name: "port-scheduler-policy"},
			{Name: "overrides"},
			{Name: "max-rate"},
			{Name: "rate-or-percent-rate"},
			{Name: "rate"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e3", Key: map[string]string{"e3-index": ""}},
			{Name: "network"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e3", Key: map[string]string{"e3-index": ""}},
			{Name: "ppp"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e3", Key: map[string]string{"e3-index": ""}},
			{Name: "ppp"},
			{Name: "keepalive"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "e3", Key: map[string]string{"e3-index": ""}},
			{Name: "report-alarm"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "tdm"},
			{Name: "hold-time"},
		},
	},
	{
		Elem: []*gnmi.PathElem{
			{Name: "port-xc"},
			{Name: "transceiver"},
		},
	},
}
var dependencyConfigurePort = []*parser.LeafRefGnmi{}
var localleafRefConfigurePort = []*parser.LeafRefGnmi{
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "queue-group"},
				{Name: "hsmda-queue-overrides"},
				{Name: "secondary-shaper"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "expanded-secondary-shaper", Key: map[string]string{"secondary-shaper-name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "dot1x"},
				{Name: "radius-server-policy-config"},
				{Name: "common"},
				{Name: "radius-server-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "aaa"},
				{Name: "radius"},
				{Name: "server-policy", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "dot1x"},
				{Name: "radius-server-policy-config"},
				{Name: "split"},
				{Name: "radius-server-policy-acct"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "aaa"},
				{Name: "radius"},
				{Name: "server-policy", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "dot1x"},
				{Name: "radius-server-policy-config"},
				{Name: "split"},
				{Name: "radius-server-policy-auth"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "aaa"},
				{Name: "radius"},
				{Name: "server-policy", Key: map[string]string{"name": ""}},
			},
		},
	},
}
var externalLeafRefConfigurePort = []*parser.LeafRefGnmi{
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "access"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "access"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "pool", Key: map[string]string{"name": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "pool", Key: map[string]string{"name": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "pool", Key: map[string]string{"name": ""}},
				{Name: "slope-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "slope-policy", Key: map[string]string{"slope-policy-name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "access"},
				{Name: "ingress"},
				{Name: "pool", Key: map[string]string{"name": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "access"},
				{Name: "ingress"},
				{Name: "pool", Key: map[string]string{"name": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "access"},
				{Name: "ingress"},
				{Name: "pool", Key: map[string]string{"name": ""}},
				{Name: "slope-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "slope-policy", Key: map[string]string{"slope-policy-name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "connector"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "connector"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "dist-cpu-protection"},
				{Name: "policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "system"},
				{Name: "security"},
				{Name: "dist-cpu-protection"},
				{Name: "policy", Key: map[string]string{"policy-name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "dwdm"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "dwdm"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "dwdm"},
				{Name: "coherent"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "dwdm"},
				{Name: "coherent"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "dwdm"},
				{Name: "wavetracker"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "dwdm"},
				{Name: "wavetracker"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "accounting-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "log"},
				{Name: "accounting-policy", Key: map[string]string{"policy-id": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "accounting-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "log"},
				{Name: "accounting-policy", Key: map[string]string{"policy-id": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "hsmda-queue-overrides"},
				{Name: "queue", Key: map[string]string{"queue-id": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "hsmda-queue-overrides"},
				{Name: "queue", Key: map[string]string{"queue-id": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "hsmda-queue-overrides"},
				{Name: "queue", Key: map[string]string{"queue-id": ""}},
				{Name: "queue-id"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "qos"},
				{Name: "queue-group-templates"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"egress-queue-group-name": ""}},
				{Name: "queue-group-name]"},
				{Name: "hsmda-queues"},
				{Name: "queue", Key: map[string]string{"queue-id": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "hsmda-queue-overrides"},
				{Name: "queue", Key: map[string]string{"queue-id": ""}},
				{Name: "slope-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "hsmda-slope-policy", Key: map[string]string{"hsmda-slope-policy-name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "hsmda-queue-overrides"},
				{Name: "wrr-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "hsmda-wrr-policy", Key: map[string]string{"hsmda-wrr-policy-name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "queue-group-name"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "queue-group-templates"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"egress-queue-group-name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "queue-overrides"},
				{Name: "queue", Key: map[string]string{"queue-id": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "queue-overrides"},
				{Name: "queue", Key: map[string]string{"queue-id": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "queue-overrides"},
				{Name: "queue", Key: map[string]string{"queue-id": ""}},
				{Name: "queue-id"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "qos"},
				{Name: "queue-group-templates"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"egress-queue-group-name": ""}},
				{Name: "queue-group-name]"},
				{Name: "queue", Key: map[string]string{"queue-id": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "scheduler-policy"},
				{Name: "overrides"},
				{Name: "scheduler", Key: map[string]string{"scheduler-name": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "scheduler-policy"},
				{Name: "overrides"},
				{Name: "scheduler", Key: map[string]string{"scheduler-name": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "scheduler-policy"},
				{Name: "policy-name"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "scheduler-policy", Key: map[string]string{"scheduler-policy-name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "virtual-port", Key: map[string]string{"vport-name": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "virtual-port", Key: map[string]string{"vport-name": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "virtual-port", Key: map[string]string{"vport-name": ""}},
				{Name: "hw-agg-shaper-scheduler-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "hw-agg-shaper-scheduler-policy", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "virtual-port", Key: map[string]string{"vport-name": ""}},
				{Name: "port-scheduler-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "port-scheduler-policy", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "egress"},
				{Name: "virtual-port", Key: map[string]string{"vport-name": ""}},
				{Name: "scheduler-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "scheduler-policy", Key: map[string]string{"scheduler-policy-name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "ingress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
				{Name: "accounting-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "log"},
				{Name: "accounting-policy", Key: map[string]string{"policy-id": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "ingress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "ingress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "ingress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
				{Name: "queue-group-name"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "queue-group-templates"},
				{Name: "ingress"},
				{Name: "queue-group", Key: map[string]string{"ingress-queue-group-name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "ingress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
				{Name: "queue-overrides"},
				{Name: "queue", Key: map[string]string{"queue-id": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "ingress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
				{Name: "queue-overrides"},
				{Name: "queue", Key: map[string]string{"queue-id": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "ingress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
				{Name: "queue-overrides"},
				{Name: "queue", Key: map[string]string{"queue-id": ""}},
				{Name: "queue-id"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "qos"},
				{Name: "queue-group-templates"},
				{Name: "ingress"},
				{Name: "queue-group", Key: map[string]string{"ingress-queue-group-name": ""}},
				{Name: "queue-group-name]"},
				{Name: "queue", Key: map[string]string{"queue-id": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "ingress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
				{Name: "scheduler-policy"},
				{Name: "overrides"},
				{Name: "scheduler", Key: map[string]string{"scheduler-name": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "ingress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
				{Name: "scheduler-policy"},
				{Name: "overrides"},
				{Name: "scheduler", Key: map[string]string{"scheduler-name": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "access"},
				{Name: "ingress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name": ""}},
				{Name: "scheduler-policy"},
				{Name: "policy-name"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "scheduler-policy", Key: map[string]string{"scheduler-policy-name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "accounting-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "log"},
				{Name: "accounting-policy", Key: map[string]string{"policy-id": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "dampening"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "dampening"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "dot1x"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "dot1x"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "dot1x"},
				{Name: "macsec"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "dot1x"},
				{Name: "macsec"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "dot1x"},
				{Name: "macsec"},
				{Name: "exclude-mac-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "macsec"},
				{Name: "mac-policy", Key: map[string]string{"mac-policy-id": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "dot1x"},
				{Name: "macsec"},
				{Name: "sub-port", Key: map[string]string{"sub-port-id": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "dot1x"},
				{Name: "macsec"},
				{Name: "sub-port", Key: map[string]string{"sub-port-id": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "dot1x"},
				{Name: "macsec"},
				{Name: "sub-port", Key: map[string]string{"sub-port-id": ""}},
				{Name: "ca-name"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "macsec"},
				{Name: "connectivity-association", Key: map[string]string{"ca-name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "dot1x"},
				{Name: "radius-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "system"},
				{Name: "security"},
				{Name: "dot1x"},
				{Name: "radius-policy", Key: map[string]string{"policy-name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "efm-oam"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "efm-oam"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "efm-oam"},
				{Name: "link-monitoring"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "efm-oam"},
				{Name: "link-monitoring"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "expanded-secondary-shaper", Key: map[string]string{"secondary-shaper-name": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "expanded-secondary-shaper", Key: map[string]string{"secondary-shaper-name": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "expanded-secondary-shaper", Key: map[string]string{"secondary-shaper-name": ""}},
				{Name: "class", Key: map[string]string{"class-number": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "expanded-secondary-shaper", Key: map[string]string{"secondary-shaper-name": ""}},
				{Name: "class", Key: map[string]string{"class-number": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "hs-port-pool-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "hs-port-pool-policy", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "hs-scheduler-policy"},
				{Name: "overrides"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "hs-scheduler-policy"},
				{Name: "overrides"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "hs-scheduler-policy"},
				{Name: "overrides"},
				{Name: "group", Key: map[string]string{"group-id": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "hs-scheduler-policy"},
				{Name: "overrides"},
				{Name: "group", Key: map[string]string{"group-id": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "hs-scheduler-policy"},
				{Name: "overrides"},
				{Name: "scheduling-class", Key: map[string]string{"class-number": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "hs-scheduler-policy"},
				{Name: "overrides"},
				{Name: "scheduling-class", Key: map[string]string{"class-number": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "hs-scheduler-policy"},
				{Name: "policy-name"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "hs-scheduler-policy", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "hs-secondary-shaper", Key: map[string]string{"secondary-shaper-name": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "hs-secondary-shaper", Key: map[string]string{"secondary-shaper-name": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "hs-secondary-shaper", Key: map[string]string{"secondary-shaper-name": ""}},
				{Name: "class", Key: map[string]string{"class-number": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "hs-secondary-shaper", Key: map[string]string{"secondary-shaper-name": ""}},
				{Name: "class", Key: map[string]string{"class-number": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "hsmda-scheduler-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "hsmda-scheduler-policy", Key: map[string]string{"hsmda-scheduler-policy-name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "port-qos-policy"},
				{Name: "policy-name"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "port-qos-policy", Key: map[string]string{"port-qos-policy-name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "level", Key: map[string]string{"priority-level": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "level", Key: map[string]string{"priority-level": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "policy-name"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "port-scheduler-policy", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "elmi"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "elmi"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "eth-cfm"},
				{Name: "mep", Key: map[string]string{"md-admin-name ma-admin-name mep-id": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "eth-cfm"},
				{Name: "mep", Key: map[string]string{"md-admin-name ma-admin-name mep-id": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "eth-cfm"},
				{Name: "mep", Key: map[string]string{"md-admin-name ma-admin-name mep-id": ""}},
				{Name: "ma-admin-name"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "eth-cfm"},
				{Name: "domain", Key: map[string]string{"md-admin-name": ""}},
				{Name: "md-admin-name]"},
				{Name: "association", Key: map[string]string{"ma-admin-name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "eth-cfm"},
				{Name: "mep", Key: map[string]string{"md-admin-name ma-admin-name mep-id": ""}},
				{Name: "md-admin-name"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "eth-cfm"},
				{Name: "domain", Key: map[string]string{"md-admin-name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "hsmda-scheduler-overrides"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "hsmda-scheduler-overrides"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "hsmda-scheduler-overrides"},
				{Name: "group", Key: map[string]string{"group-id": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "hsmda-scheduler-overrides"},
				{Name: "group", Key: map[string]string{"group-id": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "hsmda-scheduler-overrides"},
				{Name: "scheduling-class", Key: map[string]string{"class-number": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "hsmda-scheduler-overrides"},
				{Name: "scheduling-class", Key: map[string]string{"class-number": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "lldp"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "lldp"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "lldp"},
				{Name: "dest-mac", Key: map[string]string{"mac-type": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "lldp"},
				{Name: "dest-mac", Key: map[string]string{"mac-type": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "lldp"},
				{Name: "dest-mac", Key: map[string]string{"mac-type": ""}},
				{Name: "tx-mgmt-address", Key: map[string]string{"mgmt-address-system-type": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "lldp"},
				{Name: "dest-mac", Key: map[string]string{"mac-type": ""}},
				{Name: "tx-mgmt-address", Key: map[string]string{"mgmt-address-system-type": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "loopback"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "loopback"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "network"},
				{Name: "accounting-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "log"},
				{Name: "accounting-policy", Key: map[string]string{"policy-id": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "network"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "network"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "network"},
				{Name: "egress"},
				{Name: "port-queues"},
				{Name: "overrides"},
				{Name: "queue", Key: map[string]string{"queue-id": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "network"},
				{Name: "egress"},
				{Name: "port-queues"},
				{Name: "overrides"},
				{Name: "queue", Key: map[string]string{"queue-id": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "network"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "accounting-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "log"},
				{Name: "accounting-policy", Key: map[string]string{"policy-id": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "network"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "network"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "network"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "policer-control-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "policer-control-policy", Key: map[string]string{"policer-control-policy-name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "network"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "queue-group-name"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "queue-group-templates"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"egress-queue-group-name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "network"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "queue-overrides"},
				{Name: "queue", Key: map[string]string{"queue-id": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "network"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "queue-overrides"},
				{Name: "queue", Key: map[string]string{"queue-id": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "network"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "queue-overrides"},
				{Name: "queue", Key: map[string]string{"queue-id": ""}},
				{Name: "queue-id"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "qos"},
				{Name: "queue-group-templates"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"egress-queue-group-name": ""}},
				{Name: "queue-group-name]"},
				{Name: "queue", Key: map[string]string{"queue-id": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "network"},
				{Name: "egress"},
				{Name: "queue-group", Key: map[string]string{"queue-group-name instance-id": ""}},
				{Name: "scheduler-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "scheduler-policy", Key: map[string]string{"scheduler-policy-name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "network"},
				{Name: "egress"},
				{Name: "queue-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "network-queue", Key: map[string]string{"network-queue-policy": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "ssm"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "ssm"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "symbol-monitor"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "ethernet"},
				{Name: "symbol-monitor"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "gnss"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "gnss"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "network"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "network"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "network"},
				{Name: "egress"},
				{Name: "pool", Key: map[string]string{"name": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "network"},
				{Name: "egress"},
				{Name: "pool", Key: map[string]string{"name": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "network"},
				{Name: "egress"},
				{Name: "pool", Key: map[string]string{"name": ""}},
				{Name: "slope-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "slope-policy", Key: map[string]string{"slope-policy-name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "otu"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "otu"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "sonet-sdh"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "sonet-sdh"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "sonet-sdh"},
				{Name: "group", Key: map[string]string{"group-index": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "sonet-sdh"},
				{Name: "group", Key: map[string]string{"group-index": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "sonet-sdh"},
				{Name: "path", Key: map[string]string{"path-index": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "sonet-sdh"},
				{Name: "path", Key: map[string]string{"path-index": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "sonet-sdh"},
				{Name: "path", Key: map[string]string{"path-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "sonet-sdh"},
				{Name: "path", Key: map[string]string{"path-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "sonet-sdh"},
				{Name: "path", Key: map[string]string{"path-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "level", Key: map[string]string{"priority-level": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "sonet-sdh"},
				{Name: "path", Key: map[string]string{"path-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "level", Key: map[string]string{"priority-level": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "sonet-sdh"},
				{Name: "path", Key: map[string]string{"path-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "policy-name"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "port-scheduler-policy", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "sonet-sdh"},
				{Name: "path", Key: map[string]string{"path-index": ""}},
				{Name: "network"},
				{Name: "accounting-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "log"},
				{Name: "accounting-policy", Key: map[string]string{"policy-id": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "sonet-sdh"},
				{Name: "path", Key: map[string]string{"path-index": ""}},
				{Name: "network"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "sonet-sdh"},
				{Name: "path", Key: map[string]string{"path-index": ""}},
				{Name: "network"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "sonet-sdh"},
				{Name: "path", Key: map[string]string{"path-index": ""}},
				{Name: "network"},
				{Name: "queue-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "network-queue", Key: map[string]string{"network-queue-policy": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "sonet-sdh"},
				{Name: "path", Key: map[string]string{"path-index": ""}},
				{Name: "ppp"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "sonet-sdh"},
				{Name: "path", Key: map[string]string{"path-index": ""}},
				{Name: "ppp"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "level", Key: map[string]string{"priority-level": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "level", Key: map[string]string{"priority-level": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "policy-name"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "port-scheduler-policy", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "network"},
				{Name: "accounting-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "log"},
				{Name: "accounting-policy", Key: map[string]string{"policy-id": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "network"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "network"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "network"},
				{Name: "queue-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "network-queue", Key: map[string]string{"network-queue-policy": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "ppp"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds1", Key: map[string]string{"ds1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "ppp"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "level", Key: map[string]string{"priority-level": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "level", Key: map[string]string{"priority-level": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "policy-name"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "port-scheduler-policy", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
				{Name: "network"},
				{Name: "accounting-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "log"},
				{Name: "accounting-policy", Key: map[string]string{"policy-id": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
				{Name: "network"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
				{Name: "network"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
				{Name: "network"},
				{Name: "queue-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "network-queue", Key: map[string]string{"network-queue-policy": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
				{Name: "ppp"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "ds3", Key: map[string]string{"ds3-index": ""}},
				{Name: "ppp"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e1", Key: map[string]string{"e1-index": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e1", Key: map[string]string{"e1-index": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e1", Key: map[string]string{"e1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e1", Key: map[string]string{"e1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e1", Key: map[string]string{"e1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e1", Key: map[string]string{"e1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e1", Key: map[string]string{"e1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "level", Key: map[string]string{"priority-level": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e1", Key: map[string]string{"e1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "level", Key: map[string]string{"priority-level": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e1", Key: map[string]string{"e1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "policy-name"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "port-scheduler-policy", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e1", Key: map[string]string{"e1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "network"},
				{Name: "accounting-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "log"},
				{Name: "accounting-policy", Key: map[string]string{"policy-id": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e1", Key: map[string]string{"e1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "network"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e1", Key: map[string]string{"e1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "network"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e1", Key: map[string]string{"e1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "network"},
				{Name: "queue-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "network-queue", Key: map[string]string{"network-queue-policy": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e1", Key: map[string]string{"e1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "ppp"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e1", Key: map[string]string{"e1-index": ""}},
				{Name: "channel-group", Key: map[string]string{"ds0-index": ""}},
				{Name: "ppp"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e3", Key: map[string]string{"e3-index": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e3", Key: map[string]string{"e3-index": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e3", Key: map[string]string{"e3-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e3", Key: map[string]string{"e3-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e3", Key: map[string]string{"e3-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "level", Key: map[string]string{"priority-level": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e3", Key: map[string]string{"e3-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "overrides"},
				{Name: "level", Key: map[string]string{"priority-level": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e3", Key: map[string]string{"e3-index": ""}},
				{Name: "egress"},
				{Name: "port-scheduler-policy"},
				{Name: "policy-name"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "port-scheduler-policy", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e3", Key: map[string]string{"e3-index": ""}},
				{Name: "network"},
				{Name: "accounting-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "log"},
				{Name: "accounting-policy", Key: map[string]string{"policy-id": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e3", Key: map[string]string{"e3-index": ""}},
				{Name: "network"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e3", Key: map[string]string{"e3-index": ""}},
				{Name: "network"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e3", Key: map[string]string{"e3-index": ""}},
				{Name: "network"},
				{Name: "queue-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "network-queue", Key: map[string]string{"network-queue-policy": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e3", Key: map[string]string{"e3-index": ""}},
				{Name: "ppp"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "tdm"},
				{Name: "e3", Key: map[string]string{"e3-index": ""}},
				{Name: "ppp"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "transceiver"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port"},
				{Name: "transceiver"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port-policy"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port-policy"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port-policy"},
				{Name: "egress-port-scheduler-policy"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "qos"},
				{Name: "port-scheduler-policy", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port-xc"},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port-xc"},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port-xc"},
				{Name: "pxc", Key: map[string]string{"pxc-id": ""}},
				{Name: "apply-groups"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port-xc"},
				{Name: "pxc", Key: map[string]string{"pxc-id": ""}},
				{Name: "apply-groups-exclude"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "groups"},
				{Name: "group", Key: map[string]string{"name": ""}},
			},
		},
	},
	{
		LocalPath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "port-xc"},
				{Name: "pxc", Key: map[string]string{"pxc-id": ""}},
				{Name: "port-id"},
			},
		},
		RemotePath: &gnmi.Path{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "port", Key: map[string]string{"port-id": ""}},
			},
		},
	},
}

// SetupConfigurePort adds a controller that reconciles ConfigurePorts.
func SetupConfigurePort(mgr ctrl.Manager, o controller.Options, l logging.Logger, poll time.Duration, namespace string) (string, chan cevent.GenericEvent, error) {

	name := managed.ControllerName(srosv1alpha1.ConfigurePortGroupKind)

	events := make(chan cevent.GenericEvent)

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(srosv1alpha1.ConfigurePortGroupVersionKind),
		managed.WithExternalConnecter(&connectorConfigurePort{
			log:         l,
			kube:        mgr.GetClient(),
			usage:       resource.NewNetworkNodeUsageTracker(mgr.GetClient(), &ndrv1.NetworkNodeUsage{}),
			newClientFn: target.NewTarget},
		),
		managed.WithParser(l),
		managed.WithValidator(&validatorConfigurePort{log: l, parser: *parser.NewParser(parser.WithLogger(l))}),
		managed.WithLogger(l.WithValues("controller", name)),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))))

	return srosv1alpha1.ConfigurePortGroupKind, events, ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o).
		For(&srosv1alpha1.SrosConfigurePort{}).
		WithEventFilter(resource.IgnoreUpdateWithoutGenerationChangePredicate()).
		Watches(
			&source.Channel{Source: events},
			&handler.EnqueueRequestForObject{},
		).
		//Watches(
		//	&source.Kind{Type: &ndrv1.NetworkNode{}},
		//	handler.EnqueueRequestsFromMapFunc(r.NetworkNodeMapFunc),
		//).
		Complete(r)
}

type validatorConfigurePort struct {
	log    logging.Logger
	parser parser.Parser
}

func (v *validatorConfigurePort) ValidateLocalleafRef(ctx context.Context, mg resource.Managed) (managed.ValidateLocalleafRefObservation, error) {
	log := v.log.WithValues("resource", mg.GetName())
	log.Debug("ValidateLocalleafRef...")

	// json unmarshal the resource
	o, ok := mg.(*srosv1alpha1.SrosConfigurePort)
	if !ok {
		return managed.ValidateLocalleafRefObservation{}, errors.New(errUnexpectedConfigurePort)
	}
	d, err := json.Marshal(&o.Spec.ForNetworkNode)
	if err != nil {
		return managed.ValidateLocalleafRefObservation{}, errors.Wrap(err, errJSONMarshal)
	}
	var x1 interface{}
	json.Unmarshal(d, &x1)

	// For local leafref validation we dont need to supply the external data so we use nil
	success, resultleafRefValidation, err := v.parser.ValidateLeafRefGnmi(
		parser.LeafRefValidationLocal, x1, nil, localleafRefConfigurePort, log)
	if err != nil {
		return managed.ValidateLocalleafRefObservation{
			Success: false,
		}, nil
	}
	if !success {
		log.Debug("ValidateLocalleafRef failed", "resultleafRefValidation", resultleafRefValidation)
		return managed.ValidateLocalleafRefObservation{
			Success:          false,
			ResolvedLeafRefs: resultleafRefValidation}, nil
	}
	log.Debug("ValidateLocalleafRef success", "resultleafRefValidation", resultleafRefValidation)
	return managed.ValidateLocalleafRefObservation{
		Success:          true,
		ResolvedLeafRefs: resultleafRefValidation}, nil
}

func (v *validatorConfigurePort) ValidateExternalleafRef(ctx context.Context, mg resource.Managed, cfg []byte) (managed.ValidateExternalleafRefObservation, error) {
	log := v.log.WithValues("resource", mg.GetName())
	log.Debug("ValidateExternalleafRef...")

	// json unmarshal the resource
	o, ok := mg.(*srosv1alpha1.SrosConfigurePort)
	if !ok {
		return managed.ValidateExternalleafRefObservation{}, errors.New(errUnexpectedConfigurePort)
	}
	d, err := json.Marshal(&o.Spec.ForNetworkNode)
	if err != nil {
		return managed.ValidateExternalleafRefObservation{}, errors.Wrap(err, errJSONMarshal)
	}
	var x1 interface{}
	json.Unmarshal(d, &x1)

	// json unmarshal the external data
	var x2 interface{}
	json.Unmarshal(cfg, &x2)

	// For local external leafref validation we need to supply the external
	// data to validate the remote leafref, we use x2 for this
	success, resultleafRefValidation, err := v.parser.ValidateLeafRefGnmi(
		parser.LeafRefValidationExternal, x1, x2, externalLeafRefConfigurePort, log)
	if err != nil {
		return managed.ValidateExternalleafRefObservation{
			Success: false,
		}, nil
	}
	if !success {
		log.Debug("ValidateExternalleafRef failed", "resultleafRefValidation", resultleafRefValidation)
		return managed.ValidateExternalleafRefObservation{
			Success:          false,
			ResolvedLeafRefs: resultleafRefValidation}, nil
	}
	log.Debug("ValidateExternalleafRef success", "resultleafRefValidation", resultleafRefValidation)
	return managed.ValidateExternalleafRefObservation{
		Success:          true,
		ResolvedLeafRefs: resultleafRefValidation}, nil
}

func (v *validatorConfigurePort) ValidateParentDependency(ctx context.Context, mg resource.Managed, cfg []byte) (managed.ValidateParentDependencyObservation, error) {
	log := v.log.WithValues("resource", mg.GetName())
	log.Debug("ValidateParentDependency...")

	// we initialize a global list for finer information on the resolution
	resultleafRefValidation := make([]*parser.ResolvedLeafRefGnmi, 0)
	log.Debug("ValidateParentDependency success", "resultParentValidation", resultleafRefValidation)
	return managed.ValidateParentDependencyObservation{
		Success:          true,
		ResolvedLeafRefs: resultleafRefValidation}, nil
}

// ValidateResourceIndexes validates if the indexes of a resource got changed
// if so we need to delete the original resource, because it will be dangling if we dont delete it
func (v *validatorConfigurePort) ValidateResourceIndexes(ctx context.Context, mg resource.Managed) (managed.ValidateResourceIndexesObservation, error) {
	log := v.log.WithValues("resosurce", mg.GetName())

	// json unmarshal the resource
	o, ok := mg.(*srosv1alpha1.SrosConfigurePort)
	if !ok {
		return managed.ValidateResourceIndexesObservation{}, errors.New(errUnexpectedConfigurePort)
	}
	log.Debug("ValidateResourceIndexes", "Spec", o.Spec)

	rootPath := []*gnmi.Path{
		{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "port"},
			},
		},
	}

	origResourceIndex := mg.GetResourceIndexes()
	// we call the CompareConfigPathsWithResourceKeys irrespective is the get resource index returns nil
	changed, deletPaths, newResourceIndex := v.parser.CompareGnmiPathsWithResourceKeys(rootPath[0], origResourceIndex)
	if changed {
		log.Debug("ValidateResourceIndexes changed", "deletPaths", deletPaths[0])
		return managed.ValidateResourceIndexesObservation{Changed: true, ResourceDeletes: deletPaths, ResourceIndexes: newResourceIndex}, nil
	}

	log.Debug("ValidateResourceIndexes success")
	return managed.ValidateResourceIndexesObservation{Changed: false, ResourceIndexes: newResourceIndex}, nil
}

// A connector is expected to produce an ExternalClient when its Connect method
// is called.
type connectorConfigurePort struct {
	log         logging.Logger
	kube        client.Client
	usage       resource.Tracker
	newClientFn func(c *gnmitypes.TargetConfig) *target.Target
	//newClientFn func(ctx context.Context, cfg ndd.Config) (config.ConfigurationClient, error)
}

// Connect produces an ExternalClient by:
// 1. Tracking that the managed resource is using a NetworkNode.
// 2. Getting the managed resource's NetworkNode with connection details
// A resource is mapped to a single target
func (c *connectorConfigurePort) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	log := c.log.WithValues("resource", mg.GetName())
	log.Debug("Connect")
	o, ok := mg.(*srosv1alpha1.SrosConfigurePort)
	if !ok {
		return nil, errors.New(errUnexpectedConfigurePort)
	}
	if err := c.usage.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackTCUsage)
	}

	// find network node that is configured status
	nn := &ndrv1.NetworkNode{}
	if err := c.kube.Get(ctx, types.NamespacedName{Name: o.GetNetworkNodeReference().Name}, nn); err != nil {
		return nil, errors.Wrap(err, errGetNetworkNode)
	}

	if nn.GetCondition(ndrv1.ConditionKindDeviceDriverConfigured).Status != corev1.ConditionTrue {
		return nil, errors.New(targetNotConfigured)
	}
	cfg := &gnmitypes.TargetConfig{
		Name:       nn.GetName(),
		Address:    ndrv1.PrefixService + "-" + nn.Name + "." + ndrv1.NamespaceLocalK8sDNS + strconv.Itoa(*nn.Spec.GrpcServerPort),
		Username:   utils.StringPtr("admin"),
		Password:   utils.StringPtr("admin"),
		Timeout:    10 * time.Second,
		SkipVerify: utils.BoolPtr(true),
		Insecure:   utils.BoolPtr(true),
		TLSCA:      utils.StringPtr(""), //TODO TLS
		TLSCert:    utils.StringPtr(""), //TODO TLS
		TLSKey:     utils.StringPtr(""),
		Gzip:       utils.BoolPtr(false),
	}

	cl := target.NewTarget(cfg)
	if err := cl.CreateGNMIClient(ctx); err != nil {
		return nil, errors.Wrap(err, errNewClient)
	}

	// we make a string here since we use a trick in registration to go to multiple targets
	// while here the object is mapped to a single target/network node
	tns := make([]string, 0)
	tns = append(tns, nn.GetName())

	return &externalConfigurePort{client: cl, targets: tns, log: log, parser: *parser.NewParser(parser.WithLogger(log))}, nil
}

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type externalConfigurePort struct {
	//client  config.ConfigurationClient
	client  *target.Target
	targets []string
	log     logging.Logger
	parser  parser.Parser
}

func (e *externalConfigurePort) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	o, ok := mg.(*srosv1alpha1.SrosConfigurePort)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errUnexpectedConfigurePort)
	}
	log := e.log.WithValues("Resource", o.GetName())
	log.Debug("Observing ...")

	// rootpath of the resource
	rootPath := []*gnmi.Path{
		{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "port"},
			},
		},
	}

	// gvk: group, version, kind, name, namespace of the resource
	gvk := &gvk.GVK{
		Group:     mg.GetObjectKind().GroupVersionKind().Group,
		Version:   mg.GetObjectKind().GroupVersionKind().Version,
		Kind:      mg.GetObjectKind().GroupVersionKind().Kind,
		Name:      mg.GetName(),
		NameSpace: mg.GetNamespace(),
	}
	gvkstring, err := gvk.String()
	if err != nil {
		return managed.ExternalObservation{}, err
	}

	// gext: gni extension information for the resource: action, gvk name and level
	gextInfo := &gext.GEXT{
		Action: gext.GEXTActionGet,
		Name:   gvkstring,
		Level:  levelConfigurePort,
	}
	gextInfoString, err := gextInfo.String()
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errGetGextInfo)
	}

	// gnmi get request
	req := &gnmi.GetRequest{
		Path:     rootPath,
		Encoding: gnmi.Encoding_JSON,
		Extension: []*gnmi_ext.Extension{
			{Ext: &gnmi_ext.Extension_RegisteredExt{
				RegisteredExt: &gnmi_ext.RegisteredExtension{Id: gnmi_ext.ExtensionID_EID_EXPERIMENTAL, Msg: []byte(gextInfoString)}}},
		},
	}

	// gnmi get response
	resp, err := e.client.Get(ctx, req)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errReadConfigurePort)
	}

	// validate if the extension matches or not
	if resp.GetExtension()[0].GetRegisteredExt().GetId() != gnmi_ext.ExtensionID_EID_EXPERIMENTAL {
		log.Debug("Observe response GNMI Extension mismatch", "Extension Info", resp.GetExtension()[0])
		return managed.ExternalObservation{}, errors.New(errGnmiExtensionMismatch)
	}

	// get gnmi extension metadata
	meta := resp.GetExtension()[0].GetRegisteredExt().GetMsg()
	respMeta := &gext.GEXT{}
	if err := json.Unmarshal(meta, &respMeta); err != nil {
		log.Debug("Observe response gext unmarshal issue", "Extension Info", meta)
		return managed.ExternalObservation{}, errors.Wrap(err, errJSONMarshal)
	}

	// prepare the input data to compare against the response data
	d, err := json.Marshal(&o.Spec.ForNetworkNode)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errJSONMarshal)
	}
	var x1 interface{}
	if err := json.Unmarshal(d, &x1); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errJSONUnMarshal)
	}

	// remove the hierarchical elements for data processing, comparison, etc
	// they are used in the provider for parent dependency resolution
	// but are not relevant in the data, they are referenced in the rootPath
	// when interacting with the device driver
	hids := make([]string, 0)
	x1 = e.parser.RemoveLeafsFromJSONData(x1, hids)

	// validate gnmi resp information
	var x2 interface{}
	if len(resp.GetNotification()) != 0 {
		if len(resp.GetNotification()[0].GetUpdate()) != 0 {
			// get value from gnmi get response
			x2, err = e.parser.GetValue(resp.GetNotification()[0].GetUpdate()[0].Val)
			if err != nil {
				log.Debug("Observe response get value issue")
				return managed.ExternalObservation{}, errors.Wrap(err, errJSONMarshal)
			}
		}
	}

	// logging information that will be used to provide the response
	log.Debug("Observer Response", "Meta", string(meta))
	log.Debug("Spec Data", "X1", x1)
	log.Debug("Resp Data", "X2", x2)

	// if the cache is not ready we back off and return
	if !respMeta.CacheReady {
		log.Debug("Cache Not Ready ...")
		return managed.ExternalObservation{
			Ready:            false,
			ResourceExists:   false,
			ResourceHasData:  true,
			ResourceUpToDate: false,
		}, nil
	}

	if !respMeta.Exists {
		// Resource Does not Exists
		if respMeta.HasData {
			// this is an umnaged resource which has data and will be moved to a managed resource

			updatesx1 := e.parser.GetUpdatesFromJSONDataGnmi(rootPath[0], e.parser.XpathToGnmiPath("/", 0), x1, resourceRefPathsConfigurePort)
			for _, update := range updatesx1 {
				log.Debug("Observe Fine Grane Updates X1", "Path", e.parser.GnmiPathToXPath(update.Path, true), "Value", update.GetVal())
			}
			// for lists with keys we need to create a list before calulating the paths since this is what
			// the object eventually happens to be based upon. We avoid having multiple entries in a list object
			// and hence we have to add this step
			x2, err = e.parser.AddJSONDataToList(x2)
			if err != nil {
				return managed.ExternalObservation{}, errors.Wrap(err, errWrongInputdata)
			}
			updatesx2 := e.parser.GetUpdatesFromJSONDataGnmi(rootPath[0], e.parser.XpathToGnmiPath("/", 0), x2, resourceRefPathsConfigurePort)
			for _, update := range updatesx2 {
				log.Debug("Observe Fine Grane Updates X2", "Path", e.parser.GnmiPathToXPath(update.Path, true), "Value", update.GetVal())
			}

			deletes, updates, err := e.parser.FindResourceDeltaGnmi(updatesx1, updatesx2, log)
			if err != nil {
				return managed.ExternalObservation{}, err
			}
			if len(deletes) != 0 || len(updates) != 0 {
				// UMR -> MR with data, which is NOT up to date
				log.Debug("Observing Response: resource NOT up to date", "Exists", false, "HasData", true, "UpToDate", false, "Response", resp, "Updates", updates, "Deletes", deletes)
				for _, del := range deletes {
					log.Debug("Observing Response: resource NOT up to date, deletes", "path", e.parser.GnmiPathToXPath(del, true))
				}
				for _, upd := range updates {
					val, _ := e.parser.GetValue(upd.GetVal())
					log.Debug("Observing Response: resource NOT up to date, updates", "path", e.parser.GnmiPathToXPath(upd.GetPath(), true), "data", val)
				}
				return managed.ExternalObservation{
					Ready:            true,
					ResourceExists:   false,
					ResourceHasData:  true,
					ResourceUpToDate: false,
					ResourceDeletes:  deletes,
					ResourceUpdates:  updates,
				}, nil
			}
			// UMR -> MR with data, which is up to date
			log.Debug("Observing Response: resource up to date", "Exists", false, "HasData", true, "UpToDate", true, "Response", resp)
			return managed.ExternalObservation{
				Ready:            true,
				ResourceExists:   false,
				ResourceHasData:  true,
				ResourceUpToDate: true,
			}, nil
		}
		// UMR -> MR without data
		log.Debug("Observing Response:", "Exists", false, "HasData", false, "UpToDate", false, "Response", resp)
		return managed.ExternalObservation{
			Ready:            true,
			ResourceExists:   false,
			ResourceHasData:  false,
			ResourceUpToDate: false,
		}, nil

	}
	// Resource Exists
	switch respMeta.Status {
	case gext.ResourceStatusSuccess:
		if respMeta.HasData {
			// data is present

			updatesx1 := e.parser.GetUpdatesFromJSONDataGnmi(rootPath[0], e.parser.XpathToGnmiPath("/", 0), x1, resourceRefPathsConfigurePort)
			for _, update := range updatesx1 {
				log.Debug("Observe Fine Grane Updates X1", "Path", e.parser.GnmiPathToXPath(update.Path, true), "Value", update.GetVal())
			}
			updatesx2 := e.parser.GetUpdatesFromJSONDataGnmi(rootPath[0], e.parser.XpathToGnmiPath("/", 0), x2, resourceRefPathsConfigurePort)
			for _, update := range updatesx2 {
				log.Debug("Observe Fine Grane Updates X2", "Path", e.parser.GnmiPathToXPath(update.Path, true), "Value", update.GetVal())
			}

			deletes, updates, err := e.parser.FindResourceDeltaGnmi(updatesx1, updatesx2, log)
			if err != nil {
				return managed.ExternalObservation{}, err
			}
			// MR -> MR, resource is NOT up to date
			if len(deletes) != 0 || len(updates) != 0 {
				// resource is NOT up to date
				log.Debug("Observing Response: resource NOT up to date", "Exists", true, "HasData", true, "UpToDate", false, "Response", resp, "Updates", updates, "Deletes", deletes)
				for _, del := range deletes {
					log.Debug("Observing Response: resource NOT up to date, deletes", "path", e.parser.GnmiPathToXPath(del, true))
				}
				for _, upd := range updates {
					val, _ := e.parser.GetValue(upd.GetVal())
					log.Debug("Observing Response: resource NOT up to date, updates", "path", e.parser.GnmiPathToXPath(upd.GetPath(), true), "data", val)
				}
				return managed.ExternalObservation{
					Ready:            true,
					ResourceExists:   true,
					ResourceHasData:  true,
					ResourceUpToDate: false,
					ResourceDeletes:  deletes,
					ResourceUpdates:  updates,
				}, nil
			}
			// MR -> MR, resource is up to date
			log.Debug("Observing Response: resource up to date", "Exists", true, "HasData", true, "UpToDate", true, "Response", resp)
			return managed.ExternalObservation{
				Ready:            true,
				ResourceExists:   true,
				ResourceHasData:  true,
				ResourceUpToDate: true,
			}, nil
		}
		// MR -> MR, resource has no data, strange, someone could have deleted the resource
		log.Debug("Observing Response", "Exists", true, "HasData", false, "UpToDate", false, "Status", respMeta.Status)
		return managed.ExternalObservation{
			Ready:            true,
			ResourceExists:   true,
			ResourceHasData:  false,
			ResourceUpToDate: false,
		}, nil

	default:
		// MR -> MR, resource is not in a success state, so the object might still be in creation phase
		log.Debug("Observing Response", "Exists", true, "HasData", false, "UpToDate", false, "Status", respMeta.Status)
		return managed.ExternalObservation{
			Ready:            true,
			ResourceExists:   true,
			ResourceHasData:  false,
			ResourceUpToDate: false,
		}, nil
	}
}

func (e *externalConfigurePort) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	o, ok := mg.(*srosv1alpha1.SrosConfigurePort)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errUnexpectedConfigurePort)
	}
	log := e.log.WithValues("Resource", o.GetName())
	log.Debug("Creating ...")

	rootPath := []*gnmi.Path{
		{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "port"},
			},
		},
	}

	d, err := json.Marshal(&o.Spec.ForNetworkNode)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errJSONMarshal)
	}

	var x1 interface{}
	if err := json.Unmarshal(d, &x1); err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errJSONUnMarshal)
	}

	// remove the hierarchical elements for data processing, comparison, etc
	// they are used in the provider for parent dependency resolution
	// but are not relevant in the data, they are referenced in the rootPath
	// when interacting with the device driver
	hids := make([]string, 0)
	x1 = e.parser.RemoveLeafsFromJSONData(x1, hids)

	updates := e.parser.GetUpdatesFromJSONDataGnmi(rootPath[0], e.parser.XpathToGnmiPath("/", 0), x1, resourceRefPathsConfigurePort)
	for _, update := range updates {
		log.Debug("Create Fine Grane Updates", "Path", update.Path, "Value", update.GetVal())
	}

	gvk := &gvk.GVK{
		Group:     mg.GetObjectKind().GroupVersionKind().Group,
		Version:   mg.GetObjectKind().GroupVersionKind().Version,
		Kind:      mg.GetObjectKind().GroupVersionKind().Kind,
		Name:      mg.GetName(),
		NameSpace: mg.GetNamespace(),
	}
	gvkstring, err := gvk.String()
	if err != nil {
		return managed.ExternalCreation{}, err
	}

	gextInfo := &gext.GEXT{
		Action:   gext.GEXTActionCreate,
		Name:     gvkstring,
		Level:    levelConfigurePort,
		RootPath: rootPath[0],
	}
	gextInfoString, err := gextInfo.String()
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errGetGextInfo)
	}

	if len(updates) == 0 {
		log.Debug("cannot create object since there are no updates present")
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateObject)
	}

	req := &gnmi.SetRequest{
		Replace: updates,
		Extension: []*gnmi_ext.Extension{
			{Ext: &gnmi_ext.Extension_RegisteredExt{
				RegisteredExt: &gnmi_ext.RegisteredExtension{Id: gnmi_ext.ExtensionID_EID_EXPERIMENTAL, Msg: []byte(gextInfoString)}}},
		},
	}

	_, err = e.client.Set(ctx, req)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errReadConfigurePort)
	}

	return managed.ExternalCreation{}, nil
}

func (e *externalConfigurePort) Update(ctx context.Context, mg resource.Managed, obs managed.ExternalObservation) (managed.ExternalUpdate, error) {
	o, ok := mg.(*srosv1alpha1.SrosConfigurePort)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errUnexpectedConfigurePort)
	}
	log := e.log.WithValues("Resource", o.GetName())
	log.Debug("Updating ...")

	for _, u := range obs.ResourceUpdates {
		log.Debug("Update -> Update", "Path", u.Path, "Value", u.GetVal())
	}
	for _, d := range obs.ResourceDeletes {
		log.Debug("Update -> Delete", "Path", d)
	}

	gvk := &gvk.GVK{
		Group:     mg.GetObjectKind().GroupVersionKind().Group,
		Version:   mg.GetObjectKind().GroupVersionKind().Version,
		Kind:      mg.GetObjectKind().GroupVersionKind().Kind,
		Name:      mg.GetName(),
		NameSpace: mg.GetNamespace(),
	}
	gvkstring, err := gvk.String()
	if err != nil {
		return managed.ExternalUpdate{}, err
	}

	gextInfo := &gext.GEXT{
		Action: gext.GEXTActionUpdate,
		Name:   gvkstring,
		Level:  levelConfigurePort,
	}
	gextInfoString, err := gextInfo.String()
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errGetGextInfo)
	}

	req := &gnmi.SetRequest{
		Update: obs.ResourceUpdates,
		Delete: obs.ResourceDeletes,
		Extension: []*gnmi_ext.Extension{
			{Ext: &gnmi_ext.Extension_RegisteredExt{
				RegisteredExt: &gnmi_ext.RegisteredExtension{Id: gnmi_ext.ExtensionID_EID_EXPERIMENTAL, Msg: []byte(gextInfoString)}}},
		},
	}

	_, err = e.client.Set(ctx, req)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errReadConfigurePort)
	}

	return managed.ExternalUpdate{}, nil
}

func (e *externalConfigurePort) Delete(ctx context.Context, mg resource.Managed) error {
	o, ok := mg.(*srosv1alpha1.SrosConfigurePort)
	if !ok {
		return errors.New(errUnexpectedConfigurePort)
	}
	log := e.log.WithValues("Resource", o.GetName())
	log.Debug("Deleting ...")

	rootPath := []*gnmi.Path{
		{
			Elem: []*gnmi.PathElem{
				{Name: "configure"},
				{Name: "port"},
			},
		},
	}

	gvk := &gvk.GVK{
		Group:     mg.GetObjectKind().GroupVersionKind().Group,
		Version:   mg.GetObjectKind().GroupVersionKind().Version,
		Kind:      mg.GetObjectKind().GroupVersionKind().Kind,
		Name:      mg.GetName(),
		NameSpace: mg.GetNamespace(),
	}
	gvkstring, err := gvk.String()
	if err != nil {
		return err
	}

	gextInfo := &gext.GEXT{
		Action: gext.GEXTActionDelete,
		Name:   gvkstring,
		Level:  levelConfigurePort,
	}
	gextInfoString, err := gextInfo.String()
	if err != nil {
		return errors.Wrap(err, errGetGextInfo)
	}

	req := gnmi.SetRequest{
		Delete: rootPath,
		Extension: []*gnmi_ext.Extension{
			{Ext: &gnmi_ext.Extension_RegisteredExt{
				RegisteredExt: &gnmi_ext.RegisteredExtension{Id: gnmi_ext.ExtensionID_EID_EXPERIMENTAL, Msg: []byte(gextInfoString)}}},
		},
	}

	_, err = e.client.Set(ctx, &req)
	if err != nil {
		return errors.Wrap(err, errDeleteConfigurePort)
	}

	return nil
}

func (e *externalConfigurePort) GetTarget() []string {
	return e.targets
}

func (e *externalConfigurePort) GetConfig(ctx context.Context) ([]byte, error) {
	e.log.Debug("Get Config ...")
	req := &gnmi.GetRequest{
		Path:     []*gnmi.Path{},
		Encoding: gnmi.Encoding_JSON,
	}

	resp, err := e.client.Get(ctx, req)
	if err != nil {
		return make([]byte, 0), errors.Wrap(err, errGetConfig)
	}

	if len(resp.GetNotification()) != 0 {
		if len(resp.GetNotification()[0].GetUpdate()) != 0 {
			x2, err := e.parser.GetValue(resp.GetNotification()[0].GetUpdate()[0].Val)
			if err != nil {
				return make([]byte, 0), errors.Wrap(err, errGetConfig)
			}

			data, err := json.Marshal(x2)
			if err != nil {
				return make([]byte, 0), errors.Wrap(err, errJSONMarshal)
			}
			return data, nil
		}
	}
	e.log.Debug("Get Config Empty response")
	return nil, nil
}

func (e *externalConfigurePort) GetResourceName(ctx context.Context, path []*gnmi.Path) (string, error) {
	e.log.Debug("Get ResourceName ...")

	gextInfo := &gext.GEXT{
		Action: gext.GEXTActionGetResourceName,
	}
	gextInfoString, err := gextInfo.String()
	if err != nil {
		return "", errors.Wrap(err, errGetGextInfo)
	}

	req := &gnmi.GetRequest{
		Path:     path,
		Encoding: gnmi.Encoding_JSON,
		Extension: []*gnmi_ext.Extension{
			{Ext: &gnmi_ext.Extension_RegisteredExt{
				RegisteredExt: &gnmi_ext.RegisteredExtension{Id: gnmi_ext.ExtensionID_EID_EXPERIMENTAL, Msg: []byte(gextInfoString)}}},
		},
	}

	resp, err := e.client.Get(ctx, req)
	if err != nil {
		return "", errors.Wrap(err, errGetResourceName)
	}

	x2, err := e.parser.GetValue(resp.GetNotification()[0].GetUpdate()[0].Val)
	if err != nil {
		return "", errors.Wrap(err, errJSONMarshal)
	}

	d, err := json.Marshal(x2)
	if err != nil {
		return "", errors.Wrap(err, errJSONMarshal)
	}

	var resourceName nddv1.ResourceName
	if err := json.Unmarshal(d, &resourceName); err != nil {
		return "", errors.Wrap(err, errJSONUnMarshal)
	}

	e.log.Debug("Get ResourceName Response", "ResourceName", resourceName)

	return resourceName.Name, nil
}

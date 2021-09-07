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

package v1alpha1

import (
	"reflect"

	nddv1 "github.com/yndd/ndd-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// ConfigurePortFinalizer is the name of the finalizer added to
	// ConfigurePort to block delete operations until the physical node can be
	// deprovisioned.
	ConfigurePortFinalizer string = "port.sros.ndd.yndd.io"
)

// ConfigurePort struct
type ConfigurePort struct {
	Access             *ConfigurePortAccess    `json:"access,omitempty"`
	AdminState         *string                 `json:"admin-state,omitempty"`
	ApplyGroups        *string                 `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                 `json:"apply-groups-exclude,omitempty"`
	Connector          *ConfigurePortConnector `json:"connector,omitempty"`
	// +kubebuilder:default:=true
	DdmEvents              *bool                                `json:"ddm-events,omitempty"`
	Description            *string                              `json:"description,omitempty"`
	DistCpuProtection      *ConfigurePortDistCpuProtection      `json:"dist-cpu-protection,omitempty"`
	Dwdm                   *ConfigurePortDwdm                   `json:"dwdm,omitempty"`
	Ethernet               *ConfigurePortEthernet               `json:"ethernet,omitempty"`
	Gnss                   *ConfigurePortGnss                   `json:"gnss,omitempty"`
	HybridBufferAllocation *ConfigurePortHybridBufferAllocation `json:"hybrid-buffer-allocation,omitempty"`
	ModifyBufferAllocation *ConfigurePortModifyBufferAllocation `json:"modify-buffer-allocation,omitempty"`
	// +kubebuilder:default:=false
	MonitorAggEgressQueueStats *bool                     `json:"monitor-agg-egress-queue-stats,omitempty"`
	Network                    *ConfigurePortNetwork     `json:"network,omitempty"`
	Otu                        *ConfigurePortOtu         `json:"otu,omitempty"`
	PortId                     *string                   `json:"port-id,omitempty"`
	SonetSdh                   *ConfigurePortSonetSdh    `json:"sonet-sdh,omitempty"`
	Tdm                        *ConfigurePortTdm         `json:"tdm,omitempty"`
	Transceiver                *ConfigurePortTransceiver `json:"transceiver,omitempty"`
}

// ConfigurePortAccess struct
type ConfigurePortAccess struct {
	ApplyGroups        *string                     `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                     `json:"apply-groups-exclude,omitempty"`
	Egress             *ConfigurePortAccessEgress  `json:"egress,omitempty"`
	Ingress            *ConfigurePortAccessIngress `json:"ingress,omitempty"`
}

// ConfigurePortAccessEgress struct
type ConfigurePortAccessEgress struct {
	Pool []*ConfigurePortAccessEgressPool `json:"pool,omitempty"`
}

// ConfigurePortAccessEgressPool struct
type ConfigurePortAccessEgressPool struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=1000
	AmberAlarmThreshold *uint32 `json:"amber-alarm-threshold,omitempty"`
	ApplyGroups         *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude  *string `json:"apply-groups-exclude,omitempty"`
	Name                *string `json:"name,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=1000
	RedAlarmThreshold *uint32                               `json:"red-alarm-threshold,omitempty"`
	ResvCbs           *ConfigurePortAccessEgressPoolResvCbs `json:"resv-cbs,omitempty"`
	SlopePolicy       *string                               `json:"slope-policy,omitempty"`
}

// ConfigurePortAccessEgressPoolResvCbs struct
type ConfigurePortAccessEgressPoolResvCbs struct {
	AmberAlarmAction *ConfigurePortAccessEgressPoolResvCbsAmberAlarmAction `json:"amber-alarm-action,omitempty"`
	// +kubebuilder:default:="auto"
	Cbs *string `json:"cbs,omitempty"`
}

// ConfigurePortAccessEgressPoolResvCbsAmberAlarmAction struct
type ConfigurePortAccessEgressPoolResvCbsAmberAlarmAction struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=100
	Max *uint32 `json:"max,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=100
	Step *uint32 `json:"step,omitempty"`
}

// ConfigurePortAccessIngress struct
type ConfigurePortAccessIngress struct {
	Pool []*ConfigurePortAccessIngressPool `json:"pool,omitempty"`
}

// ConfigurePortAccessIngressPool struct
type ConfigurePortAccessIngressPool struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=1000
	AmberAlarmThreshold *uint32 `json:"amber-alarm-threshold,omitempty"`
	ApplyGroups         *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude  *string `json:"apply-groups-exclude,omitempty"`
	Name                *string `json:"name,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=1000
	RedAlarmThreshold *uint32                                `json:"red-alarm-threshold,omitempty"`
	ResvCbs           *ConfigurePortAccessIngressPoolResvCbs `json:"resv-cbs,omitempty"`
	SlopePolicy       *string                                `json:"slope-policy,omitempty"`
}

// ConfigurePortAccessIngressPoolResvCbs struct
type ConfigurePortAccessIngressPoolResvCbs struct {
	AmberAlarmAction *ConfigurePortAccessIngressPoolResvCbsAmberAlarmAction `json:"amber-alarm-action,omitempty"`
	// +kubebuilder:default:="auto"
	Cbs *string `json:"cbs,omitempty"`
}

// ConfigurePortAccessIngressPoolResvCbsAmberAlarmAction struct
type ConfigurePortAccessIngressPoolResvCbsAmberAlarmAction struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=100
	Max *uint32 `json:"max,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=100
	Step *uint32 `json:"step,omitempty"`
}

// ConfigurePortConnector struct
type ConfigurePortConnector struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:validation:Enum=`c1-100g`;`c1-10g`;`c1-25g`;`c1-400g`;`c1-40g`;`c1-50g`;`c10-10g`;`c2-100g`;`c4-100g`;`c4-10g`;`c4-25g`;`c8-50g`
	Breakout *string `json:"breakout,omitempty"`
	// +kubebuilder:validation:Enum=`cl91-514-528`;`cl91-514-544`
	RsFecMode *string `json:"rs-fec-mode,omitempty"`
}

// ConfigurePortDistCpuProtection struct
type ConfigurePortDistCpuProtection struct {
	Policy *string `json:"policy,omitempty"`
}

// ConfigurePortDwdm struct
type ConfigurePortDwdm struct {
	ApplyGroups        *string                       `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                       `json:"apply-groups-exclude,omitempty"`
	Channel            *string                       `json:"channel,omitempty"`
	Coherent           *ConfigurePortDwdmCoherent    `json:"coherent,omitempty"`
	RxdtvAdjust        *bool                         `json:"rxdtv-adjust,omitempty"`
	Wavetracker        *ConfigurePortDwdmWavetracker `json:"wavetracker,omitempty"`
}

// ConfigurePortDwdmCoherent struct
type ConfigurePortDwdmCoherent struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	Channel            *string `json:"channel,omitempty"`
	Compatibility      *string `json:"compatibility,omitempty"`
	// kubebuilder:validation:Minimum=2
	// kubebuilder:validation:Maximum=2
	// +kubebuilder:default:=32
	CprWindowSize *uint32 `json:"cpr-window-size,omitempty"`
	// kubebuilder:validation:Minimum=50000
	// kubebuilder:validation:Maximum=50000
	Dispersion    *int32                                `json:"dispersion,omitempty"`
	Mode          *string                               `json:"mode,omitempty"`
	ReportAlarm   *ConfigurePortDwdmCoherentReportAlarm `json:"report-alarm,omitempty"`
	RxLosReaction *string                               `json:"rx-los-reaction,omitempty"`
	// kubebuilder:validation:Minimum=3000
	// kubebuilder:validation:Maximum=1300
	// +kubebuilder:default:="-23"
	RxLosThresh *string                         `json:"rx-los-thresh,omitempty"`
	Sweep       *ConfigurePortDwdmCoherentSweep `json:"sweep,omitempty"`
	// kubebuilder:validation:Minimum=2000
	// kubebuilder:validation:Maximum=300
	// +kubebuilder:default:="1"
	TargetPower *string `json:"target-power,omitempty"`
}

// ConfigurePortDwdmCoherentReportAlarm struct
type ConfigurePortDwdmCoherentReportAlarm struct {
	// +kubebuilder:default:=true
	Hosttx *bool `json:"hosttx,omitempty"`
	// +kubebuilder:default:=true
	Mod *bool `json:"mod,omitempty"`
	// +kubebuilder:default:=true
	Modflt *bool `json:"modflt,omitempty"`
	// +kubebuilder:default:=true
	Netrx *bool `json:"netrx,omitempty"`
	// +kubebuilder:default:=true
	Nettx *bool `json:"nettx,omitempty"`
}

// ConfigurePortDwdmCoherentSweep struct
type ConfigurePortDwdmCoherentSweep struct {
	// kubebuilder:validation:Minimum=50000
	// kubebuilder:validation:Maximum=50000
	// +kubebuilder:default:=2000
	End *int32 `json:"end,omitempty"`
	// kubebuilder:validation:Minimum=50000
	// kubebuilder:validation:Maximum=50000
	// +kubebuilder:default:=-25500
	Start *int32 `json:"start,omitempty"`
}

// ConfigurePortDwdmWavetracker struct
type ConfigurePortDwdmWavetracker struct {
	ApplyGroups        *string                                   `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                                   `json:"apply-groups-exclude,omitempty"`
	Encode             *ConfigurePortDwdmWavetrackerEncode       `json:"encode,omitempty"`
	PowerControl       *ConfigurePortDwdmWavetrackerPowerControl `json:"power-control,omitempty"`
	ReportAlarm        *ConfigurePortDwdmWavetrackerReportAlarm  `json:"report-alarm,omitempty"`
}

// ConfigurePortDwdmWavetrackerEncode struct
type ConfigurePortDwdmWavetrackerEncode struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=4095
	Key1 *uint32 `json:"key1"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=4095
	Key2 *uint32 `json:"key2"`
}

// ConfigurePortDwdmWavetrackerPowerControl struct
type ConfigurePortDwdmWavetrackerPowerControl struct {
	// kubebuilder:validation:Minimum=2200
	// kubebuilder:validation:Maximum=300
	// +kubebuilder:default:="-20"
	TargetPower *string `json:"target-power,omitempty"`
}

// ConfigurePortDwdmWavetrackerReportAlarm struct
type ConfigurePortDwdmWavetrackerReportAlarm struct {
	// +kubebuilder:default:=true
	EncoderDegrade *bool `json:"encoder-degrade,omitempty"`
	// +kubebuilder:default:=true
	EncoderFailure *bool `json:"encoder-failure,omitempty"`
	// +kubebuilder:default:=true
	MissingPluggableVoa *bool `json:"missing-pluggable-voa,omitempty"`
	// +kubebuilder:default:=true
	PowerControlDegrade *bool `json:"power-control-degrade,omitempty"`
	// +kubebuilder:default:=true
	PowerControlFailure *bool `json:"power-control-failure,omitempty"`
	// +kubebuilder:default:=true
	PowerControlHighLimit *bool `json:"power-control-high-limit,omitempty"`
	// +kubebuilder:default:=true
	PowerControlLowLimit *bool `json:"power-control-low-limit,omitempty"`
}

// ConfigurePortEthernet struct
type ConfigurePortEthernet struct {
	Access             *ConfigurePortEthernetAccess `json:"access,omitempty"`
	AccountingPolicy   *string                      `json:"accounting-policy,omitempty"`
	ApplyGroups        *string                      `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                      `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:validation:Enum=`false`;`limited`;`true`
	Autonegotiate *string `json:"autonegotiate,omitempty"`
	// +kubebuilder:default:=false
	CollectStats *bool                            `json:"collect-stats,omitempty"`
	CrcMonitor   *ConfigurePortEthernetCrcMonitor `json:"crc-monitor,omitempty"`
	Dampening    *ConfigurePortEthernetDampening  `json:"dampening,omitempty"`
	// +kubebuilder:default:=false
	DiscardRxPauseFrames    *bool                                         `json:"discard-rx-pause-frames,omitempty"`
	Dot1qEtype              *string                                       `json:"dot1q-etype,omitempty"`
	Dot1x                   *ConfigurePortEthernetDot1x                   `json:"dot1x,omitempty"`
	DownOnInternalError     *ConfigurePortEthernetDownOnInternalError     `json:"down-on-internal-error,omitempty"`
	DownWhenLooped          *ConfigurePortEthernetDownWhenLooped          `json:"down-when-looped,omitempty"`
	Duplex                  *string                                       `json:"duplex,omitempty"`
	EfmOam                  *ConfigurePortEthernetEfmOam                  `json:"efm-oam,omitempty"`
	Egress                  *ConfigurePortEthernetEgress                  `json:"egress,omitempty"`
	Elmi                    *ConfigurePortEthernetElmi                    `json:"elmi,omitempty"`
	EncapType               *string                                       `json:"encap-type,omitempty"`
	EthCfm                  *ConfigurePortEthernetEthCfm                  `json:"eth-cfm,omitempty"`
	HoldTime                *ConfigurePortEthernetHoldTime                `json:"hold-time,omitempty"`
	HsmdaSchedulerOverrides *ConfigurePortEthernetHsmdaSchedulerOverrides `json:"hsmda-scheduler-overrides,omitempty"`
	Ingress                 *ConfigurePortEthernetIngress                 `json:"ingress,omitempty"`
	// +kubebuilder:default:=false
	LacpTunnel             *bool                          `json:"lacp-tunnel,omitempty"`
	Lldp                   *ConfigurePortEthernetLldp     `json:"lldp,omitempty"`
	LoadBalancingAlgorithm *string                        `json:"load-balancing-algorithm,omitempty"`
	Loopback               *ConfigurePortEthernetLoopback `json:"loopback,omitempty"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`[0-9a-fA-F]{2}(:[0-9a-fA-F]{2}){5}`
	// +kubebuilder:default:="00:00:00:00:00:00"
	MacAddress *string `json:"mac-address,omitempty"`
	// kubebuilder:validation:Minimum=64
	// kubebuilder:validation:Maximum=64
	// +kubebuilder:default:=64
	MinFrameLength *uint32 `json:"min-frame-length,omitempty"`
	Mode           *string `json:"mode,omitempty"`
	// kubebuilder:validation:Minimum=512
	// kubebuilder:validation:Maximum=9800
	Mtu      *uint32                       `json:"mtu,omitempty"`
	Network  *ConfigurePortEthernetNetwork `json:"network,omitempty"`
	PbbEtype *string                       `json:"pbb-etype,omitempty"`
	// kubebuilder:validation:Minimum=2147483648
	// kubebuilder:validation:Maximum=2147483647
	PtpAsymmetry *int32                            `json:"ptp-asymmetry,omitempty"`
	QinqEtype    *string                           `json:"qinq-etype,omitempty"`
	ReportAlarm  *ConfigurePortEthernetReportAlarm `json:"report-alarm,omitempty"`
	// +kubebuilder:validation:Enum=`cl108`;`cl74`;`cl91-514-528`
	RsFecMode *string `json:"rs-fec-mode,omitempty"`
	// +kubebuilder:default:=false
	SingleFiber *bool `json:"single-fiber,omitempty"`
	// kubebuilder:validation:Minimum=10
	// kubebuilder:validation:Maximum=10
	Speed         *uint32                             `json:"speed,omitempty"`
	Ssm           *ConfigurePortEthernetSsm           `json:"ssm,omitempty"`
	SymbolMonitor *ConfigurePortEthernetSymbolMonitor `json:"symbol-monitor,omitempty"`
	// kubebuilder:validation:Minimum=30
	// kubebuilder:validation:Maximum=600
	// +kubebuilder:default:=300
	UtilStatsInterval *uint32 `json:"util-stats-interval,omitempty"`
	// +kubebuilder:validation:Enum=`lan`;`wan`
	Xgig *string `json:"xgig,omitempty"`
}

// ConfigurePortEthernetAccess struct
type ConfigurePortEthernetAccess struct {
	AccountingPolicy   *string `json:"accounting-policy,omitempty"`
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=6400000000
	Bandwidth *uint64 `json:"bandwidth,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=1000
	// +kubebuilder:default:=100
	BookingFactor *uint32 `json:"booking-factor,omitempty"`
	// +kubebuilder:default:=false
	CollectStats *bool                               `json:"collect-stats,omitempty"`
	Egress       *ConfigurePortEthernetAccessEgress  `json:"egress,omitempty"`
	Ingress      *ConfigurePortEthernetAccessIngress `json:"ingress,omitempty"`
}

// ConfigurePortEthernetAccessEgress struct
type ConfigurePortEthernetAccessEgress struct {
	QueueGroup  []*ConfigurePortEthernetAccessEgressQueueGroup  `json:"queue-group,omitempty"`
	VirtualPort []*ConfigurePortEthernetAccessEgressVirtualPort `json:"virtual-port,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroup struct
type ConfigurePortEthernetAccessEgressQueueGroup struct {
	AccountingPolicy   *string                                                   `json:"accounting-policy,omitempty"`
	AggregateRate      *ConfigurePortEthernetAccessEgressQueueGroupAggregateRate `json:"aggregate-rate,omitempty"`
	ApplyGroups        *string                                                   `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                                                   `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:default:=false
	CollectStats *bool                                                 `json:"collect-stats,omitempty"`
	Description  *string                                               `json:"description,omitempty"`
	HostMatch    *ConfigurePortEthernetAccessEgressQueueGroupHostMatch `json:"host-match,omitempty"`
	// +kubebuilder:default:=false
	HsTurbo             *bool                                                           `json:"hs-turbo,omitempty"`
	HsmdaQueueOverrides *ConfigurePortEthernetAccessEgressQueueGroupHsmdaQueueOverrides `json:"hsmda-queue-overrides,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=65535
	InstanceId      *uint16                                                     `json:"instance-id,omitempty"`
	QueueGroupName  *string                                                     `json:"queue-group-name,omitempty"`
	QueueOverrides  *ConfigurePortEthernetAccessEgressQueueGroupQueueOverrides  `json:"queue-overrides,omitempty"`
	SchedulerPolicy *ConfigurePortEthernetAccessEgressQueueGroupSchedulerPolicy `json:"scheduler-policy,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroupAggregateRate struct
type ConfigurePortEthernetAccessEgressQueueGroupAggregateRate struct {
	// +kubebuilder:default:=false
	LimitUnusedBandwidth *bool `json:"limit-unused-bandwidth,omitempty"`
	// +kubebuilder:default:=false
	QueueFrameBasedAccounting *bool `json:"queue-frame-based-accounting,omitempty"`
	// +kubebuilder:default:="max"
	Rate *string `json:"rate,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroupHostMatch struct
type ConfigurePortEthernetAccessEgressQueueGroupHostMatch struct {
	IntDestId []*ConfigurePortEthernetAccessEgressQueueGroupHostMatchIntDestId `json:"int-dest-id,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroupHostMatchIntDestId struct
type ConfigurePortEthernetAccessEgressQueueGroupHostMatchIntDestId struct {
	DestinationString *string `json:"destination-string,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroupHsmdaQueueOverrides struct
type ConfigurePortEthernetAccessEgressQueueGroupHsmdaQueueOverrides struct {
	PacketByteOffset *string                                                                `json:"packet-byte-offset,omitempty"`
	Queue            []*ConfigurePortEthernetAccessEgressQueueGroupHsmdaQueueOverridesQueue `json:"queue,omitempty"`
	SecondaryShaper  *string                                                                `json:"secondary-shaper,omitempty"`
	WrrPolicy        *string                                                                `json:"wrr-policy,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroupHsmdaQueueOverridesQueue struct
type ConfigurePortEthernetAccessEgressQueueGroupHsmdaQueueOverridesQueue struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	Mbs                *string `json:"mbs,omitempty"`
	QueueId            *string `json:"queue-id,omitempty"`
	Rate               *string `json:"rate,omitempty"`
	SlopePolicy        *string `json:"slope-policy,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=32
	WrrWeight *int32 `json:"wrr-weight,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroupQueueOverrides struct
type ConfigurePortEthernetAccessEgressQueueGroupQueueOverrides struct {
	Queue []*ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueue `json:"queue,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueue struct
type ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueue struct {
	AdaptationRule     *ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueAdaptationRule `json:"adaptation-rule,omitempty"`
	ApplyGroups        *string                                                                       `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                                                                       `json:"apply-groups-exclude,omitempty"`
	BurstLimit         *string                                                                       `json:"burst-limit,omitempty"`
	Cbs                *string                                                                       `json:"cbs,omitempty"`
	DropTail           *ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueDropTail       `json:"drop-tail,omitempty"`
	Mbs                *string                                                                       `json:"mbs,omitempty"`
	// +kubebuilder:default:=false
	MonitorDepth      *bool                                                                            `json:"monitor-depth,omitempty"`
	MonitorQueueDepth *ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueMonitorQueueDepth `json:"monitor-queue-depth,omitempty"`
	Parent            *ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueParent            `json:"parent,omitempty"`
	QueueId           *string                                                                          `json:"queue-id,omitempty"`
	QueueOverrideRate *ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueQueueOverrideRate `json:"queue-override-rate,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueAdaptationRule struct
type ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueAdaptationRule struct {
	Cir *string `json:"cir,omitempty"`
	Pir *string `json:"pir,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueDropTail struct
type ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueDropTail struct {
	Low *ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueDropTailLow `json:"low,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueDropTailLow struct
type ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueDropTailLow struct {
	PercentReductionFromMbs *string `json:"percent-reduction-from-mbs,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueMonitorQueueDepth struct
type ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueMonitorQueueDepth struct {
	// +kubebuilder:default:=false
	FastPolling *bool `json:"fast-polling,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=9999
	ViolationThreshold *string `json:"violation-threshold,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueParent struct
type ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueParent struct {
	CirWeight *string `json:"cir-weight,omitempty"`
	Weight    *string `json:"weight,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueQueueOverrideRate struct
type ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueQueueOverrideRate struct {
	PercentRate *ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueQueueOverrideRatePercentRate `json:"percent-rate,omitempty"`
	Rate        *ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueQueueOverrideRateRate        `json:"rate,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueQueueOverrideRatePercentRate struct
type ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueQueueOverrideRatePercentRate struct {
	PercentRate *ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueQueueOverrideRatePercentRatePercentRate `json:"percent-rate,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueQueueOverrideRatePercentRatePercentRate struct
type ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueQueueOverrideRatePercentRatePercentRate struct {
	Cir *string `json:"cir,omitempty"`
	Pir *string `json:"pir,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueQueueOverrideRateRate struct
type ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueQueueOverrideRateRate struct {
	Rate *ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueQueueOverrideRateRateRate `json:"rate,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueQueueOverrideRateRateRate struct
type ConfigurePortEthernetAccessEgressQueueGroupQueueOverridesQueueQueueOverrideRateRateRate struct {
	Cir *string `json:"cir,omitempty"`
	Pir *string `json:"pir,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroupSchedulerPolicy struct
type ConfigurePortEthernetAccessEgressQueueGroupSchedulerPolicy struct {
	Overrides  *ConfigurePortEthernetAccessEgressQueueGroupSchedulerPolicyOverrides `json:"overrides,omitempty"`
	PolicyName *string                                                              `json:"policy-name,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroupSchedulerPolicyOverrides struct
type ConfigurePortEthernetAccessEgressQueueGroupSchedulerPolicyOverrides struct {
	Scheduler []*ConfigurePortEthernetAccessEgressQueueGroupSchedulerPolicyOverridesScheduler `json:"scheduler,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroupSchedulerPolicyOverridesScheduler struct
type ConfigurePortEthernetAccessEgressQueueGroupSchedulerPolicyOverridesScheduler struct {
	ApplyGroups        *string                                                                             `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                                                                             `json:"apply-groups-exclude,omitempty"`
	Parent             *ConfigurePortEthernetAccessEgressQueueGroupSchedulerPolicyOverridesSchedulerParent `json:"parent,omitempty"`
	Rate               *ConfigurePortEthernetAccessEgressQueueGroupSchedulerPolicyOverridesSchedulerRate   `json:"rate,omitempty"`
	SchedulerName      *string                                                                             `json:"scheduler-name,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroupSchedulerPolicyOverridesSchedulerParent struct
type ConfigurePortEthernetAccessEgressQueueGroupSchedulerPolicyOverridesSchedulerParent struct {
	CirWeight *string `json:"cir-weight,omitempty"`
	Weight    *string `json:"weight,omitempty"`
}

// ConfigurePortEthernetAccessEgressQueueGroupSchedulerPolicyOverridesSchedulerRate struct
type ConfigurePortEthernetAccessEgressQueueGroupSchedulerPolicyOverridesSchedulerRate struct {
	Cir *string `json:"cir,omitempty"`
	Pir *string `json:"pir,omitempty"`
}

// ConfigurePortEthernetAccessEgressVirtualPort struct
type ConfigurePortEthernetAccessEgressVirtualPort struct {
	AggregateRate              *ConfigurePortEthernetAccessEgressVirtualPortAggregateRate `json:"aggregate-rate,omitempty"`
	ApplyGroups                *string                                                    `json:"apply-groups,omitempty"`
	ApplyGroupsExclude         *string                                                    `json:"apply-groups-exclude,omitempty"`
	Description                *string                                                    `json:"description,omitempty"`
	HostMatch                  *ConfigurePortEthernetAccessEgressVirtualPortHostMatch     `json:"host-match,omitempty"`
	HwAggShaperSchedulerPolicy *string                                                    `json:"hw-agg-shaper-scheduler-policy,omitempty"`
	// +kubebuilder:default:=false
	MonitorHwAggShaperScheduler *bool `json:"monitor-hw-agg-shaper-scheduler,omitempty"`
	// +kubebuilder:default:=false
	MonitorPortScheduler *bool `json:"monitor-port-scheduler,omitempty"`
	// +kubebuilder:default:=false
	MulticastHqosAdjustment *bool   `json:"multicast-hqos-adjustment,omitempty"`
	PortSchedulerPolicy     *string `json:"port-scheduler-policy,omitempty"`
	SchedulerPolicy         *string `json:"scheduler-policy,omitempty"`
	VportName               *string `json:"vport-name,omitempty"`
}

// ConfigurePortEthernetAccessEgressVirtualPortAggregateRate struct
type ConfigurePortEthernetAccessEgressVirtualPortAggregateRate struct {
	// +kubebuilder:default:=false
	LimitUnusedBandwidth *bool `json:"limit-unused-bandwidth,omitempty"`
	// +kubebuilder:default:="max"
	Rate *string `json:"rate,omitempty"`
}

// ConfigurePortEthernetAccessEgressVirtualPortHostMatch struct
type ConfigurePortEthernetAccessEgressVirtualPortHostMatch struct {
	IntDestId []*ConfigurePortEthernetAccessEgressVirtualPortHostMatchIntDestId `json:"int-dest-id,omitempty"`
}

// ConfigurePortEthernetAccessEgressVirtualPortHostMatchIntDestId struct
type ConfigurePortEthernetAccessEgressVirtualPortHostMatchIntDestId struct {
	DestinationString *string `json:"destination-string,omitempty"`
}

// ConfigurePortEthernetAccessIngress struct
type ConfigurePortEthernetAccessIngress struct {
	QueueGroup []*ConfigurePortEthernetAccessIngressQueueGroup `json:"queue-group,omitempty"`
}

// ConfigurePortEthernetAccessIngressQueueGroup struct
type ConfigurePortEthernetAccessIngressQueueGroup struct {
	AccountingPolicy   *string `json:"accounting-policy,omitempty"`
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:default:=false
	CollectStats    *bool                                                        `json:"collect-stats,omitempty"`
	Description     *string                                                      `json:"description,omitempty"`
	QueueGroupName  *string                                                      `json:"queue-group-name,omitempty"`
	QueueOverrides  *ConfigurePortEthernetAccessIngressQueueGroupQueueOverrides  `json:"queue-overrides,omitempty"`
	SchedulerPolicy *ConfigurePortEthernetAccessIngressQueueGroupSchedulerPolicy `json:"scheduler-policy,omitempty"`
}

// ConfigurePortEthernetAccessIngressQueueGroupQueueOverrides struct
type ConfigurePortEthernetAccessIngressQueueGroupQueueOverrides struct {
	Queue []*ConfigurePortEthernetAccessIngressQueueGroupQueueOverridesQueue `json:"queue,omitempty"`
}

// ConfigurePortEthernetAccessIngressQueueGroupQueueOverridesQueue struct
type ConfigurePortEthernetAccessIngressQueueGroupQueueOverridesQueue struct {
	AdaptationRule     *ConfigurePortEthernetAccessIngressQueueGroupQueueOverridesQueueAdaptationRule `json:"adaptation-rule,omitempty"`
	ApplyGroups        *string                                                                        `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                                                                        `json:"apply-groups-exclude,omitempty"`
	Cbs                *string                                                                        `json:"cbs,omitempty"`
	DropTail           *ConfigurePortEthernetAccessIngressQueueGroupQueueOverridesQueueDropTail       `json:"drop-tail,omitempty"`
	Mbs                *string                                                                        `json:"mbs,omitempty"`
	// +kubebuilder:default:=false
	MonitorDepth      *bool                                                                             `json:"monitor-depth,omitempty"`
	MonitorQueueDepth *ConfigurePortEthernetAccessIngressQueueGroupQueueOverridesQueueMonitorQueueDepth `json:"monitor-queue-depth,omitempty"`
	QueueId           *string                                                                           `json:"queue-id,omitempty"`
	Rate              *ConfigurePortEthernetAccessIngressQueueGroupQueueOverridesQueueRate              `json:"rate,omitempty"`
}

// ConfigurePortEthernetAccessIngressQueueGroupQueueOverridesQueueAdaptationRule struct
type ConfigurePortEthernetAccessIngressQueueGroupQueueOverridesQueueAdaptationRule struct {
	Cir *string `json:"cir,omitempty"`
	Pir *string `json:"pir,omitempty"`
}

// ConfigurePortEthernetAccessIngressQueueGroupQueueOverridesQueueDropTail struct
type ConfigurePortEthernetAccessIngressQueueGroupQueueOverridesQueueDropTail struct {
	Low *ConfigurePortEthernetAccessIngressQueueGroupQueueOverridesQueueDropTailLow `json:"low,omitempty"`
}

// ConfigurePortEthernetAccessIngressQueueGroupQueueOverridesQueueDropTailLow struct
type ConfigurePortEthernetAccessIngressQueueGroupQueueOverridesQueueDropTailLow struct {
	PercentReductionFromMbs *string `json:"percent-reduction-from-mbs,omitempty"`
}

// ConfigurePortEthernetAccessIngressQueueGroupQueueOverridesQueueMonitorQueueDepth struct
type ConfigurePortEthernetAccessIngressQueueGroupQueueOverridesQueueMonitorQueueDepth struct {
	// +kubebuilder:default:=false
	FastPolling *bool `json:"fast-polling,omitempty"`
}

// ConfigurePortEthernetAccessIngressQueueGroupQueueOverridesQueueRate struct
type ConfigurePortEthernetAccessIngressQueueGroupQueueOverridesQueueRate struct {
	Cir *string `json:"cir,omitempty"`
	Pir *string `json:"pir,omitempty"`
}

// ConfigurePortEthernetAccessIngressQueueGroupSchedulerPolicy struct
type ConfigurePortEthernetAccessIngressQueueGroupSchedulerPolicy struct {
	Overrides  *ConfigurePortEthernetAccessIngressQueueGroupSchedulerPolicyOverrides `json:"overrides,omitempty"`
	PolicyName *string                                                               `json:"policy-name,omitempty"`
}

// ConfigurePortEthernetAccessIngressQueueGroupSchedulerPolicyOverrides struct
type ConfigurePortEthernetAccessIngressQueueGroupSchedulerPolicyOverrides struct {
	Scheduler []*ConfigurePortEthernetAccessIngressQueueGroupSchedulerPolicyOverridesScheduler `json:"scheduler,omitempty"`
}

// ConfigurePortEthernetAccessIngressQueueGroupSchedulerPolicyOverridesScheduler struct
type ConfigurePortEthernetAccessIngressQueueGroupSchedulerPolicyOverridesScheduler struct {
	ApplyGroups        *string                                                                              `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                                                                              `json:"apply-groups-exclude,omitempty"`
	Parent             *ConfigurePortEthernetAccessIngressQueueGroupSchedulerPolicyOverridesSchedulerParent `json:"parent,omitempty"`
	Rate               *ConfigurePortEthernetAccessIngressQueueGroupSchedulerPolicyOverridesSchedulerRate   `json:"rate,omitempty"`
	SchedulerName      *string                                                                              `json:"scheduler-name,omitempty"`
}

// ConfigurePortEthernetAccessIngressQueueGroupSchedulerPolicyOverridesSchedulerParent struct
type ConfigurePortEthernetAccessIngressQueueGroupSchedulerPolicyOverridesSchedulerParent struct {
	CirWeight *string `json:"cir-weight,omitempty"`
	Weight    *string `json:"weight,omitempty"`
}

// ConfigurePortEthernetAccessIngressQueueGroupSchedulerPolicyOverridesSchedulerRate struct
type ConfigurePortEthernetAccessIngressQueueGroupSchedulerPolicyOverridesSchedulerRate struct {
	Cir *string `json:"cir,omitempty"`
	Pir *string `json:"pir,omitempty"`
}

// ConfigurePortEthernetCrcMonitor struct
type ConfigurePortEthernetCrcMonitor struct {
	SignalDegrade *ConfigurePortEthernetCrcMonitorSignalDegrade `json:"signal-degrade,omitempty"`
	SignalFailure *ConfigurePortEthernetCrcMonitorSignalFailure `json:"signal-failure,omitempty"`
	// kubebuilder:validation:Minimum=5
	// kubebuilder:validation:Maximum=60
	// +kubebuilder:default:=10
	WindowSize *uint32 `json:"window-size,omitempty"`
}

// ConfigurePortEthernetCrcMonitorSignalDegrade struct
type ConfigurePortEthernetCrcMonitorSignalDegrade struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=9
	// +kubebuilder:default:=1
	Multiplier *uint32 `json:"multiplier,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=9
	Threshold *uint32 `json:"threshold,omitempty"`
}

// ConfigurePortEthernetCrcMonitorSignalFailure struct
type ConfigurePortEthernetCrcMonitorSignalFailure struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=9
	// +kubebuilder:default:=1
	Multiplier *uint32 `json:"multiplier,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=9
	Threshold *uint32 `json:"threshold,omitempty"`
}

// ConfigurePortEthernetDampening struct
type ConfigurePortEthernetDampening struct {
	AdminState         *string `json:"admin-state,omitempty"`
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=2000
	// +kubebuilder:default:=5
	HalfLife *uint32 `json:"half-life,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=43200
	// +kubebuilder:default:=20
	MaxSuppressTime *uint32 `json:"max-suppress-time,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=20000
	// +kubebuilder:default:=1000
	ReuseThreshold *uint32 `json:"reuse-threshold,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=20000
	// +kubebuilder:default:=2000
	SuppressThreshold *uint32 `json:"suppress-threshold,omitempty"`
}

// ConfigurePortEthernetDot1x struct
type ConfigurePortEthernetDot1x struct {
	ApplyGroups        *string                           `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                           `json:"apply-groups-exclude,omitempty"`
	Macsec             *ConfigurePortEthernetDot1xMacsec `json:"macsec,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=10
	// +kubebuilder:default:=2
	MaxAuthenticationRequests *uint32                                          `json:"max-authentication-requests,omitempty"`
	PerHostAuthentication     *ConfigurePortEthernetDot1xPerHostAuthentication `json:"per-host-authentication,omitempty"`
	PortControl               *string                                          `json:"port-control,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=3600
	// +kubebuilder:default:=60
	QuietPeriod              *uint32                                             `json:"quiet-period,omitempty"`
	RadiusPolicy             *string                                             `json:"radius-policy,omitempty"`
	RadiusServerPolicyConfig *ConfigurePortEthernetDot1xRadiusServerPolicyConfig `json:"radius-server-policy-config,omitempty"`
	ReAuthentication         *ConfigurePortEthernetDot1xReAuthentication         `json:"re-authentication,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=300
	// +kubebuilder:default:=30
	ServerTimeout *uint32 `json:"server-timeout,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=300
	// +kubebuilder:default:=30
	SupplicantTimeout *uint32 `json:"supplicant-timeout,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=3600
	// +kubebuilder:default:=30
	TransmitPeriod *uint32 `json:"transmit-period,omitempty"`
	// +kubebuilder:default:=true
	TunnelDot1q *bool `json:"tunnel-dot1q,omitempty"`
	// +kubebuilder:default:=true
	TunnelQinq *bool `json:"tunnel-qinq,omitempty"`
	// +kubebuilder:default:=false
	Tunneling *bool `json:"tunneling,omitempty"`
}

// ConfigurePortEthernetDot1xMacsec struct
type ConfigurePortEthernetDot1xMacsec struct {
	ApplyGroups        *string                                          `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                                          `json:"apply-groups-exclude,omitempty"`
	ExcludeMacPolicy   *string                                          `json:"exclude-mac-policy,omitempty"`
	ExcludeProtocol    *ConfigurePortEthernetDot1xMacsecExcludeProtocol `json:"exclude-protocol,omitempty"`
	// +kubebuilder:default:=false
	RxMustBeEncrypted *bool                                      `json:"rx-must-be-encrypted,omitempty"`
	SubPort           []*ConfigurePortEthernetDot1xMacsecSubPort `json:"sub-port,omitempty"`
}

// ConfigurePortEthernetDot1xMacsecExcludeProtocol struct
type ConfigurePortEthernetDot1xMacsecExcludeProtocol struct {
	// +kubebuilder:default:=false
	Cdp *bool `json:"cdp,omitempty"`
	// +kubebuilder:default:=false
	EapolStart *bool `json:"eapol-start,omitempty"`
	// +kubebuilder:default:=false
	EfmOam *bool `json:"efm-oam,omitempty"`
	// +kubebuilder:default:=false
	EthCfm *bool `json:"eth-cfm,omitempty"`
	// +kubebuilder:default:=false
	Lacp *bool `json:"lacp,omitempty"`
	// +kubebuilder:default:=false
	Lldp *bool `json:"lldp,omitempty"`
	// +kubebuilder:default:=false
	Ptp *bool `json:"ptp,omitempty"`
	// +kubebuilder:default:=false
	Ubfd *bool `json:"ubfd,omitempty"`
}

// ConfigurePortEthernetDot1xMacsecSubPort struct
type ConfigurePortEthernetDot1xMacsecSubPort struct {
	AdminState         *string `json:"admin-state,omitempty"`
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	CaName             *string `json:"ca-name,omitempty"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`[0-9a-fA-F]{2}(:[0-9a-fA-F]{2}){5}`
	EapolDestinationAddress *string                                            `json:"eapol-destination-address,omitempty"`
	EncapMatch              *ConfigurePortEthernetDot1xMacsecSubPortEncapMatch `json:"encap-match,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=32
	MaxPeers *uint32 `json:"max-peers,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=1023
	SubPortId *uint32 `json:"sub-port-id,omitempty"`
}

// ConfigurePortEthernetDot1xMacsecSubPortEncapMatch struct
type ConfigurePortEthernetDot1xMacsecSubPortEncapMatch struct {
	Encap *ConfigurePortEthernetDot1xMacsecSubPortEncapMatchEncap `json:"encap,omitempty"`
}

// ConfigurePortEthernetDot1xMacsecSubPortEncapMatchEncap struct
type ConfigurePortEthernetDot1xMacsecSubPortEncapMatchEncap struct {
	AllMatch  *ConfigurePortEthernetDot1xMacsecSubPortEncapMatchEncapAllMatch  `json:"all-match,omitempty"`
	DoubleTag *ConfigurePortEthernetDot1xMacsecSubPortEncapMatchEncapDoubleTag `json:"double-tag,omitempty"`
	SingleTag *ConfigurePortEthernetDot1xMacsecSubPortEncapMatchEncapSingleTag `json:"single-tag,omitempty"`
	Untagged  *ConfigurePortEthernetDot1xMacsecSubPortEncapMatchEncapUntagged  `json:"untagged,omitempty"`
}

// ConfigurePortEthernetDot1xMacsecSubPortEncapMatchEncapAllMatch struct
type ConfigurePortEthernetDot1xMacsecSubPortEncapMatchEncapAllMatch struct {
	// +kubebuilder:default:=true
	AllMatch *bool `json:"all-match,omitempty"`
}

// ConfigurePortEthernetDot1xMacsecSubPortEncapMatchEncapDoubleTag struct
type ConfigurePortEthernetDot1xMacsecSubPortEncapMatchEncapDoubleTag struct {
	DoubleTag *string `json:"double-tag,omitempty"`
}

// ConfigurePortEthernetDot1xMacsecSubPortEncapMatchEncapSingleTag struct
type ConfigurePortEthernetDot1xMacsecSubPortEncapMatchEncapSingleTag struct {
	SingleTag *string `json:"single-tag,omitempty"`
}

// ConfigurePortEthernetDot1xMacsecSubPortEncapMatchEncapUntagged struct
type ConfigurePortEthernetDot1xMacsecSubPortEncapMatchEncapUntagged struct {
	Untagged *bool `json:"untagged,omitempty"`
}

// ConfigurePortEthernetDot1xPerHostAuthentication struct
type ConfigurePortEthernetDot1xPerHostAuthentication struct {
	AdminState        *string                                                           `json:"admin-state,omitempty"`
	AllowedSourceMacs *ConfigurePortEthernetDot1xPerHostAuthenticationAllowedSourceMacs `json:"allowed-source-macs,omitempty"`
	// +kubebuilder:default:=true
	AuthenticatorInit *bool `json:"authenticator-init,omitempty"`
}

// ConfigurePortEthernetDot1xPerHostAuthenticationAllowedSourceMacs struct
type ConfigurePortEthernetDot1xPerHostAuthenticationAllowedSourceMacs struct {
	MacAddress []*ConfigurePortEthernetDot1xPerHostAuthenticationAllowedSourceMacsMacAddress `json:"mac-address,omitempty"`
}

// ConfigurePortEthernetDot1xPerHostAuthenticationAllowedSourceMacsMacAddress struct
type ConfigurePortEthernetDot1xPerHostAuthenticationAllowedSourceMacsMacAddress struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`[0-9a-fA-F]{2}(:[0-9a-fA-F]{2}){5}`
	Mac *string `json:"mac,omitempty"`
}

// ConfigurePortEthernetDot1xRadiusServerPolicyConfig struct
type ConfigurePortEthernetDot1xRadiusServerPolicyConfig struct {
	Common *ConfigurePortEthernetDot1xRadiusServerPolicyConfigCommon `json:"common,omitempty"`
	Split  *ConfigurePortEthernetDot1xRadiusServerPolicyConfigSplit  `json:"split,omitempty"`
}

// ConfigurePortEthernetDot1xRadiusServerPolicyConfigCommon struct
type ConfigurePortEthernetDot1xRadiusServerPolicyConfigCommon struct {
	RadiusServerPolicy *string `json:"radius-server-policy,omitempty"`
}

// ConfigurePortEthernetDot1xRadiusServerPolicyConfigSplit struct
type ConfigurePortEthernetDot1xRadiusServerPolicyConfigSplit struct {
	RadiusServerPolicyAcct *string `json:"radius-server-policy-acct,omitempty"`
	RadiusServerPolicyAuth *string `json:"radius-server-policy-auth,omitempty"`
}

// ConfigurePortEthernetDot1xReAuthentication struct
type ConfigurePortEthernetDot1xReAuthentication struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=9000
	// +kubebuilder:default:=3600
	Period *uint32 `json:"period,omitempty"`
}

// ConfigurePortEthernetDownOnInternalError struct
type ConfigurePortEthernetDownOnInternalError struct {
	// +kubebuilder:validation:Enum=`off`;`on`
	// +kubebuilder:default:="on"
	TxLaser *string `json:"tx-laser,omitempty"`
}

// ConfigurePortEthernetDownWhenLooped struct
type ConfigurePortEthernetDownWhenLooped struct {
	AdminState *string `json:"admin-state,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=120
	// +kubebuilder:default:=10
	KeepAlive *uint32 `json:"keep-alive,omitempty"`
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=0
	// +kubebuilder:default:=120
	RetryTimeout *uint32 `json:"retry-timeout,omitempty"`
	// +kubebuilder:default:=false
	UseBroadcastAddress *bool `json:"use-broadcast-address,omitempty"`
}

// ConfigurePortEthernetEfmOam struct
type ConfigurePortEthernetEfmOam struct {
	// +kubebuilder:default:=false
	AcceptRemoteLoopback *bool                                 `json:"accept-remote-loopback,omitempty"`
	AdminState           *string                               `json:"admin-state,omitempty"`
	ApplyGroups          *string                               `json:"apply-groups,omitempty"`
	ApplyGroupsExclude   *string                               `json:"apply-groups-exclude,omitempty"`
	Discovery            *ConfigurePortEthernetEfmOamDiscovery `json:"discovery,omitempty"`
	// +kubebuilder:default:=true
	DyingGaspTxOnReset *bool `json:"dying-gasp-tx-on-reset,omitempty"`
	GraceTx            *bool `json:"grace-tx,omitempty"`
	// +kubebuilder:default:="00:16:4D"
	GraceVendorOui *string `json:"grace-vendor-oui,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=50
	HoldTime *uint32 `json:"hold-time,omitempty"`
	// +kubebuilder:default:=false
	IgnoreEfmState *bool                                      `json:"ignore-efm-state,omitempty"`
	LinkMonitoring *ConfigurePortEthernetEfmOamLinkMonitoring `json:"link-monitoring,omitempty"`
	// +kubebuilder:validation:Enum=`active`;`passive`
	// +kubebuilder:default:="active"
	Mode *string `json:"mode,omitempty"`
	// kubebuilder:validation:Minimum=2
	// kubebuilder:validation:Maximum=5
	// +kubebuilder:default:=5
	Multiplier *uint32                               `json:"multiplier,omitempty"`
	PeerRdiRx  *ConfigurePortEthernetEfmOamPeerRdiRx `json:"peer-rdi-rx,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=600
	// +kubebuilder:default:=10
	TransmitInterval *uint32 `json:"transmit-interval,omitempty"`
	// +kubebuilder:validation:Enum=`critical-event`;`dying-gasp`
	TriggerFault *string `json:"trigger-fault,omitempty"`
	// +kubebuilder:default:=false
	Tunneling *bool `json:"tunneling,omitempty"`
}

// ConfigurePortEthernetEfmOamDiscovery struct
type ConfigurePortEthernetEfmOamDiscovery struct {
	AdvertiseCapabilities *ConfigurePortEthernetEfmOamDiscoveryAdvertiseCapabilities `json:"advertise-capabilities,omitempty"`
}

// ConfigurePortEthernetEfmOamDiscoveryAdvertiseCapabilities struct
type ConfigurePortEthernetEfmOamDiscoveryAdvertiseCapabilities struct {
	// +kubebuilder:default:=true
	LinkMonitoring *bool `json:"link-monitoring,omitempty"`
}

// ConfigurePortEthernetEfmOamLinkMonitoring struct
type ConfigurePortEthernetEfmOamLinkMonitoring struct {
	AdminState          *string                                                       `json:"admin-state,omitempty"`
	ApplyGroups         *string                                                       `json:"apply-groups,omitempty"`
	ApplyGroupsExclude  *string                                                       `json:"apply-groups-exclude,omitempty"`
	ErroredFrame        *ConfigurePortEthernetEfmOamLinkMonitoringErroredFrame        `json:"errored-frame,omitempty"`
	ErroredFramePeriod  *ConfigurePortEthernetEfmOamLinkMonitoringErroredFramePeriod  `json:"errored-frame-period,omitempty"`
	ErroredFrameSeconds *ConfigurePortEthernetEfmOamLinkMonitoringErroredFrameSeconds `json:"errored-frame-seconds,omitempty"`
	ErroredSymbols      *ConfigurePortEthernetEfmOamLinkMonitoringErroredSymbols      `json:"errored-symbols,omitempty"`
	LocalSfAction       *ConfigurePortEthernetEfmOamLinkMonitoringLocalSfAction       `json:"local-sf-action,omitempty"`
}

// ConfigurePortEthernetEfmOamLinkMonitoringErroredFrame struct
type ConfigurePortEthernetEfmOamLinkMonitoringErroredFrame struct {
	AdminState *string `json:"admin-state,omitempty"`
	// +kubebuilder:default:=true
	EventNotification *bool `json:"event-notification,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=1000000
	SdThreshold *uint32 `json:"sd-threshold,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=1000000
	// +kubebuilder:default:=1
	SfThreshold *uint32 `json:"sf-threshold,omitempty"`
	// kubebuilder:validation:Minimum=10
	// kubebuilder:validation:Maximum=600
	// +kubebuilder:default:=10
	Window *uint32 `json:"window,omitempty"`
}

// ConfigurePortEthernetEfmOamLinkMonitoringErroredFramePeriod struct
type ConfigurePortEthernetEfmOamLinkMonitoringErroredFramePeriod struct {
	AdminState *string `json:"admin-state,omitempty"`
	// +kubebuilder:default:=true
	EventNotification *bool `json:"event-notification,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=1000000
	SdThreshold *uint32 `json:"sd-threshold,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=1000000
	// +kubebuilder:default:=1
	SfThreshold *uint32 `json:"sf-threshold,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=4294967295
	// +kubebuilder:default:=1488095
	Window *uint32 `json:"window,omitempty"`
}

// ConfigurePortEthernetEfmOamLinkMonitoringErroredFrameSeconds struct
type ConfigurePortEthernetEfmOamLinkMonitoringErroredFrameSeconds struct {
	AdminState *string `json:"admin-state,omitempty"`
	// +kubebuilder:default:=true
	EventNotification *bool `json:"event-notification,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=900
	SdThreshold *uint32 `json:"sd-threshold,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=900
	// +kubebuilder:default:=1
	SfThreshold *uint32 `json:"sf-threshold,omitempty"`
	// kubebuilder:validation:Minimum=100
	// kubebuilder:validation:Maximum=9000
	// +kubebuilder:default:=600
	Window *uint32 `json:"window,omitempty"`
}

// ConfigurePortEthernetEfmOamLinkMonitoringErroredSymbols struct
type ConfigurePortEthernetEfmOamLinkMonitoringErroredSymbols struct {
	AdminState *string `json:"admin-state,omitempty"`
	// +kubebuilder:default:=true
	EventNotification *bool `json:"event-notification,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=1000000
	SdThreshold *uint32 `json:"sd-threshold,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=1000000
	// +kubebuilder:default:=1
	SfThreshold *uint32 `json:"sf-threshold,omitempty"`
	// kubebuilder:validation:Minimum=10
	// kubebuilder:validation:Maximum=600
	// +kubebuilder:default:=10
	Window *uint32 `json:"window,omitempty"`
}

// ConfigurePortEthernetEfmOamLinkMonitoringLocalSfAction struct
type ConfigurePortEthernetEfmOamLinkMonitoringLocalSfAction struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=5
	// +kubebuilder:default:=1
	EventNotificationBurst *uint32                                                                 `json:"event-notification-burst,omitempty"`
	InfoNotification       *ConfigurePortEthernetEfmOamLinkMonitoringLocalSfActionInfoNotification `json:"info-notification,omitempty"`
	LocalPortAction        *string                                                                 `json:"local-port-action,omitempty"`
}

// ConfigurePortEthernetEfmOamLinkMonitoringLocalSfActionInfoNotification struct
type ConfigurePortEthernetEfmOamLinkMonitoringLocalSfActionInfoNotification struct {
	// +kubebuilder:default:=false
	CriticalEvent *bool `json:"critical-event,omitempty"`
	// +kubebuilder:default:=false
	DyingGasp *bool `json:"dying-gasp,omitempty"`
}

// ConfigurePortEthernetEfmOamPeerRdiRx struct
type ConfigurePortEthernetEfmOamPeerRdiRx struct {
	CriticalEvent     *string `json:"critical-event,omitempty"`
	DyingGasp         *string `json:"dying-gasp,omitempty"`
	EventNotification *string `json:"event-notification,omitempty"`
	LinkFault         *string `json:"link-fault,omitempty"`
}

// ConfigurePortEthernetEgress struct
type ConfigurePortEthernetEgress struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:default:=false
	EthBnRateChanges        *bool                                                 `json:"eth-bn-rate-changes,omitempty"`
	ExpandedSecondaryShaper []*ConfigurePortEthernetEgressExpandedSecondaryShaper `json:"expanded-secondary-shaper,omitempty"`
	HsPortPoolPolicy        *string                                               `json:"hs-port-pool-policy,omitempty"`
	HsSchedulerPolicy       *ConfigurePortEthernetEgressHsSchedulerPolicy         `json:"hs-scheduler-policy,omitempty"`
	HsSecondaryShaper       []*ConfigurePortEthernetEgressHsSecondaryShaper       `json:"hs-secondary-shaper,omitempty"`
	HsmdaSchedulerPolicy    *string                                               `json:"hsmda-scheduler-policy,omitempty"`
	// +kubebuilder:default:=false
	MonitorPortScheduler *bool                                           `json:"monitor-port-scheduler,omitempty"`
	PortQosPolicy        *ConfigurePortEthernetEgressPortQosPolicy       `json:"port-qos-policy,omitempty"`
	PortSchedulerPolicy  *ConfigurePortEthernetEgressPortSchedulerPolicy `json:"port-scheduler-policy,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=400000000
	Rate *int32 `json:"rate,omitempty"`
}

// ConfigurePortEthernetEgressExpandedSecondaryShaper struct
type ConfigurePortEthernetEgressExpandedSecondaryShaper struct {
	AggregateBurst     *ConfigurePortEthernetEgressExpandedSecondaryShaperAggregateBurst `json:"aggregate-burst,omitempty"`
	ApplyGroups        *string                                                           `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                                                           `json:"apply-groups-exclude,omitempty"`
	Class              []*ConfigurePortEthernetEgressExpandedSecondaryShaperClass        `json:"class,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=8
	// +kubebuilder:default:=8
	LowBurstMaxClass    *uint32 `json:"low-burst-max-class,omitempty"`
	MonitorThreshold    *string `json:"monitor-threshold,omitempty"`
	Rate                *string `json:"rate,omitempty"`
	SecondaryShaperName *string `json:"secondary-shaper-name,omitempty"`
}

// ConfigurePortEthernetEgressExpandedSecondaryShaperAggregateBurst struct
type ConfigurePortEthernetEgressExpandedSecondaryShaperAggregateBurst struct {
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=65528
	HighBurstIncrease *int32  `json:"high-burst-increase,omitempty"`
	LowBurstLimit     *string `json:"low-burst-limit,omitempty"`
}

// ConfigurePortEthernetEgressExpandedSecondaryShaperClass struct
type ConfigurePortEthernetEgressExpandedSecondaryShaperClass struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	BurstLimit         *string `json:"burst-limit,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=8
	ClassNumber      *uint32 `json:"class-number,omitempty"`
	MonitorThreshold *string `json:"monitor-threshold,omitempty"`
	Rate             *string `json:"rate,omitempty"`
}

// ConfigurePortEthernetEgressHsSchedulerPolicy struct
type ConfigurePortEthernetEgressHsSchedulerPolicy struct {
	Overrides  *ConfigurePortEthernetEgressHsSchedulerPolicyOverrides `json:"overrides,omitempty"`
	PolicyName *string                                                `json:"policy-name,omitempty"`
}

// ConfigurePortEthernetEgressHsSchedulerPolicyOverrides struct
type ConfigurePortEthernetEgressHsSchedulerPolicyOverrides struct {
	ApplyGroups        *string                                                                 `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                                                                 `json:"apply-groups-exclude,omitempty"`
	Group              []*ConfigurePortEthernetEgressHsSchedulerPolicyOverridesGroup           `json:"group,omitempty"`
	MaxRate            *string                                                                 `json:"max-rate,omitempty"`
	SchedulingClass    []*ConfigurePortEthernetEgressHsSchedulerPolicyOverridesSchedulingClass `json:"scheduling-class,omitempty"`
}

// ConfigurePortEthernetEgressHsSchedulerPolicyOverridesGroup struct
type ConfigurePortEthernetEgressHsSchedulerPolicyOverridesGroup struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=1
	GroupId *uint32 `json:"group-id,omitempty"`
	Rate    *string `json:"rate,omitempty"`
}

// ConfigurePortEthernetEgressHsSchedulerPolicyOverridesSchedulingClass struct
type ConfigurePortEthernetEgressHsSchedulerPolicyOverridesSchedulingClass struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=6
	ClassNumber *uint32 `json:"class-number,omitempty"`
	Rate        *string `json:"rate,omitempty"`
	Weight      *string `json:"weight,omitempty"`
}

// ConfigurePortEthernetEgressHsSecondaryShaper struct
type ConfigurePortEthernetEgressHsSecondaryShaper struct {
	Aggregate           *ConfigurePortEthernetEgressHsSecondaryShaperAggregate `json:"aggregate,omitempty"`
	ApplyGroups         *string                                                `json:"apply-groups,omitempty"`
	ApplyGroupsExclude  *string                                                `json:"apply-groups-exclude,omitempty"`
	Class               []*ConfigurePortEthernetEgressHsSecondaryShaperClass   `json:"class,omitempty"`
	Description         *string                                                `json:"description,omitempty"`
	SecondaryShaperName *string                                                `json:"secondary-shaper-name,omitempty"`
}

// ConfigurePortEthernetEgressHsSecondaryShaperAggregate struct
type ConfigurePortEthernetEgressHsSecondaryShaperAggregate struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=6
	// +kubebuilder:default:=6
	LowBurstMaxClass *uint32 `json:"low-burst-max-class,omitempty"`
	// +kubebuilder:default:="max"
	Rate *string `json:"rate,omitempty"`
}

// ConfigurePortEthernetEgressHsSecondaryShaperClass struct
type ConfigurePortEthernetEgressHsSecondaryShaperClass struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=6
	ClassNumber *uint32 `json:"class-number,omitempty"`
	// +kubebuilder:default:="max"
	Rate *string `json:"rate,omitempty"`
}

// ConfigurePortEthernetEgressPortQosPolicy struct
type ConfigurePortEthernetEgressPortQosPolicy struct {
	PolicyName *string `json:"policy-name,omitempty"`
}

// ConfigurePortEthernetEgressPortSchedulerPolicy struct
type ConfigurePortEthernetEgressPortSchedulerPolicy struct {
	Overrides  *ConfigurePortEthernetEgressPortSchedulerPolicyOverrides `json:"overrides,omitempty"`
	PolicyName *string                                                  `json:"policy-name,omitempty"`
}

// ConfigurePortEthernetEgressPortSchedulerPolicyOverrides struct
type ConfigurePortEthernetEgressPortSchedulerPolicyOverrides struct {
	ApplyGroups        *string                                                         `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                                                         `json:"apply-groups-exclude,omitempty"`
	Level              []*ConfigurePortEthernetEgressPortSchedulerPolicyOverridesLevel `json:"level,omitempty"`
	MaxRate            *ConfigurePortEthernetEgressPortSchedulerPolicyOverridesMaxRate `json:"max-rate,omitempty"`
}

// ConfigurePortEthernetEgressPortSchedulerPolicyOverridesLevel struct
type ConfigurePortEthernetEgressPortSchedulerPolicyOverridesLevel struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=8
	PriorityLevel     *uint32                                                                        `json:"priority-level,omitempty"`
	RateOrPercentRate *ConfigurePortEthernetEgressPortSchedulerPolicyOverridesLevelRateOrPercentRate `json:"rate-or-percent-rate,omitempty"`
}

// ConfigurePortEthernetEgressPortSchedulerPolicyOverridesLevelRateOrPercentRate struct
type ConfigurePortEthernetEgressPortSchedulerPolicyOverridesLevelRateOrPercentRate struct {
	PercentRate *ConfigurePortEthernetEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRate `json:"percent-rate,omitempty"`
	Rate        *ConfigurePortEthernetEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRate        `json:"rate,omitempty"`
}

// ConfigurePortEthernetEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRate struct
type ConfigurePortEthernetEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRate struct {
	PercentRate *ConfigurePortEthernetEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRatePercentRate `json:"percent-rate,omitempty"`
}

// ConfigurePortEthernetEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRatePercentRate struct
type ConfigurePortEthernetEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRatePercentRate struct {
	Cir *string `json:"cir,omitempty"`
	Pir *string `json:"pir,omitempty"`
}

// ConfigurePortEthernetEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRate struct
type ConfigurePortEthernetEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRate struct {
	Rate *ConfigurePortEthernetEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRateRate `json:"rate,omitempty"`
}

// ConfigurePortEthernetEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRateRate struct
type ConfigurePortEthernetEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRateRate struct {
	Cir *string `json:"cir,omitempty"`
	Pir *string `json:"pir,omitempty"`
}

// ConfigurePortEthernetEgressPortSchedulerPolicyOverridesMaxRate struct
type ConfigurePortEthernetEgressPortSchedulerPolicyOverridesMaxRate struct {
	RateOrPercentRate *ConfigurePortEthernetEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRate `json:"rate-or-percent-rate,omitempty"`
}

// ConfigurePortEthernetEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRate struct
type ConfigurePortEthernetEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRate struct {
	PercentRate *ConfigurePortEthernetEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRatePercentRate `json:"percent-rate,omitempty"`
	Rate        *ConfigurePortEthernetEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRateRate        `json:"rate,omitempty"`
}

// ConfigurePortEthernetEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRatePercentRate struct
type ConfigurePortEthernetEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRatePercentRate struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=10000
	// +kubebuilder:default:="100"
	PercentRate *string `json:"percent-rate,omitempty"`
}

// ConfigurePortEthernetEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRateRate struct
type ConfigurePortEthernetEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRateRate struct {
	Rate *string `json:"rate,omitempty"`
}

// ConfigurePortEthernetElmi struct
type ConfigurePortEthernetElmi struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:validation:Enum=`uni-n`
	Mode *string `json:"mode,omitempty"`
	// kubebuilder:validation:Minimum=2
	// kubebuilder:validation:Maximum=10
	// +kubebuilder:default:=4
	N393 *uint32 `json:"n393,omitempty"`
	// kubebuilder:validation:Minimum=5
	// kubebuilder:validation:Maximum=30
	// +kubebuilder:default:=10
	T391 *uint32 `json:"t391,omitempty"`
	// kubebuilder:validation:Minimum=5
	// kubebuilder:validation:Maximum=30
	// +kubebuilder:default:=15
	T392 *uint32 `json:"t392,omitempty"`
}

// ConfigurePortEthernetEthCfm struct
type ConfigurePortEthernetEthCfm struct {
	Mep []*ConfigurePortEthernetEthCfmMep `json:"mep,omitempty"`
}

// ConfigurePortEthernetEthCfmMep struct
type ConfigurePortEthernetEthCfmMep struct {
	AdminState         *string                                          `json:"admin-state,omitempty"`
	Ais                *ConfigurePortEthernetEthCfmMepAis               `json:"ais,omitempty"`
	AlarmNotification  *ConfigurePortEthernetEthCfmMepAlarmNotification `json:"alarm-notification,omitempty"`
	ApplyGroups        *string                                          `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                                          `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:default:=false
	Ccm            *bool   `json:"ccm,omitempty"`
	CcmLtmPriority *string `json:"ccm-ltm-priority,omitempty"`
	// kubebuilder:validation:Minimum=3
	// kubebuilder:validation:Maximum=1500
	CcmPaddingSize *uint32 `json:"ccm-padding-size,omitempty"`
	CcmTlvIgnore   *string `json:"ccm-tlv-ignore,omitempty"`
	// +kubebuilder:default:=false
	CollectLmmStats *bool                                  `json:"collect-lmm-stats,omitempty"`
	Csf             *ConfigurePortEthernetEthCfmMepCsf     `json:"csf,omitempty"`
	Description     *string                                `json:"description,omitempty"`
	EthBn           *ConfigurePortEthernetEthCfmMepEthBn   `json:"eth-bn,omitempty"`
	EthTest         *ConfigurePortEthernetEthCfmMepEthTest `json:"eth-test,omitempty"`
	// +kubebuilder:default:=false
	FacilityFault *bool                                `json:"facility-fault,omitempty"`
	Grace         *ConfigurePortEthernetEthCfmMepGrace `json:"grace,omitempty"`
	// +kubebuilder:default:=false
	InstallMep        *bool   `json:"install-mep,omitempty"`
	LowPriorityDefect *string `json:"low-priority-defect,omitempty"`
	MaAdminName       *string `json:"ma-admin-name,omitempty"`
	MacAddress        *string `json:"mac-address,omitempty"`
	MdAdminName       *string `json:"md-admin-name,omitempty"`
	MepId             *string `json:"mep-id,omitempty"`
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=600
	// +kubebuilder:default:=3
	OneWayDelayThreshold *uint32 `json:"one-way-delay-threshold,omitempty"`
	Vlan                 *string `json:"vlan,omitempty"`
}

// ConfigurePortEthernetEthCfmMepAis struct
type ConfigurePortEthernetEthCfmMepAis struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=7
	ClientMegLevel *uint32 `json:"client-meg-level,omitempty"`
	// +kubebuilder:default:=false
	InterfaceSupport *bool `json:"interface-support,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=1
	// +kubebuilder:default:=1
	Interval *uint32 `json:"interval,omitempty"`
	// +kubebuilder:validation:Enum=`all-def`;`mac-rem-err-xcon`
	// +kubebuilder:default:="all-def"
	LowPriorityDefect *string `json:"low-priority-defect,omitempty"`
	Priority          *string `json:"priority,omitempty"`
}

// ConfigurePortEthernetEthCfmMepAlarmNotification struct
type ConfigurePortEthernetEthCfmMepAlarmNotification struct {
	// kubebuilder:validation:Minimum=250
	// kubebuilder:validation:Maximum=250
	FngAlarmTime *int32 `json:"fng-alarm-time,omitempty"`
	// kubebuilder:validation:Minimum=250
	// kubebuilder:validation:Maximum=250
	FngResetTime *int32 `json:"fng-reset-time,omitempty"`
}

// ConfigurePortEthernetEthCfmMepCsf struct
type ConfigurePortEthernetEthCfmMepCsf struct {
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=0
	// +kubebuilder:default:="3.5"
	Multiplier *string `json:"multiplier,omitempty"`
}

// ConfigurePortEthernetEthCfmMepEthBn struct
type ConfigurePortEthernetEthCfmMepEthBn struct {
	// +kubebuilder:default:=false
	Receive *bool `json:"receive,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=600
	// +kubebuilder:default:=5
	RxUpdatePacing *uint32 `json:"rx-update-pacing,omitempty"`
}

// ConfigurePortEthernetEthCfmMepEthTest struct
type ConfigurePortEthernetEthCfmMepEthTest struct {
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=11840
	// +kubebuilder:default:=1
	BitErrorThreshold *uint32                                           `json:"bit-error-threshold,omitempty"`
	TestPattern       *ConfigurePortEthernetEthCfmMepEthTestTestPattern `json:"test-pattern,omitempty"`
}

// ConfigurePortEthernetEthCfmMepEthTestTestPattern struct
type ConfigurePortEthernetEthCfmMepEthTestTestPattern struct {
	// +kubebuilder:default:=false
	CrcTlv *bool `json:"crc-tlv,omitempty"`
	// +kubebuilder:validation:Enum=`all-ones`;`all-zeros`
	// +kubebuilder:default:="all-zeros"
	Pattern *string `json:"pattern,omitempty"`
}

// ConfigurePortEthernetEthCfmMepGrace struct
type ConfigurePortEthernetEthCfmMepGrace struct {
	EthEd       *ConfigurePortEthernetEthCfmMepGraceEthEd       `json:"eth-ed,omitempty"`
	EthVsmGrace *ConfigurePortEthernetEthCfmMepGraceEthVsmGrace `json:"eth-vsm-grace,omitempty"`
}

// ConfigurePortEthernetEthCfmMepGraceEthEd struct
type ConfigurePortEthernetEthCfmMepGraceEthEd struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=86400
	MaxRxDefectWindow *uint32 `json:"max-rx-defect-window,omitempty"`
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=7
	Priority *int32 `json:"priority,omitempty"`
	// +kubebuilder:default:=true
	RxEthEd *bool `json:"rx-eth-ed,omitempty"`
	// +kubebuilder:default:=false
	TxEthEd *bool `json:"tx-eth-ed,omitempty"`
}

// ConfigurePortEthernetEthCfmMepGraceEthVsmGrace struct
type ConfigurePortEthernetEthCfmMepGraceEthVsmGrace struct {
	// +kubebuilder:default:=true
	RxEthVsmGrace *bool `json:"rx-eth-vsm-grace,omitempty"`
	// +kubebuilder:default:=true
	TxEthVsmGrace *bool `json:"tx-eth-vsm-grace,omitempty"`
}

// ConfigurePortEthernetHoldTime struct
type ConfigurePortEthernetHoldTime struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=3600000
	Down *uint32 `json:"down,omitempty"`
	// +kubebuilder:validation:Enum=`centiseconds`;`seconds`
	// +kubebuilder:default:="seconds"
	Units *string `json:"units,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=3600000
	Up *uint32 `json:"up,omitempty"`
}

// ConfigurePortEthernetHsmdaSchedulerOverrides struct
type ConfigurePortEthernetHsmdaSchedulerOverrides struct {
	ApplyGroups        *string                                                        `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                                                        `json:"apply-groups-exclude,omitempty"`
	Group              []*ConfigurePortEthernetHsmdaSchedulerOverridesGroup           `json:"group,omitempty"`
	MaxRate            *string                                                        `json:"max-rate,omitempty"`
	SchedulingClass    []*ConfigurePortEthernetHsmdaSchedulerOverridesSchedulingClass `json:"scheduling-class,omitempty"`
}

// ConfigurePortEthernetHsmdaSchedulerOverridesGroup struct
type ConfigurePortEthernetHsmdaSchedulerOverridesGroup struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=2
	GroupId *uint32 `json:"group-id,omitempty"`
	Rate    *string `json:"rate,omitempty"`
}

// ConfigurePortEthernetHsmdaSchedulerOverridesSchedulingClass struct
type ConfigurePortEthernetHsmdaSchedulerOverridesSchedulingClass struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=8
	ClassNumber *uint32 `json:"class-number,omitempty"`
	Rate        *string `json:"rate,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=100
	WeightInGroup *int32 `json:"weight-in-group,omitempty"`
}

// ConfigurePortEthernetIngress struct
type ConfigurePortEthernetIngress struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=400000
	Rate *int32 `json:"rate,omitempty"`
}

// ConfigurePortEthernetLldp struct
type ConfigurePortEthernetLldp struct {
	ApplyGroups        *string                             `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                             `json:"apply-groups-exclude,omitempty"`
	DestMac            []*ConfigurePortEthernetLldpDestMac `json:"dest-mac,omitempty"`
}

// ConfigurePortEthernetLldpDestMac struct
type ConfigurePortEthernetLldpDestMac struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	MacType            *string `json:"mac-type,omitempty"`
	// +kubebuilder:default:=false
	Notification *bool `json:"notification,omitempty"`
	// +kubebuilder:validation:Enum=`tx-if-alias`;`tx-if-name`;`tx-local`
	// +kubebuilder:default:="tx-local"
	PortIdSubtype *string `json:"port-id-subtype,omitempty"`
	// +kubebuilder:default:=false
	Receive *bool `json:"receive,omitempty"`
	// +kubebuilder:default:=false
	Transmit *bool `json:"transmit,omitempty"`
	// +kubebuilder:default:=false
	TunnelNearestBridge *bool                                            `json:"tunnel-nearest-bridge,omitempty"`
	TxMgmtAddress       []*ConfigurePortEthernetLldpDestMacTxMgmtAddress `json:"tx-mgmt-address,omitempty"`
	TxTlvs              *ConfigurePortEthernetLldpDestMacTxTlvs          `json:"tx-tlvs,omitempty"`
}

// ConfigurePortEthernetLldpDestMacTxMgmtAddress struct
type ConfigurePortEthernetLldpDestMacTxMgmtAddress struct {
	AdminState         *string `json:"admin-state,omitempty"`
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:validation:Enum=`oob`;`oob-ipv6`;`system`;`system-ipv6`
	MgmtAddressSystemType *string `json:"mgmt-address-system-type,omitempty"`
}

// ConfigurePortEthernetLldpDestMacTxTlvs struct
type ConfigurePortEthernetLldpDestMacTxTlvs struct {
	// +kubebuilder:default:=false
	PortDesc *bool `json:"port-desc,omitempty"`
	// +kubebuilder:default:=false
	SysCap *bool `json:"sys-cap,omitempty"`
	// +kubebuilder:default:=false
	SysDesc *bool `json:"sys-desc,omitempty"`
	// +kubebuilder:default:=false
	SysName *bool `json:"sys-name,omitempty"`
}

// ConfigurePortEthernetLoopback struct
type ConfigurePortEthernetLoopback struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:validation:Enum=`internal`;`line`
	Direction *string `json:"direction,omitempty"`
	// +kubebuilder:default:=false
	SwapSrcDstMac *bool `json:"swap-src-dst-mac,omitempty"`
}

// ConfigurePortEthernetNetwork struct
type ConfigurePortEthernetNetwork struct {
	AccountingPolicy   *string `json:"accounting-policy,omitempty"`
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:default:=false
	CollectStats *bool                               `json:"collect-stats,omitempty"`
	Egress       *ConfigurePortEthernetNetworkEgress `json:"egress,omitempty"`
}

// ConfigurePortEthernetNetworkEgress struct
type ConfigurePortEthernetNetworkEgress struct {
	PortQueues  *ConfigurePortEthernetNetworkEgressPortQueues   `json:"port-queues,omitempty"`
	QueueGroup  []*ConfigurePortEthernetNetworkEgressQueueGroup `json:"queue-group,omitempty"`
	QueuePolicy *string                                         `json:"queue-policy,omitempty"`
}

// ConfigurePortEthernetNetworkEgressPortQueues struct
type ConfigurePortEthernetNetworkEgressPortQueues struct {
	Overrides *ConfigurePortEthernetNetworkEgressPortQueuesOverrides `json:"overrides,omitempty"`
}

// ConfigurePortEthernetNetworkEgressPortQueuesOverrides struct
type ConfigurePortEthernetNetworkEgressPortQueuesOverrides struct {
	Queue []*ConfigurePortEthernetNetworkEgressPortQueuesOverridesQueue `json:"queue,omitempty"`
}

// ConfigurePortEthernetNetworkEgressPortQueuesOverridesQueue struct
type ConfigurePortEthernetNetworkEgressPortQueuesOverridesQueue struct {
	ApplyGroups        *string                                                                      `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                                                                      `json:"apply-groups-exclude,omitempty"`
	MonitorQueueDepth  *ConfigurePortEthernetNetworkEgressPortQueuesOverridesQueueMonitorQueueDepth `json:"monitor-queue-depth,omitempty"`
	QueueId            *string                                                                      `json:"queue-id,omitempty"`
}

// ConfigurePortEthernetNetworkEgressPortQueuesOverridesQueueMonitorQueueDepth struct
type ConfigurePortEthernetNetworkEgressPortQueuesOverridesQueueMonitorQueueDepth struct {
	// +kubebuilder:default:=false
	FastPolling *bool `json:"fast-polling,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=9999
	ViolationThreshold *string `json:"violation-threshold,omitempty"`
}

// ConfigurePortEthernetNetworkEgressQueueGroup struct
type ConfigurePortEthernetNetworkEgressQueueGroup struct {
	AccountingPolicy   *string                                                    `json:"accounting-policy,omitempty"`
	AggregateRate      *ConfigurePortEthernetNetworkEgressQueueGroupAggregateRate `json:"aggregate-rate,omitempty"`
	ApplyGroups        *string                                                    `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                                                    `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:default:=false
	CollectStats *bool   `json:"collect-stats,omitempty"`
	Description  *string `json:"description,omitempty"`
	// +kubebuilder:default:=false
	HsTurbo *bool `json:"hs-turbo,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=65535
	InstanceId           *uint16                                                     `json:"instance-id,omitempty"`
	PolicerControlPolicy *string                                                     `json:"policer-control-policy,omitempty"`
	QueueGroupName       *string                                                     `json:"queue-group-name,omitempty"`
	QueueOverrides       *ConfigurePortEthernetNetworkEgressQueueGroupQueueOverrides `json:"queue-overrides,omitempty"`
	SchedulerPolicy      *string                                                     `json:"scheduler-policy,omitempty"`
}

// ConfigurePortEthernetNetworkEgressQueueGroupAggregateRate struct
type ConfigurePortEthernetNetworkEgressQueueGroupAggregateRate struct {
	// +kubebuilder:default:=false
	LimitUnusedBandwidth *bool `json:"limit-unused-bandwidth,omitempty"`
	// +kubebuilder:default:=false
	QueueFrameBasedAccounting *bool `json:"queue-frame-based-accounting,omitempty"`
	// +kubebuilder:default:="max"
	Rate *string `json:"rate,omitempty"`
}

// ConfigurePortEthernetNetworkEgressQueueGroupQueueOverrides struct
type ConfigurePortEthernetNetworkEgressQueueGroupQueueOverrides struct {
	Queue []*ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueue `json:"queue,omitempty"`
}

// ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueue struct
type ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueue struct {
	AdaptationRule     *ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueAdaptationRule `json:"adaptation-rule,omitempty"`
	ApplyGroups        *string                                                                        `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                                                                        `json:"apply-groups-exclude,omitempty"`
	Cbs                *string                                                                        `json:"cbs,omitempty"`
	DropTail           *ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueDropTail       `json:"drop-tail,omitempty"`
	Mbs                *string                                                                        `json:"mbs,omitempty"`
	// +kubebuilder:default:=false
	MonitorDepth      *bool                                                                             `json:"monitor-depth,omitempty"`
	MonitorQueueDepth *ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueMonitorQueueDepth `json:"monitor-queue-depth,omitempty"`
	QueueId           *string                                                                           `json:"queue-id,omitempty"`
	QueueOverrideRate *ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueQueueOverrideRate `json:"queue-override-rate,omitempty"`
}

// ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueAdaptationRule struct
type ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueAdaptationRule struct {
	Cir *string `json:"cir,omitempty"`
	Pir *string `json:"pir,omitempty"`
}

// ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueDropTail struct
type ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueDropTail struct {
	Low *ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueDropTailLow `json:"low,omitempty"`
}

// ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueDropTailLow struct
type ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueDropTailLow struct {
	PercentReductionFromMbs *string `json:"percent-reduction-from-mbs,omitempty"`
}

// ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueMonitorQueueDepth struct
type ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueMonitorQueueDepth struct {
	// +kubebuilder:default:=false
	FastPolling *bool `json:"fast-polling,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=9999
	ViolationThreshold *string `json:"violation-threshold,omitempty"`
}

// ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueQueueOverrideRate struct
type ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueQueueOverrideRate struct {
	PercentRate *ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueQueueOverrideRatePercentRate `json:"percent-rate,omitempty"`
	Rate        *ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueQueueOverrideRateRate        `json:"rate,omitempty"`
}

// ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueQueueOverrideRatePercentRate struct
type ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueQueueOverrideRatePercentRate struct {
	PercentRate *ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueQueueOverrideRatePercentRatePercentRate `json:"percent-rate,omitempty"`
}

// ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueQueueOverrideRatePercentRatePercentRate struct
type ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueQueueOverrideRatePercentRatePercentRate struct {
	Cir *string `json:"cir,omitempty"`
	Pir *string `json:"pir,omitempty"`
}

// ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueQueueOverrideRateRate struct
type ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueQueueOverrideRateRate struct {
	Rate *ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueQueueOverrideRateRateRate `json:"rate,omitempty"`
}

// ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueQueueOverrideRateRateRate struct
type ConfigurePortEthernetNetworkEgressQueueGroupQueueOverridesQueueQueueOverrideRateRateRate struct {
	Cir *string `json:"cir,omitempty"`
	Pir *string `json:"pir,omitempty"`
}

// ConfigurePortEthernetReportAlarm struct
type ConfigurePortEthernetReportAlarm struct {
	AlignmentMarkerNotLocked *bool `json:"alignment-marker-not-locked,omitempty"`
	BlockNotLocked           *bool `json:"block-not-locked,omitempty"`
	DuplicateLane            *bool `json:"duplicate-lane,omitempty"`
	FrameNotLocked           *bool `json:"frame-not-locked,omitempty"`
	HighBer                  *bool `json:"high-ber,omitempty"`
	Local                    *bool `json:"local,omitempty"`
	Remote                   *bool `json:"remote,omitempty"`
	SignalFail               *bool `json:"signal-fail,omitempty"`
}

// ConfigurePortEthernetSsm struct
type ConfigurePortEthernetSsm struct {
	AdminState         *string `json:"admin-state,omitempty"`
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:validation:Enum=`sdh`;`sonet`
	// +kubebuilder:default:="sdh"
	CodeType *string `json:"code-type,omitempty"`
	// +kubebuilder:default:=false
	EsmcTunnel *bool `json:"esmc-tunnel,omitempty"`
	TxDus      *bool `json:"tx-dus,omitempty"`
}

// ConfigurePortEthernetSymbolMonitor struct
type ConfigurePortEthernetSymbolMonitor struct {
	AdminState         *string                                          `json:"admin-state,omitempty"`
	ApplyGroups        *string                                          `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                                          `json:"apply-groups-exclude,omitempty"`
	SignalDegrade      *ConfigurePortEthernetSymbolMonitorSignalDegrade `json:"signal-degrade,omitempty"`
	SignalFailure      *ConfigurePortEthernetSymbolMonitorSignalFailure `json:"signal-failure,omitempty"`
	// kubebuilder:validation:Minimum=5
	// kubebuilder:validation:Maximum=60
	// +kubebuilder:default:=10
	WindowSize *uint32 `json:"window-size,omitempty"`
}

// ConfigurePortEthernetSymbolMonitorSignalDegrade struct
type ConfigurePortEthernetSymbolMonitorSignalDegrade struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=9
	// +kubebuilder:default:=1
	Multiplier *uint32 `json:"multiplier,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=9
	Threshold *uint32 `json:"threshold,omitempty"`
}

// ConfigurePortEthernetSymbolMonitorSignalFailure struct
type ConfigurePortEthernetSymbolMonitorSignalFailure struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=9
	// +kubebuilder:default:=1
	Multiplier *uint32 `json:"multiplier,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=9
	Threshold *uint32 `json:"threshold,omitempty"`
}

// ConfigurePortGnss struct
type ConfigurePortGnss struct {
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=32767
	// +kubebuilder:default:=0
	AntennaCableDelay  *uint32                         `json:"antenna-cable-delay,omitempty"`
	ApplyGroups        *string                         `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                         `json:"apply-groups-exclude,omitempty"`
	Constellation      *ConfigurePortGnssConstellation `json:"constellation,omitempty"`
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=89
	// +kubebuilder:default:=10
	ElevationMaskAngle *uint32 `json:"elevation-mask-angle,omitempty"`
}

// ConfigurePortGnssConstellation struct
type ConfigurePortGnssConstellation struct {
	// +kubebuilder:default:=false
	Glonass *bool `json:"glonass,omitempty"`
	// +kubebuilder:default:=true
	Gps *bool `json:"gps,omitempty"`
}

// ConfigurePortHybridBufferAllocation struct
type ConfigurePortHybridBufferAllocation struct {
	EgressWeight  *ConfigurePortHybridBufferAllocationEgressWeight  `json:"egress-weight,omitempty"`
	IngressWeight *ConfigurePortHybridBufferAllocationIngressWeight `json:"ingress-weight,omitempty"`
}

// ConfigurePortHybridBufferAllocationEgressWeight struct
type ConfigurePortHybridBufferAllocationEgressWeight struct {
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=100
	// +kubebuilder:default:=50
	Access *uint32 `json:"access,omitempty"`
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=100
	// +kubebuilder:default:=50
	Network *uint32 `json:"network,omitempty"`
}

// ConfigurePortHybridBufferAllocationIngressWeight struct
type ConfigurePortHybridBufferAllocationIngressWeight struct {
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=100
	// +kubebuilder:default:=50
	Access *uint32 `json:"access,omitempty"`
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=100
	// +kubebuilder:default:=50
	Network *uint32 `json:"network,omitempty"`
}

// ConfigurePortModifyBufferAllocation struct
type ConfigurePortModifyBufferAllocation struct {
	PercentageOfRate *ConfigurePortModifyBufferAllocationPercentageOfRate `json:"percentage-of-rate,omitempty"`
}

// ConfigurePortModifyBufferAllocationPercentageOfRate struct
type ConfigurePortModifyBufferAllocationPercentageOfRate struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=1000
	Egress *uint32 `json:"egress,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=1000
	// +kubebuilder:default:=100
	Ingress *uint32 `json:"ingress,omitempty"`
}

// ConfigurePortNetwork struct
type ConfigurePortNetwork struct {
	ApplyGroups        *string                     `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                     `json:"apply-groups-exclude,omitempty"`
	Egress             *ConfigurePortNetworkEgress `json:"egress,omitempty"`
}

// ConfigurePortNetworkEgress struct
type ConfigurePortNetworkEgress struct {
	Pool []*ConfigurePortNetworkEgressPool `json:"pool,omitempty"`
}

// ConfigurePortNetworkEgressPool struct
type ConfigurePortNetworkEgressPool struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=1000
	AmberAlarmThreshold *uint32 `json:"amber-alarm-threshold,omitempty"`
	ApplyGroups         *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude  *string `json:"apply-groups-exclude,omitempty"`
	Name                *string `json:"name,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=1000
	RedAlarmThreshold *uint32                                `json:"red-alarm-threshold,omitempty"`
	ResvCbs           *ConfigurePortNetworkEgressPoolResvCbs `json:"resv-cbs,omitempty"`
	SlopePolicy       *string                                `json:"slope-policy,omitempty"`
}

// ConfigurePortNetworkEgressPoolResvCbs struct
type ConfigurePortNetworkEgressPoolResvCbs struct {
	AmberAlarmAction *ConfigurePortNetworkEgressPoolResvCbsAmberAlarmAction `json:"amber-alarm-action,omitempty"`
	// +kubebuilder:default:="auto"
	Cbs *string `json:"cbs,omitempty"`
}

// ConfigurePortNetworkEgressPoolResvCbsAmberAlarmAction struct
type ConfigurePortNetworkEgressPoolResvCbsAmberAlarmAction struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=100
	Max *uint32 `json:"max,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=100
	Step *uint32 `json:"step,omitempty"`
}

// ConfigurePortOtu struct
type ConfigurePortOtu struct {
	ApplyGroups                *string                                     `json:"apply-groups,omitempty"`
	ApplyGroupsExclude         *string                                     `json:"apply-groups-exclude,omitempty"`
	AsyncMapping               *bool                                       `json:"async-mapping,omitempty"`
	Fec                        *string                                     `json:"fec,omitempty"`
	FineGranularityBer         *ConfigurePortOtuFineGranularityBer         `json:"fine-granularity-ber,omitempty"`
	Otu2LanDataRate            *string                                     `json:"otu2-lan-data-rate,omitempty"`
	PathMonitoring             *ConfigurePortOtuPathMonitoring             `json:"path-monitoring,omitempty"`
	PayloadStructureIdentifier *ConfigurePortOtuPayloadStructureIdentifier `json:"payload-structure-identifier,omitempty"`
	ReportAlarm                *ConfigurePortOtuReportAlarm                `json:"report-alarm,omitempty"`
	// kubebuilder:validation:Minimum=5
	// kubebuilder:validation:Maximum=9
	// +kubebuilder:default:=7
	SdThreshold       *uint32                            `json:"sd-threshold,omitempty"`
	SectionMonitoring *ConfigurePortOtuSectionMonitoring `json:"section-monitoring,omitempty"`
	SfSdMethod        *string                            `json:"sf-sd-method,omitempty"`
	// kubebuilder:validation:Minimum=3
	// kubebuilder:validation:Maximum=6
	// +kubebuilder:default:=5
	SfThreshold *uint32 `json:"sf-threshold,omitempty"`
}

// ConfigurePortOtuFineGranularityBer struct
type ConfigurePortOtuFineGranularityBer struct {
	SignalDegrade *ConfigurePortOtuFineGranularityBerSignalDegrade `json:"signal-degrade,omitempty"`
	SignalFailure *ConfigurePortOtuFineGranularityBerSignalFailure `json:"signal-failure,omitempty"`
}

// ConfigurePortOtuFineGranularityBerSignalDegrade struct
type ConfigurePortOtuFineGranularityBerSignalDegrade struct {
	Clear *ConfigurePortOtuFineGranularityBerSignalDegradeClear `json:"clear,omitempty"`
	Raise *ConfigurePortOtuFineGranularityBerSignalDegradeRaise `json:"raise,omitempty"`
}

// ConfigurePortOtuFineGranularityBerSignalDegradeClear struct
type ConfigurePortOtuFineGranularityBerSignalDegradeClear struct {
	// kubebuilder:validation:Minimum=10
	// kubebuilder:validation:Maximum=99
	// +kubebuilder:default:=10
	Multiplier *uint32 `json:"multiplier,omitempty"`
	// kubebuilder:validation:Minimum=3
	// kubebuilder:validation:Maximum=10
	// +kubebuilder:default:=8
	Threshold *uint32 `json:"threshold,omitempty"`
}

// ConfigurePortOtuFineGranularityBerSignalDegradeRaise struct
type ConfigurePortOtuFineGranularityBerSignalDegradeRaise struct {
	// kubebuilder:validation:Minimum=10
	// kubebuilder:validation:Maximum=99
	// +kubebuilder:default:=10
	Multiplier *uint32 `json:"multiplier,omitempty"`
	// kubebuilder:validation:Minimum=3
	// kubebuilder:validation:Maximum=9
	// +kubebuilder:default:=7
	Threshold *uint32 `json:"threshold,omitempty"`
}

// ConfigurePortOtuFineGranularityBerSignalFailure struct
type ConfigurePortOtuFineGranularityBerSignalFailure struct {
	Clear *ConfigurePortOtuFineGranularityBerSignalFailureClear `json:"clear,omitempty"`
	Raise *ConfigurePortOtuFineGranularityBerSignalFailureRaise `json:"raise,omitempty"`
}

// ConfigurePortOtuFineGranularityBerSignalFailureClear struct
type ConfigurePortOtuFineGranularityBerSignalFailureClear struct {
	// kubebuilder:validation:Minimum=10
	// kubebuilder:validation:Maximum=99
	// +kubebuilder:default:=10
	Multiplier *uint32 `json:"multiplier,omitempty"`
	// kubebuilder:validation:Minimum=3
	// kubebuilder:validation:Maximum=9
	// +kubebuilder:default:=6
	Threshold *uint32 `json:"threshold,omitempty"`
}

// ConfigurePortOtuFineGranularityBerSignalFailureRaise struct
type ConfigurePortOtuFineGranularityBerSignalFailureRaise struct {
	// kubebuilder:validation:Minimum=10
	// kubebuilder:validation:Maximum=99
	// +kubebuilder:default:=10
	Multiplier *uint32 `json:"multiplier,omitempty"`
	// kubebuilder:validation:Minimum=3
	// kubebuilder:validation:Maximum=8
	// +kubebuilder:default:=5
	Threshold *uint32 `json:"threshold,omitempty"`
}

// ConfigurePortOtuPathMonitoring struct
type ConfigurePortOtuPathMonitoring struct {
	TrailTraceIdentifier *ConfigurePortOtuPathMonitoringTrailTraceIdentifier `json:"trail-trace-identifier,omitempty"`
}

// ConfigurePortOtuPathMonitoringTrailTraceIdentifier struct
type ConfigurePortOtuPathMonitoringTrailTraceIdentifier struct {
	Expected         *ConfigurePortOtuPathMonitoringTrailTraceIdentifierExpected `json:"expected,omitempty"`
	MismatchReaction *string                                                     `json:"mismatch-reaction,omitempty"`
	Transmit         *ConfigurePortOtuPathMonitoringTrailTraceIdentifierTransmit `json:"transmit,omitempty"`
}

// ConfigurePortOtuPathMonitoringTrailTraceIdentifierExpected struct
type ConfigurePortOtuPathMonitoringTrailTraceIdentifierExpected struct {
	Expected *ConfigurePortOtuPathMonitoringTrailTraceIdentifierExpectedExpected `json:"expected,omitempty"`
}

// ConfigurePortOtuPathMonitoringTrailTraceIdentifierExpectedExpected struct
type ConfigurePortOtuPathMonitoringTrailTraceIdentifierExpectedExpected struct {
	AutoGenerated *ConfigurePortOtuPathMonitoringTrailTraceIdentifierExpectedExpectedAutoGenerated `json:"auto-generated,omitempty"`
	Bytes         *ConfigurePortOtuPathMonitoringTrailTraceIdentifierExpectedExpectedBytes         `json:"bytes,omitempty"`
	String        *ConfigurePortOtuPathMonitoringTrailTraceIdentifierExpectedExpectedString        `json:"string,omitempty"`
}

// ConfigurePortOtuPathMonitoringTrailTraceIdentifierExpectedExpectedAutoGenerated struct
type ConfigurePortOtuPathMonitoringTrailTraceIdentifierExpectedExpectedAutoGenerated struct {
	AutoGenerated *string `json:"auto-generated,omitempty"`
}

// ConfigurePortOtuPathMonitoringTrailTraceIdentifierExpectedExpectedBytes struct
type ConfigurePortOtuPathMonitoringTrailTraceIdentifierExpectedExpectedBytes struct {
	// kubebuilder:validation:MinLength=0
	// kubebuilder:validation:MaxLength=192
	Bytes *string `json:"bytes,omitempty"`
}

// ConfigurePortOtuPathMonitoringTrailTraceIdentifierExpectedExpectedString struct
type ConfigurePortOtuPathMonitoringTrailTraceIdentifierExpectedExpectedString struct {
	// kubebuilder:validation:MinLength=0
	// kubebuilder:validation:MaxLength=64
	String *string `json:"string,omitempty"`
}

// ConfigurePortOtuPathMonitoringTrailTraceIdentifierTransmit struct
type ConfigurePortOtuPathMonitoringTrailTraceIdentifierTransmit struct {
	Transmit *ConfigurePortOtuPathMonitoringTrailTraceIdentifierTransmitTransmit `json:"transmit,omitempty"`
}

// ConfigurePortOtuPathMonitoringTrailTraceIdentifierTransmitTransmit struct
type ConfigurePortOtuPathMonitoringTrailTraceIdentifierTransmitTransmit struct {
	AutoGenerated *ConfigurePortOtuPathMonitoringTrailTraceIdentifierTransmitTransmitAutoGenerated `json:"auto-generated,omitempty"`
	Bytes         *ConfigurePortOtuPathMonitoringTrailTraceIdentifierTransmitTransmitBytes         `json:"bytes,omitempty"`
	String        *ConfigurePortOtuPathMonitoringTrailTraceIdentifierTransmitTransmitString        `json:"string,omitempty"`
}

// ConfigurePortOtuPathMonitoringTrailTraceIdentifierTransmitTransmitAutoGenerated struct
type ConfigurePortOtuPathMonitoringTrailTraceIdentifierTransmitTransmitAutoGenerated struct {
	AutoGenerated *string `json:"auto-generated,omitempty"`
}

// ConfigurePortOtuPathMonitoringTrailTraceIdentifierTransmitTransmitBytes struct
type ConfigurePortOtuPathMonitoringTrailTraceIdentifierTransmitTransmitBytes struct {
	// kubebuilder:validation:MinLength=0
	// kubebuilder:validation:MaxLength=192
	Bytes *string `json:"bytes,omitempty"`
}

// ConfigurePortOtuPathMonitoringTrailTraceIdentifierTransmitTransmitString struct
type ConfigurePortOtuPathMonitoringTrailTraceIdentifierTransmitTransmitString struct {
	// kubebuilder:validation:MinLength=0
	// kubebuilder:validation:MaxLength=64
	String *string `json:"string,omitempty"`
}

// ConfigurePortOtuPayloadStructureIdentifier struct
type ConfigurePortOtuPayloadStructureIdentifier struct {
	Payload *ConfigurePortOtuPayloadStructureIdentifierPayload `json:"payload,omitempty"`
}

// ConfigurePortOtuPayloadStructureIdentifierPayload struct
type ConfigurePortOtuPayloadStructureIdentifierPayload struct {
	Expected         *string `json:"expected,omitempty"`
	MismatchReaction *string `json:"mismatch-reaction,omitempty"`
	Transmit         *string `json:"transmit,omitempty"`
}

// ConfigurePortOtuReportAlarm struct
type ConfigurePortOtuReportAlarm struct {
	// +kubebuilder:default:=false
	FecFail *bool `json:"fec-fail,omitempty"`
	// +kubebuilder:default:=false
	FecSd *bool `json:"fec-sd,omitempty"`
	// +kubebuilder:default:=true
	FecSf *bool `json:"fec-sf,omitempty"`
	// +kubebuilder:default:=false
	FecUncorr *bool `json:"fec-uncorr,omitempty"`
	// +kubebuilder:default:=true
	Loc *bool `json:"loc,omitempty"`
	// +kubebuilder:default:=true
	Lof *bool `json:"lof,omitempty"`
	// +kubebuilder:default:=true
	Lom *bool `json:"lom,omitempty"`
	// +kubebuilder:default:=true
	Los *bool `json:"los,omitempty"`
	// +kubebuilder:default:=false
	OduAis *bool `json:"odu-ais,omitempty"`
	// +kubebuilder:default:=false
	OduBdi *bool `json:"odu-bdi,omitempty"`
	// +kubebuilder:default:=false
	OduLck *bool `json:"odu-lck,omitempty"`
	// +kubebuilder:default:=false
	OduOci *bool `json:"odu-oci,omitempty"`
	// +kubebuilder:default:=false
	OduTim *bool `json:"odu-tim,omitempty"`
	// +kubebuilder:default:=false
	OpuPlm *bool `json:"opu-plm,omitempty"`
	// +kubebuilder:default:=false
	OtuAis *bool `json:"otu-ais,omitempty"`
	// +kubebuilder:default:=true
	OtuBdi *bool `json:"otu-bdi,omitempty"`
	// +kubebuilder:default:=false
	OtuBerSd *bool `json:"otu-ber-sd,omitempty"`
	// +kubebuilder:default:=true
	OtuBerSf *bool `json:"otu-ber-sf,omitempty"`
	// +kubebuilder:default:=false
	OtuBiae *bool `json:"otu-biae,omitempty"`
	// +kubebuilder:default:=false
	OtuIae *bool `json:"otu-iae,omitempty"`
	// +kubebuilder:default:=false
	OtuTim *bool `json:"otu-tim,omitempty"`
}

// ConfigurePortOtuSectionMonitoring struct
type ConfigurePortOtuSectionMonitoring struct {
	TrailTraceIdentifier *ConfigurePortOtuSectionMonitoringTrailTraceIdentifier `json:"trail-trace-identifier,omitempty"`
}

// ConfigurePortOtuSectionMonitoringTrailTraceIdentifier struct
type ConfigurePortOtuSectionMonitoringTrailTraceIdentifier struct {
	Expected         *ConfigurePortOtuSectionMonitoringTrailTraceIdentifierExpected `json:"expected,omitempty"`
	MismatchReaction *string                                                        `json:"mismatch-reaction,omitempty"`
	Transmit         *ConfigurePortOtuSectionMonitoringTrailTraceIdentifierTransmit `json:"transmit,omitempty"`
}

// ConfigurePortOtuSectionMonitoringTrailTraceIdentifierExpected struct
type ConfigurePortOtuSectionMonitoringTrailTraceIdentifierExpected struct {
	Expected *ConfigurePortOtuSectionMonitoringTrailTraceIdentifierExpectedExpected `json:"expected,omitempty"`
}

// ConfigurePortOtuSectionMonitoringTrailTraceIdentifierExpectedExpected struct
type ConfigurePortOtuSectionMonitoringTrailTraceIdentifierExpectedExpected struct {
	AutoGenerated *ConfigurePortOtuSectionMonitoringTrailTraceIdentifierExpectedExpectedAutoGenerated `json:"auto-generated,omitempty"`
	Bytes         *ConfigurePortOtuSectionMonitoringTrailTraceIdentifierExpectedExpectedBytes         `json:"bytes,omitempty"`
	String        *ConfigurePortOtuSectionMonitoringTrailTraceIdentifierExpectedExpectedString        `json:"string,omitempty"`
}

// ConfigurePortOtuSectionMonitoringTrailTraceIdentifierExpectedExpectedAutoGenerated struct
type ConfigurePortOtuSectionMonitoringTrailTraceIdentifierExpectedExpectedAutoGenerated struct {
	AutoGenerated *string `json:"auto-generated,omitempty"`
}

// ConfigurePortOtuSectionMonitoringTrailTraceIdentifierExpectedExpectedBytes struct
type ConfigurePortOtuSectionMonitoringTrailTraceIdentifierExpectedExpectedBytes struct {
	// kubebuilder:validation:MinLength=0
	// kubebuilder:validation:MaxLength=192
	Bytes *string `json:"bytes,omitempty"`
}

// ConfigurePortOtuSectionMonitoringTrailTraceIdentifierExpectedExpectedString struct
type ConfigurePortOtuSectionMonitoringTrailTraceIdentifierExpectedExpectedString struct {
	// kubebuilder:validation:MinLength=0
	// kubebuilder:validation:MaxLength=64
	String *string `json:"string,omitempty"`
}

// ConfigurePortOtuSectionMonitoringTrailTraceIdentifierTransmit struct
type ConfigurePortOtuSectionMonitoringTrailTraceIdentifierTransmit struct {
	Transmit *ConfigurePortOtuSectionMonitoringTrailTraceIdentifierTransmitTransmit `json:"transmit,omitempty"`
}

// ConfigurePortOtuSectionMonitoringTrailTraceIdentifierTransmitTransmit struct
type ConfigurePortOtuSectionMonitoringTrailTraceIdentifierTransmitTransmit struct {
	AutoGenerated *ConfigurePortOtuSectionMonitoringTrailTraceIdentifierTransmitTransmitAutoGenerated `json:"auto-generated,omitempty"`
	Bytes         *ConfigurePortOtuSectionMonitoringTrailTraceIdentifierTransmitTransmitBytes         `json:"bytes,omitempty"`
	String        *ConfigurePortOtuSectionMonitoringTrailTraceIdentifierTransmitTransmitString        `json:"string,omitempty"`
}

// ConfigurePortOtuSectionMonitoringTrailTraceIdentifierTransmitTransmitAutoGenerated struct
type ConfigurePortOtuSectionMonitoringTrailTraceIdentifierTransmitTransmitAutoGenerated struct {
	AutoGenerated *string `json:"auto-generated,omitempty"`
}

// ConfigurePortOtuSectionMonitoringTrailTraceIdentifierTransmitTransmitBytes struct
type ConfigurePortOtuSectionMonitoringTrailTraceIdentifierTransmitTransmitBytes struct {
	// kubebuilder:validation:MinLength=0
	// kubebuilder:validation:MaxLength=192
	Bytes *string `json:"bytes,omitempty"`
}

// ConfigurePortOtuSectionMonitoringTrailTraceIdentifierTransmitTransmitString struct
type ConfigurePortOtuSectionMonitoringTrailTraceIdentifierTransmitTransmitString struct {
	// kubebuilder:validation:MinLength=0
	// kubebuilder:validation:MaxLength=64
	String *string `json:"string,omitempty"`
}

// ConfigurePortSonetSdh struct
type ConfigurePortSonetSdh struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:validation:Enum=`loop-timed`;`node-timed`
	ClockSource *string `json:"clock-source,omitempty"`
	// +kubebuilder:validation:Enum=`sdh`;`sonet`
	// +kubebuilder:default:="sonet"
	Framing  *string                        `json:"framing,omitempty"`
	Group    []*ConfigurePortSonetSdhGroup  `json:"group,omitempty"`
	HoldTime *ConfigurePortSonetSdhHoldTime `json:"hold-time,omitempty"`
	// +kubebuilder:validation:Enum=`internal`;`line`
	Loopback    *string                           `json:"loopback,omitempty"`
	Path        []*ConfigurePortSonetSdhPath      `json:"path,omitempty"`
	ReportAlarm *ConfigurePortSonetSdhReportAlarm `json:"report-alarm,omitempty"`
	// +kubebuilder:default:=false
	ResetPortOnPathDown *bool `json:"reset-port-on-path-down,omitempty"`
	// kubebuilder:validation:Minimum=3
	// kubebuilder:validation:Maximum=9
	// +kubebuilder:default:=6
	SdThreshold  *uint32                            `json:"sd-threshold,omitempty"`
	SectionTrace *ConfigurePortSonetSdhSectionTrace `json:"section-trace,omitempty"`
	// kubebuilder:validation:Minimum=3
	// kubebuilder:validation:Maximum=6
	// +kubebuilder:default:=3
	SfThreshold *uint32 `json:"sf-threshold,omitempty"`
	// +kubebuilder:default:=false
	SingleFiber *bool `json:"single-fiber,omitempty"`
	// +kubebuilder:validation:Enum=`oc1`;`oc12`;`oc192`;`oc3`;`oc48`;`oc768`
	Speed *string `json:"speed,omitempty"`
	// +kubebuilder:default:=false
	SuppressLowOrderAlarms *bool `json:"suppress-low-order-alarms,omitempty"`
	// +kubebuilder:default:=false
	TxDus *bool `json:"tx-dus,omitempty"`
}

// ConfigurePortSonetSdhGroup struct
type ConfigurePortSonetSdhGroup struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	GroupIndex         *string `json:"group-index,omitempty"`
	Payload            *string `json:"payload,omitempty"`
}

// ConfigurePortSonetSdhHoldTime struct
type ConfigurePortSonetSdhHoldTime struct {
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=100
	Down *uint32 `json:"down,omitempty"`
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=100
	// +kubebuilder:default:=5
	Up *uint32 `json:"up,omitempty"`
}

// ConfigurePortSonetSdhPath struct
type ConfigurePortSonetSdhPath struct {
	AdminState         *string `json:"admin-state,omitempty"`
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// kubebuilder:validation:Minimum=16
	// kubebuilder:validation:Maximum=16
	Crc                    *uint32                          `json:"crc,omitempty"`
	Description            *string                          `json:"description,omitempty"`
	Egress                 *ConfigurePortSonetSdhPathEgress `json:"egress,omitempty"`
	EncapType              *string                          `json:"encap-type,omitempty"`
	LoadBalancingAlgorithm *string                          `json:"load-balancing-algorithm,omitempty"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`[0-9a-fA-F]{2}(:[0-9a-fA-F]{2}){5}`
	// +kubebuilder:default:="00:00:00:00:00:00"
	MacAddress *string `json:"mac-address,omitempty"`
	// +kubebuilder:validation:Enum=`access`;`network`
	Mode *string `json:"mode,omitempty"`
	// kubebuilder:validation:Minimum=512
	// kubebuilder:validation:Maximum=9208
	Mtu         *uint32                               `json:"mtu,omitempty"`
	Network     *ConfigurePortSonetSdhPathNetwork     `json:"network,omitempty"`
	PathIndex   *string                               `json:"path-index,omitempty"`
	Payload     *string                               `json:"payload,omitempty"`
	Ppp         *ConfigurePortSonetSdhPathPpp         `json:"ppp,omitempty"`
	ReportAlarm *ConfigurePortSonetSdhPathReportAlarm `json:"report-alarm,omitempty"`
	Scramble    *bool                                 `json:"scramble,omitempty"`
	// kubebuilder:validation:MinLength=1
	// kubebuilder:validation:MaxLength=4
	SignalLabel *string `json:"signal-label,omitempty"`
	TraceString *string `json:"trace-string,omitempty"`
}

// ConfigurePortSonetSdhPathEgress struct
type ConfigurePortSonetSdhPathEgress struct {
	PortSchedulerPolicy *ConfigurePortSonetSdhPathEgressPortSchedulerPolicy `json:"port-scheduler-policy,omitempty"`
}

// ConfigurePortSonetSdhPathEgressPortSchedulerPolicy struct
type ConfigurePortSonetSdhPathEgressPortSchedulerPolicy struct {
	Overrides  *ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverrides `json:"overrides,omitempty"`
	PolicyName *string                                                      `json:"policy-name,omitempty"`
}

// ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverrides struct
type ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverrides struct {
	ApplyGroups        *string                                                             `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                                                             `json:"apply-groups-exclude,omitempty"`
	Level              []*ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesLevel `json:"level,omitempty"`
	MaxRate            *ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesMaxRate `json:"max-rate,omitempty"`
}

// ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesLevel struct
type ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesLevel struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=8
	PriorityLevel     *uint32                                                                            `json:"priority-level,omitempty"`
	RateOrPercentRate *ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesLevelRateOrPercentRate `json:"rate-or-percent-rate,omitempty"`
}

// ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesLevelRateOrPercentRate struct
type ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesLevelRateOrPercentRate struct {
	PercentRate *ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRate `json:"percent-rate,omitempty"`
	Rate        *ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRate        `json:"rate,omitempty"`
}

// ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRate struct
type ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRate struct {
	PercentRate *ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRatePercentRate `json:"percent-rate,omitempty"`
}

// ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRatePercentRate struct
type ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRatePercentRate struct {
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=10000
	// +kubebuilder:default:="100"
	Cir *string `json:"cir,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=10000
	// +kubebuilder:default:="100"
	Pir *string `json:"pir,omitempty"`
}

// ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRate struct
type ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRate struct {
	Rate *ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRateRate `json:"rate,omitempty"`
}

// ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRateRate struct
type ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRateRate struct {
	Cir *string `json:"cir,omitempty"`
	Pir *string `json:"pir,omitempty"`
}

// ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesMaxRate struct
type ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesMaxRate struct {
	RateOrPercentRate *ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRate `json:"rate-or-percent-rate,omitempty"`
}

// ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRate struct
type ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRate struct {
	PercentRate *ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRatePercentRate `json:"percent-rate,omitempty"`
	Rate        *ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRateRate        `json:"rate,omitempty"`
}

// ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRatePercentRate struct
type ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRatePercentRate struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=10000
	// +kubebuilder:default:="100"
	PercentRate *string `json:"percent-rate,omitempty"`
}

// ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRateRate struct
type ConfigurePortSonetSdhPathEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRateRate struct {
	Rate *string `json:"rate,omitempty"`
}

// ConfigurePortSonetSdhPathNetwork struct
type ConfigurePortSonetSdhPathNetwork struct {
	AccountingPolicy   *string `json:"accounting-policy,omitempty"`
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:default:=false
	CollectStats *bool   `json:"collect-stats,omitempty"`
	QueuePolicy  *string `json:"queue-policy,omitempty"`
}

// ConfigurePortSonetSdhPathPpp struct
type ConfigurePortSonetSdhPathPpp struct {
	ApplyGroups        *string                                `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                                `json:"apply-groups-exclude,omitempty"`
	Keepalive          *ConfigurePortSonetSdhPathPppKeepalive `json:"keepalive,omitempty"`
}

// ConfigurePortSonetSdhPathPppKeepalive struct
type ConfigurePortSonetSdhPathPppKeepalive struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=255
	// +kubebuilder:default:=3
	DropCount *uint32 `json:"drop-count,omitempty"`
	// +kubebuilder:default:="10"
	Interval *string `json:"interval,omitempty"`
}

// ConfigurePortSonetSdhPathReportAlarm struct
type ConfigurePortSonetSdhPathReportAlarm struct {
	// +kubebuilder:default:=false
	Pais *bool `json:"pais,omitempty"`
	// +kubebuilder:default:=false
	Plcd *bool `json:"plcd,omitempty"`
	// +kubebuilder:default:=true
	Plop *bool `json:"plop,omitempty"`
	// +kubebuilder:default:=true
	Pplm *bool `json:"pplm,omitempty"`
	// +kubebuilder:default:=false
	Prdi *bool `json:"prdi,omitempty"`
	// +kubebuilder:default:=false
	Prei *bool `json:"prei,omitempty"`
	// +kubebuilder:default:=true
	Puneq *bool `json:"puneq,omitempty"`
}

// ConfigurePortSonetSdhReportAlarm struct
type ConfigurePortSonetSdhReportAlarm struct {
	// +kubebuilder:default:=false
	Lais *bool `json:"lais,omitempty"`
	// +kubebuilder:default:=false
	Lb2erSd *bool `json:"lb2er-sd,omitempty"`
	// +kubebuilder:default:=true
	Lb2erSf *bool `json:"lb2er-sf,omitempty"`
	// +kubebuilder:default:=true
	Loc *bool `json:"loc,omitempty"`
	// +kubebuilder:default:=true
	Lrdi *bool `json:"lrdi,omitempty"`
	// +kubebuilder:default:=false
	Lrei *bool `json:"lrei,omitempty"`
	// +kubebuilder:default:=true
	Slof *bool `json:"slof,omitempty"`
	// +kubebuilder:default:=true
	Slos *bool `json:"slos,omitempty"`
	// +kubebuilder:default:=false
	Ss1f *bool `json:"ss1f,omitempty"`
}

// ConfigurePortSonetSdhSectionTrace struct
type ConfigurePortSonetSdhSectionTrace struct {
	SectionTrace *ConfigurePortSonetSdhSectionTraceSectionTrace `json:"section-trace,omitempty"`
}

// ConfigurePortSonetSdhSectionTraceSectionTrace struct
type ConfigurePortSonetSdhSectionTraceSectionTrace struct {
	Byte        *ConfigurePortSonetSdhSectionTraceSectionTraceByte        `json:"byte,omitempty"`
	IncrementZ0 *ConfigurePortSonetSdhSectionTraceSectionTraceIncrementZ0 `json:"increment-z0,omitempty"`
	String      *ConfigurePortSonetSdhSectionTraceSectionTraceString      `json:"string,omitempty"`
}

// ConfigurePortSonetSdhSectionTraceSectionTraceByte struct
type ConfigurePortSonetSdhSectionTraceSectionTraceByte struct {
	// kubebuilder:validation:MinLength=1
	// kubebuilder:validation:MaxLength=4
	// +kubebuilder:default:="1"
	Byte *string `json:"byte,omitempty"`
}

// ConfigurePortSonetSdhSectionTraceSectionTraceIncrementZ0 struct
type ConfigurePortSonetSdhSectionTraceSectionTraceIncrementZ0 struct {
	IncrementZ0 *string `json:"increment-z0,omitempty"`
}

// ConfigurePortSonetSdhSectionTraceSectionTraceString struct
type ConfigurePortSonetSdhSectionTraceSectionTraceString struct {
	// kubebuilder:validation:MinLength=0
	// kubebuilder:validation:MaxLength=16
	String *string `json:"string,omitempty"`
}

// ConfigurePortTdm struct
type ConfigurePortTdm struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:validation:Enum=`long`;`short`
	// +kubebuilder:default:="short"
	Buildout *string                `json:"buildout,omitempty"`
	Ds1      []*ConfigurePortTdmDs1 `json:"ds1,omitempty"`
	Ds3      []*ConfigurePortTdmDs3 `json:"ds3,omitempty"`
	E1       []*ConfigurePortTdmE1  `json:"e1,omitempty"`
	E3       []*ConfigurePortTdmE3  `json:"e3,omitempty"`
	// +kubebuilder:validation:Enum=`ami`;`b8zs`;`hdb3`
	Encoding *string                   `json:"encoding,omitempty"`
	HoldTime *ConfigurePortTdmHoldTime `json:"hold-time,omitempty"`
	// kubebuilder:validation:Minimum=75
	// kubebuilder:validation:Maximum=75
	LineImpedance *uint32 `json:"line-impedance,omitempty"`
}

// ConfigurePortTdmDs1 struct
type ConfigurePortTdmDs1 struct {
	AdminState         *string                            `json:"admin-state,omitempty"`
	ApplyGroups        *string                            `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                            `json:"apply-groups-exclude,omitempty"`
	BerThreshold       *ConfigurePortTdmDs1BerThreshold   `json:"ber-threshold,omitempty"`
	ChannelGroup       []*ConfigurePortTdmDs1ChannelGroup `json:"channel-group,omitempty"`
	// +kubebuilder:validation:Enum=`adaptive`;`differential`;`loop-timed`;`node-timed`
	ClockSource *string `json:"clock-source,omitempty"`
	Ds1Index    *string `json:"ds1-index,omitempty"`
	// +kubebuilder:validation:Enum=`ds1-unframed`;`extended-super-frame`;`super-frame`
	// +kubebuilder:default:="extended-super-frame"
	Framing  *string                      `json:"framing,omitempty"`
	HoldTime *ConfigurePortTdmDs1HoldTime `json:"hold-time,omitempty"`
	// +kubebuilder:validation:Enum=`fdl-ansi`;`fdl-bellcore`;`inband-ansi`;`inband-bellcore`;`internal`;`line`;`payload-ansi`
	Loopback *string `json:"loopback,omitempty"`
	// +kubebuilder:default:=false
	RemoteLoopRespond *bool                           `json:"remote-loop-respond,omitempty"`
	ReportAlarm       *ConfigurePortTdmDs1ReportAlarm `json:"report-alarm,omitempty"`
	// +kubebuilder:validation:Enum=`channel-associated-signaling`
	SignalMode *string `json:"signal-mode,omitempty"`
}

// ConfigurePortTdmDs1BerThreshold struct
type ConfigurePortTdmDs1BerThreshold struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=1
	// +kubebuilder:default:=5
	SignalDegrade *uint32 `json:"signal-degrade,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=1
	// +kubebuilder:default:=50
	SignalFailure *uint32 `json:"signal-failure,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroup struct
type ConfigurePortTdmDs1ChannelGroup struct {
	AdminState         *string `json:"admin-state,omitempty"`
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// kubebuilder:validation:Minimum=16
	// kubebuilder:validation:Maximum=16
	Crc         *uint32 `json:"crc,omitempty"`
	Description *string `json:"description,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=24
	Ds0Index               *uint32                                         `json:"ds0-index,omitempty"`
	Egress                 *ConfigurePortTdmDs1ChannelGroupEgress          `json:"egress,omitempty"`
	EncapType              *string                                         `json:"encap-type,omitempty"`
	IdleCycleFlag          *string                                         `json:"idle-cycle-flag,omitempty"`
	IdlePayloadFill        *ConfigurePortTdmDs1ChannelGroupIdlePayloadFill `json:"idle-payload-fill,omitempty"`
	IdleSignalFill         *ConfigurePortTdmDs1ChannelGroupIdleSignalFill  `json:"idle-signal-fill,omitempty"`
	LoadBalancingAlgorithm *string                                         `json:"load-balancing-algorithm,omitempty"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`[0-9a-fA-F]{2}(:[0-9a-fA-F]{2}){5}`
	// +kubebuilder:default:="00:00:00:00:00:00"
	MacAddress *string `json:"mac-address,omitempty"`
	// +kubebuilder:validation:Enum=`access`;`network`
	Mode *string `json:"mode,omitempty"`
	// kubebuilder:validation:Minimum=512
	// kubebuilder:validation:Maximum=9208
	Mtu     *uint32                                 `json:"mtu,omitempty"`
	Network *ConfigurePortTdmDs1ChannelGroupNetwork `json:"network,omitempty"`
	Ppp     *ConfigurePortTdmDs1ChannelGroupPpp     `json:"ppp,omitempty"`
	// kubebuilder:validation:Minimum=56
	// kubebuilder:validation:Maximum=56
	// +kubebuilder:default:=64
	Speed *uint32 `json:"speed,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=24
	Timeslot *uint32 `json:"timeslot,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupEgress struct
type ConfigurePortTdmDs1ChannelGroupEgress struct {
	PortSchedulerPolicy *ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicy `json:"port-scheduler-policy,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicy struct
type ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicy struct {
	Overrides  *ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverrides `json:"overrides,omitempty"`
	PolicyName *string                                                            `json:"policy-name,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverrides struct
type ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverrides struct {
	ApplyGroups        *string                                                                   `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                                                                   `json:"apply-groups-exclude,omitempty"`
	Level              []*ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesLevel `json:"level,omitempty"`
	MaxRate            *ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRate `json:"max-rate,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesLevel struct
type ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesLevel struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=8
	PriorityLevel     *uint32                                                                                  `json:"priority-level,omitempty"`
	RateOrPercentRate *ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRate `json:"rate-or-percent-rate,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRate struct
type ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRate struct {
	PercentRate *ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRate `json:"percent-rate,omitempty"`
	Rate        *ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRate        `json:"rate,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRate struct
type ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRate struct {
	PercentRate *ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRatePercentRate `json:"percent-rate,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRatePercentRate struct
type ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRatePercentRate struct {
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=10000
	// +kubebuilder:default:="100"
	Cir *string `json:"cir,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=10000
	// +kubebuilder:default:="100"
	Pir *string `json:"pir,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRate struct
type ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRate struct {
	Rate *ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRateRate `json:"rate,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRateRate struct
type ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRateRate struct {
	Cir *string `json:"cir,omitempty"`
	Pir *string `json:"pir,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRate struct
type ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRate struct {
	RateOrPercentRate *ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRate `json:"rate-or-percent-rate,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRate struct
type ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRate struct {
	PercentRate *ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRatePercentRate `json:"percent-rate,omitempty"`
	Rate        *ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRateRate        `json:"rate,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRatePercentRate struct
type ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRatePercentRate struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=10000
	// +kubebuilder:default:="100"
	PercentRate *string `json:"percent-rate,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRateRate struct
type ConfigurePortTdmDs1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRateRate struct {
	Rate *string `json:"rate,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupIdlePayloadFill struct
type ConfigurePortTdmDs1ChannelGroupIdlePayloadFill struct {
	IdlePayloadFillChoice *ConfigurePortTdmDs1ChannelGroupIdlePayloadFillIdlePayloadFillChoice `json:"idle-payload-fill-choice,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupIdlePayloadFillIdlePayloadFillChoice struct
type ConfigurePortTdmDs1ChannelGroupIdlePayloadFillIdlePayloadFillChoice struct {
	AllOnes *ConfigurePortTdmDs1ChannelGroupIdlePayloadFillIdlePayloadFillChoiceAllOnes `json:"all-ones,omitempty"`
	Pattern *ConfigurePortTdmDs1ChannelGroupIdlePayloadFillIdlePayloadFillChoicePattern `json:"pattern,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupIdlePayloadFillIdlePayloadFillChoiceAllOnes struct
type ConfigurePortTdmDs1ChannelGroupIdlePayloadFillIdlePayloadFillChoiceAllOnes struct {
	AllOnes *string `json:"all-ones,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupIdlePayloadFillIdlePayloadFillChoicePattern struct
type ConfigurePortTdmDs1ChannelGroupIdlePayloadFillIdlePayloadFillChoicePattern struct {
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=255
	Pattern *uint32 `json:"pattern,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupIdleSignalFill struct
type ConfigurePortTdmDs1ChannelGroupIdleSignalFill struct {
	IdleSignalFillChoice *ConfigurePortTdmDs1ChannelGroupIdleSignalFillIdleSignalFillChoice `json:"idle-signal-fill-choice,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupIdleSignalFillIdleSignalFillChoice struct
type ConfigurePortTdmDs1ChannelGroupIdleSignalFillIdleSignalFillChoice struct {
	AllOnes *ConfigurePortTdmDs1ChannelGroupIdleSignalFillIdleSignalFillChoiceAllOnes `json:"all-ones,omitempty"`
	Pattern *ConfigurePortTdmDs1ChannelGroupIdleSignalFillIdleSignalFillChoicePattern `json:"pattern,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupIdleSignalFillIdleSignalFillChoiceAllOnes struct
type ConfigurePortTdmDs1ChannelGroupIdleSignalFillIdleSignalFillChoiceAllOnes struct {
	AllOnes *string `json:"all-ones,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupIdleSignalFillIdleSignalFillChoicePattern struct
type ConfigurePortTdmDs1ChannelGroupIdleSignalFillIdleSignalFillChoicePattern struct {
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=15
	Pattern *uint32 `json:"pattern,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupNetwork struct
type ConfigurePortTdmDs1ChannelGroupNetwork struct {
	AccountingPolicy   *string `json:"accounting-policy,omitempty"`
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:default:=false
	CollectStats *bool   `json:"collect-stats,omitempty"`
	QueuePolicy  *string `json:"queue-policy,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupPpp struct
type ConfigurePortTdmDs1ChannelGroupPpp struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:default:=false
	BerSfLinkDown *bool                                        `json:"ber-sf-link-down,omitempty"`
	Compress      *ConfigurePortTdmDs1ChannelGroupPppCompress  `json:"compress,omitempty"`
	Keepalive     *ConfigurePortTdmDs1ChannelGroupPppKeepalive `json:"keepalive,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupPppCompress struct
type ConfigurePortTdmDs1ChannelGroupPppCompress struct {
	// +kubebuilder:default:=false
	Acfc *bool `json:"acfc,omitempty"`
	// +kubebuilder:default:=false
	Pfc *bool `json:"pfc,omitempty"`
}

// ConfigurePortTdmDs1ChannelGroupPppKeepalive struct
type ConfigurePortTdmDs1ChannelGroupPppKeepalive struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=255
	// +kubebuilder:default:=3
	DropCount *uint32 `json:"drop-count,omitempty"`
	// +kubebuilder:default:="10"
	Interval *string `json:"interval,omitempty"`
}

// ConfigurePortTdmDs1HoldTime struct
type ConfigurePortTdmDs1HoldTime struct {
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=100
	// +kubebuilder:default:=0
	Down *uint32 `json:"down,omitempty"`
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=100
	// +kubebuilder:default:=0
	Up *uint32 `json:"up,omitempty"`
}

// ConfigurePortTdmDs1ReportAlarm struct
type ConfigurePortTdmDs1ReportAlarm struct {
	// +kubebuilder:default:=true
	Ais *bool `json:"ais,omitempty"`
	// +kubebuilder:default:=false
	BerSd *bool `json:"ber-sd,omitempty"`
	// +kubebuilder:default:=false
	BerSf *bool `json:"ber-sf,omitempty"`
	// +kubebuilder:default:=false
	Looped *bool `json:"looped,omitempty"`
	// +kubebuilder:default:=true
	Los *bool `json:"los,omitempty"`
	// +kubebuilder:default:=false
	Oof *bool `json:"oof,omitempty"`
	// +kubebuilder:default:=false
	Rai *bool `json:"rai,omitempty"`
}

// ConfigurePortTdmDs3 struct
type ConfigurePortTdmDs3 struct {
	AdminState         *string `json:"admin-state,omitempty"`
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:validation:Enum=`ds1`;`e1`
	Channelized *string `json:"channelized,omitempty"`
	// +kubebuilder:validation:Enum=`loop-timed`;`node-timed`
	// +kubebuilder:default:="node-timed"
	ClockSource *string `json:"clock-source,omitempty"`
	// kubebuilder:validation:Minimum=16
	// kubebuilder:validation:Maximum=16
	Crc         *uint32                    `json:"crc,omitempty"`
	Description *string                    `json:"description,omitempty"`
	Ds3Index    *string                    `json:"ds3-index,omitempty"`
	Egress      *ConfigurePortTdmDs3Egress `json:"egress,omitempty"`
	EncapType   *string                    `json:"encap-type,omitempty"`
	// +kubebuilder:default:=false
	FeacLoopRespond *bool `json:"feac-loop-respond,omitempty"`
	// +kubebuilder:validation:Enum=`c-bit`;`ds3-unframed`;`m23`
	// +kubebuilder:default:="c-bit"
	Framing                *string `json:"framing,omitempty"`
	IdleCycleFlag          *string `json:"idle-cycle-flag,omitempty"`
	LoadBalancingAlgorithm *string `json:"load-balancing-algorithm,omitempty"`
	// +kubebuilder:validation:Enum=`internal`;`line`;`remote`
	Loopback *string `json:"loopback,omitempty"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`[0-9a-fA-F]{2}(:[0-9a-fA-F]{2}){5}`
	// +kubebuilder:default:="00:00:00:00:00:00"
	MacAddress          *string                                 `json:"mac-address,omitempty"`
	MaintenanceDataLink *ConfigurePortTdmDs3MaintenanceDataLink `json:"maintenance-data-link,omitempty"`
	// +kubebuilder:validation:Enum=`access`;`network`
	Mode *string `json:"mode,omitempty"`
	// kubebuilder:validation:Minimum=512
	// kubebuilder:validation:Maximum=9208
	Mtu         *uint32                         `json:"mtu,omitempty"`
	Network     *ConfigurePortTdmDs3Network     `json:"network,omitempty"`
	Ppp         *ConfigurePortTdmDs3Ppp         `json:"ppp,omitempty"`
	ReportAlarm *ConfigurePortTdmDs3ReportAlarm `json:"report-alarm,omitempty"`
	Scramble    *bool                           `json:"scramble,omitempty"`
	Subrate     *ConfigurePortTdmDs3Subrate     `json:"subrate,omitempty"`
}

// ConfigurePortTdmDs3Egress struct
type ConfigurePortTdmDs3Egress struct {
	PortSchedulerPolicy *ConfigurePortTdmDs3EgressPortSchedulerPolicy `json:"port-scheduler-policy,omitempty"`
}

// ConfigurePortTdmDs3EgressPortSchedulerPolicy struct
type ConfigurePortTdmDs3EgressPortSchedulerPolicy struct {
	Overrides  *ConfigurePortTdmDs3EgressPortSchedulerPolicyOverrides `json:"overrides,omitempty"`
	PolicyName *string                                                `json:"policy-name,omitempty"`
}

// ConfigurePortTdmDs3EgressPortSchedulerPolicyOverrides struct
type ConfigurePortTdmDs3EgressPortSchedulerPolicyOverrides struct {
	ApplyGroups        *string                                                       `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                                                       `json:"apply-groups-exclude,omitempty"`
	Level              []*ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesLevel `json:"level,omitempty"`
	MaxRate            *ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesMaxRate `json:"max-rate,omitempty"`
}

// ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesLevel struct
type ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesLevel struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=8
	PriorityLevel     *uint32                                                                      `json:"priority-level,omitempty"`
	RateOrPercentRate *ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRate `json:"rate-or-percent-rate,omitempty"`
}

// ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRate struct
type ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRate struct {
	PercentRate *ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRate `json:"percent-rate,omitempty"`
	Rate        *ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRate        `json:"rate,omitempty"`
}

// ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRate struct
type ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRate struct {
	PercentRate *ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRatePercentRate `json:"percent-rate,omitempty"`
}

// ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRatePercentRate struct
type ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRatePercentRate struct {
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=10000
	// +kubebuilder:default:="100"
	Cir *string `json:"cir,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=10000
	// +kubebuilder:default:="100"
	Pir *string `json:"pir,omitempty"`
}

// ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRate struct
type ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRate struct {
	Rate *ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRateRate `json:"rate,omitempty"`
}

// ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRateRate struct
type ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRateRate struct {
	Cir *string `json:"cir,omitempty"`
	Pir *string `json:"pir,omitempty"`
}

// ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesMaxRate struct
type ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesMaxRate struct {
	RateOrPercentRate *ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRate `json:"rate-or-percent-rate,omitempty"`
}

// ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRate struct
type ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRate struct {
	PercentRate *ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRatePercentRate `json:"percent-rate,omitempty"`
	Rate        *ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRateRate        `json:"rate,omitempty"`
}

// ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRatePercentRate struct
type ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRatePercentRate struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=10000
	// +kubebuilder:default:="100"
	PercentRate *string `json:"percent-rate,omitempty"`
}

// ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRateRate struct
type ConfigurePortTdmDs3EgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRateRate struct {
	Rate *string `json:"rate,omitempty"`
}

// ConfigurePortTdmDs3MaintenanceDataLink struct
type ConfigurePortTdmDs3MaintenanceDataLink struct {
	// kubebuilder:validation:MinLength=0
	// kubebuilder:validation:MaxLength=10
	EquipmentIdCode *string `json:"equipment-id-code,omitempty"`
	// kubebuilder:validation:MinLength=0
	// kubebuilder:validation:MaxLength=38
	FacilityIdCode *string `json:"facility-id-code,omitempty"`
	// kubebuilder:validation:MinLength=0
	// kubebuilder:validation:MaxLength=10
	FrameIdCode *string `json:"frame-id-code,omitempty"`
	// kubebuilder:validation:MinLength=0
	// kubebuilder:validation:MaxLength=38
	GeneratorString *string `json:"generator-string,omitempty"`
	// kubebuilder:validation:MinLength=0
	// kubebuilder:validation:MaxLength=11
	LocationIdCode *string `json:"location-id-code,omitempty"`
	// kubebuilder:validation:MinLength=0
	// kubebuilder:validation:MaxLength=38
	PortString          *string                                                    `json:"port-string,omitempty"`
	TransmitMessageType *ConfigurePortTdmDs3MaintenanceDataLinkTransmitMessageType `json:"transmit-message-type,omitempty"`
	// kubebuilder:validation:MinLength=0
	// kubebuilder:validation:MaxLength=6
	UnitIdCode *string `json:"unit-id-code,omitempty"`
}

// ConfigurePortTdmDs3MaintenanceDataLinkTransmitMessageType struct
type ConfigurePortTdmDs3MaintenanceDataLinkTransmitMessageType struct {
	// +kubebuilder:default:=false
	IdleSignal *bool `json:"idle-signal,omitempty"`
	// +kubebuilder:default:=false
	Path *bool `json:"path,omitempty"`
	// +kubebuilder:default:=false
	TestSignal *bool `json:"test-signal,omitempty"`
}

// ConfigurePortTdmDs3Network struct
type ConfigurePortTdmDs3Network struct {
	AccountingPolicy   *string `json:"accounting-policy,omitempty"`
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:default:=false
	CollectStats *bool   `json:"collect-stats,omitempty"`
	QueuePolicy  *string `json:"queue-policy,omitempty"`
}

// ConfigurePortTdmDs3Ppp struct
type ConfigurePortTdmDs3Ppp struct {
	ApplyGroups        *string                          `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                          `json:"apply-groups-exclude,omitempty"`
	Keepalive          *ConfigurePortTdmDs3PppKeepalive `json:"keepalive,omitempty"`
}

// ConfigurePortTdmDs3PppKeepalive struct
type ConfigurePortTdmDs3PppKeepalive struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=255
	// +kubebuilder:default:=3
	DropCount *uint32 `json:"drop-count,omitempty"`
	// +kubebuilder:default:="10"
	Interval *string `json:"interval,omitempty"`
}

// ConfigurePortTdmDs3ReportAlarm struct
type ConfigurePortTdmDs3ReportAlarm struct {
	// +kubebuilder:default:=true
	Ais *bool `json:"ais,omitempty"`
	// +kubebuilder:default:=false
	Looped *bool `json:"looped,omitempty"`
	// +kubebuilder:default:=true
	Los *bool `json:"los,omitempty"`
	// +kubebuilder:default:=false
	Oof *bool `json:"oof,omitempty"`
	// +kubebuilder:default:=false
	Rai *bool `json:"rai,omitempty"`
}

// ConfigurePortTdmDs3Subrate struct
type ConfigurePortTdmDs3Subrate struct {
	// +kubebuilder:validation:Enum=`digital-link`;`larscom`
	CsuMode *string `json:"csu-mode,omitempty"`
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=147
	RateStep *uint32 `json:"rate-step,omitempty"`
}

// ConfigurePortTdmE1 struct
type ConfigurePortTdmE1 struct {
	AdminState         *string                           `json:"admin-state,omitempty"`
	ApplyGroups        *string                           `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                           `json:"apply-groups-exclude,omitempty"`
	BerThreshold       *ConfigurePortTdmE1BerThreshold   `json:"ber-threshold,omitempty"`
	ChannelGroup       []*ConfigurePortTdmE1ChannelGroup `json:"channel-group,omitempty"`
	// +kubebuilder:validation:Enum=`adaptive`;`differential`;`loop-timed`;`node-timed`
	ClockSource *string `json:"clock-source,omitempty"`
	E1Index     *string `json:"e1-index,omitempty"`
	// +kubebuilder:validation:Enum=`e1-unframed`;`g704`;`no-crc-g704`
	// +kubebuilder:default:="g704"
	Framing  *string                     `json:"framing,omitempty"`
	HoldTime *ConfigurePortTdmE1HoldTime `json:"hold-time,omitempty"`
	// +kubebuilder:validation:Enum=`internal`;`line`
	Loopback     *string                         `json:"loopback,omitempty"`
	NationalBits *ConfigurePortTdmE1NationalBits `json:"national-bits,omitempty"`
	ReportAlarm  *ConfigurePortTdmE1ReportAlarm  `json:"report-alarm,omitempty"`
	// +kubebuilder:validation:Enum=`channel-associated-signaling`
	SignalMode *string `json:"signal-mode,omitempty"`
}

// ConfigurePortTdmE1BerThreshold struct
type ConfigurePortTdmE1BerThreshold struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=1
	// +kubebuilder:default:=5
	SignalDegrade *uint32 `json:"signal-degrade,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=1
	// +kubebuilder:default:=50
	SignalFailure *uint32 `json:"signal-failure,omitempty"`
}

// ConfigurePortTdmE1ChannelGroup struct
type ConfigurePortTdmE1ChannelGroup struct {
	AdminState         *string `json:"admin-state,omitempty"`
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// kubebuilder:validation:Minimum=16
	// kubebuilder:validation:Maximum=16
	Crc         *uint32 `json:"crc,omitempty"`
	Description *string `json:"description,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=32
	Ds0Index               *uint32                                        `json:"ds0-index,omitempty"`
	Egress                 *ConfigurePortTdmE1ChannelGroupEgress          `json:"egress,omitempty"`
	EncapType              *string                                        `json:"encap-type,omitempty"`
	IdleCycleFlag          *string                                        `json:"idle-cycle-flag,omitempty"`
	IdlePayloadFill        *ConfigurePortTdmE1ChannelGroupIdlePayloadFill `json:"idle-payload-fill,omitempty"`
	IdleSignalFill         *ConfigurePortTdmE1ChannelGroupIdleSignalFill  `json:"idle-signal-fill,omitempty"`
	LoadBalancingAlgorithm *string                                        `json:"load-balancing-algorithm,omitempty"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`[0-9a-fA-F]{2}(:[0-9a-fA-F]{2}){5}`
	// +kubebuilder:default:="00:00:00:00:00:00"
	MacAddress *string `json:"mac-address,omitempty"`
	// +kubebuilder:validation:Enum=`access`;`network`
	Mode *string `json:"mode,omitempty"`
	// kubebuilder:validation:Minimum=512
	// kubebuilder:validation:Maximum=9208
	Mtu     *uint32                                `json:"mtu,omitempty"`
	Network *ConfigurePortTdmE1ChannelGroupNetwork `json:"network,omitempty"`
	Ppp     *ConfigurePortTdmE1ChannelGroupPpp     `json:"ppp,omitempty"`
	// kubebuilder:validation:Minimum=56
	// kubebuilder:validation:Maximum=56
	// +kubebuilder:default:=64
	Speed *uint32 `json:"speed,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=32
	Timeslot *uint32 `json:"timeslot,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupEgress struct
type ConfigurePortTdmE1ChannelGroupEgress struct {
	PortSchedulerPolicy *ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicy `json:"port-scheduler-policy,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicy struct
type ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicy struct {
	Overrides  *ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverrides `json:"overrides,omitempty"`
	PolicyName *string                                                           `json:"policy-name,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverrides struct
type ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverrides struct {
	ApplyGroups        *string                                                                  `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                                                                  `json:"apply-groups-exclude,omitempty"`
	Level              []*ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesLevel `json:"level,omitempty"`
	MaxRate            *ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRate `json:"max-rate,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesLevel struct
type ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesLevel struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=8
	PriorityLevel     *uint32                                                                                 `json:"priority-level,omitempty"`
	RateOrPercentRate *ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRate `json:"rate-or-percent-rate,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRate struct
type ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRate struct {
	PercentRate *ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRate `json:"percent-rate,omitempty"`
	Rate        *ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRate        `json:"rate,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRate struct
type ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRate struct {
	PercentRate *ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRatePercentRate `json:"percent-rate,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRatePercentRate struct
type ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRatePercentRate struct {
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=10000
	// +kubebuilder:default:="100"
	Cir *string `json:"cir,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=10000
	// +kubebuilder:default:="100"
	Pir *string `json:"pir,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRate struct
type ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRate struct {
	Rate *ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRateRate `json:"rate,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRateRate struct
type ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRateRate struct {
	Cir *string `json:"cir,omitempty"`
	Pir *string `json:"pir,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRate struct
type ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRate struct {
	RateOrPercentRate *ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRate `json:"rate-or-percent-rate,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRate struct
type ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRate struct {
	PercentRate *ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRatePercentRate `json:"percent-rate,omitempty"`
	Rate        *ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRateRate        `json:"rate,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRatePercentRate struct
type ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRatePercentRate struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=10000
	// +kubebuilder:default:="100"
	PercentRate *string `json:"percent-rate,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRateRate struct
type ConfigurePortTdmE1ChannelGroupEgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRateRate struct {
	Rate *string `json:"rate,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupIdlePayloadFill struct
type ConfigurePortTdmE1ChannelGroupIdlePayloadFill struct {
	IdlePayloadFillChoice *ConfigurePortTdmE1ChannelGroupIdlePayloadFillIdlePayloadFillChoice `json:"idle-payload-fill-choice,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupIdlePayloadFillIdlePayloadFillChoice struct
type ConfigurePortTdmE1ChannelGroupIdlePayloadFillIdlePayloadFillChoice struct {
	AllOnes *ConfigurePortTdmE1ChannelGroupIdlePayloadFillIdlePayloadFillChoiceAllOnes `json:"all-ones,omitempty"`
	Pattern *ConfigurePortTdmE1ChannelGroupIdlePayloadFillIdlePayloadFillChoicePattern `json:"pattern,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupIdlePayloadFillIdlePayloadFillChoiceAllOnes struct
type ConfigurePortTdmE1ChannelGroupIdlePayloadFillIdlePayloadFillChoiceAllOnes struct {
	AllOnes *string `json:"all-ones,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupIdlePayloadFillIdlePayloadFillChoicePattern struct
type ConfigurePortTdmE1ChannelGroupIdlePayloadFillIdlePayloadFillChoicePattern struct {
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=255
	Pattern *uint32 `json:"pattern,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupIdleSignalFill struct
type ConfigurePortTdmE1ChannelGroupIdleSignalFill struct {
	IdleSignalFillChoice *ConfigurePortTdmE1ChannelGroupIdleSignalFillIdleSignalFillChoice `json:"idle-signal-fill-choice,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupIdleSignalFillIdleSignalFillChoice struct
type ConfigurePortTdmE1ChannelGroupIdleSignalFillIdleSignalFillChoice struct {
	AllOnes *ConfigurePortTdmE1ChannelGroupIdleSignalFillIdleSignalFillChoiceAllOnes `json:"all-ones,omitempty"`
	Pattern *ConfigurePortTdmE1ChannelGroupIdleSignalFillIdleSignalFillChoicePattern `json:"pattern,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupIdleSignalFillIdleSignalFillChoiceAllOnes struct
type ConfigurePortTdmE1ChannelGroupIdleSignalFillIdleSignalFillChoiceAllOnes struct {
	AllOnes *string `json:"all-ones,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupIdleSignalFillIdleSignalFillChoicePattern struct
type ConfigurePortTdmE1ChannelGroupIdleSignalFillIdleSignalFillChoicePattern struct {
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=15
	Pattern *uint32 `json:"pattern,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupNetwork struct
type ConfigurePortTdmE1ChannelGroupNetwork struct {
	AccountingPolicy   *string `json:"accounting-policy,omitempty"`
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:default:=false
	CollectStats *bool   `json:"collect-stats,omitempty"`
	QueuePolicy  *string `json:"queue-policy,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupPpp struct
type ConfigurePortTdmE1ChannelGroupPpp struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:default:=false
	BerSfLinkDown *bool                                       `json:"ber-sf-link-down,omitempty"`
	Compress      *ConfigurePortTdmE1ChannelGroupPppCompress  `json:"compress,omitempty"`
	Keepalive     *ConfigurePortTdmE1ChannelGroupPppKeepalive `json:"keepalive,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupPppCompress struct
type ConfigurePortTdmE1ChannelGroupPppCompress struct {
	// +kubebuilder:default:=false
	Acfc *bool `json:"acfc,omitempty"`
	// +kubebuilder:default:=false
	Pfc *bool `json:"pfc,omitempty"`
}

// ConfigurePortTdmE1ChannelGroupPppKeepalive struct
type ConfigurePortTdmE1ChannelGroupPppKeepalive struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=255
	// +kubebuilder:default:=3
	DropCount *uint32 `json:"drop-count,omitempty"`
	// +kubebuilder:default:="10"
	Interval *string `json:"interval,omitempty"`
}

// ConfigurePortTdmE1HoldTime struct
type ConfigurePortTdmE1HoldTime struct {
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=100
	// +kubebuilder:default:=0
	Down *uint32 `json:"down,omitempty"`
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=100
	// +kubebuilder:default:=0
	Up *uint32 `json:"up,omitempty"`
}

// ConfigurePortTdmE1NationalBits struct
type ConfigurePortTdmE1NationalBits struct {
	// +kubebuilder:default:=false
	Sa4 *bool `json:"sa4,omitempty"`
	// +kubebuilder:default:=false
	Sa5 *bool `json:"sa5,omitempty"`
	// +kubebuilder:default:=false
	Sa6 *bool `json:"sa6,omitempty"`
	// +kubebuilder:default:=false
	Sa7 *bool `json:"sa7,omitempty"`
	// +kubebuilder:default:=false
	Sa8 *bool `json:"sa8,omitempty"`
}

// ConfigurePortTdmE1ReportAlarm struct
type ConfigurePortTdmE1ReportAlarm struct {
	// +kubebuilder:default:=true
	Ais *bool `json:"ais,omitempty"`
	// +kubebuilder:default:=false
	BerSd *bool `json:"ber-sd,omitempty"`
	// +kubebuilder:default:=false
	BerSf *bool `json:"ber-sf,omitempty"`
	// +kubebuilder:default:=false
	Looped *bool `json:"looped,omitempty"`
	// +kubebuilder:default:=true
	Los *bool `json:"los,omitempty"`
	// +kubebuilder:default:=false
	Oof *bool `json:"oof,omitempty"`
	// +kubebuilder:default:=false
	Rai *bool `json:"rai,omitempty"`
}

// ConfigurePortTdmE3 struct
type ConfigurePortTdmE3 struct {
	AdminState         *string `json:"admin-state,omitempty"`
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:validation:Enum=`loop-timed`;`node-timed`
	// +kubebuilder:default:="node-timed"
	ClockSource *string `json:"clock-source,omitempty"`
	// kubebuilder:validation:Minimum=16
	// kubebuilder:validation:Maximum=16
	Crc         *uint32                   `json:"crc,omitempty"`
	Description *string                   `json:"description,omitempty"`
	E3Index     *string                   `json:"e3-index,omitempty"`
	Egress      *ConfigurePortTdmE3Egress `json:"egress,omitempty"`
	EncapType   *string                   `json:"encap-type,omitempty"`
	// +kubebuilder:validation:Enum=`e3-unframed`;`g751`;`g832`
	// +kubebuilder:default:="g751"
	Framing                *string `json:"framing,omitempty"`
	IdleCycleFlag          *string `json:"idle-cycle-flag,omitempty"`
	LoadBalancingAlgorithm *string `json:"load-balancing-algorithm,omitempty"`
	// +kubebuilder:validation:Enum=`internal`;`line`
	Loopback *string `json:"loopback,omitempty"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`[0-9a-fA-F]{2}(:[0-9a-fA-F]{2}){5}`
	// +kubebuilder:default:="00:00:00:00:00:00"
	MacAddress *string `json:"mac-address,omitempty"`
	// +kubebuilder:validation:Enum=`access`;`network`
	Mode *string `json:"mode,omitempty"`
	// kubebuilder:validation:Minimum=512
	// kubebuilder:validation:Maximum=9208
	Mtu         *uint32                        `json:"mtu,omitempty"`
	Network     *ConfigurePortTdmE3Network     `json:"network,omitempty"`
	Ppp         *ConfigurePortTdmE3Ppp         `json:"ppp,omitempty"`
	ReportAlarm *ConfigurePortTdmE3ReportAlarm `json:"report-alarm,omitempty"`
	Scramble    *bool                          `json:"scramble,omitempty"`
}

// ConfigurePortTdmE3Egress struct
type ConfigurePortTdmE3Egress struct {
	PortSchedulerPolicy *ConfigurePortTdmE3EgressPortSchedulerPolicy `json:"port-scheduler-policy,omitempty"`
}

// ConfigurePortTdmE3EgressPortSchedulerPolicy struct
type ConfigurePortTdmE3EgressPortSchedulerPolicy struct {
	Overrides  *ConfigurePortTdmE3EgressPortSchedulerPolicyOverrides `json:"overrides,omitempty"`
	PolicyName *string                                               `json:"policy-name,omitempty"`
}

// ConfigurePortTdmE3EgressPortSchedulerPolicyOverrides struct
type ConfigurePortTdmE3EgressPortSchedulerPolicyOverrides struct {
	ApplyGroups        *string                                                      `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                                                      `json:"apply-groups-exclude,omitempty"`
	Level              []*ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesLevel `json:"level,omitempty"`
	MaxRate            *ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesMaxRate `json:"max-rate,omitempty"`
}

// ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesLevel struct
type ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesLevel struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=8
	PriorityLevel     *uint32                                                                     `json:"priority-level,omitempty"`
	RateOrPercentRate *ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRate `json:"rate-or-percent-rate,omitempty"`
}

// ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRate struct
type ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRate struct {
	PercentRate *ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRate `json:"percent-rate,omitempty"`
	Rate        *ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRate        `json:"rate,omitempty"`
}

// ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRate struct
type ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRate struct {
	PercentRate *ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRatePercentRate `json:"percent-rate,omitempty"`
}

// ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRatePercentRate struct
type ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRatePercentRatePercentRate struct {
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=10000
	// +kubebuilder:default:="100"
	Cir *string `json:"cir,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=10000
	// +kubebuilder:default:="100"
	Pir *string `json:"pir,omitempty"`
}

// ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRate struct
type ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRate struct {
	Rate *ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRateRate `json:"rate,omitempty"`
}

// ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRateRate struct
type ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesLevelRateOrPercentRateRateRate struct {
	Cir *string `json:"cir,omitempty"`
	Pir *string `json:"pir,omitempty"`
}

// ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesMaxRate struct
type ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesMaxRate struct {
	RateOrPercentRate *ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRate `json:"rate-or-percent-rate,omitempty"`
}

// ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRate struct
type ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRate struct {
	PercentRate *ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRatePercentRate `json:"percent-rate,omitempty"`
	Rate        *ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRateRate        `json:"rate,omitempty"`
}

// ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRatePercentRate struct
type ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRatePercentRate struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=10000
	// +kubebuilder:default:="100"
	PercentRate *string `json:"percent-rate,omitempty"`
}

// ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRateRate struct
type ConfigurePortTdmE3EgressPortSchedulerPolicyOverridesMaxRateRateOrPercentRateRate struct {
	Rate *string `json:"rate,omitempty"`
}

// ConfigurePortTdmE3Network struct
type ConfigurePortTdmE3Network struct {
	AccountingPolicy   *string `json:"accounting-policy,omitempty"`
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:default:=false
	CollectStats *bool   `json:"collect-stats,omitempty"`
	QueuePolicy  *string `json:"queue-policy,omitempty"`
}

// ConfigurePortTdmE3Ppp struct
type ConfigurePortTdmE3Ppp struct {
	ApplyGroups        *string                         `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string                         `json:"apply-groups-exclude,omitempty"`
	Keepalive          *ConfigurePortTdmE3PppKeepalive `json:"keepalive,omitempty"`
}

// ConfigurePortTdmE3PppKeepalive struct
type ConfigurePortTdmE3PppKeepalive struct {
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=255
	// +kubebuilder:default:=3
	DropCount *uint32 `json:"drop-count,omitempty"`
	// +kubebuilder:default:="10"
	Interval *string `json:"interval,omitempty"`
}

// ConfigurePortTdmE3ReportAlarm struct
type ConfigurePortTdmE3ReportAlarm struct {
	// +kubebuilder:default:=true
	Ais *bool `json:"ais,omitempty"`
	// +kubebuilder:default:=false
	Looped *bool `json:"looped,omitempty"`
	// +kubebuilder:default:=true
	Los *bool `json:"los,omitempty"`
	// +kubebuilder:default:=false
	Oof *bool `json:"oof,omitempty"`
	// +kubebuilder:default:=false
	Rai *bool `json:"rai,omitempty"`
}

// ConfigurePortTdmHoldTime struct
type ConfigurePortTdmHoldTime struct {
	// kubebuilder:validation:Minimum=0
	// kubebuilder:validation:Maximum=100
	// +kubebuilder:default:=5
	Down *uint32 `json:"down,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=100
	Up *uint32 `json:"up,omitempty"`
}

// ConfigurePortTransceiver struct
type ConfigurePortTransceiver struct {
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	// +kubebuilder:default:=false
	DigitalCoherentOptics *bool `json:"digital-coherent-optics,omitempty"`
}

// ConfigurePortPolicy struct
type ConfigurePortPolicy struct {
	ApplyGroups               *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude        *string `json:"apply-groups-exclude,omitempty"`
	Description               *string `json:"description,omitempty"`
	EgressPortSchedulerPolicy *string `json:"egress-port-scheduler-policy,omitempty"`
	Name                      *string `json:"name,omitempty"`
}

// ConfigurePortXc struct
type ConfigurePortXc struct {
	ApplyGroups        *string               `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string               `json:"apply-groups-exclude,omitempty"`
	Pxc                []*ConfigurePortXcPxc `json:"pxc,omitempty"`
}

// ConfigurePortXcPxc struct
type ConfigurePortXcPxc struct {
	AdminState         *string `json:"admin-state,omitempty"`
	ApplyGroups        *string `json:"apply-groups,omitempty"`
	ApplyGroupsExclude *string `json:"apply-groups-exclude,omitempty"`
	Description        *string `json:"description,omitempty"`
	PortId             *string `json:"port-id,omitempty"`
	// kubebuilder:validation:Minimum=1
	// kubebuilder:validation:Maximum=64
	PxcId *uint32 `json:"pxc-id,omitempty"`
}

// ConfigurePortParameters are the parameter fields of a ConfigurePort.
type ConfigurePortParameters struct {
	SrosConfigurePort *ConfigurePort `json:"port,omitempty"`
}

// ConfigurePortObservation are the observable fields of a ConfigurePort.
type ConfigurePortObservation struct {
}

// A ConfigurePortSpec defines the desired state of a ConfigurePort.
type ConfigurePortSpec struct {
	nddv1.ResourceSpec `json:",inline"`
	ForNetworkNode     ConfigurePortParameters `json:"forNetworkNode"`
}

// A ConfigurePortStatus represents the observed state of a ConfigurePort.
type ConfigurePortStatus struct {
	nddv1.ResourceStatus `json:",inline"`
	AtNetworkNode        ConfigurePortObservation `json:"atNetworkNode,omitempty"`
}

// +kubebuilder:object:root=true

// SrosConfigurePort is the Schema for the ConfigurePort API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="TARGET",type="string",JSONPath=".status.conditions[?(@.kind=='TargetFound')].status"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.conditions[?(@.kind=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNC",type="string",JSONPath=".status.conditions[?(@.kind=='Synced')].status"
// +kubebuilder:printcolumn:name="LOCALLEAFREF",type="string",JSONPath=".status.conditions[?(@.kind=='InternalLeafrefValidationSuccess')].status"
// +kubebuilder:printcolumn:name="EXTLEAFREF",type="string",JSONPath=".status.conditions[?(@.kind=='ExternalLeafrefValidationSuccess')].status"
// +kubebuilder:printcolumn:name="PARENTDEP",type="string",JSONPath=".status.conditions[?(@.kind=='ParentValidationSuccess')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories={ndd,srl}
type SrosConfigurePort struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConfigurePortSpec   `json:"spec,omitempty"`
	Status ConfigurePortStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SrosConfigurePortList contains a list of ConfigurePorts
type SrosConfigurePortList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SrosConfigurePort `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SrosConfigurePort{}, &SrosConfigurePortList{})
}

// ConfigurePort type metadata.
var (
	ConfigurePortKind             = reflect.TypeOf(SrosConfigurePort{}).Name()
	ConfigurePortGroupKind        = schema.GroupKind{Group: Group, Kind: ConfigurePortKind}.String()
	ConfigurePortKindAPIVersion   = ConfigurePortKind + "." + GroupVersion.String()
	ConfigurePortGroupVersionKind = GroupVersion.WithKind(ConfigurePortKind)
)

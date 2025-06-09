package models

import (
	"strings"
)

const (
	JOB_NAMESPACE_CPU = "SolarWinds.Orion.Core.Pollers"
	JOB_TYPE_CPU      = "SolarWinds.Orion.Core.Pollers.CorePollersJob, SolarWinds.Orion.Core.Pollers, Version=2025.2.0.0, Culture=neutral, PublicKeyToken=null"

	JOB_NAMESPACE_DISCOVERY = "orion"
	JOB_TYPE_DISCOVERY      = "SolarWinds.Orion.Discovery.Job.OrionDiscoveryJob, SolarWinds.Orion.Discovery.Job, Version=2025.2.0.0, Culture=neutral, PublicKeyToken=null"
)

const node_inventory_job_description = `
{
    "Pollers": [
		{
            "Assignment": {
                "NetObjectType": "N",
                "NetObjectID": $NetObjectId$,
                "PollerType": "N.Details.SNMP.Generic",
                "Enabled": true
            },
			"Settings": {
				"$type": "SolarWinds.Orion.Core.Pollers.Node.Settings.NodeDetailsPollerSettings, SolarWinds.Orion.Core.Pollers",
				"CollectMAC": false
			},
            "SettingID": 1
        }
    ],
    "Settings": {
        1: {
            "$type": "SolarWinds.Orion.Pollers.Framework.SNMP.SnmpSettings, SolarWinds.Orion.Pollers.Framework",
            "IP": "$IP$",
            "AgentPort": 161,
            "ProtocolVersion": "None",
            "InterQueryDelayMs": 0
        }
    },
    "GlobalSettings": {
        "$type": "SolarWinds.Orion.Pollers.Framework.SNMP.SnmpGlobalSettings, SolarWinds.Orion.Pollers.Framework",
        "RequestTimeout": 2500,
        "RequestRetries": 2,
        "MaxReplies": 5,
        "HsrpEnabled": true,
        "VulnerabilityCheckDisabled": false
    }
}
`

const icmp_job_description = `
{
    "Pollers": [
		{
            "Assignment": {
                "NetObjectType": "N",
                "NetObjectID": $NetObjectId$,
                "PollerType": "N.StatusAndResponseTime.ICMP.SendEcho",
                "Enabled": true
            },
			"Settings": {
				"$type": "SolarWinds.Orion.Core.Pollers.Node.ResponseTime.Settings.NodeResponseTimePollerSettings, SolarWinds.Orion.Core.Pollers",
				"PollingSettings": {
				},
				"PolledData": 6
			},
            "SettingID": 1
        }
    ],
    "Settings": {
		"1": {
			"$type": "SolarWinds.Orion.Pollers.Framework.ICMP.IcmpSettings, SolarWinds.Orion.Pollers.Framework",
			"IP": "$IP$"
		}
    },
    "GlobalSettings": {
        "$type": "SolarWinds.Orion.Pollers.Framework.ICMP.IcmpGlobalSettings, SolarWinds.Orion.Pollers.Framework",
    }
}
`

const snmp_job_description = `
{
    "Pollers": [{
            "Assignment": {
                "NetObjectType": "N",
                "NetObjectID": $NetObjectId$,
                "PollerType": "$CpuPollerType$",
                "Enabled": true
            },
			"Settings": {
				"$type": "SolarWinds.Orion.Core.Pollers.Cpu.Settings.CpuPollerSetting, SolarWinds.Orion.Core.Pollers",
            },
            "SettingID": 1
        },
		{
            "Assignment": {
                "NetObjectType": "N",
                "NetObjectID": $NetObjectId$,
                "PollerType": "$MemoryPollerType$",
                "Enabled": true
            },
            "SettingID": 1
        }
    ],
    "Settings": {
        1: {
            "$type": "SolarWinds.Orion.Pollers.Framework.SNMP.SnmpSettings, SolarWinds.Orion.Pollers.Framework",
            "IP": "$IP$",
            "AgentPort": 161,
            "ProtocolVersion": "None",
            "InterQueryDelayMs": 0
        }
    },
    "GlobalSettings": {
        "$type": "SolarWinds.Orion.Pollers.Framework.SNMP.SnmpGlobalSettings, SolarWinds.Orion.Pollers.Framework",
        "RequestTimeout": 2500,
        "RequestRetries": 2,
        "MaxReplies": 5,
        "HsrpEnabled": true,
        "VulnerabilityCheckDisabled": false
    }
}
`

type vars = map[string]string

func applyVariables(input string, variables vars) string {
	for key, value := range variables {
		input = strings.ReplaceAll(input, "$"+key+"$", value)
	}
	return input
}

func GetCoreSnmpJobDescription(variables vars) string {
	return applyVariables(snmp_job_description, variables)
}

func GetCoreIcmpJobDescription(variables vars) string {
	return applyVariables(icmp_job_description, variables)
}

func GetCoreInventoryJobDescription(variables vars) string {
	return applyVariables(node_inventory_job_description, variables)
}

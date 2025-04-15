package models

const (
	JOB_NAMESPACE_PCU = "SolarWinds.PCU.Pollers"
	JOB_TYPE_PCU      = "SolarWinds.PCU.Pollers.PCUPollerJob, SolarWinds.PCU.Pollers, Version=2025.2.0.0, Culture=neutral, PublicKeyToken=null"
)

const pcu_job_description = `
{
    "Pollers": [{
            "Assignment": {
                "NetObjectType": "PCU",
                "NetObjectID": $NetObjectId$,
                "PollerType": "PCU.Statistics.SNMP.Generic",
                "Enabled": true
            },
            "Settings": {
                "$type": "SolarWinds.PCU.Pollers.PCUPollerSettings, SolarWinds.PCU.Pollers",
                "EntityId": 1
            },
            "SettingID": 1
        }
    ],
    "Settings": {
        "1": {
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

func GetPCUJobDescription(variables vars) string {
	return applyVariables(pcu_job_description, variables)
}

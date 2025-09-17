{
    "Pollers": [
		{
            "Assignment": {
                "NetObjectType": "I",
                "NetObjectID": {{ .NetObjectId }},
                "PollerType": "{{ .PollerType }}",
                "Enabled": true
            },
			"Settings": {
				"$type": "SolarWinds.Interfaces.Pollers.InterfacesPollerSetting, SolarWinds.Interfaces.Pollers",
				"InterfaceIndex": {{ .NetObjectId }},
                "Allow64BitCounters": true
			},
            "SettingID": 1
        }
    ],
    "Settings": {
        1: {
            "$type": "SolarWinds.Orion.Pollers.Framework.SNMP.SnmpSettings, SolarWinds.Orion.Pollers.Framework",
            "IP": "{{ .IP }}",
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
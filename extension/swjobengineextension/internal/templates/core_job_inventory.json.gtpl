{
    "Pollers": [
		{
            "Assignment": {
                "NetObjectType": "N",
                "NetObjectID": {{ .NetObjectId }},
                "PollerType": "{{ .PollerType }}",
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
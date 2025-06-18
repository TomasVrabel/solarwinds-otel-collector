{
    "Pollers": [
		{
            "Assignment": {
                "NetObjectType": "N",
                "NetObjectID": {{ .NetObjectId }},
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
			"IP": "{{ .IP }}"
		}
    },
    "GlobalSettings": {
        "$type": "SolarWinds.Orion.Pollers.Framework.ICMP.IcmpGlobalSettings, SolarWinds.Orion.Pollers.Framework",
    }
}
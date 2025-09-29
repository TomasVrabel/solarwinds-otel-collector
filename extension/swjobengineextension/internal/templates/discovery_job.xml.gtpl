<root>
	<data0>
		<DiscoveryPluginInfoCollection
			xmlns="http://schemas.solarwinds.com/2008/Orion"
			xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
			<PluginInfos>
				<DiscoveryPluginInfo>
					<ModuleName i:nil="true"/>
					<PluginAssemblyName>SolarWinds.Orion.Core.DiscoveryJobPlugin</PluginAssemblyName>
					<PluginTypeName>SolarWinds.Orion.Core.DiscoveryJobPlugin.CoreDiscoveryJobPlugin</PluginTypeName>
					<ProcessingOrder>1</ProcessingOrder>
					<RootPath i:nil="true"/>
					<SupportedPollingEngineTypes
						xmlns:a="http://schemas.datacontract.org/2004/07/SolarWinds.Orion.Core.Models.Discovery">
						<a:DiscoveryPollingEngineType>Primary</a:DiscoveryPollingEngineType>
						<a:DiscoveryPollingEngineType>Additional</a:DiscoveryPollingEngineType>
					</SupportedPollingEngineTypes>
				</DiscoveryPluginInfo>
				<DiscoveryPluginInfo>
					<ModuleName i:nil="true"/>
					<PluginAssemblyName>SolarWinds.Interfaces.SNMP.DiscoveryJobPlugin</PluginAssemblyName>
      				<PluginTypeName>SolarWinds.Interfaces.SNMP.DiscoveryJobPlugin.InterfacesSnmpDiscoveryJobPlugin</PluginTypeName>
					<ProcessingOrder>2</ProcessingOrder>
      				<RootPath i:nil="true"/>
					<SupportedPollingEngineTypes
						xmlns:a="http://schemas.datacontract.org/2004/07/SolarWinds.Orion.Core.Models.Discovery">
						<a:DiscoveryPollingEngineType>Primary</a:DiscoveryPollingEngineType>
						<a:DiscoveryPollingEngineType>Additional</a:DiscoveryPollingEngineType>
					</SupportedPollingEngineTypes>
				</DiscoveryPluginInfo>
				<DiscoveryPluginInfo>
					<ModuleName i:nil="true"/>
					<PluginAssemblyName>SolarWinds.Interfaces.DiscoveryJobPlugin</PluginAssemblyName>
					<PluginTypeName>SolarWinds.Interfaces.DiscoveryJobPlugin.InterfacesDiscoveryJobPlugin</PluginTypeName>
					<ProcessingOrder>3</ProcessingOrder>
					<RootPath i:nil="true"/>
					<SupportedPollingEngineTypes
						xmlns:a="http://schemas.datacontract.org/2004/07/SolarWinds.Orion.Core.Models.Discovery">
						<a:DiscoveryPollingEngineType>Primary</a:DiscoveryPollingEngineType>
						<a:DiscoveryPollingEngineType>Additional</a:DiscoveryPollingEngineType>
					</SupportedPollingEngineTypes>
				</DiscoveryPluginInfo>
			</PluginInfos>
		</DiscoveryPluginInfoCollection>
	</data0>
	<data1>
		<DiscoveryJobDescription
			xmlns="http://schemas.solarwinds.com/2008/Orion"
			xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
			<DefaultProbes>Icmp Snmp</DefaultProbes>
			<DisableICMP>false</DisableICMP>
			<DiscoveryPluginJobDescriptions	xmlns:a="http://schemas.datacontract.org/2004/07/SolarWinds.Orion.Core.Models.Discovery">
				<a:DiscoveryPluginJobDescriptionBase i:type="CoreDiscoveryPluginJobDescription">
					<ActiveDirectoryList/>
					<AgentsAddresses xmlns:b="http://schemas.microsoft.com/2003/10/Serialization/Arrays"/>
                    <AgentsFilterQuery i:nil="true"/>
                    <CheckOnCertificateChainErrors>false</CheckOnCertificateChainErrors>
                    <CheckOnCertificateNameMismatch>false</CheckOnCertificateNameMismatch>
                    <CheckOnCertificateRevocation>false</CheckOnCertificateRevocation>
                    <Credentials>
                        <credentials>
                            <knownTypes>
                                <ArrayOfstring
                                    xmlns="http://schemas.microsoft.com/2003/10/Serialization/Arrays">
                                    <string>SolarWinds.Orion.Core.Models.Credentials.SnmpCredentialsV2,SolarWinds.Orion.Core.SharedCredentials.Provider</string>
                                </ArrayOfstring>
                            </knownTypes>
                            <pluginItem>
									<ArrayOfCredential>
										{{range .CredentialSnmpV2}}
										<Credential i:type="SnmpCredentialsV2">
											<ID>{{.Id}}</ID>
											<Name>{{.Name}}</Name>
											<Community>{{.Community}}</Community>
											<Description>Test credential</Description>
											<Owner>Core</Owner>
											<IsBroken>false</IsBroken>
										</Credential>
										{{end}}
									</ArrayOfCredential>
                            </pluginItem>
                        </credentials>
                    </Credentials>
                    <DiscoverAgentNodes>false</DiscoverAgentNodes>
                    <ExecutionTimeoutInMilliseconds>0</ExecutionTimeoutInMilliseconds>
                    <Port>0</Port>
                    <Targets>
							{{range .IpAddresses}}
							<DiscoveryTargetBase i:type="SingleIpTarget">
								<IP>{{.}}</IP>
							</DiscoveryTargetBase>
							{{end}}
                            {{range .Subnets}}
							<DiscoveryTargetBase i:type="SubnetTarget">
								<IP>{{.IpAddress}}</IP>
								<Mask>{{.Mask}}</Mask>
							</DiscoveryTargetBase>
							{{end}}
							{{range .IpRanges}}
							<DiscoveryTargetBase i:type="IpRangeTarget">
								<FirstIP>{{.FirstIpAddress}}</FirstIP>
								<LastIP>{{.LastIpAddress}}</LastIP>
							</DiscoveryTargetBase>
                            {{end}}
                    </Targets>
                    <UrlPrefix i:nil="true"/>
                    <UseHttps>false</UseHttps>
                    <WinRmAuthenticationMechanism>Default</WinRmAuthenticationMechanism>
                    <WindowsConnectionMode>WmiOnly</WindowsConnectionMode>
                    <WmiAuthenticationMode>Default</WmiAuthenticationMode>
                    <WmiAutoCorrectRDNSInconsistencies>false</WmiAutoCorrectRDNSInconsistencies>
                    <WmiRetries>0</WmiRetries>
                    <WmiRetryInterval>PT4S</WmiRetryInterval>
                    <WmiRootNamespaceOverrideAddresses i:nil="true" xmlns:b="http://schemas.microsoft.com/2003/10/Serialization/Arrays"/>
                </a:DiscoveryPluginJobDescriptionBase>
				<a:DiscoveryPluginJobDescriptionBase xmlns="http://schemas.solarwinds.com/2008/Interfaces" i:type="InterfacesDiscoveryPluginJobDescription">
				</a:DiscoveryPluginJobDescriptionBase>
				<a:DiscoveryPluginJobDescriptionBase xmlns="http://schemas.solarwinds.com/2008/Interfaces.SNMP" i:type="InterfacesSnmpDiscoveryPluginJobDescription">
				</a:DiscoveryPluginJobDescriptionBase>
            </DiscoveryPluginJobDescriptions>
            <EngineId>0</EngineId>
            <HopCount>0</HopCount>
            <IcmpTimeout>PT4S</IcmpTimeout>
            <MaxThreadsInDetectionPhase>20</MaxThreadsInDetectionPhase>
            <MaxThreadsInInventoryPhase>20</MaxThreadsInInventoryPhase>
            <PreferredDnsAddressFamily>0</PreferredDnsAddressFamily>
            <PreferredPollingMethod>SNMP</PreferredPollingMethod>
            <ProfileId i:nil="true"/>
            <SnmpConfiguration>
                <MaxSnmpReplies>2</MaxSnmpReplies>
                <PreferredSnmpVersion>SNMP2c</PreferredSnmpVersion>
                <SnmpPort>161</SnmpPort>
                <SnmpRetries>1</SnmpRetries>
                <SnmpTimeout>PT2S</SnmpTimeout>
            </SnmpConfiguration>
            <TagFilter xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays"/>
            <VulnerabilityCheckDisabled>false</VulnerabilityCheckDisabled>
        </DiscoveryJobDescription>
    </data1>
</root>
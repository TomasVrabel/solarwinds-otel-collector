package models

const credential_snmpv2 = `
<PollersJobCredentials xmlns="http://schemas.solarwinds.com/2008/Orion" xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
  <Credentials xmlns:a="http://schemas.microsoft.com/2003/10/Serialization/Arrays">
	<a:KeyValueOflongCredentialART3mlPg>
	  <a:Key>1</a:Key>
	  <a:Value i:type="SnmpCredentialsV2">
        <Description i:nil="true"/><ID i:nil="true"/>
		<IsBroken>false</IsBroken><Name i:nil="true"/>
        <Owner>Core</Owner>
		<Community>$Community$</Community>
      </a:Value>
    </a:KeyValueOflongCredentialART3mlPg>
  </Credentials>
</PollersJobCredentials>
`

func GetSnmpV2Credentials(variables vars) string {
	return applyVariables(credential_snmpv2, variables)
}

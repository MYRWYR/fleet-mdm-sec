package main

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/crewjam/saml"
	dsig "github.com/russellhaering/goxmldsig"
)

// In a real deployment, this key would need to be set on a per-instance basis (provided via
// config, or generated by Fleet) and kept secret.
var key = func() crypto.PrivateKey {
	b, _ := pem.Decode([]byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA0OhbMuizgtbFOfwbK7aURuXhZx6VRuAs3nNibiuifwCGz6u9
yy7bOR0P+zqN0YkjxaokqFgra7rXKCdeABmoLqCC0U+cGmLNwPOOA0PaD5q5xKhQ
4Me3rt/R9C4Ca6k3/OnkxnKwnogcsmdgs2l8liT3qVHP04Oc7Uymq2v09bGb6nPu
fOrkXS9F6mSClxHG/q59AGOWsXK1xzIRV1eu8W2SNdyeFVU1JHiQe444xLoPul5t
InWasKayFsPlJfWNc8EoU8COjNhfo/GovFTHVjh9oUR/gwEFVwifIHihRE0Hazn2
EQSLaOr2LM0TsRsQroFjmwSGgI+X2bfbMTqWOQIDAQABAoIBAFWZwDTeESBdrLcT
zHZe++cJLxE4AObn2LrWANEv5AeySYsyzjRBYObIN9IzrgTb8uJ900N/zVr5VkxH
xUa5PKbOcowd2NMfBTw5EEnaNbILLm+coHdanrNzVu59I9TFpAFoPavrNt/e2hNo
NMGPSdOkFi81LLl4xoadz/WR6O/7N2famM+0u7C2uBe+TrVwHyuqboYoidJDhO8M
w4WlY9QgAUhkPyzZqrl+VfF1aDTGVf4LJgaVevfFCas8Ws6DQX5q4QdIoV6/0vXi
B1M+aTnWjHuiIzjBMWhcYW2+I5zfwNWRXaxdlrYXRukGSdnyO+DH/FhHePJgmlkj
NInADDkCgYEA6MEQFOFSCc/ELXYWgStsrtIlJUcsLdLBsy1ocyQa2lkVUw58TouW
RciE6TjW9rp31pfQUnO2l6zOUC6LT9Jvlb9PSsyW+rvjtKB5PjJI6W0hjX41wEO6
fshFELMJd9W+Ezao2AsP2hZJ8McCF8no9e00+G4xTAyxHsNI2AFTCQcCgYEA5cWZ
JwNb4t7YeEajPt9xuYNUOQpjvQn1aGOV7KcwTx5ELP/Hzi723BxHs7GSdrLkkDmi
Gpb+mfL4wxCt0fK0i8GFQsRn5eusyq9hLqP/bmjpHoXe/1uajFbE1fZQR+2LX05N
3ATlKaH2hdfCJedFa4wf43+cl6Yhp6ZA0Yet1r8CgYEAwiu1j8W9G+RRA5/8/DtO
yrUTOfsbFws4fpLGDTA0mq0whf6Soy/96C90+d9qLaC3srUpnG9eB0CpSOjbXXbv
kdxseLkexwOR3bD2FHX8r4dUM2bzznZyEaxfOaQypN8SV5ME3l60Fbr8ajqLO288
wlTmGM5Mn+YCqOg/T7wjGmcCgYBpzNfdl/VafOROVbBbhgXWtzsz3K3aYNiIjbp+
MunStIwN8GUvcn6nEbqOaoiXcX4/TtpuxfJMLw4OvAJdtxUdeSmEee2heCijV6g3
ErrOOy6EqH3rNWHvlxChuP50cFQJuYOueO6QggyCyruSOnDDuc0BM0SGq6+5g5s7
H++S/wKBgQDIkqBtFr9UEf8d6JpkxS0RXDlhSMjkXmkQeKGFzdoJcYVFIwq8jTNB
nJrVIGs3GcBkqGic+i7rTO1YPkquv4dUuiIn+vKZVoO6b54f+oPBXd4S0BnuEqFE
rdKNuCZhiaE2XD9L/O9KP1fh5bfEcKwazQ23EvpJHBMm8BGC+/YZNw==
-----END RSA PRIVATE KEY-----`))
	k, _ := x509.ParsePKCS1PrivateKey(b.Bytes)
	return k
}()

var cert = func() *x509.Certificate {
	b, _ := pem.Decode([]byte(`-----BEGIN CERTIFICATE-----
MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNV
BAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5
NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEB
BQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8A
hs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+a
ucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWx
m+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6
D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURN
B2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0O
BBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56
zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5
pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uv
NONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEf
y/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL
/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsb
GFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTL
UzreO96WzlBBMtY=
-----END CERTIFICATE-----`))
	c, _ := x509.ParseCertificate(b.Bytes)
	return c
}()

// In a real deployment this would need to be configurable to use the actual metadata from the Okta
// instance. Right now it's hardcoded to the metadata from a development instance.
var spMeta = func() *saml.EntityDescriptor {
	var res saml.EntityDescriptor
	err := xml.Unmarshal([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<md:EntityDescriptor entityID="https://www.okta.com/saml2/service-provider/spdhkwkrikxbahheysxi" xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata"><md:SPSSODescriptor AuthnRequestsSigned="true" WantAssertionsSigned="true" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol"><md:KeyDescriptor use="encryption"><ds:KeyInfo xmlns:ds="http://www.w3.org/2000/09/xmldsig#"><ds:X509Data><ds:X509Certificate>MIIDqjCCApKgAwIBAgIGAY5iiKkkMA0GCSqGSIb3DQEBCwUAMIGVMQswCQYDVQQGEwJVUzETMBEG
A1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzENMAsGA1UECgwET2t0YTEU
MBIGA1UECwwLU1NPUHJvdmlkZXIxFjAUBgNVBAMMDXRyaWFsLTIxNzE2MDUxHDAaBgkqhkiG9w0B
CQEWDWluZm9Ab2t0YS5jb20wHhcNMjQwMzIxMTk0MDQ3WhcNMzQwMzIxMTk0MTQ3WjCBlTELMAkG
A1UEBhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28xDTAL
BgNVBAoMBE9rdGExFDASBgNVBAsMC1NTT1Byb3ZpZGVyMRYwFAYDVQQDDA10cmlhbC0yMTcxNjA1
MRwwGgYJKoZIhvcNAQkBFg1pbmZvQG9rdGEuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB
CgKCAQEAulSGykIpWxpSPYbFWi535nt5zxXx7ku+h6e4XX7W3DBIFb+D6MGldfridUHkuXrF1b9f
ZP+9mbD6D6jNLvjTgBeK3y2+Ty1NyFEoS5T5U9m37p0NE8r9kOUkjtnKpo9b3I1pQw9XRaQk0JYY
ZINhzd5K7xSa8tHf0F7PzWtzOC1SEuwXzm+bWgyLRaeBVWc1M6lLkar+uFpSBtUcUcV88DgdPx5g
8M8waNdvv92Qs+mL5PL6CKiWrfatTbXgSYbmf0ssZSgk4BOOSLeL49NL3ED0rEVUlL0TbypPgEIO
+UtGFgrcbCjfBiRI/4thV/Sup04bpfRE7s666QN5SJNzgwIDAQABMA0GCSqGSIb3DQEBCwUAA4IB
AQCf2z9mp/lwaD7AnoR2Xa3WcJCTKON02xrftlCpDBfj6thPYhZr1he0RZXXlzAn7pzj5la46Bwq
LdK6yKvzy8bMY0phaIeiNuTXXK62hj2SJuRDgyVJZyfYmw582h44g51Xt0qSq6pNiClaUGYf9WEa
dirsr/V9GtrMLSFuaqnDpceKoILg2F+YBTXWXu9vwSbtf1gccbjtgTc41jrsUfEMol3Yr/iBhXIg
j6+NHot/QROcy1tB7CWPFyJ8SxzcUe8cJoueBLKBWDF6Jdpy6y8TDdaDYwXKsCJBXWSHd9C6celq
DDJo/AsWSArqnGsEy7QjitygdVc6y+4X4ee3h+fM</ds:X509Certificate></ds:X509Data></ds:KeyInfo><md:EncryptionMethod Algorithm="http://www.w3.org/2001/04/xmlenc#aes128-cbc"/><md:EncryptionMethod Algorithm="http://www.w3.org/2001/04/xmlenc#aes192-cbc"/><md:EncryptionMethod Algorithm="http://www.w3.org/2001/04/xmlenc#aes256-cbc"/><md:EncryptionMethod Algorithm="http://www.w3.org/2009/xmlenc11#aes128-gcm"/><md:EncryptionMethod Algorithm="http://www.w3.org/2009/xmlenc11#aes256-gcm"/><md:EncryptionMethod Algorithm="http://www.w3.org/2001/04/xmlenc#tripledes-cbc"/></md:KeyDescriptor><md:KeyDescriptor use="signing"><ds:KeyInfo xmlns:ds="http://www.w3.org/2000/09/xmldsig#"><ds:X509Data><ds:X509Certificate>MIIDqjCCApKgAwIBAgIGAY5iiKkkMA0GCSqGSIb3DQEBCwUAMIGVMQswCQYDVQQGEwJVUzETMBEG
A1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzENMAsGA1UECgwET2t0YTEU
MBIGA1UECwwLU1NPUHJvdmlkZXIxFjAUBgNVBAMMDXRyaWFsLTIxNzE2MDUxHDAaBgkqhkiG9w0B
CQEWDWluZm9Ab2t0YS5jb20wHhcNMjQwMzIxMTk0MDQ3WhcNMzQwMzIxMTk0MTQ3WjCBlTELMAkG
A1UEBhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28xDTAL
BgNVBAoMBE9rdGExFDASBgNVBAsMC1NTT1Byb3ZpZGVyMRYwFAYDVQQDDA10cmlhbC0yMTcxNjA1
MRwwGgYJKoZIhvcNAQkBFg1pbmZvQG9rdGEuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB
CgKCAQEAulSGykIpWxpSPYbFWi535nt5zxXx7ku+h6e4XX7W3DBIFb+D6MGldfridUHkuXrF1b9f
ZP+9mbD6D6jNLvjTgBeK3y2+Ty1NyFEoS5T5U9m37p0NE8r9kOUkjtnKpo9b3I1pQw9XRaQk0JYY
ZINhzd5K7xSa8tHf0F7PzWtzOC1SEuwXzm+bWgyLRaeBVWc1M6lLkar+uFpSBtUcUcV88DgdPx5g
8M8waNdvv92Qs+mL5PL6CKiWrfatTbXgSYbmf0ssZSgk4BOOSLeL49NL3ED0rEVUlL0TbypPgEIO
+UtGFgrcbCjfBiRI/4thV/Sup04bpfRE7s666QN5SJNzgwIDAQABMA0GCSqGSIb3DQEBCwUAA4IB
AQCf2z9mp/lwaD7AnoR2Xa3WcJCTKON02xrftlCpDBfj6thPYhZr1he0RZXXlzAn7pzj5la46Bwq
LdK6yKvzy8bMY0phaIeiNuTXXK62hj2SJuRDgyVJZyfYmw582h44g51Xt0qSq6pNiClaUGYf9WEa
dirsr/V9GtrMLSFuaqnDpceKoILg2F+YBTXWXu9vwSbtf1gccbjtgTc41jrsUfEMol3Yr/iBhXIg
j6+NHot/QROcy1tB7CWPFyJ8SxzcUe8cJoueBLKBWDF6Jdpy6y8TDdaDYwXKsCJBXWSHd9C6celq
DDJo/AsWSArqnGsEy7QjitygdVc6y+4X4ee3h+fM</ds:X509Certificate></ds:X509Data></ds:KeyInfo></md:KeyDescriptor><md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat><md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress</md:NameIDFormat><md:NameIDFormat>urn:oasis:names:tc:SAML:2.0:nameid-format:persistent</md:NameIDFormat><md:NameIDFormat>urn:oasis:names:tc:SAML:2.0:nameid-format:transient</md:NameIDFormat><md:AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="https://trial-2171605.okta.com/sso/saml2/0oackhonw04QjK6ki697" index="0" isDefault="true"/><md:AttributeConsumingService index="0"><md:ServiceName xml:lang="en">Fleet</md:ServiceName><md:RequestedAttribute FriendlyName="First Name" Name="firstName" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:uri" isRequired="true"/><md:RequestedAttribute FriendlyName="Last Name" Name="lastName" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:uri" isRequired="true"/><md:RequestedAttribute FriendlyName="Email" Name="email" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:uri" isRequired="true"/><md:RequestedAttribute FriendlyName="Mobile Phone" Name="mobilePhone" NameFormat="urn:oasis:names:tc:SAML:2.0:attrname-format:uri" isRequired="false"/></md:AttributeConsumingService></md:SPSSODescriptor><md:Organization><md:OrganizationName xml:lang="en">trial-2171605</md:OrganizationName><md:OrganizationDisplayName xml:lang="en">fleetdm-trial-2171605</md:OrganizationDisplayName><md:OrganizationURL xml:lang="en">https://fleetdm.com</md:OrganizationURL></md:Organization></md:EntityDescriptor>`),
		&res)
	if err != nil {
		panic(err)
	}
	return &res
}()

type identityProvider struct {
	// embed IdentityProvider to override methods
	saml.IdentityProvider
	// this ServiceProvider is used to do the MFA with Okta
	mfaServiceProvider saml.ServiceProvider
	// requestMap holds a map with keys that are the service provider request ID (the ID that we
	// generate before sending a MFA service provider request TO Okta), and values that are the
	// original identity provider request (FROM Okta). This is used to look up the original request
	// when we receive the MFA response (with InResponseTo).
	requestMap map[string]*saml.IdpAuthnRequest
}

func (i *identityProvider) GetServiceProvider(r *http.Request, serviceProviderID string) (*saml.EntityDescriptor, error) {
	log.Printf("request id: %s", serviceProviderID)
	return spMeta, nil
}

func (i *identityProvider) GetSession(w http.ResponseWriter, r *http.Request, req *saml.IdpAuthnRequest) *saml.Session {
	email := req.Request.Subject.NameID.Value
	log.Printf("auth email: %+v", email)

	if err := checkForEmail(email); err != nil {
		// Check failed -- write HTML showing what the user needs to do. In a future iteration, we
		// would probably want to manually return the SAML redirect with assertion so that the user
		// doesn't have to reload the login page to complete the login after remediation.

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<html>
		<body>
		<img src="/img/failure.png" style="width:80%; display: block;  margin-left: auto;  margin-right: auto;" />
		</body>
		</html>`)
		return nil
	}

	// Check succeeded -- let the saml library write the SAML redirect with the assertion.
	return &saml.Session{
		NameID: email,
	}
}

func getRequestValues(r *http.Request) url.Values {
	switch r.Method {
	case "GET":
		return r.URL.Query()
	case "POST":
		if err := r.ParseForm(); err != nil {
			panic(err)
		}
		return r.PostForm
	default:
		panic("unsupported request method")
	}
}

// ServeSSO handles the initial SSO request.
//
// First we validate the SAML request and then we return HTML/JS that will collect the device
// identifier from the locally running agent. That JS then redirects back to the secondary SSO
// request where we can process both the user's email and the device identifier.
func (idp *identityProvider) ServeSSO(w http.ResponseWriter, r *http.Request) {
	req, err := saml.NewIdpAuthnRequest(&idp.IdentityProvider, r)
	if err != nil {
		idp.Logger.Printf("failed to parse request: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		idp.Logger.Printf("failed to validate request: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	type tmplParams struct {
		Email       string
		SAMLRequest string
		RelayState  string
	}
	t, err := template.New("").Parse(`
<html>
	<head>

	</head>
	<body>
		If you are seeing this, javascript is disabled or an error occurred. See the console for more.
	</body>
		<script type="text/javascript">
	(async () => {
		document.body.innerHTML = ""; // Clear JS disabled message

		const email = {{ .Email }};
		const saml = {{ .SAMLRequest }};
		const relay = {{ .RelayState }};
		// This is where we would want to try all the different predetermined localhost addresses in
		// case the default port is bound
		const resp = await fetch("http://localhost:9339/identifier");
		const json = await resp.json();
		console.log("got parameters: ", json.identifier, email, saml);

		// Now that we have the parameters, let's redirect to the second step
		const form = document.createElement('form');
  		form.method = "post";
  		form.action = "/ssowithidentifier";

		const identifierField = document.createElement("input");
		identifierField.type = "hidden";
		identifierField.name = "FleetIdentifier";
		identifierField.value = json.identifier;
		form.appendChild(identifierField);

		const samlField = document.createElement("input");
		samlField.type = "hidden";
		samlField.name = "SAMLRequest";
		samlField.value = saml;
		form.appendChild(samlField);

		const relayField = document.createElement("input");
		relayField.type = "hidden";
		relayField.name = "RelayState";
		relayField.value = relay;
		form.appendChild(relayField);

		document.body.appendChild(form);
  		form.submit();
	})();
	</script>
</html>
	`)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	values := getRequestValues(r)

	err = t.Execute(w, tmplParams{
		Email:       req.Request.Subject.NameID.Value,
		SAMLRequest: values.Get("SAMLRequest"),
		RelayState:  values.Get("RelayState"),
	})
	if err != nil {
		panic(err)
	}
}

// ServeSSO handles SAML auth requests after we've retrieved the device identifier
func (idp *identityProvider) ServeSSOWithIdentifier(w http.ResponseWriter, r *http.Request) {
	req, err := saml.NewIdpAuthnRequest(&idp.IdentityProvider, r)
	if err != nil {
		idp.Logger.Printf("failed to parse request: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		idp.Logger.Printf("failed to validate request: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	///////////

	// Here's where we know the device identifier and the user's email

	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	fleetID := r.PostForm.Get("FleetIdentifier")
	log.Println("got fleet ID: ", fleetID)

	// Show the error if the user needs to resolve policies
	if err := checkByIdentifier(fleetID); err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<html>
		<body>
		<img src="/img/failure.png" style="width:80%; display: block;  margin-left: auto;  margin-right: auto;" />
		</body>
		</html>`)
		return
	}

	// If we were not doing MFA, we could redirect to Okta now with a successful assertion.

	// Redirect to Okta as a service provider for MFA
	// Assume we are using HTTP Post binding. We could support other bindings in the future.
	bindingLocation := idp.mfaServiceProvider.GetSSOBindingLocation(saml.HTTPPostBinding)
	authReq, err := idp.mfaServiceProvider.MakeAuthenticationRequest(bindingLocation, saml.HTTPPostBinding, saml.HTTPPostBinding)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Track this request (in memory here, though we would probably want to do this in the database
	// for a production implementation). For production: Do we need to do anything else to prevent
	// replay attacks?
	idp.requestMap[authReq.ID] = req
	w.Header().Add("Content-Security-Policy", ""+
		"default-src; "+
		"script-src 'sha256-AjPdJSbZmeWHnEc5ykvJFay8FTWeTeRbs9dutfZ0HqE='; "+
		"reflected-xss block; referrer no-referrer;")
	w.Header().Add("Content-type", "text/html")
	var buf bytes.Buffer
	buf.WriteString(`<!DOCTYPE html><html><body>`)
	buf.Write(authReq.Post(authReq.ID))
	buf.WriteString(`</body></html>`)
	if _, err := w.Write(buf.Bytes()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func (idp *identityProvider) ServeMFAResponse(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// It seems a little sketchy to allow the RelayState from the response to be used as one of the
	// possible InResponseTo IDs (because the request could be crafted to make those match). I think
	// this is okay because we check below that the RelayState actually corresponds to a request
	// that we initiated by looking it up in the requestMap. Someone should verify this.
	relayState := r.Form.Get("RelayState")
	assertion, err := idp.mfaServiceProvider.ParseResponse(r, []string{relayState})
	if err != nil {
		if ierr, ok := err.(*saml.InvalidResponseError); ok {
			log.Println(ierr.PrivateErr.Error())
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check that we initiated the MFA request and look up the originating IdP request
	req, ok := idp.requestMap[relayState]
	if !ok {
		http.Error(w, "did not find matching request", http.StatusForbidden)
		return
	}
	delete(idp.requestMap, relayState)

	// Validate that the same user completed the MFA request as initiated the original IdP request.
	if assertion.Subject.NameID.Value != req.Request.Subject.NameID.Value {
		http.Error(w, "name IDs do not match", http.StatusForbidden)
		return
	}

	// If everything is good up to this point, complete the original IdP request back to Okta

	session := &saml.Session{
		NameID: req.Request.Subject.NameID.Value,
	}

	assertionMaker := idp.AssertionMaker
	if assertionMaker == nil {
		assertionMaker = saml.DefaultAssertionMaker{}
	}
	if err := assertionMaker.MakeAssertion(req, session); err != nil {
		idp.Logger.Printf("failed to make assertion: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if err := req.WriteResponse(w); err != nil {
		idp.Logger.Printf("failed to write response: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func getIDP() *identityProvider {
	u, err := url.Parse("http://" + addr)
	if err != nil {
		panic(err)
	}
	metadataURL := *u
	metadataURL.Path = metadataURL.Path + "/metadata"
	ssoURL := *u
	ssoURL.Path = ssoURL.Path + "/sso"

	i := &identityProvider{
		mfaServiceProvider: getSP(),
		requestMap:         make(map[string]*saml.IdpAuthnRequest),
	}
	idp := saml.IdentityProvider{
		Key:                     key,
		SignatureMethod:         dsig.RSASHA256SignatureMethod,
		Logger:                  log.New(os.Stderr, "idp: ", 0),
		Certificate:             cert,
		MetadataURL:             metadataURL,
		SSOURL:                  ssoURL,
		ServiceProviderProvider: i,
		SessionProvider:         i,
	}
	i.IdentityProvider = idp
	return i
}
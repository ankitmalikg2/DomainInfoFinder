package models

//whois 020ywzs.com | grep -e 'Registrant Name' -e 'Registrant Organization' -e 'Registrant Street' -e 'Registrant City' -e 'Registrant State/Province' -e 'Registrant Postal Code' -e 'Registrant Country' -e 'Registrant Phone Ext' -e 'Registrant Fax' -e 'Registrant Fax Ext' -e 'Registrant Email'
type DomainInfo struct {
	RegistrantName       string
	RegistrantOrg        string
	RegistrantStreet     string
	RegistrantCity       string
	RegistrantProvince   string
	RegistrantPostalCode string
	RegistrantCountry    string
	RegistrantPhone      string
	RegistrantPhoneExt   string
	RegistrantFax        string
	RegistrantFaxExt     string
	RegistrantEmail      string
}

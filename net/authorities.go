package net

// Authorities is the hardcoded set of GOBL Net Addresses considered
// trusted KYC vendors. An identity returned by Client.WhoIs is only
// considered endorsed if its /who envelope is signed by one of these.
//
// The list is intentionally empty in this release; entries are added
// here as vendors come online.
var Authorities = []Address{}

// RegisterAuthority adds an address to the global set of trusted
// authority addresses.
func RegisterAuthority(addr Address) {
	Authorities = append(Authorities, addr)
}

package untdid

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

const (
	// ExtKeyCharge is used to identify the UNTDID 7161 charge codes.
	ExtKeyCharge cbc.Key = "untdid-charge"
)

var extCharge = &cbc.KeyDefinition{
	Key:  ExtKeyCharge,
	Name: i18n.NewString("UNTDID 7161 Charge"),
	Desc: i18n.String{
		i18n.EN: here.Doc(`
			UNTDID 7161 code used to describe the charge. List is based on the
			EN16931 code lists with extensions for taxes and duties.
		`),
	},
	Values: []*cbc.ValueDefinition{
		{
			Value: "AA",
			Name:  i18n.NewString("Advertising"),
		},
		{
			Value: "AAA",
			Name:  i18n.NewString("Telecommunication"),
		},
		{
			Value: "AAC",
			Name:  i18n.NewString("Technical modification"),
		},
		{
			Value: "AAD",
			Name:  i18n.NewString("Job-order production"),
		},
		{
			Value: "AAE",
			Name:  i18n.NewString("Outlays"),
		},
		{
			Value: "AAF",
			Name:  i18n.NewString("Off-premises"),
		},
		{
			Value: "AAH",
			Name:  i18n.NewString("Additional processing"),
		},
		{
			Value: "AAI",
			Name:  i18n.NewString("Attesting"),
		},
		{
			Value: "AAS",
			Name:  i18n.NewString("Acceptance"),
		},
		{
			Value: "AAT",
			Name:  i18n.NewString("Rush delivery"),
		},
		{
			Value: "AAV",
			Name:  i18n.NewString("Special construction"),
		},
		{
			Value: "AAY",
			Name:  i18n.NewString("Airport facilities"),
		},
		{
			Value: "AAZ",
			Name:  i18n.NewString("Concession"),
		},
		{
			Value: "ABA",
			Name:  i18n.NewString("Compulsory storage"),
		},
		{
			Value: "ABB",
			Name:  i18n.NewString("Fuel removal"),
		},
		{
			Value: "ABC",
			Name:  i18n.NewString("Into plane"),
		},
		{
			Value: "ABD",
			Name:  i18n.NewString("Overtime"),
		},
		{
			Value: "ABF",
			Name:  i18n.NewString("Tooling"),
		},
		{
			Value: "ABK",
			Name:  i18n.NewString("Miscellaneous"),
		},
		{
			Value: "ABL",
			Name:  i18n.NewString("Additional packaging"),
		},
		{
			Value: "ABN",
			Name:  i18n.NewString("Dunnage"),
		},
		{
			Value: "ABR",
			Name:  i18n.NewString("Containerisation"),
		},
		{
			Value: "ABS",
			Name:  i18n.NewString("Carton packing"),
		},
		{
			Value: "ABT",
			Name:  i18n.NewString("Hessian wrapped"),
		},
		{
			Value: "ABU",
			Name:  i18n.NewString("Polyethylene wrap packing"),
		},
		{
			Value: "ABW", // not in EN16931
			Name:  i18n.NewString("Customs duty charge"),
		},
		{
			Value: "ACF",
			Name:  i18n.NewString("Miscellaneous treatment"),
		},
		{
			Value: "ACG",
			Name:  i18n.NewString("Enamelling treatment"),
		},
		{
			Value: "ACH",
			Name:  i18n.NewString("Heat treatment"),
		},
		{
			Value: "ACI",
			Name:  i18n.NewString("Plating treatment"),
		},
		{
			Value: "ACJ",
			Name:  i18n.NewString("Painting"),
		},
		{
			Value: "ACK",
			Name:  i18n.NewString("Polishing"),
		},
		{
			Value: "ACL",
			Name:  i18n.NewString("Priming"),
		},
		{
			Value: "ACM",
			Name:  i18n.NewString("Preservation treatment"),
		},
		{
			Value: "ACS",
			Name:  i18n.NewString("Fitting"),
		},
		{
			Value: "ADC",
			Name:  i18n.NewString("Consolidation"),
		},
		{
			Value: "ADE",
			Name:  i18n.NewString("Bill of lading"),
		},
		{
			Value: "ADJ",
			Name:  i18n.NewString("Airbag"),
		},
		{
			Value: "ADK",
			Name:  i18n.NewString("Transfer"),
		},
		{
			Value: "ADL",
			Name:  i18n.NewString("Slipsheet"),
		},
		{
			Value: "ADM",
			Name:  i18n.NewString("Binding"),
		},
		{
			Value: "ADN",
			Name:  i18n.NewString("Repair or replacement of broken returnable package"),
		},
		{
			Value: "ADO",
			Name:  i18n.NewString("Efficient logistics"),
		},
		{
			Value: "ADP",
			Name:  i18n.NewString("Merchandising"),
		},
		{
			Value: "ADQ",
			Name:  i18n.NewString("Product mix"),
		},
		{
			Value: "ADR",
			Name:  i18n.NewString("Other services"),
		},
		{
			Value: "ADT",
			Name:  i18n.NewString("Pick-up"),
		},
		{
			Value: "ADW",
			Name:  i18n.NewString("Chronic illness"),
		},
		{
			Value: "ADY",
			Name:  i18n.NewString("New product introduction"),
		},
		{
			Value: "ADZ",
			Name:  i18n.NewString("Direct delivery"),
		},
		{
			Value: "AEA",
			Name:  i18n.NewString("Diversion"),
		},
		{
			Value: "AEB",
			Name:  i18n.NewString("Disconnect"),
		},
		{
			Value: "AEC",
			Name:  i18n.NewString("Distribution"),
		},
		{
			Value: "AED",
			Name:  i18n.NewString("Handling of hazardous cargo"),
		},
		{
			Value: "AEF",
			Name:  i18n.NewString("Rents and leases"),
		},
		{
			Value: "AEH",
			Name:  i18n.NewString("Location differential"),
		},
		{
			Value: "AEI",
			Name:  i18n.NewString("Aircraft refueling"),
		},
		{
			Value: "AEJ",
			Name:  i18n.NewString("Fuel shipped into storage"),
		},
		{
			Value: "AEK",
			Name:  i18n.NewString("Cash on delivery"),
		},
		{
			Value: "AEL",
			Name:  i18n.NewString("Small order processing service"),
		},
		{
			Value: "AEM",
			Name:  i18n.NewString("Clerical or administrative services"),
		},
		{
			Value: "AEN",
			Name:  i18n.NewString("Guarantee"),
		},
		{
			Value: "AEO",
			Name:  i18n.NewString("Collection and recycling"),
		},
		{
			Value: "AEP",
			Name:  i18n.NewString("Copyright fee collection"),
		},
		{
			Value: "AES",
			Name:  i18n.NewString("Veterinary inspection service"),
		},
		{
			Value: "AET",
			Name:  i18n.NewString("Pensioner service"),
		},
		{
			Value: "AEU",
			Name:  i18n.NewString("Medicine free pass holder"),
		},
		{
			Value: "AEV",
			Name:  i18n.NewString("Environmental protection service"),
		},
		{
			Value: "AEW",
			Name:  i18n.NewString("Environmental clean-up service"),
		},
		{
			Value: "AEX",
			Name:  i18n.NewString("National cheque processing service outside account area"),
		},
		{
			Value: "AEY",
			Name:  i18n.NewString("National payment service outside account area"),
		},
		{
			Value: "AEZ",
			Name:  i18n.NewString("National payment service within account area"),
		},
		{
			Value: "AJ",
			Name:  i18n.NewString("Adjustments"),
		},
		{
			Value: "AU",
			Name:  i18n.NewString("Authentication"),
		},
		{
			Value: "CA",
			Name:  i18n.NewString("Cataloguing"),
		},
		{
			Value: "CAB",
			Name:  i18n.NewString("Cartage"),
		},
		{
			Value: "CAD",
			Name:  i18n.NewString("Certification"),
		},
		{
			Value: "CAE",
			Name:  i18n.NewString("Certificate of conformance"),
		},
		{
			Value: "CAF",
			Name:  i18n.NewString("Certificate of origin"),
		},
		{
			Value: "CAI",
			Name:  i18n.NewString("Cutting"),
		},
		{
			Value: "CAJ",
			Name:  i18n.NewString("Consular service"),
		},
		{
			Value: "CAK",
			Name:  i18n.NewString("Customer collection"),
		},
		{
			Value: "CAL",
			Name:  i18n.NewString("Payroll payment service"),
		},
		{
			Value: "CAM",
			Name:  i18n.NewString("Cash transportation"),
		},
		{
			Value: "CAN",
			Name:  i18n.NewString("Home banking service"),
		},
		{
			Value: "CAO",
			Name:  i18n.NewString("Bilateral agreement service"),
		},
		{
			Value: "CAP",
			Name:  i18n.NewString("Insurance brokerage service"),
		},
		{
			Value: "CAQ",
			Name:  i18n.NewString("Cheque generation"),
		},
		{
			Value: "CAR",
			Name:  i18n.NewString("Preferential merchandising location"),
		},
		{
			Value: "CAS",
			Name:  i18n.NewString("Crane"),
		},
		{
			Value: "CAT",
			Name:  i18n.NewString("Special colour service"),
		},
		{
			Value: "CAU",
			Name:  i18n.NewString("Sorting"),
		},
		{
			Value: "CAV",
			Name:  i18n.NewString("Battery collection and recycling"),
		},
		{
			Value: "CAW",
			Name:  i18n.NewString("Product take back fee"),
		},
		{
			Value: "CAX",
			Name:  i18n.NewString("Quality control released"),
		},
		{
			Value: "CAY",
			Name:  i18n.NewString("Quality control held"),
		},
		{
			Value: "CAZ",
			Name:  i18n.NewString("Quality control embargo"),
		},
		{
			Value: "CD",
			Name:  i18n.NewString("Car loading"),
		},
		{
			Value: "CG",
			Name:  i18n.NewString("Cleaning"),
		},
		{
			Value: "CS",
			Name:  i18n.NewString("Cigarette stamping"),
		},
		{
			Value: "CT",
			Name:  i18n.NewString("Count and recount"),
		},
		{
			Value: "DAB",
			Name:  i18n.NewString("Layout/design"),
		},
		{
			Value: "DAC",
			Name:  i18n.NewString("Assortment allowance"),
		},
		{
			Value: "DAD",
			Name:  i18n.NewString("Driver assigned unloading"),
		},
		{
			Value: "DAF",
			Name:  i18n.NewString("Debtor bound"),
		},
		{
			Value: "DAG",
			Name:  i18n.NewString("Dealer allowance"),
		},
		{
			Value: "DAH",
			Name:  i18n.NewString("Allowance transferable to the consumer"),
		},
		{
			Value: "DAI",
			Name:  i18n.NewString("Growth of business"),
		},
		{
			Value: "DAJ",
			Name:  i18n.NewString("Introduction allowance"),
		},
		{
			Value: "DAK",
			Name:  i18n.NewString("Multi-buy promotion"),
		},
		{
			Value: "DAL",
			Name:  i18n.NewString("Partnership"),
		},
		{
			Value: "DAM",
			Name:  i18n.NewString("Return handling"),
		},
		{
			Value: "DAN",
			Name:  i18n.NewString("Minimum order not fulfilled charge"),
		},
		{
			Value: "DAO",
			Name:  i18n.NewString("Point of sales threshold allowance"),
		},
		{
			Value: "DAP",
			Name:  i18n.NewString("Wholesaling discount"),
		},
		{
			Value: "DAQ",
			Name:  i18n.NewString("Documentary credits transfer commission"),
		},
		{
			Value: "DL",
			Name:  i18n.NewString("Delivery"),
		},
		{
			Value: "EG",
			Name:  i18n.NewString("Engraving"),
		},
		{
			Value: "EP",
			Name:  i18n.NewString("Expediting"),
		},
		{
			Value: "ER",
			Name:  i18n.NewString("Exchange rate guarantee"),
		},
		{
			Value: "FAA",
			Name:  i18n.NewString("Fabrication"),
		},
		{
			Value: "FAB",
			Name:  i18n.NewString("Freight equalization"),
		},
		{
			Value: "FAC",
			Name:  i18n.NewString("Freight extraordinary handling"),
		},
		{
			Value: "FC",
			Name:  i18n.NewString("Freight service"),
		},
		{
			Value: "FH",
			Name:  i18n.NewString("Filling/handling"),
		},
		{
			Value: "FI",
			Name:  i18n.NewString("Financing"),
		},
		{
			Value: "GAA",
			Name:  i18n.NewString("Grinding"),
		},
		{
			Value: "HAA",
			Name:  i18n.NewString("Hose"),
		},
		{
			Value: "HD",
			Name:  i18n.NewString("Handling"),
		},
		{
			Value: "HH",
			Name:  i18n.NewString("Hoisting and hauling"),
		},
		{
			Value: "IAA",
			Name:  i18n.NewString("Installation"),
		},
		{
			Value: "IAB",
			Name:  i18n.NewString("Installation and warranty"),
		},
		{
			Value: "ID",
			Name:  i18n.NewString("Inside delivery"),
		},
		{
			Value: "IF",
			Name:  i18n.NewString("Inspection"),
		},
		{
			Value: "IN", // not in EN16931
			Name:  i18n.NewString("Insurance"),
		},
		{
			Value: "IR",
			Name:  i18n.NewString("Installation and training"),
		},
		{
			Value: "IS",
			Name:  i18n.NewString("Invoicing"),
		},
		{
			Value: "KO",
			Name:  i18n.NewString("Koshering"),
		},
		{
			Value: "L1",
			Name:  i18n.NewString("Carrier count"),
		},
		{
			Value: "LA",
			Name:  i18n.NewString("Labelling"),
		},
		{
			Value: "LAA",
			Name:  i18n.NewString("Labour"),
		},
		{
			Value: "LAB",
			Name:  i18n.NewString("Repair and return"),
		},
		{
			Value: "LF",
			Name:  i18n.NewString("Legalisation"),
		},
		{
			Value: "MAE",
			Name:  i18n.NewString("Mounting"),
		},
		{
			Value: "MI",
			Name:  i18n.NewString("Mail invoice"),
		},
		{
			Value: "ML",
			Name:  i18n.NewString("Mail invoice to each location"),
		},
		{
			Value: "NAA",
			Name:  i18n.NewString("Non-returnable containers"),
		},
		{
			Value: "OA",
			Name:  i18n.NewString("Outside cable connectors"),
		},
		{
			Value: "PA",
			Name:  i18n.NewString("Invoice with shipment"),
		},
		{
			Value: "PAA",
			Name:  i18n.NewString("Phosphatizing (steel treatment)"),
		},
		{
			Value: "PC",
			Name:  i18n.NewString("Packing"),
		},
		{
			Value: "PL",
			Name:  i18n.NewString("Palletizing"),
		},
		{
			Value: "PRV",
			Name:  i18n.NewString("Price variation"),
		},
		{
			Value: "RAB",
			Name:  i18n.NewString("Repacking"),
		},
		{
			Value: "RAC",
			Name:  i18n.NewString("Repair"),
		},
		{
			Value: "RAD",
			Name:  i18n.NewString("Returnable container"),
		},
		{
			Value: "RAF",
			Name:  i18n.NewString("Restocking"),
		},
		{
			Value: "RE",
			Name:  i18n.NewString("Re-delivery"),
		},
		{
			Value: "RF",
			Name:  i18n.NewString("Refurbishing"),
		},
		{
			Value: "RH",
			Name:  i18n.NewString("Rail wagon hire"),
		},
		{
			Value: "RV",
			Name:  i18n.NewString("Loading"),
		},
		{
			Value: "SA",
			Name:  i18n.NewString("Salvaging"),
		},
		{
			Value: "SAA",
			Name:  i18n.NewString("Shipping and handling"),
		},
		{
			Value: "SAD",
			Name:  i18n.NewString("Special packaging"),
		},
		{
			Value: "SAE",
			Name:  i18n.NewString("Stamping"),
		},
		{
			Value: "SAI",
			Name:  i18n.NewString("Consignee unload"),
		},
		{
			Value: "SG",
			Name:  i18n.NewString("Shrink-wrap"),
		},
		{
			Value: "SH",
			Name:  i18n.NewString("Special handling"),
		},
		{
			Value: "SM",
			Name:  i18n.NewString("Special finish"),
		},
		{
			Value: "ST", // not in EN16931
			Name:  i18n.NewString("Stamp duties"),
		},
		{
			Value: "SU",
			Name:  i18n.NewString("Set-up"),
		},
		{
			Value: "TAB",
			Name:  i18n.NewString("Tank renting"),
		},
		{
			Value: "TAC",
			Name:  i18n.NewString("Testing"),
		},
		{
			Value: "TT",
			Name:  i18n.NewString("Transportation - third party billing"),
		},
		{
			Value: "TV",
			Name:  i18n.NewString("Transportation by vendor"),
		},
		{
			Value: "TX", // not in EN16931
			Name:  i18n.NewString("Tax"),
		},
		{
			Value: "V1",
			Name:  i18n.NewString("Drop yard"),
		},
		{
			Value: "V2",
			Name:  i18n.NewString("Drop dock"),
		},
		{
			Value: "WH",
			Name:  i18n.NewString("Warehousing"),
		},
		{
			Value: "XAA",
			Name:  i18n.NewString("Combine all same day shipment"),
		},
		{
			Value: "YY",
			Name:  i18n.NewString("Split pick-up"),
		},
		{
			Value: "ZZZ",
			Name:  i18n.NewString("Mutually defined"),
		},
	},
}

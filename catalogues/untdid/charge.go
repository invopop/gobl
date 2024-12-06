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

var extCharge = &cbc.Definition{
	Key:  ExtKeyCharge,
	Name: i18n.NewString("UNTDID 7161 Charge"),
	Desc: i18n.String{
		i18n.EN: here.Doc(`
			UNTDID 7161 code used to describe the charge. List is based on the
			EN16931 code lists with extensions for taxes and duties.
		`),
	},
	Values: []*cbc.Definition{
		{
			Code: "AA",
			Name: i18n.NewString("Advertising"),
		},
		{
			Code: "AAA",
			Name: i18n.NewString("Telecommunication"),
		},
		{
			Code: "AAC",
			Name: i18n.NewString("Technical modification"),
		},
		{
			Code: "AAD",
			Name: i18n.NewString("Job-order production"),
		},
		{
			Code: "AAE",
			Name: i18n.NewString("Outlays"),
		},
		{
			Code: "AAF",
			Name: i18n.NewString("Off-premises"),
		},
		{
			Code: "AAH",
			Name: i18n.NewString("Additional processing"),
		},
		{
			Code: "AAI",
			Name: i18n.NewString("Attesting"),
		},
		{
			Code: "AAS",
			Name: i18n.NewString("Acceptance"),
		},
		{
			Code: "AAT",
			Name: i18n.NewString("Rush delivery"),
		},
		{
			Code: "AAV",
			Name: i18n.NewString("Special construction"),
		},
		{
			Code: "AAY",
			Name: i18n.NewString("Airport facilities"),
		},
		{
			Code: "AAZ",
			Name: i18n.NewString("Concession"),
		},
		{
			Code: "ABA",
			Name: i18n.NewString("Compulsory storage"),
		},
		{
			Code: "ABB",
			Name: i18n.NewString("Fuel removal"),
		},
		{
			Code: "ABC",
			Name: i18n.NewString("Into plane"),
		},
		{
			Code: "ABD",
			Name: i18n.NewString("Overtime"),
		},
		{
			Code: "ABF",
			Name: i18n.NewString("Tooling"),
		},
		{
			Code: "ABK",
			Name: i18n.NewString("Miscellaneous"),
		},
		{
			Code: "ABL",
			Name: i18n.NewString("Additional packaging"),
		},
		{
			Code: "ABN",
			Name: i18n.NewString("Dunnage"),
		},
		{
			Code: "ABR",
			Name: i18n.NewString("Containerisation"),
		},
		{
			Code: "ABS",
			Name: i18n.NewString("Carton packing"),
		},
		{
			Code: "ABT",
			Name: i18n.NewString("Hessian wrapped"),
		},
		{
			Code: "ABU",
			Name: i18n.NewString("Polyethylene wrap packing"),
		},
		{
			Code: "ABW", // not in EN16931
			Name: i18n.NewString("Customs duty charge"),
		},
		{
			Code: "ACF",
			Name: i18n.NewString("Miscellaneous treatment"),
		},
		{
			Code: "ACG",
			Name: i18n.NewString("Enamelling treatment"),
		},
		{
			Code: "ACH",
			Name: i18n.NewString("Heat treatment"),
		},
		{
			Code: "ACI",
			Name: i18n.NewString("Plating treatment"),
		},
		{
			Code: "ACJ",
			Name: i18n.NewString("Painting"),
		},
		{
			Code: "ACK",
			Name: i18n.NewString("Polishing"),
		},
		{
			Code: "ACL",
			Name: i18n.NewString("Priming"),
		},
		{
			Code: "ACM",
			Name: i18n.NewString("Preservation treatment"),
		},
		{
			Code: "ACS",
			Name: i18n.NewString("Fitting"),
		},
		{
			Code: "ADC",
			Name: i18n.NewString("Consolidation"),
		},
		{
			Code: "ADE",
			Name: i18n.NewString("Bill of lading"),
		},
		{
			Code: "ADJ",
			Name: i18n.NewString("Airbag"),
		},
		{
			Code: "ADK",
			Name: i18n.NewString("Transfer"),
		},
		{
			Code: "ADL",
			Name: i18n.NewString("Slipsheet"),
		},
		{
			Code: "ADM",
			Name: i18n.NewString("Binding"),
		},
		{
			Code: "ADN",
			Name: i18n.NewString("Repair or replacement of broken returnable package"),
		},
		{
			Code: "ADO",
			Name: i18n.NewString("Efficient logistics"),
		},
		{
			Code: "ADP",
			Name: i18n.NewString("Merchandising"),
		},
		{
			Code: "ADQ",
			Name: i18n.NewString("Product mix"),
		},
		{
			Code: "ADR",
			Name: i18n.NewString("Other services"),
		},
		{
			Code: "ADT",
			Name: i18n.NewString("Pick-up"),
		},
		{
			Code: "ADW",
			Name: i18n.NewString("Chronic illness"),
		},
		{
			Code: "ADY",
			Name: i18n.NewString("New product introduction"),
		},
		{
			Code: "ADZ",
			Name: i18n.NewString("Direct delivery"),
		},
		{
			Code: "AEA",
			Name: i18n.NewString("Diversion"),
		},
		{
			Code: "AEB",
			Name: i18n.NewString("Disconnect"),
		},
		{
			Code: "AEC",
			Name: i18n.NewString("Distribution"),
		},
		{
			Code: "AED",
			Name: i18n.NewString("Handling of hazardous cargo"),
		},
		{
			Code: "AEF",
			Name: i18n.NewString("Rents and leases"),
		},
		{
			Code: "AEH",
			Name: i18n.NewString("Location differential"),
		},
		{
			Code: "AEI",
			Name: i18n.NewString("Aircraft refueling"),
		},
		{
			Code: "AEJ",
			Name: i18n.NewString("Fuel shipped into storage"),
		},
		{
			Code: "AEK",
			Name: i18n.NewString("Cash on delivery"),
		},
		{
			Code: "AEL",
			Name: i18n.NewString("Small order processing service"),
		},
		{
			Code: "AEM",
			Name: i18n.NewString("Clerical or administrative services"),
		},
		{
			Code: "AEN",
			Name: i18n.NewString("Guarantee"),
		},
		{
			Code: "AEO",
			Name: i18n.NewString("Collection and recycling"),
		},
		{
			Code: "AEP",
			Name: i18n.NewString("Copyright fee collection"),
		},
		{
			Code: "AES",
			Name: i18n.NewString("Veterinary inspection service"),
		},
		{
			Code: "AET",
			Name: i18n.NewString("Pensioner service"),
		},
		{
			Code: "AEU",
			Name: i18n.NewString("Medicine free pass holder"),
		},
		{
			Code: "AEV",
			Name: i18n.NewString("Environmental protection service"),
		},
		{
			Code: "AEW",
			Name: i18n.NewString("Environmental clean-up service"),
		},
		{
			Code: "AEX",
			Name: i18n.NewString("National cheque processing service outside account area"),
		},
		{
			Code: "AEY",
			Name: i18n.NewString("National payment service outside account area"),
		},
		{
			Code: "AEZ",
			Name: i18n.NewString("National payment service within account area"),
		},
		{
			Code: "AJ",
			Name: i18n.NewString("Adjustments"),
		},
		{
			Code: "AU",
			Name: i18n.NewString("Authentication"),
		},
		{
			Code: "CA",
			Name: i18n.NewString("Cataloguing"),
		},
		{
			Code: "CAB",
			Name: i18n.NewString("Cartage"),
		},
		{
			Code: "CAD",
			Name: i18n.NewString("Certification"),
		},
		{
			Code: "CAE",
			Name: i18n.NewString("Certificate of conformance"),
		},
		{
			Code: "CAF",
			Name: i18n.NewString("Certificate of origin"),
		},
		{
			Code: "CAI",
			Name: i18n.NewString("Cutting"),
		},
		{
			Code: "CAJ",
			Name: i18n.NewString("Consular service"),
		},
		{
			Code: "CAK",
			Name: i18n.NewString("Customer collection"),
		},
		{
			Code: "CAL",
			Name: i18n.NewString("Payroll payment service"),
		},
		{
			Code: "CAM",
			Name: i18n.NewString("Cash transportation"),
		},
		{
			Code: "CAN",
			Name: i18n.NewString("Home banking service"),
		},
		{
			Code: "CAO",
			Name: i18n.NewString("Bilateral agreement service"),
		},
		{
			Code: "CAP",
			Name: i18n.NewString("Insurance brokerage service"),
		},
		{
			Code: "CAQ",
			Name: i18n.NewString("Cheque generation"),
		},
		{
			Code: "CAR",
			Name: i18n.NewString("Preferential merchandising location"),
		},
		{
			Code: "CAS",
			Name: i18n.NewString("Crane"),
		},
		{
			Code: "CAT",
			Name: i18n.NewString("Special colour service"),
		},
		{
			Code: "CAU",
			Name: i18n.NewString("Sorting"),
		},
		{
			Code: "CAV",
			Name: i18n.NewString("Battery collection and recycling"),
		},
		{
			Code: "CAW",
			Name: i18n.NewString("Product take back fee"),
		},
		{
			Code: "CAX",
			Name: i18n.NewString("Quality control released"),
		},
		{
			Code: "CAY",
			Name: i18n.NewString("Quality control held"),
		},
		{
			Code: "CAZ",
			Name: i18n.NewString("Quality control embargo"),
		},
		{
			Code: "CD",
			Name: i18n.NewString("Car loading"),
		},
		{
			Code: "CG",
			Name: i18n.NewString("Cleaning"),
		},
		{
			Code: "CS",
			Name: i18n.NewString("Cigarette stamping"),
		},
		{
			Code: "CT",
			Name: i18n.NewString("Count and recount"),
		},
		{
			Code: "DAB",
			Name: i18n.NewString("Layout/design"),
		},
		{
			Code: "DAC",
			Name: i18n.NewString("Assortment allowance"),
		},
		{
			Code: "DAD",
			Name: i18n.NewString("Driver assigned unloading"),
		},
		{
			Code: "DAF",
			Name: i18n.NewString("Debtor bound"),
		},
		{
			Code: "DAG",
			Name: i18n.NewString("Dealer allowance"),
		},
		{
			Code: "DAH",
			Name: i18n.NewString("Allowance transferable to the consumer"),
		},
		{
			Code: "DAI",
			Name: i18n.NewString("Growth of business"),
		},
		{
			Code: "DAJ",
			Name: i18n.NewString("Introduction allowance"),
		},
		{
			Code: "DAK",
			Name: i18n.NewString("Multi-buy promotion"),
		},
		{
			Code: "DAL",
			Name: i18n.NewString("Partnership"),
		},
		{
			Code: "DAM",
			Name: i18n.NewString("Return handling"),
		},
		{
			Code: "DAN",
			Name: i18n.NewString("Minimum order not fulfilled charge"),
		},
		{
			Code: "DAO",
			Name: i18n.NewString("Point of sales threshold allowance"),
		},
		{
			Code: "DAP",
			Name: i18n.NewString("Wholesaling discount"),
		},
		{
			Code: "DAQ",
			Name: i18n.NewString("Documentary credits transfer commission"),
		},
		{
			Code: "DL",
			Name: i18n.NewString("Delivery"),
		},
		{
			Code: "EG",
			Name: i18n.NewString("Engraving"),
		},
		{
			Code: "EP",
			Name: i18n.NewString("Expediting"),
		},
		{
			Code: "ER",
			Name: i18n.NewString("Exchange rate guarantee"),
		},
		{
			Code: "FAA",
			Name: i18n.NewString("Fabrication"),
		},
		{
			Code: "FAB",
			Name: i18n.NewString("Freight equalization"),
		},
		{
			Code: "FAC",
			Name: i18n.NewString("Freight extraordinary handling"),
		},
		{
			Code: "FC",
			Name: i18n.NewString("Freight service"),
		},
		{
			Code: "FH",
			Name: i18n.NewString("Filling/handling"),
		},
		{
			Code: "FI",
			Name: i18n.NewString("Financing"),
		},
		{
			Code: "GAA",
			Name: i18n.NewString("Grinding"),
		},
		{
			Code: "HAA",
			Name: i18n.NewString("Hose"),
		},
		{
			Code: "HD",
			Name: i18n.NewString("Handling"),
		},
		{
			Code: "HH",
			Name: i18n.NewString("Hoisting and hauling"),
		},
		{
			Code: "IAA",
			Name: i18n.NewString("Installation"),
		},
		{
			Code: "IAB",
			Name: i18n.NewString("Installation and warranty"),
		},
		{
			Code: "ID",
			Name: i18n.NewString("Inside delivery"),
		},
		{
			Code: "IF",
			Name: i18n.NewString("Inspection"),
		},
		{
			Code: "IN", // not in EN16931
			Name: i18n.NewString("Insurance"),
		},
		{
			Code: "IR",
			Name: i18n.NewString("Installation and training"),
		},
		{
			Code: "IS",
			Name: i18n.NewString("Invoicing"),
		},
		{
			Code: "KO",
			Name: i18n.NewString("Koshering"),
		},
		{
			Code: "L1",
			Name: i18n.NewString("Carrier count"),
		},
		{
			Code: "LA",
			Name: i18n.NewString("Labelling"),
		},
		{
			Code: "LAA",
			Name: i18n.NewString("Labour"),
		},
		{
			Code: "LAB",
			Name: i18n.NewString("Repair and return"),
		},
		{
			Code: "LF",
			Name: i18n.NewString("Legalisation"),
		},
		{
			Code: "MAE",
			Name: i18n.NewString("Mounting"),
		},
		{
			Code: "MI",
			Name: i18n.NewString("Mail invoice"),
		},
		{
			Code: "ML",
			Name: i18n.NewString("Mail invoice to each location"),
		},
		{
			Code: "NAA",
			Name: i18n.NewString("Non-returnable containers"),
		},
		{
			Code: "OA",
			Name: i18n.NewString("Outside cable connectors"),
		},
		{
			Code: "PA",
			Name: i18n.NewString("Invoice with shipment"),
		},
		{
			Code: "PAA",
			Name: i18n.NewString("Phosphatizing (steel treatment)"),
		},
		{
			Code: "PC",
			Name: i18n.NewString("Packing"),
		},
		{
			Code: "PL",
			Name: i18n.NewString("Palletizing"),
		},
		{
			Code: "PRV",
			Name: i18n.NewString("Price variation"),
		},
		{
			Code: "RAB",
			Name: i18n.NewString("Repacking"),
		},
		{
			Code: "RAC",
			Name: i18n.NewString("Repair"),
		},
		{
			Code: "RAD",
			Name: i18n.NewString("Returnable container"),
		},
		{
			Code: "RAF",
			Name: i18n.NewString("Restocking"),
		},
		{
			Code: "RE",
			Name: i18n.NewString("Re-delivery"),
		},
		{
			Code: "RF",
			Name: i18n.NewString("Refurbishing"),
		},
		{
			Code: "RH",
			Name: i18n.NewString("Rail wagon hire"),
		},
		{
			Code: "RV",
			Name: i18n.NewString("Loading"),
		},
		{
			Code: "SA",
			Name: i18n.NewString("Salvaging"),
		},
		{
			Code: "SAA",
			Name: i18n.NewString("Shipping and handling"),
		},
		{
			Code: "SAD",
			Name: i18n.NewString("Special packaging"),
		},
		{
			Code: "SAE",
			Name: i18n.NewString("Stamping"),
		},
		{
			Code: "SAI",
			Name: i18n.NewString("Consignee unload"),
		},
		{
			Code: "SG",
			Name: i18n.NewString("Shrink-wrap"),
		},
		{
			Code: "SH",
			Name: i18n.NewString("Special handling"),
		},
		{
			Code: "SM",
			Name: i18n.NewString("Special finish"),
		},
		{
			Code: "ST", // not in EN16931
			Name: i18n.NewString("Stamp duties"),
		},
		{
			Code: "SU",
			Name: i18n.NewString("Set-up"),
		},
		{
			Code: "TAB",
			Name: i18n.NewString("Tank renting"),
		},
		{
			Code: "TAC",
			Name: i18n.NewString("Testing"),
		},
		{
			Code: "TT",
			Name: i18n.NewString("Transportation - third party billing"),
		},
		{
			Code: "TV",
			Name: i18n.NewString("Transportation by vendor"),
		},
		{
			Code: "TX", // not in EN16931
			Name: i18n.NewString("Tax"),
		},
		{
			Code: "V1",
			Name: i18n.NewString("Drop yard"),
		},
		{
			Code: "V2",
			Name: i18n.NewString("Drop dock"),
		},
		{
			Code: "WH",
			Name: i18n.NewString("Warehousing"),
		},
		{
			Code: "XAA",
			Name: i18n.NewString("Combine all same day shipment"),
		},
		{
			Code: "YY",
			Name: i18n.NewString("Split pick-up"),
		},
		{
			Code: "ZZZ",
			Name: i18n.NewString("Mutually defined"),
		},
	},
}

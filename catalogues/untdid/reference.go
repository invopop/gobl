package untdid

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

const (
	// ExtKeyReference is used to identify the UNTDID 1153 reference codes
	// qualifiers.
	ExtKeyReference cbc.Key = "untdid-reference"
)

var extReference = &cbc.Definition{
	Key: ExtKeyReference,
	Name: i18n.String{
		i18n.EN: "UNTDID 1153 Reference Code Qualifier",
	},
	Desc: i18n.String{
		i18n.EN: here.Doc(`
			UNTDID 1153 code used to describe the reference code qualifier. This list is based on the
			[EN16931 code list](https://ec.europa.eu/digital-building-blocks/sites/display/DIGITAL/Registry+of+supporting+artefacts+to+implement+EN16931#RegistryofsupportingartefactstoimplementEN16931-Codelists)
			values table which focusses on invoices and payments.
		`),
	},
	Values: []*cbc.Definition{
		{
			Code: "AAA",
			Name: i18n.NewString("Order acknowledgement document identifier"),
		},
		{
			Code: "AAB",
			Name: i18n.NewString("Proforma invoice document identifier"),
		},
		{
			Code: "AAC",
			Name: i18n.NewString("Documentary credit identifier"),
		},
		{
			Code: "AAD",
			Name: i18n.NewString("Contract document addendum identifier"),
		},
		{
			Code: "AAE",
			Name: i18n.NewString("Goods declaration number"),
		},
		{
			Code: "AAF",
			Name: i18n.NewString("Debit card number"),
		},
		{
			Code: "AAG",
			Name: i18n.NewString("Offer number"),
		},
		{
			Code: "AAH",
			Name: i18n.NewString("Bank's batch interbank transaction reference number"),
		},
		{
			Code: "AAI",
			Name: i18n.NewString("Bank's individual interbank transaction reference number"),
		},
		{
			Code: "AAJ",
			Name: i18n.NewString("Delivery order number"),
		},
		{
			Code: "AAK",
			Name: i18n.NewString("Despatch advice number"),
		},
		{
			Code: "AAL",
			Name: i18n.NewString("Drawing number"),
		},
		{
			Code: "AAM",
			Name: i18n.NewString("Waybill number"),
		},
		{
			Code: "AAN",
			Name: i18n.NewString("Delivery schedule number"),
		},
		{
			Code: "AAO",
			Name: i18n.NewString("Consignment identifier, consignee assigned"),
		},
		{
			Code: "AAP",
			Name: i18n.NewString("Partial shipment identifier"),
		},
		{
			Code: "AAQ",
			Name: i18n.NewString("Transport equipment identifier"),
		},
		{
			Code: "AAR",
			Name: i18n.NewString("Municipality assigned business registry number"),
		},
		{
			Code: "AAS",
			Name: i18n.NewString("Transport contract document identifier"),
		},
		{
			Code: "AAT",
			Name: i18n.NewString("Master label number"),
		},
		{
			Code: "AAU",
			Name: i18n.NewString("Despatch note document identifier"),
		},
		{
			Code: "AAV",
			Name: i18n.NewString("Enquiry number"),
		},
		{
			Code: "AAW",
			Name: i18n.NewString("Docket number"),
		},
		{
			Code: "AAX",
			Name: i18n.NewString("Civil action number"),
		},
		{
			Code: "AAY",
			Name: i18n.NewString("Carrier's agent reference number"),
		},
		{
			Code: "AAZ",
			Name: i18n.NewString("Standard Carrier Alpha Code (SCAC) number"),
		},
		{
			Code: "ABA",
			Name: i18n.NewString("Customs valuation decision number"),
		},
		{
			Code: "ABB",
			Name: i18n.NewString("End use authorization number"),
		},
		{
			Code: "ABC",
			Name: i18n.NewString("Anti-dumping case number"),
		},
		{
			Code: "ABD",
			Name: i18n.NewString("Customs tariff number"),
		},
		{
			Code: "ABE",
			Name: i18n.NewString("Declarant's reference number"),
		},
		{
			Code: "ABF",
			Name: i18n.NewString("Repair estimate number"),
		},
		{
			Code: "ABG",
			Name: i18n.NewString("Customs decision request number"),
		},
		{
			Code: "ABH",
			Name: i18n.NewString("Sub-house bill of lading number"),
		},
		{
			Code: "ABI",
			Name: i18n.NewString("Tax payment identifier"),
		},
		{
			Code: "ABJ",
			Name: i18n.NewString("Quota number"),
		},
		{
			Code: "ABK",
			Name: i18n.NewString("Transit (onward carriage) guarantee (bond) number"),
		},
		{
			Code: "ABL",
			Name: i18n.NewString("Customs guarantee number"),
		},
		{
			Code: "ABM",
			Name: i18n.NewString("Replacing part number"),
		},
		{
			Code: "ABN",
			Name: i18n.NewString("Seller's catalogue number"),
		},
		{
			Code: "ABO",
			Name: i18n.NewString("Originator's reference"),
		},
		{
			Code: "ABP",
			Name: i18n.NewString("Declarant's Customs identity number"),
		},
		{
			Code: "ABQ",
			Name: i18n.NewString("Importer reference number"),
		},
		{
			Code: "ABR",
			Name: i18n.NewString("Export clearance instruction reference number"),
		},
		{
			Code: "ABS",
			Name: i18n.NewString("Import clearance instruction reference number"),
		},
		{
			Code: "ABT",
			Name: i18n.NewString("Goods declaration document identifier, Customs"),
		},
		{
			Code: "ABU",
			Name: i18n.NewString("Article number"),
		},
		{
			Code: "ABV",
			Name: i18n.NewString("Intra-plant routing"),
		},
		{
			Code: "ABW",
			Name: i18n.NewString("Stock keeping unit number"),
		},
		{
			Code: "ABX",
			Name: i18n.NewString("Text Element Identifier deletion reference"),
		},
		{
			Code: "ABY",
			Name: i18n.NewString("Allotment identification (Air)"),
		},
		{
			Code: "ABZ",
			Name: i18n.NewString("Vehicle licence number"),
		},
		{
			Code: "AC",
			Name: i18n.NewString("Air cargo transfer manifest"),
		},
		{
			Code: "ACA",
			Name: i18n.NewString("Cargo acceptance order reference number"),
		},
		{
			Code: "ACB",
			Name: i18n.NewString("US government agency number"),
		},
		{
			Code: "ACC",
			Name: i18n.NewString("Shipping unit identification"),
		},
		{
			Code: "ACD",
			Name: i18n.NewString("Additional reference number"),
		},
		{
			Code: "ACE",
			Name: i18n.NewString("Related document number"),
		},
		{
			Code: "ACF",
			Name: i18n.NewString("Addressee reference"),
		},
		{
			Code: "ACG",
			Name: i18n.NewString("ATA carnet number"),
		},
		{
			Code: "ACH",
			Name: i18n.NewString("Packaging unit identification"),
		},
		{
			Code: "ACI",
			Name: i18n.NewString("Outerpackaging unit identification"),
		},
		{
			Code: "ACJ",
			Name: i18n.NewString("Customer material specification number"),
		},
		{
			Code: "ACK",
			Name: i18n.NewString("Bank reference"),
		},
		{
			Code: "ACL",
			Name: i18n.NewString("Principal reference number"),
		},
		{
			Code: "ACN",
			Name: i18n.NewString("Collection advice document identifier"),
		},
		{
			Code: "ACO",
			Name: i18n.NewString("Iron charge number"),
		},
		{
			Code: "ACP",
			Name: i18n.NewString("Hot roll number"),
		},
		{
			Code: "ACQ",
			Name: i18n.NewString("Cold roll number"),
		},
		{
			Code: "ACR",
			Name: i18n.NewString("Railway wagon number"),
		},
		{
			Code: "ACT",
			Name: i18n.NewString("Unique claims reference number of the sender"),
		},
		{
			Code: "ACU",
			Name: i18n.NewString("Loss/event number"),
		},
		{
			Code: "ACV",
			Name: i18n.NewString("Estimate order reference number"),
		},
		{
			Code: "ACW",
			Name: i18n.NewString("Reference number to previous message"),
		},
		{
			Code: "ACX",
			Name: i18n.NewString("Banker's acceptance"),
		},
		{
			Code: "ACY",
			Name: i18n.NewString("Duty memo number"),
		},
		{
			Code: "ACZ",
			Name: i18n.NewString("Equipment transport charge number"),
		},
		{
			Code: "ADA",
			Name: i18n.NewString("Buyer's item number"),
		},
		{
			Code: "ADB",
			Name: i18n.NewString("Matured certificate of deposit"),
		},
		{
			Code: "ADC",
			Name: i18n.NewString("Loan"),
		},
		{
			Code: "ADD",
			Name: i18n.NewString("Analysis number/test number"),
		},
		{
			Code: "ADE",
			Name: i18n.NewString("Account number"),
		},
		{
			Code: "ADF",
			Name: i18n.NewString("Treaty number"),
		},
		{
			Code: "ADG",
			Name: i18n.NewString("Catastrophe number"),
		},
		{
			Code: "ADI",
			Name: i18n.NewString("Bureau signing (statement reference)"),
		},
		{
			Code: "ADJ",
			Name: i18n.NewString("Company / syndicate reference 1"),
		},
		{
			Code: "ADK",
			Name: i18n.NewString("Company / syndicate reference 2"),
		},
		{
			Code: "ADL",
			Name: i18n.NewString("Ordering customer consignment reference number"),
		},
		{
			Code: "ADM",
			Name: i18n.NewString("Shipowner's authorization number"),
		},
		{
			Code: "ADN",
			Name: i18n.NewString("Inland transport order number"),
		},
		{
			Code: "ADO",
			Name: i18n.NewString("Container work order reference number"),
		},
		{
			Code: "ADP",
			Name: i18n.NewString("Statement number"),
		},
		{
			Code: "ADQ",
			Name: i18n.NewString("Unique market reference"),
		},
		{
			Code: "ADT",
			Name: i18n.NewString("Group accounting"),
		},
		{
			Code: "ADU",
			Name: i18n.NewString("Broker reference 1"),
		},
		{
			Code: "ADV",
			Name: i18n.NewString("Broker reference 2"),
		},
		{
			Code: "ADW",
			Name: i18n.NewString("Lloyd's claims office reference"),
		},
		{
			Code: "ADX",
			Name: i18n.NewString("Secure delivery terms and conditions agreement reference"),
		},
		{
			Code: "ADY",
			Name: i18n.NewString("Report number"),
		},
		{
			Code: "ADZ",
			Name: i18n.NewString("Trader account number"),
		},
		{
			Code: "AE",
			Name: i18n.NewString("Authorization for expense (AFE) number"),
		},
		{
			Code: "AEA",
			Name: i18n.NewString("Government agency reference number"),
		},
		{
			Code: "AEB",
			Name: i18n.NewString("Assembly number"),
		},
		{
			Code: "AEC",
			Name: i18n.NewString("Symbol number"),
		},
		{
			Code: "AED",
			Name: i18n.NewString("Commodity number"),
		},
		{
			Code: "AEE",
			Name: i18n.NewString("Eur 1 certificate number"),
		},
		{
			Code: "AEF",
			Name: i18n.NewString("Customer process specification number"),
		},
		{
			Code: "AEG",
			Name: i18n.NewString("Customer specification number"),
		},
		{
			Code: "AEH",
			Name: i18n.NewString("Applicable instructions or standards"),
		},
		{
			Code: "AEI",
			Name: i18n.NewString("Registration number of previous Customs declaration"),
		},
		{
			Code: "AEJ",
			Name: i18n.NewString("Post-entry reference"),
		},
		{
			Code: "AEK",
			Name: i18n.NewString("Payment order number"),
		},
		{
			Code: "AEL",
			Name: i18n.NewString("Delivery number (transport)"),
		},
		{
			Code: "AEM",
			Name: i18n.NewString("Transport route"),
		},
		{
			Code: "AEN",
			Name: i18n.NewString("Customer's unit inventory number"),
		},
		{
			Code: "AEO",
			Name: i18n.NewString("Product reservation number"),
		},
		{
			Code: "AEP",
			Name: i18n.NewString("Project number"),
		},
		{
			Code: "AEQ",
			Name: i18n.NewString("Drawing list number"),
		},
		{
			Code: "AER",
			Name: i18n.NewString("Project specification number"),
		},
		{
			Code: "AES",
			Name: i18n.NewString("Primary reference"),
		},
		{
			Code: "AET",
			Name: i18n.NewString("Request for cancellation number"),
		},
		{
			Code: "AEU",
			Name: i18n.NewString("Supplier's control number"),
		},
		{
			Code: "AEV",
			Name: i18n.NewString("Shipping note number"),
		},
		{
			Code: "AEW",
			Name: i18n.NewString("Empty container bill number"),
		},
		{
			Code: "AEX",
			Name: i18n.NewString("Non-negotiable maritime transport document number"),
		},
		{
			Code: "AEY",
			Name: i18n.NewString("Substitute air waybill number"),
		},
		{
			Code: "AEZ",
			Name: i18n.NewString("Despatch note (post parcels) number"),
		},
		{
			Code: "AF",
			Name: i18n.NewString("Airlines flight identification number"),
		},
		{
			Code: "AFA",
			Name: i18n.NewString("Through bill of lading number"),
		},
		{
			Code: "AFB",
			Name: i18n.NewString("Cargo manifest number"),
		},
		{
			Code: "AFC",
			Name: i18n.NewString("Bordereau number"),
		},
		{
			Code: "AFD",
			Name: i18n.NewString("Customs item number"),
		},
		{
			Code: "AFE",
			Name: i18n.NewString("Export Control Commodity number (ECCN)"),
		},
		{
			Code: "AFF",
			Name: i18n.NewString("Marking/label reference"),
		},
		{
			Code: "AFG",
			Name: i18n.NewString("Tariff number"),
		},
		{
			Code: "AFH",
			Name: i18n.NewString("Replenishment purchase order number"),
		},
		{
			Code: "AFI",
			Name: i18n.NewString("Immediate transportation no. for in bond movement"),
		},
		{
			Code: "AFJ",
			Name: i18n.NewString("Transportation exportation no. for in bond movement"),
		},
		{
			Code: "AFK",
			Name: i18n.NewString("Immediate exportation no. for in bond movement"),
		},
		{
			Code: "AFL",
			Name: i18n.NewString("Associated invoices"),
		},
		{
			Code: "AFM",
			Name: i18n.NewString("Secondary Customs reference"),
		},
		{
			Code: "AFN",
			Name: i18n.NewString("Account party's reference"),
		},
		{
			Code: "AFO",
			Name: i18n.NewString("Beneficiary's reference"),
		},
		{
			Code: "AFP",
			Name: i18n.NewString("Second beneficiary's reference"),
		},
		{
			Code: "AFQ",
			Name: i18n.NewString("Applicant's bank reference"),
		},
		{
			Code: "AFR",
			Name: i18n.NewString("Issuing bank's reference"),
		},
		{
			Code: "AFS",
			Name: i18n.NewString("Beneficiary's bank reference"),
		},
		{
			Code: "AFT",
			Name: i18n.NewString("Direct payment valuation number"),
		},
		{
			Code: "AFU",
			Name: i18n.NewString("Direct payment valuation request number"),
		},
		{
			Code: "AFV",
			Name: i18n.NewString("Quantity valuation number"),
		},
		{
			Code: "AFW",
			Name: i18n.NewString("Quantity valuation request number"),
		},
		{
			Code: "AFX",
			Name: i18n.NewString("Bill of quantities number"),
		},
		{
			Code: "AFY",
			Name: i18n.NewString("Payment valuation number"),
		},
		{
			Code: "AFZ",
			Name: i18n.NewString("Situation number"),
		},
		{
			Code: "AGA",
			Name: i18n.NewString("Agreement to pay number"),
		},
		{
			Code: "AGB",
			Name: i18n.NewString("Contract party reference number"),
		},
		{
			Code: "AGC",
			Name: i18n.NewString("Account party's bank reference"),
		},
		{
			Code: "AGD",
			Name: i18n.NewString("Agent's bank reference"),
		},
		{
			Code: "AGE",
			Name: i18n.NewString("Agent's reference"),
		},
		{
			Code: "AGF",
			Name: i18n.NewString("Applicant's reference"),
		},
		{
			Code: "AGG",
			Name: i18n.NewString("Dispute number"),
		},
		{
			Code: "AGH",
			Name: i18n.NewString("Credit rating agency's reference number"),
		},
		{
			Code: "AGI",
			Name: i18n.NewString("Request number"),
		},
		{
			Code: "AGJ",
			Name: i18n.NewString("Single transaction sequence number"),
		},
		{
			Code: "AGK",
			Name: i18n.NewString("Application reference number"),
		},
		{
			Code: "AGL",
			Name: i18n.NewString("Delivery verification certificate"),
		},
		{
			Code: "AGM",
			Name: i18n.NewString("Number of temporary importation document"),
		},
		{
			Code: "AGN",
			Name: i18n.NewString("Reference number quoted on statement"),
		},
		{
			Code: "AGO",
			Name: i18n.NewString("Sender's reference to the original message"),
		},
		{
			Code: "AGP",
			Name: i18n.NewString("Company issued equipment ID"),
		},
		{
			Code: "AGQ",
			Name: i18n.NewString("Domestic flight number"),
		},
		{
			Code: "AGR",
			Name: i18n.NewString("International flight number"),
		},
		{
			Code: "AGS",
			Name: i18n.NewString("Employer identification number of service bureau"),
		},
		{
			Code: "AGT",
			Name: i18n.NewString("Service group identification number"),
		},
		{
			Code: "AGU",
			Name: i18n.NewString("Member number"),
		},
		{
			Code: "AGV",
			Name: i18n.NewString("Previous member number"),
		},
		{
			Code: "AGW",
			Name: i18n.NewString("Scheme/plan number"),
		},
		{
			Code: "AGX",
			Name: i18n.NewString("Previous scheme/plan number"),
		},
		{
			Code: "AGY",
			Name: i18n.NewString("Receiving party's member identification"),
		},
		{
			Code: "AGZ",
			Name: i18n.NewString("Payroll number"),
		},
		{
			Code: "AHA",
			Name: i18n.NewString("Packaging specification number"),
		},
		{
			Code: "AHB",
			Name: i18n.NewString("Authority issued equipment identification"),
		},
		{
			Code: "AHC",
			Name: i18n.NewString("Training flight number"),
		},
		{
			Code: "AHD",
			Name: i18n.NewString("Fund code number"),
		},
		{
			Code: "AHE",
			Name: i18n.NewString("Signal code number"),
		},
		{
			Code: "AHF",
			Name: i18n.NewString("Major force program number"),
		},
		{
			Code: "AHG",
			Name: i18n.NewString("Nomination number"),
		},
		{
			Code: "AHH",
			Name: i18n.NewString("Laboratory registration number"),
		},
		{
			Code: "AHI",
			Name: i18n.NewString("Transport contract reference number"),
		},
		{
			Code: "AHJ",
			Name: i18n.NewString("Payee's reference number"),
		},
		{
			Code: "AHK",
			Name: i18n.NewString("Payer's reference number"),
		},
		{
			Code: "AHL",
			Name: i18n.NewString("Creditor's reference number"),
		},
		{
			Code: "AHM",
			Name: i18n.NewString("Debtor's reference number"),
		},
		{
			Code: "AHN",
			Name: i18n.NewString("Joint venture reference number"),
		},
		{
			Code: "AHO",
			Name: i18n.NewString("Chamber of Commerce registration number"),
		},
		{
			Code: "AHP",
			Name: i18n.NewString("Tax registration number"),
		},
		{
			Code: "AHQ",
			Name: i18n.NewString("Wool identification number"),
		},
		{
			Code: "AHR",
			Name: i18n.NewString("Wool tax reference number"),
		},
		{
			Code: "AHS",
			Name: i18n.NewString("Meat processing establishment registration number"),
		},
		{
			Code: "AHT",
			Name: i18n.NewString("Quarantine/treatment status reference number"),
		},
		{
			Code: "AHU",
			Name: i18n.NewString("Request for quote number"),
		},
		{
			Code: "AHV",
			Name: i18n.NewString("Manual processing authority number"),
		},
		{
			Code: "AHX",
			Name: i18n.NewString("Rate note number"),
		},
		{
			Code: "AHY",
			Name: i18n.NewString("Freight Forwarder number"),
		},
		{
			Code: "AHZ",
			Name: i18n.NewString("Customs release code"),
		},
		{
			Code: "AIA",
			Name: i18n.NewString("Compliance code number"),
		},
		{
			Code: "AIB",
			Name: i18n.NewString("Department of transportation bond number"),
		},
		{
			Code: "AIC",
			Name: i18n.NewString("Export establishment number"),
		},
		{
			Code: "AID",
			Name: i18n.NewString("Certificate of conformity"),
		},
		{
			Code: "AIE",
			Name: i18n.NewString("Ministerial certificate of homologation"),
		},
		{
			Code: "AIF",
			Name: i18n.NewString("Previous delivery instruction number"),
		},
		{
			Code: "AIG",
			Name: i18n.NewString("Passport number"),
		},
		{
			Code: "AIH",
			Name: i18n.NewString("Common transaction reference number"),
		},
		{
			Code: "AII",
			Name: i18n.NewString("Bank's common transaction reference number"),
		},
		{
			Code: "AIJ",
			Name: i18n.NewString("Customer's individual transaction reference number"),
		},
		{
			Code: "AIK",
			Name: i18n.NewString("Bank's individual transaction reference number"),
		},
		{
			Code: "AIL",
			Name: i18n.NewString("Customer's common transaction reference number"),
		},
		{
			Code: "AIM",
			Name: i18n.NewString("Individual transaction reference number"),
		},
		{
			Code: "AIN",
			Name: i18n.NewString("Product sourcing agreement number"),
		},
		{
			Code: "AIO",
			Name: i18n.NewString("Customs transhipment number"),
		},
		{
			Code: "AIP",
			Name: i18n.NewString("Customs preference inquiry number"),
		},
		{
			Code: "AIQ",
			Name: i18n.NewString("Packing plant number"),
		},
		{
			Code: "AIR",
			Name: i18n.NewString("Original certificate number"),
		},
		{
			Code: "AIS",
			Name: i18n.NewString("Processing plant number"),
		},
		{
			Code: "AIT",
			Name: i18n.NewString("Slaughter plant number"),
		},
		{
			Code: "AIU",
			Name: i18n.NewString("Charge card account number"),
		},
		{
			Code: "AIV",
			Name: i18n.NewString("Event reference number"),
		},
		{
			Code: "AIW",
			Name: i18n.NewString("Transport section reference number"),
		},
		{
			Code: "AIX",
			Name: i18n.NewString("Referred product for mechanical analysis"),
		},
		{
			Code: "AIY",
			Name: i18n.NewString("Referred product for chemical analysis"),
		},
		{
			Code: "AIZ",
			Name: i18n.NewString("Consolidated invoice number"),
		},
		{
			Code: "AJA",
			Name: i18n.NewString("Part reference indicator in a drawing"),
		},
		{
			Code: "AJB",
			Name: i18n.NewString("U.S. Code of Federal Regulations (CFR)"),
		},
		{
			Code: "AJC",
			Name: i18n.NewString("Purchasing activity clause number"),
		},
		{
			Code: "AJD",
			Name: i18n.NewString("U.S. Defense Federal Acquisition Regulation Supplement"),
		},
		{
			Code: "AJE",
			Name: i18n.NewString("Agency clause number"),
		},
		{
			Code: "AJF",
			Name: i18n.NewString("Circular publication number"),
		},
		{
			Code: "AJG",
			Name: i18n.NewString("U.S. Federal Acquisition Regulation"),
		},
		{
			Code: "AJH",
			Name: i18n.NewString("U.S. General Services Administration Regulation"),
		},
		{
			Code: "AJI",
			Name: i18n.NewString("U.S. Federal Information Resources Management Regulation"),
		},
		{
			Code: "AJJ",
			Name: i18n.NewString("Paragraph"),
		},
		{
			Code: "AJK",
			Name: i18n.NewString("Special instructions number"),
		},
		{
			Code: "AJL",
			Name: i18n.NewString("Site specific procedures, terms, and conditions number"),
		},
		{
			Code: "AJM",
			Name: i18n.NewString("Master solicitation procedures, terms, and conditions"),
		},
		{
			Code: "AJN",
			Name: i18n.NewString("U.S. Department of Veterans Affairs Acquisition Regulation"),
		},
		{
			Code: "AJO",
			Name: i18n.NewString("Military Interdepartmental Purchase Request (MIPR) number"),
		},
		{
			Code: "AJP",
			Name: i18n.NewString("Foreign military sales number"),
		},
		{
			Code: "AJQ",
			Name: i18n.NewString("Defense priorities allocation system priority rating"),
		},
		{
			Code: "AJR",
			Name: i18n.NewString("Wage determination number"),
		},
		{
			Code: "AJS",
			Name: i18n.NewString("Agreement number"),
		},
		{
			Code: "AJT",
			Name: i18n.NewString("Standard Industry Classification (SIC) number"),
		},
		{
			Code: "AJU",
			Name: i18n.NewString("End item number"),
		},
		{
			Code: "AJV",
			Name: i18n.NewString("Federal supply schedule item number"),
		},
		{
			Code: "AJW",
			Name: i18n.NewString("Technical document number"),
		},
		{
			Code: "AJX",
			Name: i18n.NewString("Technical order number"),
		},
		{
			Code: "AJY",
			Name: i18n.NewString("Suffix"),
		},
		{
			Code: "AJZ",
			Name: i18n.NewString("Transportation account number"),
		},
		{
			Code: "AKA",
			Name: i18n.NewString("Container disposition order reference number"),
		},
		{
			Code: "AKB",
			Name: i18n.NewString("Container prefix"),
		},
		{
			Code: "AKC",
			Name: i18n.NewString("Transport equipment return reference"),
		},
		{
			Code: "AKD",
			Name: i18n.NewString("Transport equipment survey reference"),
		},
		{
			Code: "AKE",
			Name: i18n.NewString("Transport equipment survey report number"),
		},
		{
			Code: "AKF",
			Name: i18n.NewString("Transport equipment stuffing order"),
		},
		{
			Code: "AKG",
			Name: i18n.NewString("Vehicle Identification Number (VIN)"),
		},
		{
			Code: "AKH",
			Name: i18n.NewString("Government bill of lading"),
		},
		{
			Code: "AKI",
			Name: i18n.NewString("Ordering customer's second reference number"),
		},
		{
			Code: "AKJ",
			Name: i18n.NewString("Direct debit reference"),
		},
		{
			Code: "AKK",
			Name: i18n.NewString("Meter reading at the beginning of the delivery"),
		},
		{
			Code: "AKL",
			Name: i18n.NewString("Meter reading at the end of delivery"),
		},
		{
			Code: "AKM",
			Name: i18n.NewString("Replenishment purchase order range start number"),
		},
		{
			Code: "AKN",
			Name: i18n.NewString("Third bank's reference"),
		},
		{
			Code: "AKO",
			Name: i18n.NewString("Action authorization number"),
		},
		{
			Code: "AKP",
			Name: i18n.NewString("Appropriation number"),
		},
		{
			Code: "AKQ",
			Name: i18n.NewString("Product change authority number"),
		},
		{
			Code: "AKR",
			Name: i18n.NewString("General cargo consignment reference number"),
		},
		{
			Code: "AKS",
			Name: i18n.NewString("Catalogue sequence number"),
		},
		{
			Code: "AKT",
			Name: i18n.NewString("Forwarding order number"),
		},
		{
			Code: "AKU",
			Name: i18n.NewString("Transport equipment survey reference number"),
		},
		{
			Code: "AKV",
			Name: i18n.NewString("Lease contract reference"),
		},
		{
			Code: "AKW",
			Name: i18n.NewString("Transport costs reference number"),
		},
		{
			Code: "AKX",
			Name: i18n.NewString("Transport equipment stripping order"),
		},
		{
			Code: "AKY",
			Name: i18n.NewString("Prior policy number"),
		},
		{
			Code: "AKZ",
			Name: i18n.NewString("Policy number"),
		},
		{
			Code: "ALA",
			Name: i18n.NewString("Procurement budget number"),
		},
		{
			Code: "ALB",
			Name: i18n.NewString("Domestic inventory management code"),
		},
		{
			Code: "ALC",
			Name: i18n.NewString("Customer reference number assigned to previous balance of"),
		},
		{
			Code: "ALD",
			Name: i18n.NewString("Previous credit advice reference number"),
		},
		{
			Code: "ALE",
			Name: i18n.NewString("Reporting form number"),
		},
		{
			Code: "ALF",
			Name: i18n.NewString("Authorization number for exception to dangerous goods"),
		},
		{
			Code: "ALG",
			Name: i18n.NewString("Dangerous goods security number"),
		},
		{
			Code: "ALH",
			Name: i18n.NewString("Dangerous goods transport licence number"),
		},
		{
			Code: "ALI",
			Name: i18n.NewString("Previous rental agreement number"),
		},
		{
			Code: "ALJ",
			Name: i18n.NewString("Next rental agreement reason number"),
		},
		{
			Code: "ALK",
			Name: i18n.NewString("Consignee's invoice number"),
		},
		{
			Code: "ALL",
			Name: i18n.NewString("Message batch number"),
		},
		{
			Code: "ALM",
			Name: i18n.NewString("Previous delivery schedule number"),
		},
		{
			Code: "ALN",
			Name: i18n.NewString("Physical inventory recount reference number"),
		},
		{
			Code: "ALO",
			Name: i18n.NewString("Receiving advice number"),
		},
		{
			Code: "ALP",
			Name: i18n.NewString("Returnable container reference number"),
		},
		{
			Code: "ALQ",
			Name: i18n.NewString("Returns notice number"),
		},
		{
			Code: "ALR",
			Name: i18n.NewString("Sales forecast number"),
		},
		{
			Code: "ALS",
			Name: i18n.NewString("Sales report number"),
		},
		{
			Code: "ALT",
			Name: i18n.NewString("Previous tax control number"),
		},
		{
			Code: "ALU",
			Name: i18n.NewString("AGERD (Aerospace Ground Equipment Requirement Data) number"),
		},
		{
			Code: "ALV",
			Name: i18n.NewString("Registered capital reference"),
		},
		{
			Code: "ALW",
			Name: i18n.NewString("Standard number of inspection document"),
		},
		{
			Code: "ALX",
			Name: i18n.NewString("Model"),
		},
		{
			Code: "ALY",
			Name: i18n.NewString("Financial management reference"),
		},
		{
			Code: "ALZ",
			Name: i18n.NewString("NOTIfication for COLlection number (NOTICOL)"),
		},
		{
			Code: "AMA",
			Name: i18n.NewString("Previous request for metered reading reference number"),
		},
		{
			Code: "AMB",
			Name: i18n.NewString("Next rental agreement number"),
		},
		{
			Code: "AMC",
			Name: i18n.NewString("Reference number of a request for metered reading"),
		},
		{
			Code: "AMD",
			Name: i18n.NewString("Hastening number"),
		},
		{
			Code: "AME",
			Name: i18n.NewString("Repair data request number"),
		},
		{
			Code: "AMF",
			Name: i18n.NewString("Consumption data request number"),
		},
		{
			Code: "AMG",
			Name: i18n.NewString("Profile number"),
		},
		{
			Code: "AMH",
			Name: i18n.NewString("Case number"),
		},
		{
			Code: "AMI",
			Name: i18n.NewString("Government quality assurance and control level Number"),
		},
		{
			Code: "AMJ",
			Name: i18n.NewString("Payment plan reference"),
		},
		{
			Code: "AMK",
			Name: i18n.NewString("Replaced meter unit number"),
		},
		{
			Code: "AML",
			Name: i18n.NewString("Replenishment purchase order range end number"),
		},
		{
			Code: "AMM",
			Name: i18n.NewString("Insurer assigned reference number"),
		},
		{
			Code: "AMN",
			Name: i18n.NewString("Canadian excise entry number"),
		},
		{
			Code: "AMO",
			Name: i18n.NewString("Premium rate table"),
		},
		{
			Code: "AMP",
			Name: i18n.NewString("Advise through bank's reference"),
		},
		{
			Code: "AMQ",
			Name: i18n.NewString("US, Department of Transportation bond surety code"),
		},
		{
			Code: "AMR",
			Name: i18n.NewString("US, Food and Drug Administration establishment indicator"),
		},
		{
			Code: "AMS",
			Name: i18n.NewString("US, Federal Communications Commission (FCC) import"),
		},
		{
			Code: "AMT",
			Name: i18n.NewString("Goods and Services Tax identification number"),
		},
		{
			Code: "AMU",
			Name: i18n.NewString("Integrated logistic support cross reference number"),
		},
		{
			Code: "AMV",
			Name: i18n.NewString("Department number"),
		},
		{
			Code: "AMW",
			Name: i18n.NewString("Buyer's catalogue number"),
		},
		{
			Code: "AMX",
			Name: i18n.NewString("Financial settlement party's reference number"),
		},
		{
			Code: "AMY",
			Name: i18n.NewString("Standard's version number"),
		},
		{
			Code: "AMZ",
			Name: i18n.NewString("Pipeline number"),
		},
		{
			Code: "ANA",
			Name: i18n.NewString("Account servicing bank's reference number"),
		},
		{
			Code: "ANB",
			Name: i18n.NewString("Completed units payment request reference"),
		},
		{
			Code: "ANC",
			Name: i18n.NewString("Payment in advance request reference"),
		},
		{
			Code: "AND",
			Name: i18n.NewString("Parent file"),
		},
		{
			Code: "ANE",
			Name: i18n.NewString("Sub file"),
		},
		{
			Code: "ANF",
			Name: i18n.NewString("CAD file layer convention"),
		},
		{
			Code: "ANG",
			Name: i18n.NewString("Technical regulation"),
		},
		{
			Code: "ANH",
			Name: i18n.NewString("Plot file"),
		},
		{
			Code: "ANI",
			Name: i18n.NewString("File conversion journal"),
		},
		{
			Code: "ANJ",
			Name: i18n.NewString("Authorization number"),
		},
		{
			Code: "ANK",
			Name: i18n.NewString("Reference number assigned by third party"),
		},
		{
			Code: "ANL",
			Name: i18n.NewString("Deposit reference number"),
		},
		{
			Code: "ANM",
			Name: i18n.NewString("Named bank's reference"),
		},
		{
			Code: "ANN",
			Name: i18n.NewString("Drawee's reference"),
		},
		{
			Code: "ANO",
			Name: i18n.NewString("Case of need party's reference"),
		},
		{
			Code: "ANP",
			Name: i18n.NewString("Collecting bank's reference"),
		},
		{
			Code: "ANQ",
			Name: i18n.NewString("Remitting bank's reference"),
		},
		{
			Code: "ANR",
			Name: i18n.NewString("Principal's bank reference"),
		},
		{
			Code: "ANS",
			Name: i18n.NewString("Presenting bank's reference"),
		},
		{
			Code: "ANT",
			Name: i18n.NewString("Consignee's reference"),
		},
		{
			Code: "ANU",
			Name: i18n.NewString("Financial transaction reference number"),
		},
		{
			Code: "ANV",
			Name: i18n.NewString("Credit reference number"),
		},
		{
			Code: "ANW",
			Name: i18n.NewString("Receiving bank's authorization number"),
		},
		{
			Code: "ANX",
			Name: i18n.NewString("Clearing reference"),
		},
		{
			Code: "ANY",
			Name: i18n.NewString("Sending bank's reference number"),
		},
		{
			Code: "AOA",
			Name: i18n.NewString("Documentary payment reference"),
		},
		{
			Code: "AOD",
			Name: i18n.NewString("Accounting file reference"),
		},
		{
			Code: "AOE",
			Name: i18n.NewString("Sender's file reference number"),
		},
		{
			Code: "AOF",
			Name: i18n.NewString("Receiver's file reference number"),
		},
		{
			Code: "AOG",
			Name: i18n.NewString("Source document internal reference"),
		},
		{
			Code: "AOH",
			Name: i18n.NewString("Principal's reference"),
		},
		{
			Code: "AOI",
			Name: i18n.NewString("Debit reference number"),
		},
		{
			Code: "AOJ",
			Name: i18n.NewString("Calendar"),
		},
		{
			Code: "AOK",
			Name: i18n.NewString("Work shift"),
		},
		{
			Code: "AOL",
			Name: i18n.NewString("Work breakdown structure"),
		},
		{
			Code: "AOM",
			Name: i18n.NewString("Organisation breakdown structure"),
		},
		{
			Code: "AON",
			Name: i18n.NewString("Work task charge number"),
		},
		{
			Code: "AOO",
			Name: i18n.NewString("Functional work group"),
		},
		{
			Code: "AOP",
			Name: i18n.NewString("Work team"),
		},
		{
			Code: "AOQ",
			Name: i18n.NewString("Department"),
		},
		{
			Code: "AOR",
			Name: i18n.NewString("Statement of work"),
		},
		{
			Code: "AOS",
			Name: i18n.NewString("Work package"),
		},
		{
			Code: "AOT",
			Name: i18n.NewString("Planning package"),
		},
		{
			Code: "AOU",
			Name: i18n.NewString("Cost account"),
		},
		{
			Code: "AOV",
			Name: i18n.NewString("Work order"),
		},
		{
			Code: "AOW",
			Name: i18n.NewString("Transportation Control Number (TCN)"),
		},
		{
			Code: "AOX",
			Name: i18n.NewString("Constraint notation"),
		},
		{
			Code: "AOY",
			Name: i18n.NewString("ETERMS reference"),
		},
		{
			Code: "AOZ",
			Name: i18n.NewString("Implementation version number"),
		},
		{
			Code: "AP",
			Name: i18n.NewString("Accounts receivable number"),
		},
		{
			Code: "APA",
			Name: i18n.NewString("Incorporated legal reference"),
		},
		{
			Code: "APB",
			Name: i18n.NewString("Payment instalment reference number"),
		},
		{
			Code: "APC",
			Name: i18n.NewString("Equipment owner reference number"),
		},
		{
			Code: "APD",
			Name: i18n.NewString("Cedent's claim number"),
		},
		{
			Code: "APE",
			Name: i18n.NewString("Reinsurer's claim number"),
		},
		{
			Code: "APF",
			Name: i18n.NewString("Price/sales catalogue response reference number"),
		},
		{
			Code: "APG",
			Name: i18n.NewString("General purpose message reference number"),
		},
		{
			Code: "APH",
			Name: i18n.NewString("Invoicing data sheet reference number"),
		},
		{
			Code: "API",
			Name: i18n.NewString("Inventory report reference number"),
		},
		{
			Code: "APJ",
			Name: i18n.NewString("Ceiling formula reference number"),
		},
		{
			Code: "APK",
			Name: i18n.NewString("Price variation formula reference number"),
		},
		{
			Code: "APL",
			Name: i18n.NewString("Reference to account servicing bank's message"),
		},
		{
			Code: "APM",
			Name: i18n.NewString("Party sequence number"),
		},
		{
			Code: "APN",
			Name: i18n.NewString("Purchaser's request reference"),
		},
		{
			Code: "APO",
			Name: i18n.NewString("Contractor request reference"),
		},
		{
			Code: "APP",
			Name: i18n.NewString("Accident reference number"),
		},
		{
			Code: "APQ",
			Name: i18n.NewString("Commercial account summary reference number"),
		},
		{
			Code: "APR",
			Name: i18n.NewString("Contract breakdown reference"),
		},
		{
			Code: "APS",
			Name: i18n.NewString("Contractor registration number"),
		},
		{
			Code: "APT",
			Name: i18n.NewString("Applicable coefficient identification number"),
		},
		{
			Code: "APU",
			Name: i18n.NewString("Special budget account number"),
		},
		{
			Code: "APV",
			Name: i18n.NewString("Authorisation for repair reference"),
		},
		{
			Code: "APW",
			Name: i18n.NewString("Manufacturer defined repair rates reference"),
		},
		{
			Code: "APX",
			Name: i18n.NewString("Original submitter log number"),
		},
		{
			Code: "APY",
			Name: i18n.NewString("Original submitter, parent Data Maintenance Request (DMR)"),
		},
		{
			Code: "APZ",
			Name: i18n.NewString("Original submitter, child Data Maintenance Request (DMR)"),
		},
		{
			Code: "AQA",
			Name: i18n.NewString("Entry point assessment log number"),
		},
		{
			Code: "AQB",
			Name: i18n.NewString("Entry point assessment log number, parent DMR"),
		},
		{
			Code: "AQC",
			Name: i18n.NewString("Entry point assessment log number, child DMR"),
		},
		{
			Code: "AQD",
			Name: i18n.NewString("Data structure tag"),
		},
		{
			Code: "AQE",
			Name: i18n.NewString("Central secretariat log number"),
		},
		{
			Code: "AQF",
			Name: i18n.NewString("Central secretariat log number, parent Data Maintenance"),
		},
		{
			Code: "AQG",
			Name: i18n.NewString("Central secretariat log number, child Data Maintenance"),
		},
		{
			Code: "AQH",
			Name: i18n.NewString("International assessment log number"),
		},
		{
			Code: "AQI",
			Name: i18n.NewString("International assessment log number, parent Data"),
		},
		{
			Code: "AQJ",
			Name: i18n.NewString("International assessment log number, child Data Maintenance"),
		},
		{
			Code: "AQK",
			Name: i18n.NewString("Status report number"),
		},
		{
			Code: "AQL",
			Name: i18n.NewString("Message design group number"),
		},
		{
			Code: "AQM",
			Name: i18n.NewString("US Customs Service (USCS) entry code"),
		},
		{
			Code: "AQN",
			Name: i18n.NewString("Beginning job sequence number"),
		},
		{
			Code: "AQO",
			Name: i18n.NewString("Sender's clause number"),
		},
		{
			Code: "AQP",
			Name: i18n.NewString("Dun and Bradstreet Canada's 8 digit Standard Industrial"),
		},
		{
			Code: "AQQ",
			Name: i18n.NewString("Activite Principale Exercee (APE) identifier"),
		},
		{
			Code: "AQR",
			Name: i18n.NewString("Dun and Bradstreet US 8 digit Standard Industrial"),
		},
		{
			Code: "AQS",
			Name: i18n.NewString("Nomenclature Activity Classification Economy (NACE)"),
		},
		{
			Code: "AQT",
			Name: i18n.NewString("Norme Activite Francaise (NAF) identifier"),
		},
		{
			Code: "AQU",
			Name: i18n.NewString("Registered contractor activity type"),
		},
		{
			Code: "AQV",
			Name: i18n.NewString("Statistic Bundes Amt (SBA) identifier"),
		},
		{
			Code: "AQW",
			Name: i18n.NewString("State or province assigned entity identification"),
		},
		{
			Code: "AQX",
			Name: i18n.NewString("Institute of Security and Future Market Development (ISFMD)"),
		},
		{
			Code: "AQY",
			Name: i18n.NewString("File identification number"),
		},
		{
			Code: "AQZ",
			Name: i18n.NewString("Bankruptcy procedure number"),
		},
		{
			Code: "ARA",
			Name: i18n.NewString("National government business identification number"),
		},
		{
			Code: "ARB",
			Name: i18n.NewString("Prior Data Universal Number System (DUNS) number"),
		},
		{
			Code: "ARC",
			Name: i18n.NewString("Companies Registry Office (CRO) number"),
		},
		{
			Code: "ARD",
			Name: i18n.NewString("Costa Rican judicial number"),
		},
		{
			Code: "ARE",
			Name: i18n.NewString("Numero de Identificacion Tributaria (NIT)"),
		},
		{
			Code: "ARF",
			Name: i18n.NewString("Patron number"),
		},
		{
			Code: "ARG",
			Name: i18n.NewString("Registro Informacion Fiscal (RIF) number"),
		},
		{
			Code: "ARH",
			Name: i18n.NewString("Registro Unico de Contribuyente (RUC) number"),
		},
		{
			Code: "ARI",
			Name: i18n.NewString("Tokyo SHOKO Research (TSR) business identifier"),
		},
		{
			Code: "ARJ",
			Name: i18n.NewString("Personal identity card number"),
		},
		{
			Code: "ARK",
			Name: i18n.NewString("Systeme Informatique pour le Repertoire des ENtreprises"),
		},
		{
			Code: "ARL",
			Name: i18n.NewString("Systeme Informatique pour le Repertoire des ETablissements"),
		},
		{
			Code: "ARM",
			Name: i18n.NewString("Publication issue number"),
		},
		{
			Code: "ARN",
			Name: i18n.NewString("Original filing number"),
		},
		{
			Code: "ARO",
			Name: i18n.NewString("Document page identifier"),
		},
		{
			Code: "ARP",
			Name: i18n.NewString("Public filing registration number"),
		},
		{
			Code: "ARQ",
			Name: i18n.NewString("Regiristo Federal de Contribuyentes"),
		},
		{
			Code: "ARR",
			Name: i18n.NewString("Social security number"),
		},
		{
			Code: "ARS",
			Name: i18n.NewString("Document volume number"),
		},
		{
			Code: "ART",
			Name: i18n.NewString("Book number"),
		},
		{
			Code: "ARU",
			Name: i18n.NewString("Stock exchange company identifier"),
		},
		{
			Code: "ARV",
			Name: i18n.NewString("Imputation account"),
		},
		{
			Code: "ARW",
			Name: i18n.NewString("Financial phase reference"),
		},
		{
			Code: "ARX",
			Name: i18n.NewString("Technical phase reference"),
		},
		{
			Code: "ARY",
			Name: i18n.NewString("Prior contractor registration number"),
		},
		{
			Code: "ARZ",
			Name: i18n.NewString("Stock adjustment number"),
		},
		{
			Code: "ASA",
			Name: i18n.NewString("Dispensation reference"),
		},
		{
			Code: "ASB",
			Name: i18n.NewString("Investment reference number"),
		},
		{
			Code: "ASC",
			Name: i18n.NewString("Assuming company"),
		},
		{
			Code: "ASD",
			Name: i18n.NewString("Budget chapter"),
		},
		{
			Code: "ASE",
			Name: i18n.NewString("Duty free products security number"),
		},
		{
			Code: "ASF",
			Name: i18n.NewString("Duty free products receipt authorisation number"),
		},
		{
			Code: "ASG",
			Name: i18n.NewString("Party information message reference"),
		},
		{
			Code: "ASH",
			Name: i18n.NewString("Formal statement reference"),
		},
		{
			Code: "ASI",
			Name: i18n.NewString("Proof of delivery reference number"),
		},
		{
			Code: "ASJ",
			Name: i18n.NewString("Supplier's credit claim reference number"),
		},
		{
			Code: "ASK",
			Name: i18n.NewString("Picture of actual product"),
		},
		{
			Code: "ASL",
			Name: i18n.NewString("Picture of a generic product"),
		},
		{
			Code: "ASM",
			Name: i18n.NewString("Trading partner identification number"),
		},
		{
			Code: "ASN",
			Name: i18n.NewString("Prior trading partner identification number"),
		},
		{
			Code: "ASO",
			Name: i18n.NewString("Password"),
		},
		{
			Code: "ASP",
			Name: i18n.NewString("Formal report number"),
		},
		{
			Code: "ASQ",
			Name: i18n.NewString("Fund account number"),
		},
		{
			Code: "ASR",
			Name: i18n.NewString("Safe custody number"),
		},
		{
			Code: "ASS",
			Name: i18n.NewString("Master account number"),
		},
		{
			Code: "AST",
			Name: i18n.NewString("Group reference number"),
		},
		{
			Code: "ASU",
			Name: i18n.NewString("Accounting transmission number"),
		},
		{
			Code: "ASV",
			Name: i18n.NewString("Product data file number"),
		},
		{
			Code: "ASW",
			Name: i18n.NewString("Cadastro Geral do Contribuinte (CGC)"),
		},
		{
			Code: "ASX",
			Name: i18n.NewString("Foreign resident identification number"),
		},
		{
			Code: "ASY",
			Name: i18n.NewString("CD-ROM"),
		},
		{
			Code: "ASZ",
			Name: i18n.NewString("Physical medium"),
		},
		{
			Code: "ATA",
			Name: i18n.NewString("Financial cancellation reference number"),
		},
		{
			Code: "ATB",
			Name: i18n.NewString("Purchase for export Customs agreement number"),
		},
		{
			Code: "ATC",
			Name: i18n.NewString("Judgment number"),
		},
		{
			Code: "ATD",
			Name: i18n.NewString("Secretariat number"),
		},
		{
			Code: "ATE",
			Name: i18n.NewString("Previous banking status message reference"),
		},
		{
			Code: "ATF",
			Name: i18n.NewString("Last received banking status message reference"),
		},
		{
			Code: "ATG",
			Name: i18n.NewString("Bank's documentary procedure reference"),
		},
		{
			Code: "ATH",
			Name: i18n.NewString("Customer's documentary procedure reference"),
		},
		{
			Code: "ATI",
			Name: i18n.NewString("Safe deposit box number"),
		},
		{
			Code: "ATJ",
			Name: i18n.NewString("Receiving Bankgiro number"),
		},
		{
			Code: "ATK",
			Name: i18n.NewString("Sending Bankgiro number"),
		},
		{
			Code: "ATL",
			Name: i18n.NewString("Bankgiro reference"),
		},
		{
			Code: "ATM",
			Name: i18n.NewString("Guarantee number"),
		},
		{
			Code: "ATN",
			Name: i18n.NewString("Collection instrument number"),
		},
		{
			Code: "ATO",
			Name: i18n.NewString("Converted Postgiro number"),
		},
		{
			Code: "ATP",
			Name: i18n.NewString("Cost centre alignment number"),
		},
		{
			Code: "ATQ",
			Name: i18n.NewString("Kamer Van Koophandel (KVK) number"),
		},
		{
			Code: "ATR",
			Name: i18n.NewString("Institut Belgo-Luxembourgeois de Codification (IBLC) number"),
		},
		{
			Code: "ATS",
			Name: i18n.NewString("External object reference"),
		},
		{
			Code: "ATT",
			Name: i18n.NewString("Exceptional transport authorisation number"),
		},
		{
			Code: "ATU",
			Name: i18n.NewString("Clave Unica de Identificacion Tributaria (CUIT)"),
		},
		{
			Code: "ATV",
			Name: i18n.NewString("Registro Unico Tributario (RUT)"),
		},
		{
			Code: "ATW",
			Name: i18n.NewString("Flat rack container bundle identification number"),
		},
		{
			Code: "ATX",
			Name: i18n.NewString("Transport equipment acceptance order reference"),
		},
		{
			Code: "ATY",
			Name: i18n.NewString("Transport equipment release order reference"),
		},
		{
			Code: "ATZ",
			Name: i18n.NewString("Ship's stay reference number"),
		},
		{
			Code: "AU",
			Name: i18n.NewString("Authorization to meet competition number"),
		},
		{
			Code: "AUA",
			Name: i18n.NewString("Place of positioning reference"),
		},
		{
			Code: "AUB",
			Name: i18n.NewString("Party reference"),
		},
		{
			Code: "AUC",
			Name: i18n.NewString("Issued prescription identification"),
		},
		{
			Code: "AUD",
			Name: i18n.NewString("Collection reference"),
		},
		{
			Code: "AUE",
			Name: i18n.NewString("Travel service"),
		},
		{
			Code: "AUF",
			Name: i18n.NewString("Consignment stock contract"),
		},
		{
			Code: "AUG",
			Name: i18n.NewString("Importer's letter of credit reference"),
		},
		{
			Code: "AUH",
			Name: i18n.NewString("Performed prescription identification"),
		},
		{
			Code: "AUI",
			Name: i18n.NewString("Image reference"),
		},
		{
			Code: "AUJ",
			Name: i18n.NewString("Proposed purchase order reference number"),
		},
		{
			Code: "AUK",
			Name: i18n.NewString("Application for financial support reference number"),
		},
		{
			Code: "AUL",
			Name: i18n.NewString("Manufacturing quality agreement number"),
		},
		{
			Code: "AUM",
			Name: i18n.NewString("Software editor reference"),
		},
		{
			Code: "AUN",
			Name: i18n.NewString("Software reference"),
		},
		{
			Code: "AUO",
			Name: i18n.NewString("Software quality reference"),
		},
		{
			Code: "AUP",
			Name: i18n.NewString("Consolidated orders' reference"),
		},
		{
			Code: "AUQ",
			Name: i18n.NewString("Customs binding ruling number"),
		},
		{
			Code: "AUR",
			Name: i18n.NewString("Customs non-binding ruling number"),
		},
		{
			Code: "AUS",
			Name: i18n.NewString("Delivery route reference"),
		},
		{
			Code: "AUT",
			Name: i18n.NewString("Net area supplier reference"),
		},
		{
			Code: "AUU",
			Name: i18n.NewString("Time series reference"),
		},
		{
			Code: "AUV",
			Name: i18n.NewString("Connecting point to central grid"),
		},
		{
			Code: "AUW",
			Name: i18n.NewString("Marketing plan identification number (MPIN)"),
		},
		{
			Code: "AUX",
			Name: i18n.NewString("Entity reference number, previous"),
		},
		{
			Code: "AUY",
			Name: i18n.NewString("International Standard Industrial Classification (ISIC)"),
		},
		{
			Code: "AUZ",
			Name: i18n.NewString("Customs pre-approval ruling number"),
		},
		{
			Code: "AV",
			Name: i18n.NewString("Account payable number"),
		},
		{
			Code: "AVA",
			Name: i18n.NewString("First financial institution's transaction reference"),
		},
		{
			Code: "AVB",
			Name: i18n.NewString("Product characteristics directory"),
		},
		{
			Code: "AVC",
			Name: i18n.NewString("Supplier's customer reference number"),
		},
		{
			Code: "AVD",
			Name: i18n.NewString("Inventory report request number"),
		},
		{
			Code: "AVE",
			Name: i18n.NewString("Metering point"),
		},
		{
			Code: "AVF",
			Name: i18n.NewString("Passenger reservation number"),
		},
		{
			Code: "AVG",
			Name: i18n.NewString("Slaughterhouse approval number"),
		},
		{
			Code: "AVH",
			Name: i18n.NewString("Meat cutting plant approval number"),
		},
		{
			Code: "AVI",
			Name: i18n.NewString("Customer travel service identifier"),
		},
		{
			Code: "AVJ",
			Name: i18n.NewString("Export control classification number"),
		},
		{
			Code: "AVK",
			Name: i18n.NewString("Broker reference 3"),
		},
		{
			Code: "AVL",
			Name: i18n.NewString("Consignment information"),
		},
		{
			Code: "AVM",
			Name: i18n.NewString("Goods item information"),
		},
		{
			Code: "AVN",
			Name: i18n.NewString("Dangerous Goods information"),
		},
		{
			Code: "AVO",
			Name: i18n.NewString("Pilotage services exemption number"),
		},
		{
			Code: "AVP",
			Name: i18n.NewString("Person registration number"),
		},
		{
			Code: "AVQ",
			Name: i18n.NewString("Place of packing approval number"),
		},
		{
			Code: "AVR",
			Name: i18n.NewString("Original Mandate Reference"),
		},
		{
			Code: "AVS",
			Name: i18n.NewString("Mandate Reference"),
		},
		{
			Code: "AVT",
			Name: i18n.NewString("Reservation station indentifier"),
		},
		{
			Code: "AVU",
			Name: i18n.NewString("Unique goods shipment identifier"),
		},
		{
			Code: "AVV",
			Name: i18n.NewString("Framework Agreement Number"),
		},
		{
			Code: "AVW",
			Name: i18n.NewString("Hash value"),
		},
		{
			Code: "AVX",
			Name: i18n.NewString("Movement reference number"),
		},
		{
			Code: "AVY",
			Name: i18n.NewString("Economic Operators Registration and Identification Number"),
		},
		{
			Code: "AVZ",
			Name: i18n.NewString("Local Reference Number"),
		},
		{
			Code: "AWA",
			Name: i18n.NewString("Rate code number"),
		},
		{
			Code: "AWB",
			Name: i18n.NewString("Air waybill number"),
		},
		{
			Code: "AWC",
			Name: i18n.NewString("Documentary credit amendment number"),
		},
		{
			Code: "AWD",
			Name: i18n.NewString("Advising bank's reference"),
		},
		{
			Code: "AWE",
			Name: i18n.NewString("Cost centre"),
		},
		{
			Code: "AWF",
			Name: i18n.NewString("Work item quantity determination"),
		},
		{
			Code: "AWG",
			Name: i18n.NewString("Internal data process number"),
		},
		{
			Code: "AWH",
			Name: i18n.NewString("Category of work reference"),
		},
		{
			Code: "AWI",
			Name: i18n.NewString("Policy form number"),
		},
		{
			Code: "AWJ",
			Name: i18n.NewString("Net area"),
		},
		{
			Code: "AWK",
			Name: i18n.NewString("Service provider"),
		},
		{
			Code: "AWL",
			Name: i18n.NewString("Error position"),
		},
		{
			Code: "AWM",
			Name: i18n.NewString("Service category reference"),
		},
		{
			Code: "AWN",
			Name: i18n.NewString("Connected location"),
		},
		{
			Code: "AWO",
			Name: i18n.NewString("Related party"),
		},
		{
			Code: "AWP",
			Name: i18n.NewString("Latest accounting entry record reference"),
		},
		{
			Code: "AWQ",
			Name: i18n.NewString("Accounting entry"),
		},
		{
			Code: "AWR",
			Name: i18n.NewString("Document reference, original"),
		},
		{
			Code: "AWS",
			Name: i18n.NewString("Hygienic Certificate number, national"),
		},
		{
			Code: "AWT",
			Name: i18n.NewString("Administrative Reference Code"),
		},
		{
			Code: "AWU",
			Name: i18n.NewString("Pick-up sheet number"),
		},
		{
			Code: "AWV",
			Name: i18n.NewString("Phone number"),
		},
		{
			Code: "AWW",
			Name: i18n.NewString("Buyer's fund number"),
		},
		{
			Code: "AWX",
			Name: i18n.NewString("Company trading account number"),
		},
		{
			Code: "AWY",
			Name: i18n.NewString("Reserved goods identifier"),
		},
		{
			Code: "AWZ",
			Name: i18n.NewString("Handling and movement reference number"),
		},
		{
			Code: "AXA",
			Name: i18n.NewString("Instruction to despatch reference number"),
		},
		{
			Code: "AXB",
			Name: i18n.NewString("Instruction for returns number"),
		},
		{
			Code: "AXC",
			Name: i18n.NewString("Metered services consumption report number"),
		},
		{
			Code: "AXD",
			Name: i18n.NewString("Order status enquiry number"),
		},
		{
			Code: "AXE",
			Name: i18n.NewString("Firm booking reference number"),
		},
		{
			Code: "AXF",
			Name: i18n.NewString("Product inquiry number"),
		},
		{
			Code: "AXG",
			Name: i18n.NewString("Split delivery number"),
		},
		{
			Code: "AXH",
			Name: i18n.NewString("Service relation number"),
		},
		{
			Code: "AXI",
			Name: i18n.NewString("Serial shipping container code"),
		},
		{
			Code: "AXJ",
			Name: i18n.NewString("Test specification number"),
		},
		{
			Code: "AXK",
			Name: i18n.NewString("Transport status report number"),
		},
		{
			Code: "AXL",
			Name: i18n.NewString("Tooling contract number"),
		},
		{
			Code: "AXM",
			Name: i18n.NewString("Formula reference number"),
		},
		{
			Code: "AXN",
			Name: i18n.NewString("Pre-agreement number"),
		},
		{
			Code: "AXO",
			Name: i18n.NewString("Product certification number"),
		},
		{
			Code: "AXP",
			Name: i18n.NewString("Consignment contract number"),
		},
		{
			Code: "AXQ",
			Name: i18n.NewString("Product specification reference number"),
		},
		{
			Code: "AXR",
			Name: i18n.NewString("Payroll deduction advice reference"),
		},
		{
			Code: "AXS",
			Name: i18n.NewString("TRACES party identification"),
		},
		{
			Code: "BA",
			Name: i18n.NewString("Beginning meter reading actual"),
		},
		{
			Code: "BC",
			Name: i18n.NewString("Buyer's contract number"),
		},
		{
			Code: "BD",
			Name: i18n.NewString("Bid number"),
		},
		{
			Code: "BE",
			Name: i18n.NewString("Beginning meter reading estimated"),
		},
		{
			Code: "BH",
			Name: i18n.NewString("House bill of lading number"),
		},
		{
			Code: "BM",
			Name: i18n.NewString("Bill of lading number"),
		},
		{
			Code: "BN",
			Name: i18n.NewString("Consignment identifier, carrier assigned"),
		},
		{
			Code: "BO",
			Name: i18n.NewString("Blanket order number"),
		},
		{
			Code: "BR",
			Name: i18n.NewString("Broker or sales office number"),
		},
		{
			Code: "BT",
			Name: i18n.NewString("Batch number/lot number"),
		},
		{
			Code: "BTP",
			Name: i18n.NewString("Battery and accumulator producer registration number"),
		},
		{
			Code: "BW",
			Name: i18n.NewString("Blended with number"),
		},
		{
			Code: "CAS",
			Name: i18n.NewString("IATA Cargo Agent CASS Address number"),
		},
		{
			Code: "CAT",
			Name: i18n.NewString("Matching of entries, balanced"),
		},
		{
			Code: "CAU",
			Name: i18n.NewString("Entry flagging"),
		},
		{
			Code: "CAV",
			Name: i18n.NewString("Matching of entries, unbalanced"),
		},
		{
			Code: "CAW",
			Name: i18n.NewString("Document reference, internal"),
		},
		{
			Code: "CAX",
			Name: i18n.NewString("European Value Added Tax identification"),
		},
		{
			Code: "CAY",
			Name: i18n.NewString("Cost accounting document"),
		},
		{
			Code: "CAZ",
			Name: i18n.NewString("Grid operator's customer reference number"),
		},
		{
			Code: "CBA",
			Name: i18n.NewString("Ticket control number"),
		},
		{
			Code: "CBB",
			Name: i18n.NewString("Order shipment grouping reference"),
		},
		{
			Code: "CD",
			Name: i18n.NewString("Credit note number"),
		},
		{
			Code: "CEC",
			Name: i18n.NewString("Ceding company"),
		},
		{
			Code: "CED",
			Name: i18n.NewString("Debit letter number"),
		},
		{
			Code: "CFE",
			Name: i18n.NewString("Consignee's further order"),
		},
		{
			Code: "CFF",
			Name: i18n.NewString("Animal farm licence number"),
		},
		{
			Code: "CFO",
			Name: i18n.NewString("Consignor's further order"),
		},
		{
			Code: "CG",
			Name: i18n.NewString("Consignee's order number"),
		},
		{
			Code: "CH",
			Name: i18n.NewString("Customer catalogue number"),
		},
		{
			Code: "CK",
			Name: i18n.NewString("Cheque number"),
		},
		{
			Code: "CKN",
			Name: i18n.NewString("Checking number"),
		},
		{
			Code: "CM",
			Name: i18n.NewString("Credit memo number"),
		},
		{
			Code: "CMR",
			Name: i18n.NewString("Road consignment note number"),
		},
		{
			Code: "CN",
			Name: i18n.NewString("Carrier's reference number"),
		},
		{
			Code: "CNO",
			Name: i18n.NewString("Charges note document attachment indicator"),
		},
		{
			Code: "COF",
			Name: i18n.NewString("Call off order number"),
		},
		{
			Code: "CP",
			Name: i18n.NewString("Condition of purchase document number"),
		},
		{
			Code: "CR",
			Name: i18n.NewString("Customer reference number"),
		},
		{
			Code: "CRN",
			Name: i18n.NewString("Transport means journey identifier"),
		},
		{
			Code: "CS",
			Name: i18n.NewString("Condition of sale document number"),
		},
		{
			Code: "CST",
			Name: i18n.NewString("Team assignment number"),
		},
		{
			Code: "CT",
			Name: i18n.NewString("Contract number"),
		},
		{
			Code: "CU",
			Name: i18n.NewString("Consignment identifier, consignor assigned"),
		},
		{
			Code: "CV",
			Name: i18n.NewString("Container operators reference number"),
		},
		{
			Code: "CW",
			Name: i18n.NewString("Package number"),
		},
		{
			Code: "CZ",
			Name: i18n.NewString("Cooperation contract number"),
		},
		{
			Code: "DA",
			Name: i18n.NewString("Deferment approval number"),
		},
		{
			Code: "DAN",
			Name: i18n.NewString("Debit account number"),
		},
		{
			Code: "DB",
			Name: i18n.NewString("Buyer's debtor number"),
		},
		{
			Code: "DI",
			Name: i18n.NewString("Distributor invoice number"),
		},
		{
			Code: "DL",
			Name: i18n.NewString("Debit note number"),
		},
		{
			Code: "DM",
			Name: i18n.NewString("Document identifier"),
		},
		{
			Code: "DQ",
			Name: i18n.NewString("Delivery note number"),
		},
		{
			Code: "DR",
			Name: i18n.NewString("Dock receipt number"),
		},
		{
			Code: "EA",
			Name: i18n.NewString("Ending meter reading actual"),
		},
		{
			Code: "EB",
			Name: i18n.NewString("Embargo permit number"),
		},
		{
			Code: "ED",
			Name: i18n.NewString("Export declaration"),
		},
		{
			Code: "EE",
			Name: i18n.NewString("Ending meter reading estimated"),
		},
		{
			Code: "EEP",
			Name: i18n.NewString("Electrical and electronic equipment producer registration"),
		},
		{
			Code: "EI",
			Name: i18n.NewString("Employer's identification number"),
		},
		{
			Code: "EN",
			Name: i18n.NewString("Embargo number"),
		},
		{
			Code: "EQ",
			Name: i18n.NewString("Equipment number"),
		},
		{
			Code: "ER",
			Name: i18n.NewString("Container/equipment receipt number"),
		},
		{
			Code: "ERN",
			Name: i18n.NewString("Exporter's reference number"),
		},
		{
			Code: "ET",
			Name: i18n.NewString("Excess transportation number"),
		},
		{
			Code: "EX",
			Name: i18n.NewString("Export permit identifier"),
		},
		{
			Code: "FC",
			Name: i18n.NewString("Fiscal number"),
		},
		{
			Code: "FF",
			Name: i18n.NewString("Consignment identifier, freight forwarder assigned"),
		},
		{
			Code: "FI",
			Name: i18n.NewString("File line identifier"),
		},
		{
			Code: "FLW",
			Name: i18n.NewString("Flow reference number"),
		},
		{
			Code: "FN",
			Name: i18n.NewString("Freight bill number"),
		},
		{
			Code: "FO",
			Name: i18n.NewString("Foreign exchange"),
		},
		{
			Code: "FS",
			Name: i18n.NewString("Final sequence number"),
		},
		{
			Code: "FT",
			Name: i18n.NewString("Free zone identifier"),
		},
		{
			Code: "FV",
			Name: i18n.NewString("File version number"),
		},
		{
			Code: "FX",
			Name: i18n.NewString("Foreign exchange contract number"),
		},
		{
			Code: "GA",
			Name: i18n.NewString("Standard's number"),
		},
		{
			Code: "GC",
			Name: i18n.NewString("Government contract number"),
		},
		{
			Code: "GD",
			Name: i18n.NewString("Standard's code number"),
		},
		{
			Code: "GDN",
			Name: i18n.NewString("General declaration number"),
		},
		{
			Code: "GN",
			Name: i18n.NewString("Government reference number"),
		},
		{
			Code: "HS",
			Name: i18n.NewString("Harmonised system number"),
		},
		{
			Code: "HWB",
			Name: i18n.NewString("House waybill number"),
		},
		{
			Code: "IA",
			Name: i18n.NewString("Internal vendor number"),
		},
		{
			Code: "IB",
			Name: i18n.NewString("In bond number"),
		},
		{
			Code: "ICA",
			Name: i18n.NewString("IATA cargo agent code number"),
		},
		{
			Code: "ICE",
			Name: i18n.NewString("Insurance certificate reference number"),
		},
		{
			Code: "ICO",
			Name: i18n.NewString("Insurance contract reference number"),
		},
		{
			Code: "II",
			Name: i18n.NewString("Initial sample inspection report number"),
		},
		{
			Code: "IL",
			Name: i18n.NewString("Internal order number"),
		},
		{
			Code: "INB",
			Name: i18n.NewString("Intermediary broker"),
		},
		{
			Code: "INN",
			Name: i18n.NewString("Interchange number new"),
		},
		{
			Code: "INO",
			Name: i18n.NewString("Interchange number old"),
		},
		{
			Code: "IP",
			Name: i18n.NewString("Import permit identifier"),
		},
		{
			Code: "IS",
			Name: i18n.NewString("Invoice number suffix"),
		},
		{
			Code: "IT",
			Name: i18n.NewString("Internal customer number"),
		},
		{
			Code: "IV",
			Name: i18n.NewString("Invoice document identifier"),
		},
		{
			Code: "JB",
			Name: i18n.NewString("Job number"),
		},
		{
			Code: "JE",
			Name: i18n.NewString("Ending job sequence number"),
		},
		{
			Code: "LA",
			Name: i18n.NewString("Shipping label serial number"),
		},
		{
			Code: "LAN",
			Name: i18n.NewString("Loading authorisation identifier"),
		},
		{
			Code: "LAR",
			Name: i18n.NewString("Lower number in range"),
		},
		{
			Code: "LB",
			Name: i18n.NewString("Lockbox"),
		},
		{
			Code: "LC",
			Name: i18n.NewString("Letter of credit number"),
		},
		{
			Code: "LI",
			Name: i18n.NewString("Document line identifier"),
		},
		{
			Code: "LO",
			Name: i18n.NewString("Load planning number"),
		},
		{
			Code: "LRC",
			Name: i18n.NewString("Reservation office identifier"),
		},
		{
			Code: "LS",
			Name: i18n.NewString("Bar coded label serial number"),
		},
		{
			Code: "MA",
			Name: i18n.NewString("Ship notice/manifest number"),
		},
		{
			Code: "MB",
			Name: i18n.NewString("Master bill of lading number"),
		},
		{
			Code: "MF",
			Name: i18n.NewString("Manufacturer's part number"),
		},
		{
			Code: "MG",
			Name: i18n.NewString("Meter unit number"),
		},
		{
			Code: "MH",
			Name: i18n.NewString("Manufacturing order number"),
		},
		{
			Code: "MR",
			Name: i18n.NewString("Message recipient"),
		},
		{
			Code: "MRN",
			Name: i18n.NewString("Mailing reference number"),
		},
		{
			Code: "MS",
			Name: i18n.NewString("Message sender"),
		},
		{
			Code: "MSS",
			Name: i18n.NewString("Manufacturer's material safety data sheet number"),
		},
		{
			Code: "MWB",
			Name: i18n.NewString("Master air waybill number"),
		},
		{
			Code: "NA",
			Name: i18n.NewString("North American hazardous goods classification number"),
		},
		{
			Code: "NF",
			Name: i18n.NewString("Nota Fiscal"),
		},
		{
			Code: "OH",
			Name: i18n.NewString("Current invoice number"),
		},
		{
			Code: "OI",
			Name: i18n.NewString("Previous invoice number"),
		},
		{
			Code: "ON",
			Name: i18n.NewString("Order document identifier, buyer assigned"),
		},
		{
			Code: "OP",
			Name: i18n.NewString("Original purchase order"),
		},
		{
			Code: "OR",
			Name: i18n.NewString("General order number"),
		},
		{
			Code: "PB",
			Name: i18n.NewString("Payer's financial institution account number"),
		},
		{
			Code: "PC",
			Name: i18n.NewString("Production code"),
		},
		{
			Code: "PD",
			Name: i18n.NewString("Promotion deal number"),
		},
		{
			Code: "PE",
			Name: i18n.NewString("Plant number"),
		},
		{
			Code: "PF",
			Name: i18n.NewString("Prime contractor contract number"),
		},
		{
			Code: "PI",
			Name: i18n.NewString("Price list version number"),
		},
		{
			Code: "PK",
			Name: i18n.NewString("Packing list number"),
		},
		{
			Code: "PL",
			Name: i18n.NewString("Price list number"),
		},
		{
			Code: "POR",
			Name: i18n.NewString("Purchase order response number"),
		},
		{
			Code: "PP",
			Name: i18n.NewString("Purchase order change number"),
		},
		{
			Code: "PQ",
			Name: i18n.NewString("Payment reference"),
		},
		{
			Code: "PR",
			Name: i18n.NewString("Price quote number"),
		},
		{
			Code: "PS",
			Name: i18n.NewString("Purchase order number suffix"),
		},
		{
			Code: "PW",
			Name: i18n.NewString("Prior purchase order number"),
		},
		{
			Code: "PY",
			Name: i18n.NewString("Payee's financial institution account number"),
		},
		{
			Code: "RA",
			Name: i18n.NewString("Remittance advice number"),
		},
		{
			Code: "RC",
			Name: i18n.NewString("Rail/road routing code"),
		},
		{
			Code: "RCN",
			Name: i18n.NewString("Railway consignment note number"),
		},
		{
			Code: "RE",
			Name: i18n.NewString("Release number"),
		},
		{
			Code: "REN",
			Name: i18n.NewString("Consignment receipt identifier"),
		},
		{
			Code: "RF",
			Name: i18n.NewString("Export reference number"),
		},
		{
			Code: "RR",
			Name: i18n.NewString("Payer's financial institution transit routing No.(ACH"),
		},
		{
			Code: "RT",
			Name: i18n.NewString("Payee's financial institution transit routing No."),
		},
		{
			Code: "SA",
			Name: i18n.NewString("Sales person number"),
		},
		{
			Code: "SB",
			Name: i18n.NewString("Sales region number"),
		},
		{
			Code: "SD",
			Name: i18n.NewString("Sales department number"),
		},
		{
			Code: "SE",
			Name: i18n.NewString("Serial number"),
		},
		{
			Code: "SEA",
			Name: i18n.NewString("Allocated seat"),
		},
		{
			Code: "SF",
			Name: i18n.NewString("Ship from"),
		},
		{
			Code: "SH",
			Name: i18n.NewString("Previous highest schedule number"),
		},
		{
			Code: "SI",
			Name: i18n.NewString("SID (Shipper's identifying number for shipment)"),
		},
		{
			Code: "SM",
			Name: i18n.NewString("Sales office number"),
		},
		{
			Code: "SN",
			Name: i18n.NewString("Transport equipment seal identifier"),
		},
		{
			Code: "SP",
			Name: i18n.NewString("Scan line"),
		},
		{
			Code: "SQ",
			Name: i18n.NewString("Equipment sequence number"),
		},
		{
			Code: "SRN",
			Name: i18n.NewString("Shipment reference number"),
		},
		{
			Code: "SS",
			Name: i18n.NewString("Sellers reference number"),
		},
		{
			Code: "STA",
			Name: i18n.NewString("Station reference number"),
		},
		{
			Code: "SW",
			Name: i18n.NewString("Swap order number"),
		},
		{
			Code: "SZ",
			Name: i18n.NewString("Specification number"),
		},
		{
			Code: "TB",
			Name: i18n.NewString("Trucker's bill of lading"),
		},
		{
			Code: "TCR",
			Name: i18n.NewString("Terminal operator's consignment reference"),
		},
		{
			Code: "TE",
			Name: i18n.NewString("Telex message number"),
		},
		{
			Code: "TF",
			Name: i18n.NewString("Transfer number"),
		},
		{
			Code: "TI",
			Name: i18n.NewString("TIR carnet number"),
		},
		{
			Code: "TIN",
			Name: i18n.NewString("Transport instruction number"),
		},
		{
			Code: "TL",
			Name: i18n.NewString("Tax exemption licence number"),
		},
		{
			Code: "TN",
			Name: i18n.NewString("Transaction reference number"),
		},
		{
			Code: "TP",
			Name: i18n.NewString("Test report number"),
		},
		{
			Code: "UAR",
			Name: i18n.NewString("Upper number of range"),
		},
		{
			Code: "UC",
			Name: i18n.NewString("Ultimate customer's reference number"),
		},
		{
			Code: "UCN",
			Name: i18n.NewString("Unique consignment reference number"),
		},
		{
			Code: "UN",
			Name: i18n.NewString("United Nations Dangerous Goods identifier"),
		},
		{
			Code: "UO",
			Name: i18n.NewString("Ultimate customer's order number"),
		},
		{
			Code: "URI",
			Name: i18n.NewString("Uniform Resource Identifier"),
		},
		{
			Code: "VA",
			Name: i18n.NewString("VAT registration number"),
		},
		{
			Code: "VC",
			Name: i18n.NewString("Vendor contract number"),
		},
		{
			Code: "VGR",
			Name: i18n.NewString("Transport equipment gross mass verification reference"),
		},
		{
			Code: "VM",
			Name: i18n.NewString("Vessel identifier"),
		},
		{
			Code: "VN",
			Name: i18n.NewString("Order number (vendor)"),
		},
		{
			Code: "VON",
			Name: i18n.NewString("Voyage number"),
		},
		{
			Code: "VOR",
			Name: i18n.NewString("Transport equipment gross mass verification order reference"),
		},
		{
			Code: "VP",
			Name: i18n.NewString("Vendor product number"),
		},
		{
			Code: "VR",
			Name: i18n.NewString("Vendor ID number"),
		},
		{
			Code: "VS",
			Name: i18n.NewString("Vendor order number suffix"),
		},
		{
			Code: "VT",
			Name: i18n.NewString("Motor vehicle identification number"),
		},
		{
			Code: "VV",
			Name: i18n.NewString("Voucher number"),
		},
		{
			Code: "WE",
			Name: i18n.NewString("Warehouse entry number"),
		},
		{
			Code: "WM",
			Name: i18n.NewString("Weight agreement number"),
		},
		{
			Code: "WN",
			Name: i18n.NewString("Well number"),
		},
		{
			Code: "WR",
			Name: i18n.NewString("Warehouse receipt number"),
		},
		{
			Code: "WS",
			Name: i18n.NewString("Warehouse storage location number"),
		},
		{
			Code: "WY",
			Name: i18n.NewString("Rail waybill number"),
		},
		{
			Code: "XA",
			Name: i18n.NewString("Company/place registration number"),
		},
		{
			Code: "XC",
			Name: i18n.NewString("Cargo control number"),
		},
		{
			Code: "XP",
			Name: i18n.NewString("Previous cargo control number"),
		},
		{
			Code: "ZZZ",
			Name: i18n.NewString("Mutually defined reference number"),
		},
	},
}

package untdid

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

const (
	// ExtKeyItemType is used to identify the UNTDID 7143 item type code.
	ExtKeyItemType cbc.Key = "untdid-item-type"
)

var extItemTypes = &cbc.Definition{
	Key: ExtKeyItemType,
	Name: i18n.String{
		i18n.EN: "UNTDID 7143 Item Type Identification Code",
	},
	Desc: i18n.String{
		i18n.EN: here.Doc(`

		`),
	},
	Values: []*cbc.Definition{
		{
			Code: "AA",
			Name: i18n.NewString("Product version number"),
		},
		{
			Code: "AB",
			Name: i18n.NewString("Assembly"),
		},
		{
			Code: "AC",
			Name: i18n.NewString("HIBC (Health Industry Bar Code)"),
		},
		{
			Code: "AD",
			Name: i18n.NewString("Cold roll number"),
		},
		{
			Code: "AE",
			Name: i18n.NewString("Hot roll number"),
		},
		{
			Code: "AF",
			Name: i18n.NewString("Slab number"),
		},
		{
			Code: "AG",
			Name: i18n.NewString("Software revision number"),
		},
		{
			Code: "AH",
			Name: i18n.NewString("UPC (Universal Product Code) Consumer package code (1-5-5)"),
		},
		{
			Code: "AI",
			Name: i18n.NewString("UPC (Universal Product Code) Consumer package code (1-5-5-1)"),
		},
		{
			Code: "AJ",
			Name: i18n.NewString("Sample number"),
		},
		{
			Code: "AK",
			Name: i18n.NewString("Pack number"),
		},
		{
			Code: "AL",
			Name: i18n.NewString("UPC (Universal Product Code) Shipping container code (1-2-5-5)"),
		},
		{
			Code: "AM",
			Name: i18n.NewString("UPC (Universal Product Code)/EAN (European article number) Shipping container code (1-2-5-5-1)"),
		},
		{
			Code: "AN",
			Name: i18n.NewString("UPC (Universal Product Code) suffix"),
		},
		{
			Code: "AO",
			Name: i18n.NewString("State label code"),
		},
		{
			Code: "AP",
			Name: i18n.NewString("Heat number"),
		},
		{
			Code: "AQ",
			Name: i18n.NewString("Coupon number"),
		},
		{
			Code: "AR",
			Name: i18n.NewString("Resource number"),
		},
		{
			Code: "AS",
			Name: i18n.NewString("Work task number"),
		},
		{
			Code: "AT",
			Name: i18n.NewString("Price look up number"),
		},
		{
			Code: "AU",
			Name: i18n.NewString("NSN (North Atlantic Treaty Organization Stock Number)"),
		},
		{
			Code: "AV",
			Name: i18n.NewString("Refined product code"),
		},
		{
			Code: "AW",
			Name: i18n.NewString("Exhibit"),
		},
		{
			Code: "AX",
			Name: i18n.NewString("End item"),
		},
		{
			Code: "AY",
			Name: i18n.NewString("Federal supply classification"),
		},
		{
			Code: "AZ",
			Name: i18n.NewString("Engineering data list"),
		},
		{
			Code: "BA",
			Name: i18n.NewString("Milestone event number"),
		},
		{
			Code: "BB",
			Name: i18n.NewString("Lot number"),
		},
		{
			Code: "BC",
			Name: i18n.NewString("National drug code 4-4-2 format"),
		},
		{
			Code: "BD",
			Name: i18n.NewString("National drug code 5-3-2 format"),
		},
		{
			Code: "BE",
			Name: i18n.NewString("National drug code 5-4-1 format"),
		},
		{
			Code: "BF",
			Name: i18n.NewString("National drug code 5-4-2 format"),
		},
		{
			Code: "BG",
			Name: i18n.NewString("National drug code"),
		},
		{
			Code: "BH",
			Name: i18n.NewString("Part number"),
		},
		{
			Code: "BI",
			Name: i18n.NewString("Local Stock Number (LSN)"),
		},
		{
			Code: "BJ",
			Name: i18n.NewString("Next higher assembly number"),
		},
		{
			Code: "BK",
			Name: i18n.NewString("Data category"),
		},
		{
			Code: "BL",
			Name: i18n.NewString("Control number"),
		},
		{
			Code: "BM",
			Name: i18n.NewString("Special material identification code"),
		},
		{
			Code: "BN",
			Name: i18n.NewString("Locally assigned control number"),
		},
		{
			Code: "BO",
			Name: i18n.NewString("Buyer's colour"),
		},
		{
			Code: "BP",
			Name: i18n.NewString("Buyer's part number"),
		},
		{
			Code: "BQ",
			Name: i18n.NewString("Variable measure product code"),
		},
		{
			Code: "BR",
			Name: i18n.NewString("Financial phase"),
		},
		{
			Code: "BS",
			Name: i18n.NewString("Contract breakdown"),
		},
		{
			Code: "BT",
			Name: i18n.NewString("Technical phase"),
		},
		{
			Code: "BU",
			Name: i18n.NewString("Dye lot number"),
		},
		{
			Code: "BV",
			Name: i18n.NewString("Daily statement of activities"),
		},
		{
			Code: "BW",
			Name: i18n.NewString("Periodical statement of activities within a bilaterally agreed time period"),
		},
		{
			Code: "BX",
			Name: i18n.NewString("Calendar week statement of activities"),
		},
		{
			Code: "BY",
			Name: i18n.NewString("Calendar month statement of activities"),
		},
		{
			Code: "BZ",
			Name: i18n.NewString("Original equipment number"),
		},
		{
			Code: "CC",
			Name: i18n.NewString("Industry commodity code"),
		},
		{
			Code: "CG",
			Name: i18n.NewString("Commodity grouping"),
		},
		{
			Code: "CL",
			Name: i18n.NewString("Colour number"),
		},
		{
			Code: "CR",
			Name: i18n.NewString("Contract number"),
		},
		{
			Code: "CV",
			Name: i18n.NewString("Customs article number"),
		},
		{
			Code: "DR",
			Name: i18n.NewString("Drawing revision number"),
		},
		{
			Code: "DW",
			Name: i18n.NewString("Drawing"),
		},
		{
			Code: "EC",
			Name: i18n.NewString("Engineering change level"),
		},
		{
			Code: "EF",
			Name: i18n.NewString("Material code"),
		},
		{
			Code: "EMD",
			Name: i18n.NewString("EMDN (European Medical Device Nomenclature)"),
		},
		{
			Code: "EN",
			Name: i18n.NewString("International Article Numbering Association (EAN)"),
		},
		{
			Code: "FS",
			Name: i18n.NewString("Fish species"),
		},
		{
			Code: "GB",
			Name: i18n.NewString("Buyer's internal product group code"),
		},
		{
			Code: "GMN",
			Name: i18n.NewString("Global model number"),
		},
		{
			Code: "GN",
			Name: i18n.NewString("National product group code"),
		},
		{
			Code: "GS",
			Name: i18n.NewString("General specification number"),
		},
		{
			Code: "HS",
			Name: i18n.NewString("Harmonised system"),
		},
		{
			Code: "IB",
			Name: i18n.NewString("ISBN (International Standard Book Number)"),
		},
		{
			Code: "IN",
			Name: i18n.NewString("Buyer's item number"),
		},
		{
			Code: "IS",
			Name: i18n.NewString("ISSN (International Standard Serial Number)"),
		},
		{
			Code: "IT",
			Name: i18n.NewString("Buyer's style number"),
		},
		{
			Code: "IZ",
			Name: i18n.NewString("Buyer's size code"),
		},
		{
			Code: "MA",
			Name: i18n.NewString("Machine number"),
		},
		{
			Code: "MF",
			Name: i18n.NewString("Manufacturer's (producer's) article number"),
		},
		{
			Code: "MN",
			Name: i18n.NewString("Model number"),
		},
		{
			Code: "MP",
			Name: i18n.NewString("Product/service identification number"),
		},
		{
			Code: "NB",
			Name: i18n.NewString("Batch number"),
		},
		{
			Code: "ON",
			Name: i18n.NewString("Customer order number"),
		},
		{
			Code: "PD",
			Name: i18n.NewString("Part number description"),
		},
		{
			Code: "PL",
			Name: i18n.NewString("Purchaser's order line number"),
		},
		{
			Code: "PO",
			Name: i18n.NewString("Purchase order number"),
		},
		{
			Code: "PV",
			Name: i18n.NewString("Promotional variant number"),
		},
		{
			Code: "QS",
			Name: i18n.NewString("Buyer's qualifier for size"),
		},
		{
			Code: "RC",
			Name: i18n.NewString("Returnable container number"),
		},
		{
			Code: "RN",
			Name: i18n.NewString("Release number"),
		},
		{
			Code: "RU",
			Name: i18n.NewString("Run number"),
		},
		{
			Code: "RY",
			Name: i18n.NewString("Record keeping of model year"),
		},
		{
			Code: "SA",
			Name: i18n.NewString("Supplier's article number"),
		},
		{
			Code: "SG",
			Name: i18n.NewString("Standard group of products (mixed assortment)"),
		},
		{
			Code: "SK",
			Name: i18n.NewString("SKU (Stock keeping unit)"),
		},
		{
			Code: "SN",
			Name: i18n.NewString("Serial number"),
		},
		{
			Code: "SRS",
			Name: i18n.NewString("RSK number"),
		},
		{
			Code: "SRT",
			Name: i18n.NewString("IFLS (Institut Francais du Libre Service) 5 digit product"),
		},
		{
			Code: "SRU",
			Name: i18n.NewString("IFLS (Institut Francais du Libre Service) 9 digit product"),
		},
		{
			Code: "SRV",
			Name: i18n.NewString("GS1 Global Trade Item Number"),
		},
		{
			Code: "SRW",
			Name: i18n.NewString("EDIS (Energy Data Identification System)"),
		},
		{
			Code: "SRX",
			Name: i18n.NewString("Slaughter number"),
		},
		{
			Code: "SRY",
			Name: i18n.NewString("Official animal number"),
		},
		{
			Code: "SRZ",
			Name: i18n.NewString("Harmonized tariff schedule"),
		},
		{
			Code: "SS",
			Name: i18n.NewString("Supplier's supplier article number"),
		},
		{
			Code: "SSA",
			Name: i18n.NewString("46 Level DOT Code"),
		},
		{
			Code: "SSB",
			Name: i18n.NewString("Airline Tariff 6D"),
		},
		{
			Code: "SSC",
			Name: i18n.NewString("Title 49 Code of Federal Regulations"),
		},
		{
			Code: "SSD",
			Name: i18n.NewString("International Civil Aviation Administration code"),
		},
		{
			Code: "SSE",
			Name: i18n.NewString("Hazardous Materials ID DOT"),
		},
		{
			Code: "SSF",
			Name: i18n.NewString("Endorsement"),
		},
		{
			Code: "SSG",
			Name: i18n.NewString("Air Force Regulation 71-4"),
		},
		{
			Code: "SSH",
			Name: i18n.NewString("Breed"),
		},
		{
			Code: "SSI",
			Name: i18n.NewString("Chemical Abstract Service (CAS) registry number"),
		},
		{
			Code: "SSJ",
			Name: i18n.NewString("Engine model designation"),
		},
		{
			Code: "SSK",
			Name: i18n.NewString("Institutional Meat Purchase Specifications (IMPS) Number"),
		},
		{
			Code: "SSL",
			Name: i18n.NewString("Price Look-Up code (PLU)"),
		},
		{
			Code: "SSM",
			Name: i18n.NewString("International Maritime Organization (IMO) Code"),
		},
		{
			Code: "SSN",
			Name: i18n.NewString("Bureau of Explosives 600-A (rail)"),
		},
		{
			Code: "SSO",
			Name: i18n.NewString("United Nations Dangerous Goods List"),
		},
		{
			Code: "SSP",
			Name: i18n.NewString("International Code of Botanical Nomenclature (ICBN)"),
		},
		{
			Code: "SSQ",
			Name: i18n.NewString("International Code of Zoological Nomenclature (ICZN)"),
		},
		{
			Code: "SSR",
			Name: i18n.NewString("International Code of Nomenclature for Cultivated Plants"),
		},
		{
			Code: "SSS",
			Name: i18n.NewString("Distributor’s article identifier"),
		},
		{
			Code: "SST",
			Name: i18n.NewString("Norwegian Classification system ENVA"),
		},
		{
			Code: "SSU",
			Name: i18n.NewString("Supplier assigned classification"),
		},
		{
			Code: "SSV",
			Name: i18n.NewString("Mexican classification system AMECE"),
		},
		{
			Code: "SSW",
			Name: i18n.NewString("German classification system CCG"),
		},
		{
			Code: "SSX",
			Name: i18n.NewString("Finnish classification system EANFIN"),
		},
		{
			Code: "SSY",
			Name: i18n.NewString("Canadian classification system ICC"),
		},
		{
			Code: "SSZ",
			Name: i18n.NewString("French classification system IFLS5"),
		},
		{
			Code: "ST",
			Name: i18n.NewString("Style number"),
		},
		{
			Code: "STA",
			Name: i18n.NewString("Dutch classification system CBL"),
		},
		{
			Code: "STB",
			Name: i18n.NewString("Japanese classification system JICFS"),
		},
		{
			Code: "STC",
			Name: i18n.NewString("European Union dairy subsidy eligibility classification"),
		},
		{
			Code: "STD",
			Name: i18n.NewString("GS1 Spain classification system"),
		},
		{
			Code: "STE",
			Name: i18n.NewString("GS1 Poland classification system"),
		},
		{
			Code: "STF",
			Name: i18n.NewString("Federal Agency on Technical Regulating and Metrology of the"),
		},
		{
			Code: "STG",
			Name: i18n.NewString("Efficient Consumer Response (ECR) Austria classification"),
		},
		{
			Code: "STH",
			Name: i18n.NewString("GS1 Italy classification system"),
		},
		{
			Code: "STI",
			Name: i18n.NewString("CPV (Common Procurement Vocabulary)"),
		},
		{
			Code: "STJ",
			Name: i18n.NewString("IFDA (International Foodservice Distributors Association)"),
		},
		{
			Code: "STK",
			Name: i18n.NewString("AHFS (American Hospital Formulary Service) pharmacologic -"),
		},
		{
			Code: "STL",
			Name: i18n.NewString("ATC (Anatomical Therapeutic Chemical) classification system"),
		},
		{
			Code: "STM",
			Name: i18n.NewString("CLADIMED (Classification des Dispositifs Médicaux)"),
		},
		{
			Code: "STN",
			Name: i18n.NewString("CMDR (Canadian Medical Device Regulations) classification"),
		},
		{
			Code: "STO",
			Name: i18n.NewString("CNDM (Classificazione Nazionale dei Dispositivi Medici)"),
		},
		{
			Code: "STP",
			Name: i18n.NewString("UK DM&D (Dictionary of Medicines & Devices) standard coding"),
		},
		{
			Code: "STQ",
			Name: i18n.NewString("eCl@ss"),
		},
		{
			Code: "STR",
			Name: i18n.NewString("EDMA (European Diagnostic Manufacturers Association)"),
		},
		{
			Code: "STS",
			Name: i18n.NewString("EGAR (European Generic Article Register)"),
		},
		{
			Code: "STT",
			Name: i18n.NewString("GMDN (Global Medical Devices Nomenclature)"),
		},
		{
			Code: "STU",
			Name: i18n.NewString("GPI (Generic Product Identifier)"),
		},
		{
			Code: "STV",
			Name: i18n.NewString("HCPCS (Healthcare Common Procedure Coding System)"),
		},
		{
			Code: "STW",
			Name: i18n.NewString("ICPS (International Classification for Patient Safety)"),
		},
		{
			Code: "STX",
			Name: i18n.NewString("MedDRA (Medical Dictionary for Regulatory Activities)"),
		},
		{
			Code: "STY",
			Name: i18n.NewString("Medical Columbus"),
		},
		{
			Code: "STZ",
			Name: i18n.NewString("NAPCS (North American Product Classification System)"),
		},
		{
			Code: "SUA",
			Name: i18n.NewString("NHS (National Health Services) eClass"),
		},
		{
			Code: "SUB",
			Name: i18n.NewString("US FDA (Food and Drug Administration) Product Code"),
		},
		{
			Code: "SUC",
			Name: i18n.NewString("SNOMED CT (Systematized Nomenclature of Medicine-Clinical"),
		},
		{
			Code: "SUD",
			Name: i18n.NewString("UMDNS (Universal Medical Device Nomenclature System)"),
		},
		{
			Code: "SUE",
			Name: i18n.NewString("GS1 Global Returnable Asset Identifier, non-serialised"),
		},
		{
			Code: "SUF",
			Name: i18n.NewString("IMEI"),
		},
		{
			Code: "SUG",
			Name: i18n.NewString("Waste Type (EMSA)"),
		},
		{
			Code: "SUH",
			Name: i18n.NewString("Ship's store classification type"),
		},
		{
			Code: "SUI",
			Name: i18n.NewString("Emergency fire code"),
		},
		{
			Code: "SUJ",
			Name: i18n.NewString("Emergency spillage code"),
		},
		{
			Code: "SUK",
			Name: i18n.NewString("IMDG packing group"),
		},
		{
			Code: "SUL",
			Name: i18n.NewString("MARPOL Code IBC"),
		},
		{
			Code: "SUM",
			Name: i18n.NewString("IMDG subsidiary risk class"),
		},
		{
			Code: "TG",
			Name: i18n.NewString("Transport group number"),
		},
		{
			Code: "TSN",
			Name: i18n.NewString("Taxonomic Serial Number"),
		},
		{
			Code: "TSO",
			Name: i18n.NewString("IMDG main hazard class"),
		},
		{
			Code: "TSP",
			Name: i18n.NewString("EU Combined Nomenclature"),
		},
		{
			Code: "TSQ",
			Name: i18n.NewString("Therapeutic classification number"),
		},
		{
			Code: "TSR",
			Name: i18n.NewString("European Waste Catalogue"),
		},
		{
			Code: "TSS",
			Name: i18n.NewString("Price grouping code"),
		},
		{
			Code: "TST",
			Name: i18n.NewString("UNSPSC"),
		},
		{
			Code: "TSU",
			Name: i18n.NewString("EU RoHS Directive"),
		},
		{
			Code: "UA",
			Name: i18n.NewString("Ultimate customer's article number"),
		},
		{
			Code: "UP",
			Name: i18n.NewString("UPC (Universal product code)"),
		},
		{
			Code: "VN",
			Name: i18n.NewString("Vendor item number"),
		},
		{
			Code: "VP",
			Name: i18n.NewString("Vendor's (seller's) part number"),
		},
		{
			Code: "VS",
			Name: i18n.NewString("Vendor's supplemental item number"),
		},
		{
			Code: "VX",
			Name: i18n.NewString("Vendor specification number"),
		},
		{
			Code: "ZZZ",
			Name: i18n.NewString("Mutually defined"),
		},
	},
}

package gr

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// Regime extension codes.
const (
	ExtKeyMyDATAVATCat     = "gr-mydata-vat-cat"
	ExtKeyMyDATAExemption  = "gr-mydata-exemption"
	ExtKeyMyDATAIncomeCat  = "gr-mydata-income-cat"
	ExtKeyMyDATAIncomeType = "gr-mydata-income-type"
)

var extensionKeys = []*cbc.KeyDefinition{
	{
		Key: ExtKeyMyDATAVATCat,
		Name: i18n.String{
			i18n.EN: "VAT category",
			i18n.EL: "Κατηγορία ΦΠΑ",
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "1",
				Name: i18n.String{
					i18n.EN: "Standard rate",
					i18n.EL: "Κανονικός συντελεστής",
				},
			},
			{
				Value: "2",
				Name: i18n.String{
					i18n.EN: "Reduced rate",
					i18n.EL: "Μειωμένος συντελεστής",
				},
			},
			{
				Value: "3",
				Name: i18n.String{
					i18n.EN: "Super-Reduced Rate",
					i18n.EL: "Υπερμειωμένος συντελεστής",
				},
			},
			{
				Value: "4",
				Name: i18n.String{
					i18n.EN: "Standard rate (Island)",
					i18n.EL: "Κανονικός συντελεστής (Νησί)",
				},
			},
			{
				Value: "5",
				Name: i18n.String{
					i18n.EN: "Reduced rate (Island)",
					i18n.EL: "Μειωμένος συντελεστής (Νησί)",
				},
			},
			{
				Value: "6",
				Name: i18n.String{
					i18n.EN: "Super-reduced rate (Island)",
					i18n.EL: "Υπερμειωμένος συντελεστής (Νησί)",
				},
			},
			{
				Value: "7",
				Name: i18n.String{
					i18n.EN: "Without VAT",
					i18n.EL: "Άνευ ΦΠΑ",
				},
			},
			{
				Value: "8",
				Name: i18n.String{
					i18n.EN: "Records without VAT (e.g. Payroll, Amortisations)",
					i18n.EL: "Εγγραφές χωρίς ΦΠΑ (πχ Μισθοδοσία, Αποσβέσεις)",
				},
			},
		},
	},
	{
		Key: ExtKeyMyDATAExemption,
		Name: i18n.String{
			i18n.EN: "VAT exemption cause",
			i18n.EL: "Κατηγορία Αιτίας Εξαίρεσης ΦΠΑ",
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "1",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 3 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 3 του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "2",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 5 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 5 του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "3",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 13 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 13 του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "4",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 14 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 14 του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "5",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 16 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 16 του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "6",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 19 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 19 του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "7",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 22 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 22 του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "8",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 24 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 24 του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "9",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 25 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 25 του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "10",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 26 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 26 του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "11",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 27 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 27 του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "12",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 27 - Seagoing Vessels of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 27 - Πλοία Ανοικτής Θαλάσσης του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "13",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 27.1.γ - Seagoing Vessels of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 27.1.γ - Πλοία Ανοικτής Θαλάσσης του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "14",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 28 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 28 του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "15",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 39 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 39 του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "16",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 39a of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 39α του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "17",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 40 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 40 του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "18",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 41 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 41 του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "19",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 47 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 47 του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "20",
				Name: i18n.String{
					i18n.EN: "VAT included - article 43 of the VAT code",
					i18n.EL: "ΦΠΑ εμπεριεχόμενος - άρθρο 43 του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "21",
				Name: i18n.String{
					i18n.EN: "VAT included - article 44 of the VAT code",
					i18n.EL: "ΦΠΑ εμπεριεχόμενος - άρθρο 44 του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "22",
				Name: i18n.String{
					i18n.EN: "VAT included - article 45 of the VAT code",
					i18n.EL: "ΦΠΑ εμπεριεχόμενος - άρθρο 45 του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "23",
				Name: i18n.String{
					i18n.EN: "VAT included - article 46 of the VAT code",
					i18n.EL: "ΦΠΑ εμπεριεχόμενος - άρθρο 46 του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "24",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 6 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 6 του Κώδικα ΦΠΑ",
				},
			},
			{
				Value: "25",
				Name: i18n.String{
					i18n.EN: "Without VAT - ΠΟΛ.1029/1995",
					i18n.EL: "Χωρίς ΦΠΑ - ΠΟΛ.1029/1995",
				},
			},
			{
				Value: "26",
				Name: i18n.String{
					i18n.EN: "Without VAT - ΠΟΛ.1167/2015",
					i18n.EL: "Χωρίς ΦΠΑ - ΠΟΛ.1167/2015",
				},
			},
			{
				Value: "27",
				Name: i18n.String{
					i18n.EN: "Without VAT - Other VAT exceptions",
					i18n.EL: "Λοιπές Εξαιρέσεις ΦΠΑ",
				},
			},
			{
				Value: "28",
				Name: i18n.String{
					i18n.EN: "Without VAT - Article 24 (b) (1) of the VAT Code (Tax Free)",
					i18n.EL: "Χωρίς ΦΠΑ – άρθρο 24 περ. β' παρ.1 του Κώδικα ΦΠΑ, (Tax Free)",
				},
			},
			{
				Value: "29",
				Name: i18n.String{
					i18n.EN: "Without VAT - Article 47b of the VAT Code (OSS non-EU scheme)",
					i18n.EL: "Χωρίς ΦΠΑ – άρθρο 47β, του Κώδικα ΦΠΑ (OSS μη ενωσιακό καθεστώς)",
				},
			},
			{
				Value: "30",
				Name: i18n.String{
					i18n.EN: "Without VAT - Article 47c of the VAT Code (OSS EU scheme)",
					i18n.EL: "Χωρίς ΦΠΑ – άρθρο 47γ, του Κώδικα ΦΠΑ (OSS ενωσιακό καθεστώς)",
				},
			},
			{
				Value: "31",
				Name: i18n.String{
					i18n.EN: "Excluding VAT - Article 47d of the VAT Code (IOSS)",
					i18n.EL: "Χωρίς ΦΠΑ – άρθρο 47δ του Κώδικα ΦΠΑ (IOSS)",
				},
			},
		},
	},
	{
		Key: ExtKeyMyDATAIncomeCat,
		Name: i18n.String{
			i18n.EN: "Income Classification Category",
			i18n.EL: "Κωδικός Κατηγορίας Χαρακτηρισμού Εσόδων",
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "category1_1",
				Name: i18n.String{
					i18n.EN: "Commodity Sale Income (+)/(-)",
					i18n.EL: "Έσοδα από Πώληση Εμπορευμάτων (+)/(-)",
				},
			},
			{
				Value: "category1_2",
				Name: i18n.String{
					i18n.EN: "Product Sale Income (+)/(-)",
					i18n.EL: "Έσοδα από Πώληση Προϊόντων (+)/(-)",
				},
			},
			{
				Value: "category1_3",
				Name: i18n.String{
					i18n.EN: "Provision of Services Income (+)/(-)",
					i18n.EL: "Έσοδα από Παροχή Υπηρεσιών (+)/(-)",
				},
			},
			{
				Value: "category1_4",
				Name: i18n.String{
					i18n.EN: "Sale of Fixed Assets Income (+)/(-)",
					i18n.EL: "Έσοδα από Πώληση Παγίων (+)/(-)",
				},
			},
			{
				Value: "category1_5",
				Name: i18n.String{
					i18n.EN: "Other Income/Profits (+)/(-)",
					i18n.EL: "Λοιπά Έσοδα/ Κέρδη (+)/(-)",
				},
			},
			{
				Value: "category1_6",
				Name: i18n.String{
					i18n.EN: "Self-Deliveries/Self-Supplies (+)/(-)",
					i18n.EL: "Αυτοπαραδόσεις / Ιδιοχρησιμοποιήσεις (+)/(-)",
				},
			},
			{
				Value: "category1_7",
				Name: i18n.String{
					i18n.EN: "Income on behalf of Third Parties (+)/(-)",
					i18n.EL: "Έσοδα για λ/σμο τρίτων (+)/(-)",
				},
			},
			{
				Value: "category1_8",
				Name: i18n.String{
					i18n.EN: "Past fiscal years income (+)/(-)",
					i18n.EL: "Έσοδα προηγούμενων χρήσεων (+)/ (-)",
				},
			},
			{
				Value: "category1_9",
				Name: i18n.String{
					i18n.EN: "Future fiscal years income (+)/(-)",
					i18n.EL: "Έσοδα επομένων χρήσεων (+)/(-)",
				},
			},
			{
				Value: "category1_10",
				Name: i18n.String{
					i18n.EN: "Other Income Adjustment/Regularisation Entries (+)/(-)",
					i18n.EL: "Λοιπές Εγγραφές Τακτοποίησης Εσόδων (+)/(-)",
				},
			},
			{
				Value: "category1_95",
				Name: i18n.String{
					i18n.EN: "Other Income-related Information (+)/(-)",
					i18n.EL: "Λοιπά Πληροφοριακά Στοιχεία Εσόδων (+)/(-)",
				},
			},
		},
	},
	{
		Key: ExtKeyMyDATAIncomeType,
		Name: i18n.String{
			i18n.EN: "Income Classification Type",
			i18n.EL: "Κωδικός Τύπου Χαρακτηρισμού Εσόδων",
		},
		Values: []*cbc.ValueDefinition{
			{
				Value: "E3_106",
				Name: i18n.String{
					i18n.EN: "Self-Production of Fixed Assets – Self-Deliveries – Destroying inventory/Commodities",
					i18n.EL: "Ιδιοπαραγωγή παγίων - Αυτοπαραδόσεις - Καταστροφές αποθεμάτων/Εμπορεύματα",
				},
			},
			{
				Value: "E3_205",
				Name: i18n.String{
					i18n.EN: "Self-Production of Fixed Assets – Self-Deliveries – Destroying inventory/Raw and other materials",
					i18n.EL: "Ιδιοπαραγωγή παγίων - Αυτοπαραδόσεις - Καταστροφές αποθεμάτων/Πρώτες ύλες και λοιπά υλικά",
				},
			},
			{
				Value: "E3_210",
				Name: i18n.String{
					i18n.EN: "Self-Production of Fixed Assets – Self-Deliveries – Destroying inventory/Products and production in progress",
					i18n.EL: "Ιδιοπαραγωγή παγίων - Αυτοπαραδόσεις - Καταστροφές αποθεμάτων/Προϊόντα και παραγωγή σε εξέλιξη",
				},
			},
			{
				Value: "E3_305",
				Name: i18n.String{
					i18n.EN: "Self-Production of Fixed Assets – Self-Deliveries – Destroying inventory/Raw and other materials",
					i18n.EL: "Ιδιοπαραγωγή παγίων - Αυτοπαραδόσεις – Καταστροφές αποθεμάτων/Πρώτες ύλες και λοιπά υλικά",
				},
			},
			{
				Value: "E3_310",
				Name: i18n.String{
					i18n.EN: "Self-Production of Fixed Assets – Self-Deliveries – Destroying inventory/Products and production in progress",
					i18n.EL: "Ιδιοπαραγωγή παγίων - Αυτοπαραδόσεις - Καταστροφές αποθεμάτων/Προϊόντα και παραγωγή σε εξέλιξη",
				},
			},
			{
				Value: "E3_318",
				Name: i18n.String{
					i18n.EN: "Self-Production of Fixed Assets – Self-Deliveries – Destroying inventory/Production expenses",
					i18n.EL: "Ιδιοπαραγωγή παγίων - Αυτοπαραδόσεις - Καταστροφές αποθεμάτων/Έξοδα παραγωγής",
				},
			},
			{
				Value: "E3_561_001",
				Name: i18n.String{
					i18n.EN: "Wholesale Sales of Goods and Services – for Traders",
					i18n.EL: "Πωλήσεις αγαθών και υπηρεσιών Χονδρικές - Επιτηδευματιών",
				},
			},
			{
				Value: "E3_561_002",
				Name: i18n.String{
					i18n.EN: "Wholesale Sales of Goods and Services pursuant to article 39a paragraph 5 of the VAT Code (Law 2859/2000)",
					i18n.EL: "Πωλήσεις αγαθών και υπηρεσιών Χονδρικές βάσει άρθρου 39α παρ 5 του Κώδικα Φ.Π.Α. (Ν.2859/2000)",
				},
			},
			{
				Value: "E3_561_003",
				Name: i18n.String{
					i18n.EN: "Retail Sales of Goods and Services – Private Clientele",
					i18n.EL: "Πωλήσεις αγαθών και υπηρεσιών Λιανικές - Ιδιωτική Πελατεία",
				},
			},
			{
				Value: "E3_561_004",
				Name: i18n.String{
					i18n.EN: "Retail Sales of Goods and Services pursuant to article 39a paragraph 5 of the VAT Code (Law 2859/2000)",
					i18n.EL: "Πωλήσεις αγαθών και υπηρεσιών Λιανικές βάσει άρθρου 39α παρ 5 του Κώδικα Φ.Π.Α. (Ν.2859/2000)",
				},
			},
			{
				Value: "E3_561_005",
				Name: i18n.String{
					i18n.EN: "Intra-Community Foreign Sales of Goods and Services",
					i18n.EL: "Πωλήσεις αγαθών και υπηρεσιών Εξωτερικού Ενδοκοινοτικές",
				},
			},
			{
				Value: "E3_561_006",
				Name: i18n.String{
					i18n.EN: "Third Country Foreign Sales of Goods and Services",
					i18n.EL: "Πωλήσεις αγαθών και υπηρεσιών Εξωτερικού Τρίτες Χώρες",
				},
			},
			{
				Value: "E3_561_007",
				Name: i18n.String{
					i18n.EN: "Other Sales of Goods and Services",
					i18n.EL: "Πωλήσεις αγαθών και υπηρεσιών Λοιπά",
				},
			},
			{
				Value: "E3_562",
				Name: i18n.String{
					i18n.EN: "Other Ordinary Income",
					i18n.EL: "Λοιπά συνήθη έσοδα",
				},
			},
			{
				Value: "E3_563",
				Name: i18n.String{
					i18n.EN: "Credit Interest and Related Income",
					i18n.EL: "Πιστωτικοί τόκοι και συναφή έσοδα",
				},
			},
			{
				Value: "E3_564",
				Name: i18n.String{
					i18n.EN: "Credit Exchange Differences",
					i18n.EL: "Πιστωτικές συναλλαγματικές διαφορές",
				},
			},
			{
				Value: "E3_565",
				Name: i18n.String{
					i18n.EN: "Income from Participations",
					i18n.EL: "Έσοδα συμμετοχών",
				},
			},
			{
				Value: "E3_566",
				Name: i18n.String{
					i18n.EN: "Profits from Disposing Non-Current Assets",
					i18n.EL: "Κέρδη από διάθεση μη κυκλοφορούντων περιουσιακών στοιχείων",
				},
			},
			{
				Value: "E3_567",
				Name: i18n.String{
					i18n.EN: "Profits from the Reversal of Provisions and Impairments",
					i18n.EL: "Κέρδη από αναστροφή προβλέψεων και απομειώσεων",
				},
			},
			{
				Value: "E3_568",
				Name: i18n.String{
					i18n.EN: "Profits from Measurement at Fair Value",
					i18n.EL: "Κέρδη από επιμέτρηση στην εύλογη αξία",
				},
			},
			{
				Value: "E3_570",
				Name: i18n.String{
					i18n.EN: "Extraordinary income and profits",
					i18n.EL: "Ασυνήθη έσοδα και κέρδη",
				},
			},
			{
				Value: "E3_595",
				Name: i18n.String{
					i18n.EN: "Self-Production Expenses",
					i18n.EL: "Έξοδα σε ιδιοπαραγωγή",
				},
			},
			{
				Value: "E3_596",
				Name: i18n.String{
					i18n.EN: "Subsidies - Grants",
					i18n.EL: "Επιδοτήσεις - Επιχορηγήσεις",
				},
			},
			{
				Value: "E3_597",
				Name: i18n.String{
					i18n.EN: "Subsidies – Grants for Investment Purposes – Expense Coverage",
					i18n.EL: "Επιδοτήσεις - Επιχορηγήσεις για επενδυτικούς σκοπούς - κάλυψη δαπανών",
				},
			},
			{
				Value: "E3_880_001",
				Name: i18n.String{
					i18n.EN: "Wholesale Sales of Fixed Assets",
					i18n.EL: "Πωλήσεις Παγίων Χονδρικές",
				},
			},
			{
				Value: "E3_880_002",
				Name: i18n.String{
					i18n.EN: "Retail Sales of Fixed Assets",
					i18n.EL: "Πωλήσεις Παγίων Λιανικές",
				},
			},
			{
				Value: "E3_880_003",
				Name: i18n.String{
					i18n.EN: "Intra-Community Foreign Sales of Fixed Assets",
					i18n.EL: "Πωλήσεις Παγίων Εξωτερικού Ενδοκοινοτικές",
				},
			},
			{
				Value: "E3_880_004",
				Name: i18n.String{
					i18n.EN: "Third Country Foreign Sales of Fixed Assets",
					i18n.EL: "Πωλήσεις Παγίων Εξωτερικού Τρίτες Χώρες",
				},
			},
			{
				Value: "E3_881_001",
				Name: i18n.String{
					i18n.EN: "Wholesale Sales on behalf of Third Parties",
					i18n.EL: "Πωλήσεις για λογ/σμο Τρίτων Χονδρικές",
				},
			},
			{
				Value: "E3_881_002",
				Name: i18n.String{
					i18n.EN: "Retail Sales on behalf of Third Parties",
					i18n.EL: "Πωλήσεις για λογ/σμο Τρίτων Λιανικές",
				},
			},
			{
				Value: "E3_881_003",
				Name: i18n.String{
					i18n.EN: "Intra-Community Foreign Sales on behalf of Third Parties",
					i18n.EL: "Πωλήσεις για λογ/σμο Τρίτων Εξωτερικού Ενδοκοινοτικές",
				},
			},
			{
				Value: "E3_881_004",
				Name: i18n.String{
					i18n.EN: "Third Country Foreign Sales on behalf of Third Parties",
					i18n.EL: "Πωλήσεις για λογ/σμο Τρίτων Εξωτερικού Τρίτες Χώρες",
				},
			},
			{
				Value: "E3_598_001",
				Name: i18n.String{
					i18n.EN: "Sales of goods belonging to excise duty",
					i18n.EL: "Πωλήσεις αγαθών που υπάγονται σε ΕΦΚ",
				},
			},
			{
				Value: "E3_598_003",
				Name: i18n.String{
					i18n.EN: "Sales on behalf of farmers through an agricultural cooperative e.t.c.",
					i18n.EL: "Πωλήσεις για λογαριασμό αγροτών μέσω αγροτικού συνεταιρισμού κ.λ.π.",
				},
			},
		},
	},
}

package mydata

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// Regime extension codes.
const (
	ExtKeyVATRate      = "gr-mydata-vat-rate"
	ExtKeyInvoiceType  = "gr-mydata-invoice-type"
	ExtKeyExemption    = "gr-mydata-exemption"
	ExtKeyIncomeCat    = "gr-mydata-income-cat"
	ExtKeyIncomeType   = "gr-mydata-income-type"
	ExtKeyPaymentMeans = "gr-mydata-payment-means"

	InvoiceTypeRetailPrefix = "11."
)

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyVATRate,
		Name: i18n.String{
			i18n.EN: "VAT rate",
			i18n.EL: "Κατηγορία ΦΠΑ",
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Standard rate",
					i18n.EL: "Κανονικός συντελεστής",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Reduced rate",
					i18n.EL: "Μειωμένος συντελεστής",
				},
			},
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "Super-Reduced Rate",
					i18n.EL: "Υπερμειωμένος συντελεστής",
				},
			},
			{
				Code: "4",
				Name: i18n.String{
					i18n.EN: "Standard rate (Island)",
					i18n.EL: "Κανονικός συντελεστής (Νησί)",
				},
			},
			{
				Code: "5",
				Name: i18n.String{
					i18n.EN: "Reduced rate (Island)",
					i18n.EL: "Μειωμένος συντελεστής (Νησί)",
				},
			},
			{
				Code: "6",
				Name: i18n.String{
					i18n.EN: "Super-reduced rate (Island)",
					i18n.EL: "Υπερμειωμένος συντελεστής (Νησί)",
				},
			},
			{
				Code: "7",
				Name: i18n.String{
					i18n.EN: "Without VAT",
					i18n.EL: "Άνευ ΦΠΑ",
				},
			},
			{
				Code: "8",
				Name: i18n.String{
					i18n.EN: "Records without VAT (e.g. Payroll, Amortisations)",
					i18n.EL: "Εγγραφές χωρίς ΦΠΑ (πχ Μισθοδοσία, Αποσβέσεις)",
				},
			},
		},
	},
	{
		Key: ExtKeyInvoiceType,
		Name: i18n.String{
			i18n.EN: "Invoice type",
			i18n.EL: "Είδος παραστατικού",
		},
		Values: []*cbc.Definition{
			{
				Code: "1.1",
				Name: i18n.String{
					i18n.EN: "Sales Invoice",
					i18n.EL: "Τιμολόγιο Πώλησης",
				},
			},
			{
				Code: "1.2",
				Name: i18n.String{
					i18n.EN: "Sales Invoice/Intra-community Supplies",
					i18n.EL: "Τιμολόγιο Πώλησης/Ενδοκοινοτικές Παραδόσεις",
				},
			},
			{
				Code: "1.3",
				Name: i18n.String{
					i18n.EN: "Sales Invoice/Third Country Supplies",
					i18n.EL: "Τιμολόγιο Πώλησης/Παραδόσεις Τρίτων Χωρών",
				},
			},
			{
				Code: "1.4",
				Name: i18n.String{
					i18n.EN: "Sales Invoice/Sale on Behalf of Third Parties",
					i18n.EL: "Τιμολόγιο Πώλησης/Πώληση για Λογαριασμό Τρίτων",
				},
			},
			{
				Code: "1.5",
				Name: i18n.String{
					i18n.EN: "Sales Invoice/Clearance of Sales on Behalf of Third Parties – Fees from Sales on Behalf of Third Parties",
					i18n.EL: "Τιμολόγιο Πώλησης/Εκκαθάριση Πωλήσεων Τρίτων - Αμοιβή από Πωλήσεις Τρίτων",
				},
			},
			{
				Code: "1.6",
				Name: i18n.String{
					i18n.EN: "Sales Invoice/Supplemental Accounting Source Document",
					i18n.EL: "Τιμολόγιο Πώλησης/Συμπληρωματικό Παραστατικό",
				},
			},
			{
				Code: "2.1",
				Name: i18n.String{
					i18n.EN: "Service Rendered Invoice",
					i18n.EL: "Τιμολόγιο Παροχής Υπηρεσιών",
				},
			},
			{
				Code: "2.2",
				Name: i18n.String{
					i18n.EN: "Intra-community Service Rendered Invoice",
					i18n.EL: "Τιμολόγιο Παροχής/Ενδοκοινοτική Παροχή Υπηρεσιών",
				},
			},
			{
				Code: "2.3",
				Name: i18n.String{
					i18n.EN: "Third Country Service Rendered Invoice",
					i18n.EL: "Τιμολόγιο Παροχής/Παροχή Υπηρεσιών σε λήπτη Τρίτης Χώρας",
				},
			},
			{
				Code: "2.4",
				Name: i18n.String{
					i18n.EN: "Service Rendered Invoice/Supplemental Accounting Source Document",
					i18n.EL: "Τιμολόγιο Παροχής/Συμπληρωματικό Παραστατικό",
				},
			},
			{
				Code: "3.1",
				Name: i18n.String{
					i18n.EN: "Proof of Expenditure (non-liable Issuer)",
					i18n.EL: "Τίτλος Κτήσης (μη υπόχρεος Εκδότης)",
				},
			},
			{
				Code: "3.2",
				Name: i18n.String{
					i18n.EN: "Proof of Expenditure (denial of issuance by liable Issuer)",
					i18n.EL: "Τίτλος Κτήσης (άρνηση έκδοσης από υπόχρεο Εκδότη)",
				},
			},
			{
				Code: "5.1",
				Name: i18n.String{
					i18n.EN: "Credit Invoice/Associated",
					i18n.EL: "Πιστωτικό Τιμολόγιο/Συσχετιζόμενο",
				},
			},
			{
				Code: "5.2",
				Name: i18n.String{
					i18n.EN: "Credit Invoice/Non-Associated",
					i18n.EL: "Πιστωτικό Τιμολόγιο/Μη Συσχετιζόμενο",
				},
			},
			{
				Code: "6.1",
				Name: i18n.String{
					i18n.EN: "Self-Delivery Record",
					i18n.EL: "Στοιχείο Αυτοπαράδοσης",
				},
			},
			{
				Code: "6.2",
				Name: i18n.String{
					i18n.EN: "Self-Supply Record",
					i18n.EL: "Στοιχείο Ιδιοχρησιμοποίησης",
				},
			},
			{
				Code: "7.1",
				Name: i18n.String{
					i18n.EN: "Contract – Income",
					i18n.EL: "Συμβόλαιο - Έσοδο",
				},
			},
			{
				Code: "8.1",
				Name: i18n.String{
					i18n.EN: "Rents – Income",
					i18n.EL: "Ενοίκια - Έσοδο",
				},
			},
			{
				Code: "8.2",
				Name: i18n.String{
					i18n.EN: "Special Record – Accommodation Tax Collection/Payment Receipt",
					i18n.EL: "Ειδικό Στοιχείο – Απόδειξης Είσπραξης Φόρου Διαμονής",
				},
			},
			{
				Code: "11.1",
				Name: i18n.String{
					i18n.EN: "Retail Sales Receipt",
					i18n.EL: "ΑΛΠ",
				},
			},
			{
				Code: "11.2",
				Name: i18n.String{
					i18n.EN: "Service Rendered Receipt",
					i18n.EL: "ΑΠΥ",
				},
			},
			{
				Code: "11.3",
				Name: i18n.String{
					i18n.EN: "Simplified Invoice",
					i18n.EL: "Απλοποιημένο Τιμολόγιο",
				},
			},
			{
				Code: "11.4",
				Name: i18n.String{
					i18n.EN: "Retail Sales Credit Note",
					i18n.EL: "Πιστωτικό Στοιχ. Λιανικής",
				},
			},
			{
				Code: "11.5",
				Name: i18n.String{
					i18n.EN: "Retail Sales Receipt on Behalf of Third Parties",
					i18n.EL: "Απόδειξη Λιανικής Πώλησης για Λογ/σμό Τρίτων",
				},
			},
			{
				Code: "13.1",
				Name: i18n.String{
					i18n.EN: "Expenses – Domestic/Foreign Retail Transaction Purchases",
					i18n.EL: "Έξοδα - Αγορές Λιανικών Συναλλαγών ημεδαπής / αλλοδαπής",
				},
			},
			{
				Code: "13.2",
				Name: i18n.String{
					i18n.EN: "Domestic/Foreign Retail Transaction Provision",
					i18n.EL: "Παροχή Λιανικών Συναλλαγών ημεδαπής / αλλοδαπής",
				},
			},
			{
				Code: "13.3",
				Name: i18n.String{
					i18n.EN: "Shared Utility Bills",
					i18n.EL: "Κοινόχρηστα",
				},
			},
			{
				Code: "13.4",
				Name: i18n.String{
					i18n.EN: "Subscriptions",
					i18n.EL: "Συνδρομές",
				},
			},
			{
				Code: "13.30",
				Name: i18n.String{
					i18n.EN: "Self-Declared Entity Accounting Source Documents (Dynamic)",
					i18n.EL: "Παραστατικά Οντότητας ως Αναγράφονται από την ίδια (Δυναμικό)",
				},
			},
			{
				Code: "13.31",
				Name: i18n.String{
					i18n.EN: "Domestic/Foreign Retail Sales Credit Note",
					i18n.EL: "Πιστωτικό Στοιχ. Λιανικής ημεδαπής / αλλοδαπής",
				},
			},
			{
				Code: "14.1",
				Name: i18n.String{
					i18n.EN: "Invoice/Intra-community Acquisitions",
					i18n.EL: "Τιμολόγιο / Ενδοκοινοτικές Αποκτήσεις",
				},
			},
			{
				Code: "14.2",
				Name: i18n.String{
					i18n.EN: "Invoice/Third Country Acquisitions",
					i18n.EL: "Τιμολόγιο / Αποκτήσεις Τρίτων Χωρών",
				},
			},
			{
				Code: "14.3",
				Name: i18n.String{
					i18n.EN: "Invoice/Intra-community Services Receipt",
					i18n.EL: "Τιμολόγιο / Ενδοκοινοτική Λήψη Υπηρεσιών",
				},
			},
			{
				Code: "14.4",
				Name: i18n.String{
					i18n.EN: "Invoice/Third Country Services Receipt",
					i18n.EL: "Τιμολόγιο / Λήψη Υπηρεσιών Τρίτων Χωρών",
				},
			},
			{
				Code: "14.5",
				Name: i18n.String{
					i18n.EN: "EFKA",
					i18n.EL: "ΕΦΚΑ και λοιποί Ασφαλιστικοί Οργανισμοί",
				},
			},
			{
				Code: "14.30",
				Name: i18n.String{
					i18n.EN: "Self-Declared Entity Accounting Source Documents (Dynamic)",
					i18n.EL: "Παραστατικά Οντότητας ως Αναγράφονται από την ίδια (Δυναμικό)",
				},
			},
			{
				Code: "14.31",
				Name: i18n.String{
					i18n.EN: "Domestic/Foreign Credit Note",
					i18n.EL: "Πιστωτικό ημεδαπής / αλλοδαπής",
				},
			},
			{
				Code: "15.1",
				Name: i18n.String{
					i18n.EN: "Contract-Expense",
					i18n.EL: "Συμβόλαιο - Έξοδο",
				},
			},
			{
				Code: "16.1",
				Name: i18n.String{
					i18n.EN: "Rent-Expense",
					i18n.EL: "Ενοίκιο Έξοδο",
				},
			},
			{
				Code: "17.1",
				Name: i18n.String{
					i18n.EN: "Payroll",
					i18n.EL: "Μισθοδοσία",
				},
			},
			{
				Code: "17.2",
				Name: i18n.String{
					i18n.EN: "Amortisations",
					i18n.EL: "Αποσβέσεις",
				},
			},
			{
				Code: "17.3",
				Name: i18n.String{
					i18n.EN: "Other Income Adjustment/Regularisation Entries – Accounting Base",
					i18n.EL: "Λοιπές Εγγραφές Τακτοποίησης Εσόδων - Λογιστική Βάση",
				},
			},
			{
				Code: "17.4",
				Name: i18n.String{
					i18n.EN: "Other Income Adjustment/Regularisation Entries – Tax Base",
					i18n.EL: "Λοιπές Εγγραφές Τακτοποίησης Εσόδων - Φορολογική Βάση",
				},
			},
			{
				Code: "17.5",
				Name: i18n.String{
					i18n.EN: "Other Expense Adjustment/Regularisation Entries – Accounting Base",
					i18n.EL: "Λοιπές Εγγραφές Τακτοποίησης Εξόδων - Λογιστική Βάση",
				},
			},
			{
				Code: "17.6",
				Name: i18n.String{
					i18n.EN: "Other Expense Adjustment/Regularisation Entries – Tax Base",
					i18n.EL: "Λοιπές Εγγραφές Τακτοποίησης Εξόδων - Φορολογική Βάση",
				},
			},
		},
	},
	{
		Key: ExtKeyPaymentMeans,
		Name: i18n.String{
			i18n.EN: "Payment means",
			i18n.EL: "Τρόπος Πληρωμής",
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Domestic Payments Account Number",
					i18n.EL: "Επαγ. Λογαριασμός Πληρωμών Ημεδαπής",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Foreign Payments Account Number",
					i18n.EL: "Επαγ. Λογαριασμός Πληρωμών Αλλοδαπής",
				},
			},
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "Cash",
					i18n.EL: "Μετρητά",
				},
			},
			{
				Code: "4",
				Name: i18n.String{
					i18n.EN: "Check",
					i18n.EL: "Επιταγή",
				},
			},
			{
				Code: "5",
				Name: i18n.String{
					i18n.EN: "On credit",
					i18n.EL: "Επί Πιστώσει",
				},
			},
			{
				Code: "6",
				Name: i18n.String{
					i18n.EN: "Web Banking",
					i18n.EL: "Web Banking",
				},
			},
			{
				Code: "7",
				Name: i18n.String{
					i18n.EN: "POS / e-POS",
					i18n.EL: "POS / e-POS",
				},
			},
		},
	},
	{
		Key: ExtKeyExemption,
		Name: i18n.String{
			i18n.EN: "VAT exemption cause",
			i18n.EL: "Κατηγορία Αιτίας Εξαίρεσης ΦΠΑ",
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 3 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 3 του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 5 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 5 του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 13 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 13 του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "4",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 14 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 14 του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "5",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 16 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 16 του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "6",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 19 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 19 του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "7",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 22 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 22 του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "8",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 24 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 24 του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "9",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 25 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 25 του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "10",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 26 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 26 του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "11",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 27 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 27 του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "12",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 27 - Seagoing Vessels of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 27 - Πλοία Ανοικτής Θαλάσσης του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "13",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 27.1.γ - Seagoing Vessels of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 27.1.γ - Πλοία Ανοικτής Θαλάσσης του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "14",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 28 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 28 του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "15",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 39 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 39 του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "16",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 39a of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 39α του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "17",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 40 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 40 του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "18",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 41 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 41 του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "19",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 47 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 47 του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "20",
				Name: i18n.String{
					i18n.EN: "VAT included - article 43 of the VAT code",
					i18n.EL: "ΦΠΑ εμπεριεχόμενος - άρθρο 43 του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "21",
				Name: i18n.String{
					i18n.EN: "VAT included - article 44 of the VAT code",
					i18n.EL: "ΦΠΑ εμπεριεχόμενος - άρθρο 44 του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "22",
				Name: i18n.String{
					i18n.EN: "VAT included - article 45 of the VAT code",
					i18n.EL: "ΦΠΑ εμπεριεχόμενος - άρθρο 45 του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "23",
				Name: i18n.String{
					i18n.EN: "VAT included - article 46 of the VAT code",
					i18n.EL: "ΦΠΑ εμπεριεχόμενος - άρθρο 46 του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "24",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 6 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 6 του Κώδικα ΦΠΑ",
				},
			},
			{
				Code: "25",
				Name: i18n.String{
					i18n.EN: "Without VAT - ΠΟΛ.1029/1995",
					i18n.EL: "Χωρίς ΦΠΑ - ΠΟΛ.1029/1995",
				},
			},
			{
				Code: "26",
				Name: i18n.String{
					i18n.EN: "Without VAT - ΠΟΛ.1167/2015",
					i18n.EL: "Χωρίς ΦΠΑ - ΠΟΛ.1167/2015",
				},
			},
			{
				Code: "27",
				Name: i18n.String{
					i18n.EN: "Without VAT - Other VAT exceptions",
					i18n.EL: "Λοιπές Εξαιρέσεις ΦΠΑ",
				},
			},
			{
				Code: "28",
				Name: i18n.String{
					i18n.EN: "Without VAT - Article 24 (b) (1) of the VAT Code (Tax Free)",
					i18n.EL: "Χωρίς ΦΠΑ – άρθρο 24 περ. β' παρ.1 του Κώδικα ΦΠΑ, (Tax Free)",
				},
			},
			{
				Code: "29",
				Name: i18n.String{
					i18n.EN: "Without VAT - Article 47b of the VAT Code (OSS non-EU scheme)",
					i18n.EL: "Χωρίς ΦΠΑ – άρθρο 47β, του Κώδικα ΦΠΑ (OSS μη ενωσιακό καθεστώς)",
				},
			},
			{
				Code: "30",
				Name: i18n.String{
					i18n.EN: "Without VAT - Article 47c of the VAT Code (OSS EU scheme)",
					i18n.EL: "Χωρίς ΦΠΑ – άρθρο 47γ, του Κώδικα ΦΠΑ (OSS ενωσιακό καθεστώς)",
				},
			},
			{
				Code: "31",
				Name: i18n.String{
					i18n.EN: "Excluding VAT - Article 47d of the VAT Code (IOSS)",
					i18n.EL: "Χωρίς ΦΠΑ – άρθρο 47δ του Κώδικα ΦΠΑ (IOSS)",
				},
			},
		},
	},
	{
		Key: ExtKeyIncomeCat,
		Name: i18n.String{
			i18n.EN: "Income Classification Category",
			i18n.EL: "Κωδικός Κατηγορίας Χαρακτηρισμού Εσόδων",
		},
		Values: []*cbc.Definition{
			{
				Code: "category1_1",
				Name: i18n.String{
					i18n.EN: "Commodity Sale Income (+)/(-)",
					i18n.EL: "Έσοδα από Πώληση Εμπορευμάτων (+)/(-)",
				},
			},
			{
				Code: "category1_2",
				Name: i18n.String{
					i18n.EN: "Product Sale Income (+)/(-)",
					i18n.EL: "Έσοδα από Πώληση Προϊόντων (+)/(-)",
				},
			},
			{
				Code: "category1_3",
				Name: i18n.String{
					i18n.EN: "Provision of Services Income (+)/(-)",
					i18n.EL: "Έσοδα από Παροχή Υπηρεσιών (+)/(-)",
				},
			},
			{
				Code: "category1_4",
				Name: i18n.String{
					i18n.EN: "Sale of Fixed Assets Income (+)/(-)",
					i18n.EL: "Έσοδα από Πώληση Παγίων (+)/(-)",
				},
			},
			{
				Code: "category1_5",
				Name: i18n.String{
					i18n.EN: "Other Income/Profits (+)/(-)",
					i18n.EL: "Λοιπά Έσοδα/ Κέρδη (+)/(-)",
				},
			},
			{
				Code: "category1_6",
				Name: i18n.String{
					i18n.EN: "Self-Deliveries/Self-Supplies (+)/(-)",
					i18n.EL: "Αυτοπαραδόσεις / Ιδιοχρησιμοποιήσεις (+)/(-)",
				},
			},
			{
				Code: "category1_7",
				Name: i18n.String{
					i18n.EN: "Income on behalf of Third Parties (+)/(-)",
					i18n.EL: "Έσοδα για λ/σμο τρίτων (+)/(-)",
				},
			},
			{
				Code: "category1_8",
				Name: i18n.String{
					i18n.EN: "Past fiscal years income (+)/(-)",
					i18n.EL: "Έσοδα προηγούμενων χρήσεων (+)/ (-)",
				},
			},
			{
				Code: "category1_9",
				Name: i18n.String{
					i18n.EN: "Future fiscal years income (+)/(-)",
					i18n.EL: "Έσοδα επομένων χρήσεων (+)/(-)",
				},
			},
			{
				Code: "category1_10",
				Name: i18n.String{
					i18n.EN: "Other Income Adjustment/Regularisation Entries (+)/(-)",
					i18n.EL: "Λοιπές Εγγραφές Τακτοποίησης Εσόδων (+)/(-)",
				},
			},
			{
				Code: "category1_95",
				Name: i18n.String{
					i18n.EN: "Other Income-related Information (+)/(-)",
					i18n.EL: "Λοιπά Πληροφοριακά Στοιχεία Εσόδων (+)/(-)",
				},
			},
		},
	},
	{
		Key: ExtKeyIncomeType,
		Name: i18n.String{
			i18n.EN: "Income Classification Type",
			i18n.EL: "Κωδικός Τύπου Χαρακτηρισμού Εσόδων",
		},
		Values: []*cbc.Definition{
			{
				Code: "E3_106",
				Name: i18n.String{
					i18n.EN: "Self-Production of Fixed Assets – Self-Deliveries – Destroying inventory/Commodities",
					i18n.EL: "Ιδιοπαραγωγή παγίων - Αυτοπαραδόσεις - Καταστροφές αποθεμάτων/Εμπορεύματα",
				},
			},
			{
				Code: "E3_205",
				Name: i18n.String{
					i18n.EN: "Self-Production of Fixed Assets – Self-Deliveries – Destroying inventory/Raw and other materials",
					i18n.EL: "Ιδιοπαραγωγή παγίων - Αυτοπαραδόσεις - Καταστροφές αποθεμάτων/Πρώτες ύλες και λοιπά υλικά",
				},
			},
			{
				Code: "E3_210",
				Name: i18n.String{
					i18n.EN: "Self-Production of Fixed Assets – Self-Deliveries – Destroying inventory/Products and production in progress",
					i18n.EL: "Ιδιοπαραγωγή παγίων - Αυτοπαραδόσεις - Καταστροφές αποθεμάτων/Προϊόντα και παραγωγή σε εξέλιξη",
				},
			},
			{
				Code: "E3_305",
				Name: i18n.String{
					i18n.EN: "Self-Production of Fixed Assets – Self-Deliveries – Destroying inventory/Raw and other materials",
					i18n.EL: "Ιδιοπαραγωγή παγίων - Αυτοπαραδόσεις – Καταστροφές αποθεμάτων/Πρώτες ύλες και λοιπά υλικά",
				},
			},
			{
				Code: "E3_310",
				Name: i18n.String{
					i18n.EN: "Self-Production of Fixed Assets – Self-Deliveries – Destroying inventory/Products and production in progress",
					i18n.EL: "Ιδιοπαραγωγή παγίων - Αυτοπαραδόσεις - Καταστροφές αποθεμάτων/Προϊόντα και παραγωγή σε εξέλιξη",
				},
			},
			{
				Code: "E3_318",
				Name: i18n.String{
					i18n.EN: "Self-Production of Fixed Assets – Self-Deliveries – Destroying inventory/Production expenses",
					i18n.EL: "Ιδιοπαραγωγή παγίων - Αυτοπαραδόσεις - Καταστροφές αποθεμάτων/Έξοδα παραγωγής",
				},
			},
			{
				Code: "E3_561_001",
				Name: i18n.String{
					i18n.EN: "Wholesale Sales of Goods and Services – for Traders",
					i18n.EL: "Πωλήσεις αγαθών και υπηρεσιών Χονδρικές - Επιτηδευματιών",
				},
			},
			{
				Code: "E3_561_002",
				Name: i18n.String{
					i18n.EN: "Wholesale Sales of Goods and Services pursuant to article 39a paragraph 5 of the VAT Code (Law 2859/2000)",
					i18n.EL: "Πωλήσεις αγαθών και υπηρεσιών Χονδρικές βάσει άρθρου 39α παρ 5 του Κώδικα Φ.Π.Α. (Ν.2859/2000)",
				},
			},
			{
				Code: "E3_561_003",
				Name: i18n.String{
					i18n.EN: "Retail Sales of Goods and Services – Private Clientele",
					i18n.EL: "Πωλήσεις αγαθών και υπηρεσιών Λιανικές - Ιδιωτική Πελατεία",
				},
			},
			{
				Code: "E3_561_004",
				Name: i18n.String{
					i18n.EN: "Retail Sales of Goods and Services pursuant to article 39a paragraph 5 of the VAT Code (Law 2859/2000)",
					i18n.EL: "Πωλήσεις αγαθών και υπηρεσιών Λιανικές βάσει άρθρου 39α παρ 5 του Κώδικα Φ.Π.Α. (Ν.2859/2000)",
				},
			},
			{
				Code: "E3_561_005",
				Name: i18n.String{
					i18n.EN: "Intra-Community Foreign Sales of Goods and Services",
					i18n.EL: "Πωλήσεις αγαθών και υπηρεσιών Εξωτερικού Ενδοκοινοτικές",
				},
			},
			{
				Code: "E3_561_006",
				Name: i18n.String{
					i18n.EN: "Third Country Foreign Sales of Goods and Services",
					i18n.EL: "Πωλήσεις αγαθών και υπηρεσιών Εξωτερικού Τρίτες Χώρες",
				},
			},
			{
				Code: "E3_561_007",
				Name: i18n.String{
					i18n.EN: "Other Sales of Goods and Services",
					i18n.EL: "Πωλήσεις αγαθών και υπηρεσιών Λοιπά",
				},
			},
			{
				Code: "E3_562",
				Name: i18n.String{
					i18n.EN: "Other Ordinary Income",
					i18n.EL: "Λοιπά συνήθη έσοδα",
				},
			},
			{
				Code: "E3_563",
				Name: i18n.String{
					i18n.EN: "Credit Interest and Related Income",
					i18n.EL: "Πιστωτικοί τόκοι και συναφή έσοδα",
				},
			},
			{
				Code: "E3_564",
				Name: i18n.String{
					i18n.EN: "Credit Exchange Differences",
					i18n.EL: "Πιστωτικές συναλλαγματικές διαφορές",
				},
			},
			{
				Code: "E3_565",
				Name: i18n.String{
					i18n.EN: "Income from Participations",
					i18n.EL: "Έσοδα συμμετοχών",
				},
			},
			{
				Code: "E3_566",
				Name: i18n.String{
					i18n.EN: "Profits from Disposing Non-Current Assets",
					i18n.EL: "Κέρδη από διάθεση μη κυκλοφορούντων περιουσιακών στοιχείων",
				},
			},
			{
				Code: "E3_567",
				Name: i18n.String{
					i18n.EN: "Profits from the Reversal of Provisions and Impairments",
					i18n.EL: "Κέρδη από αναστροφή προβλέψεων και απομειώσεων",
				},
			},
			{
				Code: "E3_568",
				Name: i18n.String{
					i18n.EN: "Profits from Measurement at Fair Value",
					i18n.EL: "Κέρδη από επιμέτρηση στην εύλογη αξία",
				},
			},
			{
				Code: "E3_570",
				Name: i18n.String{
					i18n.EN: "Extraordinary income and profits",
					i18n.EL: "Ασυνήθη έσοδα και κέρδη",
				},
			},
			{
				Code: "E3_595",
				Name: i18n.String{
					i18n.EN: "Self-Production Expenses",
					i18n.EL: "Έξοδα σε ιδιοπαραγωγή",
				},
			},
			{
				Code: "E3_596",
				Name: i18n.String{
					i18n.EN: "Subsidies - Grants",
					i18n.EL: "Επιδοτήσεις - Επιχορηγήσεις",
				},
			},
			{
				Code: "E3_597",
				Name: i18n.String{
					i18n.EN: "Subsidies – Grants for Investment Purposes – Expense Coverage",
					i18n.EL: "Επιδοτήσεις - Επιχορηγήσεις για επενδυτικούς σκοπούς - κάλυψη δαπανών",
				},
			},
			{
				Code: "E3_880_001",
				Name: i18n.String{
					i18n.EN: "Wholesale Sales of Fixed Assets",
					i18n.EL: "Πωλήσεις Παγίων Χονδρικές",
				},
			},
			{
				Code: "E3_880_002",
				Name: i18n.String{
					i18n.EN: "Retail Sales of Fixed Assets",
					i18n.EL: "Πωλήσεις Παγίων Λιανικές",
				},
			},
			{
				Code: "E3_880_003",
				Name: i18n.String{
					i18n.EN: "Intra-Community Foreign Sales of Fixed Assets",
					i18n.EL: "Πωλήσεις Παγίων Εξωτερικού Ενδοκοινοτικές",
				},
			},
			{
				Code: "E3_880_004",
				Name: i18n.String{
					i18n.EN: "Third Country Foreign Sales of Fixed Assets",
					i18n.EL: "Πωλήσεις Παγίων Εξωτερικού Τρίτες Χώρες",
				},
			},
			{
				Code: "E3_881_001",
				Name: i18n.String{
					i18n.EN: "Wholesale Sales on behalf of Third Parties",
					i18n.EL: "Πωλήσεις για λογ/σμο Τρίτων Χονδρικές",
				},
			},
			{
				Code: "E3_881_002",
				Name: i18n.String{
					i18n.EN: "Retail Sales on behalf of Third Parties",
					i18n.EL: "Πωλήσεις για λογ/σμο Τρίτων Λιανικές",
				},
			},
			{
				Code: "E3_881_003",
				Name: i18n.String{
					i18n.EN: "Intra-Community Foreign Sales on behalf of Third Parties",
					i18n.EL: "Πωλήσεις για λογ/σμο Τρίτων Εξωτερικού Ενδοκοινοτικές",
				},
			},
			{
				Code: "E3_881_004",
				Name: i18n.String{
					i18n.EN: "Third Country Foreign Sales on behalf of Third Parties",
					i18n.EL: "Πωλήσεις για λογ/σμο Τρίτων Εξωτερικού Τρίτες Χώρες",
				},
			},
			{
				Code: "E3_598_001",
				Name: i18n.String{
					i18n.EN: "Sales of goods belonging to excise duty",
					i18n.EL: "Πωλήσεις αγαθών που υπάγονται σε ΕΦΚ",
				},
			},
			{
				Code: "E3_598_003",
				Name: i18n.String{
					i18n.EN: "Sales on behalf of farmers through an agricultural cooperative e.t.c.",
					i18n.EL: "Πωλήσεις για λογαριασμό αγροτών μέσω αγροτικού συνεταιρισμού κ.λ.π.",
				},
			},
		},
	},
}

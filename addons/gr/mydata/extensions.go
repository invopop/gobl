package mydata

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Regime extension codes.
const (
	ExtKeyVATRate      = "gr-mydata-vat-rate"
	ExtKeyInvoiceType  = "gr-mydata-invoice-type"
	ExtKeyExemption    = "gr-mydata-exemption"
	ExtKeyIncomeCat    = "gr-mydata-income-cat"
	ExtKeyIncomeType   = "gr-mydata-income-type"
	ExtKeyPaymentMeans = "gr-mydata-payment-means"
	ExtKeyTaxType      = "gr-mydata-tax-type"
	ExtKeyOtherTax     = "gr-mydata-other-tax"
	ExtKeyFee          = "gr-mydata-fee"
	ExtKeyStampDuty    = "gr-mydata-stamp-duty"
)

// Tax type codes.
const (
	TaxTypeWithholding = "1"
	TaxTypeFee         = "2"
	TaxTypeOtherTax    = "3"
	TaxTypeStampDuty   = "4"
	TaxTypeDeduction   = "5"
)

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyVATRate,
		Name: i18n.String{
			i18n.EN: "VAT rate",
			i18n.EL: "Κατηγορία ΦΠΑ",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Greece has three VAT rates: standard, reduced and super-reduced. Each of these rates are reduced by
				30% on the islands of Leros, Lesbos, Kos, Samos and Chios. The tax authority identifies each rate
				with a specific VAT category.

				The IAPR VAT category code must be set using the ~gr-mydata-vat-rate~ extension of
				a line's tax to one of the codes.

				| Code | Description                 | GOBL Rate              |
				| ---- | --------------------------- | ---------------------- |
				| ~1~  | Standard rate               | ~standard~             |
				| ~2~  | Reduced rate                | ~reduced~              |
				| ~3~  | Super-reduced rate          | ~super-reduced~        |
				| ~4~  | Standard rate (Island)      | ~standard+island~      |
				| ~5~  | Reduced rate (Island)       | ~reduced+island~       |
				| ~6~  | Super-reduced rate (Island) | ~super-reduced+island~ |
				| ~7~  | Without VAT                 | ~exempt~               |
				| ~8~  | Records without VAT         |                        |

				Please, note that GOBL will automatically set the proper ~gr-mydata-vat-rate~ code and tax percent automatically when the line tax uses any of the GOBL rates specified in the table above. For example:

				~~~js
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					// ...
					"lines": [
						{
							"i": 1,
							"quantity": "20",
							"item": {
								"name": "Υπηρεσίες Ανάπτυξης",
								"price": "90.00",
							},
							"sum": "1800.00",
							"taxes": [
								{
									"cat": "VAT",
									"rate": "standard+island"
								}
							],
							"total": "1800.00"
						}
					],
				}
				~~~
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Technical description of REST API interfaces for sending & receiving data for ERP users (Version 2.0.0)",
					i18n.EL: "Τεχνική περιγραφή διεπαφών REST API για διαβίβαση & λήψη δεδομένων για χρήστες ERP (Έκδοση 2.0.0)",
				},
				URL: "https://www.aade.gr/sites/default/files/2025-12/myDATA%20API%20Documentation%20v2.0.0_official_erp.pdf",
			},
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
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The Greek tax authority (IAPR) requires an invoice type code to be specified as part of the invoice. GOBL will
				automatically set the correct code based on the invoice's ~type~ and ~$tags~ values.

				However, you can also set the code manually using the ~gr-mydata-invoice-type~ extension in the tax
				section of the invoice, and setting the invoice's ~type~ to ~other~.

				The following table lists how the combination of ~type~ and ~$tags~ values are mapped to the
				IAPR MyDATA invoice type code:

				| Type   | Description                                     | GOBL Type     | GOBL Tags                  |
				| ------ | ----------------------------------------------- | ------------- |----------------------------|
				| ~1.1~  | Sales Invoice                                   | ~standard~    | ~goods~                    |
				| ~1.2~  | Sales Invoice/Intra-community Supplies          | ~standard~    | ~goods~, ~export~, ~eu~    |
				| ~1.3~  | Sales Invoice/Third Country Supplies            | ~standard~    | ~goods~, ~export~          |
				| ~1.4~  | Sales Invoice/Sale on Behalf of Third Parties   | ~standard~    | ~goods~, ~self-billed~     |
				| ~2.1~  | Service Rendered Invoice                        | ~standard~    | ~services~                 |
				| ~2.2~  | Intra-community Service Rendered Invoice        | ~standard~    | ~services~, ~export~, ~eu~ |
				| ~2.3~  | Third Country Service Rendered Invoice          | ~standard~    | ~services~, ~export~       |
				| ~5.1~  | Credit Invoice/Associated                       | ~credit-note~ |                            |
				| ~11.1~ | Retail Sales Receipt                            | ~standard~    | ~goods~, ~simplified~      |
				| ~11.2~ | Service Rendered Receipt                        | ~standard~    | ~services~, ~simplified~   |
				| ~11.3~ | Simplified Invoice                              | ~standard~    | ~simplified~               |
				| ~11.4~ | Retail Sales Credit Note                        | ~credit-note~ | ~simplified~               |
				| ~11.5~ | Retail Sales Receipt on Behalf of Third Parties | ~credit-note~ | ~goods~, ~simplified~, ~self-billed~ |

				For example, this is how you set the IAPR invoice type explicitly:

				~~~json
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					// ...
					"type": "other",
					"tax": {
						"ext": {
							"gr-mydata-invoice-type": "2.3"
						}
					}
				}
				~~~

				And this is how you'll get the same result by using the GOBL type and tags:

				~~~json
				{
					"$schema": "https://gobl.org/draft-0/bill/invoice",
					"$addons": ["gr-mydata-v1"],
					"$tags": ["services", "export"],
					// ...
					"type": "standard",
				}
				~~~
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Technical description of REST API interfaces for sending & receiving data for ERP users (Version 2.0.0)",
					i18n.EL: "Τεχνική περιγραφή διεπαφών REST API για διαβίβαση & λήψη δεδομένων για χρήστες ERP (Έκδοση 2.0.0)",
				},
				URL: "https://www.aade.gr/sites/default/files/2025-12/myDATA%20API%20Documentation%20v2.0.0_official_erp.pdf",
			},
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
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The IAPR requires invoices to specify a payment method code. In a GOBL invoice,
				the payment means is set using the ~key~ field in the payment instructions.
				The following table lists all the IAPR payment methods and how GOBL will map from
				the payment instructions key to each of them:

				| Code | Name                             | GOBL Payment Instruction Key |
				| ---- | -------------------------------- | ---------------------------- |
				| ~1~  | Domestic Payments Account Number | ~credit-transfer~            |
				| ~2~  | Foreign Payments Account Number  | ~credit-transfer+foreign~    |
				| ~3~  | Cash                             | ~cash~                       |
				| ~4~  | Check                            | ~cheque~                     |
				| ~5~  | On credit                        | ~promissory-note~            |
				| ~6~  | Web Banking                      | ~online~                     |
				| ~7~  | POS / e-POS                      | ~card~                       |

				For example:

				~~~js
				"payment": {
					"instructions": {
						"key": "credit-transfer+foreign" // Will set the IAPR Payment Method to "2"
					}
				}
				~~~
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Technical description of REST API interfaces for sending & receiving data for ERP users (Version 2.0.0)",
					i18n.EL: "Τεχνική περιγραφή διεπαφών REST API για διαβίβαση & λήψη δεδομένων για χρήστες ERP (Έκδοση 2.0.0)",
				},
				URL: "https://www.aade.gr/sites/default/files/2025-12/myDATA%20API%20Documentation%20v2.0.0_official_erp.pdf",
			},
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
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Greece invoices can be exempt of VAT for different causes and the tax authority
				require a specific cause code to be provided.

				GOBL will try to automatically map the tax combo ~key~ to an appropriate exemption
				code, but given the number of options you will no doubt need to consult with
				a tax professional to determine the correct code for your case.

				Example usage:

				~~~js
				"lines": [
					{
						"i": 1,
						"quantity": "20",
						"item": {
							"name": "Υπηρεσίες Ανάπτυξης",
							"price": "90.00",
						},
						"sum": "1800.00",
						"taxes": [
							{
								"cat": "VAT",
								"key": "exempt",
								"ext": {
									"gr-mydata-exemption": "7"
								}
							}
						],
						"total": "1800.00"
					}
				]
				~~~

				Note that descriptions provided in the ~gr-mydata-exemption~ extension are
				generated using AI and may not be 100% accurate. Always consult with a tax
				professional to ensure you're choosing the correct code.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Technical description of REST API interfaces for sending & receiving data for ERP users (Version 2.0.0)",
					i18n.EL: "Τεχνική περιγραφή διεπαφών REST API για διαβίβαση & λήψη δεδομένων για χρήστες ERP (Έκδοση 2.0.0)",
				},
				URL: "https://www.aade.gr/sites/default/files/2025-12/myDATA%20API%20Documentation%20v2.0.0_official_erp.pdf",
			},
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 3 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 3 του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Outside the scope of VAT per Art. 3 (transaction not within the VAT “object”)."),
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 5 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 5 του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Outside Greek VAT due to the legal cases in Art. 5 (e.g., not treated as a taxable supply under that article)."),
			},
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 13 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 13 του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Use when the transaction is exempt under Art. 13 (as defined in Greek VAT Code)."),
			},
			{
				Code: "4",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 14 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 14 του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Use when Art. 14 grants an exemption (specific cases set in law)."),
			},
			{
				Code: "5",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 16 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 16 του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Use for the specific non-taxable/exempt cases referenced in Art. 16."),
			},
			{
				Code: "6",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 19 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 19 του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Exempt cases per Art. 19 (e.g., special valuation/administrative rules that make the item non-taxable)."),
			},
			{
				Code: "7",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 22 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 22 του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Classic VAT-exempt activities (healthcare, education, culture, financial/insurance, rentals, etc.)."),
			},
			{
				Code: "8",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 24 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 24 του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Exports of goods outside the EU (exempt with right to deduct)."),
			},
			{
				Code: "9",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 25 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 25 του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Free zones/warehousing / international transport‐related exemptions covered by Art. 25."),
			},
			{
				Code: "10",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 26 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 26 του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Exemption per Art. 26 (misc. cases defined in law)."),
			},
			{
				Code: "11",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 27 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 27 του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Supplies relating to vessels/aircraft and similar Art. 27 exemptions."),
			},
			{
				Code: "12",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 27 - Seagoing Vessels of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 27 - Πλοία Ανοικτής Θαλάσσης του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Supplies for qualifying seagoing vessels (navigation on high seas) – exempt with right to deduct."),
			},
			{
				Code: "13",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 27.1.γ - Seagoing Vessels of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 27.1.γ - Πλοία Ανοικτής Θαλάσσης του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Specific seagoing-vessel case under Art. 27(1)(c)."),
			},
			{
				Code: "14",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 28 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 28 του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Intra-Community supply of goods (B2B, goods dispatched to another EU MS)."),
			},
			{
				Code: "15",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 39 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 39 του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Small-business exemption regime."),
			},
			{
				Code: "16",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 39a of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 39α του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Domestic reverse charge (e.g., mobiles/tablets/consoles, scrap, etc.)."),
			},
			{
				Code: "17",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 40 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 40 του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Flat-rate (lump-sum) special scheme: no output VAT charged."),
			},
			{
				Code: "18",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 41 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 41 του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Farmers’ special scheme (no VAT charged; compensation scheme applies)."),
			},
			{
				Code: "19",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 47 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 47 του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Investment gold special regime (generally exempt, option to tax possible)."),
			},
			{
				Code: "20",
				Name: i18n.String{
					i18n.EN: "VAT included - article 43 of the VAT code",
					i18n.EL: "ΦΠΑ εμπεριεχόμενος - άρθρο 43 του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Travel Agents’ Margin Scheme (TOMS) – VAT is on the margin."),
			},
			{
				Code: "21",
				Name: i18n.String{
					i18n.EN: "VAT included - article 44 of the VAT code",
					i18n.EL: "ΦΠΑ εμπεριεχόμενος - άρθρο 44 του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Margin scheme – second-hand goods."),
			},
			{
				Code: "22",
				Name: i18n.String{
					i18n.EN: "VAT included - article 45 of the VAT code",
					i18n.EL: "ΦΠΑ εμπεριεχόμενος - άρθρο 45 του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Margin scheme – works of art."),
			},
			{
				Code: "23",
				Name: i18n.String{
					i18n.EN: "VAT included - article 46 of the VAT code",
					i18n.EL: "ΦΠΑ εμπεριεχόμενος - άρθρο 46 του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Margin scheme – collectors’ items & antiques."),
			},
			{
				Code: "24",
				Name: i18n.String{
					i18n.EN: "Without VAT - article 6 of the VAT code",
					i18n.EL: "Χωρίς ΦΠΑ - άρθρο 6 του Κώδικα ΦΠΑ",
				},
				Desc: i18n.NewString("Out-of-scope/“excluding VAT” per Art. 6 (place-of-supply style exclusions)."),
			},
			{
				Code: "25",
				Name: i18n.String{
					i18n.EN: "Without VAT - ΠΟΛ.1029/1995",
					i18n.EL: "Χωρίς ΦΠΑ - ΠΟΛ.1029/1995",
				},
				Desc: i18n.NewString("Exemption per Circular ΠΟΛ.1029/1995 (narrow administrative exemption)."),
			},
			{
				Code: "26",
				Name: i18n.String{
					i18n.EN: "Without VAT - ΠΟΛ.1167/2015",
					i18n.EL: "Χωρίς ΦΠΑ - ΠΟΛ.1167/2015",
				},
				Desc: i18n.NewString("Exemption per Circular ΠΟΛ.1167/2015 (narrow administrative exemption)."),
			},
			{
				Code: "27",
				Name: i18n.String{
					i18n.EN: "Without VAT - Other VAT exceptions",
					i18n.EL: "Λοιπές Εξαιρέσεις ΦΠΑ",
				},
				Desc: i18n.NewString("Catch-all for other statutory exemptions not listed above."),
			},
			{
				Code: "28",
				Name: i18n.String{
					i18n.EN: "Without VAT - Article 24 (b) (1) of the VAT Code (Tax Free)",
					i18n.EL: "Χωρίς ΦΠΑ – άρθρο 24 περ. β' παρ.1 του Κώδικα ΦΠΑ, (Tax Free)",
				},
				Desc: i18n.NewString("Retail “tax-free shopping” to third-country travelers."),
			},
			{
				Code: "29",
				Name: i18n.String{
					i18n.EN: "Without VAT - Article 47b of the VAT Code (OSS non-EU scheme)",
					i18n.EL: "Χωρίς ΦΠΑ – άρθρο 47β, του Κώδικα ΦΠΑ (OSS μη ενωσιακό καθεστώς)",
				},
				Desc: i18n.NewString("Non-Union OSS services to EU consumers — VAT due in MS of consumption, not Greece."),
			},
			{
				Code: "30",
				Name: i18n.String{
					i18n.EN: "Without VAT - Article 47c of the VAT Code (OSS EU scheme)",
					i18n.EL: "Χωρίς ΦΠΑ – άρθρο 47γ, του Κώδικα ΦΠΑ (OSS ενωσιακό καθεστώς)",
				},
				Desc: i18n.NewString("Union OSS (intra-EU B2C distance sales) — VAT due in MS of consumption."),
			},
			{
				Code: "31",
				Name: i18n.String{
					i18n.EN: "Excluding VAT - Article 47d of the VAT Code (IOSS)",
					i18n.EL: "Χωρίς ΦΠΑ – άρθρο 47δ του Κώδικα ΦΠΑ (IOSS)",
				},
				Desc: i18n.NewString("IOSS for low-value imports — VAT accounted via IOSS in MS of consumption."),
			},
		},
	},
	{
		Key: ExtKeyIncomeCat,
		Name: i18n.String{
			i18n.EN: "Income Classification Category",
			i18n.EL: "Κωδικός Κατηγορίας Χαρακτηρισμού Εσόδων",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Invoices reported to the Greek tax authority via myDATA can optionally include information
				about the income classification of each invoice item.

				In a GOBL invoice, the ~gr-mydata-income-cat~ and ~gr-mydata-income-type~ extensions can be
				set at the item level to any of the values expected by the IAPR. For example:

				~~~json
				"lines": [
					{
						"i": 1,
						"quantity": "20",
						"item": {
							"name": "Υπηρεσίες Ανάπτυξης",
							"price": "90.00",
							"ext": {
								"gr-mydata-income-cat": "category1_1",
								"gr-mydata-income-type": "E3_561_001",
							}
						}
					}
				]
				~~~
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Technical description of REST API interfaces for sending & receiving data for ERP users (Version 2.0.0)",
					i18n.EL: "Τεχνική περιγραφή διεπαφών REST API για διαβίβαση & λήψη δεδομένων για χρήστες ERP (Έκδοση 2.0.0)",
				},
				URL: "https://www.aade.gr/sites/default/files/2025-12/myDATA%20API%20Documentation%20v2.0.0_official_erp.pdf",
			},
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
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				See the Income Classification Category for more information.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Technical description of REST API interfaces for sending & receiving data for ERP users (Version 2.0.0)",
					i18n.EL: "Τεχνική περιγραφή διεπαφών REST API για διαβίβαση & λήψη δεδομένων για χρήστες ERP (Έκδοση 2.0.0)",
				},
				URL: "https://www.aade.gr/sites/default/files/2025-12/myDATA%20API%20Documentation%20v2.0.0_official_erp.pdf",
			},
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
	{
		Key: ExtKeyTaxType,
		Name: i18n.String{
			i18n.EN: "Tax type",
			i18n.EL: "Είδος Φόρου",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Greek invoices may contain charges or allowances for taxes and fees different from VAT. The
				~gr-mydata-tax-type~ extension at charge level must be used to specify the tax type.

				GOBL will attempt to automatically set the ~gr-mydata-tax-type~ extension. The following
				table shows the mappings used: first, GOBL will check the charge's ~key~ field; if that
				does not match, it will then check for the presence of the related extension:

				| Code | Name       | GOBL Charge Key | Extension              |
				| ---- | ---------- | --------------- | ---------------------- |
				| ~2~  | Fee        |                 | ~gr-mydata-fee~        |
				| ~3~  | Other Tax  | ~tax~           | ~gr-mydata-other-tax~  |
				| ~4~  | Stamp Duty | ~stamp-duty~    | ~gr-mydata-stamp-duty~ |

				For example (manual):

				~~~json
				"charges": [
					{
						"amount": "10.00",
						"reason": "Stamp duty fee",
						"ext": {
							"gr-mydata-tax-type": "2",
							"gr-mydata-fee": "13",
						}
					}
				]
				~~~
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Technical description of REST API interfaces for sending & receiving data for ERP users (Version 2.0.0)",
					i18n.EL: "Τεχνική περιγραφή διεπαφών REST API για διαβίβαση & λήψη δεδομένων για χρήστες ERP (Έκδοση 2.0.0)",
				},
				URL: "https://www.aade.gr/sites/default/files/2025-12/myDATA%20API%20Documentation%20v2.0.0_official_erp.pdf",
			},
		},
		Values: []*cbc.Definition{
			{
				Code: TaxTypeWithholding,
				Name: i18n.String{
					i18n.EN: "Withholding Tax",
					i18n.EL: "Παρακρατούμενος Φόρος",
				},
			},
			{
				Code: TaxTypeFee,
				Name: i18n.String{
					i18n.EN: "Fee",
					i18n.EL: "Τέλη",
				},
			},
			{
				Code: TaxTypeOtherTax,
				Name: i18n.String{
					i18n.EN: "Other Tax",
					i18n.EL: "Λοιποί Φόροι",
				},
			},
			{
				Code: TaxTypeStampDuty,
				Name: i18n.String{
					i18n.EN: "Stamp Duty",
					i18n.EL: "Ψηφιακού Τέλος συναλλαγής",
				},
			},
			{
				Code: TaxTypeDeduction,
				Name: i18n.String{
					i18n.EN: "Deduction",
					i18n.EL: "Κρατήσεις",
				},
			},
		},
	},
	{
		Key: ExtKeyFee,
		Name: i18n.String{
			i18n.EN: "Fee category",
			i18n.EL: "Κατηγορία Τελών",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Greek invoices may contain charges for different types of fees. In GOBL, you must use
				the ~gr-mydata-fee~ extension at charge level to specify the specific fee category
				for the charge.

				For example:

				~~~json
				"charges": [
					{
						"amount": "10.00",
						"reason": "Subscription fee",
						"ext": {
							"gr-mydata-tax-type": "2",
							"gr-mydata-fee": "13",
						}
					}
				]
				~~~
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Technical description of REST API interfaces for sending & receiving data for ERP users (Version 2.0.0)",
					i18n.EL: "Τεχνική περιγραφή διεπαφών REST API για διαβίβαση & λήψη δεδομένων για χρήστες ERP (Έκδοση 2.0.0)",
				},
				URL: "https://www.aade.gr/sites/default/files/2025-12/myDATA%20API%20Documentation%20v2.0.0_official_erp.pdf",
			},
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Monthly account up to 50 euros 12%",
					i18n.EL: "Για μηνιαίο λογαριασμό μέχρι και 50 ευρώ 12%",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Monthly account from 50.01 to 100 euros 15%",
					i18n.EL: "Για μηνιαίο λογαριασμό από 50,01 μέχρι και 100 ευρώ 15%",
				},
			},
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "Monthly account from 100.01 to 150 euros 18%",
					i18n.EL: "Για μηνιαίο λογαριασμό από 100,01 μέχρι και 150 ευρώ 18%",
				},
			},
			{
				Code: "4",
				Name: i18n.String{
					i18n.EN: "Monthly account from 150.01 euros and above 20%",
					i18n.EL: "Για μηνιαίο λογαριασμό από 150,01 ευρώ και άνω 20%",
				},
			},
			{
				Code: "5",
				Name: i18n.String{
					i18n.EN: "Mobile telephony fee based on call time value (12%)",
					i18n.EL: "Τέλος καρτοκινητής επί της αξίας του χρόνου ομιλίας (12%)",
				},
			},
			{
				Code: "6",
				Name: i18n.String{
					i18n.EN: "Subscription television fee 10%",
					i18n.EL: "Τέλος στη συνδρομητική τηλεόραση 10%",
				},
			},
			{
				Code: "7",
				Name: i18n.String{
					i18n.EN: "Fixed-line subscription fee 5%",
					i18n.EL: "Τέλος συνδρομητών σταθερής τηλεφωνίας 5%",
				},
			},
			{
				Code: "8",
				Name: i18n.String{
					i18n.EN: "Environmental fee & plastic bag v. 2339/2001 art. 6a 0.07 euros per item",
					i18n.EL: "Περιβαλλοντικό Τέλος & πλαστικής σακούλας ν. 2339/2001 αρ. 6α 0,07 ευρώ ανά τεμάχιο",
				},
			},
			{
				Code: "9",
				Name: i18n.String{
					i18n.EN: "Postal service contribution 2%",
					i18n.EL: "Εισφορά δακοκτονίας 2%",
				},
			},
			{
				Code: "10",
				Name: i18n.String{
					i18n.EN: "Other fees",
					i18n.EL: "Λοιπά τέλη",
				},
			},
			{
				Code: "11",
				Name: i18n.String{
					i18n.EN: "Petroleum fees",
					i18n.EL: "Τέλη Λοιπών Φόρων",
				},
			},
			{
				Code: "12",
				Name: i18n.String{
					i18n.EN: "Postal service contribution",
					i18n.EL: "Εισφορά δακοκτονίας",
				},
			},
			{
				Code: "13",
				Name: i18n.String{
					i18n.EN: "Monthly account per connection (10%)",
					i18n.EL: "Για μηνιαίο λογαριασμό κάθε σύνδεσης (10%)",
				},
			},
			{
				Code: "14",
				Name: i18n.String{
					i18n.EN: "Mobile and mobile telephony fee based on call time value (10%)",
					i18n.EL: "Τέλος καρτοκινητής επί της αξίας του χρόνου ομιλίας (10%)",
				},
			},
			{
				Code: "15",
				Name: i18n.String{
					i18n.EN: "Mobile and mobile telephony fee for natural persons aged 15 to 29 years (0%)",
					i18n.EL: "Τέλος κινητής και καρτοκινητής για φυσικά πρόσωπα ηλικίας 15 έως και 29 ετών (0%)",
				},
			},
			{
				Code: "16",
				Name: i18n.String{
					i18n.EN: "Environmental protection contribution for plastic products 0.04 cents per item [article 4 v. 4736/2020]",
					i18n.EL: "Εισφορά προστασίας περιβάλλοντος πλαστικών προϊόντων 0,04 λεπτά ανά τεμάχιο [άρθρο 4 ν. 4736/2020]",
				},
			},
			{
				Code: "17",
				Name: i18n.String{
					i18n.EN: "Recycling fee 0.08 cents per item [article 80 v. 4819/2021]",
					i18n.EL: "Τέλος ανακύκλωσης 0,08 λεπτά ανά τεμάχιο [άρθρο 80 ν. 4819/2021]",
				},
			},
			{
				Code: "18",
				Name: i18n.String{
					i18n.EN: "Accommodation fee",
					i18n.EL: "Τέλος διαμονής παρεπιδημούντων",
				},
			},
			{
				Code: "19",
				Name: i18n.String{
					i18n.EN: "Fee on gross revenues of restaurants and related establishments",
					i18n.EL: "Τέλος επί των ακαθάριστων εσόδων των εστιατορίων και συναφών καταστημάτων",
				},
			},
			{
				Code: "20",
				Name: i18n.String{
					i18n.EN: "Fee on gross revenues of entertainment centers",
					i18n.EL: "Τέλος επί των ακαθάριστων εσόδων των κέντρων διασκέδασης",
				},
			},
			{
				Code: "21",
				Name: i18n.String{
					i18n.EN: "Fee on gross revenues of casinos",
					i18n.EL: "Τέλος επί των ακαθάριστων εσόδων των καζίνο",
				},
			},
			{
				Code: "22",
				Name: i18n.String{
					i18n.EN: "Other fees on gross revenues",
					i18n.EL: "Λοιπά τέλη επί των ακαθάριστων εσόδων",
				},
			},
		},
	},
	{
		Key: ExtKeyOtherTax,
		Name: i18n.String{
			i18n.EN: "Other taxes category",
			i18n.EL: "Κατηγορία Λοιπών Φόρων",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Greek invoices may contain charges for different taxes other than VAT. The
				~gr-mydata-other-tax~ extension at charge level must be used to specify the specific
				other tax category for the charge.

				For example:

				~~~json
				"charges": [
					{
						"amount": "3.00",
						"reason": "Climate Crisis Fee",
						"ext": {
							"gr-mydata-tax-type": "3",
							"gr-mydata-other-tax": "8",
						}
					}
				]
				~~~
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Technical description of REST API interfaces for sending & receiving data for ERP users (Version 2.0.0)",
					i18n.EL: "Τεχνική περιγραφή διεπαφών REST API για διαβίβαση & λήψη δεδομένων για χρήστες ERP (Έκδοση 2.0.0)",
				},
				URL: "https://www.aade.gr/sites/default/files/2025-12/myDATA%20API%20Documentation%20v2.0.0_official_erp.pdf",
			},
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "a1) fire insurance premiums 20%",
					i18n.EL: "α1) ασφάλιστρα κλάδου πυρός 20%",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "a2) fire insurance premiums 20%",
					i18n.EL: "α2) ασφάλιστρα κλάδου πυρός 20%",
				},
			},
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "b) life insurance premiums 4%",
					i18n.EL: "β) ασφάλιστρα κλάδου ζωής 4%",
				},
			},
			{
				Code: "4",
				Name: i18n.String{
					i18n.EN: "c) insurance premiums for other branches 15%",
					i18n.EL: "γ) ασφάλιστρα λοιπών κλάδων 15%",
				},
			},
			{
				Code: "5",
				Name: i18n.String{
					i18n.EN: "d) tax-exempt insurance premiums 0%",
					i18n.EL: "δ) απαλλασσόμενα φόρου ασφάλιστρα 0%",
				},
			},
			{
				Code: "6",
				Name: i18n.String{
					i18n.EN: "Hotels 1-2 stars 0.50 €",
					i18n.EL: "Ξενοδοχεία 1-2 αστέρων 0,50 €",
				},
			},
			{
				Code: "7",
				Name: i18n.String{
					i18n.EN: "Hotels 3 stars 1.50 €",
					i18n.EL: "Ξενοδοχεία 3 αστέρων 1,50 €",
				},
			},
			{
				Code: "8",
				Name: i18n.String{
					i18n.EN: "Hotels 4 stars 3.00 €",
					i18n.EL: "Ξενοδοχεία 4 αστέρων 3,00 €",
				},
			},
			{
				Code: "9",
				Name: i18n.String{
					i18n.EN: "Hotels 4 stars 4.00 €",
					i18n.EL: "Ξενοδοχεία 4 αστέρων 4,00 €",
				},
			},
			{
				Code: "10",
				Name: i18n.String{
					i18n.EN: "Rented - furnished rooms - apartments 0.50 €",
					i18n.EL: "Ενοικιαζόμενα - επιπλωμένα δωμάτια - διαμερίσματα 0,50 €",
				},
			},
			{
				Code: "11",
				Name: i18n.String{
					i18n.EN: "Special tax on advertisements broadcast on television (EFTD) 5%",
					i18n.EL: "Ειδικός Φόρος στις διαφημίσεις που προβάλλονται από την τηλεόραση (ΕΦΔΤ) 5%",
				},
			},
			{
				Code: "12",
				Name: i18n.String{
					i18n.EN: "3.1 Luxury tax 10% on the taxable value for intra-Community acquisitions and imported from third countries 10%",
					i18n.EL: "3.1 Φόρος πολυτελείας 10% επί της φορολογητέας αξίας για τα ενδοκοινοτικώς αποκτούμενα και εισαγόμενα από τρίτες χώρες 10%",
				},
			},
			{
				Code: "13",
				Name: i18n.String{
					i18n.EN: "3.2 Luxury tax 10% on the selling price before VAT for domestically produced goods 10%",
					i18n.EL: "3.2 Φόρος πολυτελείας 10% επί της τιμής πώλησης προ Φ.Π.Α. για τα εγχωρίως παραγόμενα είδη 10%",
				},
			},
			{
				Code: "14",
				Name: i18n.String{
					i18n.EN: "State rights on casino tickets (80% on the ticket)",
					i18n.EL: "Δικαίωμα του Δημοσίου στα εισιτήρια των καζίνο (80% επί του εισιτηρίου)",
				},
			},
			{
				Code: "15",
				Name: i18n.String{
					i18n.EN: "Fire insurance premiums 20%",
					i18n.EL: "ασφάλιστρα κλάδου πυρός 20%",
				},
			},
			{
				Code: "16",
				Name: i18n.String{
					i18n.EN: "Other Customs Duties-Taxes",
					i18n.EL: "Λοιποί Τελωνειακοί Δασμοί-Φόροι",
				},
			},
			{
				Code: "17",
				Name: i18n.String{
					i18n.EN: "Other Taxes",
					i18n.EL: "Λοιποί Φόροι",
				},
			},
			{
				Code: "18",
				Name: i18n.String{
					i18n.EN: "Charges of Other Taxes",
					i18n.EL: "Επιβαρύνσεις Λοιπών Φόρων",
				},
			},
			{
				Code: "19",
				Name: i18n.String{
					i18n.EN: "EFK",
					i18n.EL: "ΕΦΚ",
				},
			},
			{
				Code: "20",
				Name: i18n.String{
					i18n.EN: "Hotels 1-2 stars 1.50€ (per Room/Night)",
					i18n.EL: "Ξενοδοχεία 1-2 αστέρων 1,50€ (ανά Δωμ./Διαμ.)",
				},
			},
			{
				Code: "21",
				Name: i18n.String{
					i18n.EN: "Hotels 3 stars 3.00€ (per Room/Night)",
					i18n.EL: "Ξενοδοχεία 3 αστέρων 3,00€ (ανά Δωμ./Διαμ.)",
				},
			},
			{
				Code: "22",
				Name: i18n.String{
					i18n.EN: "Hotels 4 stars 7.00€ (per Room/Night)",
					i18n.EL: "Ξενοδοχεία 4 αστέρων 7,00€ (ανά Δωμ./Διαμ.)",
				},
			},
			{
				Code: "23",
				Name: i18n.String{
					i18n.EN: "Hotels 5 stars 10.00€ (per Room/Night)",
					i18n.EL: "Ξενοδοχεία 5 αστέρων 10,00€ (ανά Δωμ./Διαμ.)",
				},
			},
			{
				Code: "24",
				Name: i18n.String{
					i18n.EN: "Rented furnished rooms – apartments 1.50€ (per Room/Night)",
					i18n.EL: "Ενοικιαζόμενα επιπλωμένα δωμάτια – διαμερίσματα 1,50€ (ανά Δωμ./Διαμ.)",
				},
			},
			{
				Code: "25",
				Name: i18n.String{
					i18n.EN: "Short-term rental properties 1.50€",
					i18n.EL: "Ακίνητα βραχυχρόνιας μίσθωσης 1,50€",
				},
			},
			{
				Code: "26",
				Name: i18n.String{
					i18n.EN: "Short-term rental properties detached houses over 80 sq.m. 10.00€",
					i18n.EL: "Ακίνητα βραχυχρόνιας μίσθωσης μονοκατοικίες άνω 80 τ.μ. 10,00€",
				},
			},
			{
				Code: "27",
				Name: i18n.String{
					i18n.EN: "Self-catering accommodation – tourist furnished villas (villas) 10.00€",
					i18n.EL: "Αυτοεξυπηρετούμενα καταλύματα – τουριστικές επιπλωμένες επαύλεις (βίλες) 10,00€",
				},
			},
			{
				Code: "28",
				Name: i18n.String{
					i18n.EN: "Short-term rental properties 0.50€",
					i18n.EL: "Ακίνητα βραχυχρόνιας μίσθωσης 0,50€",
				},
			},
			{
				Code: "29",
				Name: i18n.String{
					i18n.EN: "Short-term rental properties detached houses over 80 sq.m. 4.00€",
					i18n.EL: "Ακίνητα βραχυχρόνιας μίσθωσης μονοκατοικίες άνω 80 τ.μ. 4,00€",
				},
			},
			{
				Code: "30",
				Name: i18n.String{
					i18n.EN: "Self-catering accommodation – tourist furnished villas (villas) 4.00€",
					i18n.EL: "Αυτοεξυπηρετούμενα καταλύματα – τουριστικές επιπλωμένες επαύλεις (βίλες) 4,00€",
				},
			},
		},
	},
	{
		Key: ExtKeyStampDuty,
		Name: i18n.String{
			i18n.EN: "Stamp duty coefficient category",
			i18n.EL: "Κατηγορία Συντελεστή Ψηφιακού Τέλους συναλλαγής",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Greek invoices may contain charges for stamp duty. The ~gr-mydata-stamp-duty~ extension
				at charge level must be used to specify the specific stamp duty coefficient category
				for the charge.

				For example:

				~~~json
				"charges": [
					{
						"amount": "12.00",
						"reason": "Stamp duty",
						"ext": {
							"gr-mydata-tax-type": "4",
							"gr-mydata-stamp-duty": "1"
						}
					}
				]
				~~~
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Technical description of REST API interfaces for sending & receiving data for ERP users (Version 2.0.0)",
					i18n.EL: "Τεχνική περιγραφή διεπαφών REST API για διαβίβαση & λήψη δεδομένων για χρήστες ERP (Έκδοση 2.0.0)",
				},
				URL: "https://www.aade.gr/sites/default/files/2025-12/myDATA%20API%20Documentation%20v2.0.0_official_erp.pdf",
			},
		},
		Values: []*cbc.Definition{
			{
				Code: "1",
				Name: i18n.String{
					i18n.EN: "Coefficient 1.2%",
					i18n.EL: "Συντελεστής 1,2 %",
				},
			},
			{
				Code: "2",
				Name: i18n.String{
					i18n.EN: "Coefficient 2.4%",
					i18n.EL: "Συντελεστής 2,4 %",
				},
			},
			{
				Code: "3",
				Name: i18n.String{
					i18n.EN: "Coefficient 3.6%",
					i18n.EL: "Συντελεστής 3,6 %",
				},
			},
			{
				Code: "4",
				Name: i18n.String{
					i18n.EN: "Other cases",
					i18n.EL: "Λοιπές περιπτώσεις",
				},
			},
		},
	},
}

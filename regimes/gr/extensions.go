package gr

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// Regime extension codes.
const (
	ExtKeyIAPRVATCat    = "gr-iapr-vat-cat"
	ExtKeyIAPRExemption = "gr-iapr-exemption"
)

var extensionKeys = []*cbc.KeyDefinition{
	{
		Key: ExtKeyIAPRVATCat,
		Name: i18n.String{
			i18n.EN: "VAT category",
			i18n.EL: "Κατηγορία ΦΠΑ",
		},
		Codes: []*cbc.CodeDefinition{
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
		Key: ExtKeyIAPRExemption,
		Name: i18n.String{
			i18n.EN: "VAT exemption cause",
			i18n.EL: "Κατηγορία Αιτίας Εξαίρεσης ΦΠΑ",
		},
		Codes: []*cbc.CodeDefinition{
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
}

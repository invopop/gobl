package es

import (
	"github.com/invopop/gobl/org"
)

// InvoiceLegalNoteExamples defines a list of notes which may be required by Spanish law.
// These are expected to be used in user interfaces as examples that can be modified
// according to the details of the invoice. Most of this data has now been moved to
// scheme definitions, but some examples require a bit more effort from the user side.
var InvoiceLegalNoteExamples = map[string]*org.Note{
	// Exempt transaction pursuant to Article [fill in with the relevant Article number] of the Law No 37/1992 of 28 December 1992 on value added tax
	"exempt": {
		Key:  org.NoteKeyLegal,
		Text: "Operación exenta por aplicación del artículo [indicar el articulo] de la Ley 37/1992, del 28 de diciembre, del Impuesto sobre el Valor Añadido.",
	},
	// Means of transport [fill in with a brief description, for example, Jaguar S-Type automobile] the date of first start-up [fill in with the date] and the distance covered or hours of navigation [fill in with the distance covered or hours of navigation, for example, 5,900 km or 48 hours]
	"transport": {
		Key:  org.NoteKeyLegal,
		Text: "Medio de transporte [describir el medio, por ejemplo automóvil turismo Seat Ibiza TDI 2.0] fecha 1ª puesta en servicio [indicar la fecha] distancias/horas recorridas [indicar la distancia o las horas, por ejemplo, 5.900 km o 48 horas].",
	},
}

package es

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// Item Keys which may be used by TicketBAI in the Basque Country.
const (
	ItemResale   cbc.Key = "resale"
	ItemServices cbc.Key = "services"
	ItemGoods    cbc.Key = "goods"
)

// Special item keys required by the TicketBAI system in the basque country
// on a per-line basis.
var itemKeyDefinitions = []*tax.KeyDefinition{
	{
		Key: ItemResale,
		Name: i18n.String{
			i18n.ES: "Reventa de bienes sin modificación por vendedor en regimen simplificado",
			i18n.EN: "Resale of goods without modification by vendor in the simplified regime",
		},
	},
	{
		Key: ItemServices,
		Name: i18n.String{
			i18n.ES: "Prestacion de servicios",
			i18n.EN: "Provision of services",
		},
	},
	{
		Key: ItemGoods,
		Name: i18n.String{
			i18n.ES: "Entrega de bienes",
			i18n.EN: "Delivery of goods",
		},
	},
}

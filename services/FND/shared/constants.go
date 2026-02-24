package shared

// CenterCurrency is the system's central/base currency (TWD for Taiwan).
const CenterCurrency = "TWD"

// SettlementCurrencies are the non-miscellaneous foreign currencies.
// A currency NOT in this map (and not CenterCurrency) is considered "miscellaneous" (雜幣).
var SettlementCurrencies = map[string]bool{
	"TWD": true, "USD": true, "EUR": true, "JPY": true,
	"GBP": true, "AUD": true, "CNY": true, "HKD": true,
	"CAD": true, "CHF": true, "SGD": true, "NZD": true,
}

// IsMiscCurrency returns true if the currency is a "miscellaneous" currency (雜幣).
// A currency is miscellaneous if it is NOT in SettlementCurrencies AND NOT CenterCurrency.
func IsMiscCurrency(cryID string) bool {
	if cryID == CenterCurrency {
		return false
	}
	return !SettlementCurrencies[cryID]
}

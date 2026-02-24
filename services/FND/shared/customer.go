package shared

// GetCusName - placeholder for customer name lookup.
func GetCusName(cusID string) string {
	switch cusID {
	case "C001":
		return "VIP Customer A"
	case "C002":
		return "Corporate Customer B"
	default:
		if cusID != "" {
			return cusID + " Name"
		}
		return ""
	}
}

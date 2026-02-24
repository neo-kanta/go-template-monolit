// Package shared provides cross-module lookup utilities for the FND service.
// TODO: Replace stubs with actual SysUtility / DB lookups when shared services are ready.
package shared

// ─── Currency Name Lookup (SysUtility.GetCry.CryIDNM) ───

var currencyNames = map[string]string{
	"TWD": "新臺幣", "USD": "美元", "EUR": "歐元", "JPY": "日圓",
	"GBP": "英鎊", "AUD": "澳幣", "CNY": "人民幣", "HKD": "港幣",
	"ZAR": "南非幣", "BRL": "巴西幣", "NZD": "紐幣", "CAD": "加幣",
	"CHF": "瑞士法郎", "SGD": "新加坡幣", "SEK": "瑞典克朗",
}

var currencyDecLen = map[string]int32{
	"TWD": 0, "USD": 2, "EUR": 2, "JPY": 0, "GBP": 2,
	"AUD": 2, "CNY": 2, "HKD": 2, "ZAR": 2, "BRL": 2,
	"NZD": 2, "CAD": 2, "CHF": 2, "SGD": 2, "SEK": 2,
}

// GetCryName returns the display name for a currency code.
func GetCryName(cryID string) string {
	if nm, ok := currencyNames[cryID]; ok {
		return nm
	}
	return cryID
}

// GetCryDecLen returns the decimal length for a currency.
func GetCryDecLen(cryID string) int32 {
	if dl, ok := currencyDecLen[cryID]; ok {
		return dl
	}
	return 2
}

// ─── Bank Name Lookup (SysUtility.GetBankV / GetBankVNM) ───

// BankInfo holds stub bank data.
type BankInfo struct {
	Name      string
	SWIFTCode string
	SName     string
	Addr      string
}

var bankData = map[string]BankInfo{
	"004": {Name: "臺灣銀行", SWIFTCode: "BKTWTWTP", SName: "Bank of Taiwan", Addr: "臺北市中正區重慶南路一段120號"},
	"005": {Name: "土地銀行", SWIFTCode: "LBOTTWTP", SName: "Land Bank of Taiwan", Addr: "臺北市館前路46號"},
	"006": {Name: "合作金庫", SWIFTCode: "TACBTWTP", SName: "TCB Bank", Addr: "臺北市館前路77號"},
	"007": {Name: "第一銀行", SWIFTCode: "FCBKTWTP", SName: "First Bank", Addr: "臺北市重慶南路一段30號"},
	"008": {Name: "華南銀行", SWIFTCode: "HNBKTWTP", SName: "Hua Nan Bank", Addr: "臺北市信義區松仁路123號"},
	"009": {Name: "彰化銀行", SWIFTCode: "CCBCTWTP", SName: "Chang Hwa Bank", Addr: "臺中市自由路二段38號"},
	"012": {Name: "台北富邦", SWIFTCode: "TPBKTWTP", SName: "Taipei Fubon Bank", Addr: "臺北市建國南路一段236號"},
	"013": {Name: "國泰世華", SWIFTCode: "UWCBTWTP", SName: "Cathay United Bank", Addr: "臺北市信義區松仁路7號"},
	"017": {Name: "兆豐銀行", SWIFTCode: "ICBCTWTP", SName: "Mega ICBC", Addr: "臺北市忠孝東路二段85號"},
	"021": {Name: "花旗銀行", SWIFTCode: "CABORC2X", SName: "Citibank Taiwan", Addr: "臺北市松智路1號"},
}

// GetBankInfo returns stub bank info for a bank ID.
func GetBankInfo(bankID string) BankInfo {
	if b, ok := bankData[bankID]; ok {
		return b
	}
	return BankInfo{Name: bankID}
}

// ─── Code Name Lookup (SysUtility.GetCode.CodeNM) ───

var tcCodeNames = map[string]map[string]string{
	"001": {"01": "股票型", "02": "債券型", "03": "平衡型", "04": "貨幣市場型", "05": "組合型", "06": "保本型", "07": "指數型", "08": "指數股票型(ETF)", "09": "不動產證券化型", "10": "其他"},
	"021": {"A": "A群組", "B": "B群組", "C": "C群組", "D": "D群組", "E": "E群組"},
	"022": {"01": "群組一", "02": "群組二", "03": "群組三"},
	"030": {"A": "A群組", "B": "B群組", "C": "C群組"},
	"031": {"A": "A級別", "B": "B級別", "C": "C級別", "I": "I級別(法人)", "N": "N級別", "S": "S級別", "P": "P級別"},
	"033": {"01": "群組一", "02": "群組二"},
}

// GetCodeName returns the display name for a TC code value.
func GetCodeName(tcCode, code string) string {
	if codes, ok := tcCodeNames[tcCode]; ok {
		if nm, ok := codes[code]; ok {
			return nm
		}
	}
	return code
}

// ─── Company Name Lookup (SysUtility.GetSysCoNM) ───

var companyNames = map[string]string{
	"SWSTD":  "範例投信股份有限公司",
	"FUND02": "台灣投信股份有限公司",
	"SITCA":  "中華民國證券投資信託暨顧問商業同業公會",
}

// GetSysCoName returns the company display name.
func GetSysCoName(sysCoID string) string {
	if nm, ok := companyNames[sysCoID]; ok {
		return nm
	}
	return sysCoID
}

// ─── Code Type Name Lookup ───

var codeTypeNames = map[string]string{
	"022": "通路服務費",
	"033": "基金停利",
}

// GetCodeTypeName returns the code type display name.
func GetCodeTypeName(codeType string) string {
	if nm, ok := codeTypeNames[codeType]; ok {
		return nm
	}
	return codeType
}

// ─── Fund Name Lookup ───

// GetDFundName returns a fund display name (stub).
func GetDFundName(fundCode string) string {
	return fundCode + "(Fund Name)"
}

// ─── Option Dictionary (SysUtility.GetOption.ItemName) ───

var optionNames = map[string]map[string]string{
	// 000003: 基金風險屬性 (Fund Risk Level — SITCA RR rating)
	"000003": {
		"1": "RR1", "2": "RR2", "3": "RR3", "4": "RR4", "5": "RR5",
	},
	// 000007: 基金狀態 (Fund Status)
	"000007": {
		"N": "正常", "S": "暫停申購", "R": "暫停買回",
		"X": "暫停交易", "C": "清算", "M": "合併",
	},
}

// GetOptName returns the display name for an option (OptID, ItemID) pair.
// Maps to C# TAFNDFundInfoMethods.GetOpt("000003,000007").
func GetOptName(optID, itemID string) string {
	if items, ok := optionNames[optID]; ok {
		if nm, ok := items[itemID]; ok {
			return nm
		}
	}
	return itemID
}

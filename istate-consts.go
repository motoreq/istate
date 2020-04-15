//

package istate

const (
	iStateTag            = "istate"
	iStatePrimaryTag     = "primary"
	iStatePrimaryTrueVal = "true"

	docTypeField   = ".docType"
	keyRefField    = ".keyref"
	valueField     = ".value"
	fieldNameField = ".fieldName"

	// Encode
	boolTrue  = "t"
	boolFalse = "f"

	// Inverting below causes malfunction
	positiveNum = "1"
	negativeNum = "0"

	star = "*"
	dot  = "."

	null      = "\000"
	seperator = null
	asciiLast = "~"
	incChar   = "\003" // This must not be less than other predefined symbols

	// Warning: These must be changed when encoding method changes
	numSym       = "\001"
	numSeparator = "\002"
	pNumPrefix   = numSym + positiveNum
	nNumPrefix   = numSym + negativeNum
	zero         = numSym + positiveNum + "01_0"
	// biggestPNum  = numSym + "099_" //+ asciiLast
	// biggnestNNum = numSym + "199_" //+ asciiLast
	biggestPNum = numSym + positiveNum + "20_18446744073709552000" // largest(uint64) + 1
	biggestNNum = numSym + negativeNum + "20_18446744073709552000" // largest(uint64) + 1

	int32Biggest = 2147483647
)

var nullByte = []byte{0x01}

var numDigits = map[int]string{
	1:  "01",
	2:  "02",
	3:  "03",
	4:  "04",
	5:  "05",
	6:  "06",
	7:  "07",
	8:  "08",
	9:  "09",
	10: "10",
	11: "11",
	12: "12",
	13: "13",
	14: "14",
	15: "15",
	16: "16",
	17: "17",
	18: "18",
	19: "19",
	20: "20",
	21: "21",
	22: "22",
	23: "23",
	24: "24",
	25: "25",
	26: "26",
	27: "27",
	28: "28",
	29: "29",
	30: "30",
	31: "31",
	32: "32",
	33: "33",
	34: "34",
	35: "35",
	36: "36",
	37: "37",
	38: "38",
	39: "39",
	40: "40",
	41: "41",
	42: "42",
	43: "43",
	44: "44",
	45: "45",
	46: "46",
	47: "47",
	48: "48",
	49: "49",
	50: "50",
	51: "51",
	52: "52",
	53: "53",
	54: "54",
	55: "55",
	56: "56",
	57: "57",
	58: "58",
	59: "59",
	60: "60",
	61: "61",
	62: "62",
	63: "63",
	64: "64",
	65: "65",
	66: "66",
	67: "67",
	68: "68",
	69: "69",
	70: "70",
	71: "71",
	72: "72",
	73: "73",
	74: "74",
	75: "75",
	76: "76",
	77: "77",
	78: "78",
	79: "79",
	80: "80",
	81: "81",
	82: "82",
	83: "83",
	84: "84",
	85: "85",
	86: "86",
	87: "87",
	88: "88",
	89: "89",
	90: "90",
	91: "91",
	92: "92",
	93: "93",
	94: "94",
	95: "95",
	96: "96",
	97: "97",
	98: "98",
	99: "99",
}

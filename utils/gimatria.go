package utils

import (
	"strconv"
	"strings"
)

func NumToHebrew(num uint16) string {
	if num == 15 {
		return "טו"
	}
	if num == 16 {
		return "טז"
	}

	digits := strconv.Itoa(int(num))
	if num > 9999 {
		return digits
	}

	res := ""
	for i := 0; i < len(digits); i++ {
		d, _ := strconv.Atoi(digits[i:i+1])
		m := len(digits) - (i + 1)
		if m == 3 {
			res = res + numToHeb(d, m) + " "
		} else if m == 1 && num > 100 && digits[i-1:i] == "0" {
			if digits[i:i+2] == "15" {
				return res + "טו"
			} else if digits[i:i+2] == "16" {
				return res + "טז"
			} else {
				res = res + numToHeb(d, m)
			}
		} else {
			res = res + numToHeb(d, m)
		}
	}

	return strings.TrimSpace(res)
}

var MAGNITUDE_MAPS = map[int]map[int]string{
	0: {
		1: "א",
		2: "ב",
		3: "ג",
		4: "ד",
		5: "ה",
		6: "ו",
		7: "ז",
		8: "ח",
		9: "ט",
	},
	1: {
		1: "י",
		2: "כ",
		3: "ל",
		4: "מ",
		5: "נ",
		6: "ס",
		7: "ע",
		8: "פ",
		9: "צ",
	},
	2: {
		1: "ק",
		2: "ר",
		3: "ש",
		4: "ת",
		5: "תק",
		6: "תר",
		7: "תש",
		8: "תת",
		9: "תתק",
	},
	3: {
		1: "א'",
		2: "ב'",
		3: "ג'",
		4: "ד'",
		5: "ה'",
		6: "ו'",
		7: "ז'",
		8: "ח'",
		9: "ט'",
	},
}

func numToHeb(num int, magnitude int) string {
	return MAGNITUDE_MAPS[magnitude][num]
}

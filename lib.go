package country_info

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

/// just dont find any native enums
const (
	ErrOpenRefFile = 0
	ErrOpenCCFile  = 1
)

func checkError(e error, own_e int) {
	if e != nil {
		switch own_e {

		case ErrOpenRefFile:
			fmt.Println("Cant open reference file. Aborted")

		case ErrOpenCCFile:
			fmt.Println("Cant open file with country code (CC). Aborted")

		default:
			fmt.Println("There are another arbitrary error was appeared. Aborted")
		}

		panic(e)
	}
}

/// get information from both sources and merge it into a new CountryCell type that represent full information
func mixCountryInfo(country_map *map[string]CountryCell, ref *[]ReferenceCountryCell, cc *[]CCCountryCell) {

	/// find entry where CodeName23 contain CСodeName3 and merge this info into the CountryCell
	/// different source have differenct numbers of entryes so I decide make something like ReferenceCountryCell LEFT JOIN CCCountryCell
	for _, v_ref := range *ref {

		is_merged := false

		for _, v_cc := range *cc {

			//fmt.Println(i_ref, v_ref, i_cc, v_cc)
			if strings.Contains(v_cc.CodeName23, v_ref.CСodeName3) {

				(*country_map)[v_ref.CСodeName3] = CountryCell{
					Name:        v_ref.NName,
					CodeName2:   v_ref.CCodeName2,
					CodeName3:   v_ref.CСodeName3,
					Numeric:     v_ref.NNumeric,
					PhonePrefix: v_cc.PhonePrefix,
				}

				is_merged = true
				break
			}
		}

		if is_merged == false {
			(*country_map)[v_ref.CСodeName3] = CountryCell{
				Name:        v_ref.NName,
				CodeName2:   v_ref.CCodeName2,
				CodeName3:   v_ref.CСodeName3,
				Numeric:     v_ref.NNumeric,
				PhonePrefix: "",
			}
		}

	}
}

type CountryCell struct {
	Name        string
	CodeName2   string
	CodeName3   string
	Numeric     int
	PhonePrefix string
}

/// still work even fields names dont match with a json enthryes till fields names start with lowercase letter
type ReferenceCountryCell struct {
	NName      string `json:"Name"`
	CCodeName2 string `json:"CodeName2"`
	CСodeName3 string `json:"CodeName3"`
	NNumeric   int    `json:"Numeric"`
}

/// still work perfect without explicit tag declaration during fields names matching with the json entries
type CCCountryCell struct {
	PhonePrefix string
	CodeName23  string
}

/// work also fine, cas (how i understand) map alloc memory on heap and wont delete it during program has at least one pointer
func Init() *map[string]CountryCell {

	reference_json, err := os.ReadFile("./table_reference.json")
	checkError(err, ErrOpenRefFile)

	cc_json, err := os.ReadFile("./table_cc.json")
	checkError(err, ErrOpenCCFile)

	var ref_country_arr []ReferenceCountryCell
	json.Unmarshal(reference_json, &ref_country_arr)

	var cc_country_arr []CCCountryCell
	json.Unmarshal(cc_json, &cc_country_arr)

	m := make(map[string]CountryCell)
	mixCountryInfo(&m, &ref_country_arr, &cc_country_arr)

	return &m
}

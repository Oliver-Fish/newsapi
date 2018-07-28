package newsapi

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type allowedParameters map[string]string
type parameters map[string]interface{}

func (ap allowedParameters) verify(p parameters) error {
	for k, v := range p {
		if _, ok := ap[k]; !ok {
			return fmt.Errorf("Invalid parameter %v", k)
		}
		if reflect.TypeOf(v).String() != ap[k] {
			return fmt.Errorf("Invalid type for parameter %v expected type %v got type %v", k, ap[k], reflect.TypeOf(v).String())
		}
	}
	return nil
}

var allowedCountryCodes = map[string]struct{}{"ae": struct{}{}, "ar": struct{}{}, "at": struct{}{}, "au": struct{}{}, "be": struct{}{}, "bg": struct{}{}, "br": struct{}{}, "ca": struct{}{}, "ch": struct{}{}, "cn": struct{}{}, "co": struct{}{}, "cu": struct{}{}, "cz": struct{}{}, "de": struct{}{}, "eg": struct{}{}, "fr": struct{}{}, "gb": struct{}{}, "gr": struct{}{}, "hk": struct{}{}, "hu": struct{}{}, "id": struct{}{}, "ie": struct{}{}, "il": struct{}{}, "in": struct{}{}, "it": struct{}{}, "jp": struct{}{}, "kr": struct{}{}, "lt": struct{}{}, "lv": struct{}{}, "ma": struct{}{}, "mx": struct{}{}, "my": struct{}{}, "ng": struct{}{}, "nl": struct{}{}, "no": struct{}{}, "nz": struct{}{}, "ph": struct{}{}, "pl": struct{}{}, "pt": struct{}{}, "ro": struct{}{}, "rs": struct{}{}, "ru": struct{}{}, "sa": struct{}{}, "se": struct{}{}, "sg": struct{}{}, "si": struct{}{}, "sk": struct{}{}, "th": struct{}{}, "tr": struct{}{}, "tw": struct{}{}, "ua": struct{}{}, "us": struct{}{}, "ve": struct{}{}, "za": struct{}{}}

var allowedCategories = map[string]struct{}{"business": struct{}{}, "entertainment": struct{}{}, "general": struct{}{}, "health": struct{}{}, "science": struct{}{}, "sports": struct{}{}, "technology": struct{}{}}

var allowedLanguages = map[string]struct{}{"ar": struct{}{}, "de": struct{}{}, "en": struct{}{}, "es": struct{}{}, "fr": struct{}{}, "he": struct{}{}, "it": struct{}{}, "nl": struct{}{}, "no": struct{}{}, "pt": struct{}{}, "ru": struct{}{}, "se": struct{}{}, "ud": struct{}{}, "zh": struct{}{}}

var allowedSortByOptions = map[string]struct{}{"publishedAt": struct{}{}, "relevancy": struct{}{}, "popularity": struct{}{}}

func (p *parameters) buildURL(bURL string, ap *allowedParameters) (string, error) {
	err := ap.verify(*p)
	if err != nil {
		return "", err
	}
	u := make(url.Values)
	for k, v := range *p {
		switch k {
		case "country":
			s := fmt.Sprintf("%v", v)
			if _, ok := allowedCountryCodes[s]; !ok {
				return "", fmt.Errorf("Unsupported country code %v", v)
			}
			u.Add(k, s)
		case "category":
			s := fmt.Sprintf("%v", v)
			if _, ok := allowedCategories[s]; !ok {
				return "", fmt.Errorf("Invalid category %v", v)
			}
			u.Add(k, s)
		case "language":
			s := fmt.Sprintf("%v", v)
			if _, ok := allowedLanguages[s]; !ok {
				return "", fmt.Errorf("Unsupported language %v", v)
			}
			u.Add(k, s)
		case "sortBy":
			s := fmt.Sprintf("%v", v)
			if _, ok := allowedSortByOptions[s]; !ok {
				return "", fmt.Errorf("Invalid sort by type %v", v)
			}
			u.Add(k, s)
		case "sources":
			s, l := interfaceToStringList(v)
			if l == 0 {
				return "", errors.New("Empty list of sources")
			}
			if l > 20 {
				return "", fmt.Errorf("Maximum of 20 sources got %v", l)
			}
			u.Add(k, s)
		case "domains":
			s, l := interfaceToStringList(v)
			if l == 0 {
				return "", errors.New("Empty list of domains")
			}
			u.Add(k, s)
		case "q":
			s := fmt.Sprintf("%v", v)
			if s == "" {
				return "", fmt.Errorf("Expected query got empty string")
			}
			u.Add(k, s)
		case "pageSize", "page":
			i, ok := v.(int)
			if ok {
				u.Add(k, strconv.Itoa(i))
			}
		default:
			return "", fmt.Errorf("Unhandled parameter %v", k)
		}
	}
	return (bURL + u.Encode()), nil
}

//interfaceToStringList converts an interface with a underlying type of []string to a comma seperated string
//This returns a comma seperated string and an int containing the number of values in the string
func interfaceToStringList(i interface{}) (string, int) {
	obj := reflect.ValueOf(i)
	l := obj.Len()

	var s strings.Builder

	for i := 0; i < l; i++ {
		s.WriteString(obj.Index(i).String())
		if !(i == l-1) {
			s.WriteString(",")
		}
	}
	return string(s.String()), l
}

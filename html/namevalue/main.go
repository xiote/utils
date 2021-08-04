package namevalue

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/text/encoding/korean"
	tf "golang.org/x/text/transform"
)

type NameValue struct {
	Name  string
	Value string
}

type IdValue struct {
	Id    string
	Value string
}

func GetNameValue(nameValues []NameValue, name string) (nameValue NameValue) {
	for _, item := range nameValues {
		if item.Name == name {
			nameValue = item
			return
		}
	}
	return
}

func GetValue(nameValues []NameValue, name string) (value string) {
	for _, item := range nameValues {
		if item.Name == name {
			value = item.Value
			return
		}
	}
	return
}

func UpdateNameValue(ptNameValues *[]NameValue, name string, value string) (err error) {
	for idx, item := range *ptNameValues {
		if item.Name == name {
			(*ptNameValues)[idx].Value = value
			break
		}
	}

	err = fmt.Errorf("NameValue not found!")
	return
}

func ExtractNameValue(html string, xpath string) (nameValues []NameValue, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(html))
	if err != nil {
		err = err
		return
	}

	list, err := htmlquery.QueryAll(doc, xpath)
	if err != nil {
		err = err
		return
	}

	if len(list) == 0 {
		err = fmt.Errorf("Cannot find node : %s", xpath)
		return
	}

	for _, input := range list {
		nameValues = append(nameValues, NameValue{htmlquery.SelectAttr(input, "name"), Transform(htmlquery.SelectAttr(input, "value"), "euc-kr")})
	}
	return

}

func ExtractIdValueUtf8(html string, xpath string) (idValues []IdValue, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(html))
	if err != nil {
		err = err
		return
	}

	list, err := htmlquery.QueryAll(doc, xpath)
	if err != nil {
		err = err
		return
	}

	if len(list) == 0 {
		err = fmt.Errorf("Cannot find node : %s", xpath)
		return
	}

	for _, input := range list {
		idValues = append(idValues, IdValue{htmlquery.SelectAttr(input, "id"), htmlquery.SelectAttr(input, "value")})
	}
	return

}

func Transform(source string, option string) string {
	if option == "euc-kr" {
		var bufs bytes.Buffer
		wr := tf.NewWriter(&bufs, korean.EUCKR.NewEncoder())
		defer wr.Close()
		wr.Write([]byte(source))
		return bufs.String()
	}
	return ""
}

func GetQueryString(params []NameValue) string {

	var paramsString string
	for _, input := range params {
		paramsString += fmt.Sprintf("%s=%s&", input.Name, url.QueryEscape(input.Value))
	}
	paramsString = strings.TrimSuffix(paramsString, "&")
	// paramsString = strings.ReplaceAll(paramsString, "%2B", "+")
	paramsString = strings.ReplaceAll(paramsString, "~", "%7E")
	return paramsString
}

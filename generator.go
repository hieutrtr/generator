package postgresql_generator

import (
	"fmt"
	"math/rand"
	"reflect"

	"github.com/hieutrtr/fake"
)

func GenInsertion(st interface{}) string {
	ref := reflect.ValueOf(st).Elem()
	queryFormat := "INSERT INTO users (%s) VALUES (%s);"
	var fields string
	var values string
	for j := 0; j < ref.NumField(); j++ {
		typeField := ref.Type().Field(j)
		tag := typeField.Tag
		fields += typeField.Name
		values += fmt.Sprint(getValueOfType(tag.Get("gentype")))
		if j < ref.NumField()-1 {
			fields += ","
			values += ","
		}
	}
	return fmt.Sprintf(queryFormat, fields, values)
}

func GenJSON(st interface{}) string {
	ref := reflect.ValueOf(st).Elem()
	var res string
	var fields string
	var values string
	for j := 0; j < ref.NumField(); j++ {
		typeField := ref.Type().Field(j)
		tag := typeField.Tag
		fields = typeField.Name
		values = fmt.Sprint(getValueOfType(tag.Get("gentype")))
		if j < ref.NumField()-1 {
			res += fmt.Sprintf("\"%s\":%s, ", fields, values)
		}
	}
	res += fmt.Sprintf("\"%s\":%s", fields, values)
	return fmt.Sprintf("{%s}", res)
}

func getValueOfType(tp string) interface{} {
	// types := []string{"varchar", "smallint", "int", "money", "cidr", "jsonb"}
	switch tp {
	case "ad_state":
		fake.UseExternalData(true)
		return fmt.Sprintf("\"%s\"", fake.AdState())
	case "varchar":
		return fmt.Sprintf("\"%s\"", fake.FullName())
	case "smallint":
		return uint8(rand.Uint32() % 100)
	case "int":
		return rand.Int() % 100
	case "money":
		return rand.Int31()
	case "cidr":
		return fmt.Sprintf("\"%s\"", fake.IPv4())
	case "jsonb":
		return fmt.Sprintf("'{\"company\":\"%s\",\"industry\":\"%s\",\"title\":\"%s\"}'", fake.Company(), fake.Industry(), fake.JobTitle())
	default:
		return nil
	}
}

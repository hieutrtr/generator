package postgresql-generator

import (
	"fmt"
	"math/rand"
	"reflect"

	"github.com/icrowley/fake"
)

func genInsertion(st interface{}) string {
	ref := reflect.ValueOf(st).Elem()
	queryFormat := "INSERT INTO users (%s) VALUES (%s);"
	var fields string
	var values string
	for j := 0; j < ref.NumField(); j++ {
		typeField := ref.Type().Field(j)
		tag := typeField.Tag
		fields += typeField.Name
		values += fmt.Sprint(getValueOfType(tag.Get("pgtype")))
		if j < ref.NumField()-1 {
			fields += ","
			values += ","
		}
	}
	return fmt.Sprintf(queryFormat, fields, values)
}

// NewSupplier push generated data into channel
func NewSupplier(st interface{}, num int, sup chan<- string) {
	for i := 0; i < num; i++ {
		sup <- genInsertion(st)
	}
	close(sup)
}

func getValueOfType(tp string) interface{} {
	// types := []string{"varchar", "smallint", "int", "money", "cidr", "jsonb"}
	switch tp {
	case "varchar":
		return fmt.Sprintf("'%s'", fake.FullName())
	case "smallint":
		return uint8(rand.Uint32() % 100)
	case "int":
		return rand.Int() % 100
	case "money":
		return rand.Int31()
	case "cidr":
		return fmt.Sprintf("'%s'", fake.IPv4())
	case "jsonb":
		return fmt.Sprintf("'{\"company\":\"%s\",\"industry\":\"%s\",\"title\":\"%s\"}'", fake.Company(), fake.Industry(), fake.JobTitle())
	default:
		return nil
	}
}

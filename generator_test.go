package postgresql_generator

import (
	"reflect"
	"strings"
	"testing"
)

type SomeStruct struct {
	name     string `gentype:"varchar"`
	age      uint8  `gentype:"smallint"`
	friends  int    `gentype:"int"`
	salary   int64  `gentype:"money"`
	ipv4     string `gentype:"cidr"`
	metadata string `gentype:"jsonb"`
}

func TestGenInsertion(t *testing.T) {
	q := GenInsertion(&SomeStruct{})
	if !validInsertionQuery(q) {
		t.Fatal("Invalid insert query " + q)
	}
}

func TestGetValueOfType(t *testing.T) {
	vc := getValueOfType("varchar")
	if reflect.TypeOf(vc).Kind() != reflect.String || vc == "" {
		t.Fatal("Error of generating varchar")
	}
	smi := getValueOfType("smallint")
	if reflect.TypeOf(smi).Kind() != reflect.Uint8 {
		t.Fatal("Error of generating smallint")
	}
	i := getValueOfType("int")
	if reflect.TypeOf(i).Kind() != reflect.Int {
		t.Fatal("Error of generating int")
	}
	mone := getValueOfType("money")
	if reflect.TypeOf(mone).Kind() != reflect.Int32 {
		t.Fatal("Error of generating money")
	}
	ipv4 := getValueOfType("cidr")
	if reflect.TypeOf(ipv4).Kind() != reflect.String || ipv4 == "" {
		t.Fatal("Error of generating cidr")
	}
	jsb := getValueOfType("jsonb")
	if reflect.TypeOf(jsb).Kind() != reflect.String || jsb == "" {
		t.Fatal("Error of generating jsonb")
	}
	shouldNil := getValueOfType("not-any-type")
	if shouldNil != nil {
		t.Fatal("Error of default generating")
	}
}

func validInsertionQuery(q string) bool {
	return strings.HasPrefix(q, "INSERT INTO")
}

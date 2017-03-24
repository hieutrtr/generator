package postgresql_generator

import "testing"

var pgCtrl *PGCtrl
var cfg *Config

func TestNewPG(t *testing.T) {
	var err error
	pgCtrl, err = NewPG()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestConnect(t *testing.T) {
	cfg = &Config{}
	err := cfg.Parse()
	if err == nil {
		err = pgCtrl.Connect(cfg)
		if err != nil {
			t.Fatal(err.Error())
		}
	}
}

func TestExecute(t *testing.T) {
	err := cfg.Parse()
	if err == nil {
		err = pgCtrl.Execute("INSERT INTO users (name) VALUES('test');")
		if err != nil {
			t.Fatal(err.Error())
		}
		err = pgCtrl.Execute("DELETE FROM users WHERE name='test';")
	}
}

package postgresql_generator

import "testing"

func TestBuildUploadEvent(t *testing.T) {
	e := &UploadEvent{
		Topic:   "test-imgevent",
		Payload: "test-id",
	}
	_, err := e.buildEvent()
	if err != nil {
		t.Fatal("Building upload event was fail on error", err.Error())
	}
}

func TestCreateProducer(t *testing.T) {
	_, err := NewProducer()
	if err != nil {
		t.Fatal("Producing event to kafka was fail on error", err.Error())
	}
}

func TestProduceUploadEvent(t *testing.T) {
	p, _ := NewProducer()
	e := &UploadEvent{
		Topic:   "test-imgevent",
		Payload: "test-id",
	}
	err := p.Produce(e)
	if err != nil {
		t.Fatal("Producing event to kafka was fail on error", err.Error())
	}
}

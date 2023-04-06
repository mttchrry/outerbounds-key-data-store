package keyDataStore

import (
	"testing"

	"github.com/google/uuid"
)

func TestPhoneNumbers_Ingest(t *testing.T) {
	kds, err := New("../../example.data") 
	if err != nil {
		t.Fatalf("failed initializing KDS: %v", err)
	}
	if len(kds.Data) != 10000 {
		t.Fatalf("expecting 10000 vals, got : %v", len(kds.Data))
	}

	key, _ := uuid.Parse("0e6f37fb-506d-47bd-a2d2-2046b1a35627")
	val := kds.Data[key]
	if val != "wazaugt wqv" {
		t.Fatalf("incorrect data, got: %v", val)
	}
}

func TestPhoneNumbers_IngestFail(t *testing.T) {
	_, err := New("notExist.data") 

	if err == nil {
		t.Fatalf("expecting missing file error")
	}
}
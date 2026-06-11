package showrunner

import (
	"encoding/json"
	"testing"
)

func TestFlexibleStringListUnmarshalString(t *testing.T) {
	var list FlexibleStringList
	if err := json.Unmarshal([]byte(`"calm, controlled"`), &list); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if len(list) != 1 || list[0] != "calm, controlled" {
		t.Fatalf("list = %#v", list)
	}
}

func TestFlexibleStringListUnmarshalArray(t *testing.T) {
	var list FlexibleStringList
	if err := json.Unmarshal([]byte(`["calm","controlled"]`), &list); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if len(list) != 2 || list[0] != "calm" || list[1] != "controlled" {
		t.Fatalf("list = %#v", list)
	}
}

func TestFlexibleStringListUnmarshalNull(t *testing.T) {
	var list FlexibleStringList
	if err := json.Unmarshal([]byte(`null`), &list); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if list == nil || len(list) != 0 {
		t.Fatalf("list = %#v, want non-nil empty list", list)
	}
}

func TestFlexibleStringListRejectsOtherTypes(t *testing.T) {
	var list FlexibleStringList
	if err := json.Unmarshal([]byte(`42`), &list); err == nil {
		t.Fatal("Unmarshal() error = nil, want clear type error")
	}
}

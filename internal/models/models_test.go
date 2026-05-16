package models

import (
	"encoding/json"
	"testing"
)

func TestWorkItemJSON(t *testing.T) {
	wi := WorkItem{
		ID: "abc-123", Name: "Test Item", State: "Backlog",
		Priority: "high", Assignees: []string{"dev@test.com"}, Labels: []string{"bug"},
	}
	data, _ := json.Marshal(wi)
	var result WorkItem
	json.Unmarshal(data, &result)
	if result.ID != "abc-123" || result.Name != "Test Item" {
		t.Errorf("Got %+v", result)
	}
}

func TestCycleJSON(t *testing.T) {
	c := Cycle{ID: "cyc-1", Name: "Sprint 1", StartDate: "2026-01-01"}
	data, _ := json.Marshal(c)
	var result Cycle
	json.Unmarshal(data, &result)
	if result.ID != "cyc-1" {
		t.Errorf("Got %+v", result)
	}
}

func TestLabelJSON(t *testing.T) {
	l := Label{ID: "lbl-1", Name: "bug", Color: "#ff0000"}
	data, _ := json.Marshal(l)
	var result Label
	json.Unmarshal(data, &result)
	if result.Color != "#ff0000" {
		t.Errorf("Got %+v", result)
	}
}

func TestMemberJSON(t *testing.T) {
	m := Member{ID: "usr-1", Email: "dev@test.com", DisplayName: "Dev", Role: "admin"}
	data, _ := json.Marshal(m)
	var result Member
	json.Unmarshal(data, &result)
	if result.Email != "dev@test.com" {
		t.Errorf("Got %+v", result)
	}
}

func TestPageJSON(t *testing.T) {
	p := Page{ID: "pg-1", Name: "Home", Content: "# Welcome"}
	data, _ := json.Marshal(p)
	var result Page
	json.Unmarshal(data, &result)
	if result.Content != "# Welcome" {
		t.Errorf("Got %+v", result)
	}
}

func TestCommentJSON(t *testing.T) {
	c := Comment{ID: "cmt-1", Content: "Looks good", CreatedBy: "dev@test.com"}
	data, _ := json.Marshal(c)
	var result Comment
	json.Unmarshal(data, &result)
	if result.Content != "Looks good" {
		t.Errorf("Got %+v", result)
	}
}

func TestCollectionJSON(t *testing.T) {
	coll := Collection{Items: []interface{}{"a", "b"}, Total: 2}
	data, _ := json.Marshal(coll)
	var result Collection
	json.Unmarshal(data, &result)
	if result.Total != 2 {
		t.Errorf("Got total %d", result.Total)
	}
}

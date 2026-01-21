package generator

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestModelTemplate_Standard_WithMethods(t *testing.T) {
	// Setup data
	data := &ModelData{
		Package:   "com.example",
		ClassName: "User",
		Fields: []*ModelField{
			{FieldName: "id", FieldType: "Long", ColumnName: "id"},
			{FieldName: "name", FieldType: "String", ColumnName: "name"},
		},
		NeedConstructors:           true,
		NeedToStringHashcodeEquals: true,
	}

	tmpl, err := template.New("model").Funcs(TemplateFuncs).Parse(modelTemplate)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	// Verify constructors
	assert.Contains(t, output, "public User() {}", "Should contain no-arg constructor")
	assert.Contains(t, output, "public User(Long id, String name)", "Should contain all-args constructor")
	assert.Contains(t, output, "this.id = id;", "Constructor should assign fields")

	// Verify toString/hashCode/equals
	assert.Contains(t, output, "public boolean equals(Object o)", "Should contain equals")
	assert.Contains(t, output, "public int hashCode()", "Should contain hashCode")
	assert.Contains(t, output, "public String toString()", "Should contain toString")
	assert.Contains(t, output, "Objects.equals", "Should use Objects.equals")
	assert.Contains(t, output, "Objects.hash", "Should use Objects.hash")
}

func TestModelTemplate_Standard_NoMethods(t *testing.T) {
	// Setup data
	data := &ModelData{
		Package:   "com.example",
		ClassName: "User",
		Fields: []*ModelField{
			{FieldName: "id", FieldType: "Long", ColumnName: "id"},
		},
		NeedConstructors:           false,
		NeedToStringHashcodeEquals: false,
	}

	tmpl, err := template.New("model").Funcs(TemplateFuncs).Parse(modelTemplate)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	assert.NotContains(t, output, "public User() {}", "Should NOT contain constructors")
	assert.NotContains(t, output, "equals(Object o)", "Should NOT contain equals")
}

func TestModelTemplate_Lombok(t *testing.T) {
	// Setup data
	data := &ModelData{
		Package:          "com.example",
		ClassName:        "User",
		Fields:           []*ModelField{},
		NeedConstructors: true,
	}

	tmpl, err := template.New("model").Funcs(TemplateFuncs).Parse(modelLombokTemplate)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	assert.Contains(t, output, "@NoArgsConstructor", "Should contain @NoArgsConstructor")
	assert.Contains(t, output, "@AllArgsConstructor", "Should contain @AllArgsConstructor")
}

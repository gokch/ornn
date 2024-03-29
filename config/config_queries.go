package config

import (
	"fmt"

	"ariga.io/atlas/sql/schema"
	"github.com/gokch/ornn/config/atlas"
)

type Queries struct {
	schema *Schema             `json:"-"`
	Class  map[string][]*Query `json:"class"` // auto generated by schema
}

func (t *Queries) init(schema *Schema) {
	t.schema = schema
	if t.Class == nil {
		t.Class = make(map[string][]*Query)
		t.InitDefaultQueryTables()
	}
}

func (t *Queries) InitDefaultQueryTables() error {
	for _, table := range t.schema.Tables {
		err := t.initDefaultQueryTable(table)
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO : 작업 예정
func (t *Queries) initDefaultQueryTable(table *schema.Table) error {
	// insert all
	var insertQuestionare string
	for i, col := range table.Columns {
		if table.PrimaryKey != nil && col.Name == table.PrimaryKey.Name {
			insertQuestionare += "NULL, "
		} else {
			if t.schema.DbType == atlas.DbTypePostgre || t.schema.DbType == atlas.DbTypeCockroachDB {
				insertQuestionare += fmt.Sprintf("$%d, ", i+1)
			} else {
				insertQuestionare += "?, "
			}
		}
	}
	insertQuestionare = insertQuestionare[:len(insertQuestionare)-2]
	t.AddQuery(table.Name, &Query{
		Name:    "insert",
		Comment: "default query - insert",
		Sql:     fmt.Sprintf("INSERT INTO %s VALUES (%s)", table.Name, insertQuestionare),
	})

	// where
	var i int
	var where string
	if table.PrimaryKey != nil && len(table.PrimaryKey.Parts) == 1 {
		pkName := table.PrimaryKey.Parts[0].C.Name // TODO
		if pkName != "" {
			if t.schema.DbType == atlas.DbTypePostgre || t.schema.DbType == atlas.DbTypeCockroachDB {
				where = fmt.Sprintf(" WHERE %s = $%d", pkName, i+1)
			} else {
				where = fmt.Sprintf(" WHERE %s = ?", pkName)
			}
		}
	}

	// select
	t.AddQuery(table.Name, &Query{
		Name:    "select",
		Comment: "default query - select",
		Sql:     fmt.Sprintf("SELECT * FROM %s%s", table.Name, where),
	})

	// delete
	t.AddQuery(table.Name, &Query{
		Name:    "delete",
		Comment: "default query - delete",
		Sql:     fmt.Sprintf("DELETE FROM %s%s", table.Name, where),
	})

	// set
	setQuestionaire := ""
	var col *schema.Column
	for i, col = range table.Columns {
		if col.Name == table.PrimaryKey.Name {
			continue
		}
		if t.schema.DbType == atlas.DbTypePostgre || t.schema.DbType == atlas.DbTypeCockroachDB {
			setQuestionaire += fmt.Sprintf("%s = $%d, ", col.Name, i+1)
		} else {
			setQuestionaire += fmt.Sprintf("%s = ?, ", col.Name)
		}
	}
	setQuestionaire = setQuestionaire[:len(setQuestionaire)-2]

	// where ( update )
	if table.PrimaryKey != nil && len(table.PrimaryKey.Parts) == 1 {
		pkName := table.PrimaryKey.Parts[0].C.Name // TODO
		if pkName != "" {
			if t.schema.DbType == atlas.DbTypePostgre || t.schema.DbType == atlas.DbTypeCockroachDB {
				where = fmt.Sprintf(" WHERE %s = $%d", pkName, i+2)
			} else {
				where = fmt.Sprintf(" WHERE %s = ?", pkName)
			}
		}
	}
	// update
	t.AddQuery(table.Name, &Query{
		Name:    "update",
		Comment: "default query - update",
		Sql:     fmt.Sprintf("UPDATE %s SET %s%s", table.Name, setQuestionaire, where),
	})

	return nil
}

func (t *Queries) AddQuery(tableName string, query *Query) {
	if t.Class == nil {
		t.Class = make(map[string][]*Query, 10)
	}
	if _, ok := t.Class[tableName]; ok == false {
		t.Class[tableName] = make([]*Query, 0, 10)
	}
	t.Class[tableName] = append(t.Class[tableName], query)
}

//------------------------------------------------------------------------------------------------//
// query

type Query struct {
	Name    string `json:"name"`
	Comment string `json:"comment,omitempty"`
	Sql     string `json:"sql"`

	// options
	CustomFieldTypes []*CustomFieldType `json:"custom_field_types,omitempty"`
	InsertMulti      bool               `json:"insert_multi,omitempty"`
	UpdateNullIgnore bool               `json:"update_null_ignore,omitempty"`
	ErrQuery         string             `json:"-"`
	ErrParser        string             `json:"-"`
}

type CustomFieldType struct {
	TableName  string `json:"table_name"`
	FieldName  string `json:"field_name"`
	CustomType string `json:"type"`
}

//------------------------------------------------------------------------------------------------//
// query

func (t *Query) AddCustomType(tableName, fieldName string, customType string) {
	if t.CustomFieldTypes == nil {
		t.CustomFieldTypes = make([]*CustomFieldType, 0, 10)
	}
	t.CustomFieldTypes = append(t.CustomFieldTypes, &CustomFieldType{
		TableName:  tableName,
		FieldName:  fieldName,
		CustomType: customType,
	})
}

func (t *Query) GetCustomType(fieldName string) (genType string) {
	for _, pt := range t.CustomFieldTypes {
		if pt.FieldName == fieldName {
			return pt.CustomType
		}
	}
	return ""
}

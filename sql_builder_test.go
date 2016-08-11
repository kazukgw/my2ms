package main

import (
	"testing"

	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/schema"
	"github.com/stretchr/testify/assert"
)

func TestBuildInsert(t *testing.T) {
	a := assert.New(t)
	sb := &SqlBuilder{}
	mmaps := NewMigrationMaps()
	mmaps.Add("test_table", "test_table")
	sb.MigrationMaps = mmaps
	tbl := GenTestTable()
	rows := [][]interface{}{
		{"0000-0001", 0, "a", "abc"},
		{"0000-0002", 20, "b", "abc"},
		{"0000-0003", 100, "c", "abc"},
		{"0000-0004", 312, "d", "abc"},
	}
	e := GenRowsEvent(tbl, canal.InsertAction, rows)
	sqlsets := sb.BuildInsert(e)

	a.Equal(1, len(*sqlsets))

	expected := "INSERT INTO test_table (pk1,col1,col2,col3) " +
		"VALUES (?,?,?,?),(?,?,?,?),(?,?,?,?),(?,?,?,?)"
	a.Equal(expected, (*sqlsets)[0].Sql)
}

func TestBuildUpdate(t *testing.T) {
	a := assert.New(t)
	sb := &SqlBuilder{}
	mmaps := NewMigrationMaps()
	mmaps.Add("test_table", "test_table")
	sb.MigrationMaps = mmaps
	tbl := GenTestTable()
	rows := [][]interface{}{
		[]interface{}{
			[]interface{}{"0000-0001", 0, "a", "abc"},
			[]interface{}{"0000-0001", 100, "a", "abc"},
		},
		[]interface{}{
			[]interface{}{"0000-0002", 20, "b", "abc"},
			[]interface{}{"0000-0002", 30, "z", "abc"},
		},
		[]interface{}{
			[]interface{}{"0000-0003", 100, "c", "abc"},
			[]interface{}{"0000-0003", 120, "d", "abc"},
		},
		[]interface{}{
			[]interface{}{"0000-0004", 312, "d", "abc"},
			[]interface{}{"0000-0004", 632, "ss", "abc"},
		},
	}
	e := GenRowsEvent(tbl, canal.UpdateAction, rows)
	sqlsetsPtr := sb.BuildUpdate(e)
	sqlsets := *sqlsetsPtr

	a.Equal(4, len(sqlsets))

	var expected string
	expected = "UPDATE test_table SET col1 = ? WHERE pk1 = ?"
	a.Equal(expected, sqlsets[0].Sql)

	expected = "UPDATE test_table " +
		"SET col1 = ?, col2 = ? WHERE pk1 = ?"
	a.Equal(expected, sqlsets[1].Sql)
}

func GenRowsEvent(
	tbl *schema.Table,
	action string,
	rows [][]interface{},
) *canal.RowsEvent {
	e := &canal.RowsEvent{}
	e.Table = tbl
	e.Action = action
	e.Rows = rows
	return e
}

func GenTestTable() *schema.Table {
	tbl := schema.Table{}
	tbl.Schema = "test"
	tbl.Name = "test_table"
	tbl.AddColumn("pk1", "varchar(255)", "")
	tbl.PKColumns = []int{0}
	tbl.AddColumn("col1", "int(11)", "")
	tbl.AddColumn("col2", "varchar(255)", "")
	tbl.AddColumn("col3", "varchar(255)", "")
	idx := tbl.AddIndex("test_table_index_col1_col2")
	idx.Columns = []string{"col1", "col2"}
	return &tbl
}

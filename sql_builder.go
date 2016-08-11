package main

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/davecgh/go-spew/spew"
	"github.com/siddontang/go-mysql/canal"
)

type SqlSet struct {
	Sql   string
	Args  []interface{}
	Error error
}

type SqlSets []SqlSet

type SqlBuilder struct {
	*MigrationMaps
}

func (sb *SqlBuilder) BuildSql(e *canal.RowsEvent) *SqlSets {
	spew.Println("rows:", e.Rows)
	switch e.Action {
	case canal.InsertAction:
		return sb.BuildInsert(e)
	case canal.UpdateAction:
		return sb.BuildUpdate(e)
	case canal.DeleteAction:
		return sb.BuildDelete(e)
	default:
		panic("unknown canal rows event type")
	}
}

func (sb *SqlBuilder) BuildInsert(e *canal.RowsEvent) *SqlSets {
	mp := sb.MigrationMaps.Get(e.Table)
	builder := sq.Insert(mp.ToTableName)
	builder = builder.Columns(mp.FromTableCols()...)
	for _, r := range e.Rows {
		builder = builder.Values(r...)
	}
	sqlstr, args, err := builder.ToSql()

	return &SqlSets{SqlSet{sqlstr, args, err}}
}

func (sb *SqlBuilder) BuildUpdate(e *canal.RowsEvent) *SqlSets {
	mp := sb.MigrationMaps.Get(e.Table)
	sqlsets := SqlSets{}

	builder := sq.Update(mp.ToTableName)

	// update event rows
	// [a_before, a_after, b_before, b_after, c_before, c_after  ...]
	for i := 0; i < len(e.Rows); i += 2 {
		before := e.Rows[i]
		after := e.Rows[i+1]

		updateIdxs := []int{}
		for i, val := range before {
			if val != after[i] {
				updateIdxs = append(updateIdxs, i)
			}
		}
		fromcols := mp.FromTableCols()
		for _, idx := range updateIdxs {
			builder = builder.Set(fromcols[idx], after[idx])
		}

		pks := mp.PkValue(before)
		builder = builder.Where(sq.Eq(pks))
		sqlstr, args, err := builder.ToSql()
		sqlsets = append(sqlsets, SqlSet{sqlstr, args, err})
	}

	return &sqlsets
}

func (sb *SqlBuilder) BuildDelete(e *canal.RowsEvent) *SqlSets {
	mp := sb.MigrationMaps.Get(e.Table)
	sqlsets := SqlSets{}
	for _, r := range e.Rows {
		builder := sq.Delete(mp.ToTableName)
		pks := mp.PkValue(r)
		builder = builder.Where(sq.Eq(pks))
		sqlstr, args, err := builder.ToSql()
		sqlsets = append(sqlsets, SqlSet{sqlstr, args, err})
	}
	return &sqlsets

}

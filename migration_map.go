package main

import "github.com/siddontang/go-mysql/schema"

type MigrationMaps struct {
	maps map[string]*MigrationMap
}

func NewMigrationMaps() *MigrationMaps {
	mmaps := &MigrationMaps{}
	mmaps.maps = make(map[string]*MigrationMap)
	return mmaps
}

func NewMigrationMapsWithMap(mmap map[string]string) *MigrationMaps {
	mmaps := NewMigrationMaps()
	for from, to := range mmap {
		mmaps.Add(from, to)
	}
	return mmaps
}

func (mmaps *MigrationMaps) Add(from, to string) {
	mmap := &MigrationMap{}
	mmap.FromTableName = from
	mmap.ToTableName = to
	mmaps.maps[from] = mmap
}

func (mmaps *MigrationMaps) Get(t *schema.Table) *MigrationMap {
	if m, ok := mmaps.maps[t.Name]; ok {
		if m.FromTable == nil {
			m.FromTable = t
		}
		return m
	}
	return nil
}

type MigrationMap struct {
	FromTableName string
	ToTableName   string
	FromTable     *schema.Table
}

func (mp *MigrationMap) FromTableCols() []string {
	cols := []string{}
	for _, c := range mp.FromTable.Columns {
		cols = append(cols, c.Name)
	}
	return cols
}

func (mp *MigrationMap) ColValueMap(r []interface{}) map[string]interface{} {
	cvmap := map[string]interface{}{}
	for i, col := range mp.FromTable.Columns {
		cvmap[col.Name] = r[i]
	}
	return cvmap
}

func (mp *MigrationMap) PkValue(r []interface{}) map[string]interface{} {
	pks := map[string]interface{}{}
	for _, pkindex := range mp.FromTable.PKColumns {
		col := mp.FromTable.Columns[pkindex]
		pks[col.Name] = r[pkindex]
	}
	return pks
}

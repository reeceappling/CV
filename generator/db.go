package main

import "strings"

var dbs = map[string]*DbPage{}

type DbPage struct {
	Name string
	*tracked
}

func (pg *DbPage) Link() string {
	if pg == nil {
		return "NO_LINK"
	}
	return linkFor(pg.Name, "cv", "db", withoutSpaces(pg.Name))
}

func (pg *DbPage) EntryType() string {
	return "Database"
}

func NewDb(name string) *DbPage {
	out := &DbPage{
		Name:    name,
		tracked: newTracked(),
	}
	//if slices.Contains([]string{"Aurora","RDS","DynamoDB","DocumentDB"}, name) {
	//	out.EquivalentLink = cloudServices[name].Link() // TODO: ?????????
	//}
	dbs[name] = out
	addLinkable(strings.ToLower(name), out) // TODO: will services like DynamoDB overwrite this?
	return out
}

//
//func dbIsOfAnotherType(dbName string) (bool, string){ // TODO: USE THIS?!
//	otherType, isDiff := map[string]string{
//		"Aurora": "Cloud Service",
//		"DynamoDB": "Cloud Service",
//		"DocumentDB": "Cloud Service",
//		"RDS": "Cloud Service",
//	}[dbName]
//	return isDiff, otherType
//}

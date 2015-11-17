package cassandra

import (
	"fmt"
	"log"
	"strings"

	"github.com/dollarshaveclub/go-lib/set"
	"github.com/gocql/gocql"
)

// CTable is a Cassandra table definition
type CTable struct {
	Name    string
	Columns []string
}

// UDT is a user-defined type
type UDT struct {
	Name    string
	Columns []string
}

// CreateTable creates a table
func CreateTable(c *gocql.ClusterConfig, t CTable) error {
	qs := fmt.Sprintf("CREATE TABLE IF NOT EXISTS%v ( %v );", t.Name, strings.Join(t.Columns, ", "))
	s, err := c.CreateSession()
	if err != nil {
		return err
	}
	defer s.Close()
	log.Printf("Creating table %v\n", t.Name)
	return s.Query(qs).Exec()
}

// CreateUDT creates a user-defined type
func CreateUDT(c *gocql.ClusterConfig, u UDT) error {
	qs := fmt.Sprintf("CREATE TYPE IF NOT EXISTS %v ( %v );", u.Name, strings.Join(u.Columns, ", "))
	s, err := c.CreateSession()
	if err != nil {
		return err
	}
	defer s.Close()
	log.Printf("Creating UDT %v\n", u.Name)
	return s.Query(qs).Exec()
}

// CreateRequiredTypes ensures all the types passed in are created if necessary
func CreateRequiredTypes(c *gocql.ClusterConfig, rt []UDT) error {
	for _, u := range rt {
		err := CreateUDT(c, u)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateRequiredTables ensures all the tables passed in are created if necessary
func CreateRequiredTables(c *gocql.ClusterConfig, rt []CTable) error {
	tl, err := GetTables(c)
	if err != nil {
		return err
	}

	tm := map[string]CTable{}
	rtl := []string{}
	for _, v := range rt {
		tm[v.Name] = v
		rtl = append(rtl, v.Name)
	}
	tset := set.NewStringSet(tl)
	rset := set.NewStringSet(rtl)
	diff := rset.Difference(tset)
	missing := diff.Items()
	if len(missing) > 0 {
		for _, t := range missing {
			ts := tm[t]
			err = CreateTable(c, ts)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CreateKeyspace creates a keyspace if nessary
func CreateKeyspace(c *gocql.ClusterConfig, ks string, rf string) error {
	kis, err := GetKeyspaces(c)
	if err != nil {
		return err
	}
	kss := set.NewStringSet(kis)
	if !kss.Contains(ks) {
		log.Printf("Creating keyspace: %v\n", ks)
		c.Keyspace = ""
		s, err := c.CreateSession()
		if err != nil {
			return err
		}
		defer s.Close()
		err = s.Query(fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %v WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': %v};", ks, rf)).Exec()
		if err != nil {
			return err
		}
	}
	c.Keyspace = ks
	return nil
}

// DropKeyspace deletes a keyspace and all data associated with it
func DropKeyspace(c *gocql.ClusterConfig, ks string) error {
	c.Keyspace = ""
	s, err := c.CreateSession()
	if err != nil {
		return err
	}
	defer s.Close()
	log.Printf("Dropping keyspace: %v\n", ks)
	err = s.Query(fmt.Sprintf("DROP KEYSPACE IF EXISTS %v\n", ks)).Exec()
	if err != nil {
		return err
	}
	return nil
}

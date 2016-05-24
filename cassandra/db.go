package cassandra

import (
	"github.com/gocql/gocql"
	consul "github.com/hashicorp/consul/api"
)

// GetTables returns all tables in configured keyspace
func GetTables(c *gocql.ClusterConfig) ([]string, error) {
	tables := []string{}
	s, err := c.CreateSession()
	if err != nil {
		return tables, err
	}
	defer s.Close()
	q := s.Query("SELECT columnfamily_name FROM system.schema_columnfamilies WHERE keyspace_name = ?;", c.Keyspace).Iter()
	var tn string
	for q.Scan(&tn) {
		tables = append(tables, tn)
	}
	return tables, q.Close()
}

// GetKeyspaces returns all extant keyspaces
func GetKeyspaces(c *gocql.ClusterConfig) ([]string, error) {
	kss := []string{}
	s, err := c.CreateSession()
	if err != nil {
		return kss, err
	}
	defer s.Close()
	q := s.Query("SELECT keyspace_name FROM system.schema_keyspaces;").Iter()
	var kn string
	for q.Scan(&kn) {
		kss = append(kss, kn)
	}
	return kss, q.Close()
}

// GetNodesFromConsul queries the local Consul agent for the "cassandra" service,
// returning the healthy nodes in ascending order of network distance/latency
func GetNodesFromConsul() ([]string, error) {
	nodes := []string{}
	c, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return nodes, err
	}
	h := c.Health()
	opts := &consul.QueryOptions{
		Near: "_agent",
	}
	se, _, err := h.Service("cassandra", "", true, opts)
	if err != nil {
		return nodes, err
	}
	for _, s := range se {
		nodes = append(nodes, s.Node.Address)
	}
	return nodes, nil
}

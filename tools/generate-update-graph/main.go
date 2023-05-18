package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/blang/semver/v4"
	"sigs.k8s.io/yaml"

	"github.com/authzed/spicedb-operator/pkg/config"
	"github.com/authzed/spicedb-operator/pkg/updates"
)

//go:generate go run main.go ../../proposed-update-graph.yaml

func main() {
	if len(os.Args) != 2 {
		fmt.Println("must provide filename")
		os.Exit(1)
	}

	opconfig := config.OperatorConfig{
		ImageName: "ghcr.io/authzed/spicedb",
		UpdateGraph: updates.UpdateGraph{Channels: []updates.Channel{
			postgresChannel(),
			crdbChannel(),
			mysqlChannel(),
			spannerChannel(),
			memoryChannel(),
		}},
	}

	yamlBytes, err := yaml.Marshal(&opconfig)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(os.Args[1], yamlBytes, 0o666); err != nil {
		panic(err)
	}
}

func postgresChannel() updates.Channel {
	releases := []updates.State{
		{ID: "v1.21.0", Tag: "v1.21.0", Migration: "add-gc-covering-index"},
		{ID: "v1.19.1", Tag: "v1.19.1", Migration: "add-gc-covering-index"},
		{ID: "v1.18.0", Tag: "v1.18.0", Migration: "drop-bigserial-ids"},
		{ID: "v1.17.0", Tag: "v1.17.0", Migration: "drop-bigserial-ids"},
		{ID: "v1.16.2", Tag: "v1.16.2", Migration: "drop-bigserial-ids"},
		{ID: "v1.16.1", Tag: "v1.16.1", Migration: "drop-bigserial-ids", Deprecated: true},
		{ID: "v1.16.0", Tag: "v1.16.0", Migration: "drop-bigserial-ids", Deprecated: true},
		{ID: "v1.15.0", Tag: "v1.15.0", Migration: "drop-bigserial-ids"},
		{ID: "v1.14.1", Tag: "v1.14.1", Migration: "drop-bigserial-ids"},
		{ID: "v1.14.0", Tag: "v1.14.0", Migration: "drop-bigserial-ids"},
		{ID: "v1.14.0-phase2", Tag: "v1.14.0", Migration: "add-xid-constraints", Phase: "write-both-read-new"},
		{ID: "v1.14.0-phase1", Tag: "v1.14.0", Migration: "add-xid-columns", Phase: "write-both-read-old"},
		{ID: "v1.13.0", Tag: "v1.13.0", Migration: "add-ns-config-id"},
		{ID: "v1.12.0", Tag: "v1.12.0", Migration: "add-ns-config-id"},
		{ID: "v1.11.0", Tag: "v1.11.0", Migration: "add-ns-config-id"},
		{ID: "v1.10.0", Tag: "v1.10.0", Migration: "add-ns-config-id"},
		{ID: "v1.9.0", Tag: "v1.9.0", Migration: "add-unique-datastore-id"},
		{ID: "v1.8.0", Tag: "v1.8.0", Migration: "add-unique-datastore-id"},
		{ID: "v1.7.1", Tag: "v1.7.1", Migration: "add-unique-datastore-id"},
		{ID: "v1.7.0", Tag: "v1.7.0", Migration: "add-unique-datastore-id", Deprecated: true},
		{ID: "v1.6.0", Tag: "v1.6.0", Migration: "add-unique-datastore-id"},
		{ID: "v1.5.0", Tag: "v1.5.0", Migration: "add-transaction-timestamp-index"},
		{ID: "v1.4.0", Tag: "v1.4.0", Migration: "add-transaction-timestamp-index"},
		{ID: "v1.3.0", Tag: "v1.3.0", Migration: "add-transaction-timestamp-index"},
		{ID: "v1.2.0", Tag: "v1.2.0", Migration: "add-transaction-timestamp-index"},
	}
	edgePatterns := map[string]string{
		"v1.19.1":        ">=1.21.0",
		"v1.18.0":        ">=1.19.1",
		"v1.17.0":        ">=1.18.0",
		"v1.16.2":        ">=1.17.0",
		"v1.16.1":        ">=1.16.2",
		"v1.16.0":        ">=1.16.2",
		"v1.15.0":        ">=1.16.2",
		"v1.14.1":        ">=1.15.0",
		"v1.14.0":        ">=1.14.1",
		"v1.14.0-phase2": "1.14.0",
		"v1.14.0-phase1": "1.14.0-phase2",
		"v1.13.0":        "1.14.0-phase1",
		"v1.12.0":        ">=1.13.0 <=1.14.0-phase1",
		"v1.11.0":        ">=1.12.0 <=1.14.0-phase1",
		"v1.10.0":        ">=1.11.0 <=1.14.0-phase1",
		"v1.9.0":         ">=1.10.0 <=1.14.0-phase1",
		"v1.8.0":         ">=1.9.0 <=1.14.0-phase1",
		"v1.7.1":         ">=1.8.0 <=1.14.0-phase1",
		"v1.7.0":         ">=1.7.1 <=1.14.0-phase1",
		"v1.6.0":         ">=1.7.1 <=1.14.0-phase1",
		"v1.5.0":         ">=1.6.0 <=1.14.0-phase1",
		"v1.4.0":         ">=1.5.0 <=1.14.0-phase1",
		"v1.3.0":         ">=1.4.0 <=1.14.0-phase1",
		"v1.2.0":         ">=1.3.0 <=1.14.0-phase1",
	}

	return updates.Channel{
		Name: "stable",
		Metadata: map[string]string{
			"datastore": "postgres",
			"default":   "true",
		},
		Nodes: releases,
		Edges: edgesFromPatterns(edgePatterns, releases),
	}
}

func crdbChannel() updates.Channel {
	releases := []updates.State{
		{ID: "v1.21.0", Tag: "v1.21.0", Migration: "add-caveats"},
		{ID: "v1.19.1", Tag: "v1.19.1", Migration: "add-caveats"},
		{ID: "v1.18.0", Tag: "v1.18.0", Migration: "add-caveats"},
		{ID: "v1.17.0", Tag: "v1.17.0", Migration: "add-caveats"},
		{ID: "v1.16.2", Tag: "v1.16.2", Migration: "add-caveats"},
		{ID: "v1.16.1", Tag: "v1.16.1", Migration: "add-caveats", Deprecated: true},
		{ID: "v1.16.0", Tag: "v1.16.0", Migration: "add-caveats", Deprecated: true},
		{ID: "v1.15.0", Tag: "v1.15.0", Migration: "add-caveats"},
		{ID: "v1.14.1", Tag: "v1.14.1", Migration: "add-caveats"},
		{ID: "v1.14.0", Tag: "v1.14.0", Migration: "add-caveats", Deprecated: true},
		{ID: "v1.13.0", Tag: "v1.13.0", Migration: "add-metadata-and-counters"},
		{ID: "v1.12.0", Tag: "v1.12.0", Migration: "add-metadata-and-counters"},
		{ID: "v1.11.0", Tag: "v1.11.0", Migration: "add-metadata-and-counters"},
		{ID: "v1.10.0", Tag: "v1.10.0", Migration: "add-metadata-and-counters"},
		{ID: "v1.9.0", Tag: "v1.9.0", Migration: "add-metadata-and-counters"},
		{ID: "v1.8.0", Tag: "v1.8.0", Migration: "add-metadata-and-counters"},
		{ID: "v1.7.1", Tag: "v1.7.1", Migration: "add-metadata-and-counters"},
		{ID: "v1.7.0", Tag: "v1.7.0", Migration: "add-metadata-and-counters", Deprecated: true},
		{ID: "v1.6.0", Tag: "v1.6.0", Migration: "add-metadata-and-counters"},
		{ID: "v1.5.0", Tag: "v1.5.0", Migration: "add-transactions-table"},
		{ID: "v1.4.0", Tag: "v1.4.0", Migration: "add-transactions-table"},
		{ID: "v1.3.0", Tag: "v1.3.0", Migration: "add-transactions-table"},
		{ID: "v1.2.0", Tag: "v1.2.0", Migration: "add-transactions-table"},
	}
	edgePatterns := map[string]string{
		"v1.19.1": ">=1.21.0",
		"v1.18.0": ">=1.19.1",
		"v1.17.0": ">=1.18.0",
		"v1.16.2": ">=1.17.0",
		"v1.16.1": ">=1.16.2",
		"v1.16.0": ">=1.16.2",
		"v1.15.0": ">=1.16.2",
		"v1.14.1": ">=1.15.0",
		"v1.14.0": ">=1.14.1",
		"v1.13.0": ">=1.14.1",
		"v1.12.0": ">=1.13.0",
		"v1.11.0": ">=1.12.0",
		"v1.10.0": ">=1.11.0",
		"v1.9.0":  ">=1.10.0",
		"v1.8.0":  ">=1.9.0",
		"v1.7.1":  ">=1.8.0",
		"v1.7.0":  ">=1.7.1",
		"v1.6.0":  ">=1.7.1",
		"v1.5.0":  ">=1.6.0",
		"v1.4.0":  ">=1.5.0",
		"v1.3.0":  ">=1.4.0",
		"v1.2.0":  ">=1.3.0",
	}
	return updates.Channel{
		Name: "stable",
		Metadata: map[string]string{
			"datastore": "cockroachdb",
			"default":   "true",
		},
		Nodes: releases,
		Edges: edgesFromPatterns(edgePatterns, releases),
	}
}

func mysqlChannel() updates.Channel {
	releases := []updates.State{
		{ID: "v1.21.0", Tag: "v1.21.0", Migration: "extend_object_id"},
		{ID: "v1.19.1", Tag: "v1.19.1", Migration: "add_caveat"},
		{ID: "v1.18.0", Tag: "v1.18.0", Migration: "add_caveat"},
		{ID: "v1.17.0", Tag: "v1.17.0", Migration: "add_caveat"},
		{ID: "v1.16.2", Tag: "v1.16.2", Migration: "add_caveat"},
		{ID: "v1.16.1", Tag: "v1.16.1", Migration: "add_caveat", Deprecated: true},
		{ID: "v1.16.0", Tag: "v1.16.0", Migration: "add_caveat", Deprecated: true},
		{ID: "v1.15.0", Tag: "v1.15.0", Migration: "add_caveat"},
		{ID: "v1.14.1", Tag: "v1.14.1", Migration: "add_caveat"},
		{ID: "v1.14.0", Tag: "v1.14.0", Migration: "add_caveat", Deprecated: true},
		{ID: "v1.13.0", Tag: "v1.13.0", Migration: "add_ns_config_id"},
		{ID: "v1.12.0", Tag: "v1.12.0", Migration: "add_ns_config_id"},
		{ID: "v1.11.0", Tag: "v1.11.0", Migration: "add_ns_config_id"},
		{ID: "v1.10.0", Tag: "v1.10.0", Migration: "add_ns_config_id"},
		{ID: "v1.9.0", Tag: "v1.9.0", Migration: "add_unique_datastore_id"},
		{ID: "v1.8.0", Tag: "v1.8.0", Migration: "add_unique_datastore_id"},
		{ID: "v1.7.1", Tag: "v1.7.1", Migration: "add_unique_datastore_id"},
		{ID: "v1.7.0", Tag: "v1.7.0", Migration: "add_unique_datastore_id", Deprecated: true},
	}
	edgePatterns := map[string]string{
		"v1.19.1": ">=1.21.0",
		"v1.18.0": ">=1.19.1",
		"v1.17.0": ">=1.18.0",
		"v1.16.2": ">=1.17.0",
		"v1.16.1": ">=1.16.2",
		"v1.16.0": ">=1.16.2",
		"v1.15.0": ">=1.16.2",
		"v1.14.1": ">=1.15.0",
		"v1.14.0": ">=1.14.1",
		"v1.13.0": ">=1.14.1",
		"v1.12.0": ">=1.13.0",
		"v1.11.0": ">=1.12.0",
		"v1.10.0": ">=1.11.0",
		"v1.9.0":  ">=1.10.0",
		"v1.8.0":  ">=1.9.0",
		"v1.7.1":  ">=1.8.0",
		"v1.7.0":  ">=1.7.1",
	}
	return updates.Channel{
		Name: "stable",
		Metadata: map[string]string{
			"datastore": "mysql",
			"default":   "true",
		},
		Nodes: releases,
		Edges: edgesFromPatterns(edgePatterns, releases),
	}
}

func spannerChannel() updates.Channel {
	releases := []updates.State{
		{ID: "v1.21.0", Tag: "v1.21.0", Migration: "add-caveats"},
		{ID: "v1.19.1", Tag: "v1.19.1", Migration: "add-caveats"},
		{ID: "v1.18.0", Tag: "v1.18.0", Migration: "add-caveats"},
		{ID: "v1.17.0", Tag: "v1.17.0", Migration: "add-caveats"},
		{ID: "v1.16.2", Tag: "v1.16.2", Migration: "add-caveats"},
		{ID: "v1.16.1", Tag: "v1.16.1", Migration: "add-caveats", Deprecated: true},
		{ID: "v1.16.0", Tag: "v1.16.0", Migration: "add-caveats", Deprecated: true},
		{ID: "v1.15.0", Tag: "v1.15.0", Migration: "add-caveats"},
		{ID: "v1.14.1", Tag: "v1.14.1", Migration: "add-caveats"},
		{ID: "v1.14.0", Tag: "v1.14.0", Migration: "add-caveats", Deprecated: true},
		{ID: "v1.13.0", Tag: "v1.13.0", Migration: "add-metadata-and-counters"},
		{ID: "v1.12.0", Tag: "v1.12.0", Migration: "add-metadata-and-counters"},
		{ID: "v1.11.0", Tag: "v1.11.0", Migration: "add-metadata-and-counters"},
		{ID: "v1.10.0", Tag: "v1.10.0", Migration: "add-metadata-and-counters"},
		{ID: "v1.9.0", Tag: "v1.9.0", Migration: "add-metadata-and-counters"},
		{ID: "v1.8.0", Tag: "v1.8.0", Migration: "add-metadata-and-counters"},
	}
	edgePatterns := map[string]string{
		"v1.19.1": ">=1.21.0",
		"v1.18.0": ">=1.19.1",
		"v1.17.0": ">=1.18.0",
		"v1.16.2": ">=1.17.0",
		"v1.16.1": ">=1.16.2",
		"v1.16.0": ">=1.16.2",
		"v1.15.0": ">=1.16.2",
		"v1.14.1": ">=1.15.0",
		"v1.14.0": ">=1.14.1",
		"v1.13.0": ">=1.14.1",
		"v1.12.0": ">=1.13.0",
		"v1.11.0": ">=1.12.0",
		"v1.10.0": ">=1.11.0",
		"v1.9.0":  ">=1.10.0",
		"v1.8.0":  ">=1.9.0",
	}
	return updates.Channel{
		Name: "stable",
		Metadata: map[string]string{
			"datastore": "spanner",
			"default":   "true",
		},
		Nodes: releases,
		Edges: edgesFromPatterns(edgePatterns, releases),
	}
}

func memoryChannel() updates.Channel {
	releases := []updates.State{
		{ID: "v1.21.0", Tag: "v1.21.0"},
		{ID: "v1.19.1", Tag: "v1.19.1"},
		{ID: "v1.18.0", Tag: "v1.18.0"},
		{ID: "v1.17.0", Tag: "v1.17.0"},
		{ID: "v1.16.2", Tag: "v1.16.2"},
		{ID: "v1.16.1", Tag: "v1.16.1", Deprecated: true},
		{ID: "v1.16.0", Tag: "v1.16.0", Deprecated: true},
		{ID: "v1.15.0", Tag: "v1.15.0"},
		{ID: "v1.14.1", Tag: "v1.14.1"},
		{ID: "v1.14.0", Tag: "v1.14.0", Deprecated: true},
		{ID: "v1.13.0", Tag: "v1.13.0"},
		{ID: "v1.12.0", Tag: "v1.12.0"},
		{ID: "v1.11.0", Tag: "v1.11.0"},
		{ID: "v1.10.0", Tag: "v1.10.0"},
		{ID: "v1.9.0", Tag: "v1.9.0"},
		{ID: "v1.8.0", Tag: "v1.8.0"},
		{ID: "v1.7.1", Tag: "v1.7.1"},
		{ID: "v1.7.0", Tag: "v1.7.0", Deprecated: true},
		{ID: "v1.6.0", Tag: "v1.6.0"},
		{ID: "v1.5.0", Tag: "v1.5.0"},
		{ID: "v1.4.0", Tag: "v1.4.0"},
		{ID: "v1.3.0", Tag: "v1.3.0"},
		{ID: "v1.2.0", Tag: "v1.2.0"},
	}
	edgePatterns := make(map[string]string)
	for _, from := range releases {
		edgePatterns[from.ID] = ">" + strings.TrimPrefix(from.ID, "v")
	}

	return updates.Channel{
		Name: "stable",
		Metadata: map[string]string{
			"datastore": "memory",
			"default":   "true",
		},
		Nodes: releases,
		Edges: edgesFromPatterns(edgePatterns, releases),
	}
}

func edgesFromPatterns(patterns map[string]string, releases []updates.State) map[string][]string {
	edges := make(map[string][]string)
	for from, to := range patterns {
		toRange := semver.MustParseRange(to)
		for _, release := range releases {
			if release.Deprecated {
				continue
			}
			if !toRange(semver.MustParse(strings.TrimPrefix(release.ID, "v"))) {
				continue
			}
			edges[from] = append(edges[from], release.ID)
		}
		sort.Slice(edges[from], func(i, j int) bool {
			v1 := semver.MustParse(strings.TrimPrefix(edges[from][i], "v"))
			v2 := semver.MustParse(strings.TrimPrefix(edges[from][j], "v"))
			return v1.LT(v2)
		})
	}
	return edges
}

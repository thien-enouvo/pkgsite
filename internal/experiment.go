// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package internal contains data used through x/pkgsite.
package internal

const (
	ExperimentAutocomplete = "autocomplete"
	ExperimentSidenav      = "sidenav"
	ExperimentUnitPage     = "unit-page"
	ExperimentUseUnits     = "use-units"
)

// Experiments represents all of the active experiments in the codebase and
// a description of each experiment.
var Experiments = map[string]string{
	ExperimentAutocomplete: "Enable autocomplete with search.",
	ExperimentSidenav:      "Display documentation index on the left sidenav.",
	ExperimentUnitPage:     "Enable the redesigned details page.",
	ExperimentUseUnits:     "Read from paths, documentation, readmes, and package_imports tables.",
}

// Experiment holds data associated with an experimental feature for frontend
// or worker.
type Experiment struct {
	// This struct is used to decode dynamic config (see
	// internal/config/dynconfig). Make sure that changes to this struct are
	// coordinated with the deployment of config files.

	// Name is the name of the feature.
	Name string

	// Rollout is the percentage of requests enrolled in the experiment.
	Rollout uint

	// Description provides a description of the experiment.
	Description string
}

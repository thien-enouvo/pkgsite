// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package etl

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v7"
	"golang.org/x/discovery/internal/complete"
	"golang.org/x/discovery/internal/postgres"
	"golang.org/x/discovery/internal/testing/sample"
)

func TestUpdateRedisIndexes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()
	defer postgres.ResetTestDB(testDB, t)
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	// Set up a simple test case with two module versions and two packages. The
	// package in v2 imports the package in v1. By setting our 'popular cutoff'
	// to 1, we can force the package in v1 to be considered popular.
	rc := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	v1 := sample.Version()
	v1.ModulePath = "github.com/something"
	v1.Packages[0].Path = v1.ModulePath + "/apples/bananas"
	v2 := sample.Version()
	v2.ModulePath = "github.com/something/else"
	v2.Packages[0].Path = v2.ModulePath + "/oranges/bananas"
	v2.Packages[0].Imports = []string{v1.Packages[0].Path}
	if err := testDB.InsertVersion(ctx, v1); err != nil {
		t.Fatal(err)
	}
	if err := testDB.InsertVersion(ctx, v2); err != nil {
		t.Fatal(err)
	}
	if _, err := testDB.UpdateSearchDocumentsImportedByCount(ctx); err != nil {
		t.Fatal(err)
	}
	if err := updateRedisIndexes(ctx, testDB.GetSQLDB(), rc, 1); err != nil {
		t.Fatal(err)
	}
	popCount, err := rc.ZCount(complete.PopularKey, "0", "0").Result()
	if err != nil {
		t.Fatal(err)
	}
	// There are 4 path components in github.com/something/apples/bananas
	if popCount != 4 {
		t.Errorf("got %d popular autocompletions, want %d", popCount, 4)
	}
	remCount, err := rc.ZCount(complete.RemainingKey, "0", "0").Result()
	if err != nil {
		t.Fatal(err)
	}
	// There are 5 path components in github.com/something/else/oranges/bananas
	if remCount != 5 {
		t.Errorf("got %d remaining autocompletions, want %d", remCount, 5)
	}
}
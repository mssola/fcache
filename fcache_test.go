// Copyright (C) 2013 Miquel Sabaté Solà
// This file is licensed under the MIT license.
// See the LICENSE file.

package fcache

import (
    "testing"
    "os"
    "time"
    "github.com/bmizerany/assert"
)

func TestNewCache(t *testing.T) {
    // I'm going to remove it, so I can test that it gets created again
    // when calling the NewCache function.
    dir := "/tmp/mssola/fcache/test"
    os.RemoveAll("/tmp/mssola/fcache")
    _, err := os.Stat(dir)
    assert.Equal(t, os.IsNotExist(err), true)
    cache := NewCache(dir, 2 * time.Second, 0774)

    // We first check that the members are properly set.
    assert.Equal(t, cache.Dir, dir)
    assert.Equal(t, cache.Expiration, 2 * time.Second)

    // We check that the directory has been created.
    _, err = os.Stat(dir)
    assert.Equal(t, err, nil)
}

func TestSet(t *testing.T) {
//     Printf("TestSet: TODO\n")
}

func TestGetNonExistent(t *testing.T) {
//     Printf("TestGetNonExistent: TODO\n")
}

func TestGetInvalid(t *testing.T) {
//     Printf("TestGetInvalid: TODO\n")
}

func TestGetValid(t *testing.T) {
//     Printf("TestGetValid: TODO\n")
}

func TestFlush(t *testing.T) {
//     Printf("TestFlush: TODO\n")
}

func TestValid(t *testing.T) {
//     Printf("TestFlush: TODO\n")
}

// Copyright (C) 2013 Miquel Sabaté Solà
// This file is licensed under the MIT license.
// See the LICENSE file.

package fcache

import (
    "os"
    "time"
    "testing"
    "io/ioutil"
    "github.com/bmizerany/assert"
)

const (
    // This path contains the directory of the cache.
    basePath = "/tmp/mssola/fcache"

    // The path to the directory of the cache.
    cachePath = basePath + "/test"
)

func TestNewCache(t *testing.T) {
    // I'm going to remove it, so I can test that it gets created again
    // when calling the NewCache function.
    os.RemoveAll(basePath)
    _, err := os.Stat(cachePath)
    assert.Equal(t, os.IsNotExist(err), true)
    cache := NewCache(cachePath, 2 * time.Second, 0774)

    // We first check that the members are properly set.
    assert.Equal(t, cache.Dir, cachePath)
    assert.Equal(t, cache.Expiration, 2 * time.Second)

    // We check that the directory has been created.
    _, err = os.Stat(cachePath)
    assert.Equal(t, err, nil)

    // Remove it again, so the environment is clean for the rest of the tests.
    os.RemoveAll(basePath)
}

func TestSet(t *testing.T) {
    os.RemoveAll(basePath)
    cache := NewCache(cachePath, 2 * time.Second, 0774)
    cache.Set("file.txt", []byte("contents"))

    // It gets created the first time and set the given contents.
    _, err := os.Stat(cachePath + "/file.txt")
    assert.Equal(t, err, nil)
    got, _ := ioutil.ReadFile(cachePath + "/file.txt")
    assert.Equal(t, got, []byte("contents"))

    // Now we'll overwrite it.
    cache.Set("file.txt", []byte("lalala"))
    got, _ = ioutil.ReadFile(cachePath + "/file.txt")
    assert.Equal(t, got, []byte("lalala"))

    os.RemoveAll(basePath)
}

func TestGetNonExistent(t *testing.T) {
    os.RemoveAll(basePath)
    cache := NewCache(cachePath, 2 * time.Second, 0774)

    // We get nothing because it doesn't exist.
    got, err := cache.Get("file.txt")
    assert.Equal(t, got, []byte(""))
    assert.NotEqual(t, err, nil)
    assert.Equal(t, err.Error(), "miss.")

    // Check the isValid function.
    assert.Equal(t, false, cache.IsValid("file.txt"))

    os.RemoveAll(basePath)
}

func TestGetInvalid(t *testing.T) {
    os.RemoveAll(basePath)
    cache := NewCache(cachePath, 1 * time.Millisecond, 0774)
    cache.Set("file.txt", []byte("contents"))

    // It really exists, but it will get cold after 1 ms.
    _, err := os.Stat(cachePath + "/file.txt")
    assert.Equal(t, err, nil)
    time.Sleep(10 * time.Millisecond)

    // We get nothing, since it's cold.
    got, err := cache.Get("file.txt")
    assert.Equal(t, got, []byte(""))
    assert.NotEqual(t, err, nil)
    assert.Equal(t, err.Error(), "miss.")

    // Check the isValid function.
    assert.Equal(t, false, cache.IsValid("file.txt"))

    // It gets removed when it's known to be cold.
    _, err = os.Stat(cachePath + "/file.txt")
    assert.Equal(t, os.IsNotExist(err), true)

    os.RemoveAll(basePath)
}

func TestGetValid(t *testing.T) {
    os.RemoveAll(basePath)
    cache := NewCache(cachePath, 1 * time.Hour, 0774)
    cache.Set("file.txt", []byte("contents"))

    // It's hot: hit !
    got, err := cache.Get("file.txt")
    assert.Equal(t, got, []byte("contents"))
    assert.Equal(t, err, nil)

    // Check the isValid function.
    assert.Equal(t, true, cache.IsValid("file.txt"))

    // It remains untouched, because everything is fine.
    _, err = os.Stat(cachePath + "/file.txt")
    assert.Equal(t, err, nil)

    os.RemoveAll(basePath)
}

func TestPath(t *testing.T) {
    cache := NewCache(cachePath, 2 * time.Second, 0774)
    path := cache.Path("file.txt")
    assert.Equal(t, path, cachePath + "/file.txt")
}

func TestFlush(t *testing.T) {
    os.RemoveAll(basePath)
    cache := NewCache(cachePath, 2 * time.Second, 0774)
    cache.Set("file.txt", []byte("contents"))
    cache.Set("another.txt", []byte("contents"))

    _, err := os.Stat(cachePath + "/file.txt")
    assert.Equal(t, err, nil)
    _, err = os.Stat(cachePath + "/another.txt")
    assert.Equal(t, err, nil)
    cache.Flush("file.txt")
    _, err = os.Stat(cachePath + "/file.txt")
    assert.Equal(t, os.IsNotExist(err), true)
    _, err = os.Stat(cachePath + "/another.txt")
    assert.Equal(t, err, nil)

    os.RemoveAll(basePath)
}

func TestFlushAll(t *testing.T) {
    os.RemoveAll(basePath)
    cache := NewCache(cachePath, 2 * time.Second, 0774)
    cache.Set("file.txt", []byte("contents"))
    cache.Set("another.txt", []byte("contents"))

    _, err := os.Stat(cachePath + "/file.txt")
    assert.Equal(t, err, nil)
    _, err = os.Stat(cachePath + "/another.txt")
    assert.Equal(t, err, nil)
    cache.FlushAll()
    _, err = os.Stat(cachePath + "/file.txt")
    assert.Equal(t, os.IsNotExist(err), true)
    _, err = os.Stat(cachePath + "/another.txt")
    assert.Equal(t, os.IsNotExist(err), true)

    os.RemoveAll(basePath)
}

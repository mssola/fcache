// Copyright (C) 2013 Miquel Sabaté Solà
// This file is licensed under the MIT license.
// See the LICENSE file.

// This package encapsulates all the methods regarding the File Cache.
package fcache

import (
    "os"
    "fmt"
    "path"
    "time"
    "errors"
    "io/ioutil"
)

// This type contains some needed info that will be used when caching.
type Cache struct {
    // The directory
    Dir string

    // The expiration time to be set for each file.
    Expiration time.Duration

    // The permissions to be set when this cache creates new files.
    Permissions os.FileMode
}

// Get a pointer to an initialized Cache structure.
//
// dir        - The path to the cache directory. If the directory does not
//              exist, it will create a new directory with permissions 0644.
// expiration - The expiration time. That is, how many nanoseconds has to pass
//              by when a cache file is no longer considered valid.
// perm       - The permissions that the cache should operate in when creating
//              new files.
//
// Returns a Cache pointer that points to an initialized Cache structure. It
// will return nil if something goes wrong.
func NewCache(dir string, expiration time.Duration, perm os.FileMode) *Cache {
    // First of all, get the directory path straight.
    if _, err := os.Stat(dir); err != nil {
        if os.IsNotExist(err) {
            if err = os.MkdirAll(dir, perm); err != nil {
                fmt.Printf("Error: %v\n", err)
                return nil
            }
        } else {
            fmt.Printf("Error: %v\n", err)
            return nil
        }
    }

    // Now it's safe to create the cache.
    cache := new(Cache)
    cache.Dir = dir
    cache.Expiration = expiration
    cache.Permissions = perm
    return cache
}

// Set the contents for a cache file. If this file doesn't exist already, it
// will be created with permissions 0644.
//
// name     - The name of the file.
// contents - The contents that the cache file has to contain after calling
//            this function.
//
// Returns nil if everything was ok. Otherwise it will return an error.
func (c *Cache) Set(name string, contents []byte) error {
    url := path.Join(c.Dir, name)
    return ioutil.WriteFile(url, contents, c.Permissions)
}

// Get the contents of a valid cache file.
//
// name - The name of the file.
//
// Returns a slice of bytes and an error. The slice of bytes contain the
// contents of the cache file. The error is set to nil if everything was fine.
func (c *Cache) Get(name string) ([]byte, error) {
    url := path.Join(c.Dir, name)
    if fi, err := os.Stat(url); err == nil {
        elapsed := time.Now().Sub(fi.ModTime())
        if c.Expiration > elapsed {
            return ioutil.ReadFile(url)
        }
        // Remove this file, its time has expired.
        os.Remove(url)
    }
    return []byte{}, errors.New("miss.")
}

// Remove a cache file.
//
// name - The name of the file.
//
// Returns nil if everything was ok. Otherwise it will return an error.
func (c *Cache) Flush(name string) error {
    url := path.Join(c.Dir, name)
    return os.Remove(url)
}

// Remove all the files from the cache.
//
// Returns nil if everything was ok. Otherwise it will return an error.
func (c *Cache) FlushAll() error {
    url := path.Join(c.Dir)
    err := os.RemoveAll(url)
    if err == nil {
        err = os.MkdirAll(url, c.Permissions)
    }
    return err
}

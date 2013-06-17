# FCache [![Build Status](https://travis-ci.org/mssola/fcache.png?branch=master)](https://travis-ci.org/mssola/fcache)

This package implements a simple file cache. That is, a cache implemented through
files. This kind of caching is useful if the data that you want to cache is
stored in an external resource (e.g. an external API). This cache also
considers an expiration time for each of the files that it stores. Let's show 
some code as an example:

    // ...
    cache := fcache.NewCache("/tmp/fcache", 1*time.Hour, 0774)
    cache.Set("file.txt", []byte("Some contents."))
    fmt.Printf("%v\n", cache.Get("file.txt"))
    // ...

In the first line we create a new cache located at "/tmp/fcache". Note that
if the directory doesn't exist, it will be created for you. The second
parameter of the NewCache function says that every file handled by this cache
will have an expiration time of 1 hour. That is, in this case if we try to get
a file that its last modification has happenned more than one hour ago, it will
be considered invalid and it will be removed. The third parameter is the
permissions that will be used by the cache when creating new files. The second
line of code sets the contents of a file called "file.txt" and the third line
gets the contents of this file again. Finally, the Cache type also implements
the Flush and FlushAll functions. The former flushes one file in the cache, 
and the latter flushes all the files from the cache.

Copyright &copy; 2013 Miquel Sabaté Solà, released under the MIT License.

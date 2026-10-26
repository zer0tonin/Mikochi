package main

import "strings"

func normalizePath(path string) string {
    if path == "" {
        path = "/"
    }

    if !strings.HasPrefix(path, "/") {
        path = "/" + path
    }

    if !strings.HasSuffix(path, "/") {
        path = path + "/"
    }

    return path
}

package utils

import (
    "io/ioutil"
	"log"
)

type DirContents struct{
	Total int `json:"total"`
	Contents []string `json:"contents"`
}

/**
 *  Check if a directory is a legal thing to view
 */
func GetDirectoriesLookup(legal string) map[string]bool {
    var listings = make(map[string]bool)
    files, _ := ioutil.ReadDir(legal)
    for _, f := range files {
        if f.IsDir() {
            listings[f.Name()] = true
        }
    }
    return listings
}

/**
 * Grab a small preview list of all items in the directory.
 */
func ListDirs(dir string, previewCount int) map[string]DirContents {
	// Get the current listings, check they passed in a legal key
    var listings = make(map[string]DirContents)
    files, _ := ioutil.ReadDir(dir)
    for _, f := range files {
        if f.IsDir() {
            listings[f.Name()] = GetDirContents(dir + f.Name(), previewCount)
        }
    }
	log.Println("Reading from: ", dir, " With preview count", previewCount)
    return listings
}

/**
 *  Get all the content in a particular directory.
 */
func GetDirContents(dir string, limit int) DirContents {
    var arr = []string{}
    imgs, _ := ioutil.ReadDir(dir)

	total := 0
    for _, img := range imgs {
        if !img.IsDir() && len(arr) < limit {
            arr = append(arr, img.Name())
        }
		total++
    }
	log.Println("Limit for content dir was.", dir, " with limit", limit)
	return DirContents{
		Total: total,
		Contents: arr,
	}
}

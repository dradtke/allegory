// Package cache provides support for loading resources into memory.
package cache

import (
    "fmt"
    al "github.com/dradtke/go-allegro/allegro"
    "github.com/dradtke/gopher"
)

var (
    _images = make(map[string]*al.Bitmap)
)

func Clear() {
    for key, val := range _images {
        val.Destroy()
        delete(_images, key)
    }
}

func LoadImage(id, path string) {
    bmp, err := al.LoadBitmap(path)
    if err != nil {
        gopher.Fatal(err)
    }
    _images[id] = bmp
}

func Image(id string) *al.Bitmap {
    if bmp, ok := _images[id]; !ok {
        gopher.Fatal(fmt.Errorf("Image \"%s\" not found!", id))
        // won't get here, but the compiler needs it
        return nil
    } else {
        return bmp
    }
}

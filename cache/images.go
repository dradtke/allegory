package cache

import (
	"fmt"
	al "github.com/dradtke/go-allegro/allegro"
	"sync"
)

var _images = make(map[string]*al.Bitmap)

// ClearImages() removes all images from the cache.
func ClearImages() {
	for key, val := range _images {
		val.Destroy()
		delete(_images, key)
	}
}

// LoadImage() loads an image into the cache.
func LoadImage(path string) error {
	bmp, err := al.LoadBitmap(path)
	if err != nil {
		return err
	}
	_images[path] = bmp
	return nil
}

// LoadImages() loads multiple images into the
// cache.
func LoadImages(paths []string) []error {
	var (
		n    = len(paths)
		errs = make([]error, 0, n)
		wg   sync.WaitGroup
	)
	wg.Add(n)
	for _, path := range paths {
		go func(path string) {
			err := LoadImage(path)
			if err != nil {
				errs = append(errs, err)
			}
			wg.Done()
		}(path)
	}
	wg.Wait()
	return errs
}

// FindImage() finds an image in the cache. The value of path should
// be the one that was passed into LoadImage() or LoadImages().
func FindImage(path string) (*al.Bitmap, error) {
	if bmp, ok := _images[path]; ok {
		return bmp, nil
	}
	return nil, fmt.Errorf("Image \"%s\" not found!", path)
}

// Image() gets an image from the cache using FindImage(),
// panicking if it isn't found.
func Image(path string) *al.Bitmap {
	bmp, err := FindImage(path)
	if err != nil {
		panic(err)
	}
	return bmp
}

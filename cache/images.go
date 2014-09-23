package cache

import (
	"fmt"
	"github.com/dradtke/go-allegro/allegro"
	"sync"
)

var _images = make(map[string]*allegro.Bitmap)

type ImageNotFound struct {
	Key string
}

func (e *ImageNotFound) Error() string {
	return fmt.Sprintf("image %s not found", e.Key)
}

// ClearImages() removes all images from the cache.
func ClearImages() {
	for key, val := range _images {
		val.Destroy()
		delete(_images, key)
	}
}

// LoadImage() loads an image into the cache.
func LoadImage(path, key string) error {
	bmp, err := allegro.LoadBitmap(path)
	if err != nil {
		return err
	}
	if key == "" {
		key = path
	}
	_images[key] = bmp
	return nil
}

// LoadImages() loads multiple images into the
// cache.
func LoadImages(paths []string, pathToKey func(string) string) []error {
	var (
		n    = len(paths)
		errs = make([]error, 0, n)
		wg   sync.WaitGroup
	)
	wg.Add(n)
	for _, path := range paths {
		go func(path string) {
			err := LoadImage(path, pathToKey(path))
			if err != nil {
				errs = append(errs, err)
			}
			wg.Done()
		}(path)
	}
	wg.Wait()
	return errs
}

// FindImage() finds an image in the cache. If it
// doesn't exist, an error of type ImageNotFound is returned.
func FindImage(key string) (*allegro.Bitmap, error) {
	if bmp, ok := _images[key]; ok {
		return bmp, nil
	}
	return nil, &ImageNotFound{key}
}

// Image() gets an image from the cache using FindImage(),
// panicking if it isn't found.
func Image(key string) *allegro.Bitmap {
	bmp, err := FindImage(key)
	if err != nil {
		panic(err)
	}
	return bmp
}

package cache

import (
	"fmt"
	"github.com/dradtke/go-allegro/allegro"
	"os"
	"path/filepath"
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

// LoadImages() walks root recursively loading all the images that it can.
// It returns the first error encountered, which may or may not be meaningful
// depending on whether or not root contains non-image files.
func LoadImages(root string) error {
	root_len := len(root)
	return filepath.Walk(root, func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil
		}
		return LoadImage(path, path[root_len+1:])
	})
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

package image_processing

import (
	"bytes"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

//type result struct {
//	Data []byte
//	Err  error
//}
//
//func DownloadImage(done <-chan struct{}, url string, wp workerpool.WorkerPool) (<-chan []byte, <-chan error) {
//	content := make(chan []byte)
//	errchan := make(chan error, 1)
//
//	go func() {
//		defer close(content)
//		ontentDownload := make(chan []byte)
//		defer close(ontentDownload)
//		wp.AddTask(func() {
//			data, err := utils.DownloadFile(url)
//			if err != nil {
//				errchan <- err
//				return
//			}
//			_, errc := guessImageFormat(data)
//			if errc != nil {
//				errchan <- errc
//				return
//			}
//			ontentDownload <- data
//			return
//
//		})
//
//		select {
//		case data, ok := <-ontentDownload:
//			if ok {
//				content <- data
//			}
//		case <-done:
//			return
//		}
//	}()
//
//	return content, errchan
//}
//
//func ProcessImage(done chan struct{}, content <-chan []byte) <-chan *result {
//	results := make(chan *result)
//
//	thumbnailer := func() {
//		for buffer := range content {
//			bImage := bimg.NewImage(buffer)
//			srcImage, err := Translate(bImage, 200, 200)
//
//			if err != nil {
//				select {
//				case results <- &result{nil, err}:
//				case <-done:
//					return
//				}
//			}
//			select {
//			case results <- &result{srcImage, nil}:
//			case <-done:
//				return
//			}
//
//		}
//	}
//	const numThunnail = 5
//	var wg sync.WaitGroup
//	wg.Add(numThunnail)
//	for i := 0; i < numThunnail; i++ {
//		go func() {
//			thumbnailer()
//			wg.Done()
//		}()
//	}
//	go func() {
//		wg.Wait()
//		close(results)
//	}()
//	return results
//}
//
//func Translate(urls *bimg.Image, maxWidth, maxHeight int) (buffer []byte, err error) {
//	imageSize, _ := urls.Size()
//	var ratio float64
//	var newSize bimg.ImageSize
//	if imageSize.Width >= imageSize.Height {
//		ratio = float64(maxHeight) / float64(imageSize.Height)
//		if ratio >= 1.0 {
//			goto JustConvert
//		}
//		newSize.Width = int(math.Floor(float64(imageSize.Width) * ratio))
//		newSize.Height = int(math.Floor(float64(imageSize.Height) * ratio))
//	} else if imageSize.Width < imageSize.Height {
//		ratio = float64(maxWidth) / float64(imageSize.Width)
//		if ratio >= 1.0 {
//			goto JustConvert
//		}
//		newSize.Width = int(math.Floor(float64(imageSize.Width) * ratio))
//		newSize.Height = int(math.Floor(float64(imageSize.Height) * ratio))
//	}
//	_, err = urls.ForceResize(newSize.Width, newSize.Height)
//	if err != nil {
//		return
//	}
//JustConvert:
//	buffer, err = urls.Convert(bimg.JPEG)
//	if err != nil {
//		return
//	}
//	return
//}

// Guess urls format from gif/jpeg/png/webp
func GuessImageFormat(data []byte) (format string, err error) {
	// convert byte slice to io.Reader
	reader := bytes.NewReader(data)
	_, format, err = image.DecodeConfig(reader)

	return
}

//// Guess urls format from gif/jpeg/png/webp
//func guessImageFormat(r io.Reader) (format string, err error) {
//	_, format, err = urls.DecodeConfig(r)
//	return
//}
//
//// Guess urls mime types from gif/jpeg/png/webp
//func guessImageMimeTypes(r io.Reader) string {
//	format, _ := guessImageFormat(r)
//	if format == "" {
//		return ""
//	}
//	return mime.TypeByExtension("." + format)
//}

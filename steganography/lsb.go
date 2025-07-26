package steganography

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
)

const (
	// EndMarker is used to mark the end of hidden message
	EndMarker = "###END###"
)

// HideMessage hides a message in an image using LSB steganography
func HideMessage(img image.Image, message string) (image.Image, error) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Create a new RGBA image
	newImg := image.NewRGBA(bounds)

	// Copy the original image
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			newImg.Set(x, y, img.At(x, y))
		}
	}

	// Add end marker to message
	fullMessage := message + EndMarker
	messageBytes := []byte(fullMessage)

	// Convert message length to bytes
	messageLen := len(messageBytes)
	lenBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenBytes, uint32(messageLen))

	// Combine length and message
	dataToHide := append(lenBytes, messageBytes...)

	// Check if image can hold the message
	totalPixels := width * height * 3 // RGB channels
	if len(dataToHide)*8 > totalPixels {
		return nil, fmt.Errorf("image too small to hide message of length %d", len(message))
	}

	// Hide data in LSBs
	dataIndex := 0
	bitIndex := 0

	for y := 0; y < height && dataIndex < len(dataToHide); y++ {
		for x := 0; x < width && dataIndex < len(dataToHide); x++ {
			pixel := newImg.RGBAAt(x, y)
			r, g, b, a := pixel.R, pixel.G, pixel.B, pixel.A

			// Hide data in red channel
			if dataIndex < len(dataToHide) {
				bit := (dataToHide[dataIndex] >> bitIndex) & 1
				r = (r & 0xFE) | bit
				bitIndex++
				if bitIndex == 8 {
					bitIndex = 0
					dataIndex++
				}
			}

			// Hide data in green channel
			if dataIndex < len(dataToHide) {
				bit := (dataToHide[dataIndex] >> bitIndex) & 1
				g = (g & 0xFE) | bit
				bitIndex++
				if bitIndex == 8 {
					bitIndex = 0
					dataIndex++
				}
			}

			// Hide data in blue channel
			if dataIndex < len(dataToHide) {
				bit := (dataToHide[dataIndex] >> bitIndex) & 1
				b = (b & 0xFE) | bit
				bitIndex++
				if bitIndex == 8 {
					bitIndex = 0
					dataIndex++
				}
			}

			newImg.SetRGBA(x, y, color.RGBA{
				R: r,
				G: g,
				B: b,
				A: a,
			})
		}
	}

	return newImg, nil
}

// ExtractMessage extracts a hidden message from an image
func ExtractMessage(img image.Image) (string, error) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// First, extract the message length (first 4 bytes)
	lenBytes := make([]byte, 4)
	dataIndex := 0
	bitIndex := 0

	for y := 0; y < height && dataIndex < 4; y++ {
		for x := 0; x < width && dataIndex < 4; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

			// Extract from red channel
			if dataIndex < 4 {
				bit := r8 & 1
				lenBytes[dataIndex] |= (bit << bitIndex)
				bitIndex++
				if bitIndex == 8 {
					bitIndex = 0
					dataIndex++
				}
			}

			// Extract from green channel
			if dataIndex < 4 {
				bit := g8 & 1
				lenBytes[dataIndex] |= (bit << bitIndex)
				bitIndex++
				if bitIndex == 8 {
					bitIndex = 0
					dataIndex++
				}
			}

			// Extract from blue channel
			if dataIndex < 4 {
				bit := b8 & 1
				lenBytes[dataIndex] |= (bit << bitIndex)
				bitIndex++
				if bitIndex == 8 {
					bitIndex = 0
					dataIndex++
				}
			}
		}
	}

	messageLen := binary.LittleEndian.Uint32(lenBytes)
	if messageLen == 0 || messageLen > 1000000 { // Sanity check
		return "", fmt.Errorf("invalid message length: %d", messageLen)
	}

	// Extract the message
	messageBytes := make([]byte, messageLen)
	dataIndex = 0
	bitIndex = 0
	totalBytesExtracted := 0

	// Continue from where we left off
	startY := 0
	startX := 0

	// Calculate starting position after length bytes
	pixelsUsed := (4 * 8) / 3 // 4 bytes * 8 bits / 3 channels per pixel
	if (4*8)%3 != 0 {
		pixelsUsed++
	}
	startY = pixelsUsed / width
	startX = pixelsUsed % width

	for y := startY; y < height && totalBytesExtracted < int(messageLen); y++ {
		startXForRow := 0
		if y == startY {
			startXForRow = startX
		}

		for x := startXForRow; x < width && totalBytesExtracted < int(messageLen); x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

			// Extract from red channel
			if totalBytesExtracted < int(messageLen) {
				bit := r8 & 1
				messageBytes[totalBytesExtracted] |= (bit << bitIndex)
				bitIndex++
				if bitIndex == 8 {
					bitIndex = 0
					totalBytesExtracted++
				}
			}

			// Extract from green channel
			if totalBytesExtracted < int(messageLen) {
				bit := g8 & 1
				messageBytes[totalBytesExtracted] |= (bit << bitIndex)
				bitIndex++
				if bitIndex == 8 {
					bitIndex = 0
					totalBytesExtracted++
				}
			}

			// Extract from blue channel
			if totalBytesExtracted < int(messageLen) {
				bit := b8 & 1
				messageBytes[totalBytesExtracted] |= (bit << bitIndex)
				bitIndex++
				if bitIndex == 8 {
					bitIndex = 0
					totalBytesExtracted++
				}
			}
		}
	}

	message := string(messageBytes)

	// Remove end marker if present
	if len(message) >= len(EndMarker) {
		if message[len(message)-len(EndMarker):] == EndMarker {
			message = message[:len(message)-len(EndMarker)]
		}
	}

	return message, nil
}

// DecodeImage decodes an image from io.Reader
func DecodeImage(r io.Reader) (image.Image, string, error) {
	// Read all data first
	buf := &bytes.Buffer{}
	_, err := io.Copy(buf, r)
	if err != nil {
		return nil, "", err
	}

	// Try PNG first
	img, err := png.Decode(bytes.NewReader(buf.Bytes()))
	if err == nil {
		return img, "png", nil
	}

	// Try JPEG
	img, err = jpeg.Decode(bytes.NewReader(buf.Bytes()))
	if err == nil {
		return img, "jpeg", nil
	}

	return nil, "", fmt.Errorf("unsupported image format")
}

// EncodeImage encodes an image to the specified format
func EncodeImage(w io.Writer, img image.Image, format string) error {
	switch format {
	case "png":
		return png.Encode(w, img)
	case "jpeg":
		return jpeg.Encode(w, img, &jpeg.Options{Quality: 95})
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

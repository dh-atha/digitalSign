package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"digitalSign/crypto"
	"digitalSign/steganography"
)

// HomeHandler serves the main page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Digital Signature & Steganography</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background-color: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { text-align: center; color: #333; margin-bottom: 30px; }
        .nav-buttons { display: flex; justify-content: space-around; flex-wrap: wrap; gap: 15px; }
        .nav-button { 
            display: inline-block; 
            padding: 15px 30px; 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); 
            color: white; 
            text-decoration: none; 
            border-radius: 8px; 
            text-align: center; 
            transition: transform 0.2s, box-shadow 0.2s;
            min-width: 200px;
        }
        .nav-button:hover { transform: translateY(-2px); box-shadow: 0 4px 15px rgba(0,0,0,0.2); }
        .description { margin: 20px 0; padding: 20px; background: #f8f9fa; border-radius: 8px; border-left: 4px solid #667eea; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Digital Signature & Steganography Application</h1>
        
        <div class="description">
            <h3>Aplikasi Kriptografi dengan RSA</h3>
            <p>Aplikasi ini menyediakan fitur-fitur:</p>
            <ul>
                <li><strong>Digital Signature:</strong> Menandatangani dan memverifikasi dokumen menggunakan RSA</li>
                <li><strong>Steganography:</strong> Menyembunyikan dan mengekstrak pesan dalam gambar</li>
                <li><strong>Combined Operations:</strong> Menggabungkan kedua operasi</li>
            </ul>
        </div>
        
        <div class="nav-buttons">
            <a href="/digital-signature" class="nav-button">
                🔐 Digital Signature
            </a>
            <a href="/steganography" class="nav-button">
                🖼️ Steganography
            </a>
            <a href="/combined" class="nav-button">
                🔧 Combined Operations
            </a>
        </div>
    </div>
</body>
</html>`

	t, _ := template.New("home").Parse(tmpl)
	t.Execute(w, nil)
}

// Response structures
type KeyPairResponse struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

type SignResponse struct {
	Signature string `json:"signature"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
}

type VerifyResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// GenerateKeysHandler generates RSA key pair
func GenerateKeysHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get key size from request, default to 2048
	keySize := 2048
	if keySizeStr := r.FormValue("keySize"); keySizeStr != "" {
		if size, err := strconv.Atoi(keySizeStr); err == nil && size >= 1024 {
			keySize = size
		}
	}

	keyPair, err := crypto.GenerateRSAKeyPair(keySize)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	response := KeyPairResponse{
		PrivateKey: keyPair.PrivateKeyToPEM(),
		PublicKey:  keyPair.PublicKeyToPEM(),
	}

	json.NewEncoder(w).Encode(response)
}

// SignDocumentHandler signs a document
func SignDocumentHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get document
	file, fileHeader, err := r.FormFile("document")
	if err != nil {
		http.Error(w, "No document provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	document, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read document", http.StatusBadRequest)
		return
	}

	// Get private key
	privateKeyPEM := r.FormValue("privateKey")
	if privateKeyPEM == "" {
		http.Error(w, "No private key provided", http.StatusBadRequest)
		return
	}

	privateKey, err := crypto.ParsePrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		http.Error(w, "Invalid private key: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Sign document
	signature, err := crypto.SignDocument(document, privateKey)
	if err != nil {
		http.Error(w, "Failed to sign document: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate filename for signature file
	originalFilename := fileHeader.Filename
	if originalFilename == "" {
		originalFilename = "document"
	}
	// Remove extension from filename
	if lastDot := strings.LastIndex(originalFilename, "."); lastDot != -1 {
		originalFilename = originalFilename[:lastDot]
	}
	signatureFilename := originalFilename + ".sig"
	encodedFilename := url.PathEscape(signatureFilename)

	// Set headers for file download
	w.Header().Set("Content-Type", "application/octet-stream")

	w.Header().Set("Content-Disposition", fmt.Sprintf(
		`attachment; filename="%s"; filename*=UTF-8''%s`,
		signatureFilename,
		encodedFilename,
	))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(signature)))

	// Write signature to response
	w.Write([]byte(signature))
}

// VerifySignatureHandler verifies a document signature
func VerifySignatureHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to parse form"})
		return
	}

	// Get document
	file, _, err := r.FormFile("document")
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "No document provided"})
		return
	}
	defer file.Close()

	document, err := io.ReadAll(file)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to read document"})
		return
	}

	// Get signature file
	sigFile, _, err := r.FormFile("signatureFile")
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "No signature file provided"})
		return
	}
	defer sigFile.Close()

	signatureBytes, err := io.ReadAll(sigFile)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to read signature file"})
		return
	}

	signature := string(signatureBytes)

	// Get public key
	publicKeyPEM := r.FormValue("publicKey")
	if publicKeyPEM == "" {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "No public key provided"})
		return
	}

	publicKey, err := crypto.ParsePublicKeyFromPEM(publicKeyPEM)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid public key: " + err.Error()})
		return
	}

	// Verify signature
	err = crypto.VerifySignature(document, signature, publicKey)
	if err != nil {
		response := VerifyResponse{
			Valid:   false,
			Message: "Signature verification failed: " + err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := VerifyResponse{
		Valid:   true,
		Message: "Signature is valid",
	}

	json.NewEncoder(w).Encode(response)
}

// HideMessageHandler hides a message in an image
func HideMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get image
	imageFile, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "No image provided", http.StatusBadRequest)
		return
	}
	defer imageFile.Close()

	// Get message
	message := r.FormValue("message")
	if message == "" {
		http.Error(w, "No message provided", http.StatusBadRequest)
		return
	}

	// Decode image
	img, format, err := steganography.DecodeImage(imageFile)
	if err != nil {
		http.Error(w, "Failed to decode image: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Hide message
	newImg, err := steganography.HideMessage(img, message)
	if err != nil {
		http.Error(w, "Failed to hide message: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Set appropriate content type
	switch format {
	case "png":
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Disposition", "attachment; filename=\"stego_image.png\"")
	case "jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Disposition", "attachment; filename=\"stego_image.jpg\"")
	}

	// Encode and send image
	err = steganography.EncodeImage(w, newImg, format)
	if err != nil {
		http.Error(w, "Failed to encode image: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// ExtractMessageHandler extracts a message from an image
func ExtractMessageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to parse form"})
		return
	}

	// Get image
	imageFile, _, err := r.FormFile("image")
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "No image provided"})
		return
	}
	defer imageFile.Close()

	// Decode image
	img, _, err := steganography.DecodeImage(imageFile)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to decode image: " + err.Error()})
		return
	}

	// Extract message
	message, err := steganography.ExtractMessage(img)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to extract message: " + err.Error()})
		return
	}

	response := map[string]string{
		"message": message,
		"success": "Message extracted successfully",
	}

	json.NewEncoder(w).Encode(response)
}

// SignAndHideHandler combines digital signature and steganography
func SignAndHideHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get document
	docFile, _, err := r.FormFile("document")
	if err != nil {
		http.Error(w, "No document provided", http.StatusBadRequest)
		return
	}
	defer docFile.Close()

	document, err := io.ReadAll(docFile)
	if err != nil {
		http.Error(w, "Failed to read document", http.StatusBadRequest)
		return
	}

	// Get image
	imageFile, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "No image provided", http.StatusBadRequest)
		return
	}
	defer imageFile.Close()

	// Get private key
	privateKeyPEM := r.FormValue("privateKey")
	if privateKeyPEM == "" {
		http.Error(w, "No private key provided", http.StatusBadRequest)
		return
	}

	privateKey, err := crypto.ParsePrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		http.Error(w, "Invalid private key: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Sign document
	signature, err := crypto.SignDocument(document, privateKey)
	if err != nil {
		http.Error(w, "Failed to sign document: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Decode image
	img, format, err := steganography.DecodeImage(imageFile)
	if err != nil {
		http.Error(w, "Failed to decode image: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Prepare message with signature
	message := fmt.Sprintf("SIGNED_DOCUMENT:%s:SIGNATURE:%s", string(document), signature)

	// Hide message in image
	newImg, err := steganography.HideMessage(img, message)
	if err != nil {
		http.Error(w, "Failed to hide message: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Set appropriate content type
	switch format {
	case "png":
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Disposition", "attachment; filename=\"signed_stego_image.png\"")
	case "jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Disposition", "attachment; filename=\"signed_stego_image.jpg\"")
	}

	// Encode and send image
	err = steganography.EncodeImage(w, newImg, format)
	if err != nil {
		http.Error(w, "Failed to encode image: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// ExtractAndVerifyHandler extracts and verifies a signed document from an image
func ExtractAndVerifyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to parse form"})
		return
	}

	// Get image
	imageFile, _, err := r.FormFile("image")
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "No image provided"})
		return
	}
	defer imageFile.Close()

	// Get public key
	publicKeyPEM := r.FormValue("publicKey")
	if publicKeyPEM == "" {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "No public key provided"})
		return
	}

	publicKey, err := crypto.ParsePublicKeyFromPEM(publicKeyPEM)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid public key: " + err.Error()})
		return
	}

	// Decode image
	img, _, err := steganography.DecodeImage(imageFile)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to decode image: " + err.Error()})
		return
	}

	// Extract message
	message, err := steganography.ExtractMessage(img)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to extract message: " + err.Error()})
		return
	}

	// Parse the message to extract document and signature
	// Expected format: "SIGNED_DOCUMENT:<document>:SIGNATURE:<signature>"
	if len(message) < 16 || message[:16] != "SIGNED_DOCUMENT:" {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid message format"})
		return
	}

	// Find signature separator
	sigIndex := bytes.Index([]byte(message), []byte(":SIGNATURE:"))
	if sigIndex == -1 {
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Signature not found in message"})
		return
	}

	document := message[16:sigIndex]
	signature := message[sigIndex+11:] // +11 for ":SIGNATURE:"

	// Verify signature
	err = crypto.VerifySignature([]byte(document), signature, publicKey)

	response := map[string]interface{}{
		"document":  document,
		"signature": signature,
		"valid":     err == nil,
	}

	if err != nil {
		response["error"] = "Signature verification failed: " + err.Error()
	} else {
		response["message"] = "Document extracted and signature verified successfully"
	}

	json.NewEncoder(w).Encode(response)
}

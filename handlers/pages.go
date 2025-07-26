package handlers

import (
	"html/template"
	"net/http"
)

// DigitalSignaturePageHandler serves the digital signature page
func DigitalSignaturePageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Digital Signature - RSA</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background-color: #f5f5f5; }
        .container { max-width: 900px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { text-align: center; color: #333; margin-bottom: 30px; }
        .section { margin: 20px 0; padding: 20px; border: 1px solid #ddd; border-radius: 8px; background: #fafafa; }
        .form-group { margin: 15px 0; }
        label { display: block; margin-bottom: 5px; font-weight: bold; color: #555; }
        input[type="file"], input[type="number"], textarea, select { 
            width: 100%; 
            padding: 10px; 
            border: 1px solid #ddd; 
            border-radius: 4px; 
            box-sizing: border-box; 
        }
        textarea { height: 100px; resize: vertical; }
        button { 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); 
            color: white; 
            padding: 12px 24px; 
            border: none; 
            border-radius: 6px; 
            cursor: pointer; 
            font-size: 16px; 
            margin: 5px;
            transition: transform 0.2s;
        }
        button:hover { transform: translateY(-1px); }
        .result { 
            margin-top: 20px; 
            padding: 15px; 
            border-radius: 4px; 
            white-space: pre-wrap; 
            word-break: break-all; 
        }
        .success { background-color: #d4edda; color: #155724; border: 1px solid #c3e6cb; }
        .error { background-color: #f8d7da; color: #721c24; border: 1px solid #f5c6cb; }
        .nav-link { display: inline-block; margin: 10px; padding: 10px 20px; background: #6c757d; color: white; text-decoration: none; border-radius: 4px; }
        .nav-link:hover { background: #5a6268; }
    </style>
</head>
<body>
    <div class="container">
        <div style="text-align: center; margin-bottom: 20px;">
            <a href="/" class="nav-link">🏠 Home</a>
            <a href="/steganography" class="nav-link">🖼️ Steganography</a>
            <a href="/combined" class="nav-link">🔧 Combined</a>
        </div>
        
        <h1>🔐 Digital Signature with RSA</h1>
        
        <!-- Key Generation Section -->
        <div class="section">
            <h2>1. Generate RSA Key Pair</h2>
            <div class="form-group">
                <label for="keySize">Key Size:</label>
                <select id="keySize">
                    <option value="1024">1024 bits</option>
                    <option value="2048" selected>2048 bits</option>
                    <option value="4096">4096 bits</option>
                </select>
            </div>
            <button onclick="generateKeys()">Generate Key Pair</button>
            <div id="keyResult" class="result" style="display: none;"></div>
        </div>
        
        <!-- Document Signing Section -->
        <div class="section">
            <h2>2. Sign Document</h2>
            <form id="signForm" enctype="multipart/form-data">
                <div class="form-group">
                    <label for="documentToSign">Document to Sign:</label>
                    <input type="file" id="documentToSign" name="document" required>
                </div>
                <div class="form-group">
                    <label for="privateKeySign">Private Key (PEM format):</label>
                    <textarea id="privateKeySign" name="privateKey" placeholder="-----BEGIN PRIVATE KEY-----" required></textarea>
                </div>
                <button type="submit">Sign Document</button>
            </form>
            <div id="signResult" class="result" style="display: none;"></div>
        </div>
        
        <!-- Signature Verification Section -->
        <div class="section">
            <h2>3. Verify Signature</h2>
            <form id="verifyForm" enctype="multipart/form-data">
                <div class="form-group">
                    <label for="documentToVerify">Document to Verify:</label>
                    <input type="file" id="documentToVerify" name="document" required>
                </div>
                <div class="form-group">
                    <label for="publicKeyVerify">Public Key (PEM format):</label>
                    <textarea id="publicKeyVerify" name="publicKey" placeholder="-----BEGIN PUBLIC KEY-----" required></textarea>
                </div>
                <div class="form-group">
                    <label for="signatureFile">Signature File (.md5):</label>
                    <input type="file" id="signatureFile" name="signatureFile" accept=".md5" required>
                </div>
                <button type="submit">Verify Signature</button>
            </form>
            <div id="verifyResult" class="result" style="display: none;"></div>
        </div>
    </div>

    <script>
        function generateKeys() {
            const keySize = document.getElementById('keySize').value;
            const resultDiv = document.getElementById('keyResult');
            
            const formData = new FormData();
            formData.append('keySize', keySize);
            
            fetch('/api/generate-keys', {
                method: 'POST',
                body: formData
            })
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    resultDiv.className = 'result error';
                    resultDiv.textContent = 'Error: ' + data.error;
                } else {
                    resultDiv.className = 'result success';
                    resultDiv.innerHTML = '<strong>Private Key:</strong>\n' + data.privateKey + '\n\n<strong>Public Key:</strong>\n' + data.publicKey;
                    
                    // Auto-fill the private key in sign form
                    document.getElementById('privateKeySign').value = data.privateKey;
                    // Auto-fill the public key in verify form
                    document.getElementById('publicKeyVerify').value = data.publicKey;
                }
                resultDiv.style.display = 'block';
            })
            .catch(error => {
                resultDiv.className = 'result error';
                resultDiv.textContent = 'Error: ' + error.message;
                resultDiv.style.display = 'block';
            });
        }
        
        document.getElementById('signForm').addEventListener('submit', function(e) {
            e.preventDefault();
            const resultDiv = document.getElementById('signResult');
            const formData = new FormData(this);
            
            resultDiv.textContent = 'Signing document...';
            resultDiv.className = 'result';
            resultDiv.style.display = 'block';
            
            fetch('/api/sign-document', {
                method: 'POST',
                body: formData
            })
            .then(response => {
                if (response.ok) {
                    return response.blob();
                } else {
                    return response.text().then(text => Promise.reject(text));
                }
            })
            .then(blob => {
                // Get the filename from the response headers or use default
                const contentDisposition = fetch.response?.headers?.get('content-disposition');
                let filename = 'document_sig.md5';
                if (contentDisposition) {
                    const filenameMatch = contentDisposition.match(/filename="(.+)"/);
                    if (filenameMatch) {
                        filename = filenameMatch[1];
                    }
                }
                
                // Create download link
                const url = URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = filename;
                document.body.appendChild(a);
                a.click();
                document.body.removeChild(a);
                URL.revokeObjectURL(url);
                
                resultDiv.className = 'result success';
                resultDiv.textContent = 'Document signed successfully! Signature file downloaded as: ' + filename;
            })
            .catch(error => {
                resultDiv.className = 'result error';
                resultDiv.textContent = 'Error: ' + error;
            });
        });
        
        document.getElementById('verifyForm').addEventListener('submit', function(e) {
            e.preventDefault();
            const resultDiv = document.getElementById('verifyResult');
            const formData = new FormData(this);
            
            fetch('/api/verify-signature', {
                method: 'POST',
                body: formData
            })
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    resultDiv.className = 'result error';
                    resultDiv.textContent = 'Error: ' + data.error;
                } else {
                    resultDiv.className = data.valid ? 'result success' : 'result error';
                    resultDiv.innerHTML = '<strong>Verification Result:</strong> ' + (data.valid ? 'VALID ✓' : 'INVALID ✗') + '\n<strong>Message:</strong> ' + data.message;
                }
                resultDiv.style.display = 'block';
            })
            .catch(error => {
                resultDiv.className = 'result error';
                resultDiv.textContent = 'Error: ' + error.message;
                resultDiv.style.display = 'block';
            });
        });
    </script>
</body>
</html>`

	t, _ := template.New("digital-signature").Parse(tmpl)
	t.Execute(w, nil)
}

// SteganographyPageHandler serves the steganography page
func SteganographyPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Steganography - LSB</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background-color: #f5f5f5; }
        .container { max-width: 900px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { text-align: center; color: #333; margin-bottom: 30px; }
        .section { margin: 20px 0; padding: 20px; border: 1px solid #ddd; border-radius: 8px; background: #fafafa; }
        .form-group { margin: 15px 0; }
        label { display: block; margin-bottom: 5px; font-weight: bold; color: #555; }
        input[type="file"], textarea { 
            width: 100%; 
            padding: 10px; 
            border: 1px solid #ddd; 
            border-radius: 4px; 
            box-sizing: border-box; 
        }
        textarea { height: 100px; resize: vertical; }
        button { 
            background: linear-gradient(135deg, #28a745 0%, #20c997 100%); 
            color: white; 
            padding: 12px 24px; 
            border: none; 
            border-radius: 6px; 
            cursor: pointer; 
            font-size: 16px; 
            margin: 5px;
            transition: transform 0.2s;
        }
        button:hover { transform: translateY(-1px); }
        .result { 
            margin-top: 20px; 
            padding: 15px; 
            border-radius: 4px; 
            white-space: pre-wrap; 
            word-break: break-all; 
        }
        .success { background-color: #d4edda; color: #155724; border: 1px solid #c3e6cb; }
        .error { background-color: #f8d7da; color: #721c24; border: 1px solid #f5c6cb; }
        .nav-link { display: inline-block; margin: 10px; padding: 10px 20px; background: #6c757d; color: white; text-decoration: none; border-radius: 4px; }
        .nav-link:hover { background: #5a6268; }
        .preview-image { max-width: 300px; max-height: 200px; margin: 10px 0; border: 1px solid #ddd; }
    </style>
</head>
<body>
    <div class="container">
        <div style="text-align: center; margin-bottom: 20px;">
            <a href="/" class="nav-link">🏠 Home</a>
            <a href="/digital-signature" class="nav-link">🔐 Digital Signature</a>
            <a href="/combined" class="nav-link">🔧 Combined</a>
        </div>
        
        <h1>🖼️ Steganography - LSB Method</h1>
        
        <!-- Hide Message Section -->
        <div class="section">
            <h2>1. Hide Message in Image</h2>
            <form id="hideForm" enctype="multipart/form-data">
                <div class="form-group">
                    <label for="imageToHide">Cover Image (PNG/JPEG):</label>
                    <input type="file" id="imageToHide" name="image" accept="image/png,image/jpeg" required onchange="previewImage(this, 'hidePreview')">
                    <img id="hidePreview" class="preview-image" style="display: none;">
                </div>
                <div class="form-group">
                    <label for="messageToHide">Message to Hide:</label>
                    <textarea id="messageToHide" name="message" placeholder="Enter your secret message here..." required></textarea>
                </div>
                <button type="submit">Hide Message</button>
            </form>
            <div id="hideResult" class="result" style="display: none;"></div>
        </div>
        
        <!-- Extract Message Section -->
        <div class="section">
            <h2>2. Extract Message from Image</h2>
            <form id="extractForm" enctype="multipart/form-data">
                <div class="form-group">
                    <label for="imageToExtract">Stego Image (PNG/JPEG):</label>
                    <input type="file" id="imageToExtract" name="image" accept="image/png,image/jpeg" required onchange="previewImage(this, 'extractPreview')">
                    <img id="extractPreview" class="preview-image" style="display: none;">
                </div>
                <button type="submit">Extract Message</button>
            </form>
            <div id="extractResult" class="result" style="display: none;"></div>
        </div>
    </div>

    <script>
        function previewImage(input, previewId) {
            const preview = document.getElementById(previewId);
            if (input.files && input.files[0]) {
                const reader = new FileReader();
                reader.onload = function(e) {
                    preview.src = e.target.result;
                    preview.style.display = 'block';
                };
                reader.readAsDataURL(input.files[0]);
            }
        }
        
        document.getElementById('hideForm').addEventListener('submit', function(e) {
            e.preventDefault();
            const resultDiv = document.getElementById('hideResult');
            const formData = new FormData(this);
            
            resultDiv.textContent = 'Processing...';
            resultDiv.className = 'result';
            resultDiv.style.display = 'block';
            
            fetch('/api/hide-message', {
                method: 'POST',
                body: formData
            })
            .then(response => {
                if (response.ok) {
                    return response.blob();
                } else {
                    return response.text().then(text => Promise.reject(text));
                }
            })
            .then(blob => {
                const url = URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = 'stego_image.' + (blob.type.includes('png') ? 'png' : 'jpg');
                document.body.appendChild(a);
                a.click();
                document.body.removeChild(a);
                URL.revokeObjectURL(url);
                
                resultDiv.className = 'result success';
                resultDiv.textContent = 'Message hidden successfully! Download started.';
            })
            .catch(error => {
                resultDiv.className = 'result error';
                resultDiv.textContent = 'Error: ' + error;
            });
        });
        
        document.getElementById('extractForm').addEventListener('submit', function(e) {
            e.preventDefault();
            const resultDiv = document.getElementById('extractResult');
            const formData = new FormData(this);
            
            resultDiv.textContent = 'Extracting...';
            resultDiv.className = 'result';
            resultDiv.style.display = 'block';
            
            fetch('/api/extract-message', {
                method: 'POST',
                body: formData
            })
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    resultDiv.className = 'result error';
                    resultDiv.textContent = 'Error: ' + data.error;
                } else {
                    resultDiv.className = 'result success';
                    resultDiv.innerHTML = '<strong>Extracted Message:</strong>\n' + data.message;
                }
            })
            .catch(error => {
                resultDiv.className = 'result error';
                resultDiv.textContent = 'Error: ' + error.message;
            });
        });
    </script>
</body>
</html>`

	t, _ := template.New("steganography").Parse(tmpl)
	t.Execute(w, nil)
}

// CombinedPageHandler serves the combined operations page
func CombinedPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Combined Operations - Digital Signature & Steganography</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background-color: #f5f5f5; }
        .container { max-width: 900px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { text-align: center; color: #333; margin-bottom: 30px; }
        .section { margin: 20px 0; padding: 20px; border: 1px solid #ddd; border-radius: 8px; background: #fafafa; }
        .form-group { margin: 15px 0; }
        label { display: block; margin-bottom: 5px; font-weight: bold; color: #555; }
        input[type="file"], textarea { 
            width: 100%; 
            padding: 10px; 
            border: 1px solid #ddd; 
            border-radius: 4px; 
            box-sizing: border-box; 
        }
        textarea { height: 100px; resize: vertical; }
        button { 
            background: linear-gradient(135deg, #dc3545 0%, #fd7e14 100%); 
            color: white; 
            padding: 12px 24px; 
            border: none; 
            border-radius: 6px; 
            cursor: pointer; 
            font-size: 16px; 
            margin: 5px;
            transition: transform 0.2s;
        }
        button:hover { transform: translateY(-1px); }
        .result { 
            margin-top: 20px; 
            padding: 15px; 
            border-radius: 4px; 
            white-space: pre-wrap; 
            word-break: break-all; 
        }
        .success { background-color: #d4edda; color: #155724; border: 1px solid #c3e6cb; }
        .error { background-color: #f8d7da; color: #721c24; border: 1px solid #f5c6cb; }
        .nav-link { display: inline-block; margin: 10px; padding: 10px 20px; background: #6c757d; color: white; text-decoration: none; border-radius: 4px; }
        .nav-link:hover { background: #5a6268; }
        .preview-image { max-width: 300px; max-height: 200px; margin: 10px 0; border: 1px solid #ddd; }
        .info-box { background: #e7f3ff; border: 1px solid #b8daff; padding: 15px; border-radius: 4px; margin: 15px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div style="text-align: center; margin-bottom: 20px;">
            <a href="/" class="nav-link">🏠 Home</a>
            <a href="/digital-signature" class="nav-link">🔐 Digital Signature</a>
            <a href="/steganography" class="nav-link">🖼️ Steganography</a>
        </div>
        
        <h1>🔧 Combined Operations</h1>
        
        <div class="info-box">
            <strong>Combined Operations:</strong> This section allows you to combine digital signature and steganography operations. 
            You can sign a document and hide both the document and its signature in an image, or extract and verify a signed document from an image.
        </div>
        
        <!-- Sign and Hide Section -->
        <div class="section">
            <h2>1. Sign Document and Hide in Image</h2>
            <form id="signHideForm" enctype="multipart/form-data">
                <div class="form-group">
                    <label for="documentToSignHide">Document to Sign:</label>
                    <input type="file" id="documentToSignHide" name="document" required>
                </div>
                <div class="form-group">
                    <label for="imageToSignHide">Cover Image (PNG/JPEG):</label>
                    <input type="file" id="imageToSignHide" name="image" accept="image/png,image/jpeg" required onchange="previewImage(this, 'signHidePreview')">
                    <img id="signHidePreview" class="preview-image" style="display: none;">
                </div>
                <div class="form-group">
                    <label for="privateKeySignHide">Private Key (PEM format):</label>
                    <textarea id="privateKeySignHide" name="privateKey" placeholder="-----BEGIN PRIVATE KEY-----" required></textarea>
                </div>
                <button type="submit">Sign and Hide</button>
            </form>
            <div id="signHideResult" class="result" style="display: none;"></div>
        </div>
        
        <!-- Extract and Verify Section -->
        <div class="section">
            <h2>2. Extract and Verify Signed Document</h2>
            <form id="extractVerifyForm" enctype="multipart/form-data">
                <div class="form-group">
                    <label for="imageToExtractVerify">Stego Image (PNG/JPEG):</label>
                    <input type="file" id="imageToExtractVerify" name="image" accept="image/png,image/jpeg" required onchange="previewImage(this, 'extractVerifyPreview')">
                    <img id="extractVerifyPreview" class="preview-image" style="display: none;">
                </div>
                <div class="form-group">
                    <label for="publicKeyExtractVerify">Public Key (PEM format):</label>
                    <textarea id="publicKeyExtractVerify" name="publicKey" placeholder="-----BEGIN PUBLIC KEY-----" required></textarea>
                </div>
                <button type="submit">Extract and Verify</button>
            </form>
            <div id="extractVerifyResult" class="result" style="display: none;"></div>
        </div>
    </div>

    <script>
        function previewImage(input, previewId) {
            const preview = document.getElementById(previewId);
            if (input.files && input.files[0]) {
                const reader = new FileReader();
                reader.onload = function(e) {
                    preview.src = e.target.result;
                    preview.style.display = 'block';
                };
                reader.readAsDataURL(input.files[0]);
            }
        }
        
        document.getElementById('signHideForm').addEventListener('submit', function(e) {
            e.preventDefault();
            const resultDiv = document.getElementById('signHideResult');
            const formData = new FormData(this);
            
            resultDiv.textContent = 'Processing...';
            resultDiv.className = 'result';
            resultDiv.style.display = 'block';
            
            fetch('/api/sign-and-hide', {
                method: 'POST',
                body: formData
            })
            .then(response => {
                if (response.ok) {
                    return response.blob();
                } else {
                    return response.text().then(text => Promise.reject(text));
                }
            })
            .then(blob => {
                const url = URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = 'signed_stego_image.' + (blob.type.includes('png') ? 'png' : 'jpg');
                document.body.appendChild(a);
                a.click();
                document.body.removeChild(a);
                URL.revokeObjectURL(url);
                
                resultDiv.className = 'result success';
                resultDiv.textContent = 'Document signed and hidden successfully! Download started.';
            })
            .catch(error => {
                resultDiv.className = 'result error';
                resultDiv.textContent = 'Error: ' + error;
            });
        });
        
        document.getElementById('extractVerifyForm').addEventListener('submit', function(e) {
            e.preventDefault();
            const resultDiv = document.getElementById('extractVerifyResult');
            const formData = new FormData(this);
            
            resultDiv.textContent = 'Processing...';
            resultDiv.className = 'result';
            resultDiv.style.display = 'block';
            
            fetch('/api/extract-and-verify', {
                method: 'POST',
                body: formData
            })
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    resultDiv.className = 'result error';
                    resultDiv.textContent = 'Error: ' + data.error;
                } else {
                    const status = data.valid ? 'VALID ✓' : 'INVALID ✗';
                    const statusClass = data.valid ? 'success' : 'error';
                    resultDiv.className = 'result ' + statusClass;
                    resultDiv.innerHTML = '<strong>Verification Status:</strong> ' + status + 
                                         '\n<strong>Extracted Document:</strong>\n' + data.document +
                                         '\n<strong>Signature:</strong>\n' + data.signature +
                                         (data.message ? '\n<strong>Message:</strong> ' + data.message : '');
                }
            })
            .catch(error => {
                resultDiv.className = 'result error';
                resultDiv.textContent = 'Error: ' + error.message;
            });
        });
    </script>
</body>
</html>`

	t, _ := template.New("combined").Parse(tmpl)
	t.Execute(w, nil)
}

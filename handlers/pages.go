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
                    <label for="signatureFile">Signature File (.rsasig):</label>
                    <input type="file" id="signatureFile" name="signatureFile" accept=".rsasig" required>
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
                let filename = 'document_sig.rsasig';
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

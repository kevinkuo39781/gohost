# gohost

🚀 **gohost** is a lightweight static page hosting server written in Go. It allows you to upload a ZIP archive of your static website and deploy it instantly under a custom path like `/g/your-site`.

---

## ✨ Features

- ✅ Upload a ZIP file and auto-deploy the content
- ✅ Serve at custom paths (e.g. `/g/my-site`)
- ✅ Pure static file hosting, no database
- ✅ Docker-ready for easy deployment

---

## 📦 Getting Started

### Option 1: Run with Go

```bash
go run cmd/gohost/main.go
```

### Option 2: Run with Docker

```bash
docker build -t gohost .
docker run -p 8080:8080 -v $(pwd)/data:/app/data gohost
```

---

## 📤 Upload Example

You can upload a ZIP file using `curl`:

```bash
curl -F "file=@your-site.zip" -F "path=my-site" http://localhost:8080/upload
```

Or open in the browser:

```
http://localhost:8080/
```

Then visit your site at:

```
http://localhost:8080/g/my-site/index.html
```

---

## 📁 Project Structure

```
.
├── cmd/gohost/main.go         # Main entrypoint
├── public/index.html          # Upload page UI
├── data/sites/                # Uploaded and extracted websites
├── Dockerfile
├── .gitignore
├── LICENSE
└── README.md
```

---

## 🔧 Roadmap

- Web UI for managing uploaded sites
- Token-based authentication
- Automatic site expiration
- HTTPS support with Let's Encrypt

---

## 📝 License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.

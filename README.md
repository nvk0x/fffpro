# fffpro

**fffpro** is a fast, resilient, and production-ready URL fetcher designed for bug bounty hunters and recon workflows.

It is an improved and more stable version of the original **fff**, built to handle large-scale URL lists without crashing, stalling, or silently failing.

---

## 🚀 Why fffpro?

The original tool is extremely fast, but can struggle with:

* timeouts causing dropped requests
* lack of retry logic
* uncontrolled concurrency (goroutine spikes)
* instability on large target lists

**fffpro fixes all of that**, while keeping the same philosophy:
👉 *simple, fast, stdin → stdout workflow*

---

## ⚡ Features

* 🔥 High-performance worker pool (no goroutine explosion)
* 🔁 Automatic retries with backoff
* ⏱ Configurable timeout handling
* 🌐 Smart dead-host skipping
* 💾 Response saving (compatible with fff-style workflow)
* 📂 Organized output structure
* 🧠 Built for large-scale recon & bug bounty pipelines

---

## 📦 Installation

```bash
go install github.com/nvk0x/fffpro@latest
```

Make sure `$GOPATH/bin` is in your `$PATH`.

---

## 🛠 Usage

```bash
cat urls.txt | fffpro
```

### Example

```bash
cat urls.txt | fffpro -w 100 -timeout 30 -retries 3
```

---

## ⚙️ Options

| Flag             | Description               | Default |
| ---------------- | ------------------------- | ------- |
| `-w`             | Number of workers         | 50      |
| `-timeout`       | Request timeout (seconds) | 25      |
| `-retries`       | Retry attempts            | 3       |
| `-o`             | Output directory          | out     |
| `-S`             | Save all responses        | true    |
| `--ignore-empty` | Skip empty responses      | true    |

---

## 📂 Output Structure

Responses are saved in a structured format:

```
out/
 └── target.com/
      ├── <hash>.body
      └── <hash>.headers
```

---

## 🔗 Recommended Workflow

Combine with tools like httpx:

```bash
cat urls.txt \
| httpx -silent -threads 50 -timeout 20 \
| fffpro
```

---

## 🙏 Credits

This tool is heavily inspired by the original:

* GitHub: https://github.com/tomnomnom/fff
* Author: Tom Hudson
* Twitter: https://twitter.com/tomnomnom

fffpro builds upon the idea and improves stability, scalability, and usability for modern recon workflows.

---

## ⚠️ Disclaimer

This tool is intended for **authorized security testing and research only**.
Do not use it against systems without proper permission.

---

## ⭐ Contributing

Pull requests, ideas, and improvements are welcome.

---

## 📜 License

MIT License


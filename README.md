# flavor‑forge  
*Build and run reproducible ML experiment pipelines with a single declarative YAML config.*

> Built by **Myroslav Mokhammad Abdeljawwad** to solve the pain of juggling data‑preparation scripts, model training, and evaluation steps across different environments.

---

## ✨ Features

- **Single source of truth:** One YAML file (`experiment.yml`) defines datasets, preprocessing, models, hyper‑parameters, and metrics.
- **Reproducible pipelines:** Every step is deterministic – versioned data files, fixed random seeds, and isolated Docker containers (via Go’s exec wrapper).
- **Extensible executor:** Plug in new stages (e.g., feature engineering, ensembling) without touching the core logic.
- **Rich logging:** Structured logs with `pkg/logger` that can be routed to stdout or a file.
- **Cross‑platform CLI:** Works on Linux/macOS/Windows; built with Go for zero runtime dependencies.

---

## 🚀 Installation

```bash
# Clone the repo
git clone https://github.com/<your-github>/flavor-forge.git
cd flavor-forge

# Build the binary (requires Go 1.22+)
go build -o flavor-forge ./cmd

# Or install via Go modules
go install ./cmd@latest
```

The `flavor-forge` executable will be available in your `$GOPATH/bin` or the current directory.

---

## 📦 Usage

### 1️⃣ Create an experiment config  
Place this under `examples/experiment.yml` (or any path you prefer):

```yaml
pipeline:
  - name: load_data
    type: csv_loader
    params:
      path: data/train.csv
  - name: preprocess
    type: standard_scaler
  - name: train
    type: logistic_regression
    params:
      max_iter: 200
  - name: evaluate
    type: accuracy
```

### 2️⃣ Run the pipeline

```bash
# Using the shipped shell helper (requires bash)
./scripts/run_experiment.sh examples/experiment.yml
```

or directly with Go:

```bash
flavor-forge run examples/experiment.yml
```

You’ll see structured logs:

```
[INFO] Starting pipeline: 4 steps
[DEBUG] Step 1/4: load_data
[DEBUG] Step 2/4: preprocess
[DEBUG] Step 3/4: train
[DEBUG] Step 4/4: evaluate
[INFO] Pipeline finished successfully. Accuracy: 0.87
```

---

## 🤝 Contributing

I created *flavor‑forge* to help my own data science team, but I’d love to see it grow:

1. Fork the repo  
2. Create a feature branch (`git checkout -b feat/your-feature`)  
3. Run `go test ./...` – make sure all tests pass  
4. Submit a pull request

Please keep the code style consistent with Go’s conventions and add documentation/comments where necessary.

---

## 📄 License

MIT © 2026 Myroslav Mokhammad Abdeljawwad

---

## 🔗 See Also

- [Boost Legacy Java Refactoring with Copilot’s AI API](https://dev.to/myroslavmokhammadabd/boost-legacy-java-refactoring-with-copilots-ai-api-2epb) – a blog post where I first discussed the challenges that led to *flavor‑forge*.

---

> “When you finally get your experiments reproducible, the real magic happens.” – **Myroslav Mokhammad Abdeljawwad**
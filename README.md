# BricoLink Image Processing Service

A standalone Go service that handles image resizing for the BricoLink platform. Receives uploaded images via HTTP, generates three optimized WebP variants (thumb, medium, large) in parallel using libvips, and saves them to shared storage.

Built as a companion service to the Laravel monolith. Laravel handles uploads and database logic; this service handles the CPU-intensive image work.

---

## Why This Exists

Laravel runs on PHP-FPM, which assigns one worker per request. Resizing images inside Laravel blocks that worker for the entire duration — under concurrent uploads, this starves unrelated requests.

This Go service offloads image processing to a separate process that:
- Uses libvips (C library), which is 3–6× faster than PHP's GD library
- Generates all three variants in parallel using goroutines
- Doesn't block any PHP-FPM workers

### Benchmark Results

Tested with 15 images, resizing each to 3 WebP variants (200px, 800px, 1600px):

| Runner              | Time    | Notes                          |
|---------------------|---------|--------------------------------|
| Go + libvips (parallel) | 1.73s  | 3 variants generated concurrently |
| Go + libvips (sequential) | 4.73s  | one at a time                  |
| PHP + Imagick       | 11.39s  | sequential, single-threaded    |
| PHP + GD            | —       | failed on some input formats   |

---

## Architecture

```
Laravel (web app)
  │
  │  POST /resize  (multipart: image file)
  │
  ▼
Go Image Service (:8081)
  │
  ├── Parses upload
  ├── Generates 3 variants in parallel (goroutines)
  │     ├── thumb   (200px wide)
  │     ├── medium  (800px wide)
  │     └── large   (1600px wide)
  ├── Saves as WebP to shared storage directory
  └── Returns JSON with filenames
          │
          ▼
      Laravel stores filenames in database
      Browser loads images via storage symlink
```

Both services read/write the same storage directory. Laravel's `php artisan storage:link` makes the files accessible to browsers.

---

## Prerequisites

- **Go 1.22+**
- **libvips** (C library for image processing)

### Install libvips

Arch Linux:
```bash
sudo pacman -S libvips
```

Ubuntu/Debian:
```bash
sudo apt install libvips-dev
```

---

## Setup

### 1. Clone the repository

```bash
git clone <repository-url>
cd image-processing-service
```

### 2. Install Go dependencies

```bash
go mod download
```

### 3. Configure environment variables

Copy the example file:
```bash
cp .env.example .env
```

Edit `.env` or export the variables directly:

| Variable            | Description                                      | Default                  |
|---------------------|--------------------------------------------------|--------------------------|
| `IMGPROC_PORT`      | Port the service listens on                      | `8081`                   |
| `IMGPROC_STORAGE_PATH` | Absolute path to the shared storage directory | `./storage/portfolios`   |

**Important for Laravel integration:** `IMGPROC_STORAGE_PATH` must point to Laravel's public storage directory so both services share the same files:

```bash
export IMGPROC_STORAGE_PATH=/path/to/bricolink-laravel/storage/app/public/portfolios
export IMGPROC_PORT=8081
```

The service creates the directory if it doesn't exist.

### 4. Run the service

```bash
go run ./cmd/imgproc
```

You should see:
```
listening on :8081 (default prot if not changed)
```

### 5. Verify it works

```bash
# Health check
curl http://localhost:8081/health
# {"status":"ok"}

# Resize an image
curl -X POST -F "image=@/path/to/photo.jpg" http://localhost:8081/resize
# {"thumb":"uuid_thumb.webp","medium":"uuid_medium.webp","large":"uuid_large.webp"}
```

---

## API

### `GET /health`

Returns service status.

**Response:**
```json
{"status": "ok"}
```

### `POST /resize`

Accepts an image file, generates three WebP variants, saves them to storage.

**Request:**
- Content-Type: `multipart/form-data`
- Field name: `image`
- Max file size: 10 MB

**Success response (200):**
```json
{
  "thumb":  "64c7239f_thumb.webp",
  "medium": "64c7239f_medium.webp",
  "large":  "64c7239f_large.webp"
}
```

**Error responses:**

| Status | Meaning                        |
|--------|--------------------------------|
| 400    | Missing image field or invalid upload |
| 413    | File exceeds 10 MB             |
| 500    | Resize or disk write failed    |

### Variant sizes

| Variant | Width  | Use case                         |
|---------|--------|----------------------------------|
| thumb   | 200px  | Artisan cards in search results  |
| medium  | 800px  | Portfolio grid on profile page   |
| large   | 1600px | Full-size image view (lightbox)  |

All variants: WebP format, quality 80, aspect ratio preserved.

---

## Laravel Integration

### 1. Add environment variable

In Laravel's `.env`:
```
IMGPROC_URL=http://localhost:8081
```

In `config/services.php`:
```php
'imgproc' => [
    'url' => env('IMGPROC_URL', 'http://localhost:8081'),
],
```

### 2. Create service class

`app/Services/ImageProcessor.php`:
```php
<?php

namespace App\Services;

use Illuminate\Http\UploadedFile;
use Illuminate\Support\Facades\Http;

class ImageProcessor
{
    public function process(UploadedFile $file): array
    {
        $response = Http::attach(
            'image',
            file_get_contents($file->path()),
            $file->getClientOriginalName()
        )->post(config('services.imgproc.url') . '/resize');

        if ($response->failed()) {
            throw new \RuntimeException('Image processing failed: ' . $response->body());
        }

        return $response->json();
    }
}
```

### 3. Use in controller

```php
public function store(Request $request, ImageProcessor $imgproc)
{
    $request->validate(['image' => 'required|image|max:10240']);

    $result = $imgproc->process($request->file('image'));

    PortfolioImage::create([
        'artisan_id'  => auth()->id(),
        'thumb_path'  => $result['thumb'],
        'medium_path' => $result['medium'],
        'large_path'  => $result['large'],
    ]);

    return back()->with('success', 'Image uploaded');
}
```

### 4. Display in Blade

```html
<!-- In search results (small) -->
<img src="{{ asset('storage/portfolios/' . $image->thumb_path) }}">

<!-- In portfolio grid -->
<img src="{{ asset('storage/portfolios/' . $image->medium_path) }}">

<!-- In lightbox / full view -->
<img src="{{ asset('storage/portfolios/' . $image->large_path) }}">
```

### 5. Migration update

Add columns to `portfolio_images` table:
```php
$table->string('thumb_path');
$table->string('medium_path');
$table->string('large_path');
```

---

## Project Structure

```
image-processing-service/
├── cmd/
│   └── imgproc/
│       └── main.go              # Entry point: config, wiring, server start
├── internal/
│   ├── config/
│   │   └── config.go            # Reads env vars with defaults
│   ├── httpx/
│   │   ├── health.go            # GET /health handler
│   │   └── middleware.go        # Request logging + panic recovery
│   └── imgproc/
│       ├── handler.go           # POST /resize HTTP handler
│       ├── resizer.go           # libvips resize logic
│       ├── resizer_test.go      # Tests
│       └── testdata/
│           └── test.jpg         # Test fixture
├── go.mod
├── go.sum
├── .env.example
└── README.md
```

---

## Running Tests

```bash
go test ./internal/imgproc/ -v
```

---

## Troubleshooting

**VIPS-WARNING: unable to load vips-openslide.so**
Harmless. Optional plugin for medical imaging. Suppress with:
```bash
VIPS_WARNING=0 go run ./cmd/imgproc
```

**"Unsupported image format" on resize**
The uploaded file isn't a valid image (or a format libvips doesn't handle). Supported: JPEG, PNG, WebP, TIFF, GIF.

**Storage directory not created**
Ensure `IMGPROC_STORAGE_PATH` is an absolute path (starts with `/`). Relative paths and `~` are resolved from the working directory, not the home directory.

**Connection refused from Laravel**
The Go service must be running before Laravel tries to call it. Start with `go run ./cmd/imgproc` in a separate terminal.

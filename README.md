# R2Mirror

## 🚧 Work in Progress

R2Mirror is a Go-based tool that mirrors Ubuntu releases to a Cloudflare R2 bucket. It periodically checks for new releases, streams them to your R2 bucket, creating a public mirror of the Ubuntu releases you specify.

<p align="center"> <img src="https://img.shields.io/github/stars/adhitht/R2Mirror?style=for-the-badge" alt="GitHub Stars" /> <img src="https://img.shields.io/github/forks/adhitht/R2Mirror?style=for-the-badge" alt="GitHub Forks" /> <img src="https://img.shields.io/github/issues/adhitht/R2Mirror?style=for-the-badge" alt="GitHub Issues" /> <img src="https://img.shields.io/github/license/adhitht/R2Mirror?style=for-the-badge" alt="License" /> <img src="https://img.shields.io/github/go-mod/go-version/adhitht/R2Mirror?style=for-the-badge" alt="Go Version" /> </p>

## Features

- **Automated Mirroring**: Automatically downloads and mirrors specified Ubuntu releases.
- **Cloudflare R2 Integration**: Seamlessly uploads files to your Cloudflare R2 bucket.
- **HTML Indexing**: Generates an index page for all mirrored releases and a separate index for each release version.
- **Configuration Watching**: Automatically updates the mirror when you change the configuration file.

## How It Works

1.  **Configuration**: You specify the Ubuntu releases you want to mirror in the `config.yaml` file and secrets in `.env`
2.  **Fetching**: R2Mirror fetches the list of files for each specified release from the official Ubuntu release server.
3.  **Downloading & Uploading**: It downloads each file and uploads it directly to your R2 bucket under a path corresponding to the release version (e.g., `25.04/`).
4.  **Indexing**: After all files for a release are uploaded, R2Mirror generates an `index.html` file for that release and uploads it to the corresponding directory in your R2 bucket. It also creates a main `index.html` file that links to all the mirrored releases.
5.  **Watching**: R2Mirror watches the `config.yaml` file for any changes and automatically re-runs the mirroring process if the file is updated.

## Prerequisites

- Go 1.18 or later
- A Cloudflare account with an R2 bucket
- AWS credentials configured for R2 access (see [Cloudflare R2 documentation](https://developers.cloudflare.com/r2/api/s3/api/))

## Installation

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/adhitht/R2Mirror.git
    cd R2Mirror
    ```

2.  **Install dependencies:**

    ```bash
    go mod tidy
    ```

3.  **Configure the application:**

    - Rename `.env.example` to `.env` and fill in your AWS/R2 credentials.
    - Edit `config.yaml` to specify the Ubuntu releases you want to mirror and your R2 bucket name.

4.  **Run the application:**

    ```bash
    go run main.go
    ```

## Configuration

### `config.yaml`

The `config.yaml` file is used to configure the releases you want to mirror and your R2 bucket details.

```yaml
releases:
  - "25.04"
  - "24.10"
  - "24.04"
  - "22.04"

# R2/S3 Configuration
bucket: "your-bucket-name"
region: "auto"  # Use "auto" for Cloudflare R2
```

## Usage

Once the application is running, it will start downloading and mirroring the specified Ubuntu releases to your R2 bucket. You can monitor the progress in the console.

The mirrored files will be available at the public URL of your R2 bucket.

## Roadmap

- [ ] Fix bugs in index.html generation 
- [ ] Cron job to update automatically
- [ ] Support for other distributions (e.g., Debian, Fedora)
- [ ] Support for other storage backends (e.g., Google Cloud Storage, Azure Blob Storage)
- [ ] Improved error handling and retries
- [ ] Performance optimizations for faster downloads and uploads
- [ ] A web interface to configure the mirror and view the status

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

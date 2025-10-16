# Digiflazz Gateway API

Aplikasi gateway untuk integrasi API transaksi ke Digiflazz menggunakan Golang. Aplikasi ini berfungsi sebagai middleware yang menghubungkan sistem internal dengan API Digiflazz untuk berbagai layanan seperti pulsa, token listrik, PDAM, dll.

## 🚀 Features

- **Balance Check**: Cek saldo akun Digiflazz
- **Price List**: Daftar harga produk Prabayar dan Pascabayar
- **Transaction Processing**: Topup pulsa, token listrik, pembayaran tagihan
- **Multi-Provider Support**: IRS, FM, Otomax, ST24, Payuni, Sipas, Tiger
- **Webhook Handling**: Proses callback dari Digiflazz
- **Status Checking**: Cek status transaksi real-time

## 📋 Prerequisites

- Go 1.21 or higher
- Git
- Docker (optional, for containerization)

## 🛠️ Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd gateway-digiflazz
```

2. Install dependencies:
```bash
go mod tidy
```

3. Copy environment configuration:
```bash
cp configs/.env.example .env
```

4. Configure your environment variables in `.env`:
```bash
# Digiflazz API Configuration
DIGIFLAZZ_USERNAME=your_username
DIGIFLAZZ_API_KEY=your_api_key
DIGIFLAZZ_BASE_URL=https://api.digiflazz.com

# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
LOG_LEVEL=info
```

5. Run the application:
```bash
go run cmd/server/main.go
```

## 🏗️ Project Structure

```
gateway-digiflazz/
├── cmd/                    # Application entry points
│   └── server/
│       └── main.go
├── internal/               # Private application code
│   ├── config/            # Configuration management
│   ├── handlers/          # HTTP handlers
│   ├── services/          # Business logic
│   ├── models/            # Data models
│   ├── repositories/      # Data access layer
│   └── middleware/        # HTTP middleware
├── pkg/                   # Public library code
│   ├── digiflazz/         # Digiflazz API client
│   └── utils/             # Utility functions
├── api/                   # API definitions
├── configs/               # Configuration files
├── docs/                  # Documentation
├── tests/                 # Test files
├── go.mod
├── go.sum
└── README.md
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DIGIFLAZZ_USERNAME` | Digiflazz username | - |
| `DIGIFLAZZ_API_KEY` | Digiflazz API key | - |
| `DIGIFLAZZ_BASE_URL` | Digiflazz API base URL | https://api.digiflazz.com |
| `SERVER_PORT` | Server port | 8080 |
| `SERVER_HOST` | Server host | 0.0.0.0 |
| `LOG_LEVEL` | Log level | info |

## 📚 API Documentation

### Endpoints

#### Balance Check
```http
GET /api/v1/balance
```

#### Price List
```http
GET /api/v1/prices?type=prabayar
GET /api/v1/prices?type=pascabayar
```

#### Transaction
```http
POST /api/v1/transaction/topup
POST /api/v1/transaction/pay
```

#### Status Check
```http
GET /api/v1/transaction/{ref_id}/status
```

## 🧪 Testing

Run tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## 🐳 Docker

Build Docker image:
```bash
docker build -t gateway-digiflazz .
```

Run with Docker:
```bash
docker run -p 8080:8080 --env-file .env gateway-digiflazz
```

## 📖 Documentation

- [Digiflazz API Documentation](https://developer.digiflazz.com/api/)
- [Project Architecture](docs/architecture.md)
- [API Reference](docs/api-reference.md)
- [Deployment Guide](docs/deployment.md)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Contact the development team
- Check the documentation

## 🔗 Links

- [Digiflazz Official Website](https://digiflazz.com)
- [Digiflazz API Documentation](https://developer.digiflazz.com/api/)
- [Go Documentation](https://golang.org/doc/)

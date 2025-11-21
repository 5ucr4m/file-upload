# File Upload System

A modern, full-stack file upload system with version control, built with Go and React.

## Features

✨ **Core Features:**
- 📤 Chunked file upload with pause/resume/cancel
- 📊 Real-time upload progress bar
- 🔄 Automatic file versioning
- 🏷️ Tag-based organization
- 🔍 Advanced filtering and search
- 🔗 Shareable links with optional password protection
- 📦 Version history for each file
- 🎨 Beautiful dark theme UI with glassmorphism effects

## Tech Stack

**Backend:**
- Go 1.21+
- SQLite database
- Gorilla Mux router
- bcrypt for password hashing

**Frontend:**
- TypeScript
- React 18
- Vite
- Tailwind CSS v4
- Shadcn UI components
- Lucide icons
- Axios for API calls

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Node.js 18+ and npm/yarn
- Git

### Installation

1. **Clone the repository:**
```bash
git clone <repository-url>
cd file-upload
```

2. **Setup Backend:**
```bash
cd backend

# Install dependencies
go mod download

# Run the server
go run main.go
```

The backend will start on `http://localhost:8080`

3. **Setup Frontend (in a new terminal):**
```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

The frontend will start on `http://localhost:5173`

### Building for Production

**Backend:**
```bash
cd backend
go build -o fileupload main.go
./fileupload
```

**Frontend:**
```bash
cd frontend
npm run build
# The built files will be in the dist/ directory
```

## Usage

### Uploading Files

1. Navigate to the Upload tab
2. Drag & drop files or click "Browse File" to select
3. (Optional) Add version number, description, and tags
4. Click "Start Upload"
5. Use Pause/Resume/Cancel buttons to control the upload

### Managing Files

1. Switch to the Files tab to view all uploaded files
2. Use the search bar to find specific files
3. Filter by tags
4. Click on a file to:
   - Download
   - Share with a link
   - View version history

### Sharing Files

1. Click the share icon on any file
2. (Optional) Set a password
3. (Optional) Set expiration time in hours
4. Generate and copy the shareable link

### Version Control

- Upload a file with the same name to create a new version
- Manual version strings (e.g., "v1.2.3") or auto-increment
- View all versions in the version history dialog
- Download any previous version

## API Endpoints

### Upload
- `POST /api/upload/init` - Initialize upload session
- `POST /api/upload/chunk` - Upload file chunk
- `POST /api/upload/complete` - Complete upload
- `DELETE /api/upload/cancel` - Cancel upload
- `GET /api/upload/status` - Get upload status

### Files
- `GET /api/files` - List all files
- `GET /api/files/file` - Get file details
- `GET /api/files/download` - Download file
- `GET /api/files/versions` - Get file versions

### Share
- `POST /api/share/create` - Create share link
- `POST /api/share/access` - Access shared file

## Configuration

### Backend

Environment variables:
- `PORT` - Server port (default: 8080)

### Frontend

The frontend proxies API requests to `http://localhost:8080` in development mode.

For production, update the proxy configuration in `vite.config.ts`.

## Database Schema

The SQLite database includes tables for:
- `files` - File metadata and versions
- `file_versions` - Version history
- `share_links` - Shareable links with passwords
- `upload_sessions` - Resumable upload tracking

## Project Structure

```
file-upload/
├── backend/
│   ├── main.go           # Entry point
│   ├── database/         # Database operations
│   ├── handlers/         # HTTP handlers
│   ├── models/           # Data models
│   ├── utils/            # Utilities
│   └── uploads/          # Upload storage
├── frontend/
│   ├── src/
│   │   ├── components/   # React components
│   │   ├── lib/          # Utilities
│   │   ├── styles/       # Global styles
│   │   ├── App.tsx       # Main app
│   │   └── main.tsx      # Entry point
│   ├── index.html
│   └── vite.config.ts
└── README.md
```

## Security Features

- Password hashing with bcrypt
- Optional password protection for shared links
- Expirable share links
- Chunked uploads for large files
- CORS configuration

## Browser Support

- Chrome/Edge (latest)
- Firefox (latest)
- Safari (latest)

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## License

MIT License

## Support

For issues and questions, please open an issue on GitHub.

import { useState } from 'react'
import { UploadZone } from './components/UploadZone'
import { FileList } from './components/FileList'
import { Upload, Files } from 'lucide-react'
import { Button } from './components/ui/button'

function App() {
  const [activeTab, setActiveTab] = useState<'upload' | 'files'>('upload')
  const [refreshKey, setRefreshKey] = useState(0)

  const handleUploadComplete = () => {
    setRefreshKey(prev => prev + 1)
    // Optionally switch to files tab after upload
    // setActiveTab('files')
  }

  return (
    <div className="min-h-screen">
      <div className="container mx-auto px-4 py-8 max-w-7xl">
        {/* Header */}
        <div className="mb-8 text-center">
          <div className="inline-flex items-center gap-3 mb-4">
            <div className="relative">
              <div className="absolute inset-0 bg-primary/30 blur-xl rounded-full" />
              <div className="relative bg-gradient-to-br from-blue-500 to-blue-600 p-4 rounded-xl">
                <Upload className="w-8 h-8 text-white" />
              </div>
            </div>
          </div>
          <h1 className="text-4xl font-bold text-white mb-2">File Upload System</h1>
          <p className="text-muted-foreground">
            Upload, version, and share your files with ease
          </p>
        </div>

        {/* Tab Navigation */}
        <div className="glass-effect rounded-lg p-1 mb-6 inline-flex w-full max-w-md mx-auto">
          <Button
            variant={activeTab === 'upload' ? 'default' : 'ghost'}
            className="flex-1"
            onClick={() => setActiveTab('upload')}
          >
            <Upload className="w-4 h-4 mr-2" />
            Upload
          </Button>
          <Button
            variant={activeTab === 'files' ? 'default' : 'ghost'}
            className="flex-1"
            onClick={() => setActiveTab('files')}
          >
            <Files className="w-4 h-4 mr-2" />
            Files
          </Button>
        </div>

        {/* Content */}
        <div className="mt-8">
          {activeTab === 'upload' ? (
            <UploadZone onUploadComplete={handleUploadComplete} />
          ) : (
            <FileList key={refreshKey} />
          )}
        </div>

        {/* Footer */}
        <div className="mt-12 text-center text-sm text-muted-foreground">
          <p>Built with Go, TypeScript, React, Tailwind CSS, and SQLite</p>
        </div>
      </div>
    </div>
  )
}

export default App

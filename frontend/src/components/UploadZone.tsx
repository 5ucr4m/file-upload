import { useCallback, useState } from 'react'
import { Upload, Pause, Play, X, FileIcon } from 'lucide-react'
import { Button } from './ui/button'
import { Input } from './ui/input'
import { Label } from './ui/label'
import { Progress } from './ui/progress'
import axios from 'axios'
import { cn, formatBytes } from '@/lib/utils'

interface UploadFile {
  file: File
  uploadId?: string
  progress: number
  status: 'pending' | 'uploading' | 'paused' | 'completed' | 'error'
  description: string
  version: string
  tags: string[]
  currentChunk: number
  totalChunks: number
}

const CHUNK_SIZE = 1024 * 1024 // 1MB chunks

export function UploadZone({ onUploadComplete }: { onUploadComplete?: () => void }) {
  const [isDragging, setIsDragging] = useState(false)
  const [uploads, setUploads] = useState<UploadFile[]>([])

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(true)
  }, [])

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(false)
  }, [])

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(false)

    const files = Array.from(e.dataTransfer.files)
    addFiles(files)
  }, [])

  const handleFileSelect = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      const files = Array.from(e.target.files)
      addFiles(files)
    }
  }, [])

  const addFiles = (files: File[]) => {
    const newUploads: UploadFile[] = files.map(file => ({
      file,
      progress: 0,
      status: 'pending',
      description: '',
      version: '',
      tags: [],
      currentChunk: 0,
      totalChunks: Math.ceil(file.size / CHUNK_SIZE),
    }))

    setUploads(prev => [...prev, ...newUploads])
  }

  const startUpload = async (index: number) => {
    const upload = uploads[index]
    if (!upload) return

    try {
      // Initialize upload
      const initResponse = await axios.post('/api/upload/init', {
        file_name: upload.file.name,
        file_size: upload.file.size,
        total_chunks: upload.totalChunks,
      })

      const uploadId = initResponse.data.upload_id

      setUploads(prev => {
        const newUploads = [...prev]
        newUploads[index] = {
          ...newUploads[index],
          uploadId,
          status: 'uploading',
        }
        return newUploads
      })

      // Upload chunks
      await uploadChunks(index, uploadId)
    } catch (error) {
      console.error('Upload error:', error)
      setUploads(prev => {
        const newUploads = [...prev]
        newUploads[index] = {
          ...newUploads[index],
          status: 'error',
        }
        return newUploads
      })
    }
  }

  const uploadChunks = async (index: number, uploadId: string) => {
    const upload = uploads[index]
    if (!upload) return

    for (let i = upload.currentChunk; i < upload.totalChunks; i++) {
      // Check if paused
      if (uploads[index]?.status === 'paused') {
        return
      }

      const start = i * CHUNK_SIZE
      const end = Math.min(start + CHUNK_SIZE, upload.file.size)
      const chunk = upload.file.slice(start, end)

      try {
        await axios.post(
          `/api/upload/chunk?upload_id=${uploadId}&chunk_index=${i}`,
          chunk,
          {
            headers: {
              'Content-Type': 'application/octet-stream',
            },
          }
        )

        const progress = ((i + 1) / upload.totalChunks) * 100

        setUploads(prev => {
          const newUploads = [...prev]
          newUploads[index] = {
            ...newUploads[index],
            currentChunk: i + 1,
            progress,
          }
          return newUploads
        })
      } catch (error) {
        console.error('Chunk upload error:', error)
        setUploads(prev => {
          const newUploads = [...prev]
          newUploads[index] = {
            ...newUploads[index],
            status: 'error',
          }
          return newUploads
        })
        return
      }
    }

    // Complete upload
    await completeUpload(index, uploadId)
  }

  const completeUpload = async (index: number, uploadId: string) => {
    const upload = uploads[index]
    if (!upload) return

    try {
      await axios.post('/api/upload/complete', {
        upload_id: uploadId,
        description: upload.description,
        version: upload.version,
        tags: upload.tags,
      })

      setUploads(prev => {
        const newUploads = [...prev]
        newUploads[index] = {
          ...newUploads[index],
          status: 'completed',
          progress: 100,
        }
        return newUploads
      })

      if (onUploadComplete) {
        onUploadComplete()
      }
    } catch (error) {
      console.error('Complete upload error:', error)
      setUploads(prev => {
        const newUploads = [...prev]
        newUploads[index] = {
          ...newUploads[index],
          status: 'error',
        }
        return newUploads
      })
    }
  }

  const pauseUpload = (index: number) => {
    setUploads(prev => {
      const newUploads = [...prev]
      newUploads[index] = {
        ...newUploads[index],
        status: 'paused',
      }
      return newUploads
    })
  }

  const resumeUpload = async (index: number) => {
    const upload = uploads[index]
    if (!upload || !upload.uploadId) return

    setUploads(prev => {
      const newUploads = [...prev]
      newUploads[index] = {
        ...newUploads[index],
        status: 'uploading',
      }
      return newUploads
    })

    await uploadChunks(index, upload.uploadId)
  }

  const cancelUpload = async (index: number) => {
    const upload = uploads[index]

    if (upload.uploadId) {
      try {
        await axios.delete(`/api/upload/cancel?upload_id=${upload.uploadId}`)
      } catch (error) {
        console.error('Cancel upload error:', error)
      }
    }

    setUploads(prev => prev.filter((_, i) => i !== index))
  }

  const updateField = (index: number, field: keyof UploadFile, value: any) => {
    setUploads(prev => {
      const newUploads = [...prev]
      newUploads[index] = {
        ...newUploads[index],
        [field]: value,
      }
      return newUploads
    })
  }

  return (
    <div className="space-y-6">
      {/* Upload Zone */}
      <div
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        className={cn(
          'glass-effect rounded-2xl p-8 transition-all duration-300',
          isDragging && 'glow-effect-strong border-primary'
        )}
      >
        <div className="flex flex-col items-center space-y-4">
          <div className="relative">
            <div className="absolute inset-0 bg-primary/20 blur-xl rounded-full" />
            <div className="relative bg-gradient-to-br from-blue-500 to-blue-600 p-6 rounded-2xl glow-effect">
              <Upload className="w-12 h-12 text-white" />
            </div>
          </div>

          <div className="text-center space-y-2">
            <h2 className="text-2xl font-bold text-white">Upload files</h2>
            <p className="text-muted-foreground">
              Select and upload the files of your choice
            </p>
          </div>

          <div className="w-full max-w-md space-y-4">
            <div className="text-center p-6 border-2 border-dashed border-muted rounded-lg">
              <FileIcon className="w-8 h-8 mx-auto mb-2 text-muted-foreground" />
              <p className="text-sm text-muted-foreground mb-1">
                Choose a file or drag & drop it here
              </p>
              <p className="text-xs text-muted-foreground">
                JPEG, PNG, PDF, and MP4 formats, up to 50MB
              </p>
            </div>

            <label htmlFor="file-upload">
              <Button className="w-full" size="lg" asChild>
                <span>
                  Browse File
                  <input
                    id="file-upload"
                    type="file"
                    multiple
                    className="hidden"
                    onChange={handleFileSelect}
                  />
                </span>
              </Button>
            </label>
          </div>
        </div>
      </div>

      {/* Upload List */}
      {uploads.length > 0 && (
        <div className="space-y-4">
          {uploads.map((upload, index) => (
            <div key={index} className="glass-effect rounded-lg p-4 space-y-3">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <p className="font-medium text-white">{upload.file.name}</p>
                  <p className="text-sm text-muted-foreground">
                    {formatBytes(upload.file.size)}
                  </p>
                </div>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => cancelUpload(index)}
                  disabled={upload.status === 'completed'}
                >
                  <X className="w-4 h-4" />
                </Button>
              </div>

              {upload.status !== 'completed' && (
                <div className="space-y-2">
                  <div className="grid grid-cols-2 gap-2">
                    <div>
                      <Label htmlFor={`version-${index}`} className="text-xs">
                        Version (optional)
                      </Label>
                      <Input
                        id={`version-${index}`}
                        placeholder="v1.0.0"
                        value={upload.version}
                        onChange={(e) => updateField(index, 'version', e.target.value)}
                        className="h-8 text-sm"
                      />
                    </div>
                    <div>
                      <Label htmlFor={`tags-${index}`} className="text-xs">
                        Tags (comma separated)
                      </Label>
                      <Input
                        id={`tags-${index}`}
                        placeholder="apk, release"
                        value={upload.tags.join(', ')}
                        onChange={(e) =>
                          updateField(
                            index,
                            'tags',
                            e.target.value.split(',').map(t => t.trim()).filter(Boolean)
                          )
                        }
                        className="h-8 text-sm"
                      />
                    </div>
                  </div>
                  <div>
                    <Label htmlFor={`desc-${index}`} className="text-xs">
                      Description (optional)
                    </Label>
                    <Input
                      id={`desc-${index}`}
                      placeholder="What's new in this version..."
                      value={upload.description}
                      onChange={(e) => updateField(index, 'description', e.target.value)}
                      className="h-8 text-sm"
                    />
                  </div>
                </div>
              )}

              {upload.status !== 'pending' && (
                <div className="space-y-2">
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-muted-foreground">
                      {upload.status === 'completed'
                        ? 'Completed'
                        : upload.status === 'error'
                        ? 'Error'
                        : `${Math.round(upload.progress)}%`}
                    </span>
                    <span className="text-muted-foreground">
                      {upload.currentChunk} / {upload.totalChunks} chunks
                    </span>
                  </div>
                  <Progress value={upload.progress} />
                </div>
              )}

              <div className="flex gap-2">
                {upload.status === 'pending' && (
                  <Button onClick={() => startUpload(index)} className="flex-1">
                    Start Upload
                  </Button>
                )}
                {upload.status === 'uploading' && (
                  <Button onClick={() => pauseUpload(index)} variant="secondary" className="flex-1">
                    <Pause className="w-4 h-4 mr-2" />
                    Pause
                  </Button>
                )}
                {upload.status === 'paused' && (
                  <Button onClick={() => resumeUpload(index)} className="flex-1">
                    <Play className="w-4 h-4 mr-2" />
                    Resume
                  </Button>
                )}
                {upload.status === 'completed' && (
                  <div className="flex-1 text-center text-sm text-green-400 font-medium">
                    ✓ Upload completed successfully
                  </div>
                )}
                {upload.status === 'error' && (
                  <div className="flex-1 text-center text-sm text-destructive font-medium">
                    ✗ Upload failed
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

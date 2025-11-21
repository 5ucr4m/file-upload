import { useState, useEffect } from 'react'
import { FileItem } from './FileItem'
import { FilterBar } from './FilterBar'
import axios from 'axios'

export interface FileData {
  id: string
  name: string
  original_name: string
  size: number
  mime_type: string
  description: string
  version: string
  version_number: number
  tags: string[]
  created_at: string
  updated_at: string
}

export function FileList() {
  const [files, setFiles] = useState<FileData[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedTag, setSelectedTag] = useState<string>('')
  const [searchQuery, setSearchQuery] = useState<string>('')

  const fetchFiles = async () => {
    try {
      setLoading(true)
      const params: any = { limit: 100 }
      if (selectedTag) {
        params.tag = selectedTag
      }
      const response = await axios.get('/api/files', { params })
      setFiles(response.data || [])
    } catch (error) {
      console.error('Error fetching files:', error)
      setFiles([])
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchFiles()
  }, [selectedTag])

  const filteredFiles = files.filter(file => {
    if (!searchQuery) return true
    const query = searchQuery.toLowerCase()
    return (
      file.name.toLowerCase().includes(query) ||
      file.description?.toLowerCase().includes(query) ||
      file.tags?.some(tag => tag.toLowerCase().includes(query))
    )
  })

  const allTags = Array.from(
    new Set(files.flatMap(file => file.tags || []))
  )

  return (
    <div className="space-y-6">
      <FilterBar
        searchQuery={searchQuery}
        onSearchChange={setSearchQuery}
        selectedTag={selectedTag}
        onTagChange={setSelectedTag}
        tags={allTags}
        onRefresh={fetchFiles}
      />

      {loading ? (
        <div className="text-center py-12">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
          <p className="mt-4 text-muted-foreground">Loading files...</p>
        </div>
      ) : filteredFiles.length === 0 ? (
        <div className="text-center py-12 glass-effect rounded-lg">
          <p className="text-muted-foreground">No files found</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {filteredFiles.map(file => (
            <FileItem key={file.id} file={file} onUpdate={fetchFiles} />
          ))}
        </div>
      )}
    </div>
  )
}

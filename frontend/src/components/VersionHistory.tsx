import { useState, useEffect } from 'react'
import { Download, Package, Clock } from 'lucide-react'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from './ui/dialog'
import { Button } from './ui/button'
import { Badge } from './ui/badge'
import { formatBytes, formatDate } from '@/lib/utils'
import axios from 'axios'
import { FileData } from './FileList'

interface VersionHistoryProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  fileName: string
}

export function VersionHistory({ open, onOpenChange, fileName }: VersionHistoryProps) {
  const [versions, setVersions] = useState<FileData[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (open) {
      fetchVersions()
    }
  }, [open, fileName])

  const fetchVersions = async () => {
    try {
      setLoading(true)
      const response = await axios.get(`/api/files/versions?name=${encodeURIComponent(fileName)}`)
      setVersions(response.data || [])
    } catch (error) {
      console.error('Error fetching versions:', error)
      setVersions([])
    } finally {
      setLoading(false)
    }
  }

  const handleDownload = (fileId: string) => {
    window.open(`/api/files/download?id=${fileId}`, '_blank')
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle>Version History</DialogTitle>
          <DialogDescription>{fileName}</DialogDescription>
        </DialogHeader>

        <div className="space-y-3 max-h-[500px] overflow-y-auto">
          {loading ? (
            <div className="text-center py-8">
              <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            </div>
          ) : versions.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              No versions found
            </div>
          ) : (
            versions.map((version, index) => (
              <div
                key={version.id}
                className="glass-effect rounded-lg p-4 space-y-2"
              >
                <div className="flex items-start justify-between">
                  <div className="space-y-1">
                    <div className="flex items-center gap-2">
                      <Badge variant={index === 0 ? 'default' : 'secondary'}>
                        <Package className="w-3 h-3 mr-1" />
                        {version.version || `v${version.version_number}`}
                      </Badge>
                      {index === 0 && (
                        <Badge variant="default" className="text-xs">
                          Latest
                        </Badge>
                      )}
                    </div>
                    <p className="text-sm text-muted-foreground">
                      {formatBytes(version.size)}
                    </p>
                  </div>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handleDownload(version.id)}
                  >
                    <Download className="w-4 h-4 mr-2" />
                    Download
                  </Button>
                </div>

                {version.description && (
                  <p className="text-sm text-foreground">{version.description}</p>
                )}

                <div className="flex items-center gap-4 text-xs text-muted-foreground pt-2 border-t border-border">
                  <div className="flex items-center gap-1">
                    <Clock className="w-3 h-3" />
                    <span>{formatDate(version.created_at)}</span>
                  </div>
                  {version.tags && version.tags.length > 0 && (
                    <div className="flex gap-1">
                      {version.tags.map(tag => (
                        <Badge key={tag} variant="outline" className="text-xs">
                          {tag}
                        </Badge>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            ))
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}

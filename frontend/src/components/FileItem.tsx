import { useState } from 'react'
import { Download, Share2, Clock, FileIcon, Package } from 'lucide-react'
import { Button } from './ui/button'
import { Badge } from './ui/badge'
import { ShareDialog } from './ShareDialog'
import { VersionHistory } from './VersionHistory'
import { formatBytes, formatDate } from '@/lib/utils'
import { FileData } from './FileList'

interface FileItemProps {
  file: FileData
  onUpdate: () => void
}

export function FileItem({ file, onUpdate }: FileItemProps) {
  const [showShare, setShowShare] = useState(false)
  const [showVersions, setShowVersions] = useState(false)

  const handleDownload = () => {
    window.open(`/api/files/download?id=${file.id}`, '_blank')
  }

  const getFileIcon = () => {
    if (file.mime_type.startsWith('image/')) {
      return '🖼️'
    } else if (file.mime_type.includes('pdf')) {
      return '📄'
    } else if (file.mime_type.includes('video')) {
      return '🎥'
    } else if (file.mime_type.includes('android')) {
      return '📱'
    }
    return '📎'
  }

  return (
    <>
      <div className="glass-effect rounded-lg p-4 space-y-3 hover:glow-effect transition-all">
        <div className="flex items-start justify-between">
          <div className="flex items-start gap-3 flex-1 min-w-0">
            <div className="text-3xl">{getFileIcon()}</div>
            <div className="flex-1 min-w-0">
              <h3 className="font-semibold text-white truncate">{file.original_name}</h3>
              <p className="text-sm text-muted-foreground">{formatBytes(file.size)}</p>
            </div>
          </div>
          <div className="flex gap-1">
            <Button variant="ghost" size="icon" onClick={handleDownload}>
              <Download className="w-4 h-4" />
            </Button>
            <Button variant="ghost" size="icon" onClick={() => setShowShare(true)}>
              <Share2 className="w-4 h-4" />
            </Button>
          </div>
        </div>

        {file.description && (
          <p className="text-sm text-muted-foreground line-clamp-2">{file.description}</p>
        )}

        <div className="flex flex-wrap gap-2">
          {file.version && (
            <Badge variant="secondary" className="text-xs">
              <Package className="w-3 h-3 mr-1" />
              {file.version}
            </Badge>
          )}
          {file.tags?.map(tag => (
            <Badge key={tag} variant="outline" className="text-xs">
              {tag}
            </Badge>
          ))}
        </div>

        <div className="flex items-center justify-between text-xs text-muted-foreground pt-2 border-t border-border">
          <div className="flex items-center gap-1">
            <Clock className="w-3 h-3" />
            <span>{formatDate(file.created_at)}</span>
          </div>
          <Button
            variant="ghost"
            size="sm"
            className="h-6 text-xs"
            onClick={() => setShowVersions(true)}
          >
            View versions
          </Button>
        </div>
      </div>

      <ShareDialog
        open={showShare}
        onOpenChange={setShowShare}
        fileId={file.id}
        fileName={file.original_name}
      />

      <VersionHistory
        open={showVersions}
        onOpenChange={setShowVersions}
        fileName={file.original_name}
      />
    </>
  )
}

import { useState } from 'react'
import { Copy, Check, Lock } from 'lucide-react'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from './ui/dialog'
import { Button } from './ui/button'
import { Input } from './ui/input'
import { Label } from './ui/label'
import axios from 'axios'

interface ShareDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  fileId: string
  fileName: string
}

export function ShareDialog({ open, onOpenChange, fileId, fileName }: ShareDialogProps) {
  const [shareLink, setShareLink] = useState<string>('')
  const [password, setPassword] = useState<string>('')
  const [expiresIn, setExpiresIn] = useState<string>('')
  const [loading, setLoading] = useState(false)
  const [copied, setCopied] = useState(false)

  const generateLink = async () => {
    try {
      setLoading(true)
      const response = await axios.post('/api/share/create', {
        file_id: fileId,
        password: password || null,
        expires_in: expiresIn ? parseInt(expiresIn) : null,
      })

      const fullUrl = `${window.location.origin}/s/${response.data.short_code}`
      setShareLink(fullUrl)
    } catch (error) {
      console.error('Error creating share link:', error)
      alert('Failed to create share link')
    } finally {
      setLoading(false)
    }
  }

  const copyToClipboard = () => {
    navigator.clipboard.writeText(shareLink)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Share File</DialogTitle>
          <DialogDescription>
            Create a shareable link for {fileName}
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          {!shareLink ? (
            <>
              <div className="space-y-2">
                <Label htmlFor="password">
                  Password (optional)
                  <Lock className="w-3 h-3 inline ml-1" />
                </Label>
                <Input
                  id="password"
                  type="password"
                  placeholder="Leave empty for no password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="expires">Expires in (hours, optional)</Label>
                <Input
                  id="expires"
                  type="number"
                  placeholder="Leave empty for no expiration"
                  value={expiresIn}
                  onChange={(e) => setExpiresIn(e.target.value)}
                />
              </div>

              <Button onClick={generateLink} disabled={loading} className="w-full">
                {loading ? 'Generating...' : 'Generate Share Link'}
              </Button>
            </>
          ) : (
            <>
              <div className="space-y-2">
                <Label>Share Link</Label>
                <div className="flex gap-2">
                  <Input value={shareLink} readOnly className="flex-1" />
                  <Button
                    variant="outline"
                    size="icon"
                    onClick={copyToClipboard}
                  >
                    {copied ? (
                      <Check className="w-4 h-4 text-green-500" />
                    ) : (
                      <Copy className="w-4 h-4" />
                    )}
                  </Button>
                </div>
              </div>

              {password && (
                <div className="p-3 bg-secondary rounded-md text-sm">
                  <p className="text-muted-foreground">
                    This link is password protected. Share the password separately.
                  </p>
                </div>
              )}

              {expiresIn && (
                <div className="p-3 bg-secondary rounded-md text-sm">
                  <p className="text-muted-foreground">
                    This link will expire in {expiresIn} hours.
                  </p>
                </div>
              )}

              <Button
                variant="outline"
                onClick={() => {
                  setShareLink('')
                  setPassword('')
                  setExpiresIn('')
                }}
                className="w-full"
              >
                Create Another Link
              </Button>
            </>
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}

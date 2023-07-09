import { useMemo, useState } from 'react'
import { Center, Stack } from '@chakra-ui/react'
import { SectionSpinner, variables } from '@koupr/ui'
import { File } from '@/client/api/file'
import { getAccessTokenOrRedirect } from '@/infra/token'

type ImageViewerProps = {
  file: File
}

const ImageViewer = ({ file }: ImageViewerProps) => {
  const [isLoading, setIsLoading] = useState(true)
  const download = useMemo(() => file.preview ?? file.original, [file])
  const path = useMemo(() => (file.preview ? 'preview' : 'original'), [file])
  const url = useMemo(() => {
    if (!download?.extension) {
      return ''
    }
    return `/proxy/api/v1/files/${file.id}/${path}${
      download.extension
    }?${new URLSearchParams({
      access_token: getAccessTokenOrRedirect(),
    })}`
  }, [file, download, path])

  if (!download) {
    return null
  }

  return (
    <Stack direction="column" w="100%" h="100%" spacing={variables.spacing}>
      <Center
        flexGrow={1}
        w="100%"
        h="100%"
        overflow="scroll"
        position="relative"
      >
        {isLoading && <SectionSpinner />}
        <img
          src={url}
          style={{
            objectFit: 'contain',
            width: isLoading ? 0 : 'auto',
            height: isLoading ? 0 : '100%',
            visibility: isLoading ? 'hidden' : 'visible',
          }}
          onLoad={() => setIsLoading(false)}
          alt={file.name}
        />
      </Center>
    </Stack>
  )
}

export default ImageViewer

import { Badge, Stat, StatLabel, StatNumber } from '@chakra-ui/react'
import cx from 'classnames'
import { File } from '@/client/api/file'

export type FileInfoExtensionProps = {
  file: File
}

const FileInfoExtension = ({ file }: FileInfoExtensionProps) => {
  if (!file.original) {
    return null
  }
  return (
    <Stat>
      <StatLabel>File type</StatLabel>
      <StatNumber className={cx('text-base')}>
        <Badge>{file.original.extension}</Badge>
      </StatNumber>
    </Stat>
  )
}

export default FileInfoExtension

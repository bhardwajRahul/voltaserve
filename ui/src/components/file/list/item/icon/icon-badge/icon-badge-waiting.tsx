// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { Circle, Tooltip } from '@chakra-ui/react'
import cx from 'classnames'
import { IconHourglass } from '@/lib/components/icons'

const IconBadgeWaiting = () => (
  <Tooltip label="Waiting for processing">
    <Circle
      className={cx(
        'text-orange-600',
        'bg-white',
        'w-[23px]',
        'h-[23px]',
        'border',
        'border-gray-200',
      )}
    >
      <IconHourglass />
    </Circle>
  </Tooltip>
)

export default IconBadgeWaiting

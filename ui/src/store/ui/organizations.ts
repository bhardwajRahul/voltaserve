// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { KeyedMutator } from 'swr'
import { PayloadAction, createSlice } from '@reduxjs/toolkit'
import { List } from '@/client/api/organization'

type OrganizationsState = {
  isInviteModalOpen: boolean
  mutate?: KeyedMutator<List>
}

const initialState: OrganizationsState = {
  isInviteModalOpen: false,
}

const slice = createSlice({
  name: 'organizations',
  initialState,
  reducers: {
    inviteModalDidOpen: (state) => {
      state.isInviteModalOpen = true
    },
    inviteModalDidClose: (state) => {
      state.isInviteModalOpen = false
    },
    mutateUpdated: (state, action: PayloadAction<KeyedMutator<List>>) => {
      state.mutate = action.payload
    },
  },
})

export const { inviteModalDidOpen, inviteModalDidClose, mutateUpdated } =
  slice.actions

export default slice.reducer

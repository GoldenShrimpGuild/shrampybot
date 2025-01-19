import { defineStore } from 'pinia'
import { TwitchUserDatum } from '../../model/utility/nosqldb'
import { User as DiscordUser } from '../../model/lib/discordgo'
import { SelfResponseBody as User } from '../../model/controller/auth'

// export interface User {
//   discord_id?: string
//   username?: string
//   is_authenticated?: boolean
//   is_superuser?: boolean
//   is_staff?: boolean
//   streamer?: TwitchUserDatum
//   discord_user?: DiscordUser
// }

export const useUserStore = defineStore('user', {
  state: () => {
    return {
      self: {
        id: '',
        username: '',
      } as User,
    }
  },
  actions: {
    setSelf(selfData: User) {
      this.$state.self = selfData
    },
  },
})

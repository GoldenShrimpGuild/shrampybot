import { defineStore } from 'pinia'
import { TwitchUserDatum } from '../../model/utility/nosqldb'
import { User as DiscordUser } from '../../model/lib/discordgo'

export interface User {
  username?: string
  is_authenticated?: boolean
  is_superuser?: boolean
  is_staff?: boolean
  streamer?: TwitchUserDatum
  discord_user?: DiscordUser
}

export const useUserStore = defineStore('user', {
  state: () => {
    return {
      self: {
        username: '',
        is_authenticated: false,
        is_superuser: false,
        is_staff: false,
        streamer: {} as TwitchUserDatum,
        discord_user: {} as DiscordUser,
      } as User,
    }
  },
  actions: {
    setSelf(selfData: User) {
      this.$state.self = selfData
    },
  },
})

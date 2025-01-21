import { defineStore } from 'pinia'
import { TwitchUserDatum } from '../../model/utility/nosqldb'
import { User as DiscordUser } from '../../model/lib/discordgo'
import { SelfResponseBody as User } from '../../model/controller/auth'
import { useGlobalStore } from './global-store'
import { useAuthStore } from './auth'
import axios from 'axios'

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
    async fetchSelf() {
      const GlobalStore = useGlobalStore()
      const AuthStore = useAuthStore()

      const self_path = '/auth/self'
      const axiosConfig = AuthStore.getAxiosConfig()

      const bearerResponse = axios.get(
        self_path,
        axiosConfig)
      .then((response) => {
        if (response.status === 200) {
          this.setSelf(response.data)
        }
      })
      .catch((err: any) => {
        if (err.response.status === 401) {
          AuthStore.callRefresh()
        }
      })
      return bearerResponse
    },
    isAdmin() {
      // Admins role
      if (this.$state.self.member?.roles.includes("732364663194648756")) {
        return true
      } else {
        return false
      }
    },
    isDevTeam() {
      // Development Team members role
      if (this.$state.self.member?.roles.includes("978811326589710446")) {
        return true
      } else {
        return false
      }
    },
    isGSGMember() {
      if (this.$state.self.member) {
        return true
      } else {
        return false
      }
    }
  },
})

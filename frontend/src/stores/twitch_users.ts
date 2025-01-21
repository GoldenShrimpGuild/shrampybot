import { defineStore } from 'pinia'
import { TwitchUserDatum } from '../../model/utility/nosqldb'
// import { User as DiscordUser } from '../../model/lib/discordgo'
// import { SelfResponseBody as User } from '../../model/controller/auth'
import { useGlobalStore } from './global-store'
import { useAuthStore } from './auth'
import axios from 'axios'

export const useTwitchUsersStore = defineStore('users', {
  state: () => {
    return {
      users: [] as Array<TwitchUserDatum>,
    }
  },
  actions: {
    async fetchUsers() {
      const AuthStore = useAuthStore()

      const users_path = '/admin/user'
      const axiosConfig = AuthStore.getAxiosConfig()

      const bearerResponse = axios.get(
        users_path,
        axiosConfig)
        .then((response) => {
          if (response.status === 200) {
            this.$state.users = response.data.data
          }
        })
        .catch((err: any) => {
          if (err.response.status === 401) {
            AuthStore.callRefresh()
          }
        })
      return bearerResponse
    },
  },
})

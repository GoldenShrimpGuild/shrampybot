import { defineStore } from 'pinia'
import { TwitchUserDatum } from '../../model/utility/nosqldb'
// import { User as DiscordUser } from '../../model/lib/discordgo'
// import { SelfResponseBody as User } from '../../model/controller/auth'
import { useGlobalStore } from './global-store'
import { useAuthStore } from './auth'
// import axios from 'axios'
import { useAxios } from '../plugins/axios'

export const useTwitchUsersStore = defineStore('users', {
  state: () => {
    return {
      users: [] as Array<TwitchUserDatum>,
    }
  },
  actions: {
    async fetchUsers() {
      const AuthStore = useAuthStore()
      const axios = useAxios()

      const users_path = '/admin/user'

      const bearerResponse = axios.get(
        users_path)
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

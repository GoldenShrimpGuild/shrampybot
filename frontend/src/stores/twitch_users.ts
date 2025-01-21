import { defineStore } from 'pinia'
import { TwitchUserDatum } from '../../model/utility/nosqldb'
// import { User as DiscordUser } from '../../model/lib/discordgo'
// import { SelfResponseBody as User } from '../../model/controller/auth'

export const useUsersStore = defineStore('users', {
  state: () => {
    return {
      users: [] as Array<TwitchUserDatum>,
    }
  },
  actions: {
    // getUsers() {
    //   this.$state.self = 
    // },
  },
})

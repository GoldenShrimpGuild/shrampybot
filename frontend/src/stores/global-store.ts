import { defineStore } from 'pinia'
import { useLocalStorage } from '@vueuse/core'

export const useGlobalStore = defineStore('global', {
  state: () => {
    const isDevEnvironment = useLocalStorage('isDevEnvironment', false)
    return {
      isSidebarMinimized: false,
      isDevEnvironment,
    }
  },
  actions: {
    toggleSidebar() {
      this.isSidebarMinimized = !this.isSidebarMinimized
    },
    setDevEnvironment(isDevEnvironment: boolean) {
      this.isDevEnvironment = isDevEnvironment
    },
    getApiBaseUrl() {
      // TODO: Assemble this more sensibly
      if (this.isDevEnvironment) {
        return 'https://tl72sifq5iu6gkzpqyyp7umsra0wjejp.lambda-url.ca-central-1.on.aws'
      } else {
        return 'https://3okyp4qsdy2xzm5cjmpw5it53u0tdnkv.lambda-url.ca-central-1.on.aws'
      }
    },
    getDiscordOAuthUrl() {
      // TODO: Assemble this more sensibly
      let client_id
      const discord_seg1 = 'https://discord.com/oauth2/authorize?client_id='
      const discord_seg2 = '&response_type=code&redirect_uri='
      const discord_seg3 = '&scope=identify+connections'

      const this_url = encodeURIComponent(window.location.origin + '/shrampybot/auth/validate_oauth')
      if (this.isDevEnvironment) {
        client_id = '1043225123395739780'
      } else {
        client_id = '1042309025506787359'
      }

      return discord_seg1 + client_id + discord_seg2 + this_url + discord_seg3
    },
  },
})

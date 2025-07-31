import { defineStore } from 'pinia'
import { useLocalStorage } from '@vueuse/core'
import { useCookies } from "vue3-cookies"
import { useRoute } from "vue-router"

// Should load these from an env config instead.
export const apiBaseUrlDev = "https://tl72sifq5iu6gkzpqyyp7umsra0wjejp.lambda-url.ca-central-1.on.aws"
export const apiBaseUrlProd = "https://3okyp4qsdy2xzm5cjmpw5it53u0tdnkv.lambda-url.ca-central-1.on.aws"

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
    getStateBlob() {
      const route = useRoute()
      const { cookies } = useCookies()

      const preBlob = {
        'csrftoken': cookies.get('csrftoken'),
        'redirect_path': route.query.redirect_path ? decodeURIComponent(String(route.query.redirect_path)) : '/streams',
      }

      const blob = btoa(JSON.stringify(preBlob))
      return blob
    },
    decodeStateBlob(stateBlob: string) {
      return JSON.parse(atob(stateBlob))
    },
    getApiBaseUrl() {
      // TODO: Assemble this more sensibly
      if (this.isDevEnvironment) {
        return apiBaseUrlDev
      } else {
        return apiBaseUrlProd
      }
    },
    getDiscordOAuthUrl() {
      // TODO: Assemble this more sensibly
      let client_id
      const discord_seg1 = 'https://discord.com/oauth2/authorize?client_id='
      const discord_seg2 = '&response_type=code&redirect_uri='
      const discord_seg3 = '&scope=identify+connections&state='
      const state_blob = this.getStateBlob()

      const this_url = encodeURIComponent(window.location.origin + '/shrampybot/auth/validate_oauth')
      if (this.isDevEnvironment) {
        client_id = '1043225123395739780'
      } else {
        client_id = '1042309025506787359'
      }

      return discord_seg1 + client_id + discord_seg2 + this_url + discord_seg3 + state_blob
    },
  },
})

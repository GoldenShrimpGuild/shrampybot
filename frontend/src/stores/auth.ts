import { defineStore } from 'pinia'
import axios, { AxiosRequestConfig } from 'axios'
import { useLocalStorage } from '@vueuse/core'
import { useGlobalStore } from './global-store'

export const useAuthStore = defineStore('auth', {
  state: () => {
    const accessToken = useLocalStorage('accessToken', '')
    const refreshToken = useLocalStorage('refreshToken', '')
    const userServicesStatus = useLocalStorage('uss', {} as Record<string, any>)
    return { accessToken, refreshToken, userServicesStatus }
  },
  actions: {
    getAxiosConfig() {
      const GlobalStore = useGlobalStore()
      return {
        baseURL: GlobalStore.getApiBaseUrl(),
        headers: {
          Authorization: `Bearer ${this.accessToken}`,
          'Content-Type': 'application/json',
        },
      } as AxiosRequestConfig
    },
    async testAndRefreshToken() {
      const GlobalStore = useGlobalStore()

      const path = '/auth/self'
      const axiosConfig = this.getAxiosConfig()

      try {
        const bearerResponse = await axios.get(
          path,
          axiosConfig,
        )
        console.log(bearerResponse)
      } catch (error: any) {
        if (error.response.status == 401) {
          const GlobalStore = useGlobalStore()

          const refresh_token = this.$state.refreshToken
          const refresh_path = '/auth/refresh'

          try {
            const refreshResponse = await axios.post(
              refresh_path,
              {
                refresh: refresh_token,
              },
              {
                baseURL: GlobalStore.getApiBaseUrl(),
                headers: {
                  'Content-Type': 'application/json',
                },
              },
            )
            this.$state.accessToken = refreshResponse.data.access
          } catch (refreshError: any) {
            this.$state.refreshToken = ''
            this.$state.accessToken = ''
          }
        } else if (error.response.status == 500) {
          this.$state.refreshToken = ''
          this.$state.accessToken = ''
        }
      }
    },
  },
})

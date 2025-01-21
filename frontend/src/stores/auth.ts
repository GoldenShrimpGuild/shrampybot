import { defineStore } from 'pinia'
import axios, { AxiosRequestConfig } from 'axios'
import { useLocalStorage } from '@vueuse/core'
import { useGlobalStore } from './global-store'

export const useAuthStore = defineStore('auth', {
  state: () => {
    const accessToken = ''
    const userServicesStatus = useLocalStorage('uss', {} as Record<string, any>)
    return { accessToken, userServicesStatus }
  },
  actions: {
    getAxiosConfig() {
      const GlobalStore = useGlobalStore()
      return {
        baseURL: GlobalStore.getApiBaseUrl(),
        withCredentials: true,
        headers: {
          Authorization: `Bearer ${this.accessToken}`,
          'Content-Type': 'application/json',
        },
      } as AxiosRequestConfig
    },
    async callLogout() {
      const GlobalStore = useGlobalStore()

      const logout_path = '/auth/logout'
      const axiosConfig = this.getAxiosConfig()

      const bearerResponse = await axios.post(
        logout_path, 
        {}, 
        axiosConfig)
      this.$state.accessToken = ""
      return bearerResponse
    },
    async callRefresh() {
      const GlobalStore = useGlobalStore()

      const refresh_path = '/auth/refresh'

      try {
        const refreshResponse = await axios.post(
          refresh_path,
          {},
          {
            baseURL: GlobalStore.getApiBaseUrl(),
            withCredentials: true,
            headers: {
              'Content-Type': 'application/json',
            },
          },
        )
        this.$state.accessToken = refreshResponse.data.access
        return refreshResponse
      } catch (refreshError: any) {
        this.$state.accessToken = ''
      }
    },
    async testAndRefreshToken() {
      const GlobalStore = useGlobalStore()

      const path = '/auth/touch'
      const axiosConfig = this.getAxiosConfig()

      try {
        const bearerResponse = await axios.get(path, axiosConfig)
      } catch (error: any) {
        if (error.response.status == 401) {
          this.callRefresh()
        } else if (error.response.status == 500) {
          // this.$state.accessToken = ''
        }
      }
    },
  },
})

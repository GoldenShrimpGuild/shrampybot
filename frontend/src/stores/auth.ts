import { defineStore } from 'pinia'
import axios, { AxiosRequestConfig } from 'axios'
import { useLocalStorage } from '@vueuse/core'
import { useGlobalStore } from './global-store'

export const useAuthStore = defineStore('auth', {
  state: () => {
    const userId = useLocalStorage('user_id', '' as String)
    const accessToken = '' as String
    const userServicesStatus = useLocalStorage('uss', {} as Record<string, any>)
    return { accessToken, userServicesStatus, userId }
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
      .then((response) => {
        this.$state.accessToken = ""
        this.$state.userId = ""
      })
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
        .then((response) => {
          this.$state.userId = response.data.user_id
          this.$state.accessToken = response.data.access
        })
        return refreshResponse
      } catch (refreshError: any) {
        this.$state.accessToken = ''
        this.$state.userId = ''
      }
    },
    async testAndRefreshToken() {
      const GlobalStore = useGlobalStore()

      const path = '/auth/touch'
      const axiosConfig = this.getAxiosConfig()

      try {
        const bearerResponse = await axios.get(path, axiosConfig)
          .then((response) => {
            if (response.status === 200) {
              this.$state.userId = response.data.user_id
            }
          })
      } catch (error: any) {
        if (error.response.status === 401) {
          this.callRefresh()
        } else if (error.response.status === 500) {
          // this.$state.accessToken = ''
        }
      }
    },
  },
})

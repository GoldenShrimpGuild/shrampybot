import { defineStore } from 'pinia'
import axios, { AxiosRequestConfig } from 'axios'
import { useLocalStorage } from '@vueuse/core'
import { useGlobalStore } from './global-store'
import { jwtDecode, JwtPayload } from 'jwt-decode'

export interface CustomJwtPayload extends JwtPayload {
  scopes: string
}

export const useAuthStore = defineStore('auth', {
  state: () => {
    const userId = useLocalStorage('user_id', '' as string)
    const accessTokenDev = '' as string
    const accessTokenProd = '' as string
    const userServicesStatus = useLocalStorage('uss', {} as Record<string, any>)
    return { accessTokenDev, accessTokenProd, userServicesStatus, userId }
  },
  actions: {
    setAccessToken(accessToken: string) {
      const GlobalStore = useGlobalStore()

      if (GlobalStore.isDevEnvironment) {
        this.$state.accessTokenDev = accessToken
      } else {
        this.$state.accessTokenProd = accessToken
      }
    },
    getAccessToken() {
      const GlobalStore = useGlobalStore()

      if (GlobalStore.isDevEnvironment) {
        return this.$state.accessTokenDev
      } else {
        return this.$state.accessTokenProd
      }
    },
    getAxiosConfig() {
      const GlobalStore = useGlobalStore()
      return {
        baseURL: GlobalStore.getApiBaseUrl(),
        withCredentials: true,
        headers: {
          Authorization: `Bearer ${this.getAccessToken()}`,
          'Content-Type': 'application/json',
        },
      } as AxiosRequestConfig
    },
    async callLogout() {
      const GlobalStore = useGlobalStore()

      const logout_path = '/auth/logout'
      const axiosConfig = this.getAxiosConfig()

      const bearerResponse = await axios.post(logout_path, {}, axiosConfig).then((response) => {
        this.setAccessToken('')
        this.$state.userId = ''
      })
      return bearerResponse
    },
    async callRefresh() {
      const GlobalStore = useGlobalStore()
      const refresh_path = '/auth/refresh'
      let success = false

      try {
        const refreshResponse = await axios
          .post(
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
            if (response.status === 200) {
              success = true
              this.$state.userId = response.data.user_id
              this.setAccessToken(response.data.access)
            }
          })
      } catch (refreshError: any) {
        this.setAccessToken('')
        this.$state.userId = ''
      }

      return success
    },
    async testAndRefreshToken() {
      const GlobalStore = useGlobalStore()

      const path = '/auth/touch'
      const axiosConfig = this.getAxiosConfig()

      try {
        await axios.get(path, axiosConfig).then((response) => {
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
    decode() {
      const accessToken = this.getAccessToken()

      if (accessToken) {
        return jwtDecode<CustomJwtPayload>(accessToken.toString())
      }
    },
    getScopes() {
      const jwt = this.decode()
      if (jwt) {
        return jwt.scopes.split(' ')
      }
      return []
    },
  },
})

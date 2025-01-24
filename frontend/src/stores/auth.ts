import { defineStore } from 'pinia'
import axios, { AxiosRequestConfig } from 'axios'
import { useLocalStorage } from '@vueuse/core'
import { useGlobalStore } from './global-store'
import { jwtDecode, JwtPayload } from "jwt-decode";

export interface CustomJwtPayload extends JwtPayload {
  scopes: string
}

export const useAuthStore = defineStore('auth', {
  state: () => {
    const userId = useLocalStorage('user_id', '' as String)
    const accessTokenDev = '' as String
    const accessTokenProd = '' as String
    const userServicesStatus = useLocalStorage('uss', {} as Record<string, any>)
    return { accessTokenDev, accessTokenProd, userServicesStatus, userId }
  },
  actions: {
    getAxiosConfig() {
      const GlobalStore = useGlobalStore()
      return {
        baseURL: GlobalStore.getApiBaseUrl(),
        withCredentials: true,
        headers: {
          Authorization: `Bearer ${GlobalStore.$state.isDevEnvironment ? this.$state.accessTokenDev : this.$state.accessTokenProd }`,
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
        if (GlobalStore.$state.isDevEnvironment) {
          this.$state.accessTokenDev = ""
        } else {
          this.$state.accessTokenProd = ""
        }
        this.$state.userId = ""
      })
      return bearerResponse
    },
    async callRefresh() {
      const GlobalStore = useGlobalStore()
      const refresh_path = '/auth/refresh'
      var success = false

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
          if (response.status === 200) {
            success = true
            this.$state.userId = response.data.user_id
            if (GlobalStore.$state.isDevEnvironment) {
              this.$state.accessTokenDev = response.data.access
            } else {
              this.$state.accessTokenProd = response.data.access
            }
          }
        })
      } catch (refreshError: any) {
        if (GlobalStore.$state.isDevEnvironment) {
          this.$state.accessTokenDev = ""
        } else {
          this.$state.accessTokenProd = ""
        }
        this.$state.userId = ''
      }

      return success
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
    decode() {
      const GlobalStore = useGlobalStore()

      if (GlobalStore.$state.isDevEnvironment) {
        if (this.$state.accessTokenDev) {
          return jwtDecode<CustomJwtPayload>(this.$state.accessTokenDev.toString())
        }
      } else {
        if (this.$state.accessTokenProd) {
          return jwtDecode<CustomJwtPayload>(this.$state.accessTokenProd.toString())
        }
      }
    },
    getScopes() {
      const jwt = this.decode()
      if (jwt) {
        return jwt.scopes.split(' ')
      }
      return []
    }
  },
})

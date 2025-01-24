import { defineStore } from 'pinia'
import { SelfResponseBody as User } from '../../model/controller/auth'
import { useAuthStore } from './auth'
import { useAxios } from '../plugins/axios'
import { forEach } from 'lodash'

export const useUserStore = defineStore('user', {
  state: () => {
    return {
      self: {
        id: '',
        username: '',
      } as User,
    }
  },
  actions: {
    // setSelf(selfData: User) {
    //   this.$state.self = selfData
    // },
    // async fetchSelf() {
    //   const AuthStore = useAuthStore()

    //   const self_path = '/auth/self'
    //   const axiosConfig = AuthStore.getAxiosConfig()

    //   const bearerResponse = axios.get(
    //     self_path,
    //     axiosConfig)
    //   .then((response) => {
    //     if (response.status === 200) {
    //       this.setSelf(response.data)
    //     }
    //   })
    //   .catch((err: any) => {
    //     if (err.response.status === 401) {
    //       AuthStore.callRefresh()
    //     }
    //   })
    //   return bearerResponse
    // },
    isAdmin() {
      const AuthStore = useAuthStore()
      var result = false
      
      forEach(AuthStore.getScopes(), (scope) => {
        const subScope = scope.split(':')
        if (subScope[0] === 'admin') {
          result = true
          return
        }
      })

      return result
    },
    isDevTeam() {
      const AuthStore = useAuthStore()
      var result = false
      
      forEach(AuthStore.getScopes(), (scope) => {
        const subScope = scope.split(':')
        if (subScope[0] === 'dev') {
          result = true
          return
        }
      })

      return result
    }
  },
})

import { defineStore } from 'pinia'
import { OutputStaticTokenInfo, NewTokenRequestBody, NewTokenResponseBody } from '../../model/controller/admin'
import { useAuthStore } from './auth'
import { useAxios } from '../plugins/axios'
import { forEach } from 'lodash'

export const useStaticTokenStore = defineStore('staticToken', {
    state: () => {
        return {
            tokens: [] as Array<OutputStaticTokenInfo>,
            newToken: {} as NewTokenResponseBody,
        }
    },
    actions: {
        async fetchTokenInfo() {
            const AuthStore = useAuthStore()
            const category_path = '/admin/token'
            const axios = useAxios()

            const bearerResponse = axios.get(
                category_path)
                .then((response) => {
                    if (response.status === 200) {
                        this.$state.tokens = response.data.data
                    }
                })
            return bearerResponse
        },
        clearNewToken() {
            this.newToken.id = ''
            this.newToken.purpose = ''
            this.newToken.created_at = ''
            this.newToken.creator_id = ''
            this.newToken.expires_at = ''
            this.newToken.revoked = false
            this.newToken.scopes = ''
            this.newToken.token = ''
        },
        async addToken(tokenRequest: NewTokenRequestBody) {
            const AuthStore = useAuthStore()

            const category_path = '/admin/token'
            const axios = useAxios()

            await axios.post(
                category_path,
                tokenRequest)
                .then((response) => {
                    if (response.status === 200) {
                        this.$state.newToken = response.data
                        const newTokenSanitized = JSON.parse(JSON.stringify(response.data))
                        newTokenSanitized.token = ''
                        console.log(newTokenSanitized)
                        this.$state.tokens.push(newTokenSanitized)
                    }
                })

            return this.$state.newToken
        },
        async revokeToken(id: string | undefined) {
            const AuthStore = useAuthStore()

            const category_path = '/admin/token/' + id
            const axios = useAxios()

            const bearerResponse = axios.delete(
                category_path)
                .then((response) => {
                    if (response.status === 200) {
                        forEach(this.$state.tokens, (item) => {
                            if (item.id === id) {
                                item.revoked = true
                                return
                            }
                        })
                    }
                })
            return bearerResponse as unknown
        },
    },
})

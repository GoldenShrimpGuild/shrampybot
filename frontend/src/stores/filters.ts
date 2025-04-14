import { defineStore } from 'pinia'
import { FilterDatum } from '../../model/utility/nosqldb'
// import { User as DiscordUser } from '../../model/lib/discordgo'
// import { SelfResponseBody as User } from '../../model/controller/auth'
import { useAuthStore } from './auth'
import { useAxios } from '../plugins/axios'
import { forEach } from 'lodash'

export const useFilterStore = defineStore('filters', {
  state: () => {
    return {
      filters: [] as Array<FilterDatum>,
    }
  },
  actions: {
    async fetchFilters() {
      const AuthStore = useAuthStore()
      const path = '/admin/filter'
      const axios = useAxios()

      const bearerResponse = axios.get(path).then((response) => {
        if (response.status === 200) {
          this.$state.filters = response.data.data
        }
      })
      return bearerResponse
    },
    async putFilter(filt: FilterDatum) {
      const AuthStore = useAuthStore()

      const path = '/admin/filter'
      const axios = useAxios()

      const bearerResponse = axios.put(path, filt).then((response) => {
        if (response.status === 200) {
          let foundItem = false

          forEach(this.$state.filters, (item) => {
            if (response.data && response.data.data && response.data.data.length > 0) {
              if (item.id === response.data.data[0].id) {
                item.keyword = response.data.data[0].keyword
                item.is_regex = response.data.data[0].is_regex
                item.case_insensitive = response.data.data[0].case_insensitive
                foundItem = true
                return
              }
            }
          })

          if (!foundItem) {
            this.$state.filters.push(response.data.data[0])
          }
        }
      })
      return bearerResponse as unknown
    },
    async deleteFilter(id: string) {
      const AuthStore = useAuthStore()

      const path = '/admin/filter/' + id
      const axios = useAxios()

      const bearerResponse = axios.delete(path).then((response) => {
        if (response.status === 200) {
          const copyList = this.$state.filters
          this.$state.filters = []

          forEach(copyList, (item) => {
            if (item.id !== id) {
              this.$state.filters.push(item)
            }
          })
        }
      })
      return bearerResponse as unknown
    },
  },
})

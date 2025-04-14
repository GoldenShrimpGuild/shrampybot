import { defineStore } from 'pinia'
import { CategoryDatum } from '../../model/utility/nosqldb'
// import { User as DiscordUser } from '../../model/lib/discordgo'
// import { SelfResponseBody as User } from '../../model/controller/auth'
import { useAuthStore } from './auth'
import { useAxios } from '../plugins/axios'
import { forEach } from 'lodash'

export const useCategoryStore = defineStore('categories', {
  state: () => {
    return {
      categories: [] as Array<CategoryDatum>,
    }
  },
  actions: {
    async fetchCategories() {
      const AuthStore = useAuthStore()
      const category_path = '/admin/category'
      const axios = useAxios()

      const bearerResponse = axios.get(category_path).then((response) => {
        if (response.status === 200) {
          this.$state.categories = response.data.data
        }
      })
      return bearerResponse
    },
    async putCategory(cat: CategoryDatum) {
      const AuthStore = useAuthStore()

      const category_path = '/admin/category'
      const axios = useAxios()

      const bearerResponse = axios.put(category_path, cat).then((response) => {
        if (response.status === 200) {
          let foundItem = false

          forEach(this.$state.categories, (item) => {
            if (response.data && response.data.data && response.data.data.length > 0) {
              if (item.twitch_category === response.data.data[0].twitch_category) {
                item.id = response.data.data[0].id
                item.mastodon_tags = response.data.data[0].mastodon_tags
                item.bluesky_tags = response.data.data[0].bluesky_tags
                foundItem = true
                return
              }
            }
          })

          if (!foundItem) {
            this.$state.categories.push(response.data.data[0])
          }
        }
      })
      return bearerResponse as unknown
    },
    async deleteCategory(id: string) {
      const AuthStore = useAuthStore()

      const category_path = '/admin/category/' + id
      const axios = useAxios()

      const bearerResponse = axios.delete(category_path).then((response) => {
        if (response.status === 200) {
          const copyList = this.$state.categories
          this.$state.categories = []

          forEach(copyList, (item) => {
            if (item.id !== id) {
              this.$state.categories.push(item)
            }
          })
        }
      })
      return bearerResponse as unknown
    },
  },
})

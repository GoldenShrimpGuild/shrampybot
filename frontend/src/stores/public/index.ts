import { defineStore } from 'pinia'
import { useLocalStorage } from '@vueuse/core'
import type { StreamHistoryDatum } from '../../../model/utility/nosqldb'
import { apiBaseUrlDev, apiBaseUrlProd } from '../global-store'
import { useMultiStore } from './multi'
import axios from 'axios'
import type { AxiosInstance } from 'axios'

// Data for axios
const staticContentType = "application/json"
const staticPublicStreamEndpoint = "/public/stream"

export type Streams = { 
  [key: string]: StreamHistoryDatum,
}

// Data for sorting streams out
export const gsgChannelLogin = "goldenshrimpguild";
export const loginMatchesGSG = (login: String) => login == gsgChannelLogin;

// This is a store for use with public pages such as gsgMultiTwitch
export const usePublicStore = defineStore('public', {
  state: () => {
    // Specific child store just for multiTwitch stuff
    // testing this pattern out
    const multi = useMultiStore()
    const streams: Streams = {}

    // More global options for streams loader
    const useDevApi = useLocalStorage('publicUseDevApi', false)
    const includeGSGChannel = useLocalStorage('publicIncludeGSGChannel', false)

    return {
      // Dynamic
      userLogins: [] as Array<String>,
      streams,
      useDevApi,
      includeGSGChannel,
      streamsLoaded: false,

      multi,
    }
  },
  getters: {
    apiBaseUrl(): string {
      return this.useDevApi ? apiBaseUrlDev : apiBaseUrlProd
    },
    axios(): AxiosInstance {
      return axios.create({
        baseURL: this.apiBaseUrl,
        withCredentials: false,
        headers: {
          'Content-Type': staticContentType,
        },
      })
    },
    testStreams(): Streams {
      // A static list of names for testing, for now
      const listOfStreamerLogins = [
        "litui", "pulsaroctopus", "actitect", "jaynothin", "betaunits", "youropponent0"
      ]

      var output = {} as Streams

      listOfStreamerLogins.forEach((v: string) => {
        output[v] = {
          user_login: v,
          user_name: v,
        } as StreamHistoryDatum
      })

      return output
    },
  },
  actions: {
    addRemoveGSG() {
      var foundGSG: String = ""

      const logins = Object.keys(this.streams)
      logins.forEach((login: String) => {
          foundGSG = loginMatchesGSG(login) ? login : ""
          return
      });

      if (foundGSG) {
        delete this.streams[gsgChannelLogin]
      } else {
        // Maybe we should act here, but I'd rather just wait for the counter
      }
    },
    async loadStreams() {
      const store = this;

      if (store.multi.disableStreamLoading) {
        Object.values(store.testStreams).forEach((stream: StreamHistoryDatum) => {
          if (!store.includeGSGChannel) {
            if (stream.user_login === gsgChannelLogin) {
                return
            }
          }
          // append to our list
          store.streams[stream.user_login] = stream
          if (!store.userLogins.includes(stream.user_login)) {
            store.userLogins.push(stream.user_login)
          }
        });

        Object.keys(store.streams).forEach((login: string) => {
            // Replace displayed entry if the user login matches.
            if (!store.testStreams[login]) {
                delete store.streams[login]

                const lIndex = store.userLogins.findIndex((v) => v == login)
                if (lIndex > -1) {
                  store.userLogins.splice(lIndex, 1)
                }
            }
        })

        store.streamsLoaded = true
        return
      }

      await this.axios.get(staticPublicStreamEndpoint)
        .then(async (response) => {
          if (response.status != 200) {
              return
          }

          if (response.data && response.data.count && response.data.data) {
              var responseUserLogins = response.data.data.map((x: StreamHistoryDatum) => x.user_login)

              response.data.data.forEach((stream: StreamHistoryDatum) => {
                  if (!store.includeGSGChannel) {
                    if (stream.user_login === gsgChannelLogin) {
                        return
                    }
                  }

                  // append to our list
                  store.streams[stream.user_login] = stream
                  if (!store.userLogins.includes(stream.user_login)) {
                    store.userLogins.push(stream.user_login)
                  }
              });

              Object.keys(store.streams).forEach((login: string) => {
                  // Replace displayed entry if the user login matches.
                  const respI = responseUserLogins.indexOf(login)
                  if (respI == -1) {
                      delete store.streams[login]

                      const lIndex = store.userLogins.findIndex((v) => v == login)
                      if (lIndex > -1) {
                        store.userLogins.splice(lIndex, 1)
                      }
                  }
              })

              store.streamsLoaded = true
          }
        })
    },
    toggleDevApi() {
      this.useDevApi = !this.useDevApi
    },
    toggleIncludeGSG() {
      this.includeGSGChannel = !this.includeGSGChannel
    },
    toggleStreamsLoaded() {
      this.streamsLoaded = !this.streamsLoaded
    },
  }
})

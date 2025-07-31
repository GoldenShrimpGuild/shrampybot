import { defineStore } from 'pinia'
import { useLocalStorage } from '@vueuse/core'
import { apiBaseUrlDev, apiBaseUrlProd } from '../global-store'
import { useMultiStore } from './multi'
import axios from 'axios'
import type { AxiosInstance } from 'axios'
import type { StreamHistoryDatum } from '../../../model/utility/nosqldb/index';
import { Streams } from './classes'

// Data for axios
const jsonContentType = "application/json"
const publicStreamEndpoint = "/public/stream"

// This is a store for use with public pages such as gsgMultiTwitch
export const usePublicStore = defineStore('public', {
  state: () => {
    // Specific child store just for multiTwitch stuff
    // testing this pattern out
    const multi = useMultiStore()

    const includeGSGChannel = useLocalStorage('publicIncludeGSGChannel', false)
    const streamsMap = new Streams<string, StreamHistoryDatum>([], includeGSGChannel.value)

    // More global options for streams loader
    const useDevApi = useLocalStorage('publicUseDevApi', false)

    return {
      streamsMap: streamsMap,
      useDevApi,
      includeGSGChannel,
      streamsLoaded: false,

      // windowHeight and windowWidth in REM
      currentWindowHeight: 0,
      currentWindowWidth: 0,

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
          'Content-Type': jsonContentType,
        },
      })
    },
    testStreams(): Streams<string, StreamHistoryDatum> {
      // A static list of names for testing, for now
      const listOfStreamerLogins = [
        "litui", "pulsaroctopus", "actitect", "jaynothin", "betaunits", "youropponent0"
      ]

      var output = [] as StreamHistoryDatum[]

      listOfStreamerLogins.forEach((v: string) => {
        output.push({
          user_login: v,
          user_name: v,
        } as StreamHistoryDatum)
      })

      return new Streams(output)
    },
    streamsList(): StreamHistoryDatum[] {
      return Array.from(this.streamsMap.values())
    },
    userLogins(): string[] {
      return Array.from(this.streamsMap.keys())
    },
  },
  actions: {
    async loadStreams() {
      const store = this;

      if (store.multi.disableStreamLoading) {
        store.streamsMap.reconcile(store.testStreams.values().toArray())
        store.streamsLoaded = true
        return
      }

      await this.axios.get(publicStreamEndpoint)
        .then(async (response) => {
          if (response.status != 200) {
              return
          }

          if (response.data && response.data.count && response.data.data) {
            store.streamsMap.reconcile(response.data.data)
          }
          store.streamsLoaded = true
        })
    },
    toggleDevApi() {
      this.useDevApi = !this.useDevApi
    },
    toggleIncludeGSG() {
      this.includeGSGChannel = !this.includeGSGChannel
      this.streamsMap.setHideGSG(this.includeGSGChannel)
    },
    toggleStreamsLoaded() {
      this.streamsLoaded = !this.streamsLoaded
    }
  }
})

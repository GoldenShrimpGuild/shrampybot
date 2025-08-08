import { defineStore } from 'pinia'
import { useLocalStorage } from '@vueuse/core'
import { apiBaseUrlDev, apiBaseUrlProd } from '../global-store'
import { useMultiStore } from './multi'
import axios from 'axios'
import type { AxiosInstance } from 'axios'
import type { StreamHistoryDatum as Stream } from '../../../model/utility/nosqldb/index';
import { Streams } from './classes'
import { faker } from '@faker-js/faker'

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
    const streamsMap = new Streams<string, Stream>([], includeGSGChannel.value)
    const disableStreamLoading = useLocalStorage('disableStreamLoading', false)
    const testStreamCount = useLocalStorage('testStreamCount', 5)

    // More global options for streams loader
    const useDevApi = useLocalStorage('publicUseDevApi', false)

    return {
      streamsMap: streamsMap,
      useDevApi,
      includeGSGChannel,
      streamsLoaded: false,
      disableStreamLoading,

      testStreamCount,
      testStreams: [] as Stream[],

      // windowHeight and windowWidth in REM
      currentWindowHeight: 0,
      currentWindowWidth: 0,

      // current hue degrees, starting on a random value
      rainbowUint8: Math.floor(Math.random()*256),
      rainbowIncrement: 2,

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
    streamsList(): Stream[] {
      return Array.from(this.streamsMap.values())
    },
    userLogins(): string[] {
      return Array.from(this.streamsMap.keys())
    },
    currentRainbowColour(): string {
      // ported from https://github.com/UnexpectedMaker/esp32s3-arduino-helper/blob/main/src/UMS3.h
      // MIT License
      // Original code Copyright (c) 2022 Unexpected Maker
      var r = 0, g = 0, b = 0
      var pos = this.rainbowUint8

      if (pos < 85) {
        r = 255 - pos * 3
        g = pos * 3
        b = 0
      } else if (pos < 170) {
        pos -= 85
        r = 0
        g = 255 - pos * 3
        b = pos * 3
      } else {
        pos -= 170
        r = pos * 3
        g = 0
        b = 255 - pos * 3
      }

      const color = `#${r.toString(16).padStart(2, '0')}${g.toString(16).padStart(2, '0')}${b.toString(16).padStart(2, '0')}`
      return color
    }
  },
  actions: {
    async loadStreams() {
      const store = this;

      if (store.disableStreamLoading) {
        store.streamsMap.reconcile(store.testStreams).then(() => {
          store.streamsLoaded = true
        })
        return
      }

      this.axios.get(publicStreamEndpoint)
        .then(async (response) => {
          if (response.status != 200) {
              return
          }

          if (response.data && response.data.count && response.data.data) {
            store.streamsMap.reconcile(response.data.data).then(() => {
              store.streamsLoaded = true
            })
          }
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
    },
    toggleStreamLoading() {
      this.disableStreamLoading = !this.disableStreamLoading
    },
    incrementRainbowColour() {
      this.rainbowUint8 = (this.rainbowUint8 + this.rainbowIncrement) % 256
    },
    async generateTestStreams() {
      var output = [] as Stream[]

      const newAnimal = () => {
        const animal = faker.animal.cat().toLowerCase().replace(/\s/g, "")
        const colour = faker.color.human().replace(/\s/g, "")
        return `gsg_${colour}_${animal}`
      }

      Array.from({length: this.testStreamCount}).forEach(() => {
        var login = newAnimal()

        // Regenerate until unique
        while (output.findIndex((v) => v.user_login === login) > -1) {
          login = newAnimal()
        }

        output.push({
          user_login: login,
          user_name: login,
          title: faker.hacker.phrase()
        } as Stream)
      })
      this.testStreams = output

      return output
    },
  }
})

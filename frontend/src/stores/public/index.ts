import { defineStore } from 'pinia'
import { useLocalStorage } from '@vueuse/core'
import { apiBaseUrlDev, apiBaseUrlProd } from '../global-store'
import { useMultiStore } from './multi'
import axios from 'axios'
import { GSGStream, GSGStreams } from './classes'
import { faker } from '@faker-js/faker'

// types
import type { AxiosInstance } from 'axios'
import type { Event } from '../../../model/gsg-event-web/event'

// Data for axios
const jsonContentType = "application/json"
const publicStreamEndpoint = "/public/stream"
const currentEventEndpoint = "https://event.gsg.live/event/current"
const currentEventLoadInterval = 600000

// stream constants
const musicCategoryId = "26936"
const musicCategoryName = "Music"

// This is a store for use with public pages such as gsgMultiTwitch
export const usePublicStore = defineStore('public', {
  state: () => {
    // Specific child store just for multiTwitch stuff
    // testing this pattern out
    const multi = useMultiStore()

    const includeGSGChannel = useLocalStorage('publicIncludeGSGChannel', false)
    const onlyMusicCategory = useLocalStorage('onlyMusicCategory', true)
    const streamsMap = new GSGStreams<string, GSGStream>([], includeGSGChannel.value, onlyMusicCategory.value)
    const disableStreamLoading = useLocalStorage('disableStreamLoading', false)
    const testStreamCount = useLocalStorage('testStreamCount', 5)
    const streamsListSort = useLocalStorage('streamsListSort', {field: "started_at", direction: 1})

    const filterPriority = useLocalStorage('filterPriority', 0)
    const includeFilterWords = useLocalStorage('includeFilterWords', [] as string[])
    const excludeFilterWords = useLocalStorage('excludeFilterWords', [] as string[])

    // More global options for streams loader
    const useCurrentEventData = useLocalStorage('useCurrentEventData', true)
    const useDevApi = useLocalStorage('publicUseDevApi', false)

    return {
      streamsListSort,
      streamsMap: streamsMap,
      useDevApi,
      includeGSGChannel,
      streamsLoaded: false,
      disableStreamLoading,
      useCurrentEventData,
      currentEventData: {} as Event,
      lastCurrentEventLoadAttempt: 0,
      onlyMusicCategory,

      filterPriority,
      includeFilterWords,
      excludeFilterWords,

      testStreamCount,
      testStreams: [] as GSGStream[],

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
    streamsList(): GSGStream[] {
      const sortField = this.streamsListSort.field
      const sortDir = this.streamsListSort.direction
      const includeFilters = this.includeFilterWords
      const excludeFilters = this.excludeFilterWords

      const results = Array.from(this.streamsMap.values())
        .filter((stream) => {
          // Always allow event streams through
          if (stream.isEventStream) {
            return true
          }

          // if include filter is empty, match all
          var matchedIncludeFilter = includeFilters.length > 0 ? false : true
          includeFilters.forEach((filt) => {
            if (stream.title.toLowerCase().search(filt) > -1) {
              matchedIncludeFilter = true
              console.log(`Include filter [${filt}] matched ${stream.user_login}`)
              return
            }
          })
          var matchedExcludeFilter = false
          excludeFilters.forEach((filt) => {
            if (stream.title.toLowerCase().search(filt) > -1) {
              matchedExcludeFilter = true
              console.log(`Exclude filter [${filt}] matched ${stream.user_login}`)
              return
            }
          })

          if (matchedIncludeFilter && !matchedExcludeFilter) {
            return true
          }

          if (matchedIncludeFilter && matchedExcludeFilter) {
            if (this.filterPriority === 0) {
              // include wins
              return true
            } else {
              // exclude wins
              return false
            }
          }
          
          return false
        })

      // sort by field
      results.sort((a: GSGStream, b: GSGStream) => a[sortField] < b[sortField] ? -sortDir : (a[sortField] > b[sortField] ? sortDir : 0))

      // sort based on event stream priority
      results.sort((a, b) => a.isEventStream && b.isEventStream ? 0 : (a.isEventStream ? -1 : 1))
      return results
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
    },
  },
  actions: {
    // Loads from currentEventData if scheduled and 
    async loadStreams() {
      const thisLoadTime = Date.now()
      const store = this;

      if (store.disableStreamLoading) {
        store.streamsMap.reconcile(store.testStreams).then(() => {
          store.streamsLoaded = true
        })
        // stop all other loading routines
        return
      }

      var tempStreams = new Map<string, GSGStream>()

      // Only load new current event list if interval has elapsed
      if (thisLoadTime >= this.lastCurrentEventLoadAttempt + currentEventLoadInterval) {
        // Only do anything if using the data is enabled
        if (this.useCurrentEventData) {
          try {
            const response = await this.axios.get(currentEventEndpoint)
            if (response.status === 200) {
              if (response.data && response.data.title) {
                this.currentEventData = response.data
                // update the last load time since we were successful
                this.lastCurrentEventLoadAttempt = thisLoadTime
              }
            }
          } catch (error) {
            // no current event. Ho hum, carry on anyway.
          }
        }
      }

      // compile current event data if enabled and populated
      if (this.useCurrentEventData && this.currentEventData) {
        const curTime = Math.floor(Date.now()/1000)
        // console.log(`Current time: ${Math.floor(Date.now()/1000)}`)
        // console.log(`Event time: ${this.currentEventData.start} - ${this.currentEventData.end}`)
        // Check event start & end times
        if (this.currentEventData.end >= curTime && this.currentEventData.start <= curTime) {
          this.currentEventData.raidTrains.forEach((raidTrain, raidIndex) => {
            // console.log(`RaidTrain ${raidIndex} time: ${raidTrain.start}`)
            // Check raidTrain start & end times
            if (raidTrain.end >= curTime && raidTrain.start <= curTime) {
              raidTrain.schedule.forEach((stream) => {
                // Check schedule entry start & end times
                if (stream.end >= curTime && stream.start <= curTime) {
                  // console.log(`Streamer ${stream.twitchName} time: ${stream.start}`)
                  if (stream.twitchName) {
                    const tempEntry: GSGStream = {
                      id: "",
                      user_id: "",
                      user_login: stream.twitchName,
                      user_name: stream.twitchName,
                      thumbnail_url: stream.avatarUrl || "",
                      title: this.currentEventData.title,
                      game_id: musicCategoryId,
                      game_name: musicCategoryName,
                      tag_ids: [],
                      tags: [],
                      is_mature: false,
                      shrampybot_filtered: false,
                      type: "live",
                      viewer_count: -1,
                      started_at: (new Date(stream.start*1000)).toISOString(),
                      language: "en",
                      // Additional currentEvent fields
                      isEventStream: true,
                      isNormalStream: false,
                      raidTrain: raidIndex,
                    }
                    // console.log(tempEntry)
                    tempStreams.set(stream.twitchName, tempEntry)
                  }
                }
              })
            }
          })
        }
      }

      this.axios.get(publicStreamEndpoint)
        .then(async (response) => {
          if (response.status != 200) {
              return
          }

          if (response.data && response.data.count && response.data.data) {
            const sl: GSGStream[] = response.data.data
            sl.sort((a, b) => a.started_at.localeCompare(b.started_at))
            sl.forEach((stream: GSGStream) => {
              // Persist flag for ready identification of event stream
              if (tempStreams.has(stream.user_login)) {
                const ts = tempStreams.get(stream.user_login)
                stream.isEventStream = ts?.isEventStream
                stream.raidTrain = ts?.raidTrain
              } else {
                stream.isEventStream = false
                stream.raidTrain = -1
              }
              stream.isNormalStream = true
              tempStreams.set(stream.user_login, stream)
            })
            store.streamsMap.reconcile(Array.from(tempStreams.values())).then(() => {
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
    toggleCurrentEventData() {
      this.useCurrentEventData = !this.useCurrentEventData
      if (!this.useCurrentEventData) {
        this.currentEventData = {} as Event
        this.lastCurrentEventLoadAttempt = 0
      }
      this.loadStreams()
    },
    toggleOnlyMusicCategory() {
      this.onlyMusicCategory = !this.onlyMusicCategory
      this.streamsMap.setOnlyMusic(this.onlyMusicCategory)
    },
    // Filter priority on set filter is a 0 or 1 representing include filters or exclude filters
    setFilter(includes: string[], excludes: string[], priority: number) {
      this.filterPriority = priority
      this.includeFilterWords = includes
      this.excludeFilterWords = excludes
    },
    async generateTestStreams() {
      var output = [] as GSGStream[]

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
          title: faker.hacker.phrase(),
          game_id: musicCategoryId,
          game_name: musicCategoryName,
        } as GSGStream)
      })
      this.testStreams = output

      return output
    },
  }
})

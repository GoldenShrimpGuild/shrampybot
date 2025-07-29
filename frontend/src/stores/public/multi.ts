import { defineStore } from 'pinia'
import { useLocalStorage } from '@vueuse/core'
import { usePublicStore } from '.'
import type { VaMenuOption } from 'vuestic-ui/dist/types/components/va-menu-list/types.js'

const fixedChatWidth = 300
const fixedCardDecorHeight = 60

export const useMultiStore = defineStore('multi', {
    state: () => {
        // GSGMultiTwitch Settings
        const centreFrames = useLocalStorage('mtCentreFrames', true)
        const disableStreamLoading = useLocalStorage('mtDisableStreamLoading', false)
        const forceSkeleton = useLocalStorage('mtForceSkeleton', false)
        const hideChat = useLocalStorage('mtHideChat', false)
        const hideDecor = useLocalStorage('mtHideDecor', false)
        const testMode = useLocalStorage('mtTestMode', false)
        const useDropdownChatSelect = useLocalStorage('mtUseDropdownChatSelect', false)

        return {
            // Locally Stored
            centreFrames,
            disableStreamLoading,
            forceSkeleton,
            hideChat,
            hideDecor,
            testMode,
            useDropdownChatSelect,

            // Element Dimensions
            wrapperHeight: 0,
            wrapperWidth: 0,
            streamsContainerWidth: 0,
            headerHeight: 0,
            sideButtonsHeight: 0,
            chatSelectorHeight: 0,
            currentWindowHeight: 0,
            chatEmbedHeight: 0,

            streamWidth: 0,
            streamHeight: 0,
            streamTopPadding: 0,
        }
    },
    getters: {
      sideMenu(): Array<VaMenuOption> {
        const ps = usePublicStore()

        return [
          {
            text: "Layout",
            value: "",
            disabled: true,
            handler: (o: object) => {},
          },
          {
            text: "Centre-Align Streams",
            value: "centreFramesToggle",
            disabled: false,
            icon: this.centreFrames ? 'check_box' : 'check_box_outline_blank',
            handler: (o: object) => {
              ps.multi.toggleCentreFrames()
            },
          },
          {
            text: "Dropdown Chat Selector",
            value: "dropdownChat",
            disabled: false,
            icon: this.useDropdownChatSelect ? 'check_box' : 'check_box_outline_blank',
            handler: (o: object) => {
              ps.multi.toggleUseDropdownChatSelect()
            },
          },
          {
            text: "",
            value: "",
            disabled: true,
            handler: (o: object) => {},
          },
          {
            text: "Data",
            value: "",
            disabled: true,
            handler: (o: object) => {},
          },
          {
            text: "Include GSG Channel",
            value: "includeGSGChannel",
            disabled: false,
            icon: ps.includeGSGChannel ? 'check_box' : 'check_box_outline_blank',
            handler: (o: object) => {
              ps.toggleIncludeGSG()
              ps.addRemoveGSG()
            }
          },
          {
            text: "",
            value: "",
            disabled: true,
            handler: (o: object) => {},
          },
          {
            text: "Testing",
            value: "",
            disabled: true,
            handler: (o: object) => {},
          },
          {
            text: "Disable Data Feed & Use Fake Data",
            value: "fakeDataMode",
            disabled: false,
            icon: this.testMode ? 'check_box' : 'check_box_outline_blank',
            handler: (o: object) => {
              ps.multi.testMode = !ps.multi.testMode
              ps.multi.disableStreamLoading = ps.multi.disableStreamLoading

              ps.loadStreams();
            },
          },
          {
            text: "Fake Display Elements",
            value: "toggleSkeletons",
            disabled: false,
            icon: this.disableStreamLoading ? 'check_box' : 'check_box_outline_blank',
            handler: (o: object) => {
              ps.multi.disableStreamLoading = !ps.multi.disableStreamLoading

              ps.loadStreams();
            },
          },
          {
            text: "",
            value: "",
            disabled: true,
            handler: (o: object) => {},
          },
          {
            text: "About the Golden Shrimp Guild",
            value: "",
            disabled: true,
            handler: (o: object) => {},
          },
          {
            text: "GSG.live Website",
            value: "gsgWebsite",
            disabled: false,
            icon: "link",
            handler: (o: object) => {
              const gsgUrl = "https://www.gsg.live"

              window.open(gsgUrl, '_blank')?.focus();
            },
          },
          {
            text: "GSG Twitch Channel",
            value: "gsgTwitchChannel",
            disabled: false,
            icon: "link",
            handler: (o: object) => {
              const gsgUrl = "https://www.twitch.tv/goldenshrimpguild"

              window.open(gsgUrl, '_blank')?.focus();
            },
          },
          {
            text: "GSG Twitch Team",
            value: "gsgTwitchTeam",
            disabled: false,
            icon: "link",
            handler: (o: object) => {
              const gsgUrl = "https://www.twitch.tv/team/gsg"

              window.open(gsgUrl, '_blank')?.focus();
            },
          },
        ]
    },
    dynamicChatWidth(): number {
      return this.hideChat ? 0 : fixedChatWidth
    },
    dynamicCardDecorHeight(): number {
      return this.hideDecor ? 0 : fixedCardDecorHeight
    },
  },
  actions: {
    calculateSizes() {
      const arbitraryEdgeSpacerValue = 5

      // First, the chat size.
      this.chatEmbedHeight = window.innerHeight - this.sideButtonsHeight - this.chatSelectorHeight - arbitraryEdgeSpacerValue

      // logic here borrowed from www.multitwitch.tv
      // licensed under terms: "The code of this project is free to use"
      // repo: https://github.com/bhamrick/multitwitch

      const ps = usePublicStore()

      var finalWidth = 0
      var finalHeight = 0
      var topPadding = 0

      const numStreams = Object.keys(ps.streams).length

      const areaHeight = this.currentWindowHeight - this.headerHeight
      const areaWidth = this.wrapperWidth

      const cardDecorHeight = this.dynamicCardDecorHeight

      for (var perRow = 1; perRow <= numStreams; perRow++) {
        const rows = Math.ceil(numStreams / perRow)
        const cardDecorHR = cardDecorHeight*rows
        var maxWidth = Math.floor(areaWidth / perRow)
        var maxHeight = Math.floor(areaHeight / rows) - cardDecorHR

        if ((maxWidth * 9/16) < maxHeight) {
            maxHeight = (maxWidth * 9/16)
        } else {
            maxWidth = (maxHeight * 16/9)
        }
        if (maxWidth > finalWidth) {
            finalWidth = maxWidth
            finalHeight = maxHeight + cardDecorHR*9/16
            topPadding = (areaHeight - (rows * maxHeight)) / 2
        }
      }

      this.streamWidth = finalWidth
      this.streamHeight = finalHeight
      this.streamTopPadding = topPadding
    },
    toggleCentreFrames() {
      this.centreFrames = !this.centreFrames
    },
    toggleChat() {
      this.hideChat = !this.hideChat
    },
    toggleDecor() {
      this.hideDecor = !this.hideDecor
    },
    toggleForceSkeleton() {
      this.forceSkeleton = !this.forceSkeleton
    },
    toggleStreamLoading() {
      this.disableStreamLoading = !this.disableStreamLoading
    },
    toggleTestMode() {
      this.testMode = !this.testMode
    },
    toggleUseDropdownChatSelect() {
       this.useDropdownChatSelect = !this.useDropdownChatSelect
    },
  },
})
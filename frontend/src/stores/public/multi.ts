import { defineStore } from 'pinia'
import { useLocalStorage } from '@vueuse/core'
import { usePublicStore } from '.'
import type { VaMenuOption } from 'vuestic-ui/dist/types/components/va-menu-list/types.js'

const fixedChatWidth = 320
const fixedButtonsHeight = 20
const fixedButtonsWidth = 160
const fixedCardDecorHeight = 100
const fixedHeaderHeight = 60
const fixedWindowPadding = 5
const fixedChatSelectorHeight = 50

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
            currentWindowWidth: 0,
            currentWindowHeight: 0,

            streamWidth: 0,
            streamHeight: 0,
            streamTopPadding: 0,
            numRows: 1,

            streamXPositions: new Map<string, number>(),
            streamYPositions: new Map<string, number>(),
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
              ps.multi.toggleCentre()
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
    dynamicHeaderHeight(): number {
      return this.hideDecor ? 0 : fixedHeaderHeight; // height in rem
    },
    dynamicWrapperHeight(): number {
      return this.currentWindowHeight - fixedHeaderHeight - fixedWindowPadding*2
    },
    dynamicWrapperWidth(): number {
      return this.currentWindowWidth - fixedWindowPadding*2 - (this.hideChat ? 0 : fixedChatWidth)
    },
    dynamicSidebarWidth(): number {
      return this.hideChat ? fixedButtonsWidth : fixedChatWidth
    },
    dynamicSidebarHeight(): number {
      return this.hideChat ? fixedButtonsHeight : this.currentWindowHeight + fixedButtonsHeight
    },
    fixedChatWidth(): number {
      return fixedChatWidth
    },
    fixedCardDecorHeight(): number {
      return fixedCardDecorHeight
    },
    dynamicCardDecorHeight(): number {
      return this.hideDecor ? 0 : fixedCardDecorHeight
    },
    fixedButtonsHeight(): number {
      return fixedButtonsHeight
    },
    fixedChatSelectorHeight(): number {
      return fixedChatSelectorHeight
    },
    dynamicChatEmbedHeight(): number {
      return this.dynamicSidebarHeight - this.fixedButtonsHeight - this.fixedChatSelectorHeight
    }
  },
  actions: {
    calculateSizes() {
      // logic here borrowed from www.multitwitch.tv
      // licensed under terms: "The code of this project is free to use"
      // repo: https://github.com/bhamrick/multitwitch

      const ps = usePublicStore()

      var finalWidth = 0
      var finalHeight = 0
      var topPadding = 0

      const numStreams = ps.streamsList.length

      const areaHeight = this.dynamicWrapperHeight
      const areaWidth = this.dynamicWrapperWidth

      var perRow = 1
      var numRows = 0

      for (perRow = 1; perRow <= numStreams; perRow++) {
        numRows = Math.ceil(numStreams / perRow)
        var maxWidth = Math.floor(areaWidth / perRow)
        var maxHeight = Math.floor((areaHeight) / numRows)

        if ((maxWidth * 9/16) < maxHeight) {
            maxHeight = Math.floor(maxWidth * 9/16)
        } else {
            maxWidth = Math.floor(maxHeight * 16/9)
        }

        if (maxWidth > finalWidth) {
            if (areaHeight < areaWidth) {
              finalHeight = maxHeight
              finalWidth = Math.floor(finalHeight * 16/9)
            } else {
              finalWidth = maxWidth
              finalHeight = Math.floor(finalWidth * 9/16)
            }
            topPadding = (areaHeight - (numRows * finalHeight)) / 2
        }
      }

      this.streamWidth = finalWidth - this.dynamicCardDecorHeight*numRows*16/9
      this.streamHeight = finalHeight - this.dynamicCardDecorHeight*numRows
      this.streamTopPadding = topPadding
      this.numRows = numRows
      // console.log(this.streamWidth)
      // console.log(this.streamHeight)
    },
    toggleCentre() {
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
    setWindowSize(width: number, height: number) {
      this.currentWindowWidth = width
      this.currentWindowHeight = height
    },
  },
})
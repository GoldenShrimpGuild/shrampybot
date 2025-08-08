import { defineStore } from 'pinia'
import { useLocalStorage } from '@vueuse/core'

export const useMultiStore = defineStore('multi', {
    state: () => {
        // GSGMultiTwitch Settings
        const centreFrames = useLocalStorage('mtCentreFrames', true)
        const forceSkeleton = useLocalStorage('mtForceSkeleton', false)
        const hideChat = useLocalStorage('mtHideChat', false)
        const hideDecor = useLocalStorage('mtHideDecor', false)
        const testMode = useLocalStorage('mtTestMode', false)
        const useDropdownChatSelect = useLocalStorage('mtUseDropdownChatSelect', false)
        const startMuted = useLocalStorage('mtStartMuted', false)

        return {
            // Locally Stored
            centreFrames,
            forceSkeleton,
            hideChat,
            hideDecor,
            testMode,
            useDropdownChatSelect,
            startMuted,

            numRows: 1,
        }
    },
  actions: {
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
    toggleTestMode() {
      this.testMode = !this.testMode
    },
    toggleUseDropdownChatSelect() {
       this.useDropdownChatSelect = !this.useDropdownChatSelect
    },
    toggleStartMuted() {
      this.startMuted = !this.startMuted
    },
  },
})
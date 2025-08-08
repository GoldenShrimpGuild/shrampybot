<script lang="ts" setup>
import { computed, onBeforeMount, onMounted, ref, watch, watchEffect, reactive } from 'vue'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import { usePublicStore } from '../../stores/public/'
import { useMultiStore } from '../../stores/public/multi'
import { useTimer } from 'vue-timer-hook'
import { useWindowSize } from '@vueuse/core'

// Types
import type { StreamHistoryDatum as Stream } from '../../../model/utility/nosqldb'
import type { Streams } from '../../stores/public/classes'

// Components
import Sidebar from '../../components/multi/Sidebar.vue'
import StreamViewer from '../../components/multi/StreamViewer.vue'

const { t } = useI18n()

const PS = usePublicStore()
const MS = useMultiStore()
const { 
  streamsList,
  streamsLoaded,
  disableStreamLoading,
  currentRainbowColour,
  useDevApi,
  testStreamCount,
  includeGSGChannel,
} = storeToRefs(PS)
const { 
  hideChat,
  hideDecor,
  centreFrames,
  forceSkeleton,
  startMuted
} = storeToRefs(MS)

const lastStreamLoadTime = Date.now()
const streamLoadInterval = 20000
const streamLoadTimer = useTimer(lastStreamLoadTime + streamLoadInterval)

const lastColourRefreshTime = Date.now()
const colourRefreshInterval = 250
const colourRefreshTimer = useTimer(lastColourRefreshTime + colourRefreshInterval)

// fixed sizes
const maxSidebarWidth = 320

// sizing calcs
const { width: windowWidth, height: windowHeight } = useWindowSize()

const wrapperPadding = ref(0) // px
const wrapperRightMargin = computed(() => hideChat.value ? 0 : maxSidebarWidth)
const wrapperWidth = computed(() => windowWidth.value - wrapperPadding.value - wrapperRightMargin.value)
const wrapperHeight = computed(() => windowHeight.value)

const callStreamRefresh = () => {
  PS.loadStreams()

  const time = Date.now()
  streamLoadTimer.restart(time + streamLoadInterval)
}

const callColourRefresh = () => {
  PS.incrementRainbowColour()
  const time = Date.now()
  colourRefreshTimer.restart(time + colourRefreshInterval)
}

const handleFakeData = (disStreamLoading: boolean, skeletons: boolean) => {
  if (disStreamLoading) {
    streamsLoaded.value = false
    disableStreamLoading.value = disStreamLoading
    forceSkeleton.value = skeletons
    PS.generateTestStreams().then(async () => {
      PS.loadStreams()
    })
  } else {
    streamsLoaded.value = false
    disableStreamLoading.value = disStreamLoading
    forceSkeleton.value = skeletons
    PS.loadStreams()
  }
}

const openUrlInNewTab = (url: string) => {
  window.open(url, '_blank')?.focus();
}

// watch(hideChat, async (newValue) => await updateSizeParams())
// watch(hideDecor, async (newValue) => await updateSizeParams())

onBeforeMount(async () => {
  watchEffect(async () => {
    if (streamLoadTimer.isExpired.value) {
      callStreamRefresh()
    }
    if (forceSkeleton.value) {
      if (colourRefreshTimer.isExpired.value) {
        callColourRefresh()
      }
    }
  })

  await PS.loadStreams()
})

onMounted(async () => {
  if (disableStreamLoading) {
    handleFakeData(disableStreamLoading.value, forceSkeleton.value)
  }
})
</script>

<style lang="css">
body {
  text-overflow: ellipsis;
  overflow: hidden;
  white-space: nowrap;
}
</style>

<template>
  <VaLayout
    :right="{fixed: true}"
  >
    <template #content>
      <StreamViewer
        :width="wrapperWidth"
        :height="wrapperHeight"
        :hide-decor="hideDecor"
        :hide-chat="hideChat"
        :force-skeleton="forceSkeleton"
        :streams-loaded="streamsLoaded"
        :centre-frames="centreFrames"
        :streams-list="streamsList"
        :current-rainbow-colour="currentRainbowColour"
        :start-muted="startMuted"
      >
      </StreamViewer>
    </template>

    <template #right>
      <Sidebar
        :max-width="maxSidebarWidth"
        :max-height="windowHeight"
        :test-stream-count="testStreamCount"
        :centre-frames="centreFrames"
        :current-rainbow-colour="currentRainbowColour"
        :disable-stream-loading="disableStreamLoading"
        :force-skeleton="forceSkeleton"
        :include-g-s-g-channel="includeGSGChannel"
        :use-dev-a-p-i="useDevApi"
        :hide-chat="hideChat"
        :hide-decor="hideDecor"
        :streams-loaded="streamsLoaded"
        :streams-list="streamsList"
        :start-muted="startMuted"
        @toggle-chat="MS.toggleChat()"
        @toggle-decor="MS.toggleDecor()"
        @toggle-centre="MS.toggleCentre()"
        @toggle-dev-a-p-i="PS.toggleDevApi()"
        @toggle-include-g-s-g="PS.toggleIncludeGSG()"
        @toggle-start-muted="MS.toggleStartMuted()"
        @set-fake-data="handleFakeData"
        @handle-url="openUrlInNewTab"
        @set-test-stream-count="(count) => {testStreamCount = count}"
      ></Sidebar>
    </template>
  </VaLayout>
</template>
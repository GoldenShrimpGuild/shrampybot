<script lang="ts" setup>
import { computed, onBeforeMount, onMounted, ref, watch, watchEffect } from 'vue'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import { usePublicStore } from '../../stores/public/'
import { useMultiStore } from '../../stores/public/multi'
import { useTimer } from 'vue-timer-hook'

// Types
import type { StreamHistoryDatum as Stream } from '../../../model/utility/nosqldb'
import type { Streams } from '../../stores/public/classes'

// Components
import Sidebar from '../../components/multi/Sidebar.vue'
import TwitchStream from '../../components/multi/TwitchStream.vue'

const { t } = useI18n()

const PS = usePublicStore()
const MS = useMultiStore()
const { streamsList, userLogins, streamsLoaded } = storeToRefs(PS)
const { 
  hideChat,
  fixedChatWidth,
  dynamicSidebarWidth,
  dynamicSidebarHeight,
  hideDecor,
  fixedCardDecorHeight,
  dynamicCardDecorHeight,
  centreFrames,
  streamWidth,
  streamHeight,
  streamTopPadding,
  useDropdownChatSelect,
  dynamicChatEmbedHeight,
  fixedChatSelectorHeight,
  currentWindowHeight,
  currentWindowWidth,
  dynamicHeaderHeight,
  dynamicWrapperWidth,
  dynamicWrapperHeight,
  fixedButtonsHeight,
  numRows,
  sideMenu
} = storeToRefs(MS)

const lastStreamLoadTime = Date.now()
const streamLoadInterval = 2000
const streamLoadTimer = useTimer(lastStreamLoadTime + streamLoadInterval)

// ID globals
const mtWrapperId = ref("mt-wrapper")
const mtHeaderId = ref("mt-header")

// class globals
const streamsClass = ref("mt-streams")
const cardClass = ref("mt-card")

var lastSizeUpdateTime = Date.now()

const updateSizeParams = async (force: boolean = false) => {
  const currentTime = Date.now()
  if (force || currentTime > lastSizeUpdateTime+200) {
    MS.setWindowSize(window.innerWidth, window.innerHeight)
    MS.calculateSizes()
    lastSizeUpdateTime = currentTime
  }
}

const callStreamRefresh = async () => {
  await PS.loadStreams()
  await updateSizeParams()

  const time = Date.now()
  streamLoadTimer.restart(time + streamLoadInterval)
}

watch(hideChat, async (newValue) => await updateSizeParams())
watch(hideDecor, async (newValue) => await updateSizeParams())

onBeforeMount(async () => {
  watchEffect(async () => {
    if (streamLoadTimer.isExpired.value) {
      await callStreamRefresh()
    }
  })

  window.addEventListener("resize", async (event) => {
    await updateSizeParams()
  })

  await PS.loadStreams()
})

onMounted(async () => {
  await updateSizeParams()
})
</script>

<style scoped>
#mt-wrapper {
  padding: none;
  margin: none;
  height: v-bind(currentWindowHeight + 'px');
  width: v-bind(dynamicWrapperWidth + 'px');
}

#mt-wrapper > div > h1 {
  height: v-bind(dynamicHeaderHeight + 'px');
  width: v-bind(dynamicWrapperWidth + 'px');
}

#mt-header {
  height: v-bind(dynamicHeaderHeight + 'px');
  width: v-bind(dynamicWrapperWidth + 'px');
}

.mt-streams {
  height: v-bind(dynamicWrapperHeight + 'px');
  width: v-bind(dynamicWrapperWidth + 'px');
}
</style>

<template>
  <VaLayout
    :right="{fixed: true}"
  >
    <template #content>
      <div :id="mtWrapperId">
        <div :id="mtHeaderId">
          <h1 :class="`pl-4 pt-2 pb-0 font-bold`" :hidden="hideDecor" onshow="console.log('showed')">
            {{ t('multiTwitch.header') }}
          </h1>
        </div>
        <div :class="streamsClass" :aria-busy="!streamsLoaded">
          <!-- ${centreFrames ? 'place-content-center' : ''} -->
          <VaSkeletonGroup
            v-if="!streamsLoaded"
            :class="``"
            animation="wave"
            :delay="0"
          >
            <VaCard
                v-for="stream in streamsList"
                :key="stream.user_login"
                :class="`h-[${streamHeight}px] w-[${streamWidth}px]`"
            >
              <VaSkeleton 
                variant="squared"
                :class="`h-[${streamHeight}px] w-[${streamWidth}px]`"
                :height="`${streamHeight}px`"
                :width="`${streamWidth}px`"/>
              <VaCardContent class="flex items-center" :hidden="hideDecor" >
                <VaSkeleton variant="text" class="ml-2 va-text" :lines="2" />
              </VaCardContent>
            </VaCard>
          </VaSkeletonGroup>
          <div v-else
            :class="`flex flex-wrap ${centreFrames ? 'place-content-center' : ''} w-[${dynamicWrapperWidth}px] h-[${dynamicWrapperHeight}px]`">
              <VaCard 
                v-for="stream in streamsList"
                :key="stream.user_login"
                :class="`${cardClass} item w-[${streamWidth}px] h-[${streamHeight}px]`"
              >
                <TwitchStream
                  :user-login="stream.user_login"
                  :width="streamWidth"
                  :height="streamHeight-dynamicCardDecorHeight"
                  @start-embed="updateSizeParams()"
                ></TwitchStream>
                <VaCardContent 
                  :style="`width: ${streamWidth}px; `"
                  :hidden="hideDecor"
                >
                  <h5 :style="`text-overflow: ellipsis; overflow: hidden; white-space: nowrap;`">{{ stream.user_name }}</h5>
                  <p :style="`text-overflow: ellipsis; overflow: hidden; white-space: nowrap;`">
                    {{ stream.title }}
                  </p>
                </VaCardContent>
              </VaCard>
          </div>
        </div>
      </div>
    </template>

    <template #right>
      <Sidebar
        :dynamic-sidebar-width
        :dynamic-sidebar-height
        :streams-list
        :hide-chat
        :hide-decor
        :fixed-chat-selector-height
        :use-dropdown-chat-select
        :dynamic-chat-embed-height
        :side-menu
        :fixed-buttons-height
        @toggle-chat="MS.toggleChat()"
        @toggle-decor="MS.toggleDecor()"
      ></Sidebar>
    </template>
  </VaLayout>
</template>
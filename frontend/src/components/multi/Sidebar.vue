<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import type { StreamHistoryDatum as Stream } from '../../../model/utility/nosqldb'

// components
import SidebarNav from './SidebarNav.vue'
import ChatTabs from './ChatTabs.vue'
import TwitchChat from '../../components/multi/TwitchChat.vue'
import { useWindowSize } from '@vueuse/core';

const { t } = useI18n()

const sidebarClass = ref("mt-sidebar")
const chatSelectorClass = ref("chat-selector")
const embedContainerClass = ref("chat-embed-container")

// fixed sizing
const minSidebarWidth = 160
const fixedChatSelectorHeight = 35
const fixedButtonsHeight = 20

const fixedModalWidth = 500

const props = defineProps<{
  maxWidth: number,
  maxHeight: number,
  streamsList: Stream[],
  streamsLoaded: boolean,
  useDevAPI: boolean,
  disableStreamLoading: boolean,
  includeGSGChannel: boolean,
  currentRainbowColour: string,
  hideChat: boolean,
  hideDecor: boolean,
  centreFrames: boolean,
  testStreamCount: number,
  forceSkeleton: boolean,
  startMuted: boolean,
}>()

const emit = defineEmits<{
  (e: "toggleChat"): void,
  (e: "toggleDecor"): void,
  (e: "toggleCentre"): void,
  (e: "toggleIncludeGSG"): void,
  (e: "toggleDevAPI"): void,
  (e: "setFakeData", disableStreamLoading: boolean, forceSkeleton: boolean): void,
  (e: "setTestStreamCount", streamCount: number): void,
  (e: "recalculateSize"): void,
  (e: "handleUrl", url: string): void,
  (e: "toggleStartMuted"): void,
}>()

const fakeDataModal = ref(false)
const fakeStreamCount = computed({
  get: () => props.testStreamCount,
  set: (newValue) => emit("setTestStreamCount", newValue)
})

const chatTabSelected = ref(0)
const currentChatStream = computed(() => props.streamsList.at(chatTabSelected.value) || {} as Stream)
const currentChatLogin = computed(() => props.streamsList.at(chatTabSelected.value)?.user_login || "")

const sidebarWidth = computed(() => props.hideChat ? minSidebarWidth : props.maxWidth)
const sidebarHeight = computed(() => props.hideChat ? fixedButtonsHeight : props.maxHeight)
const chatEmbedHeight = computed(() => props.maxHeight - fixedButtonsHeight - fixedChatSelectorHeight)

const {width: windowWidth, height: windowHeight} = useWindowSize()
const modalWidth = computed(() => windowWidth.value >= fixedModalWidth ? fixedModalWidth : windowWidth.value)

const showFakeDataModal = () => {
  fakeDataModal.value = true
}

</script>

<style scoped>
.mt-sidebar {
  float: right;
  margin: none;
  padding: none;
  text-overflow: ellipsis;
  overflow: hidden;
  white-space: nowrap; 
}
.va-skeleton * {
  margin: 0;
  padding: 0;
}
.va-skeleton div.skeleton-header {
  width: 100%;
  font-weight: bold;
  font-size: 24pt;
  text-align: center; 
}
.va-skeleton div.skeleton-content {
  width: 100%;
  padding-top: 5px;
  font-weight: bold;
  font-size: 12pt;
  text-align: center; 
}
</style>

<template>
  <div 
    :class="`${sidebarClass}`"
    :style="`height: ${sidebarHeight}px; width: ${sidebarWidth}px`"
  >
    <SidebarNav
      :height="fixedButtonsHeight"
      :width="sidebarWidth"
      :hide-chat="hideChat"
      :hide-decor="hideDecor"
      :centre-frames="centreFrames"
      :disable-stream-loading="disableStreamLoading"
      :include-g-s-g-channel="includeGSGChannel"
      :start-muted="startMuted"
      :use-dev-a-p-i="useDevAPI"
      @toggle-centre="emit('toggleCentre')"
      @toggle-chat="emit('toggleChat')"
      @toggle-decor="emit('toggleDecor')"
      @toggle-fake-data="showFakeDataModal()"
      @toggle-include-g-s-g="emit('toggleIncludeGSG')"
      @toggle-dev-a-p-i="emit('toggleDevAPI')"
      @handle-url="(url) => { emit('handleUrl', url)}"
      @toggle-start-muted="emit('toggleStartMuted')"
    ></SidebarNav>
    <div
      :hidden="hideChat"
      :class="`${chatSelectorClass} bg-[] row`"
    >
      <ChatTabs
        :streams-list="streamsList"
        :height="fixedChatSelectorHeight"
        :width="sidebarWidth"
        background-color="gsgDarkYellow"
        @select="(tab) => chatTabSelected = tab"
      ></ChatTabs>
    </div>
    <div
      :class="embedContainerClass"
      :hidden="hideChat"
      :style="`height: ${chatEmbedHeight}px; width: ${sidebarWidth}`"
    >
      <template
        v-for="stream in streamsList"
        :key="stream.user_login"
      >
        <VaSkeleton
          v-if="forceSkeleton || !streamsLoaded"
          variant="squared"
          :height="`${chatEmbedHeight}px`"
          :color="currentRainbowColour"
          :aria-label="t('multiTwitch.chatPlaceholder')"
        >
          <div
            class="skeleton-header"
            :style="`padding-top: ${chatEmbedHeight / 2 - 20}px;`"
          >
            {{ t('multiTwitch.chatPlaceholder') }}
          </div>
          <div
            class="skeleton-content"
          >
            {{ streamsList.at(chatTabSelected)?.user_login }}
          </div>
        </VaSkeleton>
        <TwitchChat
          v-else
          :startMuted="startMuted"
          :hidden="currentChatLogin != stream.user_login"
          :user-login="stream.user_login"
          :height="`${chatEmbedHeight}px`"
        ></TwitchChat>
      </template>
    </div>
    <VaModal
      v-model="fakeDataModal"
      :ok-text="t('multiTwitch.enable')"
      :cancel-text="t('multiTwitch.disable')"
      :max-width="`${modalWidth}px`"
      @ok="emit('setFakeData', true, true)"
      @cancel="emit('setFakeData', false, false)"
    >
      <h4 class="va-h4">
        {{ t("multiTwitch.fakeData") }}
      </h4>
      <p class="m-0">
        {{ t("multiTwitch.fakeDataGenText") }}
      </p>
      <VaForm>
        <VaFormField>
          <VaSlider
            v-model="fakeStreamCount"
            class="mt-8"
            :step="1"
            :max="15"
            :min="1"
            track-label-visible
          />
        </VaFormField>
      </VaForm>
    </VaModal>
  </div>
</template>
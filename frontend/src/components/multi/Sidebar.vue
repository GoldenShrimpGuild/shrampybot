<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'

// types
import type { GSGStream as Stream } from '../../stores/public/classes'
import type { SelectableOption } from 'vuestic-ui/dist/types/composables'

// components
import { VaSkeleton, VaModal, VaInput, VaSwitch, VaChip, VaForm, VaFormField, VaSlider } from 'vuestic-ui'
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

const isDropdownOpen = false

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
  onlyMusicCategory: boolean,
  useCurrentEventData: boolean,
  includeFilterWords: string[],
  excludeFilterWords: string[],
  filterPriority: number,
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
  (e: "toggleCurrentEventData"): void,
  (e: "toggleOnlyMusicCategory"): void,
  (e: "setFilter", includes: string[], excludes: string[], priority: number): void,
}>()

const fakeDataModal = ref(false)
const fakeStreamCount = computed({
  get: () => props.testStreamCount,
  set: (newValue) => emit("setTestStreamCount", newValue)
})

interface filterObj {
  includes: string[],
  excludes: string[],
}
const filterModal = ref(false)
const filterPriorityRef = ref(props.filterPriority === 0 ? false : true)
const includeFilter = ref(props.includeFilterWords.map((v) => v) as string[])
const excludeFilter = ref(props.excludeFilterWords.map((v) => v) as string[])
const filterIncludeInput = ref("")
const filterExcludeInput = ref("")

const processFilterKey = (event: any) => {
  const endKeyCodes = [
    13, // Enter
  ]
  if (endKeyCodes.includes(event.keyCode)) {
    if (event.target.name === 'filterInclude') {
      const valStr = filterIncludeInput.value.toLowerCase().trim()
      if (!includeFilter.value.includes(valStr)) {
        includeFilter.value.push(valStr)
        filterIncludeInput.value = ""
      }
    } else if (event.target.name === 'filterExclude') {
      const valStr = filterExcludeInput.value.toLowerCase().trim()
      if (!excludeFilter.value.includes(valStr)) {
        excludeFilter.value.push(valStr)
        filterExcludeInput.value = ""
      }
    }
  }
}

const finalizeFilterSetting = () => {
  emit("setFilter",
    includeFilter.value.map((v) => v.valueOf()),
    excludeFilter.value.map((v) => v.valueOf()),
    filterPriorityRef.value ? 1 : 0,
  )
}

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

const showFilterModal = () => {
  filterModal.value = true
}

</script>

<style scoped>
.mt-sidebar {
  margin: none;
  padding: none;
  text-overflow: ellipsis;
  overflow: hidden;
  white-space: nowrap;
  position: absolute;
  top: 0px;
  right: 0px;
  z-index: 999;
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
    :style="`height: ${sidebarHeight}px; width: ${sidebarWidth}px;`"
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
      :use-current-event-data="useCurrentEventData"
      :only-music-category="onlyMusicCategory"
      @toggle-centre="emit('toggleCentre')"
      @toggle-chat="emit('toggleChat')"
      @toggle-decor="emit('toggleDecor')"
      @toggle-fake-data="showFakeDataModal()"
      @toggle-include-g-s-g="emit('toggleIncludeGSG')"
      @toggle-dev-a-p-i="emit('toggleDevAPI')"
      @handle-url="(url) => { emit('handleUrl', url)}"
      @toggle-start-muted="emit('toggleStartMuted')"
      @toggle-only-music-category="emit('toggleOnlyMusicCategory')"
      @toggle-current-event-data="emit('toggleCurrentEventData')"
      @adjust-filter="showFilterModal()"
    ></SidebarNav>
    <div
      :hidden="hideChat"
      :class="`${chatSelectorClass} row`"
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
      :no-outside-dismiss="true"
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
    <VaModal
      v-model="filterModal"
      :ok-text="t('multiTwitch.set')"
      :no-outside-dismiss="true"
      :max-width="`${modalWidth}px`"
      @ok="finalizeFilterSetting()"
      @cancel=""
    >
      <h4 class="va-h4">
        {{ t("multiTwitch.adjustFilters") }}
      </h4>
      <p class="mt-4 mb-4 text-wrap whitespace-pre-line">
        {{ t("multiTwitch.adjustFiltersDescription") }}
      </p>
      <div>
        <VaSwitch
          v-model="filterPriorityRef"
          class="mb-4"
          size="small"
          color="gsgRed"
          offColor="gsgYellow"
          :true-label="t('multiTwitch.adjustFilterPriorityExclude')"
          :false-label="t('multiTwitch.adjustFilterPriorityInclude')"
        />
      </div>
      <VaInput
        v-model="filterIncludeInput"
        name="filterInclude"
        :label="t('multiTwitch.adjustIncludeFilterLabel')"
        :placeholder="t('multiTwitch.adjustIncludeFilterPlaceholder')"
        @keydown="processFilterKey"
      >
      </VaInput>
      <div class="pt-2 pb-2 flex flex-wrap gap-1">
        <VaChip 
          v-for="entry, index in includeFilter"
          :key="index"
          class="item font-bold"
          color="gsgYellow"
          closeable
          size="small"
          @update:model-value="includeFilter.splice(index, 1)"
        >
          {{ entry }}
        </VaChip>
      </div>
      <VaInput
        v-model="filterExcludeInput"
        name="filterExclude"
        :label="t('multiTwitch.adjustExcludeFilterLabel')"
        :placeholder="t('multiTwitch.adjustExcludeFilterPlaceholder')"
        @keydown="processFilterKey"
      >
      </VaInput>
      <div class="pt-2 pb-2 flex flex-wrap gap-1">
        <VaChip 
          v-for="entry, index in excludeFilter"
          :key="index"
          class="item font-bold"
          color="gsgRed"
          closeable
          size="small"
          @update:model-value="excludeFilter.splice(index, 1)"
        >
          {{ entry }}
        </VaChip>
      </div>
    </VaModal>
  </div>
</template>
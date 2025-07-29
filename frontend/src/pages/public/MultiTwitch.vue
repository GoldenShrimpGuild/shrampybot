<script lang="ts" setup>
import { computed, onBeforeMount, onMounted, ref, watch, watchEffect } from 'vue'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import StreamCard from '../../components/multi/StreamCard.vue'
import TwitchChat from '../../components/multi/TwitchChat.vue'
import { usePublicStore, type Streams } from '../../stores/public/'
import { useTimer } from 'vue-timer-hook'
import type { StreamHistoryDatum } from '../../../model/utility/nosqldb'

const { t } = useI18n()

const PS = usePublicStore()
const MS = PS.multi
const { streams, userLogins, streamsLoaded } = storeToRefs(PS)
const { 
  hideChat,
  dynamicChatWidth,
  hideDecor,
  dynamicCardDecorHeight,
  centreFrames,
  streamWidth,
  streamHeight,
  streamTopPadding,
  useDropdownChatSelect,
  chatEmbedHeight,
  wrapperHeight,
  wrapperWidth,
  streamsContainerWidth,
  headerHeight,
  sideButtonsHeight,
  chatSelectorHeight,
  currentWindowHeight,

  sideMenu,
} = storeToRefs(MS)

const time = Date.now()

const refreshDelay = 2000
const timer = useTimer(time + refreshDelay)

const chatTabSelected = ref(0)
const currentChatLogin = computed(() => userLogins.value[chatTabSelected.value] as string)
const currentChatStream = computed(() => streams.value[currentChatLogin.value])

const updateSizeParams = async () => {
  const gsgButtons = document.getElementById('gsgButtons');
  sideButtonsHeight.value = gsgButtons ? gsgButtons.getBoundingClientRect().height : 0

  const chatTabs = document.getElementById('chatTabs');
  chatSelectorHeight.value = chatTabs ? chatTabs.getBoundingClientRect().height : 0

  const mtHeader = document.getElementById("mtHeader")
  headerHeight.value = mtHeader ? mtHeader.getBoundingClientRect().height : 0

  currentWindowHeight.value = window.innerHeight
  wrapperHeight.value = currentWindowHeight.value
  wrapperWidth.value = window.innerWidth
  streamsContainerWidth.value = wrapperWidth.value - dynamicChatWidth.value - 5

  MS.calculateSizes()
}

const callStreamRefresh = async () => {
  await PS.loadStreams()
  await updateSizeParams()

  const time = Date.now()
  timer.restart(time + refreshDelay)
}

onBeforeMount(async () => {
  watchEffect(async () => {
    if (timer.isExpired.value) {
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

<template>
  <VaLayout :class="`h-full w-full`" style="overflow: hidden;">
    <template #content>
      <div id="mtWrapper" 
        :class="`box-content static p-0 m-0 h-[${wrapperHeight}px] w-[${wrapperWidth}px]`"
        :style="`overflow: hidden;`"
      >
        <h1 id="mtHeader" :class="`pl-4 pt-2 pb-0 font-bold h-[${headerHeight}px]`" :hidden="hideDecor">{{ t('multiTwitch.header') }}</h1>
        <div id="mtStreams" :class="`h-[${wrapperHeight - headerHeight}px] w-[${streamsContainerWidth}px]`" :aria-busy="!streamsLoaded"
          :style="`overflow: hidden; padding-top: ${streamTopPadding}px; padding-right: ${dynamicChatWidth}px;`"
        >
          <VaSkeletonGroup
            v-if="!streamsLoaded"
            :class="`w-full ${centreFrames ? 'place-content-center' : ''} gap-2`"
            animation="wave"
            :delay="0"
          >
            <VaCard
                v-for="stream in Object.values(streams)"
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
          <div 
            v-else
            :class="`flex flex-wrap w-[${streamsContainerWidth}px] h-[${wrapperHeight-headerHeight}px] ${centreFrames ? 'place-content-center' : ''} gap-2`"
          >
              <StreamCard
                v-for="stream in Object.values(streams)"
                :key="stream.user_login"
                :stream="stream"
                :height="streamHeight"
                :width="streamWidth"
                :cardOffset="dynamicCardDecorHeight"
                :decorations="hideDecor"
                @start-embed="updateSizeParams()"
              ></StreamCard>
          </div>
        </div>
      </div>
    </template>

    <template #right>
      <div 
        id="mtSidebar"
        :class="`fixed top-0 right-0 h-${hideChat ? 0 : 'full'} row`"
      >
          <div 
            id="gsgButtons"
            class="pl-2 pr-1 pt-2 pb-2 inline-grid grid-cols-3 w-full content-box"
          >
            <VaButton
              round
              class="item mr-1"
              size="small"
              color="gsgYellow"
              :icon="hideDecor ? 'check_box_outline_blank' : 'check_box'"
              v-on:click="MS.toggleDecor()"
            >
              {{ t("multiTwitch.decor") }}
            </VaButton>
            <VaButton
              round
              class="item mr-1"
              size="small"
              color="gsgYellow"
              :icon="hideChat ? 'check_box_outline_blank' : 'check_box'"
              v-on:click="MS.toggleChat()"
            >
              {{ t("multiTwitch.chat") }}
            </VaButton>
            <VaMenu
              :options="sideMenu"
              @selected="(v) => v.handler()"
              disabled-by="disabled"
            >
              <template #anchor>
                <VaButton
                  round
                  outline
                  class="item mr-1"
                  size="small"
                  color="gsgYellow"
                  :icon="'menu'"
                >
                </VaButton>
              </template>
            </VaMenu>
        </div>
        <div :hidden="hideChat" class="pt-1 h-full bg-[var(--va-background-element)]">
          <VaTabs 
            v-if="!useDropdownChatSelect"
            id="chatTabs" 
            :class="`w-[${() => dynamicChatWidth}]`"
            v-model="chatTabSelected"
            :hide-pagination="false"
          >
            <template #tabs>
              <VaTab
                v-for="stream in Object.values(streams)"
                :key="stream.user_login"
                :label="stream.user_name"
              >
                {{ stream.user_name }}
              </VaTab>
            </template>
          </VaTabs>
          <VaSelect
            v-else
            v-model="chatTabSelected"
            placeholder="Colored"
            color="#FFFFFF"
            :options="Object.values(streams)"
            inner-label
          />
          <div :class="`w-[${dynamicChatWidth}px] h-full`">
            <VaSkeletonGroup
                v-if="!streams"
                :class="`h-full w-[${dynamicChatWidth}px]`"
                animation="wave"
                :delay="0"
              >
              <div
                :style="`width: ${dynamicChatWidth}px;`"
              >
                  <VaSkeleton variant="text" class="ml-2 va-text" :lines="100" />
              </div>
            </VaSkeletonGroup>
            <template v-else v-for="stream in Object.values(streams)" :key="stream.user_login">
              <TwitchChat
                :hidden="currentChatLogin != stream.user_login"
                :user-login="stream.user_login"
                :height="chatEmbedHeight"
              ></TwitchChat>
            </template>
          </div>
        </div>
      </div>
    </template>
  </VaLayout>
</template>

<style lang="css">

#mtHeader {
  height: 45px;
}

#mtStreams {
    /* text-align: center; */
    float: left;
    margin: 0;
    /* margin-right: -100px; */
    padding: 0;
}

#mtStreams div {
  float: left;
}

#mtStreams .item {
  float: left;
}

iframe {
    border:0 none;
}

.fullwidth {
    width: 100%;
}

.left {
    float: left;
}

.right {
    float: right;
}

.centering {
    text-align: center;
}

.clear {
    clear: both;
}

</style>
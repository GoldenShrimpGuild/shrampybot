<script lang="ts" setup>
import { onBeforeMount, ref, watch, watchEffect } from 'vue'
import { storeToRefs } from 'pinia'
import { useAxios } from '../../plugins/axios'
import { useI18n } from 'vue-i18n'
import { useGlobalStore } from '../../stores/global-store'
import { useTimer } from 'vue-timer-hook'

// components
import StreamCard from '../../components/stream/StreamCard.vue'

const { t } = useI18n()

const streams = ref([])

const time = Date.now()
const timer = useTimer(time + 60000)

const GlobalStore = useGlobalStore()
const { isDevEnvironment } = storeToRefs(GlobalStore)

const axios = useAxios()

const loadStreams = async () => {
  await axios.get('/public/stream').then((response) => {
    streams.value = response.data.data
    // console.log(response.data.data[0])
    // console.log(resizeThumbnailUrl(response.data.data[0]["thumbnail_url"], 1280, 720))
  })
}

const heartbeatTimerRestart = async () => {
  // Monitor the current state of the authorization and refresh if it needs it
  loadStreams()

  const time = Date.now()
  timer.restart(time + 60000)
}

watch(isDevEnvironment, async (newValue, oldValue) => {
  loadStreams()
})

onBeforeMount(() => {
  watchEffect(async () => {
    if (timer.isExpired.value) {
      await heartbeatTimerRestart()
    }
  })

  loadStreams()
})

const resizeThumbnailUrl = (url: string, width: number, height: number) => {
  return url.replace('{width}', width.toString()).replace('{height}', height.toString())
}
</script>

<template>
  <h1 class="page-title font-bold">{{ t('menu.activeStreams') }}</h1>
  <div class="row">
    <div class="flex flex-wrap gap-5">
      <div v-for="stream in streams" class="item">
        <StreamCard
          :streamer="stream['user_name']"
          :title="stream['title']"
          :image_url="resizeThumbnailUrl(stream['thumbnail_url'], 1280, 720)"
        ></StreamCard>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { onBeforeMount, ref, watch, watchEffect } from 'vue'
import { storeToRefs } from 'pinia'
import axios from 'axios'
import { useI18n } from 'vue-i18n'
import StreamCard from '../../components/stream/StreamCard.vue'
import { useAuthStore } from '../../stores/auth'
import { useGlobalStore } from '../../stores/global-store'
import { useTimer } from 'vue-timer-hook'

const AuthStore = useAuthStore()
const { t } = useI18n()

const streams = ref([])

const time = Date.now()
const timer = useTimer(time + 60000)

const GlobalStore = useGlobalStore()
const { isDevEnvironment } = storeToRefs(GlobalStore)

const loadStreams = async () => {
  await axios.get('/public/stream', AuthStore.getAxiosConfig()).then((response) => {
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
  <!-- <section class="flex flex-col gap-4">
    <div class="flex flex-col sm:flex-row gap-4">
      <RevenueUpdates class="w-full sm:w-[70%]" />
      <div class="flex flex-col gap-4 w-full sm:w-[30%]">
        <YearlyBreakup class="h-full" />
        <MonthlyEarnings />
      </div>
    </div>
    <DataSection />
    <div class="flex flex-col md:flex-row gap-4">
      <RevenueByLocationMap class="w-full md:w-4/6" />
      <RegionRevenue class="w-full md:w-2/6" />
    </div>
    <div class="flex flex-col md:flex-row gap-4">
      <ProjectTable class="w-full md:w-1/2" />
      <Timeline class="w-full md:w-1/2" />
    </div>
  </section> -->
</template>

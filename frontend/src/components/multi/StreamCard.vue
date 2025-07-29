<template>
  <VaCard class="stream-card item" :style="`width: ${width}px;`">
      <iframe
          :id="iframeId"
          class="stream"
          :src="iframeUrl"
          :allowfullscreen="true"
          :style="`width: ${width}px; height: ${height-cardOffset}px;`"
          :height
          :width
          v-on:loadstart="emit('startEmbed')"
          v-on:load="emit('loadedEmbed')"
      ></iframe>
      <VaCardContent class="items-center" :hidden="decorations">
        <h5 style="text-overflow: ellipsis; overflow: hidden; white-space: nowrap;">{{ stream.user_name }}</h5>
        <p style="text-overflow: ellipsis; overflow: hidden; white-space: nowrap;">
          {{ stream.title }}
        </p>
      </VaCardContent>
  </VaCard>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import {  } from 'vue-router'
import { StreamHistoryDatum } from '../../../model/utility/nosqldb/index';

const thisHost = new URL(window.location.href).hostname
const embedUrl = `https://player.twitch.tv/?muted=true&channel=USER_LOGIN&parent=${thisHost}`

const props = defineProps<{
  stream: StreamHistoryDatum,
  width: number,
  height: number,
  cardOffset: number,
  decorations: boolean,
}>()

const emit = defineEmits<{
  startEmbed: [],
  loadedEmbed: []
}>()

const iframeId = computed(() => `embed_${props.stream.user_login}`)
const iframeUrl = computed(() => embedUrl.replace('USER_LOGIN', props.stream.user_login))

const iframeRef = ref({} as HTMLIFrameElement | null)

onMounted(async () => {
    // Get a handle to the iframe element
    iframeRef.value = document.getElementById(iframeId.value) as HTMLIFrameElement | null
})

</script>

<style lang="css">
</style>
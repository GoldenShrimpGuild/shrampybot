<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

const aspectRatio = 9/16;

const props = defineProps<{
  userLogin: string,
  width: number,
  startMuted: boolean,
}>()

const emit = defineEmits<{
  startEmbed: [],
  loadedEmbed: []
}>()

const thisHost = new URL(window.location.href).hostname
const embedUrl = `https://player.twitch.tv/?muted=START_MUTED_STATE&channel=USER_LOGIN&parent=${thisHost}`

const iframeId = computed(() => `embed_${props.userLogin}`)
const iframeUrl = computed(() => embedUrl.replace('USER_LOGIN', props.userLogin).replace('START_MUTED_STATE', props.startMuted ? 'true' : 'false'))

const iframeRef = ref({} as HTMLIFrameElement | null)

onMounted(async () => {
    // Get a handle to the iframe element
    iframeRef.value = document.getElementById(iframeId.value) as HTMLIFrameElement | null
})
</script>

<template>
  <iframe
      :id="iframeId"
      :class="`stream`"
      :src="iframeUrl"
      :allowfullscreen="true"
      :width="width"
      :height="width * aspectRatio"
      :style="`width: ${width}px; height: ${width * aspectRatio}px;`"
      allow="autoplay"
      v-on:loadstart="emit('startEmbed')"
      v-on:load="emit('loadedEmbed')"
  ></iframe>
</template>
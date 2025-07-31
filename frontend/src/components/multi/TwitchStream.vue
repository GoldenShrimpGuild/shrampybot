<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'

const thisHost = new URL(window.location.href).hostname
const embedUrl = `https://player.twitch.tv/?muted=true&channel=USER_LOGIN&parent=${thisHost}`

const props = defineProps<{
  userLogin: string,
  width: number,
  height: number,
}>()

const emit = defineEmits<{
  startEmbed: [],
  loadedEmbed: []
}>()

const iframeId = computed(() => `embed_${props.userLogin}`)
const iframeUrl = computed(() => embedUrl.replace('USER_LOGIN', props.userLogin))

const iframeRef = ref({} as HTMLIFrameElement | null)

onMounted(async () => {
    // Get a handle to the iframe element
    iframeRef.value = document.getElementById(iframeId.value) as HTMLIFrameElement | null
})
</script>

<template>
  <iframe
      :id="iframeId"
      :class="`stream w-full`"
      :src="iframeUrl"
      :allowfullscreen="true"
      :height="height"
      :width="width"
      v-on:loadstart="emit('startEmbed')"
      v-on:load="emit('loadedEmbed')"
  ></iframe>
</template>
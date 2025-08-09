<script setup lang="ts">
import { computed } from 'vue'

const thisHost = new URL(window.location.href).hostname
const embedUrl = `https://twitch.tv/embed/USER_LOGIN/chat?darkpopout&parent=${thisHost}`

const props = defineProps<{
  userLogin: string,
  height?: string | number,
}>()

const iframeId = computed(() => `embed_chat_${props.userLogin}`)
const iframeUrl = computed(() => embedUrl.replace('USER_LOGIN', props.userLogin))

</script>

<template>
  <iframe
    :id="iframeId"
    :class="`stream-chat h-full w-full`"
    :src="iframeUrl"
    :height="height"
    frameborder="0"
    scrolling="no"
  >
  </iframe>
</template>
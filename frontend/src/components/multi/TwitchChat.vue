<script setup lang="ts">
import { ref, computed } from 'vue'

const thisHost = new URL(window.location.href).hostname
const embedUrl = `https://twitch.tv/embed/USER_LOGIN/chat?darkpopout&parent=${thisHost}`

const props = defineProps<{
  userLogin: string,
  height: number,
}>()

const iframeId = computed(() => `embed_chat_${props.userLogin}`)
const iframeUrl = computed(() => embedUrl.replace('USER_LOGIN', props.userLogin))

</script>

<template>
  <iframe
    :class="`stream-chat v-full w-full`"
    :id="iframeId"
    :src="iframeUrl"
    width="100%"
    :height
    frameborder="0"
    scrolling="no"
  >
  </iframe>
</template>
<template>
  <span class="logo" color="primary">Shrampy<span v-if="!GlobalStore.$state.isDevEnvironment" class="bot_text">bot</span><span v-else class="bot_text">dev</span></span>
</template>

<script setup lang="ts">
import { useGlobalStore } from '../../stores/global-store';

const GlobalStore = useGlobalStore();
</script>

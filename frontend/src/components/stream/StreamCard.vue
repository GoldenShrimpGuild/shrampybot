<template>
  <div class="flex flex-col flex-nowrap gap-5">
    <VaCard style="width: 20rem">
      <VaImage class="h-52" :ratio="16 / 9" :src="image_url" />
      <VaCardBlock class="flex-nowrap">
        <div class="flex-auto">
          <VaCardTitle>{{ streamer }}</VaCardTitle>
          <VaCardContent>
            {{ title }}
          </VaCardContent>
        </div>
      </VaCardBlock>
    </VaCard>
  </div>
</template>

<script setup lang="ts">
const props = defineProps({
  image_url: {
    type: String,
    default: '',
  },
  title: {
    type: String,
    default: '',
  },
  streamer: {
    type: String,
    default: '',
  },
})
</script>

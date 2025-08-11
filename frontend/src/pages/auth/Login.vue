<template>
  <VaForm ref="form">
    <p class="text-base mb-4 leading-5">
      <ShrampybotLogo />
    </p>
    <div class="flex justify-center mt-4">
      <VaButton class="w-full" color="discordBlurple" :href="GlobalStore.getDiscordOAuthUrl()">
        <VaIcon :component="VaIconDiscord" />
        <span style="padding-left: 0.3rem">{{ t('auth.discordSignIn') }}</span>
      </VaButton>
    </div>
  </VaForm>
</template>

<script lang="ts" setup>
import { onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useGlobalStore } from '../../stores/global-store'

// components
import { VaForm, VaButton, VaIcon } from 'vuestic-ui'
import ShrampybotLogo from '../../components/logos/ShrampybotLogo.vue'
import VaIconDiscord from '../../components/icons/VaIconDiscord.vue'

const GlobalStore = useGlobalStore()

const { t } = useI18n()

onMounted(() => {
  // Force to prod environment to override previous localStorage setting
  GlobalStore.isDevEnvironment = false
})
</script>

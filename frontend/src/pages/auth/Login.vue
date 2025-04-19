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
    <div
      class="auth-layout__options flex flex-col sm:flex-row items-start sm:items-center justify-between"
      style="margin-top: 1rem"
    >
      <VaCollapse class="min-w-96" :header="t('auth.developerOptions')">
        <VaRadio v-model="environment" :options="[t('devTestEnvironment'), t('prodEnvironment')]" />
      </VaCollapse>
    </div>
  </VaForm>
</template>

<script lang="ts" setup>
import { reactive, ref, watch, onBeforeMount, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useForm, useToast } from 'vuestic-ui'
import { validators } from '../../services/utils'
import { useI18n } from 'vue-i18n'
import ShrampybotLogo from '../../components/logos/ShrampybotLogo.vue'
import VaIconDiscord from '../../components/icons/VaIconDiscord.vue'
import { useAuthStore } from '../../stores/auth'
import { useUserStore } from '../../stores/user'
import { useGlobalStore } from '../../stores/global-store'

const AuthStore = useAuthStore()
const UserStore = useUserStore()
const GlobalStore = useGlobalStore()

const { t } = useI18n()

const { push } = useRouter()
const { init } = useToast()

const environment = ref(GlobalStore.isDevEnvironment ? t('devTestEnvironment') : t('prodEnvironment'))

onMounted(() => {
  const { init, close, closeAll } = useToast()

  init({
    title: 'Cookie Agreement',
    message: `
        <p>ShrampyBot only makes minimal use of browser cookies for some of its core functionality but does use third party features (notably Twitch and Discord) that may involve cookies for tracking and other purposes.</p>
        <p>By continuing to log in to ShrampyBot you agree to the terms of our cookie usage as well as the terms of those third parties.</p>
      `,
    dangerouslyUseHtmlString: true,
    closeable: false,
    customClass: 'cookie-agreement',
    duration: 0,
    position: 'bottom-left',
  })
})

watch(environment, async (newVal) => {
  if (newVal === t('devTestEnvironment')) {
    GlobalStore.$state.isDevEnvironment = true
  } else {
    GlobalStore.$state.isDevEnvironment = false
  }
})
</script>

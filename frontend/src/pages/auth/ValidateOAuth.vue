<template>
  <div class="row">
    <div class="flex" width="100%">
      <div class="item">
        <VaModal v-model="show_modal" hide-default-actions no-dismiss blur>
          <template #default>
            <VaCardTitle>{{ progress_title }}</VaCardTitle>
            <VaCardContent>
              <div class="va-spacer"></div>
              <VaProgressBar :model-value="oauth_progress" size="large" class="oauth_progress" />
            </VaCardContent>
          </template>
        </VaModal>
        <VaModal v-model="show_error_modal" hide-default-actions no-dismiss blur>
          <template #default>
            <VaCardTitle>{{ t('auth.error') }}</VaCardTitle>
            <VaCardContent>
              <VaProgressBar
                :model-value="error_timeout"
                size="large"
                class="oauth_progress"
                color="danger"
                content-inside
              >
                {{ error_timeout_seconds }}
              </VaProgressBar>
            </VaCardContent>
          </template>
        </VaModal>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import router from '../../router'
import { useAuthStore } from '../../stores/auth'
import { useGlobalStore } from '../../stores/global-store'
import { isString } from 'lodash'
import axios from 'axios'

// components
import { VaModal, VaCardTitle, VaCardContent, VaProgressBar } from 'vuestic-ui'

const AuthStore = useAuthStore()
const GlobalStore = useGlobalStore()

// const axios = inject('axios') as AxiosInstance

const { t } = useI18n()

const route = useRoute()
const show_modal = ref(true)
const show_error_modal = ref(false)
const error_timeout = ref(0)
const error_timeout_seconds = ref(10)

const oauth_progress = ref(0)
const progress_title = ref(t('auth.title_oauth_validating'))

onMounted(() => {
  getQueryParams()
})

const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms))

const getQueryParams = async () => {
  await router.isReady()
  const code = route.query.code
  const error = route.query.error
  const state = GlobalStore.decodeStateBlob(String(route.query.state))

  if (error === 'access_denied') {
    encountered_error(route.query.error_description)
    return false
  }

  if (!isString(code)) {
    router.replace({ name: 'logout' })
    return false
  }
  const path = '/auth/validate'

  // Handle tokens
  await axios
    .post(
      path,
      {
        code: code,
      },
      {
        baseURL: GlobalStore.getApiBaseUrl(),
      },
    )
    .then(async (response) => {
      oauth_progress.value = 100
      progress_title.value = t('auth.title_oauth_synchronizing')
      if (GlobalStore.$state.isDevEnvironment) {
        AuthStore.$state.accessTokenDev = response.data.access
      } else {
        AuthStore.$state.accessTokenProd = response.data.access
      }
    })
    .catch((reason) => {
      encountered_error(reason)
    })

  router.push(state.redirect_path)
}

const encountered_error = async (reason: any) => {
  show_modal.value = false
  show_error_modal.value = true

  while (error_timeout_seconds.value > 0) {
    error_timeout_seconds.value -= 1
    error_timeout.value += 100 / 10
  }

  show_error_modal.value = false
  router.replace({ name: 'logout' })
}
</script>

<style>
.oauth_progress {
  --va-progress-bar-width: 300px;
}

.oauth_runner {
  --va-progress-bar-width: 290px;
  align-self: right;
}
</style>

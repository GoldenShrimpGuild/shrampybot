<script lang="ts" setup>
import { ref, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { useI18n } from 'vue-i18n'
import { storeToRefs } from 'pinia'
import { useGlobalStore } from '../../stores/global-store';
import { useAxios } from '../../plugins/axios'

const { t } = useI18n()

const eventId = ref('' as string | undefined)

const showICModal = ref(false)
const disableICModalButton = ref(true)
const cacheFlushSuccess = ref(false)

const GlobalStore = useGlobalStore()

const { isDevEnvironment } = storeToRefs(GlobalStore)
const axios = useAxios();

const route = useRoute();

const invalidateAPICache = async () => {
  if (isDevEnvironment.value) {
    disableICModalButton.value = true
    showICModal.value = true
    cacheFlushSuccess.value = false

    axios.get('/admin/event/' + eventId.value)
      .then((response) => {
        if (response.status === 200) {
          disableICModalButton.value = false
          cacheFlushSuccess.value = true
        }
      })
      .catch((error) => {
        disableICModalButton.value = false
        cacheFlushSuccess.value = false
      })
  }
}

onMounted(async () => {
  eventId.value = route.params.eventId as string
})

</script>

<template>
  <h1 class="page-title font-bold">{{ t('admin.events') }}</h1>

  <VaModal v-model="showICModal" hide-default-actions no-dismiss blur>
    <template #default>
      <VaCardTitle>{{ t('admin.event.modalTitle') }}</VaCardTitle>
      <VaCardContent>
        <VaProgressBar v-if="disableICModalButton" indeterminate size="large" class="oauth_progress" />
        <div v-else>
          <p v-if="cacheFlushSuccess">{{ t('admin.event.icSuccessFor') + eventId }}</p>
          <p v-else>{{ t('admin.event.icFailureFor') + eventId }}</p>
        </div>
      </VaCardContent>
    </template>
    <template #footer>
      <div class="flex gap-2">
        <VaButton color="gsgYellow" :disabled="disableICModalButton" @click="showICModal = !showICModal">OK</VaButton>
      </div>
    </template>
  </VaModal>

  <VaCard>
    <VaCardTitle>{{ t('admin.event.eventsApi') }}</VaCardTitle>
    <VaCardContent>
      <VaInput label="Event ID" v-model="eventId" style="width: 20rem;"></VaInput>

      <VaButton icon="add" style="margin: 1rem;" color="gsgYellow" round @click="invalidateAPICache">
        {{ t('admin.event.icButton') }}
      </VaButton>
    </VaCardContent>
  </VaCard>
</template>

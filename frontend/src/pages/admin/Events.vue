<script lang="ts" setup>
import { ref, onMounted, watch, getCurrentInstance } from 'vue';
import { useToast } from 'vuestic-ui';
import { useRoute } from 'vue-router';
import { useI18n } from 'vue-i18n'
import { useAxios } from '../../plugins/axios'
import { CurrentEventDatum } from '../../../model/utility/nosqldb'
import { CurrentEventPutRequestBody } from '../../../model/controller/admin'

const { t } = useI18n()
const toast = useToast()

const eventId = ref('' as string | undefined)
const currentEvent = ref({} as CurrentEventDatum)

const ICButtonLoading = ref(false)
const setCEButtonLoading = ref(false)
const clearCEButtonDisabled = ref(false)
const axios = useAxios();
const route = useRoute();

const invalidateAPICache = async () => {
  ICButtonLoading.value = true

  axios.get('/admin/event/' + eventId.value)
    .then((response) => {
      if (response.status === 200) {
        sendNormalToast(t("admin.event.icSuccess"))
        ICButtonLoading.value = false
      }
    })
    .catch((error) => {
      sendErrorToast(t("admin.event.icFailure"))
      ICButtonLoading.value = false
    })
}

const getCurrentEvent = async () => {
  axios.get('/admin/current_event/1')
    .then((response) => {
      if (response.status === 200) {
        if (response.data.status === "success") {
          currentEvent.value = response.data.currentEvent
        }
      }
    })
    .catch((error) => {
    })
}

const clearCurrentEvent = async () => {
  clearCEButtonDisabled.value = true

  const reqBody = {
    eventId: eventId.value
  } as CurrentEventPutRequestBody

  axios.delete('/admin/current_event/1')
    .then((response) => {
      if (response.status === 200) {
        if (response.data.status === "success") {
          sendNormalToast(t("admin.event.clearCESuccess"))
          currentEvent.value.eventId = ""
          currentEvent.value.title = ""
          currentEvent.value.description = ""
        } else {
          sendErrorToast(t("admin.event.clearCEFailure"))
        }
        clearCEButtonDisabled.value = false
      }
    })
    .catch((error) => {
      if (error.response.status === 502) {
        // Retry eternally when hitting 502s... seems to be AWS's fault
        clearCurrentEvent()
      } else {
        sendErrorToast(t("admin.event.clearCEFailure"))
        clearCEButtonDisabled.value = false
      }
    })
}

const setCurrentEvent = async () => {
  setCEButtonLoading.value = true

  const reqBody = {
    eventId: eventId.value
  } as CurrentEventPutRequestBody

  axios.put('/admin/current_event/1', reqBody)
    .then((response) => {
      if (response.status === 200) {
        if (response.data.status === "success") {
          currentEvent.value = response.data.currentEvent
          sendNormalToast(t("admin.event.setCESuccess"))
        } else {
          sendErrorToast(t("admin.event.setCEFailure"))
        }
        setCEButtonLoading.value = false
      }
    })
    .catch((error) => {
      if (error.response.status === 502) {
        // Retry eternally when hitting 502s... seems to be AWS's fault
        setCurrentEvent()
      } else {
        sendErrorToast(t("admin.event.setCEFailure"))
        setCEButtonLoading.value = false
      }
    })
}

const sendNormalToast = (message: string) => {
  toast.init({
    message: message,
    color: "gsgYellow"
  })
}

const sendErrorToast = (message: string) => {
  toast.init({
    message: message,
    color: "gsgYellow"
  })
}

onMounted(async () => {
  eventId.value = route.params.eventId as string

  currentEvent.value.eventId = ""
  currentEvent.value.title = ""
  currentEvent.value.description = ""
  getCurrentEvent()
})

</script>

<template>
  <h1 class="page-title font-bold">{{ t('admin.events') }}</h1>

  <div class="row">
    <div class="flex flex-col md12">
      <div class="item">
        <VaCard>
          <VaCardTitle>{{ t('admin.event.eventOperations') }}</VaCardTitle>
          <VaCardContent>
            <VaInput label="Event ID" v-model="eventId" style="width: 20rem;"></VaInput>

            <VaButton icon="va-clear" style="margin-top: 1rem; margin-left: 0.5rem;" color="gsgYellow" round @click="invalidateAPICache"
              :loading="ICButtonLoading"
              :disabled="ICButtonLoading"
            >
              {{ t('admin.event.icButton') }}
            </VaButton>
            <VaButton icon="star" style="margin: 1rem; margin-left: 0.5rem;" color="gsgYellow" round @click="setCurrentEvent"
              :loading="setCEButtonLoading"
              :disabled="setCEButtonLoading"
            >
              {{ t('admin.event.setCurrentEvent') }}
            </VaButton>
          </VaCardContent>
        </VaCard>
      </div>
    </div>
  </div>
  <br/>
  <div class="row">
    <div class="flex flex-col md12">
      <div class="item">
        <VaCard>
          <VaCardTitle>{{ t('admin.event.currentEvent') }}</VaCardTitle>
          <VaCardContent>
            <VaList class="indented-list" v-if="currentEvent.eventId">
              <VaListItemLabel>{{ t("admin.event.eventId") }}</VaListItemLabel>
              <VaListItem>{{ currentEvent.eventId }}</VaListItem>
              <VaListItemLabel>{{ t("admin.event.name") }}</VaListItemLabel>
              <VaListItem>{{ currentEvent.title }}</VaListItem>
              <VaListItemLabel>{{ t("admin.event.desc") }}</VaListItemLabel>
              <VaListItem>{{ currentEvent.description }}</VaListItem>
            </VaList>
            <p v-else>{{ t("admin.event.noCurrentEvent") }}</p>

            <VaButton icon="clear" style="margin-top: 1rem;" color="gsgRed" round @click="clearCurrentEvent"
              :disabled="!currentEvent.eventId"
            >
              {{ t('admin.event.clearCurrentEvent') }}
            </VaButton>
          </VaCardContent>
        </VaCard>
      </div>
    </div>
  </div>
</template>

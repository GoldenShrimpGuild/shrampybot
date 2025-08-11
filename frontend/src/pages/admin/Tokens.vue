<script lang="ts" setup>
import { ref, watch, onBeforeMount, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { storeToRefs } from 'pinia'
import { useStaticTokenStore } from '../../stores/static_token'
import { forEach } from 'lodash'

// components
import { VaCard, VaCardTitle, VaCardContent, VaButton, VaSwitch, VaForm, VaTextarea, VaIcon, VaCheckbox, VaDateInput, VaListLabel, VaOptionList, VaModal } from 'vuestic-ui'

// types
import type { NewTokenRequestBody, OutputStaticTokenInfo } from '../../../model/controller/admin'
import type { DateInputModelValue } from 'vuestic-ui/dist/types/components/va-date-input/types.js'

const StaticTokenStore = useStaticTokenStore()
const { tokens } = storeToRefs(StaticTokenStore)

const { t } = useI18n()

const addModalShow = ref(false)
const tokenRevealShow = ref(false)
const revokeModalShow = ref(false)
const expiresCheckbox = ref(false)
const showRevoked = ref(false)
const revokeModalId = ref('' as string | undefined)
const currentDate = ref(new Date())

onBeforeMount(() => {
  StaticTokenStore.fetchTokenInfo()
})

const scopeSelector = ref([
    {
        text: 'dev',
        value: 'dev',
        disabled: false,
    },
    {
        text: 'gsg',
        value: 'gsg',
        disabled: false,
    },
    {
        text: 'gsg:streamer',
        value: 'gsg:streamer',
        disabled: false,
    },
    {
        text: 'admin',
        value: 'admin',
        disabled: false,
    },
    {
        text: 'admin:categories',
        value: 'admin:categories',
        disabled: false
    },
    {
        text: 'admin:collection',
        value: 'admin:collection',
        disabled: false
    },
    {
        text: 'admin:events',
        value: 'admin:events',
        disabled: false
    },
    {
        text: 'admin:filters',
        value: 'admin:filters',
        disabled: false
    },
    {
        text: 'admin:users',
        value: 'admin:users',
        disabled: false
    }
] as Array<Record<string, any>>)

const toggleNarrowAdminScopes = (enabled: boolean) => {
  forEach(scopeSelector.value, (scope) => {
    if (scope.value.startsWith('admin:')) {
      scope.disabled = !enabled
    }
  })
}

const toggleNarrowGSGScopes = (enabled: boolean) => {
  forEach(scopeSelector.value, (scope) => {
    if (scope.value.startsWith('gsg:')) {
      scope.disabled = !enabled
    }
  })
}

const toggleNarrowDevScopes = (enabled: boolean) => {
  forEach(scopeSelector.value, (scope) => {
    if (scope.value.startsWith('dev:')) {
      scope.disabled = !enabled
    }
  })
}

const showModalAdd = (param: any) => {
  expiresCheckbox.value = false
  addModalForm.value.expires_at = null
  addModalForm.value.purpose = ''
  addModalForm.value.scopes = []
  addModalShow.value = !addModalShow.value
}
const addModalForm = ref({} as Record<string, any>)

const dateValidationRules = [
  (v: DateInputModelValue) => {
    if (!v) {
      return false || 'Cannot be blank.'
    }
    return v >= new Date() || 'Expiry cannot be in the past.'
  },
]

watch(addModalForm.value, (newVal) => {
  let adminSelected = false
  let devSelected = false
  let gsgSelected = false

  const dupeScopes = []

  forEach(newVal.scopes, (selectedScope) => {
    if (selectedScope.value == 'admin') {
      adminSelected = true
    }
    if (selectedScope.value == 'gsg') {
      gsgSelected = true
    }
    if (selectedScope.value == 'dev') {
      devSelected = true
    }
  })

  if (adminSelected) {
    toggleNarrowAdminScopes(false)
  } else {
    toggleNarrowAdminScopes(true)
  }

  if (gsgSelected) {
    toggleNarrowGSGScopes(false)
  } else {
    toggleNarrowGSGScopes(true)
  }

  if (devSelected) {
    toggleNarrowDevScopes(false)
  } else {
    toggleNarrowDevScopes(true)
  }
})

const sortedTokens = computed(() => {
  if (tokens) {
    const sortedCopy = [...StaticTokenStore.$state.tokens]
    function compare(a: OutputStaticTokenInfo, b: OutputStaticTokenInfo) {
      const a_created_date = new Date(a.created_at)
      const b_created_date = new Date(b.created_at)
      if (a_created_date < b_created_date) return -1
      if (a_created_date > b_created_date) return 1
      return 0
    }
    return sortedCopy.sort(compare)
  } else {
    return []
  }
})

const generateToken = async () => {
  const outRequest = {} as NewTokenRequestBody

  outRequest.purpose = addModalForm.value.purpose

  // Ensure date is correct
  if (addModalForm.value.expires_at) {
    outRequest.expires_at = addModalForm.value.expires_at.toISOString()
  } else {
    // Crank it to max for the lulz
    outRequest.expires_at = new Date(8640000000000).toISOString()
  }

  // Figure out scopes

  // Add login so this token is actually useful
  outRequest.scopes = []
  outRequest.scopes.push('login')

  forEach(addModalForm.value.scopes, (scope) => {
    if (!scope.disabled) {
      outRequest.scopes.push(scope.value)
    }
  })

  await StaticTokenStore.addToken(outRequest).then((token) => {
    tokenRevealShow.value = true
  })
}

const showModalRevoke = (id: string) => {
  revokeModalId.value = id

  revokeModalShow.value = true
}

const revokeToken = () => {
  StaticTokenStore.revokeToken(revokeModalId.value)

  revokeModalShow.value = false
}
</script>

<template>
  <h1 class="page-title font-bold">{{ t('admin.tokens') }}</h1>
  <section class="flex flex-col gap-4">
    <div class="row flex flex-col sm:flex-row gap-4">
      <VaCard>
        <VaCardTitle> Static Tokens </VaCardTitle>
        <VaCardContent>
          <div>
            <VaButton size="small" icon="add" style="margin-bottom: 1rem" color="gsgYellow" round @click="showModalAdd">
              Generate New Token
            </VaButton>
          </div>
          <VaSwitch
            v-model="showRevoked"
            label="Show revoked tokens"
            color="gsgYellow"
            size="small"
          />
          <div class="va-table-responsive">
            <table class="va-table">
              <thead>
                <tr>
                  <th>Valid</th>
                  <th>Created At (UTC)</th>
                  <th>Created By</th>
                  <th>Expires At (UTC)</th>
                  <th>Scopes</th>
                  <th>Purpose</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="token in sortedTokens" :key="token.id" rowspan="2">
                  <template v-if="showRevoked || showRevoked == token.revoked">
                    <td>
                      <VaIcon
                        v-if="!token.revoked && currentDate < new Date(token.expires_at)"
                        name="check"
                        color="#00ff00"
                      ></VaIcon>
                      <VaIcon v-else name="stop_circle" color="gsgRed"></VaIcon>
                    </td>
                    <td>{{ new Date(token.created_at).toDateString() }}</td>
                    <td>{{ token.creator_id }}</td>
                    <td>{{ new Date(token.expires_at).toDateString() }}</td>
                    <td>{{ token.scopes }}</td>
                    <td>{{ token.purpose }}</td>
                    <td>
                      <VaButton
                        v-if="!token.revoked"
                        size="small"
                        round
                        icon="stop_circle"
                        color="gsgRed"
                        @click="showModalRevoke(token.id)"
                        >Revoke</VaButton
                      >
                    </td>
                  </template>
                </tr>
              </tbody>
            </table>
          </div>
        </VaCardContent>
      </VaCard>
    </div>
  </section>
  <VaModal v-model="addModalShow" hide-default-actions blur>
    <template #header>
      <h5 class="va-h5">Generate New Token</h5>
    </template>
    <div class="flex flex-col items-start gap-2">
      <VaForm ref="addForm" immediate hide-error-messages class="flex flex-col gap-2 mb-2">
        <VaTextarea
          v-model="addModalForm.purpose"
          label="Purpose"
          name="StaticTokenPurpose"
          placeholder="describe how this token will be used"
        >
        </VaTextarea>
        <VaCheckbox v-model="expiresCheckbox" label="Set an expiry date"></VaCheckbox>
        <VaDateInput
          v-if="expiresCheckbox"
          v-model="addModalForm.expires_at"
          name="StaticTokenExpiryDate"
          label="Expiry Date"
          :rules="dateValidationRules"
        ></VaDateInput>
        <VaListLabel style="text-align: left">Scopes</VaListLabel>
        <VaOptionList v-model="addModalForm.scopes" :options="scopeSelector"></VaOptionList>
      </VaForm>
    </div>
    <template #footer>
      <div class="flex gap-2">
        <VaButton color="gsgYellow" @click="addModalShow = !addModalShow">Cancel</VaButton>
        <VaButton color="gsgYellow" @click="generateToken">Generate</VaButton>
      </div>
    </template>
  </VaModal>
  <VaModal v-model="tokenRevealShow" ok-text="OK" size="small" hide-default-actions blur no-dismiss>
    <p>Your new token is:</p>
    <VaTextarea v-model="StaticTokenStore.$state.newToken.token" style="width: 100%" autosize min-rows="6" readonly>
    </VaTextarea>
    <p>This is the only time you will be shown this token.</p>
    <template #footer>
      <div class="flex gap-2">
        <VaButton
          color="gsgYellow"
          @click="addModalShow = false; tokenRevealShow = false; StaticTokenStore.clearNewToken()"
          >OK
        </VaButton>
      </div>
    </template>
  </VaModal>
  <VaModal v-model="revokeModalShow" hide-default-actions blur ok-text="Apply">
    <h4 class="va-h5">Confirm Token Revocation</h4>
    <p>Are you sure you want to revoke the selected token?</p>
    <template #footer>
      <div class="flex gap-2">
        <VaButton color="gsgYellow" @click="revokeModalShow = !revokeModalShow">Cancel</VaButton>
        <VaButton color="gsgRed" @click="revokeToken">Revoke</VaButton>
      </div>
    </template>
  </VaModal>
</template>

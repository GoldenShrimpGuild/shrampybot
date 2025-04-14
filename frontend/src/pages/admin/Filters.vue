<script lang="ts" setup>
import { ref, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { storeToRefs } from 'pinia'
import { FilterDatum } from '../../../model/utility/nosqldb'
import { useGlobalStore } from '../../stores/global-store'
import { useFilterStore } from '../../stores/filters'

const { t } = useI18n()

const addEditModalShow = ref(false)
const deleteModalShow = ref(false)
const addEditModalForm = ref({
  id: '' as string | undefined,
  keyword: '' as string,
  isRegex: false as boolean | undefined,
  caseInsensitive: false as boolean,
})
const deleteModalId = ref('' as string | undefined)

const GlobalStore = useGlobalStore()
const FilterStore = useFilterStore()

const { isDevEnvironment } = storeToRefs(GlobalStore)

const showModalAdd = (param: any) => {
  addEditModalForm.value.id = ''
  addEditModalForm.value.keyword = ''
  addEditModalForm.value.isRegex = false
  addEditModalForm.value.caseInsensitive = false
  addEditModalShow.value = true
}

const showModalEdit = (param: FilterDatum) => {
  addEditModalForm.value.id = param.id
  addEditModalForm.value.keyword = param.keyword
  addEditModalForm.value.isRegex = param.is_regex
  addEditModalForm.value.caseInsensitive = param.case_insensitive
  addEditModalShow.value = !addEditModalShow.value
}

const showModalDelete = (param: FilterDatum) => {
  deleteModalId.value = param.id
  deleteModalShow.value = !deleteModalShow.value
}

watch(isDevEnvironment, (newValue, oldValue) => {
  FilterStore.fetchFilters()
})

onMounted(async () => {
  await FilterStore.fetchFilters()
})

const saveFilter = async () => {
  const prepFilter = {
    id: addEditModalForm.value.id,
    keyword: addEditModalForm.value.keyword,
    is_regex: addEditModalForm.value.isRegex,
    case_insensitive: addEditModalForm.value.caseInsensitive,
  } as FilterDatum

  await FilterStore.putFilter(prepFilter)

  addEditModalShow.value = !addEditModalShow.value
}

const deleteFilter = async () => {
  await FilterStore.deleteFilter(deleteModalId.value as string)

  deleteModalShow.value = !deleteModalShow.value
}
</script>

<template>
  <h1 class="page-title font-bold">{{ t('admin.filters') }}</h1>
  <VaButton size="small" icon="add" style="margin-bottom: 1rem" color="gsgYellow" round @click="showModalAdd">
    Add filter
  </VaButton>
  <div class="va-table-responsive">
    <table class="va-table va-table--hoverable">
      <thead>
        <tr>
          <th>Filter Keyword</th>
          <th>Regex</th>
          <th>Case Insensitive</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="filter in FilterStore.$state.filters" :key="filter.id">
          <td>
            <pre>{{ filter.keyword }}</pre>
          </td>
          <td>
            <VaIcon v-if="filter.is_regex" name="check" color="#00ff00"></VaIcon>
            <VaIcon v-else name="close" color="#ff0000"></VaIcon>
          </td>
          <td>
            <span v-if="!filter.is_regex">
              <VaIcon v-if="filter.case_insensitive" name="check" color="#00ff00"></VaIcon>
              <VaIcon v-else name="close" color="#ff0000"></VaIcon>
            </span>
            <span v-else> regex-determined </span>
          </td>
          <td>
            <div class="flex gap-2">
              <VaButton size="small" round icon="edit" color="gsgYellow" @click="showModalEdit(filter)"> Edit</VaButton>
              <VaButton size="small" round icon="delete" color="gsgRed" @click="showModalDelete(filter)"
                >Delete
              </VaButton>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
  <VaModal v-model="addEditModalShow" hide-default-actions blur>
    <template #header>
      <h4 v-if="addEditModalForm.id" class="va-h5">Edit Filter</h4>
      <h4 v-else class="va-h5">Add Filter</h4>
    </template>

    <div class="flex flex-col items-start gap-2">
      <VaForm ref="addForm" immediate hide-error-messages class="flex flex-col gap-2 mb-2">
        <VaInput
          v-model="addEditModalForm.keyword"
          :label="addEditModalForm.isRegex ? 'Regex' : 'Keyword'"
          name="FilterKeyword"
          :rules="[(v) => Boolean(v) || 'Keyword required']"
        />
        <VaCheckbox v-model="addEditModalForm.isRegex" label="Regex" name="IsRegex"></VaCheckbox>
        <VaCheckbox
          v-model="addEditModalForm.caseInsensitive"
          label="Case-Insensitive"
          name="CaseInsensitive"
          :disabled="addEditModalForm.isRegex"
        ></VaCheckbox>
      </VaForm>
    </div>

    <template #footer>
      <div class="flex gap-2">
        <VaButton color="gsgYellow" @click="addEditModalShow = !addEditModalShow">Cancel</VaButton>
        <VaButton color="gsgYellow" @click="saveFilter">Save</VaButton>
      </div>
    </template>
  </VaModal>
  <VaModal v-model="deleteModalShow" hide-default-actions blur ok-text="Apply">
    <h4 class="va-h5">Confirm Deletion</h4>
    <p>Are you sure you'd like to delete this filter?</p>
    <template #footer>
      <div class="flex gap-2">
        <VaButton color="gsgYellow" @click="deleteModalShow = !deleteModalShow">Cancel</VaButton>
        <VaButton color="gsgRed" @click="deleteFilter">Delete</VaButton>
      </div>
    </template>
  </VaModal>
</template>

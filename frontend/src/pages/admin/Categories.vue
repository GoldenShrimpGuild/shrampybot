<script lang="ts" setup>
import { watch, onMounted, computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { storeToRefs } from 'pinia'
import { useGlobalStore } from '../../stores/global-store'
import { useCategoryStore } from '../../stores/categories'
import { forEach } from 'lodash'

// components
import { VaModal, VaButton, VaForm, VaInput, VaTextarea, VaChip } from 'vuestic-ui'

// types
import type { CategoryDatum } from '../../../model/utility/nosqldb'
import type { AxiosResponse } from 'axios'

const { t } = useI18n()

const addEditModalShow = ref(false)
const deleteModalShow = ref(false)
const addEditModalForm = ref({
  id: '' as string | undefined,
  twitch_category: '',
  mastodon_tag_entry: '',
  mastodon_tags: [] as string[] | undefined,
  bluesky_tag_entry: '',
  bluesky_tags: [] as string[] | undefined,
})
const deleteModalId = ref('' as string | undefined)
const chipref = ref(null)

const route = useRoute()

const GlobalStore = useGlobalStore()
const CategoryStore = useCategoryStore()
const { isDevEnvironment } = storeToRefs(GlobalStore)
const { categories } = storeToRefs(CategoryStore)

const showModalAdd = (param: any) => {
  addEditModalForm.value.id = ''
  addEditModalForm.value.twitch_category = ''
  addEditModalForm.value.mastodon_tag_entry = ''
  addEditModalForm.value.mastodon_tags = []
  addEditModalForm.value.bluesky_tag_entry = ''
  addEditModalForm.value.bluesky_tags = []
  addEditModalShow.value = !addEditModalShow.value
}

const showModalEdit = (param: CategoryDatum) => {
  addEditModalForm.value.id = param.id
  addEditModalForm.value.twitch_category = param.twitch_category
  addEditModalForm.value.mastodon_tag_entry = param.mastodon_tags ? param.mastodon_tags.join(' ') : ''
  addEditModalForm.value.mastodon_tags = param.mastodon_tags
  addEditModalForm.value.bluesky_tag_entry = param.bluesky_tags ? param.bluesky_tags.join(' ') : ''
  addEditModalForm.value.bluesky_tags = param.bluesky_tags
  addEditModalShow.value = !addEditModalShow.value
}

const showModalDelete = (param: CategoryDatum) => {
  deleteModalId.value = param.id
  deleteModalShow.value = !deleteModalShow.value
}

watch(chipref, (v) => {
  console.log(chipref.value)
})

watch(isDevEnvironment, (newValue, oldValue) => {
  CategoryStore.fetchCategories()
})

onMounted(async () => {
  await CategoryStore.fetchCategories()
})

const updateTags = () => {
  var mtags = addEditModalForm.value.mastodon_tag_entry.split(' ')
  addEditModalForm.value.mastodon_tags = []
  forEach(mtags, (item) => {
    if (item.trim() != '' && (item.startsWith('#') || item.startsWith('@'))) {
      addEditModalForm.value.mastodon_tags?.push(item.trim())
    }
  })
  var mtags = addEditModalForm.value.bluesky_tag_entry.split(' ')
  addEditModalForm.value.bluesky_tags = []
  forEach(mtags, (item) => {
    if (item.trim() != '' && item.startsWith('#')) {
      addEditModalForm.value.bluesky_tags?.push(item.trim())
    }
  })
}

const saveCategory = async () => {
  const prepCategory = {
    id: addEditModalForm.value.id,
    twitch_category: addEditModalForm.value.twitch_category,
    mastodon_tags: addEditModalForm.value.mastodon_tags,
    bluesky_tags: addEditModalForm.value.bluesky_tags,
  } as CategoryDatum

  await CategoryStore.putCategory(prepCategory)

  addEditModalShow.value = !addEditModalShow.value
}

const deleteCategory = async () => {
  await CategoryStore.deleteCategory(deleteModalId.value as string)

  deleteModalShow.value = !deleteModalShow.value
}

const sortedCategories = computed(() => {
  if (categories) {
    const sortedCopy = [...CategoryStore.$state.categories]
    function compare(a: CategoryDatum, b: CategoryDatum) {
      if (a.twitch_category < b.twitch_category) return -1
      if (a.twitch_category > b.twitch_category) return 1
      return 0
    }
    return sortedCopy.sort(compare)
  } else {
    return []
  }
})
</script>

<template>
  <h1 class="page-title font-bold">{{ t('admin.categories') }}</h1>
  <VaButton size="small" icon="add" style="margin-bottom: 1rem" color="gsgYellow" round @click="showModalAdd">
    {{ t('admin.category.addCategory') }}
  </VaButton>
  <div class="va-table-responsive">
    <table class="va-table va-table--hoverable">
      <thead>
        <tr>
          <th>{{ t('admin.category.twitchCategoryName') }}</th>
          <th>{{ t('admin.category.mastodonTags') }}</th>
          <th>{{ t('admin.category.blueskyTags') }}</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="category in sortedCategories" :key="category.id">
          <td>
            <span style="font-weight: bold">{{ category.twitch_category }}</span>
          </td>
          <td>
            <ul>
              <li v-for="tag in category.mastodon_tags">
                {{ tag }}
              </li>
            </ul>
          </td>
          <td>
            <ul>
              <li v-for="tag in category.bluesky_tags">
                {{ tag }}
              </li>
            </ul>
          </td>
          <td>
            <div class="flex gap-2">
              <VaButton size="small" round icon="edit" color="gsgYellow" @click="showModalEdit(category)">
                {{ t('admin.category.edit') }}</VaButton>
              <VaButton size="small" round icon="delete" color="gsgRed" @click="showModalDelete(category)">{{
                t('admin.category.delete')
                }}</VaButton>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
  <VaModal v-model="deleteModalShow" hide-default-actions blur ok-text="Apply">
    <h4 class="va-h5">
      {{ t('admin.category.confirmDelete') }}
    </h4>
    <p>
      {{ t('admin.category.areYouSure') }}
    </p>
    <template #footer>
      <div class="flex gap-2">
        <VaButton color="gsgYellow" @click="deleteModalShow = !deleteModalShow">{{
          t('admin.category.cancel')
          }}</VaButton>
        <VaButton color="gsgRed" @click="deleteCategory">{{ t('admin.category.delete') }}</VaButton>
      </div>
    </template>
  </VaModal>
  <VaModal v-model="addEditModalShow" hide-default-actions blur>
    <template #header>
      <h4 v-if="addEditModalForm.id" class="va-h5">
        {{ t('admin.category.editCategory') }}
      </h4>
      <h4 v-else class="va-h5">
        {{ t('admin.category.addCategory') }}
      </h4>
    </template>

    <div class="flex flex-col items-start gap-2">
      <VaForm ref="addForm" immediate hide-error-messages class="flex flex-col gap-2 mb-2">
        <VaInput v-model="addEditModalForm.twitch_category" :label="t('admin.category.twitchCategory')"
          name="TwitchCategory" :rules="[(v) => Boolean(v) || t('admin.category.twitchCategoryRequired')]" />
        <VaTextarea v-model="addEditModalForm.mastodon_tag_entry" :label="t('admin.category.mastodonTags')"
          name="MastodonTags" :placeholder="t('admin.category.spaceSeparatedList')" @keyup="updateTags" />
        <div class="flex gap-2">
          <VaChip v-for="tag in addEditModalForm.mastodon_tags" color="mastodonLight" size="small">{{ tag }} </VaChip>
        </div>
        <VaTextarea v-model="addEditModalForm.bluesky_tag_entry" :label="t('admin.category.blueskyTags')"
          name="BlueskyTags" :placeholder="t('admin.category.spaceSeparatedList')" @keyup="updateTags" />
        <div class="flex gap-2">
          <VaChip v-for="tag in addEditModalForm.bluesky_tags" color="blueskyBlue" size="small">{{ tag }} </VaChip>
        </div>
      </VaForm>
    </div>

    <template #footer>
      <div class="flex gap-2">
        <VaButton color="gsgYellow" @click="addEditModalShow = !addEditModalShow">{{
          t('admin.category.cancel')
          }}</VaButton>
        <VaButton color="gsgYellow" @click="saveCategory">{{ t('admin.category.save') }}</VaButton>
      </div>
    </template>
  </VaModal>
</template>

<script lang="ts" setup>
import { watch, onMounted, computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { storeToRefs } from 'pinia'
import { useGlobalStore } from '../../stores/global-store';
import { useCategoryStore } from '../../stores/categories';
import { CategoryDatum } from '../../../model/utility/nosqldb';
import { forEach, keys } from 'lodash';
import { AxiosResponse } from 'axios';

const { t } = useI18n()

const addEditModalShow = ref(false)
const deleteModalShow = ref(false)
const addEditModalForm = ref({
    id: "" as String | undefined,
    twitch_category: "",
    mastodon_tag_entry: "",
    mastodon_tags: [] as string[] | undefined,
    bluesky_tag_entry: "",
    bluesky_tags: [] as string[] | undefined
})
const deleteModalId = ref("" as string | undefined)
const chipref = ref(null)

const route = useRoute()

const GlobalStore = useGlobalStore()
const CategoryStore = useCategoryStore()
const { isDevEnvironment } = storeToRefs(GlobalStore)
const { categories } = storeToRefs(CategoryStore)

const showModalAdd = (param: any) => {
    addEditModalForm.value.id = ""
    addEditModalForm.value.twitch_category = ""
    addEditModalForm.value.mastodon_tag_entry = ""
    addEditModalForm.value.mastodon_tags = []
    addEditModalForm.value.bluesky_tag_entry = ""
    addEditModalForm.value.bluesky_tags = []
    addEditModalShow.value = !addEditModalShow.value
}

const showModalEdit = (param: CategoryDatum) => {
    addEditModalForm.value.id = param.id
    addEditModalForm.value.twitch_category = param.twitch_category
    addEditModalForm.value.mastodon_tag_entry = param.mastodon_tags ? param.mastodon_tags.join(' ') : ""
    addEditModalForm.value.mastodon_tags = param.mastodon_tags
    addEditModalForm.value.bluesky_tag_entry = param.bluesky_tags ? param.bluesky_tags.join(' ') : ""
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

watch(categories, (v) => {

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
        if (item.trim() != "" && (item.startsWith('#') || item.startsWith('@'))) {
            addEditModalForm.value.mastodon_tags?.push(item.trim())
        }
    })
    var mtags = addEditModalForm.value.bluesky_tag_entry.split(' ')
    addEditModalForm.value.bluesky_tags = []
    forEach(mtags, (item) => {
        if (item.trim() != "" && (item.startsWith('#'))) {
            addEditModalForm.value.bluesky_tags?.push(item.trim())
        }
    })
}

const saveCategory = async () => {
    var prepCategory = {
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
        const sortedCopy = [...CategoryStore.$state.categories];
        function compare(a: CategoryDatum, b: CategoryDatum) {
            if (a.twitch_category < b.twitch_category)
                return -1;
            if (a.twitch_category > b.twitch_category)
                return 1;
            return 0;
        }
        return sortedCopy.sort(compare);
    } else {
        return []
    }
});

</script>

<template>
    <h1 class="page-title font-bold">{{ t('admin.categories') }}</h1>
    <VaButton size="small" icon="add" style="margin-bottom: 1rem;" color="gsgYellow" round @click="showModalAdd">
        Add Category
    </VaButton>
    <div class="va-table-responsive">
        <table class="va-table va-table--hoverable">
            <thead>
                <tr>
                    <th>Twitch Category Name</th>
                    <th>Mastodon Tags</th>
                    <th>Bluesky Tags</th>
                    <th></th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="category in sortedCategories" :key="category.id">
                    <td><span style="font-weight: bold;">{{ category.twitch_category }}</span></td>
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
                                Edit</VaButton>
                            <VaButton size="small" round icon="delete" color="gsgRed"
                                @click="showModalDelete(category)">Delete</VaButton>
                        </div>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
    <VaModal v-model="deleteModalShow" hide-default-actions blur ok-text="Apply">
        <h3 class="va-h5">
            Delete confirmation
        </h3>
        <p>
            Are you sure you'd like to delete this category?
        </p>
        <template #footer>
            <div class="flex gap-2">
                <VaButton color="gsgYellow" @click="deleteModalShow = !deleteModalShow">Cancel</VaButton>
                <VaButton color="gsgRed" @click="deleteCategory">Delete</VaButton>
            </div>
        </template>
    </VaModal>
    <VaModal v-model="addEditModalShow" hide-default-actions blur>
        <template #header>
            <h3 v-if="addEditModalForm.id" class="va-h5">
                Edit Category
            </h3>
            <h3 v-else class="va-h5">
                Add Category
            </h3>
        </template>

        <div class="flex flex-col items-start gap-2">

            <VaForm ref="addForm" immediate hide-error-messages class="flex flex-col gap-2 mb-2">
                <VaInput v-model="addEditModalForm.twitch_category" label="Twitch Category" name="TwitchCategory"
                    :rules="[(v) => Boolean(v) || 'Twitch Category is required']" />
                <VaTextarea v-model="addEditModalForm.mastodon_tag_entry" label="Mastodon Tags" name="MastodonTags"
                    placeholder="space-separated list" @keyup="updateTags" />
                <div class="flex gap-2">
                    <VaChip color="mastodonLight" size="small" v-for="tag in addEditModalForm.mastodon_tags">{{ tag }}
                    </VaChip>
                </div>
                <VaTextarea v-model="addEditModalForm.bluesky_tag_entry" label="Bluesky Tags" name="BlueskyTags"
                    placeholder="space-separated list" @keyup="updateTags" />
                <div class="flex gap-2">
                    <VaChip color="blueskyBlue" size="small" v-for="tag in addEditModalForm.bluesky_tags">{{ tag }}
                    </VaChip>
                </div>
            </VaForm>
        </div>

        <template #footer>
            <div class="flex gap-2">
                <VaButton color="gsgYellow" @click="addEditModalShow = !addEditModalShow">Cancel</VaButton>
                <VaButton color="gsgYellow" @click="saveCategory">Save</VaButton>
            </div>
        </template>
    </VaModal>
</template>

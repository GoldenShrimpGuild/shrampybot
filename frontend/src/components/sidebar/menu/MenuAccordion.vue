<template>
  <VaAccordion v-model="accordionValue" class="sidebar-accordion va-sidebar__menu__inner" multiple>
    <VaCollapse v-for="(route, idx) in items" :key="idx">
      <template #header>
        <VaSidebarItem :active="isRouteActive(route)" :to="route.children ? undefined : { name: route.name }">
          <VaSidebarItemContent>
            <VaIcon :name="route.meta?.nav?.icon" class="va-sidebar-item__icon" />

            <VaSidebarItemTitle>
              {{ t(route.meta?.nav?.displayName as string) }}
            </VaSidebarItemTitle>

            <VaIcon v-if="route.children" :name="accordionValue[idx] ? 'fa-angle-up' : 'fa-angle-down'" />
          </VaSidebarItemContent>
        </VaSidebarItem>
      </template>
      <template v-for="(child, index) in route.children" :key="index">
        <VaSidebarItem :active="isRouteActive(child)" :to="{ name: child.name }">
          <VaSidebarItemContent>
            <div class="va-sidebar-item__icon" />

            <VaSidebarItemTitle>
              {{ t(child.meta?.nav?.displayName as string) }}
            </VaSidebarItemTitle>
          </VaSidebarItemContent>
        </VaSidebarItem>
      </template>
    </VaCollapse>
  </VaAccordion>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { INavigationRoute } from '../../../router'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()

const props = withDefaults(
  defineProps<{
    items?: INavigationRoute[]
  }>(),
  {
    items: () => [],
  },
)

const accordionValue = ref<boolean[]>([])

onMounted(() => {
  accordionValue.value = props.items.map((item) => isItemExpanded(item))
})

// function isGroup(item: INavigationRoute) {
//   return !!item.children
// }

function isRouteActive(item: INavigationRoute) {
  return item.name === useRoute().name
}

function isItemExpanded(item: INavigationRoute): boolean {
  if (!item.children) {
    return false
  }

  const isCurrentItemActive = isRouteActive(item)
  const isChildActive = !!item.children.find((child) => (child.children ? isItemExpanded(child) : isRouteActive(child)))

  return isCurrentItemActive || isChildActive
}
</script>

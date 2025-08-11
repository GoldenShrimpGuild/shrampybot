<template>
  <VaSidebar
    v-model="writableVisible"
    :width="sidebarWidth"
    :color="color"
    minimized-width="0"
    :class="GlobalStore.$state.isDevEnvironment ? 'va-sidebar-dev no-animate' : ''"
  >
    <VaAccordion v-model="value" multiple>
      <VaCollapse v-for="(route, index) in visibleNavRoutes(navRoutes)" :key="index">
        <template #header="{ value: isCollapsed }">
          <VaSidebarItem
            :to="route.children ? undefined : { name: route.name }"
            :active="routeHasActiveChild(route)"
            :active-color="activeColor"
            :text-color="textColor(route)"
            :aria-label="`${visibleNavRoutes(route.children).length > 0 ? 'Open category ' : 'Visit'} ${t(route.meta.nav.displayName)}`"
            role="button"
            hover-opacity="0.10"
          >
            <VaSidebarItemContent class="py-3 pr-2 pl-4">
              <VaIcon
                v-if="route.meta.nav.icon"
                aria-hidden="true"
                :name="route.meta.nav.icon"
                size="20px"
                :color="iconColor(route)"
              />
              <VaSidebarItemTitle class="flex justify-between items-center leading-5 font-semibold">
                {{ t(route.meta.nav.displayName) }}
                <VaIcon
                  v-if="visibleNavRoutes(route.children).length > 0"
                  :name="arrowDirection(isCollapsed)"
                  size="20px"
                />
              </VaSidebarItemTitle>
            </VaSidebarItemContent>
          </VaSidebarItem>
        </template>
        <template #body>
          <div v-for="(childRoute, index2) in visibleNavRoutes(route.children)" :key="index2">
            <VaSidebarItem
              v-if="!childRoute.meta.nav.hidden && !childRoute.meta.nav.disabled"
              :to="{ name: childRoute.name }"
              :active="isActiveChildRoute(childRoute)"
              :active-color="activeColor"
              :text-color="textColor(childRoute)"
              :aria-label="`Visit ${t(route.meta.nav.displayName)}`"
              hover-opacity="0.10"
            >
              <VaSidebarItemContent class="py-3 pr-2 pl-11">
                <VaSidebarItemTitle class="leading-5 font-semibold">
                  {{ t(childRoute.meta.nav.displayName) }}
                </VaSidebarItemTitle>
              </VaSidebarItemContent>
            </VaSidebarItem>
          </div>
        </template>
      </VaCollapse>
    </VaAccordion>
  </VaSidebar>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { navRoutes, INavigationRoute } from '../../router'
import { useGlobalStore } from '../../stores/global-store'
import { useUserStore } from '../../stores/user'
import { useColors } from 'vuestic-ui'
import { useI18n } from 'vue-i18n'

// components
import { VaSidebar, VaSidebarItem, VaAccordion, VaCollapse, VaSidebarItemContent, VaSidebarItemTitle, VaIcon } from 'vuestic-ui'

const GlobalStore = useGlobalStore()

const { getColor, colorToRgba } = useColors()

const route = useRoute()
const { t } = useI18n()

const props = defineProps({
  visible: {
    type: Boolean,
    default: true,
  },
  mobile: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:visible'])

const writableVisible = computed({
  get: () => props.visible,
  set: (v: boolean) => emit('update:visible', v),
})

const isActiveChildRoute = (child: INavigationRoute) => route.name === child.name

const routeHasActiveChild = (section: INavigationRoute) => {
  if (!section.children) {
    return route.path.endsWith(`${section.name}`)
  }

  return section.children.some(({ name }) => route.path.endsWith(`${name}`))
}

const value = ref<boolean[]>([])

const visibleNavRoutes = (items: INavigationRoute[] | undefined) => {
  const viz = [] as INavigationRoute[]

  if (items) {
    items.forEach((item) => {
      if (!item.meta.nav.hidden && !item.meta.nav.disabled) {
        if (item.meta.perms.requiresAuth) {
          if (UserStore.scopeMatch(item.meta.perms.requiresScopes)) {
            viz.push(item)
          }
        } else {
          viz.push(item)
        }
      }
    })
  }

  return viz
}

const setActiveExpand = () => (value.value = navRoutes.map((route: INavigationRoute) => routeHasActiveChild(route)))

const sidebarWidth = computed(() => (props.mobile ? '100vw' : '280px'))
const color = computed(() => getColor('background-secondary'))
const activeColor = computed(() => colorToRgba(getColor('focus'), 0.1))

const iconColor = (route: INavigationRoute) => (routeHasActiveChild(route) ? 'primary' : 'secondary')
const textColor = (route: INavigationRoute) => (routeHasActiveChild(route) ? 'primary' : 'textPrimary')
const arrowDirection = (state: boolean) => (state ? 'va-arrow-up' : 'va-arrow-down')

watch(() => route.fullPath, setActiveExpand, { immediate: true })

const UserStore = useUserStore()

const items = ref([] as Array<INavigationRoute>)
</script>

<style lang="scss">
.va-sidebar {
  &__menu {
    padding: 2rem 0;
  }

  &-item {
    &__icon {
      width: 1.5rem;
      height: 1.5rem;
      display: flex;
      justify-content: center;
      align-items: center;
    }
  }

  &__title {
    font-family: 'Revalia';
  }
}
.va-sidebar-dev {
  background-image: url('/construction.jpg');
  background-position-y: -64px;
}
</style>

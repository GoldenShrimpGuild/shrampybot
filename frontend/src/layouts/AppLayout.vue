<template>
  <VaLayout
    :top="{ fixed: true, order: 2 }"
    :left="{ fixed: true, absolute: breakpoints.mdDown, order: 1, overlay: breakpoints.mdDown && !isSidebarMinimized }"
    @leftOverlayClick="isSidebarMinimized = true"
  >
    <template #top>
      <AppNavbar :is-mobile="isMobile" />
    </template>

    <template #left>
      <Sidebar :minimized="isSidebarMinimized" :animated="!isMobile" :mobile="isMobile" />
    </template>

    <template #content>
      <div :class="{ minimized: isSidebarMinimized }" class="app-layout__sidebar-wrapper">
        <div v-if="isFullScreenSidebar" class="flex justify-end">
          <VaButton class="px-4 py-4" icon="md_close" preset="plain" @click="onCloseSidebarButtonClick" />
        </div>
      </div>
      <AppLayoutNavigation v-if="!isMobile" class="p-4" />
      <main class="p-4 pt-0">
        <article>
          <RouterView />
        </article>
      </main>
    </template>
  </VaLayout>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onBeforeMount, onMounted, ref, computed, watch, watchEffect } from 'vue'
import { storeToRefs } from 'pinia'
import { onBeforeRouteUpdate, useRouter } from 'vue-router'
import { useBreakpoint } from 'vuestic-ui'
import { useTimer } from 'vue-timer-hook'

import { useAuthStore } from '../stores/auth'
import { useGlobalStore } from '../stores/global-store'
import { useUserStore } from '../stores/user'

import AppLayoutNavigation from '../components/app-layout-navigation/AppLayoutNavigation.vue'
import AppNavbar from '../components/navbar/AppNavbar.vue'
import Sidebar from '../components/sidebar/Sidebar.vue'

const AuthStore = useAuthStore()
const GlobalStore = useGlobalStore()
const UserStore = useUserStore()

const router = useRouter()

const breakpoints = useBreakpoint()

const sidebarWidth = ref('16rem')
const sidebarMinimizedWidth = ref('')

const isMobile = ref(false)
const isTablet = ref(false)
const { isSidebarMinimized, isDevEnvironment } = storeToRefs(GlobalStore)

const onResize = () => {
  isSidebarMinimized.value = breakpoints.mdDown
  isMobile.value = breakpoints.smDown
  isTablet.value = breakpoints.mdDown
  sidebarMinimizedWidth.value = isMobile.value ? '0' : '4.5rem'
  sidebarWidth.value = isTablet.value ? '100%' : '16rem'
}

const time = Date.now()
const timer = useTimer(time)
// const heartbeatTimerRestart = async () => {
//   // Monitor the current state of the authorization and refresh if it needs it
//   await AuthStore.testAndRefreshToken()

//   const time = Date.now()
//   timer.restart(time + 60000)
// }

// onBeforeMount(() => {
//   watchEffect(async () => {
//     if (timer.isExpired.value) {
//       await heartbeatTimerRestart()
//     }
//   })
// })

watch(isDevEnvironment, async (newValue, oldValue) => {
  // router.go(0)
})

onMounted(async () => {
  // if (!UserStore.$state.self?.id) {
  //   await UserStore.fetchSelf()
  //   console.log(UserStore.$state.self)
  // }

  window.addEventListener('resize', onResize)
  onResize()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', onResize)
})

onBeforeRouteUpdate(() => {
  if (breakpoints.mdDown) {
    // Collapse sidebar after route change for Mobile
    isSidebarMinimized.value = true
  }
})

const isFullScreenSidebar = computed(() => isTablet.value && !isSidebarMinimized.value)

const onCloseSidebarButtonClick = () => {
  isSidebarMinimized.value = true
}
</script>

<style lang="scss" scoped>
// Prevent icon jump on animation
.va-sidebar {
  width: unset !important;
  min-width: unset !important;
}
</style>

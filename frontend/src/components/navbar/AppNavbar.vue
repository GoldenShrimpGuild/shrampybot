<template>
  <VaNavbar class="app-layout-navbar py-2 px-0">
    <template #left>
      <div class="left">
        <Transition v-if="isMobile" name="icon-fade" mode="out-in">
          <VaIcon color="primary" :name="isSidebarMinimized ? 'menu' : 'close'" size="24px" style="margin-top: 3px"
            @click="isSidebarMinimized = !isSidebarMinimized" />
        </Transition>
        <RouterLink to="/" aria-label="Visit home page">
          <div>
            <ShrampybotLogo />
          </div>
        </RouterLink>
      </div>
      <div v-if="!isMobile"  style="margin-left: 1.5rem;">
      <VaButtonToggle class="" v-model="isDevEnvironment"
        size="small"
        toggle-color="#997112"
        color="#ffbb22"
        round
        :options="[
          { label: 'Dev', value: true },
          { label: 'Prod', value: false },
        ]" />
        <VaSpacer></VaSpacer>
      </div>
    </template>
    <template #right>
      <AppNavbarActions class="app-navbar__actions" :is-mobile="isMobile" />
    </template>
  </VaNavbar>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useGlobalStore } from '../../stores/global-store'
import { useAuthStore } from '../../stores/auth'
import AppNavbarActions from './components/AppNavbarActions.vue'
import ShrampybotLogo from '../logos/ShrampybotLogo.vue'
import { VaSpacer } from 'vuestic-ui'
import { watch } from 'vue'
import router from '../../router'

defineProps({
  isMobile: { type: Boolean, default: false },
})

const AuthStore = useAuthStore()
const GlobalStore = useGlobalStore()
const { isSidebarMinimized, isDevEnvironment } = storeToRefs(GlobalStore)

watch(isDevEnvironment, async (newValue, oldValue) => {
  if (newValue === false) {
    AuthStore.testAndRefreshToken()
  }
})

</script>

<style lang="scss" scoped>
.va-navbar {
  z-index: 2;

  @media screen and (max-width: 950px) {
    .left {
      width: 100%;
    }

    .app-navbar__actions {
      display: flex;
      justify-content: space-between;
    }
  }
}

.left {
  display: flex;
  align-items: center;
  margin-left: 1rem;

  &>* {
    margin-right: 1rem;
  }

  &>*:last-child {
    margin-right: 0;
  }
}

.icon-fade-enter-active,
.icon-fade-leave-active {
  transition: transform 0.5s ease;
}

.icon-fade-enter,
.icon-fade-leave-to {
  transform: scale(0.5);
}
</style>

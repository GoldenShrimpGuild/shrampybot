<script lang="ts" setup>
import { useBreakpoint } from 'vuestic-ui'
import VuesticLogo from '../components/logos/VuesticLogo.vue'
import GSGShrimp from '../components/logos/GSGShrimp.vue'
import ShrampybotLogo from '../components/logos/ShrampybotLogo.vue'

import { useGlobalStore } from '../stores/global-store'

const GlobalStore = useGlobalStore()

const breakpoint = useBreakpoint()
</script>

<template>
  <VaLayout class="h-screen bg-[var(--va-background-secondary)]">
    <template #left>
      <RouterLink
        class="bg-primary h-full flex items-center justify-center"
        style="width: 35vw"
        to="/"
        aria-label="Visit homepage"
      >
        <GSGShrimp />
      </RouterLink>
    </template>
    <template #content>
      <div
        class="h-full flex items-center justify-center"
        :class="GlobalStore.$state.isDevEnvironment ? 'va-navbar-dev no-animate' : ''"
      >
        <main class="mx-auto max-w-[420px]">
          <h1 style="font-size: 10rem; margin-bottom: 2rem;">404</h1>
          <!-- <p>Page not found</p> -->
        </main>
      </div>
    </template>
  </VaLayout>
</template>

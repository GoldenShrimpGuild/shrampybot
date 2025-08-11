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
          <RouterView />
        </main>
      </div>
    </template>
  </VaLayout>
</template>

<script lang="ts" setup>
import { useBreakpoint } from 'vuestic-ui'
import { useGlobalStore } from '../stores/global-store'

// components
import { VaLayout } from 'vuestic-ui'
import { RouterLink } from 'vue-router'
import GSGShrimp from '../components/logos/GSGShrimp.vue'

const GlobalStore = useGlobalStore()

const breakpoint = useBreakpoint()
</script>

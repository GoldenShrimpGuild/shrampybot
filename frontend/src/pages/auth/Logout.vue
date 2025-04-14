<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '../../stores/auth'
import { useRouter } from 'vue-router'

const AuthStore = useAuthStore()
const router = useRouter()

onMounted(() => {
  AuthStore.callLogout()
  router.push('login')
})
</script>

<template></template>

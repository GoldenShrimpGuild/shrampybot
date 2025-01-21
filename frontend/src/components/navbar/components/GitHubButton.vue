<template>
  <VaButton
    preset="secondary"
    color="textPrimary"
    href="https://github.com/GoldenShrimpGuild/shrampybot"
    target="_blank"
    aria-label="ShrampyBot Github"
  >
    <VaIcon :component="VaIconGitHub" />
  </VaButton>
</template>

<script lang="ts" setup>
import VaIconGitHub from '../../icons/VaIconGitHub.vue'
</script>

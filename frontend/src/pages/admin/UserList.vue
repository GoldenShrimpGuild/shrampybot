<script lang="ts" setup>
import { watch, onMounted, ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { storeToRefs } from 'pinia'
import { useGlobalStore } from '../../stores/global-store'
import { useTwitchUsersStore } from '../../stores/twitch_users'
import { TwitchUserDatum } from '../../../model/utility/nosqldb'

// components
import { VaTabs, VaTab, VaSwitch, VaIcon, VaButton, VaBadge } from "vuestic-ui"
import Bluesky from '../../components/icons/BlueskyIcon.vue'
import MastodonIcon from '../../components/icons/MastodonIcon.vue'
import TwitchIcon from '../../components/icons/TwitchIcon.vue'
import VaIconDiscord from '../../components/icons/VaIconDiscord.vue'
import YoutubeIcon from '../../components/icons/YoutubeIcon.vue'
import SteamIcon from '../../components/icons/SteamIcon.vue'
import VaIconGitHub from '../../components/icons/VaIconGitHub.vue'

const { t } = useI18n()

const GlobalStore = useGlobalStore()
const TwitchUsersStore = useTwitchUsersStore()
const { isDevEnvironment } = storeToRefs(GlobalStore)

watch(isDevEnvironment, (newValue, oldValue) => {
  TwitchUsersStore.fetchUsers()
})

onMounted(async () => {
  await TwitchUsersStore.fetchUsers()
})

const listType = ref(0)
const showIds = ref(false)

const sortedUsers = computed(() => {
  if (TwitchUsersStore.$state.users) {
    const sortedCopy = [...TwitchUsersStore.$state.users]
    function compare(a: TwitchUserDatum, b: TwitchUserDatum) {
      if (a.login < b.login) return -1
      if (a.login > b.login) return 1
      return 0
    }
    return sortedCopy.sort(compare)
  } else {
    return []
  }
})
</script>

<template>
  <h1 class="page-title font-bold">{{ t('admin.user_list') }}</h1>
  <VaTabs v-model="listType" color="gsgYellow" style="margin-bottom: 1rem">
    <template #tabs>
      <VaTab v-for="tab in [t('admin.spreadsheet'), t('admin.social_map')]" :key="tab" color="gsgYellow">
        {{ tab }}
      </VaTab>
    </template>
  </VaTabs>
  <p v-if="listType == 0" style="margin-bottom: 1rem">{{ t('admin.ml_spreadsheet_comment') }}</p>
  <div class="flex gap-8 flex-wrap" style="margin-bottom: 1rem">
    <VaSwitch v-model="showIds" color="gsgYellow" true-inner-label="IDs On" false-inner-label="IDs Off" size="small" />
  </div>
  <div class="va-table-responsive">
    <table class="va-table va-table--hoverable">
      <thead>
        <tr>
          <th v-if="listType == 0">Artist Name & Location</th>
          <th>Twitch Channel</th>
          <th>Discord Username</th>
          <th v-if="listType == 0">Email</th>
          <th v-if="listType == 1">Team/Mast.</th>
          <th v-if="listType == 1">Mastodon</th>
          <th v-if="listType == 1">Bluesky</th>
          <th v-if="listType == 1">YouTube</th>
          <th v-if="listType == 1">Steam</th>
          <th v-if="listType == 1">Github</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="user in sortedUsers" :key="user.id">
          <td v-if="listType == 0">
            <span style="font-weight: bold">{{ user.shrampybot_artist_name }}</span
            ><br />
            {{ user.shrampybot_location }}
          </td>
          <td style="text-align: left">
            <VaButton
              :href="'https://twitch.tv/' + user.login"
              round
              size="small"
              color="#000000"
              border-color="twitchPurple"
              target="_blank"
            >
              <TwitchIcon style="margin-right: 0.2rem" />
              {{ user.display_name }}
            </VaButton>
            <div v-if="showIds">
              <VaBadge :text="user.id" color="twitchPurple" />
            </div>
          </td>
          <td>
            <VaButton
              v-if="user.discord_username"
              :href="'https://discord.com/users/' + user.discord_user_id"
              target="_blank"
              round
              size="small"
              color="#000000"
              border-color="discordBlurple"
            >
              <VaIconDiscord style="margin-right: 0.2rem" />
              {{ user.discord_username }}
            </VaButton>
            <div v-if="showIds">
              <VaBadge :text="user.discord_user_id" color="discordBlurple" />
            </div>
          </td>
          <td v-if="listType == 0">
            <VaButton
              v-if="user.shrampybot_email"
              :href="'mailto:' + user.shrampybot_email"
              target="_blank"
              size="small"
              round
              color="#000000"
              border-color="gsgYellow"
            >
              {{ user.shrampybot_email }}
            </VaButton>
          </td>
          <td v-if="listType == 1">
            <VaIcon v-if="user.shrampybot_active == true" name="check" color="#00ff00"></VaIcon>
            <!-- <VaChip v-if="user.broadcaster_type" outline size="small" color="gsgYellow">{{ user.broadcaster_type }}
            </VaChip> -->
          </td>
          <td v-if="listType == 1" style="border-left: 2px solid #333333">
            <VaButton
              v-if="user.mastodon_user_id"
              :href="'https://soc.gsg.live/@' + user.mastodon_user_id"
              round
              color="#000000"
              border-color="mastodonLight"
              target="_blank"
              size="small"
            >
              <MastodonIcon style="margin-right: 0.2rem" />
              {{ user.mastodon_user_id }}
            </VaButton>
          </td>
          <td v-if="listType == 1" style="border-left: 2px solid #333333">
            <VaButton
              v-if="user.bluesky_username"
              :href="'https://bsky.app/profile/' + user.bluesky_username"
              target="_blank"
              color="#000000"
              border-color="blueskyBlue"
              size="small"
              outline
              round
            >
              <Bluesky style="margin-right: 0.2rem" />
              {{ user.bluesky_username ? '@' + user.bluesky_username : '' }}
            </VaButton>
            <div v-if="showIds">
              <VaBadge :text="user.bluesky_user_id" color="blueskyBlue" />
            </div>
          </td>
          <td v-if="listType == 1" style="border-left: 2px solid #333333">
            <VaButton
              v-if="user.youtube_username"
              :href="'https://www.youtube.com/channel/' + user.youtube_user_id"
              target="_blank"
              size="small"
              round
              color="#000000"
              border-color="youtubeRed"
            >
              <YoutubeIcon style="margin-right: 0.2rem" />
              {{ user.youtube_username }}
            </VaButton>
            <div v-if="showIds">
              <VaBadge :text="user.youtube_user_id" color="youtubeRed" />
            </div>
          </td>
          <td v-if="listType == 1" style="border-left: 2px solid #333333">
            <VaButton
              v-if="user.steam_username"
              :href="'https://steamcommunity.com/id/' + user.steam_username"
              target="_blank"
              color="#000000"
              border-color="#ffffff"
              size="small"
              outline
              round
            >
              <SteamIcon style="margin-right: 0.2rem" />
              {{ user.steam_username }}
            </VaButton>
            <div v-if="showIds">
              <VaBadge :text="user.steam_user_id" color="steamBlue" />
            </div>
          </td>
          <td v-if="listType == 1" style="border-left: 2px solid #333333">
            <VaButton
              v-if="user.github_username"
              :href="'https://steamcommunity.com/id/' + user.steam_username"
              target="_blank"
              color="#000000"
              border-color="#ffffff"
              size="small"
              outline
              round
            >
              <VaIconGitHub style="margin-right: 0.2rem" />
              {{ user.github_username }}
            </VaButton>
            <div v-if="showIds">
              <VaBadge :text="user.github_user_id" color="#ffffff" />
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style lang="css">
pre {
  font-size: 10pt;
}

.va-badge {
  display: inline-block;
}
</style>

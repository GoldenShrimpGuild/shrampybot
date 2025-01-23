<script lang="ts" setup>
import { watch, onBeforeMount, ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { storeToRefs } from 'pinia'
import { useGlobalStore } from '../../stores/global-store';
import { useTwitchUsersStore } from '../../stores/twitch_users';
import { TwitchUserDatum } from '../../../model/utility/nosqldb';
import { VaIcon } from 'vuestic-ui';
import Bluesky from '../../components/icons/BlueskyIcon.vue';
import MastodonIcon from '../../components/icons/MastodonIcon.vue';
import TwitchIcon from '../../components/icons/TwitchIcon.vue';
import VaIconDiscord from '../../components/icons/VaIconDiscord.vue';
import YoutubeIcon from '../../components/icons/YoutubeIcon.vue';
import SteamIcon from '../../components/icons/SteamIcon.vue';
import VaIconGitHub from '../../components/icons/VaIconGitHub.vue';

const { t } = useI18n()

const GlobalStore = useGlobalStore()
const TwitchUsersStore = useTwitchUsersStore()
const { isDevEnvironment } = storeToRefs(GlobalStore)

watch(isDevEnvironment, async (newValue, oldValue) => {
  TwitchUsersStore.fetchUsers()
})

onBeforeMount(() => {
  TwitchUsersStore.fetchUsers()
})

const sortOrder = ref('asc');

const sortedUsers = computed(() => {
  const sortedCopy = [...TwitchUsersStore.$state.users];
  function compare(a: TwitchUserDatum, b: TwitchUserDatum) {
    if (a.login < b.login)
      return -1;
    if (a.login > b.login)
      return 1;
    return 0;
  }
  console.log(sortedCopy[0])
  return sortedCopy.sort(compare);
});

</script>

<template>
  <h1 class="page-title font-bold">{{ t('admin.social_map') }}</h1>
  <div class="va-table-responsive">
    <table class="va-table va-table--hoverable">
      <thead>
        <tr>
          <th>Discord ID</th>
          <th>Discord Username</th>
          <th>Twitch ID</th>
          <th>Twitch Channel</th>
          <th>Team/Mast.</th>
          <th>Mastodon</th>
          <th>Bluesky</th>
          <th>YouTube</th>
          <th>Steam</th>
          <th>Github</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="user in sortedUsers" :key="user.id">
          <td>
            <pre>{{ user.discord_user_id }}</pre>
          </td>
          <td>
            <VaButton v-if="user.discord_username" round size="small" color="#000000" borderColor="discordBlurple">
              <VaIconDiscord style="margin-right: 0.2rem;" />
              {{ user.discord_username }}
            </VaButton>
          </td>
          <td style="border-left: 2px solid #333333; text-align: right">
            <pre>{{ user.id }}</pre>
          </td>
          <td style="text-align: left">
            <VaButton :href="'https://twitch.tv/' + user.login" round size="small" color="#000000"
              borderColor="twitchPurple" target="_blank">
              <TwitchIcon style="margin-right: 0.2rem;" />
              {{ user.display_name }}
            </VaButton>
          </td>
          <td>
            <VaIcon v-if="user.shrampybot_active == true" name="check" color="#00ff00"></VaIcon>
            <!-- <VaChip v-if="user.broadcaster_type" outline size="small" color="gsgYellow">{{ user.broadcaster_type }}
            </VaChip> -->
          </td>
          <td style="border-left: 2px solid #333333">
            <VaButton v-if="user.mastodon_user_id" :href="'https://soc.gsg.live/@' + user.mastodon_user_id" round
              color="#000000" borderColor="mastodonLight" target="_blank" size="small">
              <MastodonIcon style="margin-right: 0.2rem" />
              {{ user.mastodon_user_id }}
            </VaButton>
          </td>
          <td style="border-left: 2px solid #333333">
            <VaButton v-if="user.bluesky_username" :href="'https://bsky.app/profile/' + user.bluesky_username"
              target="_blank" color="#000000" borderColor="blueskyBlue" size="small" outline round>
              <Bluesky style="margin-right: 0.2rem;" />
              {{ user.bluesky_username ? '@' + user.bluesky_username : '' }}
            </VaButton>
          </td>
          <td style="border-left: 2px solid #333333">
            <VaButton v-if="user.youtube_username" :href="'https://www.youtube.com/channel/' + user.youtube_user_id"
              target="_blank" size="small" round color="#000000" borderColor="youtubeRed">
              <YoutubeIcon style="margin-right: 0.2rem;" />
              {{ user.youtube_username }}
            </VaButton>
          </td>
          <td style="border-left: 2px solid #333333">
            <VaButton v-if="user.steam_username" :href="'https://steamcommunity.com/id/' + user.steam_username"
              target="_blank" color="#000000" borderColor="#ffffff" size="small" outline round>
              <SteamIcon style="margin-right: 0.2rem;" />
              {{ user.steam_username }}
            </VaButton>
          </td>
          <td style="border-left: 2px solid #333333">
            <VaButton v-if="user.github_username" :href="'https://steamcommunity.com/id/' + user.steam_username"
              target="_blank" color="#000000" borderColor="#ffffff" size="small" outline round>
              <VaIconGitHub style="margin-right: 0.2rem;" />
              {{ user.github_username }}
            </VaButton>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

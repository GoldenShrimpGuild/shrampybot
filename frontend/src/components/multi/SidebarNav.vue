<template>
  <VaAppBar
    color="secondary"
    :class="sidebarButtonsClass"
    :style="`height: ${height}px; width: ${width}px;`"
  >
    <VaButton
      color="textInverted"
      size="medium"
      preset="plain"
      :icon="hideDecor ? 'check_box_outline_blank' : 'check_box'"
      v-on:click="emit('toggleDecor')"
    >
      {{ t("multiTwitch.decor") }}
    </VaButton>
    <VaButton
      class="ml-2"
      color="textInverted"
      size="medium"
      preset="plain"
      :icon="hideChat ? 'check_box_outline_blank' : 'check_box'"
      v-on:click="emit('toggleChat');"
    >
      {{ t("multiTwitch.chat") }}
    </VaButton>
    <VaSpacer />
    <VaMenu
      trigger="hover"
      :options="sideMenu"
      @selected="(v) => v.handler()"
      disabled-by="disabled"
    >
      <template #anchor>
        <GSGShrimpIcon
          :class="shrimpIconClass"
          body-color="#000000"
        >
        </GSGShrimpIcon>
      </template>
    </VaMenu>
  </VaAppBar>
  <div :class="`${sidebarButtonsClass} grid-cols-3 grid hidden`">
    <VaButton
      round
      class="item mr-1"
      size="medium"
      color="gsgYellow"
      :icon="hideDecor ? 'check_box_outline_blank' : 'check_box'"
      v-on:click="emit('toggleDecor')"
    >
      {{ t("multiTwitch.decor") }}
    </VaButton>
    <VaButton
      round
      class="item mr-1"
      size="medium"
      color="gsgYellow"
      :icon="hideChat ? 'check_box_outline_blank' : 'check_box'"
      v-on:click="emit('toggleChat');"
    >
      {{ t("multiTwitch.chat") }}
    </VaButton>
    <div>
      <VaMenu
        :options="sideMenu"
        @selected="(v) => v.handler()"
        disabled-by="disabled"
      >
        <template #anchor>
          <GSGShrimpIcon :class="shrimpIconClass">
          </GSGShrimpIcon>
        </template>
      </VaMenu>
    </div>
  </div>
</template>

<style scoped>
.shrimp-icon {
  position: absolute;
  top: 0px;
  right: 0px;
  margin-right: 3px;
  height: v-bind(height + 'px');
  width: v-bind(height + 'px');
}
.sidebar-buttons {
  padding: 0px;
  transition-duration: 0;
  transition-property: none;
}
</style>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'

// Components
import GSGShrimpIcon from '../icons/GSGShrimpIcon.vue'

const { t } = useI18n()

const shrimpIconClass = ref("shrimp-icon")
const sidebarButtonsClass = ref("sidebar-buttons")

const props = defineProps<{
  height: number,
  width: number,
  hideDecor: boolean,
  hideChat: boolean,
  centreFrames: boolean,
  includeGSGChannel: boolean,
  disableStreamLoading: boolean,
  startMuted: boolean,
  useDevAPI: boolean,
  onlyMusicCategory: boolean,
  useCurrentEventData: boolean,
}>()

const emit = defineEmits<{
  (e: "toggleChat"): void,
  (e: "toggleDecor"): void,
  (e: "toggleCentre"): void,
  (e: "toggleIncludeGSG"): void,
  (e: "toggleFakeData"): void,
  (e: "toggleDevAPI"): void,
  (e: "toggleStartMuted"): void,
  (e: "handleUrl", url: string): void,
  (e: "toggleCurrentEventData"): void,
  (e: "toggleOnlyMusicCategory"): void,
  (e: "adjustFilter"): void,
}>()

const sideMenu = computed(() => [
    {
        text: t("multiTwitch.preferences"),
        value: "",
        disabled: true,
        handler: (o: object) => {},
    },
    {
        text: t("multiTwitch.adjustFilters"),
        value: "adjustFilter",
        disabled: false,
        icon: 'filter_list',
        handler: (o: object) => {
          emit("adjustFilter")
        },
    },
    {
        text: t("multiTwitch.centreAlignStreams"),
        value: "centreFramesToggle",
        disabled: false,
        icon: props.centreFrames ? 'check_box' : 'check_box_outline_blank',
        handler: (o: object) => {
            emit("toggleCentre")
        },
    },
    {
        text: t("multiTwitch.streamsStartMuted"),
        value: "startMuted",
        disabled: false,
        icon: props.startMuted ? 'check_box' : 'check_box_outline_blank',
        handler: (o: object) => {
            emit("toggleStartMuted")
        }
    },
    {
        text: t("multiTwitch.includeGSGChannel"),
        value: "includeGSGChannel",
        disabled: false,
        icon: props.includeGSGChannel ? 'check_box' : 'check_box_outline_blank',
        handler: (o: object) => {
            emit("toggleIncludeGSG")
        }
    },
    {
        text: t("multiTwitch.onlyMusicCategory"),
        value: "onlyMusicCategory",
        disabled: false,
        icon: props.onlyMusicCategory ? 'check_box' : 'check_box_outline_blank',
        handler: (o: object) => {
            emit("toggleOnlyMusicCategory")
        }
    },
    {
        text: t("multiTwitch.includeCurrentEvent"),
        value: "useCurrentEventData",
        disabled: false,
        icon: props.useCurrentEventData ? 'check_box' : 'check_box_outline_blank',
        handler: (o: object) => {
            emit("toggleCurrentEventData")
        }
    },
    {
        text: "",
        value: "",
        disabled: true,
        handler: (o: object) => {},
    },
    {
        text: t("multiTwitch.testing"),
        value: "",
        disabled: true,
        handler: (o: object) => {},
    },
    {
        text: t("multiTwitch.useDevAPI"),
        value: "useDevApi",
        disabled: false,
        icon: props.useDevAPI ? 'check_box' : 'check_box_outline_blank',
        handler: async (o: object) => {
          emit("toggleDevAPI")
        }
    },
    {
        text: t("multiTwitch.setFakeData"),
        value: "fakeDataMode",
        disabled: false,
        icon: props.disableStreamLoading ? 'check_box' : 'check_box_outline_blank',
        handler: async (o: object) => {
            emit("toggleFakeData")
        },
    },
    {
      text: "",
      value: "",
      disabled: true,
      handler: (o: object) => {},
    },
    {
      text: t("multiTwitch.aboutGSG"),
      value: "",
      disabled: true,
      handler: (o: object) => {},
    },
    {
      text: t("multiTwitch.gsgWebsite"),
      value: "gsgWebsite",
      disabled: false,
      icon: "link",
      handler: (o: object) => {
          emit("handleUrl", "https://www.gsg.live")
      },
    },
    {
      text: t("multiTwitch.gsgDiscord"),
      value: "gsgDiscord",
      disabled: false,
      icon: "link",
      handler: (o: object) => {
          emit("handleUrl", "https://discord.com/invite/Ahvc7ZjCUA")
      },
    },
    {
      text: t("multiTwitch.gsgChannel"),
      value: "gsgTwitchChannel",
      disabled: false,
      icon: "link",
      handler: (o: object) => {
          emit("handleUrl", "https://www.twitch.tv/goldenshrimpguild")
      },
    },
    {
      text: t("multiTwitch.gsgTwitchTeam"),
      value: "gsgTwitchTeam",
      disabled: false,
      icon: "link",
      handler: (o: object) => {
          emit("handleUrl", "https://www.twitch.tv/team/gsg")
      },
    },
])
</script>
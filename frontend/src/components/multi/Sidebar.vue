<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import TwitchChat from '../../components/multi/TwitchChat.vue'
import type { VaMenuOption } from 'vuestic-ui/dist/types/components/va-menu-list/types.js'
import type { StreamHistoryDatum as Stream } from '../../../model/utility/nosqldb'
import GSGShrimpIcon from '../icons/GSGShrimpIcon.vue'

const { t } = useI18n()

const sidebarClass = ref("mt-sidebar")
const sidebarButtonsClass = ref("sidebar-buttons")
const chatSelectorClass = ref("chat-selector")
const shrimpIconClass = ref("shrimp-icon")
const embedContainerClass = ref("embed-container")

const props = defineProps<{
  streamsList: Stream[],
  // userLogins: string[],
  // streamsLoaded: boolean,
  
  hideChat: boolean,
  hideDecor: boolean,
  dynamicSidebarWidth: number,
  dynamicSidebarHeight: number,
  useDropdownChatSelect: boolean,
  fixedButtonsHeight: number,
  dynamicChatEmbedHeight: number,
  fixedChatSelectorHeight: number,
  sideMenu: VaMenuOption[],
}>()

const emit = defineEmits<{
  (e: "toggleChat"): void,
  (e: "toggleDecor"): void,
  (e: "toggleCentre"): void,
}>()

const userLogins = computed(() => props.streamsList.map((v) => v.user_login))

const chatTabSelected = ref(0)
const currentChatStream = computed(() => props.streamsList.at(chatTabSelected.value))
const currentChatLogin = computed(() => props.streamsList.at(chatTabSelected.value)?.user_login)

</script>

<style scoped>
.mt-sidebar {
  height: v-bind(dynamicSidebarHeight);
  width: v-bind(dynamicSidebarWidth + 'px');
  float: right;
}
.sidebar-buttons {
  height: v-bind(fixedButtonsHeight + 'px');
  width: v-bind(dynamicSidebarWidth + 'px');
  padding: 0px;
  transition-duration: 0;
  transition-property: none;
  /* padding: 2px 1px 2px 2px; pl-2 pr-1 pt-2 pb-2 */
}
.shrimp-icon {
  position: absolute;
  top: 0px;
  right: 0px;
  margin-right: 3px;
  height: v-bind(fixedButtonsHeight + 'px');
  width: v-bind(fixedButtonsHeight + 'px');
}
.chat-selector {
  height: v-bind(fixedChatSelectorHeight + 'px');
  width: v-bind(dynamicSidebarWidth + 'px');
}
.chat-selector .va-tabs {
  background-color: #715411;
}
.chat-selector .va-tab * {
  margin-top: 3px;
}
.embed-container {
  height: v-bind(dynamicChatEmbedHeight + 'px');
}

</style>

<template>
  <div :class="`${sidebarClass}`">
    <VaAppBar
      color="secondary"
      :class="`${sidebarButtonsClass}`"
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
    <div
      :hidden="hideChat"
      :class="`${chatSelectorClass} h-full bg-[]`"
    >
      <VaTabs 
        v-if="!useDropdownChatSelect"
        v-model="chatTabSelected"
      >
        <template #tabs>
          <VaTab
            v-for="stream in Object.values(streamsList)"
            :key="stream.user_login"
            :label="stream.user_name"
          >
            {{ stream.user_name }}
          </VaTab>
        </template>
      </VaTabs>
      <VaSelect
        v-else
        :class="chatSelectorClass"
        v-model="chatTabSelected"
        placeholder="Colored"
        color="#FFFFFF"
        :options="streamsList"
        inner-label
      />
      <div :class="embedContainerClass">
        <VaSkeletonGroup
            v-if="!streamsList"
            :class="`h-full w-[${dynamicSidebarWidth}px]`"
            animation="wave"
            :delay="0"
          >
          <div
            :style="`width: ${dynamicSidebarWidth}px;`"
          >
              <VaSkeleton variant="text" class="ml-2 va-text" :lines="100" />
          </div>
        </VaSkeletonGroup>
        <template
          v-else
          v-for="stream in streamsList"
          :key="stream.user_login">
          <TwitchChat
            :hidden="currentChatLogin != stream.user_login"
            :user-login="stream.user_login"
            :height="dynamicChatEmbedHeight"
          ></TwitchChat>
        </template>
      </div>
    </div>
  </div>
</template>
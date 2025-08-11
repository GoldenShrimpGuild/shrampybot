<template>
  <div
    class="mt-wrapper"
    :style="`height: ${height}px; width: ${areaWidth}px;`">
    <div class="mt-header" :hidden="hideDecor">
      <h1 
        :class="`pl-3 pt-3 pb-0 font-bold`"
        :style="`height: ${fixedHeaderHeight}px;`"
      >
        {{ t('multiTwitch.header') }}
      </h1>
    </div>
    <div 
      :class="`mt-streams-container row`"
      :style="`width: ${areaWidth}px; height: ${areaHeight}px; padding: 0 0 0 0;`"
      :aria-busy="!streamsLoaded"
    >
      <!-- ${centreFrames ? 'place-content-center' : ''} -->
      <div
        :class="`card-flexbox flex flex-wrap`"
        :style="`width: ${areaWidth}px; height: ${areaHeight}px; gap: ${spacing}px; place-content: ${centreFrames ? 'center' : 'start'};`"
      >
          <VaCard 
            v-for="stream in streamsList"
            :key="stream.user_login"
            :class="`mt-card flex-none item`"
            :style="`
              height: ${cardHeight}px;
              width: ${cardWidth}px;
              background-color: ${stream.isEventStream ? getColor('gsgDarkYellow') : getColor('backgroundSecondary')}`"
          >
            <VaSkeleton
              v-if="forceSkeleton"
              animation="wave"
              :style="`height: ${cardHeight-decorationHeight}px;`"
              :color="currentRainbowColour"
              variant="squared"
            />
            <TwitchStream
              v-else
              :user-login="stream.user_login"
              :width="cardWidth"
              :start-muted="startMuted"
              :aspect-ratio="areaRatio"
              @start-embed=""
            ></TwitchStream>
            <VaCardContent
              :hidden="hideDecor"
            >
              <h5 :style="`height: 1.2rem; width: ${textWidth}px; text-overflow: ellipsis; overflow: hidden; white-space: nowrap;`">{{ stream.user_name }}</h5>
              <p :style="`height: 1.2rem; width: ${textWidth}px; text-overflow: ellipsis; overflow: hidden; white-space: nowrap;`">
                {{ stream.title }}
              </p>
            </VaCardContent>
          </VaCard>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useColors } from 'vuestic-ui'
import { useI18n } from 'vue-i18n'

// Types
import type { GSGStream as Stream } from '../../stores/public/classes'

// Components
import TwitchStream from '../../components/multi/TwitchStream.vue'

const { t } = useI18n()

const { getColor } = useColors()

const fixedHeaderHeight = 50
const fixedDecorationHeight = 65
const fixedSpacingAmt = 6
const streamRatio = 16/9

const props = defineProps<{
  height: number,
  width: number,
  hideDecor: boolean,
  hideChat: boolean,
  forceSkeleton: boolean,
  streamsLoaded: boolean,
  centreFrames: boolean,
  startMuted: boolean,
  streamsList: Stream[],
  currentRainbowColour: string,
}>()

const spacing = ref(fixedSpacingAmt)
const headerHeight = computed(() => props.hideDecor ? 0 : fixedHeaderHeight)
const areaHeight = computed(() => props.height - headerHeight.value*2)
const areaWidth = computed(() => props.hideChat ? props.width : props.width)
const areaRatio = computed(() => areaWidth.value / areaHeight.value)
const decorationHeight = computed(() => props.hideDecor ? 0 : fixedDecorationHeight)

const cardCount = ref(0)
const cardCols = ref(1)
const cardRows = ref(1)
const cardHeight = ref(0)
const cardWidth = ref(0)
const textWidth = ref(0)
const topPadding = ref(0)

const updateCardLayout = () => {
  cardCount.value = props.streamsList.length

  var bestHeight = 0;
  var bestWidth = 0;
  var finalRows = 1;
  var finalCols = 1;
  var topWrapperPadding = 0;

  for (var cols = 1; cols <= cardCount.value; cols++) {
      const rows = Math.ceil(cardCount.value / cols);
      var maxWidth = Math.floor(areaWidth.value / cols) - spacing.value;
      var maxHeight = Math.floor(areaHeight.value / rows) - spacing.value - decorationHeight.value;
      if (maxWidth * 1/streamRatio < maxHeight) {
          maxHeight = maxWidth * 1/streamRatio;
      } else {
          maxWidth = (maxHeight) * streamRatio;
      }
      if (maxWidth > bestWidth) {
          bestWidth = maxWidth;
          bestHeight = maxHeight + decorationHeight.value;
          topWrapperPadding = (areaHeight.value - (rows * (maxHeight + spacing.value + decorationHeight.value)))/2;
          finalRows = rows;
          finalCols = cols;
      }
  }

  cardCols.value = finalCols
  cardRows.value = finalRows
  cardHeight.value = bestHeight
  cardWidth.value = bestWidth
  textWidth.value = Math.floor(cardWidth.value * 0.85)
  topPadding.value = topWrapperPadding
}

// Monitor prop as computed ref
const streamsRef = computed(() => props.streamsList)
watch(streamsRef, () => {
  updateCardLayout()
})

watch(areaWidth, () => {
  updateCardLayout()
})

watch(areaHeight, () => {
  updateCardLayout()
})

watch(decorationHeight, () => {
  updateCardLayout()
})

onMounted(() => {
  updateCardLayout()
})

</script>
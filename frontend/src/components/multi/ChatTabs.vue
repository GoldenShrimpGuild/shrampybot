<script setup lang="ts">
import { ref, useTemplateRef, computed, watch, ComponentPublicInstance } from 'vue'
import { helper } from '../../plugins/helpers'
import { useColors } from "vuestic-ui"
import { clamp, useElementSize, useMutationObserver, useParentElement } from '@vueuse/core'

// Types
import { ComputedRef, VueElement, ShallowRef, Ref } from 'vue'
import type { StreamHistoryDatum as Stream } from '../../../model/utility/nosqldb'

const { getColor } = useColors()

const props = defineProps<{
  streamsList: Stream[],
  height: string | number,
  width: string | number,
  backgroundColor: string,
  hidePagination?: boolean,
  forcedSelection?: 0,
}>()

const emit = defineEmits<{
  (e: "select", tab: number): void,
}>()

const fixedTabGraceArea = 250

const tabWrapperClass = "va-tabs__tabs"

// Template refs
const chatTabsRoot = useTemplateRef("chat-tabs")
const chatTab = useTemplateRef("chat-tab")
const pagButtonLeft = useTemplateRef("tab-pag-button-left")
const pagButtonRight = useTemplateRef("tab-pag-button-right")

// Computed refs from props
const bgColor = computed(() => getColor(props.backgroundColor, "#000000"))
const selectorHeight = computed(() => helper.parsePixels(props.height))
const selectorWidth = computed(() => helper.parsePixels(props.width))
const pagButtonWidth = useElementSize(pagButtonLeft || pagButtonRight)?.width as ShallowRef<number>
const tabSelectorWidth = computed(() => selectorWidth.value - pagButtonWidth.value*2)

const chatTabList = ref([] as ComponentPublicInstance[])
const chatTabWidths = computed(() => chatTabList.value.map((v) => v ? useElementSize(v.$el).width.value : 0))
const chatTabSum = computed(() => chatTabWidths.value.reduce((p, v) => p += v, 0) as number)
const wrapperMax = computed(() => chatTabSum ? Math.floor(Math.max(0, chatTabSum.value - fixedTabGraceArea)) : 0)
const tabWrapper: ComputedRef<VueElement | null> = computed(() => {
  if (chatTabList.value) {
    const first = chatTabList.value.at(0)
    const tt = first ? useParentElement(useParentElement(first.$el)) : null
    if (tt && tt.value) return tt.value as VueElement
  }
  return {} as VueElement
})

// Retrieve and stow for transform to tab selector
const transformRegex = /translateX\(?([0-9\.\-e]+)?px\)/i
const transformSubst = "translateX({VALUE}px)"
const pagX = ref(0) // stored as a positive integer for ease
const pagIncSize = computed(() => chatTabSum.value / chatTabList.value.length)

const updateChatTabs = () => {
  chatTabList.value = chatTab.value ? chatTab.value.values().toArray() as ComponentPublicInstance[] : []
}

// Tab selection logic, initial value to -1 (nothing selected)
const tabSelected = ref(-1)
const getExactSelectedLeft = computed(() => {
  if (tabSelected.value < 0) return 0
  return chatTabWidths.value ? chatTabWidths.value.slice(0, Math.max(0, tabSelected.value)).reduce((p, v) => p += v, 0) : 0
})

// Direction is -1 or 1, 0 for realignment
const paginate = (direction: number) => {
  updateChatTabs()
  if (!tabWrapper || !tabWrapper.value) return

  // Lock direction between -1 and 1
  const dirLock = Math.ceil(clamp(direction, -1, 1))

  // if (dirLock === 0) {
  //   if (tabSelected.value === -1) {
  //     pagX.value = 0
  //     return
  //   }

  //   // // if (!(exactSelectedLeft.value >= pagX.value && exactSelectedLeft.value < pagX.value + fixedTabGraceArea)) {
  //   //   pagX.value = exactSelectedLeft.value
  //   // // }
  //   // return
  // }
  const size = dirLock * (dirLock > 0
    // Scrolling to the right
    ? Math.max(pagIncSize.value, pagX.value)
    // Scrolling to the left
    : Math.min(pagIncSize.value, pagX.value)
  )

  pagX.value = Math.floor(clamp(pagX.value + size, 0, wrapperMax.value))
}

// Mapping streams to computed for reactivity in watch
const streamsRef = computed(() => props.streamsList)

// Means of keeping tabPagX fully up to date.
var pausePagXObserver = ref(false)
useMutationObserver(
  chatTabsRoot, 
  (mutations: any) => {
    mutations.forEach((m: any) => {
      // Watch specifically for style changes (we're looking for transform)
      if (m.type == "attributes" && m.attributeName == "style") {
        const target = m.target as VueElement
        if (target.nodeType === 1 && target.nodeName === "DIV" && target.classList.contains(tabWrapperClass)) {
          // Parse out the number, if any, in the transform
          const newValue = Math.abs(parseInt(target.style.transform.replace(transformRegex, "$1")))
          if (!newValue) return

          // Only write this value if 
          if (!pausePagXObserver && newValue !== pagX.value) {
            pagX.value = Math.abs(newValue)
          } else {
            pausePagXObserver.value = false
          }
        }
      }
    })
  }, 
  {attributes: true, childList: true, subtree: true}
)

// Pay attention if the selected tab changes
watch(tabSelected, (v, oV) => {
  emit("select", v)
})

watch(streamsRef, (v, oV) => {
  updateChatTabs()

  // Empty list? Reset tabSelected to -1 and exit
  if (!v.length) {
    tabSelected.value = -1
    return
  }

  // If no tabs are selected, default to the first one
  if (tabSelected.value === -1) {
    tabSelected.value = 0
    return
  }

  // Check if login index changed (stream in the middle went offline)
  const origLoginI = oV.at(tabSelected.value)
  const matchedLoginI = v.findIndex((s) => origLoginI ? s.user_login === origLoginI.user_login : false)
  if (matchedLoginI !== -1) {
    // Set tabSelected value to follow the matchedLogin
    tabSelected.value = matchedLoginI
    return
  }

  // If zero at this point, just stay put
  if (tabSelected.value === 0) {
    return
  }

  // Handle new list being shorter than the old list if selected value is beyond the limit
  // (sticky to the last possible value)
  if (tabSelected.value > v.length - 1) {
    tabSelected.value = v.length - 1
    return
  }
})

// Means of auto-updating the tab paginator when pagX changes
watch(pagX, (v, oV) => {
  if (!tabWrapper || !tabWrapper.value) return
  // // Only update the paginator if the value differs (prevent possible loops?)
  const transform = Math.abs(parseInt((tabWrapper.value.style.transform).replace(transformRegex, "$1")))
  if (v !== transform) {
    // prevent the PagXObserver from registering this change
    // the mutationObserver will turn it on again after skipping the next mutation
    // Note: if you decide to 
    pausePagXObserver.value = true

    const subst = transformSubst.replace(/{VALUE}/, `${-v}`)
    tabWrapper.value.style.setProperty('transform', subst)
  }
})

</script>

<template>
  <div 
    class="flex flex-cols gap-0"
    ref="chat-tabs"
    :width="`width: ${selectorWidth}px; height: ${selectorHeight}px`"
  >
    <VaButton
      size="small"
      ref="tab-pag-button-left"
      color="gsgDarkYellow"
      icon="va-arrow-left"
      :hidden="hidePagination"
      :style="`width: ${pagButtonWidth}px; height: ${selectorHeight}px`"
      @click="paginate(-1)"
      @keydown-enter="paginate(-1)"
    />
    <VaTabs 
      v-model="tabSelected"
      :hide-pagination="true"
      :style="`padding-top: 3px; background-color: ${bgColor}; width: ${tabSelectorWidth}px; height: ${selectorHeight}px`"
    >
      <template #tabs>
        <VaTab
          ref="chat-tab"
          v-for="stream, i in streamsList"
          :key="stream.user_login"
          :label="stream.user_name"
          @click="paginate(0)"
          @keydown-enter="paginate(0)"
        >
          {{ stream.user_name }}
        </VaTab>
      </template>
    </VaTabs>
    <VaButton
      size="small"
      ref="tab-pag-button-right"
      color="gsgDarkYellow"
      icon="va-arrow-right"
      :hidden="hidePagination"
      :style="`width: ${pagButtonWidth}px; height: ${selectorHeight}px`"
      @click="paginate(1)"
      @keydown-enter="paginate(1)"
    />
  </div>
</template>
<script lang="ts" setup>
import { useI18n } from 'vue-i18n'
import { VaCard, VaCardBlock } from 'vuestic-ui';

const { t } = useI18n()
// import RevenueUpdates from './cards/RevenueReport.vue'
// import ProjectTable from './cards/ProjectTable.vue'
// import RevenueByLocationMap from './cards/RevenueByLocationMap.vue'
// import DataSection from './DataSection.vue'
// import YearlyBreakup from './cards/YearlyBreakup.vue'
// import MonthlyEarnings from './cards/MonthlyEarnings.vue'
// import RegionRevenue from './cards/RegionRevenue.vue'
// import Timeline from './cards/Timeline.vue'
</script>

<template>
  <h1 class="page-title font-bold">{{ t('menu.users') }}</h1>
  <div class="va-table-responsive">
    <table class="va-table va-table--striped">
      <thead>
        <tr>
          <th>Twitch Name</th>
          <th>Country</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="user in users" :key="user.id">
          <td>{{ user.fullName }}</td>
          <td>{{ user.country }}</td>
          <td>
            <VaBadge :text="user.status" :color="user.status" />
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

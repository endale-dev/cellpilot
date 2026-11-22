<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { Name } from '../wailsjs/go/main/App'

const appName = ref('CellPilot')

onMounted(async () => {
  const name = await Name()
  if (name) {
    appName.value = name
  }
})
</script>

<template>
  <main class="window">
    <h1>{{ appName }}</h1>
    <p>Ready</p>
  </main>
</template>

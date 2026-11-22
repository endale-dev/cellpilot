<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { Name } from '../wailsjs/go/main/App'

const appName = ref('CellPilot')
const status = ref('Ready')
const failed = ref(false)

onMounted(async () => {
  try {
    const name = await Name()
    if (!name) {
      throw new Error('application name is empty')
    }
    appName.value = name
  } catch {
    failed.value = true
    status.value = 'The application could not load its name.'
  }
})
</script>

<template>
  <main class="window">
    <h1>{{ appName }}</h1>
    <p role="status" :class="{ error: failed }">{{ status }}</p>
  </main>
</template>

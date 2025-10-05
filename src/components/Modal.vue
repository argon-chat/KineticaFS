<template>
  <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-40">
    <div class="bg-white rounded-lg shadow-lg p-6 min-w-[300px] max-w-[90vw] relative">
      <button class="absolute top-2 right-2 text-purple-600 hover:text-purple-800" @click="close">Close</button>
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue';
const props = defineProps<{ open: boolean }>();
const emit = defineEmits(['close']);
function close() {
  emit('close');
}
</script>

<style scoped>
</style>

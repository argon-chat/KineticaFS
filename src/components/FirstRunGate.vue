<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useFirstRunStore } from '@/stores/firstRun';
import Modal from '@/components/Modal.vue';

const firstRunStore = useFirstRunStore();
const bootstrapLoading = ref(false);
const bootstrapError = ref<string|null>(null);
const showModal = ref(false);
const modalContent = ref<any>(null);

onMounted(() => {
  firstRunStore.checkFirstRun();
});

async function handleBootstrapAdmin() {
  bootstrapLoading.value = true;
  bootstrapError.value = null;
  try {
    const { bootstrapAdminToken } = await import('@/client');
    const response = await bootstrapAdminToken();
    modalContent.value = response;
    showModal.value = true;
    await firstRunStore.checkFirstRun();
  } catch (err: any) {
    bootstrapError.value = err?.message || 'Failed to bootstrap admin token';
  } finally {
    bootstrapLoading.value = false;
  }
}

function closeModal() {
  showModal.value = false;
  modalContent.value = null;
}
</script>

<template>
  <div class="flex flex-col items-center justify-center min-h-screen bg-white text-purple-900">
    <Modal :open="showModal" @close="closeModal">
      <div>
        <h2 class="text-xl font-bold mb-2">Admin Token Created</h2>
        <div class="mb-2 text-sm text-gray-700">This information will <span class="font-bold text-red-600">not be shown anywhere/anytime again</span>. Please copy and store it securely.</div>
        <pre class="bg-gray-100 p-2 rounded text-xs overflow-x-auto">{{ JSON.stringify(modalContent, null, 2) }}</pre>
      </div>
    </Modal>
    <div v-if="firstRunStore.loading" class="text-lg">Checking setup...</div>
    <div v-else-if="firstRunStore.error" class="text-red-600">{{ firstRunStore.error }}</div>
    <div v-else>
      <div v-if="firstRunStore.firstRun === true">
        <h1 class="text-2xl font-bold mb-4">Welcome! First time setup</h1>
        <button class="bg-purple-600 text-white px-4 py-2 rounded" @click="handleBootstrapAdmin" :disabled="bootstrapLoading">
          <span v-if="bootstrapLoading">Bootstrapping...</span>
          <span v-else>Bootstrap Admin Token</span>
        </button>
        <div v-if="bootstrapError" class="text-red-600 mt-2">{{ bootstrapError }}</div>
      </div>
      <div v-else-if="firstRunStore.firstRun === false">
        <h1 class="text-2xl font-bold mb-4">Enter Admin Token</h1>
        <form class="flex flex-col gap-2">
          <input type="password" required placeholder="Admin Token" class="border rounded px-2 py-1" />
          <button class="bg-purple-600 text-white px-4 py-2 rounded" type="submit">Login</button>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
</style>

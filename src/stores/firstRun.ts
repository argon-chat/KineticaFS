import { defineStore } from 'pinia';

export const useFirstRunStore = defineStore('firstRun', {
  state: () => ({
    firstRun: null as null | boolean,
    loading: false,
    error: null as null | string,
  }),
  actions: {
    async checkFirstRun() {
      this.loading = true;
      this.error = null;
      try {
        const { getV1StFirstRun } = await import('@/client');
        const response = await getV1StFirstRun({});
        if (response?.data && typeof response.data.first_run === 'boolean') {
          this.firstRun = response.data.first_run;
        } else {
          this.firstRun = null;
          this.error = 'Unexpected response';
        }
      } catch (err: any) {
        this.error = err?.message || 'Failed to check first run';
        this.firstRun = null;
      } finally {
        this.loading = false;
      }
    },
  },
});

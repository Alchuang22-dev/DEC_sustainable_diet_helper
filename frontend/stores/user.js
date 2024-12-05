// src/stores/user.js
import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useUserStore = defineStore('user', () => {
  const uid = ref(null);
  const isLoggedIn = ref(false);

  function setUid(newUid) {
    uid.value = newUid;
    isLoggedIn.value = !!newUid;
  }

  function clearUid() {
    uid.value = null;
    isLoggedIn.value = false;
  }

  return { uid, isLoggedIn, setUid, clearUid };
});

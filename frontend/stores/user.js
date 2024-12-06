// src/stores/user.js
import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useUserStore = defineStore('user', () => {
  // 用户信息
  const uid = ref(null);
  const isLoggedIn = ref(false);
  const token = ref(''); // 存储 token
  const avatarUrl = ref('');
  const nickName = ref('');

  // 设置用户信息
  function setUid(newUid) {
    uid.value = newUid;
    isLoggedIn.value = !!newUid;
  }

  // 清除用户信息
  function clearUid() {
    uid.value = null;
    isLoggedIn.value = false;
    token.value = '';
    avatarUrl.value = '';
    nickName.value = '';
  }

  // 设置用户的登录状态
  function setIsLoggedIn(status) {
    isLoggedIn.value = status;
  }

  // 设置用户信息，包括 token, avatarUrl, nickName
  function setUserInfo({ token: newToken, user }) {
    token.value = newToken;
    uid.value = user.id;
    nickName.value = user.nickname;
    avatarUrl.value = user.avatar_url || '/static/images/index/background_img.jpg'; // 默认头像
  }

  // 从本地存储加载用户数据
  function loadFromLocalStorage() {
    const storedUserInfo = uni.getStorageSync('userInfo');
    if (storedUserInfo) {
      setUserInfo({
        token: uni.getStorageSync('token'),
        user: storedUserInfo,
      });
    }

    const storedIsLoggedIn = uni.getStorageSync('isLoggedIn');
    if (storedIsLoggedIn !== undefined) {
      setIsLoggedIn(storedIsLoggedIn);
    }
  }

  // 保存用户信息到本地存储
  function saveToLocalStorage() {
    uni.setStorageSync('token', token.value);
    uni.setStorageSync('userInfo', {
      id: uid.value,
      nickname: nickName.value,
      avatar_url: avatarUrl.value,
    });
    uni.setStorageSync('isLoggedIn', isLoggedIn.value);
  }

  // 获取 token
  function getToken() {
    return token.value; // 返回 token
  }

  // 显式保存 token
  function setToken(newToken) {
    token.value = newToken;
    uni.setStorageSync('token', newToken); // 同时保存到本地存储
  }

  return { 
    uid, 
    isLoggedIn, 
    token, 
    avatarUrl, 
    nickName, 
    setUid, 
    clearUid, 
    setIsLoggedIn, 
    setUserInfo, 
    loadFromLocalStorage, 
    saveToLocalStorage, 
    getToken, 
    setToken 
  };
});


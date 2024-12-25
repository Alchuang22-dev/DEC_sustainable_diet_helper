// user.js
import { defineStore } from 'pinia';
import { reactive, watch } from 'vue';

const BASE_URL = 'http://xcxcs.uwdjl.cn:8080';

// 定义用户状态枚举（可根据需要扩展）
export const UserStatus = {
  LOGGED_OUT: 'logged_out', // 未登录
  LOGGED_IN: 'logged_in'    // 已登录
};

const STORAGE_KEY = 'user_store_data';

// 安全存储的 key
const SECURE_STORAGE_KEY = 'refresh_token_secure';

// 获取本地存储的 token 和 refresh_token
const storedData = uni.getStorageSync(STORAGE_KEY);
const initialToken = storedData ? JSON.parse(storedData).token : '';
console.log('token:', initialToken);

// 定义 access_token 有效期（单位：秒），这里修改为 30 分钟
const ACCESS_TOKEN_EXPIRES_IN = 30 * 60; // 1800 秒

// 定义一个全局定时器 ID
let tokenRefreshTimer = null;

// 定义一个标志位，防止多次刷新
let isRefreshing = false;

const request = (config) => {
  return new Promise((resolve, reject) => {
    uni.request({
      ...config,
      success: (res) => {
        if (res.statusCode === 401) {
          // 处理未授权，比如刷新token
          if (!isRefreshing) {
            isRefreshing = true;
            // 刷新token的逻辑
            refreshToken()
              .then(() => {
                isRefreshing = false;
                // 重新发送原始请求
                resolve(request(config));
              })
              .catch((err) => {
                isRefreshing = false;
                reject(err);
              });
          }
        } else {
          resolve(res);
        }
      },
      fail: (err) => reject(err)
    });
  });
};

// 安全存储的辅助函数
const secureStorage = {
  setRefreshToken: (token) => {
    if (uni.canIUse('setStorage')) {
      // 使用加密或其他安全方法存储
      // 这里假设使用简单的存储，实际项目中应使用加密
      uni.setStorageSync(SECURE_STORAGE_KEY, token);
      // 存储 refresh_token 的获取时间，以便后续校验有效期
      uni.setStorageSync(`${SECURE_STORAGE_KEY}_timestamp`, Date.now());
    } else {
      console.error('当前平台不支持安全存储');
    }
  },
  getRefreshToken: () => {
    if (uni.canIUse('getStorage')) {
      const token = uni.getStorageSync(SECURE_STORAGE_KEY);
      const timestamp = uni.getStorageSync(`${SECURE_STORAGE_KEY}_timestamp`);
      if (token && timestamp) {
        const currentTime = Date.now();
        const sevenDaysInMs = 7 * 24 * 60 * 60 * 1000;
        if (currentTime - timestamp < sevenDaysInMs) {
          return token;
        } else {
          // refresh_token 已过期
          secureStorage.removeRefreshToken();
          return null;
        }
      }
      return token;
    } else {
      console.error('当前平台不支持安全存储');
      return null;
    }
  },
  removeRefreshToken: () => {
    if (uni.canIUse('removeStorage')) {
      uni.removeStorageSync(SECURE_STORAGE_KEY);
      uni.removeStorageSync(`${SECURE_STORAGE_KEY}_timestamp`);
    } else {
      console.error('当前平台不支持安全存储');
    }
  }
};

export const useUserStore = defineStore('user', () => {
  const getInitialState = () => {
    try {
      const storedData = uni.getStorageSync(STORAGE_KEY);
      return storedData ? JSON.parse(storedData) : {
        uid: null,
        isLoggedIn: false,
        token: '',
        avatarUrl: '',
        nickName: '',
        status: UserStatus.LOGGED_OUT,
        tokenExpiry: null, // 新增字段，记录 access_token 过期时间
        registered_days: 0 // 新增字段，用于存储注册天数
      };
    } catch (error) {
      console.error('获取存储的用户数据失败:', error);
      return {
        uid: null,
        isLoggedIn: false,
        token: '',
        avatarUrl: '',
        nickName: '',
        status: UserStatus.LOGGED_OUT,
        tokenExpiry: null, // 新增字段
        registered_days: 0 // 新增字段
      };
    }
  };

  const user = reactive(getInitialState());

  const saveToStorage = () => {
    try {
      uni.setStorageSync(STORAGE_KEY, JSON.stringify(user));
      console.log('已保存用户数据:', user);
    } catch (error) {
      console.error('保存用户数据失败:', error);
    }
  };

  const watchUser = () => {
    const watchKeys = ['uid', 'isLoggedIn', 'token', 'avatarUrl', 'nickName', 'status', 'tokenExpiry', 'registered_days']; // 添加 'registered_days'
    watchKeys.forEach(key => {
      watch(() => user[key], () => {
        saveToStorage();
      }, { deep: true });
    });
  };

  const createRequestConfig = (config) => {
    return {
      ...config,
      header: {
        'Authorization': `Bearer ${user.token || initialToken}`,
        ...(config.header || {})
      }
    };
  };

  /**
   * 用户登录/注册（微信登录）
   */
  const login = async () => {
    try {
      // 调用 uni.login 获取微信登录凭证
      const loginRes = await new Promise((resolve, reject) => {
        uni.login({
          "provider": "weixin",
          "onlyAuthorize": true,
          success: resolve,
          fail: reject,
        });
      });

      if (!loginRes.code) {
        throw new Error('微信登录失败，请重试');
      }

      const authRes = await new Promise((resolve, reject) => {
        console.log('url:', `${BASE_URL}/users/auth`);
        uni.request({
          url: `${BASE_URL}/users/auth`, // 确保端点正确
          method: 'POST',
          data: {
            code: loginRes.code,
          },
          // filePath: avatarUrl,
          name: 'avatar', // 对应后端表单文件字段名
          success: resolve,
          fail: reject,
        });
      });

      // 打印返回的数据，调试用
      console.log("返回的完整响应：", authRes);
      console.log("后端返回的数据：", authRes.data);

      const returnData = authRes.data;
      console.log("returnData", returnData);

      // 检查 code 是否是 200
      if (authRes.statusCode !== 200) {
        throw new Error(returnData.message || '登录失败');
      }

      // 更新用户状态
      user.token = returnData.access_token;
      // 使用安全存储存储 refresh_token
      secureStorage.setRefreshToken(returnData.refresh_token);
      user.uid = returnData.user.id;
      user.nickName = returnData.user.nickname;
      user.avatarUrl = returnData.user.avatar_url || '/static/images/default_avatar.jpg';
      user.isLoggedIn = true;
      user.status = UserStatus.LOGGED_IN;
      user.tokenExpiry = Date.now() + ACCESS_TOKEN_EXPIRES_IN * 1000; // 设置过期时间
      user.registered_days = returnData.registered_days || 0; // 设置注册天数

      saveToStorage();
      scheduleTokenRefresh();

      return returnData;
    } catch (error) {
      console.error('登录失败:', error);
      throw error;
    }
  };

  // 用户登出
  const logout = async () => {
    console.log('url logout:', `${BASE_URL}/users/logout`);
    try {
      const response = await request(createRequestConfig({
        url: `${BASE_URL}/users/logout`,
        method: 'POST',
        data: {
          // 从安全存储中获取 refresh_token
          refresh_token: secureStorage.getRefreshToken()
        }
      }));
      console.log('登出成功:', response.data);
      reset();
      return response.data;
    } catch (error) {
      console.error('登出失败:', error);
      throw error;
    }
  };

  // 设置用户名
  const setNickname = async (newNickname) => {
    try {
      const response = await request(createRequestConfig({
        url: `${BASE_URL}/users/set_nickname`,
        method: 'PUT',
        data: {
          nickname: newNickname
        }
      }));
      console.log('设置昵称成功:', response.data);
      user.nickName = response.data.nickname;
      return response.data;
    } catch (error) {
      console.error('设置昵称失败:', error);
      throw error;
    }
  };

  // 设置头像
  const setAvatar = async (avatarFile) => {
    try {
      // 使用 uni.uploadFile 上传头像文件
      const uploadRes = await new Promise((resolve, reject) => {
        uni.uploadFile({
          url: `${BASE_URL}/users/set_avatar`, // 确保接口地址正确
          filePath: avatarFile, // 本地文件路径
          name: 'avatar', // 后端字段名称
          header: {
            'Authorization': `Bearer ${user.token || initialToken}` // 添加授权头部
          },
          success: resolve,
          fail: reject
        });
      });

      console.log('上传响应：', uploadRes);

      // 解析后端返回数据
      if (uploadRes.statusCode !== 200) {
        throw new Error(`头像上传失败：${uploadRes.data}`);
      }

      const responseData = JSON.parse(uploadRes.data);

      // 更新用户状态
      user.avatarUrl = responseData.avatar_url; // 后端返回的头像地址
      saveToStorage();

      console.log('设置头像成功:', responseData);

      return responseData;
    } catch (error) {
      console.error('设置头像失败:', error);
      throw error;
    }
  };

  // 刷新令牌
  const refreshToken = async () => {
    try {
      const currentRefreshToken = secureStorage.getRefreshToken();
      if (!currentRefreshToken) {
        throw new Error('refresh_token 已过期或不存在');
      }

      const response = await request(createRequestConfig({
        url: `${BASE_URL}/users/refresh`,
        method: 'POST',
        data: {
          // 从安全存储中获取 refresh_token
          refresh_token: currentRefreshToken
        },
        header: {
          'Content-Type': 'application/json'
        }
      }));

      if (response.statusCode === 200) {
        const { access_token, refresh_token: newRefreshToken, registered_days } = response.data;
        user.token = access_token;
        // 更新 refresh_token 至安全存储
        secureStorage.setRefreshToken(newRefreshToken);
        user.tokenExpiry = Date.now() + ACCESS_TOKEN_EXPIRES_IN * 1000; // 更新过期时间
        user.isLoggedIn = true;
        user.status = UserStatus.LOGGED_IN;
        user.registered_days = registered_days || user.registered_days; // 更新注册天数

        saveToStorage();
        scheduleTokenRefresh();
        console.log('刷新令牌成功:', response.data);
        return response.data;
      } else {
        throw new Error(response.data.message || '刷新令牌失败');
      }
    } catch (error) {
      console.error('刷新令牌失败:', error);
      reset(); // 刷新失败时重置用户状态
      // 跳转到登录界面
      uni.navigateTo({
        url: '/pagesMy/login/login',
      });
      throw error;
    }
  };

  // 返回用户uid
  const getUserID = () => {
    return user.uid;
  };

  // 清除本地存储数据
  const clearStorage = () => {
    try {
      uni.removeStorageSync(STORAGE_KEY);
      secureStorage.removeRefreshToken();
    } catch (error) {
      console.error('清除用户存储数据失败:', error);
    }
  };

  // 重置状态
  const reset = () => {
    user.uid = null;
    user.isLoggedIn = false;
    user.token = '';
    user.avatarUrl = '';
    user.nickName = '';
    user.status = UserStatus.LOGGED_OUT;
    user.tokenExpiry = null;
    user.registered_days = 0; // 重置注册天数
    clearStorage();
    clearTokenRefreshTimer();
  };

  // 加载本地存储的数据
  const loadFromLocalStorage = () => {
    try {
      const storedData = uni.getStorageSync(STORAGE_KEY);
      if (storedData) {
        const parsedData = JSON.parse(storedData);
        user.uid = parsedData.uid;
        user.isLoggedIn = parsedData.isLoggedIn;
        user.token = parsedData.token;
        user.avatarUrl = parsedData.avatarUrl;
        user.nickName = parsedData.nickName;
        user.status = parsedData.status;
        user.tokenExpiry = parsedData.tokenExpiry;
        user.registered_days = parsedData.registered_days || 0; // 加载注册天数
        if (user.isLoggedIn && user.token && user.tokenExpiry) {
          console.log('已加载本地用户数据:', user);
        }
      }
    } catch (error) {
      console.error('加载本地用户数据失败:', error);
    }
  };

  // 获取当前状态的可读文本
  const getStatusText = () => {
    const statusTexts = {
      [UserStatus.LOGGED_OUT]: '未登录',
      [UserStatus.LOGGED_IN]: '已登录'
    };
    return statusTexts[user.status] || '未知状态';
  };

  /**
   * 获取用户基本信息
   */
  const fetchBasicDetails = async () => {
    try {
      const response = await request(createRequestConfig({
        url: `${BASE_URL}/users/basic_details`,
        method: 'GET',
        header: {
          'Content-Type': 'application/json'
        }
      }));

      if (response.statusCode === 200) {
        const data = response.data;
        // 假设返回的数据结构为 { id, nickname, avatar_url, registered_days }
        user.uid = data.id;
        user.nickName = data.nickname;
        user.avatarUrl = data.avatar_url || '/static/images/default_avatar.jpg';
        user.registered_days = data.registered_days || 0; // 设置注册天数
        user.isLoggedIn = true;
        user.status = UserStatus.LOGGED_IN;
        user.tokenExpiry = Date.now() + ACCESS_TOKEN_EXPIRES_IN * 1000; // 更新过期时间

        saveToStorage();
        console.log('用户基本信息已更新:', data);
      } else {
        throw new Error(response.data.message || '获取用户基本信息失败');
      }
    } catch (error) {
      console.error('获取用户基本信息失败:', error);
      // 可选择在此处处理错误，例如重置用户状态
      // reset();
      throw error;
    }
  };

  // 定时器相关函数
  const scheduleTokenRefresh = () => {
    clearTokenRefreshTimer();
    if (user.tokenExpiry) {
      const currentTime = Date.now();
      const timeToExpiry = user.tokenExpiry - currentTime;
      const timeBeforeRefresh = timeToExpiry - (20 * 1000); // 20 秒前刷新

      if (timeBeforeRefresh > 0) {
        tokenRefreshTimer = setTimeout(() => {
          refreshToken().catch(err => {
            console.error('自动刷新token失败:', err);
          });
        }, timeBeforeRefresh);
        console.log(`将在 ${Math.floor(timeBeforeRefresh / 1000)} 秒后刷新token`);
      } else {
        // 如果时间已经不足，立即刷新
        refreshToken().catch(err => {
          console.error('自动刷新token失败:', err);
        });
      }
    }
  };

  const clearTokenRefreshTimer = () => {
    if (tokenRefreshTimer) {
      clearTimeout(tokenRefreshTimer);
      tokenRefreshTimer = null;
    }
  };

  watchUser();

  // 初始化时加载本地存储的数据
  loadFromLocalStorage();

  // 在 store 初始化时调用 refreshToken，确保 token 的有效性
  const initialize = async () => {
    if (user.isLoggedIn && secureStorage.getRefreshToken()) {
      try {
        await refreshToken();
        console.log('应用启动时刷新token成功');
      } catch (error) {
        console.error('应用启动时刷新token失败:', error);
        // 已在 refreshToken 中处理跳转到登录
      }
    }
  };

  fetchBasicDetails(); // 获取用户基本信息
  initialize(); // 调用初始化方法

  return {
    user,
    UserStatus,
    getStatusText,
    login,
    setNickname,
    setAvatar,
    refreshToken,
    getUserID,
    logout,
    reset,
    clearStorage,
    loadFromLocalStorage,
    fetchBasicDetails,
  };
});
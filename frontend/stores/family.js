import { defineStore } from 'pinia';
import { reactive, watch } from 'vue';

const BASE_URL = 'http://122.51.231.155:8080';

// 定义家庭状态枚举
export const FamilyStatus = {
    NOT_JOINED: 'empty',           // 未加入
    PENDING_APPROVAL: 'waiting',   // 申请加入待审核
    JOINED: 'family'               // 已加入
};

const STORAGE_KEY = 'family_store_data';
const token = uni.getStorageSync('token');
console.log('token:', token);

// 封装request为Promise
const request = (config) => {
    return new Promise((resolve, reject) => {
        uni.request({
            ...config,
            success: (res) => resolve(res),
            fail: (err) => reject(err)
        });
    });
};

export const useFamilyStore = defineStore('family', () => {
    const getInitialState = () => {
        try {
            const storedData = uni.getStorageSync(STORAGE_KEY);
            return storedData ? JSON.parse(storedData) : {
                id: '',
                name: '',
                familyId: '',
                memberCount: 0,
                admins: [],
                members: [],
                status: FamilyStatus.NOT_JOINED,
            };
        } catch (error) {
            console.error('Failed to get stored family data:', error);
            return {
                id: '',
                name: '',
                familyId: '',
                memberCount: 0,
                admins: [],
                members: [],
                status: FamilyStatus.NOT_JOINED,
            };
        }
    };

    const family = reactive(getInitialState());

    const saveToStorage = () => {
        try {
            uni.setStorageSync(STORAGE_KEY, JSON.stringify(family));
            console.log('Saved family data:', family);
        } catch (error) {
            console.error('Failed to save family data:', error);
        }
    };

    const watchFamily = () => {
        const watchKeys = ['id', 'name', 'familyId', 'memberCount', 'admins', 'members', 'status'];
        watchKeys.forEach(key => {
            watch(() => family[key], () => {
                saveToStorage();
            }, { deep: true });
        });
    };

    const createRequestConfig = (config) => {
        return {
            ...config,
            header: {
                'Authorization': `Bearer ${token}`
            }
        };
    };

    // 创建家庭
    const createFamily = async (familyName) => {
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/families/create`,
                method: 'POST',
                data: {
                    name: familyName
                }
            }));
            console.log('createFamily:', response.data);

            // 在每个赋值后都加上日志
            console.log('Setting id:', response.data.family.id);
            family.id = response.data.family.id;

            console.log('Setting name:', response.data.family.name);
            family.name = response.data.family.name;

            console.log('Setting familyId:', response.data.family.family_id);
            family.familyId = response.data.family.family_id;

            console.log('Setting status');
            family.status = FamilyStatus.JOINED;

            console.log('Before getFamilyDetails'); // 添加这行
            const details = await getFamilyDetails();
            console.log('After getFamilyDetails', details); // 添加这行

            saveToStorage();
            return response.data;
        } catch (error) {
            throw error;
        }
    };

    // 获取家庭详情
    const getFamilyDetails = async () => {
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/families/details`,
                method: 'GET'
            }));

            console.log('getFamilyDetails:', response);
            // 假设后端返回数据格式为 response.data 内包含家庭信息和status
            const data = response.data;

            family.status = data.status;

            if (data.status === FamilyStatus.JOINED) {
                // 已加入家庭
                family.id = data.id;
                family.name = data.name;
                family.familyId = data.family_id;
                family.memberCount = data.member_count;
                family.admins = data.admins.map(admin => ({
                    id: admin.id,
                    nickname: admin.nickname,
                    avatarUrl: admin.avatar_url
                }));
                family.members = data.members.map(member => ({
                    id: member.id,
                    nickname: member.nickname,
                    avatarUrl: member.avatar_url
                }));
                console.log('family:', family);
            } else if (data.status === FamilyStatus.PENDING_APPROVAL) {
                // 待审核状态
                family.id = data.id;
                family.name = data.name;
                family.familyId = data.family_id;
                family.memberCount = 0;
                family.admins = [];
                family.members = [];
            } else {
                // 未加入状态，清空数据
                family.id = '';
                family.name = '';
                family.familyId = '';
                family.memberCount = 0;
                family.admins = [];
                family.members = [];
            }

            saveToStorage();
            return response;
        } catch (error) {
            console.log('getFamilyDetails error:', error);
            throw error;
        }
    };

    // 搜索家庭
    const searchFamily = async (familyId) => {
        try {
            console.log('searchFamily:', familyId);
            // 将familyId作为query参数拼接到URL
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/families/search?family_id=${familyId}`,
                method: 'GET'
            }));
            return response;
        } catch (error) {
            throw error;
        }
    };

    // 申请加入家庭
    const joinFamily = async (familyId) => {
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/families/${familyId}/join`,
                method: 'POST'
            }));
            await getFamilyDetails(); // 重新获取状态
            return response;
        } catch (error) {
            throw error;
        }
    };

    // 取消加入申请
    const cancelJoinRequest = async () => {
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/families/cancel_join`,
                method: 'POST'
            }));
            await getFamilyDetails(); // 重新获取状态
            return response;
        } catch (error) {
            throw error;
        }
    };

    // 管理员：同意加入申请
    const admitJoinRequest = async (userId) => {
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/families/admit`,
                method: 'POST',
                data: { user_id: userId }
            }));
            await getFamilyDetails();
            return response;
        } catch (error) {
            throw error;
        }
    };

    // 管理员：拒绝加入申请
    const rejectJoinRequest = async (userId) => {
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/families/reject`,
                method: 'POST',
                data: { user_id: userId }
            }));
            await getFamilyDetails();
            return response;
        } catch (error) {
            throw error;
        }
    };

    // 判断用户是否是管理员
    const isAdmin = (userId) => {
        return family.admins.some(admin => admin.id === userId);
    };

    // 清除本地存储数据
    const clearStorage = () => {
        try {
            uni.removeStorageSync(STORAGE_KEY);
        } catch (error) {
            console.error('Failed to clear family storage:', error);
        }
    };

    // 重置状态
    const reset = () => {
        family.id = '';
        family.name = '';
        family.familyId = '';
        family.memberCount = 0;
        family.admins = [];
        family.members = [];
        family.status = FamilyStatus.NOT_JOINED;
        clearStorage();
    };

    // 获取当前状态的可读文本
    const getStatusText = () => {
        const statusTexts = {
            [FamilyStatus.NOT_JOINED]: '未加入',
            [FamilyStatus.PENDING_APPROVAL]: '待审核',
            [FamilyStatus.JOINED]: '已加入'
        };
        return statusTexts[family.status] || '未知状态';
    };

    watchFamily();

    return {
        family,
        FamilyStatus,
        useFamilyStore,
        getStatusText,
        createFamily,
        getFamilyDetails,
        searchFamily,
        joinFamily,
        cancelJoinRequest,
        admitJoinRequest,
        rejectJoinRequest,
        reset,
        clearStorage,
        isAdmin
    };
});

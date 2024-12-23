// family.js
import { defineStore } from 'pinia';
import { reactive, watch, computed } from 'vue';
import { useUserStore } from "./user.js";

const BASE_URL = 'http://122.51.231.155:8080';

// 定义家庭状态枚举
export const FamilyStatus = {
    NOT_JOINED: 'empty',           // 未加入
    PENDING_APPROVAL: 'waiting',   // 申请加入待审核
    JOINED: 'family'               // 已加入
};

const STORAGE_KEY = 'family_store_data';

// 获取当前时区
function getTimeZone() {
    return Intl.DateTimeFormat().resolvedOptions().timeZone;
}

// 封装request为Promise，并处理401状态码
const request = (config) => {
    // 获取 userStore 实例
    const userStore = useUserStore();
    return new Promise((resolve, reject) => {
        uni.request({
            ...config,
            success: async (res) => {
                if (res.statusCode === 401) {
                    // 遇到未授权，尝试刷新token
                    console.log('401 Unauthorized, trying to refresh token...');
                    try {
                        await userStore.refreshToken();
                        // 使用新的 token 重试请求
                        uni.request({
                            ...config,
                            header: {
                                ...config.header,
                                'Authorization': `Bearer ${userStore.user.token}`
                            },
                            success: (res2) => {
                                if (res2.statusCode === 401) {
                                    reject(new Error('Unauthorized'));
                                } else {
                                    resolve(res2);
                                }
                            },
                            fail: (err2) => reject(err2)
                        });
                    } catch (error) {
                        console.error('Token 刷新失败:', error);
                        // 刷新失败，跳转登录或采取其他措施
                        uni.navigateTo({
                            url: '/pagesMy/login/login',
                        });
                        reject(new Error('Unauthorized'));
                    }
                } else {
                    resolve(res);
                }
            },
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
                allMembers: [],
                waiting_members: [],
                dishProposals: [],
                status: FamilyStatus.NOT_JOINED,
                // 新增存储字段（如果本地无则默认空）
                memberDailyData: [],
                familySums: {}
            };
        } catch (error) {
            console.error('Failed to get stored family data:', error);
            return {
                id: '',
                name: '',
                familyId: '',
                memberCount: 0,
                allMembers: [],
                waiting_members: [],
                dishProposals: [],
                status: FamilyStatus.NOT_JOINED,
                memberDailyData: [],
                familySums: {}
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
        const watchKeys = [
            'id', 'name', 'familyId', 'memberCount',
            'allMembers', 'waiting_members', 'status',
            'dishProposals', 'memberDailyData', 'familySums'
        ];
        watchKeys.forEach(key => {
            watch(() => family[key], () => {
                saveToStorage();
            }, { deep: true });
        });
    };

    const createRequestConfig = (config) => {
        // 获取 userStore 实例
        const userStore = useUserStore();
        return {
            ...config,
            header: {
                'Authorization': `Bearer ${userStore.user.token || ''}`,
                ...(config.header || {})
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

            family.id = response.data.family.id;
            family.name = response.data.family.name;
            family.familyId = response.data.family.family_id;
            family.status = FamilyStatus.JOINED;

            await getFamilyDetails();

            return response.data;
        } catch (error) {
            console.error('创建家庭失败:', error);
            throw error;
        }
    };

    // 获取家庭详情
    const getFamilyDetails = async () => {
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/families/details?timezone=${encodeURIComponent(getTimeZone())}`,
                method: 'GET'
            }));

            console.log('getFamilyDetails:', response);
            const data = response.data;

            family.status = data.status;

            if (data.status === FamilyStatus.JOINED) {
                family.id = data.id;
                family.name = data.name;
                family.familyId = data.family_id;
                family.memberCount = data.member_count;

                // 合并管理员和普通成员
                const adminsWithRole = data.admins.map(admin => ({
                    id: admin.id,
                    nickname: admin.nickname,
                    avatarUrl: admin.avatar_url,
                    role: 'admin'
                }));
                const membersWithRole = data.members.map(member => ({
                    id: member.id,
                    nickname: member.nickname,
                    avatarUrl: member.avatar_url,
                    role: 'member'
                }));
                family.allMembers = [...adminsWithRole, ...membersWithRole];

                // 处理等待加入的成员
                if (data.waiting_members && Array.isArray(data.waiting_members)) {
                    family.waiting_members = data.waiting_members.map(member => ({
                        id: member.id,
                        nickname: member.nickname,
                        avatarUrl: member.avatar_url
                    }));
                } else {
                    family.waiting_members = [];
                }

                // ============= 新增：存储每日数据与总汇数据 =============
                family.memberDailyData = data.member_daily_data || [];
                family.familySums = data.family_sums || {};

            } else if (data.status === FamilyStatus.PENDING_APPROVAL) {
                family.id = data.id;
                family.name = data.name;
                family.familyId = data.family_id;
                family.memberCount = 0;
                family.allMembers = [];
                family.waiting_members = [];
            } else {
                family.id = '';
                family.name = '';
                family.familyId = '';
                family.memberCount = 0;
                family.allMembers = [];
                family.waiting_members = [];
                family.dishProposals = [];
                // 重置每日数据和汇总数据
                family.memberDailyData = [];
                family.familySums = {};
            }

            return response;
        } catch (error) {
            console.log('getFamilyDetails error:', error);
            throw error;
        }
    };

    // 获取想吃的菜品列表
    const getDesiredDishes = async () => {
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/families/desired_dishes`,
                method: 'GET'
            }));
            family.dishProposals = response.data; // 直接存储后端返回的数组
            console.log('getDesiredDishes:', family.dishProposals);
            return response.data;
        } catch (error) {
            console.error('getDesiredDishes error:', error);
            throw error;
        }
    };

    // 添加想吃的菜品
    const addDishProposal = async ({ dishId, preference }) => {
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/families/add_desired_dish`,
                method: 'POST',
                data: {
                    dish_id: dishId,
                    level_of_desire: preference
                }
            }));

            if (response.statusCode === 200) {
                // 添加成功后更新列表
                await getDesiredDishes();
                return response.data;
            } else if (response.statusCode === 400 && response.data.error === "You have already desired this dish") {
                throw new Error('DISH_ALREADY_EXISTS');
            } else {
                throw new Error(response.data.error || 'Unknown error');
            }
        } catch (error) {
            console.error('addDishProposal error:', error);
            throw error;
        }
    };

    // 删除想吃的菜品
    const deleteDesiredDish = async (dishId) => {
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/families/desired_dishes`,
                method: 'DELETE',
                data: {
                    dish_id: dishId
                }
            }));
            await getDesiredDishes();
            return response.data;
        } catch (error) {
            console.error('deleteDesiredDish error:', error);
            throw error;
        }
    };

    // 搜索家庭
    const searchFamily = async (familyId) => {
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/families/search?family_id=${familyId}`,
                method: 'GET'
            }));
            return response.data;
        } catch (error) {
            console.error('searchFamily error:', error);
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
            await getFamilyDetails();
            return response.data;
        } catch (error) {
            console.error('joinFamily error:', error);
            throw error;
        }
    };

    // 取消加入申请
    const cancelJoinRequest = async () => {
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/families/cancel_join`,
                method: 'DELETE'
            }));
            await getFamilyDetails();
            return response.data;
        } catch (error) {
            console.error('cancelJoinRequest error:', error);
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
            return response.data;
        } catch (error) {
            console.error('admitJoinRequest error:', error);
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
            return response.data;
        } catch (error) {
            console.error('rejectJoinRequest error:', error);
            throw error;
        }
    };

    // 离开家庭
    const leaveFamily = async () => {
        try {
            // 获取 userStore 实例
            const userStore = useUserStore();
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/families/leave_family`,
                method: 'DELETE',
                data: {
                    user_id: userStore.user.uid
                }
            }));
            reset();
            return response.data;
        } catch (error) {
            console.error('leaveFamily error:', error);
            throw error;
        }
    };

    // 解散家庭
    const breakFamily = async () => {
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/families/break`,
                method: 'DELETE'
            }));
            reset();
            return response.data;
        } catch (error) {
            console.error('breakFamily error:', error);
            throw error;
        }
    };

    // 设为管理员
    const setAdmin = async (userId) => {
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/families/set_admin`,
                method: 'PUT',
                data: { user_id: userId }
            }));
            await getFamilyDetails();
            return response.data;
        } catch (error) {
            console.error('setAdmin error:', error);
            throw error;
        }
    };

    // 删除家庭成员
    const removeFamilyMember = async (userId) => {
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/families/delete_family_member`,
                method: 'DELETE',
                data: { user_id: userId }
            }));
            await getFamilyDetails();
            return response.data;
        } catch (error) {
            console.error('removeFamilyMember error:', error);
            throw error;
        }
    };

    // 判断用户是否是管理员
    const isAdmin = (userId) => {
        return family.allMembers.some(member => member.id === userId && member.role === 'admin');
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
        family.allMembers = [];
        family.waiting_members = [];
        family.dishProposals = [];
        family.status = FamilyStatus.NOT_JOINED;
        family.memberDailyData = [];
        family.familySums = {};
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
        getStatusText,
        createFamily,
        getFamilyDetails,
        searchFamily,
        joinFamily,
        cancelJoinRequest,
        admitJoinRequest,
        rejectJoinRequest,
        setAdmin,
        removeFamilyMember,
        leaveFamily,
        breakFamily,
        reset,
        clearStorage,
        isAdmin,
        addDishProposal,
        getDesiredDishes,
        deleteDesiredDish
    };
});

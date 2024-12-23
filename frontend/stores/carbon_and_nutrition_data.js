// carbon_and_nutrition_data.js
import { defineStore } from 'pinia';
import { reactive, watch } from 'vue';
import { useUserStore } from "./user.js";

const BASE_URL = 'http://122.51.231.155:8095';
const STORAGE_KEY = 'carbon_and_nutrition_store_data';

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

export const useCarbonAndNutritionStore = defineStore('carbon_and_nutrition', () => {
    const getInitialState = () => {
        try {
            const storedData = uni.getStorageSync(STORAGE_KEY);
            return storedData ? JSON.parse(storedData) : {
                nutritionGoals: [],
                carbonGoals: [],
                nutritionIntakes: [],
                carbonIntakes: [],
                sharedNutritionCarbonIntakes: [] // 新增共享记录
            };
        } catch (error) {
            console.error('Failed to get stored carbon_and_nutrition data:', error);
            return {
                nutritionGoals: [],
                carbonGoals: [],
                nutritionIntakes: [],
                carbonIntakes: [],
                sharedNutritionCarbonIntakes: [] // 新增共享记录
            };
        }
    };

    const state = reactive(getInitialState());

    const saveToStorage = () => {
        try {
            uni.setStorageSync(STORAGE_KEY, JSON.stringify(state));
            console.log('Saved carbon_and_nutrition data:', state);
        } catch (error) {
            console.error('Failed to save carbon_and_nutrition data:', error);
        }
    };

    const watchState = () => {
        const watchKeys = ['nutritionGoals', 'carbonGoals', 'nutritionIntakes', 'carbonIntakes', 'sharedNutritionCarbonIntakes']; // 包含共享记录
        watchKeys.forEach(key => {
            watch(() => state[key], () => {
                saveToStorage();
            }, { deep: true });
        });
    };

    const createRequestConfig = (config) => {
        const userStore = useUserStore();
        return {
            ...config,
            header: {
                'Authorization': `Bearer ${userStore.user.token || ''}`,
                'Content-Type': 'application/json', // 确保内容类型为JSON
                ...(config.header || {})
            }
        };
    };

    // 设置营养目标 (POST /nutrition-carbon/nutrition/goals)
    const setNutritionGoals = async (goals) => {
        // goals为数组，格式参考API文档
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/nutrition-carbon/nutrition/goals`,
                method: 'POST',
                data: goals
            }));
            console.log('setNutritionGoals:', response.data);
            // 设置成功后，可选择立即获取最新的营养目标列表，以更新 state
            await getNutritionGoals();
            return response.data;
        } catch (error) {
            console.error('setNutritionGoals error:', error);
            throw error;
        }
    };

    // 获取营养目标 (GET /nutrition-carbon/nutrition/goals)
    const getNutritionGoals = async () => {
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/nutrition-carbon/nutrition/goals`,
                method: 'GET'
            }));
            console.log('getNutritionGoals:', response.data);
            if (response.statusCode === 200 && response.data && response.data.data) {
                state.nutritionGoals = response.data.data;
            }
            return response.data;
        } catch (error) {
            console.error('getNutritionGoals error:', error);
            throw error;
        }
    };

    // 设置碳排放目标 (POST /nutrition-carbon/carbon/goals)
    const setCarbonGoals = async (goals) => {
        // goals为数组，格式参考API文档
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/nutrition-carbon/carbon/goals`,
                method: 'POST',
                data: goals
            }));
            console.log('setCarbonGoals:', response.data);
            // 设置成功后更新本地数据
            await getCarbonGoals();
            return response.data;
        } catch (error) {
            console.error('setCarbonGoals error:', error);
            throw error;
        }
    };

    // 获取碳排放目标 (GET /nutrition-carbon/carbon/goals)
    const getCarbonGoals = async () => {
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/nutrition-carbon/carbon/goals`,
                method: 'GET'
            }));
            console.log('getCarbonGoals:', response.data);
            if (response.statusCode === 200 && response.data && response.data.data) {
                state.carbonGoals = response.data.data;
            }
            return response.data;
        } catch (error) {
            console.error('getCarbonGoals error:', error);
            throw error;
        }
    };

    // 获取实际营养摄入 (GET /nutrition-carbon/nutrition/intakes)
    const getNutritionIntakes = async () => {
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/nutrition-carbon/nutrition/intakes`,
                method: 'GET'
            }));
            console.log('getNutritionIntakes:', response.data);
            if (response.statusCode === 200 && Array.isArray(response.data)) {
                state.nutritionIntakes = response.data;
            }
            return response.data;
        } catch (error) {
            console.error('getNutritionIntakes error:', error);
            throw error;
        }
    };

    // 获取实际碳排放摄入 (GET /nutrition-carbon/carbon/intakes)
    const getCarbonIntakes = async () => {
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/nutrition-carbon/carbon/intakes`,
                method: 'GET'
            }));
            console.log('getCarbonIntakes:', response.data);
            if (response.statusCode === 200 && Array.isArray(response.data)) {
                state.carbonIntakes = response.data;
            }
            return response.data;
        } catch (error) {
            console.error('getCarbonIntakes error:', error);
            throw error;
        }
    };

    // 共享营养碳排放记录 (POST /nutrition-carbon/shared/nutrition-carbon)
    const setSharedNutritionCarbonIntake = async (sharedData) => {
        /**
         * sharedData 格式:
         * {
         *   date: "2024-03-21T00:00:00Z",
         *   meal_type: "breakfast",
         *   calories: 1000,
         *   protein: 30,
         *   fat: 40,
         *   carbohydrates: 120,
         *   sodium: 1000,
         *   emission: 5.0,
         *   user_shares: [
         *     { user_id: 1, ratio: 0.6 },
         *     { user_id: 2, ratio: 0.4 }
         *   ]
         * }
         */
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/nutrition-carbon/shared/nutrition-carbon`,
                method: 'POST',
                data: sharedData
            }));
            console.log('setSharedNutritionCarbonIntake:', response.data);
            // 设置成功后，可选择立即获取最新的共享记录，以更新 state
            await getSharedNutritionCarbonIntakes();
            return response.data;
        } catch (error) {
            console.error('setSharedNutritionCarbonIntake error:', error);
            // 根据错误信息进行处理
            if (error.response && error.response.data && error.response.data.error) {
                throw new Error(error.response.data.error);
            } else {
                throw error;
            }
        }
    };

    // 获取共享营养碳排放记录 (假设有对应的 GET 接口，如果没有，需要根据实际情况调整)
    const getSharedNutritionCarbonIntakes = async () => {
        try {
            const response = await request(createRequestConfig({
                url: `${BASE_URL}/nutrition-carbon/shared/nutrition-carbon`, // 假设的GET端点
                method: 'GET'
            }));
            console.log('getSharedNutritionCarbonIntakes:', response.data);
            if (response.statusCode === 200 && Array.isArray(response.data)) {
                state.sharedNutritionCarbonIntakes = response.data;
            }
            return response.data;
        } catch (error) {
            console.error('getSharedNutritionCarbonIntakes error:', error);
            throw error;
        }
    };

    // 重置状态
    const reset = () => {
        state.nutritionGoals = [];
        state.carbonGoals = [];
        state.nutritionIntakes = [];
        state.carbonIntakes = [];
        state.sharedNutritionCarbonIntakes = []; // 重置共享记录
        clearStorage();
    };

    // 清除本地存储数据
    const clearStorage = () => {
        try {
            uni.removeStorageSync(STORAGE_KEY);
        } catch (error) {
            console.error('Failed to clear carbon_and_nutrition storage:', error);
        }
    };

    watchState();


    // 添加辅助方法 getDataByDate
    const getDataByDate = (dateString) => {
        const nutritionGoal = state.nutritionGoals.find(g => g.date.startsWith(dateString))
        const carbonGoal = state.carbonGoals.find(g => g.date.startsWith(dateString))

        const dailyNutritionIntakes = state.nutritionIntakes.filter(i => i.date.startsWith(dateString))
        const dailyCarbonIntakes = state.carbonIntakes.filter(i => i.date.startsWith(dateString))

        const meals = { breakfast: {}, lunch: {}, dinner: {}, other: {} }

        const totalNutrients = { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 }
        let totalCarbonEmission = 0

        for (const intake of dailyNutritionIntakes) {
            let mealType = intake.meal_type || 'other'
            if (!meals[mealType].nutrients) {
                meals[mealType].nutrients = { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 }
            }

            meals[mealType].nutrients.calories += intake.calories || 0
            meals[mealType].nutrients.protein += intake.protein || 0
            meals[mealType].nutrients.fat += intake.fat || 0
            meals[mealType].nutrients.carbohydrates += intake.carbohydrates || 0
            meals[mealType].nutrients.sodium += intake.sodium || 0

            totalNutrients.calories += intake.calories || 0
            totalNutrients.protein += intake.protein || 0
            totalNutrients.fat += intake.fat || 0
            totalNutrients.carbohydrates += intake.carbohydrates || 0
            totalNutrients.sodium += intake.sodium || 0
        }

        for (const cIntake of dailyCarbonIntakes) {
            let mealType = cIntake.meal_type || 'other'
            if (!meals[mealType].carbonEmission) {
                meals[mealType].carbonEmission = 0
            }
            meals[mealType].carbonEmission += cIntake.emission || 0
            totalCarbonEmission += cIntake.emission || 0
        }

        // 如果不存在数据（包括明天的日期或任何后端未提供的数据），返回全为0的默认值
        // 这样无论后端是否提供明天的数据，前端都能获得明天为0的数据
        if (!nutritionGoal && dailyNutritionIntakes.length === 0 && dailyCarbonIntakes.length === 0) {
            return {
                nutrients: {
                    actual: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 },
                    target: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 }
                },
                carbonEmission: {
                    actual: 0,
                    target: 0
                },
                meals: {
                    breakfast: { nutrients: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 }, carbonEmission: 0 },
                    lunch: { nutrients: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 }, carbonEmission: 0 },
                    dinner: { nutrients: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 }, carbonEmission: 0 },
                    other: { nutrients: { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 }, carbonEmission: 0 }
                }
            }
        }

        return {
            nutrients: {
                actual: totalNutrients,
                target: nutritionGoal ? {
                    calories: nutritionGoal.calories,
                    protein: nutritionGoal.protein,
                    fat: nutritionGoal.fat,
                    carbohydrates: nutritionGoal.carbohydrates,
                    sodium: nutritionGoal.sodium
                } : { calories: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 }
            },
            carbonEmission: {
                actual: totalCarbonEmission,
                target: carbonGoal ? carbonGoal.emission : 0
            },
            meals
        }
    }

    return {
        state,
        setNutritionGoals,
        getNutritionGoals,
        setCarbonGoals,
        getCarbonGoals,
        getNutritionIntakes,
        getCarbonIntakes,
        setSharedNutritionCarbonIntake, // 新增方法
        getSharedNutritionCarbonIntakes, // 新增方法（假设有对应的GET接口）
        reset,
        clearStorage,
        getDataByDate // 导出辅助方法
    };
});
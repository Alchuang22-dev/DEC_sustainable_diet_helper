// stores/carbon_and_nutrition_data.js
import { defineStore } from 'pinia';
import { reactive } from 'vue';

export const useCarbonAndNutritionStore = defineStore('carbonAndNutrition', () => {
    // 生成从六天前到六天后的日期数组
    const generateDateRange = () => {
        const dates = [];
        const today = new Date();
        for (let i = -6; i <= 6; i++) {
            const date = new Date();
            date.setDate(today.getDate() + i);
            const dateString = date.toISOString().split('T')[0];
            dates.push(dateString);
        }
        return dates;
    };

    // 初始化数据
    const dateRange = generateDateRange();
    const data = reactive({});

    dateRange.forEach(date => {
        data[date] = {
            carbonEmission: { actual: 0, target: 0 },
            nutrients: {
                actual: { energy: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 },
                target: { energy: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 }
            },
            meals: {
                breakfast: {
                    carbonEmission: 0,
                    nutrients: { energy: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 }
                },
                lunch: {
                    carbonEmission: 0,
                    nutrients: { energy: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 }
                },
                dinner: {
                    carbonEmission: 0,
                    nutrients: { energy: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 }
                },
                others: {
                    carbonEmission: 0,
                    nutrients: { energy: 0, protein: 0, fat: 0, carbohydrates: 0, sodium: 0 }
                }
            }
        };
    });

    // 从本地存储加载数据
    const loadData = () => {
        const storedData = uni.getStorageSync('carbonAndNutritionData');
        if (storedData) {
            Object.assign(data, storedData);
        }
    };

    // 保存数据到本地存储
    const saveData = () => {
        uni.setStorageSync('carbonAndNutritionData', data);
    };

    // 更新每日数据
    const updateDailyData = (date, newData) => {
        if (data[date]) {
            Object.assign(data[date], newData);
            saveData();
        }
    };

    // 更新餐食数据
    const updateMealData = (date, mealType, newMealData) => {
        if (data[date] && data[date].meals[mealType]) {
            Object.assign(data[date].meals[mealType], newMealData);
            saveData();
        }
    };

    // 根据日期获取数据
    const getDataByDate = (date) => {
        return data[date];
    };

    // 设置目标值
    const setTargetValues = (date, targetData) => {
        if (data[date]) {
            data[date].carbonEmission.target = targetData.carbonEmission || data[date].carbonEmission.target;
            data[date].nutrients.target = targetData.nutrients || data[date].nutrients.target;
            saveData();
        }
    };

    // 初始化
    loadData();

    return {
        data,
        updateDailyData,
        updateMealData,
        getDataByDate,
        setTargetValues,
        saveData,
        loadData
    };
});

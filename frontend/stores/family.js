// /stores/family.js
import { defineStore } from 'pinia';
import { reactive } from 'vue';

export const useFamilyStore = defineStore('family', () => {
    const family = reactive({
        id: '',
        name: '',
        members: [],
        dishProposals: [],
    });

    // 创建家庭
    const createFamily = (familyName) => {
        // 模拟创建家庭ID
        family.id = 'F' + Date.now();
        family.name = familyName;
        family.members = [
            {
                id: 'U1',
                name: 'user_name',
                family_name: 'You',
                avatar: '/static/images/user_avatar.png',
            },
        ];
        family.dishProposals = [];
    };

    // 加入家庭
    const joinFamily = (familyId) => {
        // 模拟加入家庭
        family.id = familyId;
        family.name = 'Joined Family';
        family.members = [
            {
                id: 'U1',
                name: 'user_name 1',
                family_name: 'You',
                avatar: '/static/images/user_avatar.png',
            },
            {
                id: 'U2',
                name: 'user_name 2',
                family_name: 'Dad',
                avatar: '/static/images/member2_avatar.png',
            },
        ];
        family.dishProposals = [];
    };

    // 添加菜品提议
    const addDishProposal = (dish) => {
        family.dishProposals.push(dish);
    };

    return {
        family,
        createFamily,
        joinFamily,
        addDishProposal,
    };
});

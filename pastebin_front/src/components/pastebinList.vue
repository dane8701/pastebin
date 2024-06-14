<script setup>
import axios from 'axios';
import { onMounted } from 'vue';
import { ref } from "vue";
import router from "../routes/index"

const baseUrl = 'http://localhost:4000';
const pastes = ref([]);

async function getAllPaste() {
    try {
        const response = await axios.get(`${baseUrl}/bins`);
        console.log("response: ", response);
        pastes.value = response.data;
    } catch (error) {
        console.error('Error fetching data:', error);
        pastes.value = []; // 
    }
}

function goToPaste(alias) {
    router.push({ name: 'pastDetails', params: { id: alias } });
}

function goToUpdatePast(alias) {
    router.push({ name: 'pastUpdate', params: { id: alias } });
}

function goToNewPast() {
    router.push('/newPast')
}

onMounted(async () => {
    await getAllPaste();
    console.log('Pastes:', pastes.value);
});

</script>
<template>
    <div class="pasteList">
        <div class="paste" v-for="(paste, index) in pastes" :key="index">
            <div>
                {{ "Paste n°" + (index + 1) }}
                <span>{{ paste.alias }}</span>
                <span>{{ paste.contain }}</span>
            </div>
            <div class="buttons">
                <button @click="goToPaste(paste.alias)">Details</button>
                <button @click="goToUpdatePast(paste.alias)">Update</button>
            </div>
        </div>
        <button class="newpast" @click="goToNewPast">Create new past</button>
    </div>
</template>
<style lang="css">
.pasteList {
    display: grid;
    gap: 8px;
    border: 1px solid;
    padding: 28px;

    .paste {
        align-items: center;
        justify-content: space-between;
        display: flex;
        gap: 6px;
        padding: 4px;
        border: 4px solid aliceblue;
    }

    .newPast {
        margin-top: 6px;
    }
}
</style>
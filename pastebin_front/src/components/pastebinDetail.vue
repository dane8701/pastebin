<script setup>
import { onMounted, ref } from 'vue';
import axios from 'axios';
import { useRoute } from 'vue-router';
import router from "../routes/index"


const baseUrl = 'http://localhost:4000';
const paste = ref(null);

const route = useRoute();
const alias = route.params.id;

async function getPaste() {
    try {
        const response = await axios.get(`${baseUrl}/bins/${alias}`);
        console.log(response);
        paste.value = response.data;
    } catch (error) {
        console.error('Error fetching paste:', error);
        paste.value = {};
    }
}

function goToPastList() {
    router.push({ name: 'pasteList'});
}

onMounted(async () => {
    await getPaste();
});
</script>
<template>
    <div v-if="paste && paste.alias" class="paste">
        <span class="alias">Alias: {{ paste.alias }}</span>
        <span class="contain">Contain: {{ paste.contain }}</span>
    </div>
    <div v-else>
        Loading...
    </div>
    <button @click="goToPastList">Go back to the list</button>
</template>
<style>
.paste {
    display: grid;
    gap: 6px;
}
</style>
<script setup>
import { onMounted, ref } from 'vue';
import axios from 'axios';
import { useRoute } from 'vue-router';
import router from "../routes/index"


const baseUrl = 'http://localhost:4000';
const paste = ref(null);
let error = ref(null)
let img = ref("");

const route = useRoute();
const alias = route.params.id;

async function getPaste() {
    try {
        const response = await axios.get(`${baseUrl}/bins/${alias}`);
        console.log(response);
        paste.value = response.data;
    } catch (err) {
        error = err;
        console.error('Error fetching paste:', err);
        paste.value = {};
    }
}

function goToPastList() {
    router.push({ name: 'pasteList'});
}

onMounted(async () => {
    await getPaste();
    img.value = "http://localhost:4000/bins/file/"  + paste.value.alias;
});
</script>
<template>
    <div v-if="paste && paste.alias" class="pasteDetails">
        <span class="alias">Alias: {{ paste.alias }}</span>
        <span class="contain">Contain: {{ paste.contain }}</span>
        <span class="clic">Click: {{ paste.clic }}</span>
        <img :src="img"/>
    </div>
    <div v-if="error">
        Error !!!!
    </div>
    <button @click="goToPastList">Go back to the list</button>
</template>
<style>
.pasteDetails {
    display: grid;
    gap: 6px;
    margin: auto;

    img {
        width: 200px;
    }
}
</style>